package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

const (
	defaultBaseURL = "https://hub.docker.com"
)

var closeFunc = func(body io.ReadCloser) {
	if body != nil {
		if errClose := body.Close(); errClose != nil {
			fmt.Print(errClose)
		}
	}
}

type Client struct {
	client      *http.Client
	dryRun      bool
	publisherID string
	baseURL     string
	bearer      string
}

func NewClient(username string, password string, publisherID string, dryRun bool, opts ...ClientOption) (*Client, error) {
	c := &Client{
		client:      http.DefaultClient,
		baseURL:     defaultBaseURL,
		publisherID: publisherID,
		dryRun:      dryRun,
	}

	for _, opt := range opts {
		opt(c)
	}

	if err := c.doLogin(context.Background(), username, password); err != nil {
		return nil, err
	}

	return c, nil
}

// ClientOption allows to configure a Client.
// It is used in NewClient.
type ClientOption func(c *Client)

func (c *Client) doLogin(ctx context.Context, username string, password string) error {
	type login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	loginJSON, err := json.Marshal(login{Username: username, Password: password})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v2/users/login/", c.baseURL), bytes.NewReader(loginJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}

	defer func() {
		if errClose := resp.Body.Close(); errClose != nil {
			fmt.Print(errClose)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		var b []byte
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("docker: non-200 response code: could not read response body: %v", err)
		}

		return fmt.Errorf("docker: non-200 response code: %v", string(b))
	}

	type token struct {
		Token string `json:"token"`
	}
	var tkn token
	if err = json.NewDecoder(resp.Body).Decode(&tkn); err != nil {
		return err
	}
	c.bearer = tkn.Token

	fmt.Println("Successfully login")
	return nil
}

func (c *Client) doReq(ctx context.Context, method, url string, body io.Reader) (io.ReadCloser, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.bearer)

	if c.dryRun && (method == http.MethodPost || method == http.MethodPut) {
		dumpRequest, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dumpRequest))
		return nil, nil
	}
	resp, err := c.client.Do(req.WithContext(ctx))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("docker: non-200 response code: could not read response body: %v", err)
		}

		return nil, fmt.Errorf("docker: non-200 response code: %v", string(b))
	}

	return resp.Body, nil
}

func (c *Client) ListProduct(ctx context.Context) ([]Product, error) {
	body, err := c.doReq(ctx, http.MethodGet, fmt.Sprintf("%s/api/publish/v1/publishers/%s/products", c.baseURL, c.publisherID), nil)
	defer closeFunc(body)
	if err != nil {
		return nil, err
	}

	var products []Product
	if err = json.NewDecoder(body).Decode(&products); err != nil {
		return nil, err
	}

	return products, nil
}

func (c *Client) ListPlans(ctx context.Context, productID string) ([]Plan, error) {
	body, err := c.doReq(ctx, http.MethodGet, fmt.Sprintf("%s/api/publish/v1/publishers/%s/products/%s/rate-plans", c.baseURL, c.publisherID, productID), nil)
	defer closeFunc(body)
	if err != nil {
		return nil, err
	}

	var plans []Plan
	if err = json.NewDecoder(body).Decode(&plans); err != nil {
		return nil, err
	}

	return plans, nil
}

func (c *Client) ListImagesTags(ctx context.Context) ([]Tag, error) {
	url := fmt.Sprintf("%s/v2/repositories/containous/traefikee/tags", c.baseURL)
	var tags []Tag

	for url != "" {
		body, err := c.doReq(ctx, http.MethodGet, url, nil)
		defer closeFunc(body)
		if err != nil {
			return nil, err
		}

		var imageTags ImagesTags
		if err = json.NewDecoder(body).Decode(&imageTags); err != nil {
			return nil, err
		}
		tags = append(tags, imageTags.Results...)
		url = imageTags.Next
	}

	return tags, nil
}

func (c *Client) GetPlan(plans []Plan, planID string) (Plan, error) {
	for _, p := range plans {
		if p.ID == planID {
			return p, nil
		}
	}

	return Plan{}, fmt.Errorf("unable to find plan with id: %q", planID)
}

func (c *Client) SavePlan(ctx context.Context, repo string, project string, productID string, plan Plan, selectedTags []string) error {
	data := make(map[string]Plan)

	for _, value := range selectedTags {
		plan.Repositories = append(plan.Repositories, Repository{
			PublishersRepoName: Repo{
				Namespace: repo,
				RepoName:  project,
				Tag:       value,
			}})
	}
	data[plan.ID] = plan

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	body, err := c.doReq(ctx, http.MethodPut, fmt.Sprintf("%s/api/publish/v1/publishers/%s/products/%s/rate-plans", c.baseURL, c.publisherID, productID), bytes.NewReader(dataJSON))
	defer closeFunc(body)

	return err
}

func (c *Client) SubmitForReview(ctx context.Context, productID string, planID string) error {
	body, err := c.doReq(ctx, http.MethodPost, fmt.Sprintf("%s/api/publish/v1/publishers/%s/products/%s/publish-steps", c.baseURL, c.publisherID, productID), bytes.NewReader([]byte(`{"publisher_action":"submitted_for_review"}`)))
	defer closeFunc(body)

	return err
}

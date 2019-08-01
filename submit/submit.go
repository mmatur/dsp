package submit

import (
	"context"
	"regexp"

	"github.com/mmatur/dsp/choose"
	"github.com/mmatur/dsp/docker"
	"github.com/mmatur/dsp/types"
)

// semVerRegex is the regular expression used to parse a semantic version.
const semVerRegex string = `v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?`

// Submit submits a new images to docker.
func Submit(config *types.Config) error {
	ctx := context.Background()
	client, err := docker.NewClient(config.Username, config.Password, config.PublisherID, config.DryRun)
	if err != nil {
		return err
	}

	products, err := client.ListProduct(ctx)
	if err != nil {
		return err
	}

	productID, err := choose.Product(products)
	if err != nil || productID == choose.ExitValue {
		return err
	}

	plans, err := client.ListPlans(ctx, productID)
	if err != nil {
		return err
	}

	planID, err := choose.Plan(plans)
	if err != nil || planID == choose.ExitValue {
		return err
	}

	plan, err := client.GetPlan(plans, planID)
	if err != nil {
		return err
	}

	tags, err := client.ListImagesTags(ctx)
	if err != nil {
		return err
	}

	selectedTags, err := choose.Tags(createTgaList(plan.Repositories, tags))
	if err != nil || len(tags) == 0 {
		return err
	}

	if err := client.SavePlan(ctx, config.Repository, config.Project, productID, plan, selectedTags); err != nil {
		return err
	}

	return client.SubmitForReview(ctx, productID, planID)
}

func createTgaList(repositories []docker.Repository, tags []docker.Tag) choose.LstTags {
	var chooseTags choose.LstTags

	for _, t := range tags {
		versionRegex := regexp.MustCompile("^" + semVerRegex + "$")
		if m := versionRegex.FindStringSubmatch(t.Name); m != nil {
			chooseTags = append(chooseTags, &choose.Tag{
				Name:      t.Name,
				CreatedAt: t.LastUpdated,
				Select:    false,
			})
		}
	}

	for _, r := range repositories {
		for _, c := range chooseTags {
			if c.Name == r.PublishersRepoName.Tag {
				c.Select = true
				break
			}
		}
	}

	return chooseTags
}

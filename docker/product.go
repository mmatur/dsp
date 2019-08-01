package docker

import "time"

type Product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Plan struct {
	ID                   string       `json:"id"`
	Name                 string       `json:"name"`
	Price                Price        `json:"price"`
	CertificationStatus  string       `json:"certification_status"`
	Description          string       `json:"description"`
	DownloadAttribute    string       `json:"download_attribute"`
	EUSA                 string       `json:"eusa"`
	EUSAType             string       `json:"eusa_type"`
	Instructions         string       `json:"instructions"`
	IsDefault            bool         `json:"is_default"`
	IsOffline            bool         `json:"is_offline"`
	OfficialStatus       string       `json:"official_status"`
	PlanType             string       `json:"plan_type"`
	Rank                 int          `json:"rank"`
	ReleaseNotes         string       `json:"release_notes"`
	Repositories         []Repository `json:"repositories"`
	CreatedAt            time.Time    `json:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at"`
	RequestCertification bool         `json:"request_certification"`
}

type Price struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	Label             string             `json:"label"`
	Duration          int                `json:"duration"`
	DurationPeriod    string             `json:"duration_period"`
	Trial             int                `json:"trial"`
	TrialPeriod       string             `json:"trial_period"`
	Expires           bool               `json:"expires"`
	PricingComponents []PricingComponent `json:"pricing_components"`
}

type PricingComponent struct {
	Name            string `json:"name"`
	Label           string `json:"label"`
	UnitOfMeasure   string `json:"unit_of_measure"`
	ChargeType      string `json:"charge_type"`
	ChargeModel     string `json:"charge_model"`
	DefaultQuantity int    `json:"default_quantity"`
	Tiers           []Tier `json:"tiers"`
}

type Tier struct {
	LowerThreshold int    `json:"lower_threshold"`
	UpperThreshold int    `json:"upper_threshold"`
	PricingType    string `json:"pricing_type"`
	Price          int    `json:"price"`
}

type Repository struct {
	PublishersRepoName Repo `json:"publishers_repo_name"`
}

type Repo struct {
	Namespace string `json:"namespace"`
	RepoName  string `json:"reponame"`
	Tag       string `json:"tag"`
}

type ImagesTags struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []Tag  `json:"results"`
}

type Tag struct {
	Name        string    `json:"name"`
	LastUpdated time.Time `json:"last_updated"`
}

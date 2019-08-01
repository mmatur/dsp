package choose

import (
	"fmt"
	"sort"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
)

type Tag struct {
	Name      string
	CreatedAt time.Time
	Select    bool
}

type LstTags []*Tag

func (t LstTags) Len() int           { return len(t) }
func (t LstTags) Less(i, j int) bool { return t[i].CreatedAt.After(t[j].CreatedAt) }
func (t LstTags) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

type answersTags struct {
	Tags []string
}

// Tags choose a tag in the list.
func Tags(tags LstTags) ([]string, error) {
	sort.Sort(tags)
	var surveyOpts []string
	var defaultOpts []int
	for i, value := range tags {
		name := fmt.Sprintf("%s - %s", value.Name, value.CreatedAt.Format("2006-01-02"))
		surveyOpts = append(surveyOpts, name)
		if value.Select {
			defaultOpts = append(defaultOpts, i)
		}
	}

	var qs = []*survey.Question{
		{
			Name: "tags",
			Prompt: &survey.MultiSelect{
				Message:  "Choose tags",
				Options:  surveyOpts,
				PageSize: 10,
				Default:  defaultOpts,
			},
		},
	}

	answers := &answersTags{}
	err := survey.Ask(qs, answers)
	if err != nil {
		return nil, err
	}

	return answers.Tags, nil
}

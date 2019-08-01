package choose

import (
	"fmt"
	"strings"

	"github.com/mmatur/dsp/docker"

	survey "github.com/AlecAivazis/survey/v2"
)

type answersPlan struct {
	Name string
}

func (a answersPlan) isExit() bool {
	return a.Name == ExitLabel
}

func (a answersPlan) getPlanID() (string, error) {
	parts := strings.SplitN(a.Name, ":", 2)

	if len(parts) != 2 {
		return ExitValue, fmt.Errorf("unable to get the product ID")
	}

	return strings.TrimSpace(parts[1]), nil
}

// Plan choose a plan in the list.
func Plan(plans []docker.Plan) (string, error) {
	var surveyOpts []string
	for _, p := range plans {
		surveyOpts = append(surveyOpts, fmt.Sprintf("%s: %s", p.Name, p.ID))
	}

	var qs = []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Select{
				Message: "Choose a plan",
				Options: append(surveyOpts, ExitLabel),
			},
		},
	}

	answers := &answersPlan{}
	err := survey.Ask(qs, answers)
	if err != nil {
		return "", err
	}

	if answers.isExit() {
		return ExitValue, nil
	}

	return answers.getPlanID()
}

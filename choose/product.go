package choose

import (
	"fmt"
	"strings"

	"github.com/mmatur/dsp/docker"

	survey "github.com/AlecAivazis/survey/v2"
)

type answersProduct struct {
	Name string
}

func (a answersProduct) isExit() bool {
	return a.Name == ExitLabel
}

func (a answersProduct) getProductID() (string, error) {
	parts := strings.SplitN(a.Name, ":", 2)

	if len(parts) != 2 {
		return ExitValue, fmt.Errorf("unable to get the product ID")
	}

	return strings.TrimSpace(parts[1]), nil
}

// Product choose a product in the list.
func Product(products []docker.Product) (string, error) {
	var surveyOpts []string
	for _, product := range products {
		surveyOpts = append(surveyOpts, fmt.Sprintf("%s: %s", product.Name, product.ID))
	}

	var qs = []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Select{
				Message: "Choose a product",
				Options: append(surveyOpts, ExitLabel),
			},
		},
	}

	answers := &answersProduct{}
	err := survey.Ask(qs, answers)
	if err != nil {
		return "", err
	}

	if answers.isExit() {
		return ExitValue, nil
	}

	return answers.getProductID()
}

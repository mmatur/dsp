package choose

import (
	survey "github.com/AlecAivazis/survey/v2"
)

// Action name
const (
	ActionSubmit = "submit"
)

type answersAction struct {
	Action string
}

func (a answersAction) isExit() bool {
	return a.Action == ExitLabel
}

// Action get PRM action
func Action() (string, error) {
	surveyOpts := []string{ActionSubmit, ExitLabel}

	var qs = []*survey.Question{
		{
			Name: "action",
			Prompt: &survey.Select{
				Message: "Choose the action",
				Options: surveyOpts,
			},
		},
	}

	answers := &answersAction{}
	err := survey.Ask(qs, answers)
	if err != nil {
		return "", err
	}

	if answers.isExit() {
		return ExitLabel, nil
	}

	return answers.Action, nil
}

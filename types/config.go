package types

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	DSPEnvFilePathEnvVar = "DSP_ENV_FILE_PATH"
	UsernameEnvVar       = "DOCKER_USER"
	PasswordEnvVar       = "DOCKER_PASSWORD"
	PublisherIDEnvVar    = "DOCKER_PUBLISHER_ID"
	RepositoryEnvVar     = "DOCKER_REPOSITORY"
	ProjectEnvVar        = "DOCKER_PROJECT"
)

type Config struct {
	Username    string `flag:"username"`
	Password    string `flag:"password"`
	PublisherID string `flag:"publisher-id"`
	Repository  string `flag:"repository"`
	Project     string `flag:"project"`
	DryRun      bool   `flag:"dry-run"`
}

func ValidateRequiredFlags(cfg Config) error {
	var missingFlagNames []string

	value := reflect.ValueOf(cfg)
	for i := 0; i < value.NumField(); i++ {
		if value.Field(i).String() != "" {
			continue
		}

		fieldType := value.Type().Field(i)
		name := fieldType.Tag.Get("flag")
		if name == "" {
			name = fieldType.Name
		}

		missingFlagNames = append(missingFlagNames, name)
	}

	if len(missingFlagNames) > 0 {
		return fmt.Errorf(`required flag(s) "%s" not set`, strings.Join(missingFlagNames, `", "`))
	}

	return nil
}

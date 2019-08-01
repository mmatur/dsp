package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mmatur/dsp/choose"
	"github.com/mmatur/dsp/cmd"
	"github.com/mmatur/dsp/submit"
	"github.com/mmatur/dsp/types"
	"github.com/spf13/cobra"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg := &types.Config{}
	// rootCmd represents the base command when called without any subcommands.
	rootCmd := &cobra.Command{
		Use:   "dsp",
		Short: "DSP - Docker store publisher",
		Long:  "DSP - Docker store publisher",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return types.ValidateRequiredFlags(*cfg)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Println("Config used:")
			fmt.Printf("    Repository:   %s\n", cfg.Repository)
			fmt.Printf("    Project:      %s\n", cfg.Project)
			fmt.Printf("    Username:     %s\n", cfg.Username)
			fmt.Printf("    Publisher ID: %s\n", cfg.PublisherID)
			fmt.Printf("    Dry Run:      %t\n", cfg.DryRun)
			return rootRun(cfg)
		},
	}

	_ = os.Unsetenv(types.UsernameEnvVar)
	_ = os.Unsetenv(types.PasswordEnvVar)

	filePath := os.Getenv(types.DSPEnvFilePathEnvVar)

	if filePath != "" {
		if err := godotenv.Load(filePath); err != nil {
			log.Printf("Unable to load %s: %v\n", filePath, err)
		}
	}

	rootCmd.AddCommand(cmd.NewVersion())
	rootCmd.AddCommand(cmd.NewDoc())
	rootCmd.AddCommand(cmd.NewSubmit(cfg))

	rootFlags := rootCmd.PersistentFlags()
	rootFlags.StringVar(&cfg.Username, "username", os.Getenv(types.UsernameEnvVar), "Docker username.")
	rootFlags.StringVar(&cfg.Password, "password", os.Getenv(types.PasswordEnvVar), "Docker password.")
	rootFlags.StringVar(&cfg.PublisherID, "publisher-id", os.Getenv(types.PublisherIDEnvVar), "Docker publisher ID.")
	rootFlags.StringVar(&cfg.Repository, "repository", os.Getenv(types.RepositoryEnvVar), "Docker hub repository.")
	rootFlags.StringVar(&cfg.Project, "project", os.Getenv(types.ProjectEnvVar), "Docker hub project.")
	rootFlags.BoolVar(&cfg.DryRun, "dry-run", true, "Dry run mode.")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func rootRun(cfg *types.Config) error {
	action, err := choose.Action()
	if err != nil {
		return err
	}

	if action == choose.ActionSubmit {
		return submit.Submit(cfg)
	}

	return nil
}

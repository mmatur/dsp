package cmd

import (
	"fmt"

	"github.com/mmatur/dsp/types"

	"github.com/mmatur/dsp/submit"
	"github.com/spf13/cobra"
)

func NewSubmit(cfg *types.Config) *cobra.Command {
	submitCmd := &cobra.Command{
		Use:   "submit",
		Short: "Submit available image.",
		Long:  "Submit available image.",
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
			return submit.Submit(cfg)
		},
	}

	submitFlags := submitCmd.Root().Flags()

	submitFlags.BoolVar(&cfg.DryRun, "dry-run", true, "Dry run")

	return submitCmd
}

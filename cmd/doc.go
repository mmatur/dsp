package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func NewDoc() *cobra.Command {
	return &cobra.Command{
		Use:    "doc",
		Short:  "Generate documentation",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := os.MkdirAll("./docs", 0755); err != nil {
				return err
			}

			return doc.GenMarkdownTree(cmd.Parent(), "./docs")
		},
	}
}

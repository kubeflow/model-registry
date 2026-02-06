package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	var (
		entityName  string
		packageName string
		outputDir   string
	)

	cmd := &cobra.Command{
		Use:   "init <name>",
		Short: "Initialize a new catalog plugin",
		Long: `Initialize a new catalog plugin for the unified catalog server.

This creates the basic directory structure and configuration file
for a new catalog plugin.

Example:
  catalog-gen init mcp-catalog --entity=MCPServer --package=github.com/kubeflow/model-registry/catalog/plugins/mcp`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			if entityName == "" {
				return fmt.Errorf("--entity is required")
			}
			if packageName == "" {
				return fmt.Errorf("--package is required")
			}

			if outputDir == "" {
				outputDir = name
			}

			return initCatalogPlugin(name, entityName, packageName, outputDir)
		},
	}

	cmd.Flags().StringVar(&entityName, "entity", "", "Name of the main entity (e.g., MCPServer)")
	cmd.Flags().StringVar(&packageName, "package", "", "Go package path (e.g., github.com/kubeflow/model-registry/catalog/plugins/mcp)")
	cmd.Flags().StringVar(&outputDir, "output", "", "Output directory (defaults to catalog name)")

	return cmd
}

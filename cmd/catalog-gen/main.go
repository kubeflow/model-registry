// catalog-gen is a scaffolding tool for creating new catalog components.
// It generates the boilerplate code needed for a new catalog service,
// similar to how kubebuilder scaffolds Kubernetes controllers.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:   "catalog-gen",
		Short: "Scaffolding tool for creating catalog components",
		Long: `catalog-gen is a scaffolding tool for creating new catalog components
in the Model Registry project. It generates the boilerplate code needed
for a new catalog service, including:

- Entity models and repositories
- REST API handlers (OpenAPI-based)
- Data providers (YAML, HTTP)
- Kustomize manifests for deployment

Usage:
  catalog-gen init <name> --entity=<EntityName> --package=<go-package>
  catalog-gen add-provider <type>
  catalog-gen add-artifact <ArtifactName>
  catalog-gen generate`,
		Version: version,
	}

	rootCmd.AddCommand(newInitCmd())
	rootCmd.AddCommand(newAddProviderCmd())
	rootCmd.AddCommand(newAddArtifactCmd())
	rootCmd.AddCommand(newGenerateCmd())
	rootCmd.AddCommand(newGenTestdataCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

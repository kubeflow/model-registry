package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func newGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate or regenerate code from catalog.yaml",
		Long: `Generate or regenerate the catalog plugin code based on the catalog.yaml configuration.

This command reads the catalog.yaml file and generates:
  - Entity models and repositories
  - Artifact models and repositories
  - Provider implementations
  - OpenAPI specification
  - Plugin registration files

Example:
  catalog-gen generate`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return generate()
		},
	}

	return cmd
}

func generate() error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	fmt.Printf("Regenerating code for catalog plugin: %s\n", config.Metadata.Name)
	fmt.Println("Note: Only non-editable files are regenerated. Editable files are created by 'catalog-gen init'.")
	fmt.Println()

	// === PLUGIN FILES ===

	if err := generatePluginFiles(config); err != nil {
		return fmt.Errorf("failed to generate plugin files: %w", err)
	}

	// === ALWAYS REGENERATED (not meant to be edited) ===

	// Generate entity model
	if err := generateEntityModel(config); err != nil {
		return fmt.Errorf("failed to generate entity model: %w", err)
	}

	// Generate artifact models (if artifacts are configured)
	for _, artifact := range config.Spec.Artifacts {
		if err := generateArtifactModel(config, artifact); err != nil {
			return fmt.Errorf("failed to generate artifact model for %s: %w", artifact.Name, err)
		}
		if err := generateArtifactRepository(config, artifact); err != nil {
			return fmt.Errorf("failed to generate artifact repository for %s: %w", artifact.Name, err)
		}
	}

	// Generate datastore spec
	if err := generateDatastoreSpec(config); err != nil {
		return fmt.Errorf("failed to generate datastore spec: %w", err)
	}

	// Generate OpenAPI components
	if err := generateOpenAPIComponents(config); err != nil {
		return fmt.Errorf("failed to generate OpenAPI components: %w", err)
	}

	// Generate OpenAPI main file
	if err := generateOpenAPIMain(config); err != nil {
		return fmt.Errorf("failed to generate OpenAPI main: %w", err)
	}
	fmt.Printf("  Generated: api/openapi/src/openapi.yaml\n")

	// Generate loader
	if err := generateLoader(config); err != nil {
		return fmt.Errorf("failed to generate loader: %w", err)
	}

	// Regenerate YAML provider if it doesn't exist
	yamlProviderPath := filepath.Join("internal", "catalog", "providers", "yaml_provider.go")
	if _, err := os.Stat(yamlProviderPath); os.IsNotExist(err) {
		for _, provider := range config.Spec.Providers {
			if provider.Type == "yaml" {
				if err := generateYAMLProvider(config); err != nil {
					return fmt.Errorf("failed to generate YAML provider: %w", err)
				}
				break
			}
		}
	}

	// Generate post-generation documentation (if artifacts are configured)
	if len(config.Spec.Artifacts) > 0 {
		if err := generatePostGenerationDocs(config); err != nil {
			return fmt.Errorf("failed to generate post-generation docs: %w", err)
		}
	}

	fmt.Println("\nGeneration complete!")
	fmt.Println("\nIf you added new properties or artifacts to catalog.yaml, you may need to manually update:")
	fmt.Println("  - internal/db/service/<entity>.go (property mapping in converters)")
	fmt.Println("  - internal/server/openapi/api_*_service_impl.go (OpenAPI conversion)")
	fmt.Println("  - internal/catalog/providers/*.go (provider parsing)")
	fmt.Println("\nThen run 'make gen/openapi-server' to regenerate OpenAPI handlers.")

	return nil
}

// =============================================================================
// Post-Generation Documentation
// =============================================================================

func generatePostGenerationDocs(config CatalogConfig) error {
	docsDir := "docs"
	if err := ensureDir(docsDir); err != nil {
		return fmt.Errorf("failed to create docs directory: %w", err)
	}

	entityName := config.Spec.Entity.Name
	lowerEntity := strings.ToLower(entityName)
	catalogName := config.Metadata.Name

	// Build artifact-specific sections
	var artifactRepoLines strings.Builder
	var artifactMethodsDoc strings.Builder
	var artifactConversionsDoc strings.Builder
	var artifactChecklist strings.Builder
	var artifactYAMLFields strings.Builder

	for _, artifact := range config.Spec.Artifacts {
		artifactName := artifact.Name
		lowerArtifact := strings.ToLower(artifactName)
		fullArtifactName := entityName + artifactName + "Artifact"
		repoName := fullArtifactName + "Repository"

		artifactRepoLines.WriteString(fmt.Sprintf("    getRepo[models.%s](repoSet),\n", repoName))

		artifactMethodsDoc.WriteString(fmt.Sprintf(`
### Get%s%ss Endpoint

Add this method to `+"`internal/server/openapi/api_%s_service_impl.go`:\n\n", entityName, artifactName, lowerEntity))
		artifactMethodsDoc.WriteString("```go\n")
		artifactMethodsDoc.WriteString(fmt.Sprintf(`// Get%s%ss implements DefaultAPIServicer.Get%s%ss
func (s *%sCatalogServiceAPIService) Get%s%ss(
    ctx context.Context,
    name string,
    pageSize int32,
    pageToken string,
) (ImplResponse, error) {
    // Implementation here
}
`, entityName, artifactName, entityName, artifactName, entityName, entityName, artifactName))
		artifactMethodsDoc.WriteString("```\n")

		artifactConversionsDoc.WriteString(fmt.Sprintf("\n### convert%sToOpenAPIModel\n\n", fullArtifactName))
		artifactConversionsDoc.WriteString("```go\n")
		artifactConversionsDoc.WriteString(fmt.Sprintf(`func convert%sToOpenAPIModel(artifact models.%s) %s {
    // Implementation here
}
`, fullArtifactName, fullArtifactName, fullArtifactName))
		artifactConversionsDoc.WriteString("```\n")

		artifactChecklist.WriteString(fmt.Sprintf("| `internal/server/openapi/api_%s_service_impl.go` | Add Get%s%ss method | ☐ Manual |\n", lowerEntity, entityName, artifactName))
		artifactChecklist.WriteString(fmt.Sprintf("| `internal/server/openapi/api_%s_service_impl.go` | Add convert%sToOpenAPIModel | ☐ Manual |\n", lowerEntity, fullArtifactName))

		artifactYAMLFields.WriteString(fmt.Sprintf("  - %sName: \"example-%s\"\n", lowerEntity, lowerEntity))
		artifactYAMLFields.WriteString(fmt.Sprintf("    name: \"example-%s\"\n", lowerArtifact))
		for _, prop := range artifact.Properties {
			artifactYAMLFields.WriteString(fmt.Sprintf("    %s: \"example-value\"\n", prop.Name))
		}
	}

	content := fmt.Sprintf(`# Post-Artifact Generation Manual Steps

> **Auto-generated by catalog-gen** - Regenerate with `+"`catalog-gen generate`"+`

After adding an artifact to `+"`catalog.yaml`"+` and running `+"`catalog-gen generate`"+`, complete these manual steps.

## Manual Steps Required

### 1. Update plugin.go

Add the artifact repository to initServices:

`+"```go"+`
services := service.NewServices(
    getRepo[models.%[3]sRepository](repoSet),
%[4]s)
`+"```"+`

### 2. Regenerate OpenAPI Server Code

`+"```bash"+`
make gen/openapi-server
`+"```"+`

### 3. Implement the Artifact List Endpoint(s)

%[5]s

### 4. Add the Artifact Conversion Function(s)

%[6]s

## File Checklist

| File | Action | Status |
|------|--------|--------|
| `+"`plugin.go`"+` | Add artifact repository to initServices | ☐ Manual |
%[9]s
`, lowerEntity, catalogName, entityName, artifactRepoLines.String(),
		artifactMethodsDoc.String(), artifactConversionsDoc.String(),
		artifactYAMLFields.String(), config.Spec.API.Port,
		artifactChecklist.String())

	docsPath := filepath.Join(docsDir, "post-artifact-generation.md")
	if err := os.WriteFile(docsPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write docs file: %w", err)
	}

	fmt.Printf("  Generated: %s\n", docsPath)
	return nil
}

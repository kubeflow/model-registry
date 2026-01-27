package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func newAddProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-provider <type>",
		Short: "Add a new provider type to the catalog",
		Long: `Add a new provider type to the catalog configuration.

Supported provider types:
  - yaml: File-based provider using YAML catalog files
  - http: HTTP-based provider for remote APIs

Example:
  catalog-gen add-provider yaml
  catalog-gen add-provider http`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			providerType := args[0]

			switch providerType {
			case "yaml", "http":
				// OK
			default:
				return fmt.Errorf("unsupported provider type: %s (supported: yaml, http)", providerType)
			}

			return addProvider(providerType)
		},
	}

	return cmd
}

func addProvider(providerType string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	// Check if provider already exists
	for _, p := range config.Spec.Providers {
		if p.Type == providerType {
			fmt.Printf("Provider '%s' already exists in catalog.yaml\n", providerType)
			return nil
		}
	}

	// Add provider
	config.Spec.Providers = append(config.Spec.Providers, ProviderConfig{Type: providerType})

	// Save config
	if err := saveConfig(config); err != nil {
		return err
	}

	// Generate provider file
	if err := generateProviderFile(config.Spec.Entity.Name, providerType); err != nil {
		return err
	}

	fmt.Printf("Added provider '%s' to catalog.yaml\n", providerType)
	fmt.Printf("Generated provider file at internal/catalog/providers/%s_provider.go\n", providerType)

	return nil
}

func newAddArtifactCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-artifact <name>",
		Short: "Add a new artifact type to the catalog",
		Long: `Add a new artifact type to the catalog configuration.

Example:
  catalog-gen add-artifact Tool
  catalog-gen add-artifact Resource`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			artifactName := args[0]
			return addArtifact(artifactName)
		},
	}

	return cmd
}

func addArtifact(artifactName string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	// Check if artifact already exists
	for _, a := range config.Spec.Artifacts {
		if a.Name == artifactName {
			fmt.Printf("Artifact '%s' already exists in catalog.yaml\n", artifactName)
			return nil
		}
	}

	// Add artifact
	config.Spec.Artifacts = append(config.Spec.Artifacts, ArtifactConfig{
		Name: artifactName,
		Properties: []PropertyConfig{
			{Name: "uri", Type: "string"},
		},
	})

	// Save config
	if err := saveConfig(config); err != nil {
		return err
	}

	// Generate artifact model stub file
	if err := generateArtifactModelStub(artifactName); err != nil {
		return err
	}

	fmt.Printf("Added artifact '%s' to catalog.yaml\n", artifactName)
	fmt.Printf("Generated artifact model at internal/db/models/%s.go\n", strings.ToLower(artifactName))
	fmt.Println("\nRun 'catalog-gen generate' to generate the full artifact implementation.")

	return nil
}

func generateArtifactModelStub(artifactName string) error {
	lowerName := strings.ToLower(artifactName)
	data := map[string]any{
		"EntityName":        "",
		"ArtifactName":      artifactName,
		"LowerEntityName":   "",
		"LowerArtifactName": lowerName,
		"Properties":        "\tURI                      *string\n",
	}

	modelsDir := filepath.Join("internal", "db", "models")
	if err := ensureDir(modelsDir); err != nil {
		return fmt.Errorf("failed to create models directory: %w", err)
	}

	return executeTemplate(TmplModelsArtifact, filepath.Join(modelsDir, fmt.Sprintf("%s.go", lowerName)), data)
}

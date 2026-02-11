package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func newGenTestdataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen-testdata",
		Short: "Generate testdata files for testing the catalog plugin",
		Long: `Generate testdata files for testing the catalog plugin locally.

This creates:
  - testdata/test-sources.yaml (catalog server config)
  - testdata/<entity>-sources.yaml (loader config)
  - testdata/<entity>s.yaml (sample entity data)

Example:
  catalog-gen gen-testdata`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateTestdata()
		},
	}

	return cmd
}

func generateTestdata() error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	entityName := config.Spec.Entity.Name
	entityNameLower := strings.ToLower(entityName)
	catalogName := filepath.Base(config.Metadata.Name)

	// Create testdata directory
	testdataDir := "testdata"
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		return fmt.Errorf("failed to create testdata directory: %w", err)
	}

	// Generate sample entity data file
	entityDataPath := filepath.Join(testdataDir, fmt.Sprintf("%ss.yaml", entityNameLower))
	if err := generateSampleEntityData(config, entityDataPath); err != nil {
		return fmt.Errorf("failed to generate sample entity data: %w", err)
	}
	fmt.Printf("  Generated: %s\n", entityDataPath)

	// Generate loader sources config
	loaderSourcesPath := filepath.Join(testdataDir, fmt.Sprintf("%s-sources.yaml", entityNameLower))
	if err := generateLoaderSourcesConfig(config, loaderSourcesPath, entityDataPath); err != nil {
		return fmt.Errorf("failed to generate loader sources config: %w", err)
	}
	fmt.Printf("  Generated: %s\n", loaderSourcesPath)

	// Generate catalog server config
	serverConfigPath := filepath.Join(testdataDir, "test-sources.yaml")
	if err := generateServerSourcesConfig(config, serverConfigPath, loaderSourcesPath, entityDataPath); err != nil {
		return fmt.Errorf("failed to generate server sources config: %w", err)
	}
	fmt.Printf("  Generated: %s\n", serverConfigPath)

	fmt.Println("\nTestdata generated successfully!")
	fmt.Println("\nTo test the plugin, run:")
	fmt.Printf("  ./catalog-server --sources=%s --listen=:8080 --db-dsn=\"host=localhost port=5432 user=postgres password=postgres dbname=model_registry sslmode=disable\"\n", serverConfigPath)
	fmt.Println("\nThen test with:")
	fmt.Printf("  curl -s http://localhost:8080/api/%s/v1alpha1/%ss | jq\n", catalogName, entityNameLower)

	return nil
}

func generateSampleEntityData(config CatalogConfig, outputPath string) error {
	entityName := config.Spec.Entity.Name
	entityNameLower := strings.ToLower(entityName)

	var content strings.Builder
	content.WriteString(fmt.Sprintf("%ss:\n", entityNameLower))

	// Track which shared fields the user already declared as custom properties
	declaredProps := make(map[string]bool)
	for _, prop := range config.Spec.Entity.Properties {
		declaredProps[strings.ToLower(prop.Name)] = true
	}

	// Generate 3 sample entities
	for i := 1; i <= 3; i++ {
		// BaseResource shared fields
		content.WriteString(fmt.Sprintf("  - name: \"sample-%s-%d\"\n", entityNameLower, i))
		content.WriteString(fmt.Sprintf("    externalId: \"ext-%s-%d\"\n", entityNameLower, i))
		content.WriteString(fmt.Sprintf("    description: \"Sample %s %d description\"\n", entityName, i))

		// Add custom properties from catalog.yaml
		for _, prop := range config.Spec.Entity.Properties {
			lower := strings.ToLower(prop.Name)
			// Skip fields already emitted as BaseResource shared fields
			if lower == "name" || lower == "externalid" || lower == "description" {
				continue
			}

			switch prop.Type {
			case "string":
				content.WriteString(fmt.Sprintf("    %s: \"sample-%s-value-%d\"\n", prop.Name, prop.Name, i))
			case "integer", "int":
				content.WriteString(fmt.Sprintf("    %s: %d\n", prop.Name, i*10))
			case "int64":
				content.WriteString(fmt.Sprintf("    %s: %d\n", prop.Name, i*1000))
			case "boolean", "bool":
				content.WriteString(fmt.Sprintf("    %s: %t\n", prop.Name, i%2 == 0))
			case "number", "float", "double":
				content.WriteString(fmt.Sprintf("    %s: %d.%d\n", prop.Name, i, i))
			case "array":
				content.WriteString(fmt.Sprintf("    %s:\n", prop.Name))
				content.WriteString(fmt.Sprintf("      - \"item-%d-a\"\n", i))
				content.WriteString(fmt.Sprintf("      - \"item-%d-b\"\n", i))
			}
		}
	}

	return os.WriteFile(outputPath, []byte(content.String()), 0644)
}

func generateLoaderSourcesConfig(_ CatalogConfig, outputPath, entityDataPath string) error {
	entityDataFile := filepath.Base(entityDataPath)

	content := fmt.Sprintf(`catalogs:
  - id: "test-source"
    name: "Test Data Source"
    type: "yaml"
    properties:
      yamlCatalogPath: "./%s"
`, entityDataFile)

	return os.WriteFile(outputPath, []byte(content), 0644)
}

func generateServerSourcesConfig(config CatalogConfig, outputPath, loaderSourcesPath, entityDataPath string) error {
	catalogName := filepath.Base(config.Metadata.Name)
	loaderSourcesFile := filepath.Base(loaderSourcesPath)
	entityDataFile := filepath.Base(entityDataPath)

	content := fmt.Sprintf(`catalogs:
  %s:
    sources:
      - id: "test-source"
        name: "Test Data Source"
        type: "yaml"
        properties:
          loaderConfigPath: "./%s"
          yamlCatalogPath: "./%s"
`, catalogName, loaderSourcesFile, entityDataFile)

	return os.WriteFile(outputPath, []byte(content), 0644)
}

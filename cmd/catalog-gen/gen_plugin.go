package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// initCatalogPlugin initializes a new catalog plugin for the unified catalog server.
func initCatalogPlugin(name, entityName, packageName, outputDir string) error {
	fmt.Printf("Initializing catalog: %s\n", name)
	fmt.Printf("  Entity: %s\n", entityName)
	fmt.Printf("  Package: %s\n", packageName)
	fmt.Printf("  Output: %s\n", outputDir)

	// Create plugin-specific directory structure (no cmd/, no manifests/)
	dirs := []string{
		outputDir,
		filepath.Join(outputDir, "internal", "catalog", "providers"),
		filepath.Join(outputDir, "internal", "db", "models"),
		filepath.Join(outputDir, "internal", "db", "service"),
		filepath.Join(outputDir, "internal", "server", "openapi"),
		filepath.Join(outputDir, "pkg", "openapi"),
		filepath.Join(outputDir, "api", "generated"),
		filepath.Join(outputDir, "docs"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create catalog.yaml config
	config := CatalogConfig{
		APIVersion: "catalog.kubeflow.org/v1alpha1",
		Kind:       "CatalogConfig",
		Metadata: CatalogMetadata{
			Name: name,
		},
		Spec: CatalogSpec{
			Package: packageName,
			Entity: EntityConfig{
				Name: entityName,
				Properties: []PropertyConfig{
					{Name: "description", Type: "string"},
				},
			},
			Providers: []ProviderConfig{
				{Type: "yaml"},
			},
			API: APIConfig{
				BasePath: fmt.Sprintf("/api/%s/v1alpha1", name),
				Port:     8081,
			},
		},
	}

	configPath := filepath.Join(outputDir, "catalog.yaml")
	configFile, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer func() { _ = configFile.Close() }()

	encoder := yaml.NewEncoder(configFile)
	encoder.SetIndent(2)
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	// Change to output directory to use generate functions
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	if err := os.Chdir(outputDir); err != nil {
		return fmt.Errorf("failed to change to output directory: %w", err)
	}
	defer func() { _ = os.Chdir(originalDir) }()

	fmt.Println("\n=== Generating plugin files ===")

	// Generate plugin.go and register.go (instead of cmd/main.go)
	if err := generatePluginFiles(config); err != nil {
		return fmt.Errorf("failed to generate plugin files: %w", err)
	}

	fmt.Println("\n=== Generating editable files (created once, you can modify) ===")

	// Generate repository (same as standalone)
	if err := generateRepository(config); err != nil {
		return fmt.Errorf("failed to generate repository: %w", err)
	}

	// Generate OpenAPI main file
	if err := generateOpenAPIMain(config); err != nil {
		return fmt.Errorf("failed to generate OpenAPI main: %w", err)
	}
	fmt.Printf("  Created: api/openapi/src/openapi.yaml\n")

	// Generate OpenAPI service implementation
	if err := generateOpenAPIServiceImpl(config); err != nil {
		return fmt.Errorf("failed to generate OpenAPI service impl: %w", err)
	}
	fmt.Printf("  Created: internal/server/openapi/api_%s_service_impl.go\n", strings.ToLower(entityName))

	// Generate YAML provider
	for _, provider := range config.Spec.Providers {
		if provider.Type == "yaml" {
			if err := generateYAMLProvider(config); err != nil {
				return fmt.Errorf("failed to generate YAML provider: %w", err)
			}
		}
	}

	// Generate Makefile (simplified for plugins)
	if err := generatePluginMakefile(config); err != nil {
		return fmt.Errorf("failed to generate Makefile: %w", err)
	}
	fmt.Printf("  Created: Makefile\n")

	// Generate README for plugin
	if err := generatePluginREADME(config); err != nil {
		return fmt.Errorf("failed to generate README: %w", err)
	}
	fmt.Printf("  Created: README.md\n")

	// Generate .gitignore
	if err := generateGitignore(); err != nil {
		return fmt.Errorf("failed to generate .gitignore: %w", err)
	}
	fmt.Printf("  Created: .gitignore\n")

	// Generate .openapi-generator-ignore
	if err := generateOpenAPIGeneratorIgnore(); err != nil {
		return fmt.Errorf("failed to generate .openapi-generator-ignore: %w", err)
	}
	fmt.Printf("  Created: .openapi-generator-ignore\n")

	// Generate Claude Code skills and commands
	if err := generateClaudeSkills(config); err != nil {
		return fmt.Errorf("failed to generate Claude skills: %w", err)
	}

	fmt.Println("\n=== Generating auto-regenerated files ===")

	// Generate entity model
	if err := generateEntityModel(config); err != nil {
		return fmt.Errorf("failed to generate entity model: %w", err)
	}

	// Generate datastore spec
	if err := generateDatastoreSpec(config); err != nil {
		return fmt.Errorf("failed to generate datastore spec: %w", err)
	}

	// Generate OpenAPI components
	if err := generateOpenAPIComponents(config); err != nil {
		return fmt.Errorf("failed to generate OpenAPI components: %w", err)
	}

	// Generate loader
	if err := generateLoader(config); err != nil {
		return fmt.Errorf("failed to generate loader: %w", err)
	}

	fmt.Printf("\nCatalog %s initialized successfully!\n", name)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Run 'make gen/openapi-server' to generate OpenAPI handlers")
	fmt.Println("  2. Import the plugin in cmd/catalog-server/main.go:")
	fmt.Printf("     _ \"%s\"\n", packageName)
	fmt.Println("  3. Add the plugin to sources.yaml under catalogs:")
	fmt.Printf("     %s:\n", name)
	fmt.Println("       sources:")
	fmt.Println("         - id: \"my-source\"")
	fmt.Println("           type: \"yaml\"")

	return nil
}

// generatePluginFiles generates the plugin.go and register.go files.
func generatePluginFiles(config CatalogConfig) error {
	entityName := config.Spec.Entity.Name

	// Always use 'any' for artifact type to avoid import cycles
	artifactType := "any"

	// Extract the last segment of the catalog name for the Go package name
	// e.g., "catalog/plugins/mcp" -> "mcp"
	packageName := filepath.Base(config.Metadata.Name)

	data := map[string]any{
		"Name":         packageName,
		"PackageName":  packageName,
		"EntityName":   entityName,
		"Package":      config.Spec.Package,
		"BasePath":     config.Spec.API.BasePath,
		"ArtifactType": artifactType,
	}

	// Generate plugin.go
	if err := executeTemplate(TmplPluginPlugin, "plugin.go", data); err != nil {
		return fmt.Errorf("failed to generate plugin.go: %w", err)
	}
	fmt.Printf("  Created: plugin.go\n")

	// Generate register.go
	if err := executeTemplate(TmplPluginRegister, "register.go", data); err != nil {
		return fmt.Errorf("failed to generate register.go: %w", err)
	}
	fmt.Printf("  Created: register.go\n")

	return nil
}

// generatePluginMakefile generates a Makefile for plugins with proper OpenAPI generation.
func generatePluginMakefile(config CatalogConfig) error {
	content := fmt.Sprintf(`# Generated by catalog-gen - you can modify this file

PACKAGE := %s
CATALOG_NAME := %s
PROJECT_ROOT := $(shell pwd)

# Use Docker by default, but allow local binary if available
OPENAPI_GENERATOR ?= docker run --rm -v $(PROJECT_ROOT):/local -w /local openapitools/openapi-generator-cli:v7.13.0

.PHONY: all build test gen/openapi gen/openapi-server gen/openapi-client clean

all: gen/openapi-server build

build:
	go build ./...

test:
	go test ./...

# Merge OpenAPI spec with generated components
api/openapi/openapi.yaml: api/openapi/src/openapi.yaml api/openapi/src/generated/components.yaml
	@echo "Merging OpenAPI specs..."
	@mkdir -p api/openapi
	@cat api/openapi/src/openapi.yaml > api/openapi/openapi.yaml
	@echo "" >> api/openapi/openapi.yaml
	@cat api/openapi/src/generated/components.yaml >> api/openapi/openapi.yaml
	@echo "Merged to api/openapi/openapi.yaml"

gen/openapi: api/openapi/openapi.yaml

# Generate OpenAPI server code (controllers, models, routers)
# The generated controller calls your service implementation in internal/server/openapi/api_*_service_impl.go
gen/openapi-server: api/openapi/openapi.yaml
	@echo "Generating OpenAPI server code..."
	@mkdir -p internal/server/openapi
	$(OPENAPI_GENERATOR) generate \
		-i api/openapi/openapi.yaml \
		-g go-server \
		-o internal/server/openapi \
		--package-name openapi \
		--ignore-file-override .openapi-generator-ignore \
		--additional-properties=outputAsLibrary=true,router=chi,sourceFolder=,onlyInterfaces=true,isGoSubmodule=true,enumClassPrefix=true
	@echo "Running goimports..."
	@command -v goimports >/dev/null 2>&1 && goimports -w internal/server/openapi || echo "goimports not found, skipping"
	@echo "Done"

# Generate OpenAPI client code (optional - for SDK generation)
gen/openapi-client: api/openapi/openapi.yaml
	@echo "Generating OpenAPI client code..."
	@mkdir -p pkg/openapi
	$(OPENAPI_GENERATOR) generate \
		-i api/openapi/openapi.yaml \
		-g go \
		-o pkg/openapi \
		--package-name openapi \
		--ignore-file-override .openapi-generator-ignore \
		--additional-properties=isGoSubmodule=true,enumClassPrefix=true
	@command -v goimports >/dev/null 2>&1 && goimports -w pkg/openapi || echo "goimports not found, skipping"
	@echo "Done"

clean:
	rm -rf internal/server/openapi/.openapi-generator
	rm -rf pkg/openapi/.openapi-generator
	rm -f api/openapi/openapi.yaml
`, config.Spec.Package, config.Metadata.Name)

	return os.WriteFile("Makefile", []byte(content), 0644)
}

// generatePluginREADME generates a README.md for the plugin.
func generatePluginREADME(config CatalogConfig) error {
	content := fmt.Sprintf(`# %s Catalog Plugin

This is a catalog plugin generated by catalog-gen for the unified catalog server.

## Overview

- **Entity**: %s
- **Package**: %s
- **API Base Path**: %s

## Usage

### 1. Generate OpenAPI Handlers

`+"```bash"+`
make gen/openapi-server
`+"```"+`

### 2. Import the Plugin

Add the plugin import to `+"`cmd/catalog-server/main.go`"+`:

`+"```go"+`
import (
    // Import plugins - their init() registers them
    _ "%s"
)
`+"```"+`

### 3. Configure Sources

Add the plugin configuration to your `+"`sources.yaml`"+`:

`+"```yaml"+`
catalogs:
  %s:
    sources:
      - id: "my-source"
        name: "My Data Source"
        type: "yaml"
        properties:
          yamlCatalogPath: "./data/%s.yaml"
`+"```"+`

### 4. Build and Run

`+"```bash"+`
go build ./cmd/catalog-server
./catalog-server --sources=./sources.yaml --listen=:8080
`+"```"+`

## Development

### Adding Properties

Edit `+"`catalog.yaml`"+` and run:

`+"```bash"+`
catalog-gen generate
`+"```"+`

### Adding Artifacts

`+"```bash"+`
catalog-gen add-artifact MyArtifact
`+"```"+`

## Files

| File | Description |
|------|-------------|
| `+"`plugin.go`"+` | Plugin implementation (auto-generated) |
| `+"`register.go`"+` | Plugin registration (auto-generated) |
| `+"`internal/db/models/`"+` | Entity models |
| `+"`internal/db/service/`"+` | Repository implementations |
| `+"`internal/catalog/`"+` | Data loader and providers |
| `+"`internal/server/openapi/`"+` | API handler implementations |
`, config.Metadata.Name, config.Spec.Entity.Name, config.Spec.Package,
		config.Spec.API.BasePath, config.Spec.Package, config.Metadata.Name,
		strings.ToLower(config.Spec.Entity.Name))

	return os.WriteFile("README.md", []byte(content), 0644)
}

// generateClaudeSkills generates Claude Code skills and commands for the plugin.
func generateClaudeSkills(config CatalogConfig) error {
	entityName := config.Spec.Entity.Name
	entityNameLower := strings.ToLower(entityName)

	data := map[string]any{
		"Name":            config.Metadata.Name,
		"EntityName":      entityName,
		"EntityNameLower": entityNameLower,
		"Package":         config.Spec.Package,
		"HasArtifacts":    len(config.Spec.Artifacts) > 0,
	}

	// Create directories
	if err := os.MkdirAll(".claude/commands", 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(".claude/skills", 0755); err != nil {
		return err
	}

	// Generate CLAUDE.md
	if err := executeTemplate(TmplAgentClaudeMD, "CLAUDE.md", data); err != nil {
		return err
	}
	fmt.Printf("  Created: CLAUDE.md\n")

	// Generate commands
	commands := []struct{ tmpl, file string }{
		{TmplAgentCmdAddProperty, ".claude/commands/add-property.md"},
		{TmplAgentCmdAddArtifact, ".claude/commands/add-artifact.md"},
		{TmplAgentCmdAddArtifactProp, ".claude/commands/add-artifact-property.md"},
		{TmplAgentCmdRegenerate, ".claude/commands/regenerate.md"},
		{TmplAgentCmdFixBuild, ".claude/commands/fix-build.md"},
	}
	for _, cmd := range commands {
		if err := executeTemplate(cmd.tmpl, cmd.file, data); err != nil {
			return err
		}
	}
	fmt.Printf("  Created: .claude/commands/ (%d commands)\n", len(commands))

	// Generate skills
	skills := []struct{ tmpl, file string }{
		{TmplAgentSkillAddProperty, ".claude/skills/add-property.md"},
		{TmplAgentSkillAddArtifact, ".claude/skills/add-artifact.md"},
		{TmplAgentSkillAddArtifactProp, ".claude/skills/add-artifact-property.md"},
		{TmplAgentSkillRegenerate, ".claude/skills/regenerate.md"},
	}
	for _, skill := range skills {
		if err := executeTemplate(skill.tmpl, skill.file, data); err != nil {
			return err
		}
	}
	fmt.Printf("  Created: .claude/skills/ (%d skills)\n", len(skills))

	return nil
}

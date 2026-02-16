package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newInitCmd() *cobra.Command {
	var (
		entityName  string
		packageName string
		outputDir   string
	)

	cmd := &cobra.Command{
		Use:   "init <name>",
		Short: "Initialize a new catalog project",
		Long: `Initialize a new catalog project with the specified name.

This creates the basic directory structure and configuration file
for a new catalog service.

Example:
  catalog-gen init mcp-catalog --entity=MCPServer --package=github.com/kubeflow/model-registry/mcp-catalog`,
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

			return initCatalog(name, entityName, packageName, outputDir)
		},
	}

	cmd.Flags().StringVar(&entityName, "entity", "", "Name of the main entity (e.g., MCPServer)")
	cmd.Flags().StringVar(&packageName, "package", "", "Go package path (e.g., github.com/kubeflow/model-registry/mcp-catalog)")
	cmd.Flags().StringVar(&outputDir, "output", "", "Output directory (defaults to catalog name)")

	return cmd
}

func initCatalog(name, entityName, packageName, outputDir string) error {
	fmt.Printf("Initializing catalog: %s\n", name)
	fmt.Printf("  Entity: %s\n", entityName)
	fmt.Printf("  Package: %s\n", packageName)
	fmt.Printf("  Output: %s\n", outputDir)

	// Create directory structure
	dirs := []string{
		outputDir,
		filepath.Join(outputDir, "cmd"),
		filepath.Join(outputDir, "internal", "catalog", "providers"),
		filepath.Join(outputDir, "internal", "db", "models"),
		filepath.Join(outputDir, "internal", "db", "service"),
		filepath.Join(outputDir, "internal", "server", "openapi"),
		filepath.Join(outputDir, "pkg", "openapi"),
		filepath.Join(outputDir, "api", "generated"),
		filepath.Join(outputDir, "manifests", "kustomize", "base"),
		filepath.Join(outputDir, "manifests", "kustomize", "overlays", "dev"),
		filepath.Join(outputDir, "docs"),
		filepath.Join(outputDir, ".claude", "commands"),
		filepath.Join(outputDir, ".claude", "skills"),
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

	fmt.Println("\n=== Generating editable files (created once, you can modify) ===")

	// Generate main.go
	if err := generateMainGo(name, entityName, packageName); err != nil {
		return fmt.Errorf("failed to generate main.go: %w", err)
	}
	fmt.Printf("  Created: cmd/%s.go\n", name)

	// Generate repository
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

	// Generate Makefile
	if err := generateMakefile(config); err != nil {
		return fmt.Errorf("failed to generate Makefile: %w", err)
	}
	fmt.Printf("  Created: Makefile\n")

	// Generate README
	if err := generateREADME(config); err != nil {
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

	// Generate Kustomize base manifests
	if err := generateKustomizeBase(config); err != nil {
		return fmt.Errorf("failed to generate kustomize base: %w", err)
	}
	fmt.Printf("  Created: manifests/kustomize/base/\n")

	fmt.Println("\n=== Generating non-editable files (regenerated by 'catalog-gen generate') ===")

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

	// Generate dev overlay
	if err := generateDevOverlay(config); err != nil {
		return fmt.Errorf("failed to generate dev overlay: %w", err)
	}
	fmt.Printf("  Created: manifests/kustomize/overlays/dev/\n")

	// Generate CLAUDE.md for AI agent context
	if err := generateAgentContext(config); err != nil {
		return fmt.Errorf("failed to generate agent context: %w", err)
	}

	// Generate agent skills and commands
	if err := generateAgentSkills(config); err != nil {
		return fmt.Errorf("failed to generate agent skills: %w", err)
	}

	fmt.Printf("\nCatalog %s initialized successfully!\n", name)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. cd", outputDir)
	fmt.Println("  2. Run 'make gen/openapi-server' to generate OpenAPI handlers")
	fmt.Println("  3. Configure database via environment variables (PGHOST, PGUSER, etc.)")
	fmt.Println("  4. Run 'go run ./cmd/...' to start the server")
	fmt.Println("\nTo add new properties, edit catalog.yaml and run 'catalog-gen generate'")
	fmt.Println("\nAI Agent Support:")
	fmt.Println("  - CLAUDE.md: Project context for AI assistants")
	fmt.Println("  - .claude/commands/: Slash commands (e.g., /run-local, /test-api)")
	fmt.Println("  - .claude/skills/: Detailed skill instructions")
	fmt.Println("  - See .claude/skills/README.md for full command reference")

	return nil
}

// generateMainGo generates the main entry point file.
func generateMainGo(name, entityName, packageName string) error {
	data := map[string]any{
		"Name":       name,
		"EntityName": entityName,
		"Package":    packageName,
		"Port":       8081,
	}

	return executeTemplate(TmplCmdMain, filepath.Join("cmd", fmt.Sprintf("%s.go", name)), data)
}

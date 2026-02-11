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
		Long: `Generate or regenerate the catalog code based on the catalog.yaml configuration.

This command reads the catalog.yaml file and generates:
  - Entity models and repositories
  - Artifact models and repositories
  - Provider implementations
  - OpenAPI specification
  - Kustomize manifests

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

	fmt.Printf("Regenerating code for catalog: %s\n", config.Metadata.Name)
	fmt.Println("Note: Only non-editable files are regenerated. Editable files are created by 'catalog-gen init'.")
	fmt.Println()

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

	// Generate post-generation documentation (if artifacts are configured)
	if len(config.Spec.Artifacts) > 0 {
		if err := generatePostGenerationDocs(config); err != nil {
			return fmt.Errorf("failed to generate post-generation docs: %w", err)
		}
	}

	// Generate CLAUDE.md for AI agent context
	if err := generateAgentContext(config); err != nil {
		return fmt.Errorf("failed to generate agent context: %w", err)
	}

	// Generate agent skills
	if err := generateAgentSkills(config); err != nil {
		return fmt.Errorf("failed to generate agent skills: %w", err)
	}

	fmt.Println("\nGeneration complete!")
	fmt.Println("\nIf you added new properties or artifacts to catalog.yaml, you may need to manually update:")
	fmt.Println("  - internal/db/service/<entity>.go (property mapping in converters)")
	fmt.Println("  - internal/server/openapi/api_*_service_impl.go (OpenAPI conversion)")
	fmt.Println("  - internal/catalog/providers/*.go (provider parsing)")
	fmt.Println("\nThen run 'make gen/openapi-all' to regenerate OpenAPI handlers and client.")

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

### 1. Update cmd/%[2]s.go

Add the artifact repository to NewServices call:

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
| `+"`cmd/%[2]s.go`"+` | Add artifact repository to NewServices | ☐ Manual |
%[9]s| `+"`manifests/kustomize/overlays/dev/artifacts.yaml`"+` | Populate with real data | ☐ Manual |
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

// =============================================================================
// Agent Context (CLAUDE.md)
// =============================================================================

func generateAgentContext(config CatalogConfig) error {
	entityName := config.Spec.Entity.Name
	lowerEntity := strings.ToLower(entityName)
	catalogName := config.Metadata.Name

	var entityProps strings.Builder
	for _, prop := range config.Spec.Entity.Properties {
		required := ""
		if prop.Required {
			required = " (required)"
		}
		entityProps.WriteString(fmt.Sprintf("  - `%s`: %s%s\n", prop.Name, prop.Type, required))
	}

	var artifactsSection strings.Builder
	if len(config.Spec.Artifacts) > 0 {
		artifactsSection.WriteString("\n## Artifacts\n\n")
		for _, artifact := range config.Spec.Artifacts {
			fullName := entityName + artifact.Name + "Artifact"
			artifactsSection.WriteString(fmt.Sprintf("### %s\n\n", artifact.Name))
			artifactsSection.WriteString(fmt.Sprintf("- **Model**: `internal/db/models/%s_%s_artifact.go`\n", lowerEntity, strings.ToLower(artifact.Name)))
			artifactsSection.WriteString(fmt.Sprintf("- **Type**: `models.%s`\n", fullName))
			artifactsSection.WriteString("- **Properties**:\n")
			for _, prop := range artifact.Properties {
				artifactsSection.WriteString(fmt.Sprintf("  - `%s`: %s\n", prop.Name, prop.Type))
			}
			artifactsSection.WriteString("\n")
		}
	}

	content := fmt.Sprintf(`# %s

> **Auto-generated by catalog-gen** - Regenerate with `+"`catalog-gen generate`"+`

This is a catalog service built with the catalog-gen scaffolding tool.

## Quick Reference

| Item | Value |
|------|-------|
| Catalog Name | %[1]s |
| Entity Type | %[2]s |
| API Port | %[3]d |
| Base Path | /api/%[1]s/v1alpha1 |

## Entity: %[2]s

- **Model**: `+"`internal/db/models/%[4]s.go`"+`
- **Repository**: `+"`internal/db/service/%[4]s.go`"+`

### Properties

%[5]s
%[6]s
## Common Tasks

### Adding a New Property

1. Edit `+"`catalog.yaml`"+`
2. Run `+"`catalog-gen generate`"+`
3. Update converters in `+"`internal/db/service/%[4]s.go`"+`

### Running Locally

`+"```bash"+`
go build ./cmd/%[1]s.go
./%[1]s --catalogs-path=manifests/kustomize/overlays/dev/sources.yaml
`+"```"+`

## Available Commands

| Command | Description |
|---------|-------------|
| `+"`/add-property`"+` | Add a new property to the entity |
| `+"`/add-artifact`"+` | Add a new artifact type |
| `+"`/regenerate`"+` | Regenerate code from catalog.yaml |
| `+"`/run-local`"+` | Start server locally |
| `+"`/test-api`"+` | Generate curl commands to test API |
| `+"`/seed-data`"+` | Populate sample data for testing |
`, catalogName, entityName, config.Spec.API.Port, lowerEntity,
		entityProps.String(), artifactsSection.String())

	if err := os.WriteFile("CLAUDE.md", []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write CLAUDE.md: %w", err)
	}

	fmt.Printf("  Updated: CLAUDE.md\n")
	return nil
}

// =============================================================================
// Agent Skills and Commands
// =============================================================================

func generateAgentSkills(config CatalogConfig) error {
	skillsDir := filepath.Join(".claude", "skills")
	if err := ensureDir(skillsDir); err != nil {
		return fmt.Errorf("failed to create skills directory: %w", err)
	}

	commandsDir := filepath.Join(".claude", "commands")
	if err := ensureDir(commandsDir); err != nil {
		return fmt.Errorf("failed to create commands directory: %w", err)
	}

	entityName := config.Spec.Entity.Name
	lowerEntity := strings.ToLower(entityName)
	catalogName := config.Metadata.Name

	// Generate skill files
	if err := generateAddPropertySkill(skillsDir, entityName, lowerEntity, catalogName); err != nil {
		return err
	}
	if err := generateAddArtifactSkill(skillsDir, entityName, catalogName); err != nil {
		return err
	}
	if err := generateRegenerateSkill(skillsDir, config); err != nil {
		return err
	}
	if err := generateRunLocalSkill(skillsDir, catalogName, config.Spec.API.Port); err != nil {
		return err
	}
	if err := generateTestAPISkill(skillsDir, entityName, lowerEntity, catalogName, config.Spec.API.Port); err != nil {
		return err
	}
	if err := generateSeedDataSkill(skillsDir, config); err != nil {
		return err
	}

	// Generate command files
	if err := generateAddPropertyCommand(commandsDir); err != nil {
		return err
	}
	if err := generateAddArtifactCommand(commandsDir); err != nil {
		return err
	}
	if err := generateRegenerateCommand(commandsDir); err != nil {
		return err
	}
	if err := generateRunLocalCommand(commandsDir); err != nil {
		return err
	}
	if err := generateTestAPICommand(commandsDir); err != nil {
		return err
	}
	if err := generateSeedDataCommand(commandsDir); err != nil {
		return err
	}

	// Generate skills README
	if err := generateSkillsReadme(skillsDir, catalogName); err != nil {
		return err
	}

	// Generate AI agent guide
	docsDir := "docs"
	if err := ensureDir(docsDir); err != nil {
		return err
	}
	if err := generateAIAgentGuide(docsDir, catalogName); err != nil {
		return err
	}

	fmt.Printf("  Updated: .claude/skills/\n")
	fmt.Printf("  Updated: .claude/commands/\n")
	fmt.Printf("  Updated: docs/ai-agent-guide.md\n")

	return nil
}

func generateAddPropertySkill(skillsDir, entityName, lowerEntity, catalogName string) error {
	content := fmt.Sprintf(`# Add Property Skill

Add a new property to the %s entity.

## Steps

1. Edit `+"`catalog.yaml`"+` and add property under `+"`spec.entity.properties`"+`:

`+"```yaml"+`
properties:
  - name: newProperty
    type: string  # string, integer, boolean
`+"```"+`

2. Run `+"`catalog-gen generate`"+`

3. Update `+"`internal/db/service/%s.go`"+` to map the new property

4. Run `+"`make gen/openapi-server`"+`
`, entityName, lowerEntity)

	return os.WriteFile(filepath.Join(skillsDir, "add-property.md"), []byte(content), 0644)
}

func generateAddArtifactSkill(skillsDir, entityName, catalogName string) error {
	content := fmt.Sprintf(`# Add Artifact Skill

Add a new artifact type linked to %s entities.

## Steps

1. Edit `+"`catalog.yaml`"+` and add artifact under `+"`spec.artifacts`"+`:

`+"```yaml"+`
artifacts:
  - name: NewArtifact
    properties:
      - name: uri
        type: string
`+"```"+`

2. Run `+"`catalog-gen generate`"+`

3. Follow steps in `+"`docs/post-artifact-generation.md`"+`
`, entityName)

	return os.WriteFile(filepath.Join(skillsDir, "add-artifact.md"), []byte(content), 0644)
}

// regenerateSkillProperty holds property data for the regenerate skill template.
type regenerateSkillProperty struct {
	Name      string
	Type      string
	FieldName string // PascalCase field name for Go structs (from camelCase: source_id -> SourceId)
	AttrName  string // Capitalized attribute name (source_id -> Source_id)
	VarName   string // camelCase variable name (source_id -> sourceId)
}

// regenerateSkillArtifact holds artifact data for the regenerate skill template.
type regenerateSkillArtifact struct {
	Name       string
	Properties []regenerateSkillProperty
}

// regenerateSkillData holds all data for the regenerate skill template.
type regenerateSkillData struct {
	EntityName      string
	EntityNameLower string
	Properties      []regenerateSkillProperty
	Artifacts       []regenerateSkillArtifact
	HasArtifacts    bool
}

func generateRegenerateSkill(skillsDir string, config CatalogConfig) error {
	entityName := config.Spec.Entity.Name

	// Build properties list (excluding 'name' which is a base attribute)
	var properties []regenerateSkillProperty
	for _, prop := range config.Spec.Entity.Properties {
		if prop.Name != "name" {
			properties = append(properties, regenerateSkillProperty{
				Name:      prop.Name,
				Type:      prop.Type,
				FieldName: capitalize(toCamelCase(prop.Name)),
				AttrName:  capitalize(prop.Name),
				VarName:   toCamelCase(prop.Name),
			})
		}
	}

	// Build artifacts list
	var artifacts []regenerateSkillArtifact
	for _, art := range config.Spec.Artifacts {
		var artProps []regenerateSkillProperty
		for _, prop := range art.Properties {
			if prop.Name != "name" {
				artProps = append(artProps, regenerateSkillProperty{
					Name:      prop.Name,
					Type:      prop.Type,
					FieldName: capitalize(toCamelCase(prop.Name)),
					AttrName:  capitalize(prop.Name),
					VarName:   toCamelCase(prop.Name),
				})
			}
		}
		artifacts = append(artifacts, regenerateSkillArtifact{
			Name:       art.Name,
			Properties: artProps,
		})
	}

	data := regenerateSkillData{
		EntityName:      entityName,
		EntityNameLower: strings.ToLower(entityName),
		Properties:      properties,
		Artifacts:       artifacts,
		HasArtifacts:    len(artifacts) > 0,
	}

	return executeTemplate(TmplAgentRegenerateSkill, filepath.Join(skillsDir, "regenerate.md"), data)
}

func generateRunLocalSkill(skillsDir, catalogName string, port int) error {
	content := fmt.Sprintf(`# Run Local Skill

Start the server locally for testing.

## Steps

1. Build: `+"`go build -o bin/%s ./cmd/`"+`

2. Run: `+"`./bin/%s --catalogs-path=manifests/kustomize/overlays/dev/sources.yaml`"+`

3. Test: `+"`curl http://localhost:%d/api/%s/v1alpha1/`"+`
`, catalogName, catalogName, port, catalogName)

	return os.WriteFile(filepath.Join(skillsDir, "run-local.md"), []byte(content), 0644)
}

func generateTestAPISkill(skillsDir, entityName, lowerEntity, catalogName string, port int) error {
	content := fmt.Sprintf(`# Test API Skill

Test the API endpoints.

## Endpoints

- List: `+"`GET /api/%s/v1alpha1/%ss`"+`
- Get: `+"`GET /api/%s/v1alpha1/%ss/{name}`"+`

## Example

`+"```bash"+`
curl http://localhost:%d/api/%s/v1alpha1/%ss
`+"```"+`
`, catalogName, lowerEntity, catalogName, lowerEntity, port, catalogName, lowerEntity)

	return os.WriteFile(filepath.Join(skillsDir, "test-api.md"), []byte(content), 0644)
}

func generateAddPropertyCommand(commandsDir string) error {
	content := `Add a new property to the entity. See .claude/skills/add-property.md for details.`
	return os.WriteFile(filepath.Join(commandsDir, "add-property.md"), []byte(content), 0644)
}

func generateAddArtifactCommand(commandsDir string) error {
	content := `Add a new artifact type. See .claude/skills/add-artifact.md for details.`
	return os.WriteFile(filepath.Join(commandsDir, "add-artifact.md"), []byte(content), 0644)
}

func generateRegenerateCommand(commandsDir string) error {
	content := `Regenerate code from catalog.yaml. See .claude/skills/regenerate.md for details.`
	return os.WriteFile(filepath.Join(commandsDir, "regenerate.md"), []byte(content), 0644)
}

func generateRunLocalCommand(commandsDir string) error {
	content := `Start the server locally. See .claude/skills/run-local.md for details.`
	return os.WriteFile(filepath.Join(commandsDir, "run-local.md"), []byte(content), 0644)
}

func generateTestAPICommand(commandsDir string) error {
	content := `Test the API endpoints. See .claude/skills/test-api.md for details.`
	return os.WriteFile(filepath.Join(commandsDir, "test-api.md"), []byte(content), 0644)
}

func generateSeedDataSkill(skillsDir string, config CatalogConfig) error {
	entityName := config.Spec.Entity.Name
	lowerEntity := strings.ToLower(entityName)
	catalogName := config.Metadata.Name

	// Build entity sample
	var entitySample strings.Builder
	entitySample.WriteString(fmt.Sprintf("  - name: example-%s-1\n", lowerEntity))
	entitySample.WriteString("    description: \"Example entity for testing\"\n")
	for _, prop := range config.Spec.Entity.Properties {
		if prop.Name == "description" {
			continue
		}
		entitySample.WriteString(fmt.Sprintf("    %s: \"sample-value\"\n", prop.Name))
	}

	// Build artifact samples
	var artifactSamples strings.Builder
	for _, artifact := range config.Spec.Artifacts {
		artifactSamples.WriteString(fmt.Sprintf("\n## %s Artifact Sample Data\n\n", artifact.Name))
		artifactSamples.WriteString("Edit `manifests/kustomize/overlays/dev/artifacts.yaml`:\n\n")
		artifactSamples.WriteString("```yaml\n")
		artifactSamples.WriteString(fmt.Sprintf("%s%sArtifacts:\n", lowerEntity, artifact.Name))
		artifactSamples.WriteString(fmt.Sprintf("  - %sName: example-%s-1\n", lowerEntity, lowerEntity))
		for _, prop := range artifact.Properties {
			artifactSamples.WriteString(fmt.Sprintf("    %s: \"sample-value\"\n", prop.Name))
		}
		artifactSamples.WriteString("```\n")
	}

	data := map[string]interface{}{
		"EntityNameLower": lowerEntity,
		"EntitySample":    entitySample.String(),
		"ArtifactSamples": artifactSamples.String(),
		"CatalogName":     catalogName,
		"Port":            config.Spec.API.Port,
	}

	content, err := executeTemplateToString(TmplAgentSeedDataSkill, data)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(skillsDir, "seed-data.md"), []byte(content), 0644)
}

func generateSeedDataCommand(commandsDir string) error {
	content, err := executeTemplateToString(TmplAgentSeedDataCmd, nil)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(commandsDir, "seed-data.md"), []byte(content), 0644)
}

func generateSkillsReadme(skillsDir, catalogName string) error {
	content := fmt.Sprintf(`# %s Skills Reference

Available skills for AI assistants:

| Skill | Description |
|-------|-------------|
| add-property | Add a new property to the entity |
| add-artifact | Add a new artifact type |
| regenerate | Regenerate code from catalog.yaml |
| run-local | Start server locally |
| test-api | Test the API endpoints |
| seed-data | Populate sample data for testing |

See individual .md files for detailed instructions.
`, catalogName)

	return os.WriteFile(filepath.Join(skillsDir, "README.md"), []byte(content), 0644)
}

func generateAIAgentGuide(docsDir, catalogName string) error {
	content := fmt.Sprintf(`# AI Agent Guide for %s

This document provides guidance for AI assistants working with this catalog.

## Project Structure

- `+"`catalog.yaml`"+` - Configuration (source of truth)
- `+"`cmd/`"+` - Entry point
- `+"`internal/db/models/`"+` - Entity models
- `+"`internal/db/service/`"+` - Repositories
- `+"`internal/server/openapi/`"+` - API handlers
- `+"`internal/catalog/providers/`"+` - Data providers

## Common Workflows

1. **Add Property**: Edit catalog.yaml, run catalog-gen generate
2. **Add Artifact**: Edit catalog.yaml, run catalog-gen generate, implement endpoint
3. **Test Changes**: Build, run locally, test with curl

## Key Files

- `+"`CLAUDE.md`"+` - Quick project reference
- `+"`.claude/skills/`"+` - Detailed skill instructions
- `+"`.claude/commands/`"+` - Available slash commands
`, catalogName)

	return os.WriteFile(filepath.Join(docsDir, "ai-agent-guide.md"), []byte(content), 0644)
}

package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

// generateEntityModel generates the entity model file.
func generateEntityModel(config CatalogConfig) error {
	entityName := config.Spec.Entity.Name
	lowerName := strings.ToLower(entityName)

	// Build properties (skip built-in fields)
	builtinFields := map[string]bool{
		"name": true, "externalid": true, "createtimesinceepoch": true,
		"lastupdatetimesinceepoch": true, "id": true,
	}
	var propDefs strings.Builder
	for _, prop := range config.Spec.Entity.Properties {
		if builtinFields[strings.ToLower(prop.Name)] {
			continue
		}
		goType := goTypeFromSpec(prop.Type)
		propDefs.WriteString(fmt.Sprintf("\t%s\t%s\n", capitalize(prop.Name), goType))
	}

	data := map[string]any{
		"EntityName":      entityName,
		"EntityNameLower": lowerName,
		"Properties":      propDefs.String(),
	}

	modelsDir := filepath.Join("internal", "db", "models")
	if err := ensureDir(modelsDir); err != nil {
		return err
	}

	outputPath := filepath.Join(modelsDir, fmt.Sprintf("%s.go", lowerName))
	fmt.Printf("  Generated: internal/db/models/%s.go\n", lowerName)
	return executeTemplate(TmplModelsEntity, outputPath, data)
}

// generateArtifactModel generates an artifact model file.
func generateArtifactModel(config CatalogConfig, artifact ArtifactConfig) error {
	entityName := config.Spec.Entity.Name
	artifactName := artifact.Name
	lowerEntityName := strings.ToLower(entityName)
	lowerArtifactName := strings.ToLower(artifactName)

	// Build properties
	var propDefs strings.Builder
	for _, prop := range artifact.Properties {
		goType := goTypeFromSpec(prop.Type)
		propDefs.WriteString(fmt.Sprintf("\t%s\t%s\n", capitalize(prop.Name), goType))
	}

	data := map[string]any{
		"EntityName":        entityName,
		"ArtifactName":      artifactName,
		"LowerEntityName":   lowerEntityName,
		"LowerArtifactName": lowerArtifactName,
		"Properties":        propDefs.String(),
	}

	modelsDir := filepath.Join("internal", "db", "models")
	if err := ensureDir(modelsDir); err != nil {
		return err
	}

	filename := fmt.Sprintf("%s_%s_artifact.go", lowerEntityName, lowerArtifactName)
	fmt.Printf("  Generated: internal/db/models/%s\n", filename)
	return executeTemplate(TmplModelsArtifact, filepath.Join(modelsDir, filename), data)
}

// generateBaseModels generates the base models file.
func generateBaseModels(outputDir string) error {
	return executeTemplate(TmplModelsBase, filepath.Join(outputDir, "internal", "db", "models", "base.go"), nil)
}

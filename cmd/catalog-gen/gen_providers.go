package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

// generateLoader generates the catalog loader file.
func generateLoader(config CatalogConfig) error {
	entityName := config.Spec.Entity.Name
	hasArtifacts := len(config.Spec.Artifacts) > 0

	// Build artifact type switch cases for SaveArtifact
	var artifactSaveCases strings.Builder
	var artifactDeleteCalls strings.Builder
	if hasArtifacts {
		for _, artifact := range config.Spec.Artifacts {
			artifactSaveCases.WriteString(fmt.Sprintf(`		case models.%s%sArtifact:
			_, err := services.%s%sArtifactRepository.Save(a, &entityID)
			return err
`, entityName, artifact.Name, entityName, artifact.Name))
			artifactDeleteCalls.WriteString(fmt.Sprintf(`	if err := services.%s%sArtifactRepository.DeleteByParentID(entityID); err != nil {
		return err
	}
`, entityName, artifact.Name))
		}
	}

	// Always use 'any' for artifact type to avoid import cycles with providers
	artifactType := "any"

	data := map[string]any{
		"EntityName":          entityName,
		"Package":             config.Spec.Package,
		"HasArtifacts":        hasArtifacts,
		"ArtifactType":        artifactType,
		"ArtifactSaveCases":   artifactSaveCases.String(),
		"ArtifactDeleteCalls": artifactDeleteCalls.String(),
	}

	catalogDir := filepath.Join("internal", "catalog")
	if err := ensureDir(catalogDir); err != nil {
		return err
	}

	fmt.Printf("  Generated: internal/catalog/loader.go\n")
	return executeTemplate(TmplCatalogLoader, filepath.Join(catalogDir, "loader.go"), data)
}

// generateYAMLProvider generates the YAML provider file.
func generateYAMLProvider(config CatalogConfig) error {
	entityName := config.Spec.Entity.Name
	lowerName := strings.ToLower(entityName)
	hasArtifacts := len(config.Spec.Artifacts) > 0

	// Create providers directory
	providersDir := filepath.Join("internal", "catalog", "providers")
	if err := ensureDir(providersDir); err != nil {
		return err
	}

	// Determine artifact type for the provider
	artifactType := "any"
	if hasArtifacts {
		artifactType = "catalog.Artifact"
	}

	// Build artifact structs and parsing code
	var artifactStructs strings.Builder
	var artifactParseCode strings.Builder
	var artifactMatchCode strings.Builder
	if hasArtifacts {
		// Generate struct for each artifact type
		for _, artifact := range config.Spec.Artifacts {
			lowerArtifactName := strings.ToLower(artifact.Name)
			artifactStructs.WriteString(fmt.Sprintf(`
// yaml%s%s represents a %s entry in the artifacts YAML file.
type yaml%s%s struct {
	%sName string  %s
	Name         string  %s
`, entityName, artifact.Name, artifact.Name, entityName, artifact.Name, entityName,
				"`json:\""+lowerName+"Name\" yaml:\""+lowerName+"Name\"`",
				"`json:\"name\" yaml:\"name\"`"))
			for _, prop := range artifact.Properties {
				goType := goTypeFromSpec(prop.Type)
				// Remove pointer for optional fields in yaml struct
				yamlGoType := strings.TrimPrefix(goType, "*")
				artifactStructs.WriteString(fmt.Sprintf("\t%s %s `json:\"%s,omitempty\" yaml:\"%s,omitempty\"`\n",
					capitalize(prop.Name), yamlGoType, prop.Name, prop.Name))
			}
			artifactStructs.WriteString("}\n")

			// Add to artifacts catalog struct
			artifactStructs.WriteString(fmt.Sprintf(`
// yaml%sArtifactsCatalog is the structure of the artifacts YAML file.
type yaml%sArtifactsCatalog struct {
	%ss []yaml%s%s %s
}
`, entityName, entityName, artifact.Name, entityName, artifact.Name,
				"`json:\""+lowerArtifactName+"s\" yaml:\""+lowerArtifactName+"s\"`"))
		}

		// Generate artifact parsing code
		artifactParseCode.WriteString(`
	// Parse artifacts file if provided
	artifactsByEntity := make(map[string][]catalog.Artifact)
	if artifactsData != nil {
`)
		for _, artifact := range config.Spec.Artifacts {
			lowerArtifactName := strings.ToLower(artifact.Name)
			artifactParseCode.WriteString(fmt.Sprintf(`		var %sArtifacts yaml%sArtifactsCatalog
		if err := k8syaml.UnmarshalStrict(artifactsData, &%sArtifacts); err == nil {
			for _, a := range %sArtifacts.%ss {
				entityName := a.%sName
				artifactName := a.Name
`, lowerArtifactName, entityName, lowerArtifactName, lowerArtifactName, artifact.Name, entityName))
			// Build property assignments
			var propAssignments strings.Builder
			for _, prop := range artifact.Properties {
				propName := capitalize(prop.Name)
				goType := goTypeFromSpec(prop.Type)
				if strings.HasPrefix(goType, "*") {
					// For pointer types, take address
					propAssignments.WriteString(fmt.Sprintf("\t\t\t\t\t%s: &a.%s,\n", propName, propName))
				} else {
					propAssignments.WriteString(fmt.Sprintf("\t\t\t\t\t%s: a.%s,\n", propName, propName))
				}
			}
			artifactParseCode.WriteString(fmt.Sprintf(`				artifact := models.New%s%sArtifact(&models.%s%sArtifactAttributes{
					Name: &artifactName,
%s				})
				artifactsByEntity[entityName] = append(artifactsByEntity[entityName], artifact)
			}
		}
`, entityName, artifact.Name, entityName, artifact.Name, propAssignments.String()))
		}
		artifactParseCode.WriteString(`	}
`)

		// Generate code to match artifacts to entities
		artifactMatchCode.WriteString(`
		// Attach artifacts to this entity
		if artifacts, ok := artifactsByEntity[name]; ok {
			record.Artifacts = artifacts
		}
`)
	}

	// Build entity property struct fields for yaml struct
	var entityPropertyFields strings.Builder
	var entityPropertyAssignments strings.Builder
	for _, prop := range config.Spec.Entity.Properties {
		if prop.Name == "description" || prop.Name == "externalId" || prop.Name == "name" || prop.Name == "customProperties" {
			continue // Skip base fields already in template
		}
		goType := goTypeFromSpec(prop.Type)
		// For yaml struct, use non-pointer types with omitempty
		yamlGoType := strings.TrimPrefix(goType, "*")
		if prop.Type == "array" {
			// Determine array item type from items field
			itemType := "string" // default
			if prop.Items != nil && prop.Items.Type != "" {
				itemType = strings.TrimPrefix(goTypeFromSpec(prop.Items.Type), "*")
			}
			yamlGoType = "[]" + itemType
		}
		propName := capitalize(prop.Name)
		entityPropertyFields.WriteString(fmt.Sprintf("\t%s %s `json:\"%s,omitempty\" yaml:\"%s,omitempty\"`\n",
			propName, yamlGoType, prop.Name, prop.Name))
		// Generate assignment code - handle pointers and arrays
		if prop.Type == "array" {
			entityPropertyAssignments.WriteString(fmt.Sprintf(`		// Handle %s (array -> comma-separated string)
		if len(item.%s) > 0 {
			%sStr := ""
			for i, v := range item.%s {
				if i > 0 {
					%sStr += ","
				}
				%sStr += v
			}
			entity.GetAttributes().%s = &%sStr
		}
`, propName, propName, prop.Name, propName, prop.Name, prop.Name, propName, prop.Name))
		} else if strings.HasPrefix(goType, "*") {
			// Pointer type - take address
			entityPropertyAssignments.WriteString(fmt.Sprintf("\t\tentity.GetAttributes().%s = &item.%s\n", propName, propName))
		} else {
			entityPropertyAssignments.WriteString(fmt.Sprintf("\t\tentity.GetAttributes().%s = item.%s\n", propName, propName))
		}
	}

	data := map[string]any{
		"Package":                     config.Spec.Package,
		"EntityName":                  entityName,
		"EntityNameLower":             lowerName,
		"HasArtifacts":                hasArtifacts,
		"ArtifactType":                artifactType,
		"ArtifactStructs":             artifactStructs.String(),
		"ArtifactParseCode":           artifactParseCode.String(),
		"ArtifactMatchCode":           artifactMatchCode.String(),
		"EntityPropertyFields":        entityPropertyFields.String(),
		"EntityPropertyAssignments":   entityPropertyAssignments.String(),
	}

	providerPath := filepath.Join(providersDir, "yaml_provider.go")
	if err := executeTemplate(TmplProvidersYAML, providerPath, data); err != nil {
		return err
	}
	fmt.Printf("  Generated: %s\n", providerPath)

	return nil
}

// generateProviderFile generates a provider file for the specified type.
func generateProviderFile(entityName, providerType string) error {
	data := map[string]any{
		"EntityName": entityName,
	}

	providersDir := filepath.Join("internal", "catalog", "providers")
	if err := ensureDir(providersDir); err != nil {
		return fmt.Errorf("failed to create providers directory: %w", err)
	}

	var templatePath string
	switch providerType {
	case "yaml":
		templatePath = TmplProvidersYAML
	case "http":
		templatePath = TmplProvidersHTTP
	default:
		return fmt.Errorf("unknown provider type: %s", providerType)
	}

	return executeTemplate(templatePath, filepath.Join(providersDir, fmt.Sprintf("%s_provider.go", providerType)), data)
}

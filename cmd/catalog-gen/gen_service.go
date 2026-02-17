package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

// generateRepository generates the entity repository file.
func generateRepository(config CatalogConfig) error {
	entityName := config.Spec.Entity.Name
	lowerName := strings.ToLower(entityName)

	// Build dynamic property mapping code
	propVarDecls := buildPropertyVarDeclarations(config.Spec.Entity.Properties)
	propReadCases := buildPropertyReadCases(config.Spec.Entity.Properties)
	propAttrAssignments := buildPropertyAttrAssignments(config.Spec.Entity.Properties)
	propWriteStatements := buildPropertyWriteStatements(config.Spec.Entity.Properties)

	data := map[string]any{
		"EntityName":          entityName,
		"EntityNameLower":     lowerName,
		"Package":             config.Spec.Package,
		"PropVarDecls":        propVarDecls,
		"PropReadCases":       propReadCases,
		"PropAttrAssignments": propAttrAssignments,
		"PropWriteStatements": propWriteStatements,
	}

	serviceDir := filepath.Join("internal", "db", "service")
	if err := ensureDir(serviceDir); err != nil {
		return err
	}

	outputPath := filepath.Join(serviceDir, fmt.Sprintf("%s.go", lowerName))
	fmt.Printf("  Generated: internal/db/service/%s.go\n", lowerName)
	return executeTemplate(TmplServiceRepository, outputPath, data)
}

// generateArtifactRepository generates an artifact repository file.
func generateArtifactRepository(config CatalogConfig, artifact ArtifactConfig) error {
	entityName := config.Spec.Entity.Name
	artifactName := artifact.Name
	lowerEntityName := strings.ToLower(entityName)
	lowerArtifactName := strings.ToLower(artifactName)

	// Build property mapping code for artifact properties
	propVarDecls := buildPropertyVarDeclarations(artifact.Properties)
	propReadCases := buildPropertyReadCases(artifact.Properties)
	propAttrAssignments := buildPropertyAttrAssignments(artifact.Properties)
	propWriteStatements := buildArtifactPropertyWriteStatements(artifact.Properties)

	data := map[string]any{
		"Package":             config.Spec.Package,
		"EntityName":          entityName,
		"ArtifactName":        artifactName,
		"LowerEntityName":     lowerEntityName,
		"LowerArtifactName":   lowerArtifactName,
		"PropVarDecls":        propVarDecls,
		"PropReadCases":       propReadCases,
		"PropAttrAssignments": propAttrAssignments,
		"PropWriteStatements": propWriteStatements,
	}

	serviceDir := filepath.Join("internal", "db", "service")
	if err := ensureDir(serviceDir); err != nil {
		return err
	}

	filename := fmt.Sprintf("%s_%s_artifact.go", lowerEntityName, lowerArtifactName)
	fmt.Printf("  Generated: internal/db/service/%s\n", filename)
	return executeTemplate(TmplServiceArtifactRepository, filepath.Join(serviceDir, filename), data)
}

// generateDatastoreSpec generates the datastore spec file.
func generateDatastoreSpec(config CatalogConfig) error {
	entityName := config.Spec.Entity.Name
	lowerEntityName := strings.ToLower(entityName)

	// Build property definitions from config
	var propDefs []string
	for _, prop := range config.Spec.Entity.Properties {
		propMethod := datastorePropertyMethod(prop.Type)
		propDefs = append(propDefs, fmt.Sprintf("\t\t\t%s(\"%s\")", propMethod, prop.Name))
	}

	// Join properties with ".\n" but don't add trailing dot
	propDefsStr := ""
	if len(propDefs) > 0 {
		propDefsStr = strings.Join(propDefs, ".\n") + ","
	}

	// Build artifact type constants
	var artifactConstants strings.Builder
	for _, artifact := range config.Spec.Artifacts {
		artifactConstants.WriteString(fmt.Sprintf("\t%s%sArtifactTypeName = \"kf.%s%sArtifact\"\n",
			entityName, artifact.Name, entityName, artifact.Name))
	}

	// Build artifact spec registrations
	var artifactSpecs strings.Builder
	for i, artifact := range config.Spec.Artifacts {
		// Build artifact property definitions
		var artifactPropDefs []string
		for _, prop := range artifact.Properties {
			propMethod := datastorePropertyMethod(prop.Type)
			artifactPropDefs = append(artifactPropDefs, fmt.Sprintf("\t\t\t%s(\"%s\")", propMethod, prop.Name))
		}
		artifactPropDefsStr := ""
		if len(artifactPropDefs) > 0 {
			artifactPropDefsStr = strings.Join(artifactPropDefs, ".\n") + ","
		}

		// Add trailing dot only if not the last artifact
		trailingDot := ""
		if i < len(config.Spec.Artifacts)-1 {
			trailingDot = "."
		}

		artifactSpecs.WriteString(fmt.Sprintf(`		AddArtifact(%s%sArtifactTypeName, datastore.NewSpecType(New%s%sArtifactRepository).
%s
		)%s
`, entityName, artifact.Name, entityName, artifact.Name, artifactPropDefsStr, trailingDot))
	}

	// Build Services struct fields for artifacts
	var artifactServiceFields strings.Builder
	var artifactServiceParams strings.Builder
	var artifactServiceAssignments strings.Builder
	for _, artifact := range config.Spec.Artifacts {
		artifactServiceFields.WriteString(fmt.Sprintf("\t%s%sArtifactRepository models.%s%sArtifactRepository\n",
			entityName, artifact.Name, entityName, artifact.Name))
		artifactServiceParams.WriteString(fmt.Sprintf("\t%s%sArtifactRepository models.%s%sArtifactRepository,\n",
			lowerEntityName, artifact.Name, entityName, artifact.Name))
		artifactServiceAssignments.WriteString(fmt.Sprintf("\t\t%s%sArtifactRepository: %s%sArtifactRepository,\n",
			entityName, artifact.Name, lowerEntityName, artifact.Name))
	}

	// If there are artifacts, the context registration needs to chain to them
	hasArtifacts := len(config.Spec.Artifacts) > 0
	contextTrailingDot := ""
	if hasArtifacts {
		contextTrailingDot = "."
	}

	data := map[string]any{
		"EntityName":                 entityName,
		"EntityNameLower":            lowerEntityName,
		"Package":                    config.Spec.Package,
		"PropertyDefs":               propDefsStr,
		"ArtifactConstants":          artifactConstants.String(),
		"ArtifactSpecs":              artifactSpecs.String(),
		"ArtifactServiceFields":      artifactServiceFields.String(),
		"ArtifactServiceParams":      artifactServiceParams.String(),
		"ArtifactServiceAssignments": artifactServiceAssignments.String(),
		"ContextTrailingDot":         contextTrailingDot,
	}

	serviceDir := filepath.Join("internal", "db", "service")
	if err := ensureDir(serviceDir); err != nil {
		return err
	}

	fmt.Printf("  Generated: internal/db/service/spec.go\n")
	return executeTemplate(TmplServiceSpec, filepath.Join(serviceDir, "spec.go"), data)
}

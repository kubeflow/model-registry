package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

// generateOpenAPIComponents generates the OpenAPI components file.
func generateOpenAPIComponents(config CatalogConfig) error {
	entityName := config.Spec.Entity.Name

	// Build properties for OpenAPI and collect required fields (skip built-in fields)
	builtinFields := map[string]bool{
		"name": true, "externalid": true, "createtimesinceepoch": true,
		"lastupdatetimesinceepoch": true, "id": true,
	}
	var propDefs strings.Builder
	var requiredFields strings.Builder
	requiredFields.WriteString("        - name\n") // name is always required
	for _, prop := range config.Spec.Entity.Properties {
		if builtinFields[strings.ToLower(prop.Name)] {
			continue
		}
		propDefs.WriteString(generateOpenAPIPropertyDef(prop, 8))
		if prop.Required {
			requiredFields.WriteString(fmt.Sprintf("        - %s\n", prop.Name))
		}
	}

	// Build artifact schemas if artifacts are configured
	var artifactSchemas strings.Builder
	if len(config.Spec.Artifacts) > 0 {
		// Generate individual artifact schemas
		for _, artifact := range config.Spec.Artifacts {
			artifactSchemas.WriteString(fmt.Sprintf("    %s%sArtifact:\n", entityName, artifact.Name))
			artifactSchemas.WriteString("      type: object\n")
			artifactSchemas.WriteString("      properties:\n")
			artifactSchemas.WriteString("        id:\n")
			artifactSchemas.WriteString("          type: string\n")
			artifactSchemas.WriteString("          readOnly: true\n")
			artifactSchemas.WriteString("        name:\n")
			artifactSchemas.WriteString("          type: string\n")
			artifactSchemas.WriteString("        artifactType:\n")
			artifactSchemas.WriteString("          type: string\n")
			for _, prop := range artifact.Properties {
				artifactSchemas.WriteString(generateOpenAPIPropertyDef(prop, 8))
			}
		}

		// Generate artifact list schema
		artifactSchemas.WriteString(fmt.Sprintf("    %sArtifactList:\n", entityName))
		artifactSchemas.WriteString("      type: object\n")
		artifactSchemas.WriteString("      properties:\n")
		artifactSchemas.WriteString("        items:\n")
		artifactSchemas.WriteString("          type: array\n")
		artifactSchemas.WriteString("          items:\n")
		if len(config.Spec.Artifacts) == 1 {
			artifactSchemas.WriteString(fmt.Sprintf("            $ref: '#/components/schemas/%s%sArtifact'\n", entityName, config.Spec.Artifacts[0].Name))
		} else {
			artifactSchemas.WriteString("            oneOf:\n")
			for _, artifact := range config.Spec.Artifacts {
				artifactSchemas.WriteString(fmt.Sprintf("              - $ref: '#/components/schemas/%s%sArtifact'\n", entityName, artifact.Name))
			}
		}
		artifactSchemas.WriteString("        nextPageToken:\n")
		artifactSchemas.WriteString("          type: string\n")
		artifactSchemas.WriteString("        size:\n")
		artifactSchemas.WriteString("          type: integer\n")
	}

	data := map[string]any{
		"EntityName":      entityName,
		"Properties":      propDefs.String(),
		"RequiredFields":  requiredFields.String(),
		"ArtifactSchemas": artifactSchemas.String(),
	}

	generatedDir := filepath.Join("api", "openapi", "src", "generated")
	if err := ensureDir(generatedDir); err != nil {
		return err
	}

	fmt.Printf("  Generated: api/openapi/src/generated/components.yaml\n")
	return executeTemplate(TmplAPIOpenAPIComponents, filepath.Join(generatedDir, "components.yaml"), data)
}

// generateOpenAPIMain generates the OpenAPI main spec file.
func generateOpenAPIMain(config CatalogConfig) error {
	entityName := config.Spec.Entity.Name
	lowerName := strings.ToLower(entityName)

	// Build artifact routes if artifacts are configured
	artifactRoutes := ""
	if len(config.Spec.Artifacts) > 0 {
		artifactRoutes = `
  /` + lowerName + `s/{name}/artifacts:
    get:
      summary: List artifacts for a ` + entityName + `
      operationId: get` + entityName + `Artifacts
      parameters:
        - name: name
          in: path
          required: true
          schema:
            type: string
        - name: pageSize
          in: query
          schema:
            type: integer
            default: 20
        - name: pageToken
          in: query
          schema:
            type: string
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/` + entityName + `ArtifactList'
        '404':
          description: Not found
`
	}

	data := map[string]any{
		"Name":            config.Metadata.Name,
		"EntityName":      entityName,
		"EntityNameLower": lowerName,
		"BasePath":        config.Spec.API.BasePath,
		"ArtifactRoutes":  artifactRoutes,
	}

	srcDir := filepath.Join("api", "openapi", "src")
	if err := ensureDir(srcDir); err != nil {
		return err
	}

	return executeTemplate(TmplAPIOpenAPIMain, filepath.Join(srcDir, "openapi.yaml"), data)
}

// generateOpenAPIServiceImpl generates the OpenAPI service implementation stub.
func generateOpenAPIServiceImpl(config CatalogConfig) error {
	entityName := config.Spec.Entity.Name
	lowerName := strings.ToLower(entityName)

	// Build dynamic property conversion code
	propConversions := buildOpenAPIPropertyConversions(config.Spec.Entity.Properties)

	data := map[string]any{
		"EntityName":      entityName,
		"EntityNameLower": lowerName,
		"Package":         config.Spec.Package,
		"PropConversions": propConversions,
	}

	openapiDir := filepath.Join("internal", "server", "openapi")
	if err := ensureDir(openapiDir); err != nil {
		return err
	}

	return executeTemplate(TmplServerOpenAPIServiceImpl, filepath.Join(openapiDir, fmt.Sprintf("api_%s_service_impl.go", lowerName)), data)
}

package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

// generateOpenAPIComponents generates the OpenAPI components file.
// The generated components use allOf composition with BaseResource from common.yaml.
func generateOpenAPIComponents(config CatalogConfig) error {
	entityName := config.Spec.Entity.Name

	// Base properties that come from BaseResource - skip these in entity definition
	// These are defined in api/openapi/src/lib/common.yaml
	baseResourceProperties := map[string]bool{
		"name":                     true,
		"id":                       true,
		"externalid":               true,
		"description":              true,
		"customproperties":         true,
		"createtimesinceepoch":     true,
		"lastupdatetimesinceepoch": true,
	}

	var propDefs strings.Builder
	var requiredFields strings.Builder
	for _, prop := range config.Spec.Entity.Properties {
		if baseResourceProperties[strings.ToLower(prop.Name)] {
			continue
		}
		// Use 12 spaces for properties inside allOf structure
		propDefs.WriteString(generateOpenAPIPropertyDef(prop, 12))
		if prop.Required {
			fmt.Fprintf(&requiredFields, "            - %s\n", prop.Name)
		}
	}

	// Build artifact schemas if artifacts are configured
	// Artifacts also use allOf composition with BaseResource
	var artifactSchemas strings.Builder
	if len(config.Spec.Artifacts) > 0 {
		// Generate individual artifact schemas using allOf composition
		for _, artifact := range config.Spec.Artifacts {
			fmt.Fprintf(&artifactSchemas, "    %s%sArtifact:\n", entityName, artifact.Name)
			artifactSchemas.WriteString("      allOf:\n")
			artifactSchemas.WriteString("        - $ref: '#/components/schemas/BaseResource'\n")
			artifactSchemas.WriteString("        - type: object\n")
			artifactSchemas.WriteString("          properties:\n")
			artifactSchemas.WriteString("            artifactType:\n")
			artifactSchemas.WriteString("              type: string\n")
			for _, prop := range artifact.Properties {
				artifactSchemas.WriteString(generateOpenAPIPropertyDef(prop, 12))
			}
		}

		// Generate artifact list schema using allOf composition
		fmt.Fprintf(&artifactSchemas, "    %sArtifactList:\n", entityName)
		artifactSchemas.WriteString("      allOf:\n")
		artifactSchemas.WriteString("        - $ref: '#/components/schemas/BaseResourceList'\n")
		artifactSchemas.WriteString("        - type: object\n")
		artifactSchemas.WriteString("          properties:\n")
		artifactSchemas.WriteString("            items:\n")
		artifactSchemas.WriteString("              type: array\n")
		artifactSchemas.WriteString("              items:\n")
		if len(config.Spec.Artifacts) == 1 {
			fmt.Fprintf(&artifactSchemas, "                $ref: '#/components/schemas/%s%sArtifact'\n", entityName, config.Spec.Artifacts[0].Name)
		} else {
			artifactSchemas.WriteString("                oneOf:\n")
			for _, artifact := range config.Spec.Artifacts {
				fmt.Fprintf(&artifactSchemas, "                  - $ref: '#/components/schemas/%s%sArtifact'\n", entityName, artifact.Name)
			}
		}
	}

	data := map[string]any{
		"EntityName":      entityName,
		"Properties":      strings.TrimSpace(propDefs.String()),
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

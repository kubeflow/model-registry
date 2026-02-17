package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// loadConfig reads and parses catalog.yaml from the current directory.
func loadConfig() (CatalogConfig, error) {
	configPath := "catalog.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return CatalogConfig{}, fmt.Errorf("catalog.yaml not found in current directory. Run 'catalog-gen init' first")
	}

	configData, err := os.ReadFile(configPath)
	if err != nil {
		return CatalogConfig{}, fmt.Errorf("failed to read catalog.yaml: %w", err)
	}

	var config CatalogConfig
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return CatalogConfig{}, fmt.Errorf("failed to parse catalog.yaml: %w", err)
	}

	return config, nil
}

// saveConfig writes the config back to catalog.yaml.
func saveConfig(config CatalogConfig) error {
	configFile, err := os.Create("catalog.yaml")
	if err != nil {
		return fmt.Errorf("failed to open catalog.yaml: %w", err)
	}
	defer func() { _ = configFile.Close() }()

	encoder := yaml.NewEncoder(configFile)
	encoder.SetIndent(2)
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to write catalog.yaml: %w", err)
	}

	return nil
}

// ensureDir creates a directory if it doesn't exist.
func ensureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// isPluginContext detects if we're running in a plugin context by checking the current working directory.
func isPluginContext() bool {
	wd, err := os.Getwd()
	if err != nil {
		return false
	}
	// We're in a plugin if the working directory contains "/catalog/plugins/"
	return strings.Contains(wd, "/catalog/plugins/")
}

// ensureCommonLibSymlink creates a symlink at api/openapi/src/lib pointing to the
// shared common schemas in the repo root. This allows plugin OpenAPI specs to reference
// shared schemas (BaseResource, etc.) via relative paths like 'lib/common.yaml'.
// The symlink is only created in plugin contexts and is a no-op otherwise.
func ensureCommonLibSymlink() error {
	if !isPluginContext() {
		return nil
	}

	symlinkPath := filepath.Join("api", "openapi", "src", "lib")

	// If the symlink (or directory) already exists, nothing to do
	if _, err := os.Lstat(symlinkPath); err == nil {
		return nil
	}

	// Find the repo root via git
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return fmt.Errorf("failed to find git repo root: %w", err)
	}
	repoRoot := strings.TrimSpace(string(out))

	// Target is the shared lib directory at the repo root
	target := filepath.Join(repoRoot, "api", "openapi", "src", "lib")
	if _, err := os.Stat(target); os.IsNotExist(err) {
		return fmt.Errorf("common lib directory not found at %s", target)
	}

	// Ensure parent directory exists
	if err := ensureDir(filepath.Dir(symlinkPath)); err != nil {
		return err
	}

	// Compute relative path from the symlink's parent to the target
	symlinkParent, err := filepath.Abs(filepath.Dir(symlinkPath))
	if err != nil {
		return fmt.Errorf("failed to resolve symlink parent: %w", err)
	}
	relTarget, err := filepath.Rel(symlinkParent, target)
	if err != nil {
		return fmt.Errorf("failed to compute relative path: %w", err)
	}

	if err := os.Symlink(relTarget, symlinkPath); err != nil {
		return fmt.Errorf("failed to create lib symlink: %w", err)
	}

	fmt.Printf("  Created symlink: %s -> %s\n", symlinkPath, relTarget)
	return nil
}

// capitalize returns the string with the first letter uppercased.
func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// goTypeFromSpec converts a spec type to a Go type.
func goTypeFromSpec(specType string) string {
	switch specType {
	case "string":
		return "*string"
	case "integer", "int":
		return "*int32"
	case "int64":
		return "*int64"
	case "boolean", "bool":
		return "*bool"
	case "number", "float", "double":
		return "*float64"
	case "array":
		return "[]string"
	default:
		return "*string"
	}
}

// openAPIType converts a spec type to an OpenAPI type.
func openAPIType(specType string) string {
	switch specType {
	case "string":
		return "string"
	case "integer", "int", "int64":
		return "integer"
	case "boolean", "bool":
		return "boolean"
	case "number", "float", "double":
		return "number"
	case "array":
		return "array"
	default:
		return "string"
	}
}

// datastorePropertyMethod returns the datastore method for a property type.
func datastorePropertyMethod(specType string) string {
	switch specType {
	case "string":
		return "AddString"
	case "integer", "int", "int64":
		return "AddInt"
	case "boolean", "bool":
		return "AddBool"
	case "number", "float", "double":
		return "AddDouble"
	case "struct", "object":
		return "AddStruct"
	default:
		return "AddString"
	}
}

// mlmdValueField returns the MLMD value field name for a property type.
func mlmdValueField(specType string) string {
	switch specType {
	case "string":
		return "StringValue"
	case "integer", "int":
		return "IntValue"
	case "int64":
		return "IntValue"
	case "boolean", "bool":
		return "IntValue" // Bools are stored as int (0/1)
	case "number", "float", "double":
		return "DoubleValue"
	default:
		return "StringValue"
	}
}

// isStructType returns true if the type is a struct/object type.
func isStructType(specType string) bool {
	switch specType {
	case "struct", "object":
		return true
	default:
		return false
	}
}

// buildPropertyVarDeclarations generates variable declarations for reading properties.
func buildPropertyVarDeclarations(props []PropertyConfig) string {
	var sb strings.Builder
	for _, prop := range props {
		goType := goTypeFromSpec(prop.Type)
		varName := strings.ToLower(prop.Name[:1]) + prop.Name[1:] // camelCase
		if isStructType(prop.Type) {
			sb.WriteString(fmt.Sprintf("\tvar %s %s // TODO: struct type requires manual handling\n", varName, goType))
		} else {
			sb.WriteString(fmt.Sprintf("\tvar %s %s\n", varName, goType))
		}
	}
	return sb.String()
}

// buildPropertyReadCases generates switch cases for reading properties from ContextProperty.
func buildPropertyReadCases(props []PropertyConfig) string {
	var sb strings.Builder
	for _, prop := range props {
		varName := strings.ToLower(prop.Name[:1]) + prop.Name[1:]
		if isStructType(prop.Type) {
			sb.WriteString(fmt.Sprintf("\t\tcase \"%s\":\n", prop.Name))
			sb.WriteString(fmt.Sprintf("\t\t\t// TODO: Struct property '%s' requires manual conversion.\n", prop.Name))
			sb.WriteString("\t\t\t// Example: json.Unmarshal([]byte(*p.StringValue), &" + varName + ")\n")
		} else {
			valueField := mlmdValueField(prop.Type)
			sb.WriteString(fmt.Sprintf("\t\tcase \"%s\":\n", prop.Name))
			sb.WriteString(fmt.Sprintf("\t\t\tif p.%s != nil {\n", valueField))
			sb.WriteString(fmt.Sprintf("\t\t\t\t%s = p.%s\n", varName, valueField))
			sb.WriteString("\t\t\t}\n")
		}
	}
	return sb.String()
}

// buildPropertyAttrAssignments generates attribute assignments for properties.
func buildPropertyAttrAssignments(props []PropertyConfig) string {
	var sb strings.Builder
	for _, prop := range props {
		varName := strings.ToLower(prop.Name[:1]) + prop.Name[1:]
		fieldName := capitalize(prop.Name)
		sb.WriteString(fmt.Sprintf("\t\t\t%s: %s,\n", fieldName, varName))
	}
	return sb.String()
}

// buildPropertyWriteStatements generates property write statements for entity properties.
func buildPropertyWriteStatements(props []PropertyConfig) string {
	var sb strings.Builder
	for _, prop := range props {
		propName := capitalize(prop.Name)
		switch prop.Type {
		case "string":
			sb.WriteString(fmt.Sprintf(`	if attrs.%s != nil {
		props = append(props, schema.ContextProperty{
			ContextID:   entityID,
			Name:        "%s",
			StringValue: attrs.%s,
		})
	}
`, propName, prop.Name, propName))
		case "integer", "int":
			sb.WriteString(fmt.Sprintf(`	if attrs.%s != nil {
		props = append(props, schema.ContextProperty{
			ContextID: entityID,
			Name:      "%s",
			IntValue:  attrs.%s,
		})
	}
`, propName, prop.Name, propName))
		case "boolean", "bool":
			sb.WriteString(fmt.Sprintf(`	if attrs.%s != nil {
		boolVal := int64(0)
		if *attrs.%s {
			boolVal = 1
		}
		props = append(props, schema.ContextProperty{
			ContextID: entityID,
			Name:      "%s",
			IntValue:  &boolVal,
		})
	}
`, propName, propName, prop.Name))
		}
	}
	return sb.String()
}

// buildArtifactPropertyWriteStatements generates property write statements for artifact properties.
func buildArtifactPropertyWriteStatements(properties []PropertyConfig) string {
	var sb strings.Builder
	for _, prop := range properties {
		propName := capitalize(prop.Name)
		switch prop.Type {
		case "string":
			sb.WriteString(fmt.Sprintf(`		if attr.%s != nil {
			properties = append(properties, schema.ArtifactProperty{
				ArtifactID:  artifactID,
				Name:        "%s",
				StringValue: attr.%s,
			})
		}
`, propName, prop.Name, propName))
		case "integer", "int":
			sb.WriteString(fmt.Sprintf(`		if attr.%s != nil {
			properties = append(properties, schema.ArtifactProperty{
				ArtifactID: artifactID,
				Name:       "%s",
				IntValue:   attr.%s,
			})
		}
`, propName, prop.Name, propName))
		case "boolean", "bool":
			sb.WriteString(fmt.Sprintf(`		if attr.%s != nil {
			boolVal := int64(0)
			if *attr.%s {
				boolVal = 1
			}
			properties = append(properties, schema.ArtifactProperty{
				ArtifactID: artifactID,
				Name:       "%s",
				IntValue:   &boolVal,
			})
		}
`, propName, propName, prop.Name))
		}
	}
	return sb.String()
}

// buildOpenAPIPropertyConversions generates OpenAPI property conversion code.
func buildOpenAPIPropertyConversions(props []PropertyConfig) string {
	var sb strings.Builder
	for _, prop := range props {
		fieldName := capitalize(prop.Name)
		sb.WriteString(fmt.Sprintf(`	if attrs.%s != nil {
		result.%s = *attrs.%s
	}
`, fieldName, fieldName, fieldName))
	}
	return sb.String()
}

// generateOpenAPIPropertyDef generates an OpenAPI property definition.
func generateOpenAPIPropertyDef(prop PropertyConfig, indent int) string {
	spaces := strings.Repeat(" ", indent)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s%s:\n%s  type: %s\n", spaces, prop.Name, spaces, openAPIType(prop.Type)))
	if prop.Type == "array" && prop.Items != nil {
		sb.WriteString(fmt.Sprintf("%s  items:\n%s    type: %s\n", spaces, spaces, openAPIType(prop.Items.Type)))
	}
	return sb.String()
}

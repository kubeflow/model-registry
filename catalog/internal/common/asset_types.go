// Package common provides shared types and utilities used across catalog packages.
// This package exists to avoid circular imports between catalog and mcp packages.
package common

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/util/yaml"
)

const (
	// YamlCatalogPathKey is the property key for the YAML catalog file path
	YamlCatalogPathKey = "yamlCatalogPath"
)

// AssetType represents the type of assets contained in a catalog source
type AssetType string

const (
	// AssetTypeModels indicates the source contains AI/ML models
	AssetTypeModels AssetType = "models"
	// AssetTypeMcpServers indicates the source contains MCP servers
	AssetTypeMcpServers AssetType = "mcp_servers"
)

// SourceProperties is an interface for accessing source configuration properties.
// It is used to avoid circular imports between catalog and mcp packages.
type SourceProperties interface {
	GetId() string
	GetProperties() map[string]any
}

// yamlAssetDetector represents a minimal YAML structure for detecting asset types
// Note: Uses json tags because k8s.io/apimachinery/pkg/util/yaml converts YAML to JSON first
type yamlAssetDetector struct {
	Models     []any `json:"models" yaml:"models"`
	McpServers []any `json:"mcp_servers" yaml:"mcp_servers"`
}

// DetectYamlAssetType reads a YAML file and detects what type of assets it contains.
// It returns an error if the file has multiple asset types.
func DetectYamlAssetType(props SourceProperties, reldir string) (AssetType, error) {
	path, exists := props.GetProperties()[YamlCatalogPathKey].(string)
	if !exists || path == "" {
		return "", fmt.Errorf("missing %s string property", YamlCatalogPathKey)
	}

	var fullPath string
	if filepath.IsAbs(path) {
		fullPath = path
	} else {
		fullPath = filepath.Join(reldir, path)
	}

	buf, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to read YAML file %s: %v", fullPath, err)
	}

	var detector yamlAssetDetector
	if err = yaml.Unmarshal(buf, &detector); err != nil {
		return "", fmt.Errorf("failed to parse YAML file %s: %v", fullPath, err)
	}

	hasModels := len(detector.Models) > 0
	hasMcpServers := len(detector.McpServers) > 0

	glog.V(2).Infof("YAML asset detection for %s: models=%v, mcp_servers=%v", fullPath, hasModels, hasMcpServers)

	if hasModels && hasMcpServers {
		return "", fmt.Errorf("YAML file %s contains multiple asset types (models and mcp_servers); each file should contain only one asset type", fullPath)
	}

	if hasModels {
		return AssetTypeModels, nil
	}

	if hasMcpServers {
		return AssetTypeMcpServers, nil
	}

	// Default to models for backwards compatibility if neither key is present
	glog.V(2).Infof("YAML file %s has no recognized asset type keys, defaulting to models", fullPath)
	return AssetTypeModels, nil
}

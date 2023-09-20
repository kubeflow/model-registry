package library

import (
	"fmt"
	"github.com/golang/glog"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:generate go-enum -type=PropertyType

type PropertyType int32

const (
	UNKNOWN PropertyType = iota
	INT
	DOUBLE
	STRING
	STRUCT
	PROTO
	BOOLEAN
)

type MetadataType struct {
	Name        *string                 `yaml:"name,omitempty"`
	Version     *string                 `yaml:"version,omitempty"`
	Description *string                 `yaml:"description,omitempty"`
	ExternalId  *string                 `yaml:"external_id,omitempty"`
	Properties  map[string]PropertyType `yaml:"properties,omitempty"`
}

type ArtifactType struct {
	MetadataType `yaml:",inline"`
	// TODO add support for base type enum
	//BaseType *ArtifactType_SystemDefinedBaseType `yaml:"base_type,omitempty"`
}

type ContextType struct {
	MetadataType `yaml:",inline"`
}

type ExecutionType struct {
	MetadataType `yaml:",inline"`
	//InputType  *ArtifactStructType                  `yaml:"input_type,omitempty"`
	//OutputType *ArtifactStructType                  `yaml:"output_type,omitempty"`
	//BaseType   *ExecutionType_SystemDefinedBaseType `yaml:"base_type,omitempty"`
}

type MetadataLibrary struct {
	ArtifactTypes  []ArtifactType  `yaml:"artifact-types,omitempty"`
	ContextTypes   []ContextType   `yaml:"context-types,omitempty"`
	ExecutionTypes []ExecutionType `yaml:"execution-types,omitempty"`
}

func LoadLibraries(dirs []string) (map[string]*MetadataLibrary, error) {
	result := make(map[string]*MetadataLibrary)
	for _, dir := range dirs {
		abs, err := filepath.Abs(dir)
		if err != nil {
			return nil, fmt.Errorf("error getting absolute library path for %s: %w", dir, err)
		}
		_, err = os.Stat(abs)
		if err != nil {
			return nil, fmt.Errorf("error opening library path for %s: %w", abs, err)
		}
		err = filepath.WalkDir(abs, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				glog.Warningf("error reading library path %s: %v", path, err)
				return filepath.SkipDir
			}
			if entry.IsDir() || !isYamlFile(path) {
				return nil
			}

			bytes, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read library file %s: %w", path, err)
			}
			lib := &MetadataLibrary{}
			err = yaml.Unmarshal(bytes, lib)
			if err != nil {
				return fmt.Errorf("failed to parse library file %s: %w", path, err)
			}
			result[path] = lib
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to read library directory %s: %w", abs, err)
		}
	}
	return result, nil
}

func isYamlFile(path string) bool {
	lowerPath := strings.ToLower(filepath.Ext(path))
	return strings.HasSuffix(lowerPath, ".yaml") || strings.HasSuffix(lowerPath, ".yml")
}

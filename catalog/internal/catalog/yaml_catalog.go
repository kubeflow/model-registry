package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/golang/glog"
	dbmodels "github.com/kubeflow/model-registry/catalog/internal/db/models"
	apimodels "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/internal/db/models"
)

const (
	yamlCatalogPathKey = "yamlCatalogPath"
	excludedModelsKey  = "excludedModels"
)

// convertMetadataValueToProperty converts a MetadataValue to a Properties object
// This helper eliminates code duplication when converting custom properties
func convertMetadataValueToProperty(key string, value apimodels.MetadataValue) models.Properties {
	// Handle different MetadataValue types
	if value.MetadataStringValue != nil {
		return models.NewStringProperty(key, value.MetadataStringValue.StringValue, true)
	} else if value.MetadataIntValue != nil {
		// MetadataIntValue.IntValue is a string, need to convert to int32
		if intVal, err := strconv.ParseInt(value.MetadataIntValue.IntValue, 10, 32); err == nil {
			return models.NewIntProperty(key, int32(intVal), true)
		} else {
			// If parsing fails, store as string
			return models.NewStringProperty(key, value.MetadataIntValue.IntValue, true)
		}
	} else if value.MetadataDoubleValue != nil {
		return models.NewDoubleProperty(key, value.MetadataDoubleValue.DoubleValue, true)
	} else if value.MetadataBoolValue != nil {
		return models.NewBoolProperty(key, value.MetadataBoolValue.BoolValue, true)
	} else {
		// For complex types, serialize to JSON
		if jsonBytes, err := json.Marshal(value); err == nil {
			return models.NewStringProperty(key, string(jsonBytes), true)
		}
		// Fallback to empty string if JSON marshaling fails
		return models.NewStringProperty(key, "", true)
	}
}

// convertCustomProperties converts a map of custom properties to a slice of Properties
func convertCustomProperties(customProps *map[string]apimodels.MetadataValue) []models.Properties {
	if customProps == nil {
		return nil
	}

	var properties []models.Properties
	for key, value := range *customProps {
		properties = append(properties, convertMetadataValueToProperty(key, value))
	}
	return properties
}

func init() {
	if err := RegisterModelProvider("yaml", newYamlModelProvider); err != nil {
		panic(err)
	}
}

type yamlModel struct {
	apimodels.CatalogModel `yaml:",inline"`
	Artifacts              []*yamlArtifact `yaml:"artifacts"`
}

type yamlArtifact struct {
	apimodels.CatalogArtifact
}

// convertModelAttributes converts basic model attributes and timestamps
func (ym *yamlModel) convertModelAttributes() *dbmodels.CatalogModelAttributes {
	attrs := &dbmodels.CatalogModelAttributes{
		Name: &ym.Name,
	}

	// Convert timestamps
	if ym.CreateTimeSinceEpoch != nil {
		if createTime, err := strconv.ParseInt(*ym.CreateTimeSinceEpoch, 10, 64); err == nil {
			attrs.CreateTimeSinceEpoch = &createTime
		}
	}

	if ym.LastUpdateTimeSinceEpoch != nil {
		if updateTime, err := strconv.ParseInt(*ym.LastUpdateTimeSinceEpoch, 10, 64); err == nil {
			attrs.LastUpdateTimeSinceEpoch = &updateTime
		}
	}

	return attrs
}

// convertModelProperties converts model properties to regular and custom properties
func (ym *yamlModel) convertModelProperties() ([]models.Properties, []models.Properties) {
	var properties []models.Properties
	var customProperties []models.Properties

	// Regular properties
	if ym.Description != nil {
		properties = append(properties, models.NewStringProperty("description", *ym.Description, false))
	}
	if ym.Readme != nil {
		properties = append(properties, models.NewStringProperty("readme", *ym.Readme, false))
	}
	if ym.Maturity != nil {
		properties = append(properties, models.NewStringProperty("maturity", *ym.Maturity, false))
	}
	if ym.Provider != nil {
		properties = append(properties, models.NewStringProperty("provider", *ym.Provider, false))
	}
	if ym.Logo != nil {
		properties = append(properties, models.NewStringProperty("logo", *ym.Logo, false))
	}
	if ym.License != nil {
		properties = append(properties, models.NewStringProperty("license", *ym.License, false))
	}
	if ym.LicenseLink != nil {
		properties = append(properties, models.NewStringProperty("license_link", *ym.LicenseLink, false))
	}
	if ym.LibraryName != nil {
		properties = append(properties, models.NewStringProperty("library_name", *ym.LibraryName, false))
	}
	if ym.SourceId != nil {
		properties = append(properties, models.NewStringProperty("source_id", *ym.SourceId, false))
	}

	// Convert array properties as struct properties
	if ym.Language == nil {
		ym.Language = []string{}
	}
	if languageJSON, err := json.Marshal(ym.Language); err == nil {
		properties = append(properties, models.NewStringProperty("language", string(languageJSON), false))
	}

	if ym.Tasks == nil {
		ym.Tasks = []string{}
	}
	if tasksJSON, err := json.Marshal(ym.Tasks); err == nil {
		properties = append(properties, models.NewStringProperty("tasks", string(tasksJSON), false))
	}

	// Convert custom properties from the YAML model
	if customProps := convertCustomProperties(&ym.CustomProperties); customProps != nil {
		customProperties = append(customProperties, customProps...)
	}

	return properties, customProperties
}

// convertModelArtifact converts a CatalogModelArtifact to database format
func convertModelArtifact(artifact *apimodels.CatalogModelArtifact) *dbmodels.CatalogArtifact {
	modelArtifact := &dbmodels.CatalogModelArtifactImpl{}

	// Set basic attributes
	attrs := &dbmodels.CatalogModelArtifactAttributes{
		URI: &artifact.Uri,
	}

	// Convert timestamps
	if artifact.CreateTimeSinceEpoch != nil {
		if createTime, err := strconv.ParseInt(*artifact.CreateTimeSinceEpoch, 10, 64); err == nil {
			attrs.CreateTimeSinceEpoch = &createTime
		}
	}
	if artifact.LastUpdateTimeSinceEpoch != nil {
		if updateTime, err := strconv.ParseInt(*artifact.LastUpdateTimeSinceEpoch, 10, 64); err == nil {
			attrs.LastUpdateTimeSinceEpoch = &updateTime
		}
	}

	modelArtifact.Attributes = attrs

	var artifactProperties []models.Properties
	artifactProperties = append(artifactProperties, models.NewStringProperty("uri", artifact.Uri, false))

	// Convert custom properties using helper function
	if customProps := convertCustomProperties(&artifact.CustomProperties); customProps != nil {
		modelArtifact.CustomProperties = &customProps
	}

	modelArtifact.Properties = &artifactProperties

	return &dbmodels.CatalogArtifact{
		CatalogModelArtifact: modelArtifact,
	}
}

// convertMetricsArtifact converts a CatalogMetricsArtifact to database format
func convertMetricsArtifact(artifact *apimodels.CatalogMetricsArtifact) *dbmodels.CatalogArtifact {
	metricsArtifact := &dbmodels.CatalogMetricsArtifactImpl{}

	// Set basic attributes
	attrs := &dbmodels.CatalogMetricsArtifactAttributes{
		MetricsType: dbmodels.MetricsType(artifact.MetricsType),
	}

	// Convert timestamps
	if artifact.CreateTimeSinceEpoch != nil {
		if createTime, err := strconv.ParseInt(*artifact.CreateTimeSinceEpoch, 10, 64); err == nil {
			attrs.CreateTimeSinceEpoch = &createTime
		}
	}
	if artifact.LastUpdateTimeSinceEpoch != nil {
		if updateTime, err := strconv.ParseInt(*artifact.LastUpdateTimeSinceEpoch, 10, 64); err == nil {
			attrs.LastUpdateTimeSinceEpoch = &updateTime
		}
	}

	metricsArtifact.Attributes = attrs

	// Handle properties
	var artifactProperties []models.Properties
	artifactProperties = append(artifactProperties, models.NewStringProperty("metricsType", artifact.MetricsType, false))

	// Convert custom properties using helper function
	if customProps := convertCustomProperties(&artifact.CustomProperties); customProps != nil {
		metricsArtifact.CustomProperties = &customProps
	}

	metricsArtifact.Properties = &artifactProperties

	return &dbmodels.CatalogArtifact{
		CatalogMetricsArtifact: metricsArtifact,
	}
}

func (ym *yamlModel) ToModelProviderRecord() ModelProviderRecord {
	model := dbmodels.CatalogModelImpl{}
	artifacts := make([]dbmodels.CatalogArtifact, len(ym.Artifacts))

	// Convert model attributes
	model.Attributes = ym.convertModelAttributes()

	// Convert model properties
	properties, customProperties := ym.convertModelProperties()
	if len(properties) > 0 {
		model.Properties = &properties
	}
	if len(customProperties) > 0 {
		model.CustomProperties = &customProperties
	}

	// Convert artifacts
	for j := range ym.Artifacts {
		if ym.Artifacts[j].CatalogModelArtifact != nil {
			artifacts[j] = *convertModelArtifact(ym.Artifacts[j].CatalogModelArtifact)
		} else if ym.Artifacts[j].CatalogMetricsArtifact != nil {
			artifacts[j] = *convertMetricsArtifact(ym.Artifacts[j].CatalogMetricsArtifact)
		}
	}

	return ModelProviderRecord{
		Model:     &model,
		Artifacts: artifacts,
	}
}

func (a *yamlArtifact) UnmarshalJSON(buf []byte) error {
	// This is very similar to generated code to unmarshal a
	// CatalogArtifact, but this version properly handles artifacts without
	// an artifactType, which is important for backwards compatibility.
	var yat struct {
		ArtifactType string `json:"artifactType"`
	}

	err := json.Unmarshal(buf, &yat)
	if err != nil {
		return err
	}

	switch yat.ArtifactType {
	case "model-artifact", "":
		err = json.Unmarshal(buf, &a.CatalogArtifact.CatalogModelArtifact)
		if a.CatalogArtifact.CatalogModelArtifact != nil {
			// Ensure artifactType is set even if it wasn't initially.
			a.CatalogArtifact.CatalogModelArtifact.ArtifactType = "model-artifact"
		}
	case "metrics-artifact":
		err = json.Unmarshal(buf, &a.CatalogArtifact.CatalogMetricsArtifact)
	default:
		return fmt.Errorf("unknown artifactType: %s", yat.ArtifactType)
	}

	return err
}

type yamlCatalog struct {
	Source string      `yaml:"source"`
	Models []yamlModel `yaml:"models"`
}

func isModelExcluded(modelName string, patterns []string) bool {
	for _, pattern := range patterns {
		if strings.HasSuffix(pattern, "*") {
			if strings.HasPrefix(modelName, strings.TrimSuffix(pattern, "*")) {
				return true
			}
		} else if modelName == pattern {
			return true
		}
	}
	return false
}

type yamlModelProvider struct {
	path           string
	excludedModels map[string]struct{}
}

func (p *yamlModelProvider) Models(ctx context.Context) (<-chan ModelProviderRecord, error) {
	// read the catalog and report errors
	catalog, err := p.read()
	if err != nil {
		return nil, err
	}

	ch := make(chan ModelProviderRecord)
	go func() {
		defer close(ch)

		// Send the initial list right away.
		p.emit(ctx, catalog, ch)

		// Watch for changes
		changes, err := getMonitor().Path(ctx, p.path)
		if err != nil {
			// Not fatal, we still have the inital load, but there
			// won't be any updates.
			glog.Errorf("unable to watch YAML catalog file: %v", err)
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-changes:
				glog.Infof("Reloading YAML catalog %s", p.path)

				catalog, err = p.read()
				if err != nil {
					glog.Errorf("unable to load YAML catalog: %v", err)
					continue
				}

				p.emit(ctx, catalog, ch)
			}
		}
	}()

	return ch, nil
}

func (p *yamlModelProvider) read() (*yamlCatalog, error) {
	buf, err := os.ReadFile(p.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s file: %v", yamlCatalogPathKey, err)
	}

	var catalog yamlCatalog
	if err = yaml.UnmarshalStrict(buf, &catalog); err != nil {
		return nil, fmt.Errorf("failed to parse %s file: %v", yamlCatalogPathKey, err)
	}

	return &catalog, nil
}

func (p *yamlModelProvider) emit(ctx context.Context, catalog *yamlCatalog, out chan<- ModelProviderRecord) {
	done := ctx.Done()
	for _, model := range catalog.Models {
		if _, excluded := p.excludedModels[model.Name]; excluded {
			continue
		}

		select {
		case out <- model.ToModelProviderRecord():
		case <-done:
			return
		}
	}
}

func newYamlModelProvider(ctx context.Context, source *Source, reldir string) (<-chan ModelProviderRecord, error) {
	p := &yamlModelProvider{}

	path, exists := source.Properties[yamlCatalogPathKey].(string)
	if !exists || path == "" {
		return nil, fmt.Errorf("missing %s string property", yamlCatalogPathKey)
	}

	if filepath.IsAbs(path) {
		p.path = path
	} else {
		p.path = filepath.Join(reldir, path)
	}

	// Excluded models is an optional source property.
	if _, exists := source.Properties[excludedModelsKey]; exists {
		excludedModels, ok := source.Properties[excludedModelsKey].([]any)
		if !ok {
			return nil, fmt.Errorf("%q property should be a list", excludedModelsKey)
		}

		p.excludedModels = make(map[string]struct{}, len(excludedModels))
		for i, name := range excludedModels {
			nameStr, ok := name.(string)
			if !ok {
				return nil, fmt.Errorf("%s: invalid list: index %d: wanted string, got %T", name, i, name)
			}
			p.excludedModels[nameStr] = struct{}{}
		}
	}

	return p.Models(ctx)
}

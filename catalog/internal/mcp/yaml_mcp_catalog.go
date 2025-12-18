package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/catalog/internal/common"
	dbmodels "github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/models"
	"k8s.io/apimachinery/pkg/util/yaml"
)

const (
	yamlCatalogPathKey = "yamlCatalogPath"
)

// McpServerProviderRecord contains one MCP server and its associated tools.
type McpServerProviderRecord struct {
	Server dbmodels.McpServer
	Tools  []dbmodels.McpServerTool
}

// McpServerProviderFunc emits MCP servers and related data in the channel it returns.
type McpServerProviderFunc func(ctx context.Context, source *McpSource, reldir string) (<-chan McpServerProviderRecord, error)

// McpSource represents a catalog source for MCP servers.
type McpSource struct {
	Id         string         `json:"id"`
	Name       string         `json:"name"`
	Type       string         `json:"type"`
	Enabled    *bool          `json:"enabled,omitempty"`
	Labels     []string       `json:"labels,omitempty"`
	Properties map[string]any `json:"properties,omitempty"`
	Origin     string         `json:"-" yaml:"-"`

	// IncludedServers is an optional list of glob patterns for MCP servers to include.
	// If specified, only servers matching at least one pattern will be included.
	// Pattern syntax: Only '*' wildcard is supported, patterns are case-insensitive.
	IncludedServers []string `json:"includedServers,omitempty" yaml:"includedServers,omitempty"`

	// ExcludedServers is an optional list of glob patterns for MCP servers to exclude.
	// Servers matching any pattern will be excluded even if they match an includedServers pattern.
	// Exclusions take precedence over inclusions.
	ExcludedServers []string `json:"excludedServers,omitempty" yaml:"excludedServers,omitempty"`
}

// GetId returns the source ID. Implements catalog.SourceProperties interface.
func (s *McpSource) GetId() string {
	return s.Id
}

// GetProperties returns the source properties map. Implements catalog.SourceProperties interface.
func (s *McpSource) GetProperties() map[string]any {
	return s.Properties
}

// yamlMcpServer represents an MCP server in YAML format.
type yamlMcpServer struct {
	Name             string            `yaml:"name" json:"name"`
	Description      string            `yaml:"description,omitempty" json:"description,omitempty"`
	Logo             string            `yaml:"logo,omitempty" json:"logo,omitempty"`
	License          string            `yaml:"license,omitempty" json:"license,omitempty"`
	LicenseLink      string            `yaml:"license_link,omitempty" json:"license_link,omitempty"`
	Provider         string            `yaml:"provider,omitempty" json:"provider,omitempty"`
	Version          string            `yaml:"version,omitempty" json:"version,omitempty"`
	Transports       []string          `yaml:"transports,omitempty" json:"transports,omitempty"`
	DocumentationUrl string            `yaml:"documentationUrl,omitempty" json:"documentationUrl,omitempty"`
	RepositoryUrl    string            `yaml:"repositoryUrl,omitempty" json:"repositoryUrl,omitempty"`
	SourceCode       string            `yaml:"sourceCode,omitempty" json:"sourceCode,omitempty"`
	Readme           string            `yaml:"readme,omitempty" json:"readme,omitempty"`
	PublishedDate    string            `yaml:"publishedDate,omitempty" json:"publishedDate,omitempty"`
	Tools            []yamlMcpTool     `yaml:"tools,omitempty" json:"tools,omitempty"`
	Artifacts        []yamlMcpArtifact `yaml:"artifacts,omitempty" json:"artifacts,omitempty"`

	// Deployment mode: "local" (default) or "remote"
	DeploymentMode string `yaml:"deploymentMode,omitempty" json:"deploymentMode,omitempty"`

	// Endpoints for remote MCP servers (different URLs per transport)
	Endpoints *yamlMcpEndpoints `yaml:"endpoints,omitempty" json:"endpoints,omitempty"`

	// CustomProperties following Model Registry pattern:
	// - Tags are MetadataStringValue entries with empty string_value
	// - Security indicators are MetadataBoolValue entries
	CustomProperties map[string]yamlMetadataValue `yaml:"customProperties,omitempty" json:"customProperties,omitempty"`

	// Timestamps for database consistency
	CreateTimeSinceEpoch     string `yaml:"createTimeSinceEpoch,omitempty" json:"createTimeSinceEpoch,omitempty"`
	LastUpdateTimeSinceEpoch string `yaml:"lastUpdateTimeSinceEpoch,omitempty" json:"lastUpdateTimeSinceEpoch,omitempty"`
}

// yamlMcpTool represents an MCP tool in YAML format.
type yamlMcpTool struct {
	Name          string                 `yaml:"name" json:"name"`
	Description   string                 `yaml:"description" json:"description"`
	AccessType    string                 `yaml:"accessType" json:"accessType"`
	Parameters    []yamlMcpToolParameter `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	Revoked       bool                   `yaml:"revoked,omitempty" json:"revoked,omitempty"`
	RevokedReason string                 `yaml:"revokedReason,omitempty" json:"revokedReason,omitempty"`
}

// yamlMcpToolParameter represents a tool parameter in YAML format.
type yamlMcpToolParameter struct {
	Name        string `yaml:"name" json:"name"`
	Type        string `yaml:"type" json:"type"`
	Description string `yaml:"description" json:"description"`
	Required    bool   `yaml:"required" json:"required"`
}

// yamlMetadataValue represents a metadata value in YAML format.
// It can be a string value, bool value, int value, or double value.
// Following Model Registry patterns:
// - Tags are MetadataStringValue with empty string_value
// - Security indicators are MetadataBoolValue
type yamlMetadataValue struct {
	MetadataType string  `yaml:"metadataType" json:"metadataType"`
	StringValue  *string `yaml:"string_value,omitempty" json:"string_value,omitempty"`
	BoolValue    *bool   `yaml:"bool_value,omitempty" json:"bool_value,omitempty"`
	IntValue     *int64  `yaml:"int_value,omitempty" json:"int_value,omitempty"`
	DoubleValue  *float64 `yaml:"double_value,omitempty" json:"double_value,omitempty"`
}

// yamlMcpArtifact represents an MCP server artifact (e.g., OCI image) in YAML format.
// Simplified format matching model artifacts: just uri + timestamps.
type yamlMcpArtifact struct {
	Uri                      string `yaml:"uri" json:"uri"`
	CreateTimeSinceEpoch     string `yaml:"createTimeSinceEpoch,omitempty" json:"createTimeSinceEpoch,omitempty"`
	LastUpdateTimeSinceEpoch string `yaml:"lastUpdateTimeSinceEpoch,omitempty" json:"lastUpdateTimeSinceEpoch,omitempty"`
}

// yamlMcpEndpoints represents network endpoints for remote MCP servers.
type yamlMcpEndpoints struct {
	Http string `yaml:"http,omitempty" json:"http,omitempty"`
	Sse  string `yaml:"sse,omitempty" json:"sse,omitempty"`
}

// yamlMcpCatalog represents the YAML catalog structure.
type yamlMcpCatalog struct {
	Source     string          `yaml:"source" json:"source"`
	McpServers []yamlMcpServer `yaml:"mcp_servers" json:"mcp_servers"`
}

// yamlMcpProvider provides MCP servers from a YAML file.
type yamlMcpProvider struct {
	path string
}

// ToMcpServerProviderRecord converts a YAML server to a provider record.
func (ys *yamlMcpServer) ToMcpServerProviderRecord() McpServerProviderRecord {
	server := &dbmodels.McpServerImpl{}

	// Convert attributes
	attrs := &dbmodels.McpServerAttributes{
		Name: &ys.Name,
	}

	// Convert timestamps
	if ys.CreateTimeSinceEpoch != "" {
		if createTime, err := strconv.ParseInt(ys.CreateTimeSinceEpoch, 10, 64); err == nil {
			attrs.CreateTimeSinceEpoch = &createTime
		}
	}
	if ys.LastUpdateTimeSinceEpoch != "" {
		if updateTime, err := strconv.ParseInt(ys.LastUpdateTimeSinceEpoch, 10, 64); err == nil {
			attrs.LastUpdateTimeSinceEpoch = &updateTime
		}
	}

	server.Attributes = attrs

	// Convert properties
	var properties []models.Properties

	if ys.Description != "" {
		properties = append(properties, models.NewStringProperty("description", ys.Description, false))
	}
	if ys.Logo != "" {
		properties = append(properties, models.NewStringProperty("logo", ys.Logo, false))
	}
	if ys.License != "" {
		properties = append(properties, models.NewStringProperty("license", ys.License, false))
	}
	if ys.LicenseLink != "" {
		properties = append(properties, models.NewStringProperty("license_link", ys.LicenseLink, false))
	}
	if ys.Provider != "" {
		properties = append(properties, models.NewStringProperty("provider", ys.Provider, false))
	}
	if ys.Version != "" {
		properties = append(properties, models.NewStringProperty("version", ys.Version, false))
	}
	// Store transports as JSON array
	if len(ys.Transports) > 0 {
		if transportsJSON, err := json.Marshal(ys.Transports); err == nil {
			properties = append(properties, models.NewStringProperty("transports", string(transportsJSON), false))
		}
	}
	if ys.DocumentationUrl != "" {
		properties = append(properties, models.NewStringProperty("documentationUrl", ys.DocumentationUrl, false))
	}
	if ys.RepositoryUrl != "" {
		properties = append(properties, models.NewStringProperty("repositoryUrl", ys.RepositoryUrl, false))
	}
	if ys.SourceCode != "" {
		properties = append(properties, models.NewStringProperty("sourceCode", ys.SourceCode, false))
	}
	if ys.Readme != "" {
		properties = append(properties, models.NewStringProperty("readme", ys.Readme, false))
	}
	if ys.PublishedDate != "" {
		properties = append(properties, models.NewStringProperty("publishedDate", ys.PublishedDate, false))
	}

	// Convert customProperties from YAML format to database properties
	// Following Model Registry patterns:
	// - Tags are MetadataStringValue entries with empty string_value
	// - Security indicators are MetadataBoolValue entries
	if len(ys.CustomProperties) > 0 {
		// Extract tags (MetadataStringValue with empty string_value)
		var tags []string
		for key, val := range ys.CustomProperties {
			if val.MetadataType == "MetadataStringValue" && val.StringValue != nil && *val.StringValue == "" {
				tags = append(tags, key)
			}
		}
		if len(tags) > 0 {
			if tagsJSON, err := json.Marshal(tags); err == nil {
				properties = append(properties, models.NewStringProperty("tags", string(tagsJSON), false))
			}
		}

		// Extract security indicators (MetadataBoolValue entries)
		for key, val := range ys.CustomProperties {
			if val.MetadataType == "MetadataBoolValue" && val.BoolValue != nil {
				properties = append(properties, models.NewBoolProperty(key, *val.BoolValue, false))
			}
		}
	}

	// Convert artifacts as JSON for storage (for local MCP servers with OCI images)
	if len(ys.Artifacts) > 0 {
		if artifactsJSON, err := json.Marshal(ys.Artifacts); err == nil {
			properties = append(properties, models.NewStringProperty("artifacts", string(artifactsJSON), false))
		}
	}

	// Convert deployment mode (default to "local" if not specified)
	deploymentMode := ys.DeploymentMode
	if deploymentMode == "" {
		deploymentMode = "local"
	}
	properties = append(properties, models.NewStringProperty("deploymentMode", deploymentMode, false))

	// Convert endpoints for remote servers
	if ys.Endpoints != nil {
		if endpointsJSON, err := json.Marshal(ys.Endpoints); err == nil {
			properties = append(properties, models.NewStringProperty("endpoints", string(endpointsJSON), false))
		}
	}

	// Convert tools as JSON for storage
	if len(ys.Tools) > 0 {
		if toolsJSON, err := json.Marshal(ys.Tools); err == nil {
			properties = append(properties, models.NewStringProperty("tools", string(toolsJSON), false))
		}
	}

	if len(properties) > 0 {
		server.Properties = &properties
	}

	// Convert tools to separate entities (for future use when tools are stored separately)
	tools := make([]dbmodels.McpServerTool, 0, len(ys.Tools))
	for _, t := range ys.Tools {
		tool := &dbmodels.McpServerToolImpl{
			Attributes: &dbmodels.McpServerToolAttributes{
				Name: &t.Name,
			},
		}

		var toolProps []models.Properties
		if t.Description != "" {
			toolProps = append(toolProps, models.NewStringProperty("description", t.Description, false))
		}
		if t.AccessType != "" {
			toolProps = append(toolProps, models.NewStringProperty("accessType", t.AccessType, false))
		}
		if len(t.Parameters) > 0 {
			if paramsJSON, err := json.Marshal(t.Parameters); err == nil {
				toolProps = append(toolProps, models.NewStringProperty("parameters", string(paramsJSON), false))
			}
		}
		// Always store revoked status (defaults to false)
		toolProps = append(toolProps, models.NewBoolProperty("revoked", t.Revoked, false))
		if t.RevokedReason != "" {
			toolProps = append(toolProps, models.NewStringProperty("revokedReason", t.RevokedReason, false))
		}

		if len(toolProps) > 0 {
			tool.Properties = &toolProps
		}

		tools = append(tools, tool)
	}

	return McpServerProviderRecord{
		Server: server,
		Tools:  tools,
	}
}

// Models returns a channel of MCP server provider records.
func (p *yamlMcpProvider) Models(ctx context.Context) (<-chan McpServerProviderRecord, error) {
	catalog, err := p.read()
	if err != nil {
		return nil, err
	}

	ch := make(chan McpServerProviderRecord)
	go func() {
		defer close(ch)
		p.emit(ctx, catalog, ch)
	}()

	return ch, nil
}

// read reads the YAML catalog file.
func (p *yamlMcpProvider) read() (*yamlMcpCatalog, error) {
	buf, err := os.ReadFile(p.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s file: %v", yamlCatalogPathKey, err)
	}

	var catalog yamlMcpCatalog
	if err = yaml.UnmarshalStrict(buf, &catalog); err != nil {
		return nil, fmt.Errorf("failed to parse %s file: %v", yamlCatalogPathKey, err)
	}

	return &catalog, nil
}

// emit sends MCP server records to the output channel.
func (p *yamlMcpProvider) emit(ctx context.Context, catalog *yamlMcpCatalog, out chan<- McpServerProviderRecord) {
	done := ctx.Done()
	for _, server := range catalog.McpServers {
		select {
		case out <- server.ToMcpServerProviderRecord():
		case <-done:
			return
		}
	}

	// Send an empty record to indicate that we're done with the batch.
	select {
	case out <- McpServerProviderRecord{}:
	case <-done:
	}
}

// NewYamlMcpProvider creates a new YAML MCP provider.
// It detects the asset type from the YAML content and only processes files containing mcp_servers.
func NewYamlMcpProvider(ctx context.Context, source *McpSource, reldir string) (<-chan McpServerProviderRecord, error) {
	// First, detect the asset type from the YAML content
	assetType, err := common.DetectYamlAssetType(source, reldir)
	if err != nil {
		return nil, err
	}

	// Only process this source if it contains MCP servers
	if assetType != common.AssetTypeMcpServers {
		glog.V(2).Infof("Skipping source %s in MCP provider: detected asset type is %s", source.Id, assetType)
		// Return an empty channel that closes immediately
		ch := make(chan McpServerProviderRecord)
		close(ch)
		return ch, nil
	}

	p := &yamlMcpProvider{}

	path, exists := source.Properties[yamlCatalogPathKey].(string)
	if !exists || path == "" {
		return nil, fmt.Errorf("missing %s string property", yamlCatalogPathKey)
	}

	if filepath.IsAbs(path) {
		p.path = path
	} else {
		p.path = filepath.Join(reldir, path)
	}

	glog.Infof("Loading MCP servers from YAML file: %s", p.path)

	return p.Models(ctx)
}

// RegisteredMcpProviders holds the registered MCP provider functions.
var RegisteredMcpProviders = map[string]McpServerProviderFunc{
	"yaml": NewYamlMcpProvider,
}

// RegisterMcpProvider registers an MCP provider function.
func RegisterMcpProvider(name string, callback McpServerProviderFunc) error {
	if _, exists := RegisteredMcpProviders[name]; exists {
		return fmt.Errorf("MCP provider type %s already exists", name)
	}
	RegisteredMcpProviders[name] = callback
	return nil
}

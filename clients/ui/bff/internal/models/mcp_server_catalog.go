package models

import "github.com/kubeflow/model-registry/pkg/openapi"

type McpDeploymentMode string

const (
	McpDeploymentModeLocal  McpDeploymentMode = "local"
	McpDeploymentModeRemote McpDeploymentMode = "remote"
)

type McpTransportType string

const (
	McpTransportTypeStdio McpTransportType = "stdio"
	McpTransportTypeSSE   McpTransportType = "sse"
	McpTransportTypeHTTP  McpTransportType = "http"
)

type McpToolAccessType string

const (
	McpToolAccessTypeReadOnly  McpToolAccessType = "read_only"
	McpToolAccessTypeReadWrite McpToolAccessType = "read_write"
	McpToolAccessTypeExecute   McpToolAccessType = "execute"
)

type McpEndpoints struct {
	HTTP *string `json:"http,omitempty"`
	SSE  *string `json:"sse,omitempty"`
}

type McpArtifact struct {
	URI                      string  `json:"uri"`
	CreateTimeSinceEpoch     *string `json:"createTimeSinceEpoch,omitempty"`
	LastUpdateTimeSinceEpoch *string `json:"lastUpdateTimeSinceEpoch,omitempty"`
}

type McpSecurityIndicator struct {
	VerifiedSource *bool `json:"verifiedSource,omitempty"`
	SecureEndpoint *bool `json:"secureEndpoint,omitempty"`
	SAST           *bool `json:"sast,omitempty"`
	ReadOnlyTools  *bool `json:"readOnlyTools,omitempty"`
}

type McpToolParameter struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Description *string `json:"description,omitempty"`
	Required    bool    `json:"required"`
}

type McpTool struct {
	Name             string                            `json:"name"`
	Description      *string                           `json:"description,omitempty"`
	AccessType       McpToolAccessType                 `json:"accessType"`
	Parameters       []McpToolParameter                `json:"parameters,omitempty"`
	Revoked          *bool                             `json:"revoked,omitempty"`
	RevokedReason    *string                           `json:"revokedReason,omitempty"`
	CustomProperties *map[string]openapi.MetadataValue `json:"customProperties,omitempty"`
}

type McpServer struct {
	ID                 int                               `json:"id"`
	Name               string                            `json:"name"`
	SourceID           *string                           `json:"source_id,omitempty"`
	Description        *string                           `json:"description,omitempty"`
	Logo               *string                           `json:"logo,omitempty"`
	License            *string                           `json:"license,omitempty"`
	LicenseLink        *string                           `json:"licenseLink,omitempty"`
	Provider           *string                           `json:"provider,omitempty"`
	Version            *string                           `json:"version,omitempty"`
	Tags               []string                          `json:"tags,omitempty"`
	ToolCount          int                               `json:"toolCount"`
	Tools              []McpTool                         `json:"tools,omitempty"`
	SecurityIndicators *McpSecurityIndicator             `json:"securityIndicators,omitempty"`
	DocumentationURL   *string                           `json:"documentationUrl,omitempty"`
	RepositoryURL      *string                           `json:"repositoryUrl,omitempty"`
	SourceCode         *string                           `json:"sourceCode,omitempty"`
	LastUpdated        *string                           `json:"lastUpdated,omitempty"`
	PublishedDate      *string                           `json:"publishedDate,omitempty"`
	Artifacts          []McpArtifact                     `json:"artifacts,omitempty"`
	Transports         []McpTransportType                `json:"transports,omitempty"`
	Readme             *string                           `json:"readme,omitempty"`
	DeploymentMode     *McpDeploymentMode                `json:"deploymentMode,omitempty"`
	Endpoints          *McpEndpoints                     `json:"endpoints,omitempty"`
	CustomProperties   *map[string]openapi.MetadataValue `json:"customProperties,omitempty"`
}

type McpServerList struct {
	NextPageToken string      `json:"nextPageToken"`
	PageSize      int32       `json:"pageSize"`
	Size          int32       `json:"size"`
	Items         []McpServer `json:"items"`
}

type McpToolWithServer struct {
	ServerID   string  `json:"serverId"`
	ServerName string  `json:"serverName"`
	Tool       McpTool `json:"tool"`
}

type McpToolList struct {
	NextPageToken string              `json:"nextPageToken"`
	PageSize      int32               `json:"pageSize"`
	Size          int32               `json:"size"`
	Items         []McpToolWithServer `json:"items"`
}

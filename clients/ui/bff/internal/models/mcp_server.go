package models

// McpToolAccessType represents the access type for an MCP tool
type McpToolAccessType string

const (
	McpToolAccessTypeReadOnly  McpToolAccessType = "read_only"
	McpToolAccessTypeReadWrite McpToolAccessType = "read_write"
	McpToolAccessTypeExecute   McpToolAccessType = "execute"
)

// McpTransportType represents the transport protocol for an MCP server
type McpTransportType string

const (
	McpTransportTypeStdio McpTransportType = "stdio"
	McpTransportTypeSSE   McpTransportType = "sse"
	McpTransportTypeHTTP  McpTransportType = "http"
)

// McpDeploymentMode represents the deployment mode for an MCP server
type McpDeploymentMode string

const (
	McpDeploymentModeLocal  McpDeploymentMode = "local"
	McpDeploymentModeRemote McpDeploymentMode = "remote"
)

// McpEndpoints represents network endpoints for remote MCP servers
type McpEndpoints struct {
	Http *string `json:"http,omitempty"`
	Sse  *string `json:"sse,omitempty"`
}

// McpArtifact represents an artifact for an MCP server (e.g., OCI image)
type McpArtifact struct {
	Uri                      string  `json:"uri"`
	CreateTimeSinceEpoch     *string `json:"createTimeSinceEpoch,omitempty"`
	LastUpdateTimeSinceEpoch *string `json:"lastUpdateTimeSinceEpoch,omitempty"`
}

// McpSecurityIndicator represents security indicators for an MCP server
type McpSecurityIndicator struct {
	VerifiedSource bool `json:"verifiedSource"`
	SecureEndpoint bool `json:"secureEndpoint"`
	Sast           bool `json:"sast"`
	ReadOnlyTools  bool `json:"readOnlyTools"`
}

// McpToolParameter represents a parameter for an MCP tool
type McpToolParameter struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

// McpMetadataStringValue represents a string custom property value
type McpMetadataStringValue struct {
	StringValue  string `json:"string_value"`
	MetadataType string `json:"metadataType"`
}

// McpMetadataBoolValue represents a boolean custom property value
type McpMetadataBoolValue struct {
	BoolValue    bool   `json:"bool_value"`
	MetadataType string `json:"metadataType"`
}

// McpCustomProperties represents custom properties following Model Registry patterns.
// Tags are stored as MetadataStringValue entries with empty string_value (label pattern).
// Security indicators are stored as MetadataBoolValue entries.
type McpCustomProperties map[string]interface{}

// McpTool represents a tool exposed by an MCP server
type McpTool struct {
	Name             string               `json:"name"`
	Description      string               `json:"description"`
	AccessType       McpToolAccessType    `json:"accessType"`
	Parameters       []McpToolParameter   `json:"parameters,omitempty"`
	Revoked          *bool                `json:"revoked,omitempty"`
	RevokedReason    *string              `json:"revokedReason,omitempty"`
	CustomProperties *McpCustomProperties `json:"customProperties,omitempty"`
}

// McpServer represents an MCP server in the catalog
type McpServer struct {
	ID                 string                `json:"id"`
	Name               string                `json:"name"`
	Description        string                `json:"description"`
	SourceId           *string               `json:"source_id,omitempty"`
	Logo               *string               `json:"logo,omitempty"`
	License            *string               `json:"license,omitempty"`
	LicenseLink        *string               `json:"license_link,omitempty"`
	Provider           *string               `json:"provider,omitempty"`
	Version            *string               `json:"version,omitempty"`
	Tags               []string              `json:"tags,omitempty"`
	Tools              []McpTool             `json:"tools,omitempty"`
	SecurityIndicators *McpSecurityIndicator `json:"securityIndicators,omitempty"`
	DocumentationUrl   *string               `json:"documentationUrl,omitempty"`
	RepositoryUrl      *string               `json:"repositoryUrl,omitempty"`
	SourceCode         *string               `json:"sourceCode,omitempty"`
	LastUpdated        *string               `json:"lastUpdated,omitempty"`
	PublishedDate      *string               `json:"publishedDate,omitempty"`
	Artifacts          []McpArtifact         `json:"artifacts,omitempty"`
	Transports         []McpTransportType    `json:"transports,omitempty"`
	Readme             *string               `json:"readme,omitempty"`
	DeploymentMode     *McpDeploymentMode    `json:"deploymentMode,omitempty"`
	Endpoints          *McpEndpoints         `json:"endpoints,omitempty"`
	CustomProperties   *McpCustomProperties  `json:"customProperties,omitempty"`
}

// McpServerList represents a paginated list of MCP servers
type McpServerList struct {
	NextPageToken string      `json:"nextPageToken"`
	PageSize      int32       `json:"pageSize"`
	Size          int32       `json:"size"`
	Items         []McpServer `json:"items"`
}

// McpCatalogSourceStatus represents the status of an MCP catalog source
type McpCatalogSourceStatus string

const (
	McpCatalogSourceStatusAvailable McpCatalogSourceStatus = "available"
	McpCatalogSourceStatusError     McpCatalogSourceStatus = "error"
	McpCatalogSourceStatusDisabled  McpCatalogSourceStatus = "disabled"
)

// CatalogAssetType represents the type of assets in a catalog source
type CatalogAssetType string

const (
	CatalogAssetTypeModels     CatalogAssetType = "models"
	CatalogAssetTypeMcpServers CatalogAssetType = "mcp_servers"
)

// McpCatalogSource represents a source of MCP servers in the catalog
type McpCatalogSource struct {
	ID        string                  `json:"id"`
	Name      string                  `json:"name"`
	Labels    []string                `json:"labels"`
	Enabled   *bool                   `json:"enabled,omitempty"`
	AssetType *CatalogAssetType       `json:"assetType,omitempty"`
	Status    *McpCatalogSourceStatus `json:"status,omitempty"`
	Error     *string                 `json:"error,omitempty"`
}

// McpCatalogSourceList represents a paginated list of MCP catalog sources
type McpCatalogSourceList struct {
	NextPageToken string             `json:"nextPageToken"`
	PageSize      int32              `json:"pageSize"`
	Size          int32              `json:"size"`
	Items         []McpCatalogSource `json:"items"`
}

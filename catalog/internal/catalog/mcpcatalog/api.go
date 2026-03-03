package mcpcatalog

import (
	"context"

	openapi "github.com/kubeflow/model-registry/catalog/pkg/openapi"
)

type ListMCPServersParams struct {
	Name, Query, FilterQuery, NamedQuery string
	IncludeTools                         bool
	ToolLimit                            int32
	PageSize                             int32
	OrderBy                              openapi.OrderByField
	SortOrder                            openapi.SortOrder
	NextPageToken                        *string
}

type ListMCPServerToolsParams struct {
	FilterQuery   string
	PageSize      int32
	OrderBy       openapi.OrderByField
	SortOrder     openapi.SortOrder
	NextPageToken *string
}

// MCPCatalogProvider implements the MCP catalog API endpoints.
type MCPCatalogProvider interface {
	// ListMCPServers returns all MCP servers according to the parameters.
	// If nothing suitable is found, it returns an empty list.
	ListMCPServers(ctx context.Context, params ListMCPServersParams) (openapi.MCPServerList, error)

	// GetMCPServer returns MCP server metadata for a single server by its ID.
	// If includeTools is true, the server's tools are also included.
	// If nothing is found with the ID provided it returns nil, without an error.
	GetMCPServer(ctx context.Context, serverID string, includeTools bool) (*openapi.MCPServer, error)

	// ListMCPServerTools returns all tools for a specific MCP server.
	// If no server is found with that ID, it returns nil. If the server is
	// found, but has no tools, an empty list is returned.
	ListMCPServerTools(ctx context.Context, serverID string, params ListMCPServerToolsParams) (openapi.MCPToolsList, error)

	// GetMCPServerTool returns a specific tool by server ID and tool name.
	// If nothing is found it returns nil, without an error.
	GetMCPServerTool(ctx context.Context, serverID string, toolName string) (*openapi.MCPTool, error)
}

package repositories

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/httpclient"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

const mcpServersPath = "/mcp_servers"
const mcpFilterOptionsPath = "/mcp_servers/filter_options"

// McpServersInterface defines operations for MCP server catalog access via catalog service
type McpServersInterface interface {
	GetAllMcpServers(client httpclient.HTTPClientInterface, pageValues url.Values) (*models.McpServerList, error)
	GetMcpServer(client httpclient.HTTPClientInterface, serverId string) (*models.McpServer, error)
	GetMcpFilterOptions(client httpclient.HTTPClientInterface) (*models.FilterOptionsList, error)
	GetAllMcpSources(client httpclient.HTTPClientInterface, pageValues url.Values) (*models.McpCatalogSourceList, error)
}

// McpServers implements McpServersInterface
type McpServers struct {
	McpServersInterface
}

// GetAllMcpServers fetches all MCP servers from the catalog service
func (m McpServers) GetAllMcpServers(client httpclient.HTTPClientInterface, pageValues url.Values) (*models.McpServerList, error) {
	responseData, err := client.GET(UrlWithPageParams(mcpServersPath, pageValues))
	if err != nil {
		return nil, fmt.Errorf("error fetching MCP servers: %w", err)
	}

	var serverList models.McpServerList

	if err := json.Unmarshal(responseData, &serverList); err != nil {
		return nil, fmt.Errorf("error decoding MCP servers response: %w", err)
	}

	return &serverList, nil
}

// GetMcpServer fetches a specific MCP server by ID from the catalog service
func (m McpServers) GetMcpServer(client httpclient.HTTPClientInterface, serverId string) (*models.McpServer, error) {
	path, err := url.JoinPath(mcpServersPath, serverId)
	if err != nil {
		return nil, err
	}

	responseData, err := client.GET(path)
	if err != nil {
		return nil, fmt.Errorf("error fetching MCP server: %w", err)
	}

	var server models.McpServer

	if err := json.Unmarshal(responseData, &server); err != nil {
		return nil, fmt.Errorf("error decoding MCP server response: %w", err)
	}

	return &server, nil
}

// GetMcpFilterOptions fetches filter options for MCP servers from the catalog service
func (m McpServers) GetMcpFilterOptions(client httpclient.HTTPClientInterface) (*models.FilterOptionsList, error) {
	responseData, err := client.GET(mcpFilterOptionsPath)
	if err != nil {
		return nil, fmt.Errorf("error fetching MCP filter options: %w", err)
	}

	var filterOptions models.FilterOptionsList

	if err := json.Unmarshal(responseData, &filterOptions); err != nil {
		return nil, fmt.Errorf("error decoding MCP filter options response: %w", err)
	}

	return &filterOptions, nil
}

// GetAllMcpSources fetches all MCP catalog sources from the catalog service
// It uses the unified /sources endpoint with assetType=mcp_servers filter
func (m McpServers) GetAllMcpSources(client httpclient.HTTPClientInterface, pageValues url.Values) (*models.McpCatalogSourceList, error) {
	// Add assetType filter for MCP servers
	if pageValues == nil {
		pageValues = url.Values{}
	}
	pageValues.Set("assetType", "mcp_servers")

	responseData, err := client.GET(UrlWithPageParams(sourcesPath, pageValues))
	if err != nil {
		return nil, fmt.Errorf("error fetching MCP sources: %w", err)
	}

	var sourceList models.McpCatalogSourceList

	if err := json.Unmarshal(responseData, &sourceList); err != nil {
		return nil, fmt.Errorf("error decoding MCP sources response: %w", err)
	}

	return &sourceList, nil
}

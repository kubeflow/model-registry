package repositories

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/httpclient"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

const mcpServerPath = "/mcp_servers"
const mcpFilterOptionPath = "/mcp_servers/filter_options"

type McpServerCatalogInterface interface {
	GetAllMcpServers(client httpclient.HTTPClientInterface, pageValues url.Values) (*models.McpServerList, error)
	GetMcpServersFilter(client httpclient.HTTPClientInterface) (*models.FilterOptionsList, error)
	GetMcpServer(client httpclient.HTTPClientInterface, serverId string) (*models.McpServer, error)
	GetMcpServersTools(client httpclient.HTTPClientInterface, serverId string) (*models.McpToolList, error)
}

type McpServerCatalog struct {
	McpServerCatalogInterface
}

func (a *CatalogSources) GetAllMcpServers(client httpclient.HTTPClientInterface, pageValues url.Values) (*models.McpServerList, error) {
	responseData, err := client.GET(UrlWithPageParams(mcpServerPath, pageValues))

	if err != nil {
		return nil, fmt.Errorf("error fetching mcpServerPath: %w", err)
	}

	var mcpServers models.McpServerList

	if err := json.Unmarshal(responseData, &mcpServers); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &mcpServers, nil

}

func (a *CatalogSources) GetMcpServersFilter(client httpclient.HTTPClientInterface) (*models.FilterOptionsList, error) {
	responseData, err := client.GET(mcpFilterOptionPath)

	if err != nil {
		return nil, fmt.Errorf("error fetching mcpFilterOptionPath: %w", err)
	}

	var filters models.FilterOptionsList

	if err := json.Unmarshal(responseData, &filters); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}
	return &filters, nil
}

func (a *CatalogSources) GetMcpServer(client httpclient.HTTPClientInterface, serverId string) (*models.McpServer, error) {
	path, err := url.JoinPath(mcpServerPath, serverId)

	if err != nil {
		return nil, err
	}

	responseData, err := client.GET(path)

	if err != nil {
		return nil, fmt.Errorf("error fetching mcpServerPath: %w", err)
	}

	var mcpServer models.McpServer

	if err := json.Unmarshal(responseData, &mcpServer); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &mcpServer, nil
}

func (a *CatalogSources) GetMcpServersTools(client httpclient.HTTPClientInterface, serverId string) (*models.McpToolList, error) {
	path, err := url.JoinPath(mcpServerPath, serverId, "tools")

	if err != nil {
		return nil, err
	}

	responseData, err := client.GET(path)

	if err != nil {
		return nil, fmt.Errorf("error fetching mcpServerPath: %w", err)
	}

	var mcpServerTools models.McpToolList

	if err := json.Unmarshal(responseData, &mcpServerTools); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &mcpServerTools, nil
}

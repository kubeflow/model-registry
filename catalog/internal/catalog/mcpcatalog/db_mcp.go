package mcpcatalog

import (
	"context"
	"fmt"

	"github.com/kubeflow/model-registry/catalog/internal/converter"
	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	dbmodels "github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	openapi "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/pkg/api"
)

type dbMCPCatalogImpl struct {
	mcpServerRepo     dbmodels.MCPServerRepository
	mcpServerToolRepo dbmodels.MCPServerToolRepository
}

func NewDBMCPCatalog(services service.Services) MCPCatalogProvider {
	return &dbMCPCatalogImpl{
		mcpServerRepo:     services.MCPServerRepository,
		mcpServerToolRepo: services.MCPServerToolRepository,
	}
}

func (d *dbMCPCatalogImpl) ListMCPServers(ctx context.Context, params ListMCPServersParams) (openapi.MCPServerList, error) {
	// Build MCPServerListOptions from params
	listOptions := models.MCPServerListOptions{
		Query:       &params.Query,
		FilterQuery: &params.FilterQuery,
		NamedQuery:  &params.NamedQuery,
	}

	if params.Name != "" {
		listOptions.Name = &params.Name
	}

	// Set pagination
	listOptions.Pagination.PageSize = &params.PageSize
	listOptions.Pagination.OrderBy = (*string)(&params.OrderBy)
	listOptions.Pagination.SortOrder = (*string)(&params.SortOrder)
	if params.NextPageToken != nil {
		listOptions.Pagination.NextPageToken = params.NextPageToken
	}

	// Call repository
	serversList, err := d.mcpServerRepo.List(listOptions)
	if err != nil {
		return openapi.MCPServerList{}, err
	}

	// Convert to OpenAPI models
	apiServers := make([]openapi.MCPServer, 0) // Initialize as empty slice, not nil
	for _, dbServer := range serversList.Items {
		var tools []models.MCPServerTool

		// Optionally include tools
		if params.IncludeTools {
			toolOptions := models.MCPServerToolListOptions{
				ParentID: *dbServer.GetID(),
			}
			// Apply tool limit if specified
			if params.ToolLimit > 0 {
				toolOptions.Pagination.PageSize = &params.ToolLimit
			}

			tools, err = d.mcpServerToolRepo.List(toolOptions)
			if err != nil {
				return openapi.MCPServerList{}, fmt.Errorf("error loading tools for server %d: %w", *dbServer.GetID(), err)
			}
		}

		apiServer := converter.ConvertDbMCPServerWithToolsToOpenapi(dbServer, tools)
		apiServers = append(apiServers, *apiServer)
	}

	return openapi.MCPServerList{
		Items:         apiServers,
		Size:          int32(len(apiServers)),
		PageSize:      params.PageSize,
		NextPageToken: serversList.NextPageToken,
	}, nil
}

func (d *dbMCPCatalogImpl) GetMCPServer(ctx context.Context, serverID string, includeTools bool) (*openapi.MCPServer, error) {
	id, err := apiutils.ValidateIDAsInt32(serverID, "server")
	if err != nil {
		return nil, fmt.Errorf("invalid server ID '%s': %w", serverID, api.ErrBadRequest)
	}

	// Get server by ID
	dbServer, err := d.mcpServerRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("server not found with ID %s: %w", serverID, api.ErrNotFound)
	}

	var tools []models.MCPServerTool

	// Optionally include tools
	if includeTools {
		toolOptions := models.MCPServerToolListOptions{
			ParentID: *dbServer.GetID(),
		}
		tools, err = d.mcpServerToolRepo.List(toolOptions)
		if err != nil {
			return nil, fmt.Errorf("error loading tools for server %s: %w", serverID, err)
		}
	}

	apiServer := converter.ConvertDbMCPServerWithToolsToOpenapi(dbServer, tools)
	return apiServer, nil
}

func (d *dbMCPCatalogImpl) ListMCPServerTools(ctx context.Context, serverID string, params ListMCPServerToolsParams) (openapi.MCPToolsList, error) {
	id, err := apiutils.ValidateIDAsInt32(serverID, "server")
	if err != nil {
		return openapi.MCPToolsList{}, fmt.Errorf("invalid server ID '%s': %w", serverID, api.ErrBadRequest)
	}

	// Verify server exists
	_, err = d.mcpServerRepo.GetByID(id)
	if err != nil {
		return openapi.MCPToolsList{}, fmt.Errorf("server not found with ID %s: %w", serverID, api.ErrNotFound)
	}

	// Build MCPServerToolListOptions from params
	listOptions := models.MCPServerToolListOptions{
		ParentID: id,
	}

	// Set pagination
	listOptions.Pagination.PageSize = &params.PageSize
	listOptions.Pagination.OrderBy = (*string)(&params.OrderBy)
	listOptions.Pagination.SortOrder = (*string)(&params.SortOrder)
	if params.NextPageToken != nil {
		listOptions.Pagination.NextPageToken = params.NextPageToken
	}

	// Set filterQuery if provided
	if params.FilterQuery != "" {
		listOptions.Pagination.FilterQuery = &params.FilterQuery
	}

	// Call repository
	tools, err := d.mcpServerToolRepo.List(listOptions)
	if err != nil {
		return openapi.MCPToolsList{}, err
	}

	// Convert to OpenAPI models
	apiTools := make([]openapi.MCPTool, 0) // Initialize as empty slice, not nil
	for _, dbTool := range tools {
		apiTool := converter.ConvertDbMCPToolToOpenapi(dbTool)
		if apiTool != nil {
			apiTools = append(apiTools, *apiTool)
		}
	}

	return openapi.MCPToolsList{
		Items:         apiTools,
		Size:          int32(len(apiTools)),
		PageSize:      params.PageSize,
		NextPageToken: "", // TODO: Implement pagination token for tools
	}, nil
}

func (d *dbMCPCatalogImpl) GetMCPServerTool(ctx context.Context, serverID string, toolName string) (*openapi.MCPTool, error) {
	id, err := apiutils.ValidateIDAsInt32(serverID, "server")
	if err != nil {
		return nil, fmt.Errorf("invalid server ID '%s': %w", serverID, api.ErrBadRequest)
	}

	// Verify server exists
	_, err = d.mcpServerRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("server not found with ID %s: %w", serverID, api.ErrNotFound)
	}

	// List all tools for the server and find the one with matching name
	toolOptions := models.MCPServerToolListOptions{
		ParentID: id,
	}
	tools, err := d.mcpServerToolRepo.List(toolOptions)
	if err != nil {
		return nil, fmt.Errorf("error loading tools for server %s: %w", serverID, err)
	}

	// Find tool by name
	for _, tool := range tools {
		attrs := tool.GetAttributes()
		if attrs != nil && attrs.Name != nil && *attrs.Name == toolName {
			apiTool := converter.ConvertDbMCPToolToOpenapi(tool)
			return apiTool, nil
		}
	}

	return nil, fmt.Errorf("tool '%s' not found in server %s: %w", toolName, serverID, api.ErrNotFound)
}

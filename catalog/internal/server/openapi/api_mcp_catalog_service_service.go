package openapi

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/kubeflow/model-registry/catalog/internal/catalog"
	"github.com/kubeflow/model-registry/catalog/internal/catalog/mcpcatalog"
	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/pkg/api"
)

// MCPCatalogServiceAPIService is a service that implements the logic for the MCPCatalogServiceAPIServicer
type MCPCatalogServiceAPIService struct {
	mcpProvider catalog.MCPProvider
	mcpSources  *mcpcatalog.MCPSourceCollection
}

var _ MCPCatalogServiceAPIServicer = &MCPCatalogServiceAPIService{}

// NewMCPCatalogServiceAPIService creates a default api service
func NewMCPCatalogServiceAPIService(mcpProvider catalog.MCPProvider, mcpSources *mcpcatalog.MCPSourceCollection) MCPCatalogServiceAPIServicer {
	return &MCPCatalogServiceAPIService{
		mcpProvider: mcpProvider,
		mcpSources:  mcpSources,
	}
}

// FindMCPServers - List MCP servers.
func (m *MCPCatalogServiceAPIService) FindMCPServers(ctx context.Context, name string, q string, sourceLabel []string, filterQuery string, namedQuery string, includeTools bool, toolLimit int32, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	pageSizeInt := int32(10)

	if pageSize != "" {
		parsed, err := strconv.ParseInt(pageSize, 10, 32)
		if err != nil {
			return Response(http.StatusBadRequest, err), err
		}
		pageSizeInt = int32(parsed)
	}

	// Convert parameters to internal format
	params := catalog.ListMCPServersParams{
		Name:          name,
		Query:         q,
		FilterQuery:   filterQuery,
		NamedQuery:    namedQuery,
		IncludeTools:  includeTools,
		ToolLimit:     toolLimit,
		PageSize:      pageSizeInt,
		OrderBy:       orderBy,
		SortOrder:     sortOrder,
		NextPageToken: &nextPageToken,
	}

	servers, err := m.mcpProvider.ListMCPServers(ctx, params)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}

	return Response(http.StatusOK, servers), nil
}

// FindMCPServersFilterOptions - Lists fields, values, and named queries that can be used in `filterQuery` on the list MCP servers endpoint.
func (m *MCPCatalogServiceAPIService) FindMCPServersFilterOptions(ctx context.Context) (ImplResponse, error) {
	if m.mcpSources == nil {
		return Response(http.StatusOK, model.FilterOptionsList{}), nil
	}

	srcQueries := m.mcpSources.GetNamedQueries()
	converted := make(map[string]map[string]model.FieldFilter, len(srcQueries))
	for queryName, fieldFilters := range srcQueries {
		inner := make(map[string]model.FieldFilter, len(fieldFilters))
		for field, ff := range fieldFilters {
			inner[field] = model.FieldFilter{
				Operator: ff.Operator,
				Value:    ff.Value,
			}
		}
		converted[queryName] = inner
	}

	return Response(http.StatusOK, model.FilterOptionsList{NamedQueries: &converted}), nil
}

// GetMCPServer - Get an `MCPServer`.
func (m *MCPCatalogServiceAPIService) GetMCPServer(ctx context.Context, serverID string, includeTools bool) (ImplResponse, error) {
	server, err := m.mcpProvider.GetMCPServer(ctx, serverID, includeTools)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}

	if server == nil {
		return ErrorResponse(http.StatusNotFound, errors.New("server not found")), nil
	}

	return Response(http.StatusOK, server), nil
}

// FindMCPServerTools - List MCP server tools.
func (m *MCPCatalogServiceAPIService) FindMCPServerTools(ctx context.Context, serverID string, filterQuery string, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	pageSizeInt := int32(10)

	if pageSize != "" {
		parsed, err := strconv.ParseInt(pageSize, 10, 32)
		if err != nil {
			return Response(http.StatusBadRequest, err), err
		}
		pageSizeInt = int32(parsed)
	}

	// Convert parameters to internal format
	params := catalog.ListMCPServerToolsParams{
		FilterQuery:   filterQuery,
		PageSize:      pageSizeInt,
		OrderBy:       orderBy,
		SortOrder:     sortOrder,
		NextPageToken: &nextPageToken,
	}

	tools, err := m.mcpProvider.ListMCPServerTools(ctx, serverID, params)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}

	return Response(http.StatusOK, tools), nil
}

// GetMCPServerTool - Get an `MCPTool` from an `MCPServer`.
func (m *MCPCatalogServiceAPIService) GetMCPServerTool(ctx context.Context, serverID string, toolName string) (ImplResponse, error) {
	tool, err := m.mcpProvider.GetMCPServerTool(ctx, serverID, toolName)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}

	if tool == nil {
		return ErrorResponse(http.StatusNotFound, errors.New("tool not found")), nil
	}

	return Response(http.StatusOK, tool), nil
}

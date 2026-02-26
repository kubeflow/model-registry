package openapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/kubeflow/model-registry/catalog/internal/catalog"
	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
)

// MCPCatalogServiceAPIService is a service that implements the logic for the MCPCatalogServiceAPIServicer
type MCPCatalogServiceAPIService struct {
	mcpProvider catalog.MCPProvider
}

var _ MCPCatalogServiceAPIServicer = &MCPCatalogServiceAPIService{}

// NewMCPCatalogServiceAPIService creates a default api service
func NewMCPCatalogServiceAPIService(mcpProvider catalog.MCPProvider) MCPCatalogServiceAPIServicer {
	return &MCPCatalogServiceAPIService{
		mcpProvider: mcpProvider,
	}
}

// FindMCPServers - List MCP servers.
func (m *MCPCatalogServiceAPIService) FindMCPServers(ctx context.Context, name string, q string, sourceLabel []string, filterQuery string, namedQuery string, includeTools bool, toolLimit int32, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	var err error
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
		return ErrorResponse(http.StatusInternalServerError, err), err
	}

	return Response(http.StatusOK, servers), nil
}

// FindMCPServersFilterOptions - Lists fields, values, and named queries that can be used in `filterQuery` on the list MCP servers endpoint.
func (m *MCPCatalogServiceAPIService) FindMCPServersFilterOptions(ctx context.Context) (ImplResponse, error) {
	return ErrorResponse(http.StatusNotImplemented, errors.New("FindMCPServersFilterOptions not implemented")), nil
}

// GetMCPServer - Get an `MCPServer`.
func (m *MCPCatalogServiceAPIService) GetMCPServer(ctx context.Context, serverID string, includeTools bool) (ImplResponse, error) {
	server, err := m.mcpProvider.GetMCPServer(ctx, serverID, includeTools)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(fmt.Sprintf("%v", err), "not found") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(fmt.Sprintf("%v", err), "invalid") {
			statusCode = http.StatusBadRequest
		}
		return ErrorResponse(statusCode, err), err
	}

	if server == nil {
		return ErrorResponse(http.StatusNotFound, errors.New("server not found")), nil
	}

	return Response(http.StatusOK, server), nil
}

// FindMCPServerTools - List MCP server tools.
func (m *MCPCatalogServiceAPIService) FindMCPServerTools(ctx context.Context, serverID string, filterQuery string, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	var err error
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
		statusCode := http.StatusInternalServerError
		if strings.Contains(fmt.Sprintf("%v", err), "not found") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(fmt.Sprintf("%v", err), "invalid") {
			statusCode = http.StatusBadRequest
		}
		return ErrorResponse(statusCode, err), err
	}

	return Response(http.StatusOK, tools), nil
}

// GetMCPServerTool - Get an `MCPTool` from an `MCPServer`.
func (m *MCPCatalogServiceAPIService) GetMCPServerTool(ctx context.Context, serverID string, toolName string) (ImplResponse, error) {
	tool, err := m.mcpProvider.GetMCPServerTool(ctx, serverID, toolName)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(fmt.Sprintf("%v", err), "not found") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(fmt.Sprintf("%v", err), "invalid") {
			statusCode = http.StatusBadRequest
		}
		return ErrorResponse(statusCode, err), err
	}

	if tool == nil {
		return ErrorResponse(http.StatusNotFound, errors.New("tool not found")), nil
	}

	return Response(http.StatusOK, tool), nil
}

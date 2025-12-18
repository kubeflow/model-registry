package openapi

import (
	"context"
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"

	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
)

// McpCatalogProvider defines the interface for MCP catalog data providers.
// This allows switching between embedded data and database-backed implementations.
type McpCatalogProvider interface {
	ListMcpServers(ctx context.Context, name string, filterQuery string, namedQuery string) ([]model.McpServer, error)
	GetMcpServer(ctx context.Context, serverId string) (*model.McpServer, error)
	GetFilterOptions(ctx context.Context) (*model.FilterOptionsList, error)
}

// McpCatalogServiceAPIService implements the McpCatalogServiceAPIServicer interface
type McpCatalogServiceAPIService struct {
	provider McpCatalogProvider
}

var _ McpCatalogServiceAPIServicer = &McpCatalogServiceAPIService{}

// NewMcpCatalogServiceAPIService creates a new MCP catalog service
func NewMcpCatalogServiceAPIService(provider McpCatalogProvider) McpCatalogServiceAPIServicer {
	return &McpCatalogServiceAPIService{
		provider: provider,
	}
}

// FindMcpServers - List All McpServers
func (s *McpCatalogServiceAPIService) FindMcpServers(ctx context.Context, name string, filterQuery string, namedQuery string, strPageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	servers, err := s.provider.ListMcpServers(ctx, name, filterQuery, namedQuery)
	if err != nil {
		return ErrorResponse(http.StatusInternalServerError, err), err
	}

	if len(servers) > math.MaxInt32 {
		err := errors.New("too many MCP servers")
		return ErrorResponse(http.StatusInternalServerError, err), err
	}

	paginator := newMcpPaginator(strPageSize, nextPageToken)

	// Sort servers
	cmpFunc, err := genMcpServerCmpFunc(orderBy, sortOrder)
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}

	// Create a copy for sorting
	sortedServers := make([]model.McpServer, len(servers))
	copy(sortedServers, servers)

	// Sort using stable sort
	mcpSortStable(sortedServers, cmpFunc)

	// Paginate
	pagedItems, nextToken := paginator.paginate(sortedServers)

	res := model.McpServerList{
		PageSize:      paginator.pageSize,
		Items:         pagedItems,
		Size:          int32(len(pagedItems)),
		NextPageToken: nextToken,
	}

	return Response(http.StatusOK, res), nil
}

// GetMcpServer - Get an McpServer
func (s *McpCatalogServiceAPIService) GetMcpServer(ctx context.Context, serverId string) (ImplResponse, error) {
	server, err := s.provider.GetMcpServer(ctx, serverId)
	if err != nil {
		return ErrorResponse(http.StatusInternalServerError, err), err
	}

	if server == nil {
		return notFound("MCP server not found"), nil
	}

	return Response(http.StatusOK, server), nil
}

// FindMcpServersFilterOptions - Lists fields and available options that can be used in `filterQuery` on the list MCP servers endpoint.
func (s *McpCatalogServiceAPIService) FindMcpServersFilterOptions(ctx context.Context) (ImplResponse, error) {
	filterOptions, err := s.provider.GetFilterOptions(ctx)
	if err != nil {
		return ErrorResponse(http.StatusInternalServerError, err), err
	}

	return Response(http.StatusOK, filterOptions), nil
}

// genMcpServerCmpFunc generates a comparison function for sorting MCP servers
func genMcpServerCmpFunc(orderBy model.OrderByField, sortOrder model.SortOrder) (func(model.McpServer, model.McpServer) int, error) {
	multiplier := 1
	switch model.SortOrder(strings.ToUpper(string(sortOrder))) {
	case model.SORTORDER_DESC:
		multiplier = -1
	case model.SORTORDER_ASC, "":
		multiplier = 1
	default:
		return nil, errors.New("unsupported sort order")
	}

	switch model.OrderByField(strings.ToUpper(string(orderBy))) {
	case model.ORDERBYFIELD_ID, "":
		return func(a, b model.McpServer) int {
			return multiplier * strings.Compare(a.Id, b.Id)
		}, nil
	case model.ORDERBYFIELD_NAME:
		return func(a, b model.McpServer) int {
			return multiplier * strings.Compare(a.Name, b.Name)
		}, nil
	default:
		return nil, errors.New("unsupported order by field for MCP servers")
	}
}

// mcpSortStable is a stable sort implementation for MCP servers
func mcpSortStable(items []model.McpServer, cmp func(model.McpServer, model.McpServer) int) {
	// Use insertion sort for stable sorting
	for i := 1; i < len(items); i++ {
		for j := i; j > 0 && cmp(items[j-1], items[j]) > 0; j-- {
			items[j-1], items[j] = items[j], items[j-1]
		}
	}
}

// mcpServerPaginator handles pagination for MCP servers
type mcpServerPaginator struct {
	pageSize int32
	offset   int32
}

// newMcpPaginator creates a paginator for MCP servers
func newMcpPaginator(strPageSize string, nextPageToken string) *mcpServerPaginator {
	pageSize := int32(10)
	if strPageSize != "" {
		parsed, err := strconv.ParseInt(strPageSize, 10, 32)
		if err == nil && parsed > 0 {
			pageSize = int32(parsed)
		}
	}

	offset := int32(0)
	if nextPageToken != "" {
		// Simple offset-based pagination
		parsed, err := strconv.ParseInt(nextPageToken, 10, 32)
		if err == nil {
			offset = int32(parsed)
		}
	}

	return &mcpServerPaginator{
		pageSize: pageSize,
		offset:   offset,
	}
}

// paginate returns a page of items and the next page token
func (p *mcpServerPaginator) paginate(items []model.McpServer) ([]model.McpServer, string) {
	start := int(p.offset)
	if start >= len(items) {
		return []model.McpServer{}, ""
	}

	end := start + int(p.pageSize)
	if end > len(items) {
		end = len(items)
	}

	result := items[start:end]

	var next string
	if end < len(items) {
		next = strconv.Itoa(end)
	}

	return result, next
}

package api

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/httpclient"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

// McpServerListEnvelope wraps the MCP server list response
type McpServerListEnvelope Envelope[*models.McpServerList, None]

// McpServerEnvelope wraps a single MCP server response
type McpServerEnvelope Envelope[*models.McpServer, None]

// McpCatalogSourceListEnvelope wraps the MCP catalog source list response
type McpCatalogSourceListEnvelope Envelope[*models.McpCatalogSourceList, None]

// McpFilterOptionsListEnvelope wraps the MCP filter options response
type McpFilterOptionsListEnvelope Envelope[*models.FilterOptionsList, None]

// GetAllMcpServersHandler handles GET /api/v1/mcp_catalog/mcp_servers
func (app *App) GetAllMcpServersHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	client, ok := r.Context().Value(constants.ModelCatalogHttpClientKey).(httpclient.HTTPClientInterface)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("catalog REST client not found"))
		return
	}

	mcpServers, err := app.repositories.ModelCatalogClient.GetAllMcpServers(client, r.URL.Query())

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	serverList := McpServerListEnvelope{
		Data: mcpServers,
	}

	err = app.WriteJSON(w, http.StatusOK, serverList, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// GetMcpServerHandler handles GET /api/v1/mcp_catalog/mcp_servers/:server_id
func (app *App) GetMcpServerHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	client, ok := r.Context().Value(constants.ModelCatalogHttpClientKey).(httpclient.HTTPClientInterface)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("catalog REST client not found"))
		return
	}

	serverId := ps.ByName(McpServerId)

	mcpServer, err := app.repositories.ModelCatalogClient.GetMcpServer(client, serverId)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if mcpServer == nil {
		app.notFoundResponse(w, r)
		return
	}

	serverEnvelope := McpServerEnvelope{
		Data: mcpServer,
	}

	err = app.WriteJSON(w, http.StatusOK, serverEnvelope, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// GetMcpFilterOptionsHandler handles GET /api/v1/mcp_catalog/filter_options
func (app *App) GetMcpFilterOptionsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	client, ok := r.Context().Value(constants.ModelCatalogHttpClientKey).(httpclient.HTTPClientInterface)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("catalog REST client not found"))
		return
	}

	filterOptions, err := app.repositories.ModelCatalogClient.GetMcpFilterOptions(client)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := McpFilterOptionsListEnvelope{
		Data: filterOptions,
	}

	err = app.WriteJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// GetAllMcpSourcesHandler handles GET /api/v1/mcp_catalog/sources
func (app *App) GetAllMcpSourcesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	client, ok := r.Context().Value(constants.ModelCatalogHttpClientKey).(httpclient.HTTPClientInterface)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("catalog REST client not found"))
		return
	}

	mcpSources, err := app.repositories.ModelCatalogClient.GetAllMcpSources(client, r.URL.Query())

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	sourcesList := McpCatalogSourceListEnvelope{
		Data: mcpSources,
	}

	err = app.WriteJSON(w, http.StatusOK, sourcesList, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

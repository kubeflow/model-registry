package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/httpclient"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type McpServerListEnvelope Envelope[*models.McpServerList, None]
type McpServerFilterOptionEnvelope Envelope[*models.FilterOption, None]
type McpServerFilterOptionsListEnvelope Envelope[*models.FilterOptionsList, None]
type McpServerEnvelope Envelope[*models.McpServer, None]
type McpServerToolsListEnvelope Envelope[*models.McpToolList, None]

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

	mcpServerList := McpServerListEnvelope{
		Data: mcpServers,
	}

	err = app.WriteJSON(w, http.StatusOK, mcpServerList, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) GetMcpServersFiltersHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	client, ok := r.Context().Value(constants.ModelCatalogHttpClientKey).(httpclient.HTTPClientInterface)

	if !ok {
		app.serverErrorResponse(w, r, errors.New("catalog REST client not found"))
		return
	}

	mcpServerFilterOptions, err := app.repositories.ModelCatalogClient.GetMcpServersFilter(client)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	mcpServerFilterOptionsList := McpServerFilterOptionsListEnvelope{
		Data: mcpServerFilterOptions,
	}

	err = app.WriteJSON(w, http.StatusOK, mcpServerFilterOptionsList, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) GetMcpServerHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	client, ok := r.Context().Value(constants.ModelCatalogHttpClientKey).(httpclient.HTTPClientInterface)

	if !ok {
		app.serverErrorResponse(w, r, errors.New("catalog REST client not found"))
		return
	}

	serverId := ps.ByName(McpServerId)

	if serverId == "" {
		app.badRequestResponse(w, r, fmt.Errorf("server_id is required"))
		return
	}

	server, err := app.repositories.ModelCatalogClient.GetMcpServer(client, serverId)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	mcpServer := McpServerEnvelope{
		Data: server,
	}

	err = app.WriteJSON(w, http.StatusOK, mcpServer, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *App) GetMcpServersToolsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	client, ok := r.Context().Value(constants.ModelCatalogHttpClientKey).(httpclient.HTTPClientInterface)

	if !ok {
		app.serverErrorResponse(w, r, errors.New("catalog REST client not found"))
		return
	}

	serverId := ps.ByName(McpServerId)

	if serverId == "" {
		app.badRequestResponse(w, r, fmt.Errorf("server_id is required"))
		return
	}

	mcpServerTools, err := app.repositories.ModelCatalogClient.GetMcpServersTools(client, serverId)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	mcpServerToolList := McpServerToolsListEnvelope{
		Data: mcpServerTools,
	}

	err = app.WriteJSON(w, http.StatusOK, mcpServerToolList, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	helper "github.com/kubeflow/model-registry/ui/bff/internal/helpers"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type CatalogFilterOptionEnvelope Envelope[models.FilterOption, None]
type CatalogFilterOptionsListEnvelope Envelope[models.FilterOptionsList, None]

func (app *App) GetCatalogFilterListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	// TODO: this is to unblock frontend development, will be removed in actual implementation
	mockCatalogFilterOption := mocks.GetFilterOptionsListMock()

	catalogFilterOptionList := CatalogFilterOptionsListEnvelope{
		Data: mockCatalogFilterOption,
	}

	err := app.WriteJSON(w, http.StatusOK, catalogFilterOptionList, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

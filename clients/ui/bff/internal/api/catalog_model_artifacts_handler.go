package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	helper "github.com/kubeflow/model-registry/ui/bff/internal/helpers"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type CatalogModelArtifactEnvelope Envelope[models.CatalogModelArtifact, None]
type CatalogModelArtifactListEnvelope Envelope[models.CatalogModelArtifactList, None]

func (app *App) GetAllCatalogModelArtifactsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// TODO: Implement actual catalog API call for model artifacts
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	// TODO: this is to unblock frontend development, will be removed in actual implementation
	mockCatalogModelArtifactList := mocks.GetCatalogModelArtifactListMock()

	catalogModelArtifactList := CatalogModelArtifactListEnvelope{
		Data: mockCatalogModelArtifactList,
	}

	err := app.WriteJSON(w, http.StatusOK, catalogModelArtifactList, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

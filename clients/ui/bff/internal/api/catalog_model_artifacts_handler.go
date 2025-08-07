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

func (app *App) GetCatalogModelArtifactsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: Implement actual catalog API call for model artifacts
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	catalogSourceID := ps.ByName(SourceId)
	modelName := ps.ByName(CatalogModelName)

	// TODO: this is to unblock frontend development, will be removed in actual implementation
	allMockModels := mocks.GetCatalogModelMocks()
	var filteredMockModelArtifacts []models.CatalogModelArtifact

	for _, model := range allMockModels {
		if model.Name == modelName && model.SourceId != nil && *model.SourceId == catalogSourceID {
			artifacts := models.CatalogModelArtifact{
				Uri:                      "",
				CreateTimeSinceEpoch:     model.CreateTimeSinceEpoch,
				LastUpdateTimeSinceEpoch: model.LastUpdateTimeSinceEpoch,
				CustomProperties:         model.CustomProperties,
			}
			filteredMockModelArtifacts = append(filteredMockModelArtifacts, artifacts)
		}
	}

	if len(filteredMockModelArtifacts) == 0 {
		app.notFoundResponse(w, r)
		return
	}

	catalogModelArtifacts := models.CatalogModelArtifactList{
		Items:         filteredMockModelArtifacts,
		Size:          int32(len(filteredMockModelArtifacts)),
		PageSize:      int32(10),
		NextPageToken: "",
	}

	catalogModelArtifactList := CatalogModelArtifactListEnvelope{
		Data: catalogModelArtifacts,
	}

	err := app.WriteJSON(w, http.StatusOK, catalogModelArtifactList, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

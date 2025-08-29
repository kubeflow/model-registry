package api

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	helper "github.com/kubeflow/model-registry/ui/bff/internal/helpers"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type CatalogSourceEnvelope Envelope[models.CatalogSource, None]
type CatalogSourceListEnvelope Envelope[models.CatalogSourceList, None]

func (app *App) GetAllCatalogSourcesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// TODO: Implement actual catalog API call for sources
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	name := r.URL.Query().Get("name")

	// TODO: Implement actual catalog API call to get sources
	allMockSources := mocks.GetCatalogSourceMocks()

	// TODO: this is to unblock frontend development, will be removed in actual implementation
	var fileredMockSources []models.CatalogSource
	if name != "" {
		nameFilterLower := strings.ToLower(name)
		for _, source := range allMockSources {
			if strings.ToLower(source.Id) == nameFilterLower || strings.ToLower(source.Name) == nameFilterLower {
				fileredMockSources = append(fileredMockSources, source)
			}
		}
	} else {
		fileredMockSources = allMockSources
	}

	catalogSourceList := models.CatalogSourceList{
		Items:         fileredMockSources,
		PageSize:      int32(10),
		NextPageToken: "",
		Size:          int32(len(fileredMockSources)),
	}

	CatalogSources := CatalogSourceListEnvelope{
		Data: catalogSourceList,
	}

	err := app.WriteJSON(w, http.StatusOK, CatalogSources, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

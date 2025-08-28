package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	helper "github.com/kubeflow/model-registry/ui/bff/internal/helpers"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type ModelCatalogSettingsSourceEnvelope Envelope[models.ConfigMapKind, None]
type ModelCatalogSettingsSourceListEnvelope Envelope[[]models.ConfigMapKind, None]
type ModelCatalogSourcePayloadEnvelope Envelope[models.CatalogSource, None]

func (app *App) GetAllCatalogSourcesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in the context"))
	}

	newCatalogSource := []models.CatalogSource{
		{
			Name:    "name",
			Id:      "catalogSourceId",
			Type:    "catalogType",
			Enabled: BoolPtr(true),
			Properties: &models.CatalogSourceProperties{
				YamlCatalogPath: stringToPointer("a"),
				Models:          []string{"rhelai1/modelcar-granite-7b-starter"},
				ExcludedModels:  []string{"model-a:1.0", "model-b:*"},
			},
		},
		{
			Name:    "name",
			Id:      "catalogSourceId",
			Type:    "catalogType",
			Enabled: BoolPtr(true),
			Properties: &models.CatalogSourceProperties{
				YamlCatalogPath: stringToPointer("a"),
				Models:          []string{"rhelai1/modelcar-granite-7b-starter"},
				ExcludedModels:  []string{"model-a:1.0", "model-b:*"},
			},
		},
		{
			Name:    "name",
			Id:      "catalogSourceId",
			Type:    "catalogType",
			Enabled: BoolPtr(true),
			Properties: &models.CatalogSourceProperties{
				YamlCatalogPath: stringToPointer("a"),
				Models:          []string{"rhelai1/modelcar-granite-7b-starter"},
				ExcludedModels:  []string{"model-a:1.0", "model-b:*"},
			},
		},
	}
	catalogSources := createCatalogSource(namespace, newCatalogSource)

	modelCatalogSource := ModelCatalogSettingsSourceEnvelope{
		Data: catalogSources,
	}

	err := app.WriteJSON(w, http.StatusOK, modelCatalogSource, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *App) GetCatalogSourceHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in the context"))
	}

	catalogSourceId := ps.ByName(SourceId)
	newCatalogSource := []models.CatalogSource{
		{
			Name:    "name",
			Id:      catalogSourceId,
			Type:    "catalogType",
			Enabled: BoolPtr(true),
			Properties: &models.CatalogSourceProperties{
				YamlCatalogPath: stringToPointer("a"),
				Models:          []string{"rhelai1/modelcar-granite-7b-starter"},
				ExcludedModels:  []string{"model-a:1.0", "model-b:*"},
			},
		},
	}
	catalogSource := createCatalogSource(namespace, newCatalogSource)

	modelCatalogSource := ModelCatalogSettingsSourceEnvelope{
		Data: catalogSource,
	}

	err := app.WriteJSON(w, http.StatusOK, modelCatalogSource, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *App) CreateCatalogSourceHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in the context"))
	}

	var envelope ModelCatalogSourcePayloadEnvelope
	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error decoding JSON:: %v", err.Error()))
		return
	}

	var sourceName = envelope.Data.Name
	var sourceId = envelope.Data.Id
	var sourceType = envelope.Data.Type

	if sourceName == "" {
		app.badRequestResponse(w, r, fmt.Errorf("source name is required"))
		return
	}
	if sourceId == "" {
		app.badRequestResponse(w, r, fmt.Errorf("source ID is required"))
		return
	}
	if sourceType == "" {
		app.badRequestResponse(w, r, fmt.Errorf("source type is required"))
		return
	}

	ctxLogger.Info("Creating catalog source", "name", sourceName)

	newCatalogSource := []models.CatalogSource{
		{
			Name:    "name",
			Id:      sourceId,
			Type:    "catalogType",
			Enabled: BoolPtr(true),
			Properties: &models.CatalogSourceProperties{
				YamlCatalogPath: stringToPointer("a"),
				Models:          []string{"rhelai1/modelcar-granite-7b-starter"},
				ExcludedModels:  []string{"model-a:1.0", "model-b:*"},
			},
		},
	}
	catalogSource := createCatalogSource(namespace, newCatalogSource)

	modelCatalogSource := ModelCatalogSettingsSourceEnvelope{
		Data: catalogSource,
	}
	n := len(modelCatalogSource.Data.Data.SourcesYaml.Catalogs) - 1

	w.Header().Set("Location", r.URL.JoinPath(modelCatalogSource.Data.Data.SourcesYaml.Catalogs[n].Id).String())
	writeErr := app.WriteJSON(w, http.StatusCreated, modelCatalogSource, nil)
	if writeErr != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error writing JSON"))
		return
	}

}

func (app *App) UpdateCatalogSourceHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in the context"))
	}

	// this is the temoprary fix to start fronetend development
	catalogSourceId := ps.ByName(SourceId)
	newCatalogSource := []models.CatalogSource{
		{
			Name:    "name",
			Id:      catalogSourceId,
			Type:    "catalogType",
			Enabled: BoolPtr(true),
			Properties: &models.CatalogSourceProperties{
				YamlCatalogPath: stringToPointer("a"),
				Models:          []string{"rhelai1/modelcar-granite-7b-starter"},
				ExcludedModels:  []string{"model-a:1.0", "model-b:*"},
			},
		},
	}
	catalogSource := createCatalogSource(namespace, newCatalogSource)

	modelCatalogSource := ModelCatalogSettingsSourceEnvelope{
		Data: catalogSource,
	}

	err := app.WriteJSON(w, http.StatusOK, modelCatalogSource, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) DeleteCatalogSourceHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in the context"))
	}

	// this is the temoprary fix to start fronetend development
	catalogSourceId := ps.ByName(SourceId)
	newCatalogSource := []models.CatalogSource{
		{
			Name:    "name",
			Id:      catalogSourceId,
			Type:    "catalogType",
			Enabled: BoolPtr(true),
			Properties: &models.CatalogSourceProperties{
				YamlCatalogPath: stringToPointer("a"),
				Models:          []string{"rhelai1/modelcar-granite-7b-starter"},
				ExcludedModels:  []string{"model-a:1.0", "model-b:*"},
			},
		},
	}
	catalogSource := createCatalogSource(namespace, newCatalogSource)

	modelCatalogSource := ModelCatalogSettingsSourceEnvelope{
		Data: catalogSource,
	}

	err := app.WriteJSON(w, http.StatusOK, modelCatalogSource, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func createCatalogSource(namespace string, newCatalogSources []models.CatalogSource) models.ConfigMapKind {
	creationTime, _ := time.Parse(time.RFC3339, "2024-03-14T08:01:42Z")
	catalogSources := []models.CatalogSource{}

	upadtedCatalogSources := append(catalogSources, newCatalogSources...)

	// and also update the type envelop and types in models if required

	return models.ConfigMapKind{
		APIVersion: "",
		Kind:       "ConfigMap",
		Metadata: models.Metadata{
			Name:              "model-catalog-sources",
			Namespace:         namespace,
			CreationTimestamp: creationTime,
			Annotations:       map[string]string{},
		},
		Data: models.ConfigMapData{
			SamplaCatalogYaml: &models.CatalogContent{
				Source: "",
				Models: []models.BaseModel{
					{
						Name: "",
					},
				},
			},
			SourcesYaml: &models.SourcesContent{
				Catalogs: upadtedCatalogSources,
			},
		},
	}
}

// func addCatalogSource(catalogSources []models.CatalogSource, name string, id string, catalogType string) []models.CatalogSource {
// 	// newCatalogSource := models.CatalogSource{
// 	// 					Name:    name,
// 	// 					Id:      id,
// 	// 					Type:    catalogType,
// 	// 					Enabled: BoolPtr(true),
// 	// 					Properties: &models.CatalogSourceProperties{
// 	// 						YamlCatalogPath: stringToPointer("a"),
// 	// 						Models:          []string{"rhelai1/modelcar-granite-7b-starter"},
// 	// 						ExcludedModels:  []string{"model-a:1.0", "model-b:*"},
// 	// 					},
// 	// }
// 	return append(catalogSources, newCatalogSource)
// }

func stringToPointer(s string) *string {
	return &s
}

func BoolPtr(b bool) *bool {
	return &b
}

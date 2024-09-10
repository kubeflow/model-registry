/*
 * Model Registry REST API
 *
 * REST API for Model Registry to create and manage ML model metadata
 *
 * API version: v1alpha3
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	model "github.com/kubeflow/model-registry/pkg/openapi"
)

// ModelRegistryServiceAPIController binds http requests to an api service and writes the service results to the http response
type ModelRegistryServiceAPIController struct {
	service      ModelRegistryServiceAPIServicer
	errorHandler ErrorHandler
}

// ModelRegistryServiceAPIOption for how the controller is set up.
type ModelRegistryServiceAPIOption func(*ModelRegistryServiceAPIController)

// WithModelRegistryServiceAPIErrorHandler inject ErrorHandler into controller
func WithModelRegistryServiceAPIErrorHandler(h ErrorHandler) ModelRegistryServiceAPIOption {
	return func(c *ModelRegistryServiceAPIController) {
		c.errorHandler = h
	}
}

// NewModelRegistryServiceAPIController creates a default api controller
func NewModelRegistryServiceAPIController(s ModelRegistryServiceAPIServicer, opts ...ModelRegistryServiceAPIOption) Router {
	controller := &ModelRegistryServiceAPIController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the ModelRegistryServiceAPIController
func (c *ModelRegistryServiceAPIController) Routes() Routes {
	return Routes{
		"CreateEnvironmentInferenceService": Route{
			strings.ToUpper("Post"),
			"/api/model_registry/v1alpha3/serving_environments/{servingenvironmentId}/inference_services",
			c.CreateEnvironmentInferenceService,
		},
		"CreateInferenceService": Route{
			strings.ToUpper("Post"),
			"/api/model_registry/v1alpha3/inference_services",
			c.CreateInferenceService,
		},
		"CreateInferenceServiceServe": Route{
			strings.ToUpper("Post"),
			"/api/model_registry/v1alpha3/inference_services/{inferenceserviceId}/serves",
			c.CreateInferenceServiceServe,
		},
		"CreateModelArtifact": Route{
			strings.ToUpper("Post"),
			"/api/model_registry/v1alpha3/model_artifacts",
			c.CreateModelArtifact,
		},
		"CreateModelVersion": Route{
			strings.ToUpper("Post"),
			"/api/model_registry/v1alpha3/model_versions",
			c.CreateModelVersion,
		},
		"CreateRegisteredModel": Route{
			strings.ToUpper("Post"),
			"/api/model_registry/v1alpha3/registered_models",
			c.CreateRegisteredModel,
		},
		"CreateRegisteredModelVersion": Route{
			strings.ToUpper("Post"),
			"/api/model_registry/v1alpha3/registered_models/{registeredmodelId}/versions",
			c.CreateRegisteredModelVersion,
		},
		"CreateServingEnvironment": Route{
			strings.ToUpper("Post"),
			"/api/model_registry/v1alpha3/serving_environments",
			c.CreateServingEnvironment,
		},
		"FindInferenceService": Route{
			strings.ToUpper("Get"),
			"/api/model_registry/v1alpha3/inference_service",
			c.FindInferenceService,
		},
		"FindModelArtifact": Route{
			strings.ToUpper("Get"),
			"/api/model_registry/v1alpha3/model_artifact",
			c.FindModelArtifact,
		},
		"FindModelVersion": Route{
			strings.ToUpper("Get"),
			"/api/model_registry/v1alpha3/model_version",
			c.FindModelVersion,
		},
		"FindRegisteredModel": Route{
			strings.ToUpper("Get"),
			"/api/model_registry/v1alpha3/registered_model",
			c.FindRegisteredModel,
		},
		"FindServingEnvironment": Route{
			strings.ToUpper("Get"),
			"/api/model_registry/v1alpha3/serving_environment",
			c.FindServingEnvironment,
		},
		"GetEnvironmentInferenceServices": Route{
			strings.ToUpper("Get"),
			"/api/model_registry/v1alpha3/serving_environments/{servingenvironmentId}/inference_services",
			c.GetEnvironmentInferenceServices,
		},
		"GetInferenceService": Route{
			strings.ToUpper("Get"),
			"/api/model_registry/v1alpha3/inference_services/{inferenceserviceId}",
			c.GetInferenceService,
		},
		"GetInferenceServiceModel": Route{
			strings.ToUpper("Get"),
			"/api/model_registry/v1alpha3/inference_services/{inferenceserviceId}/model",
			c.GetInferenceServiceModel,
		},
		"GetInferenceServiceServes": Route{
			strings.ToUpper("Get"),
			"/api/model_registry/v1alpha3/inference_services/{inferenceserviceId}/serves",
			c.GetInferenceServiceServes,
		},
		"GetInferenceServiceVersion": Route{
			strings.ToUpper("Get"),
			"/api/model_registry/v1alpha3/inference_services/{inferenceserviceId}/version",
			c.GetInferenceServiceVersion,
		},
		"GetInferenceServices": Route{
			strings.ToUpper("Get"),
			"/api/model_registry/v1alpha3/inference_services",
			c.GetInferenceServices,
		},
		"GetModelArtifact": Route{
			strings.ToUpper("Get"),
			"/api/model_registry/v1alpha3/model_artifacts/{modelartifactId}",
			c.GetModelArtifact,
		},
		"GetModelArtifacts": Route{
			strings.ToUpper("Get"),
			"/api/model_registry/v1alpha3/model_artifacts",
			c.GetModelArtifacts,
		},
		"GetModelVersion": Route{
			strings.ToUpper("Get"),
			"/api/model_registry/v1alpha3/model_versions/{modelversionId}",
			c.GetModelVersion,
		},
		"GetModelVersionArtifacts": Route{
			strings.ToUpper("Get"),
			"/api/model_registry/v1alpha3/model_versions/{modelversionId}/artifacts",
			c.GetModelVersionArtifacts,
		},
		"GetModelVersions": Route{
			strings.ToUpper("Get"),
			"/api/model_registry/v1alpha3/model_versions",
			c.GetModelVersions,
		},
		"GetRegisteredModel": Route{
			strings.ToUpper("Get"),
			"/api/model_registry/v1alpha3/registered_models/{registeredmodelId}",
			c.GetRegisteredModel,
		},
		"GetRegisteredModelVersions": Route{
			strings.ToUpper("Get"),
			"/api/model_registry/v1alpha3/registered_models/{registeredmodelId}/versions",
			c.GetRegisteredModelVersions,
		},
		"GetRegisteredModels": Route{
			strings.ToUpper("Get"),
			"/api/model_registry/v1alpha3/registered_models",
			c.GetRegisteredModels,
		},
		"GetServingEnvironment": Route{
			strings.ToUpper("Get"),
			"/api/model_registry/v1alpha3/serving_environments/{servingenvironmentId}",
			c.GetServingEnvironment,
		},
		"GetServingEnvironments": Route{
			strings.ToUpper("Get"),
			"/api/model_registry/v1alpha3/serving_environments",
			c.GetServingEnvironments,
		},
		"UpdateInferenceService": Route{
			strings.ToUpper("Patch"),
			"/api/model_registry/v1alpha3/inference_services/{inferenceserviceId}",
			c.UpdateInferenceService,
		},
		"UpdateModelArtifact": Route{
			strings.ToUpper("Patch"),
			"/api/model_registry/v1alpha3/model_artifacts/{modelartifactId}",
			c.UpdateModelArtifact,
		},
		"UpdateModelVersion": Route{
			strings.ToUpper("Patch"),
			"/api/model_registry/v1alpha3/model_versions/{modelversionId}",
			c.UpdateModelVersion,
		},
		"UpdateRegisteredModel": Route{
			strings.ToUpper("Patch"),
			"/api/model_registry/v1alpha3/registered_models/{registeredmodelId}",
			c.UpdateRegisteredModel,
		},
		"UpdateServingEnvironment": Route{
			strings.ToUpper("Patch"),
			"/api/model_registry/v1alpha3/serving_environments/{servingenvironmentId}",
			c.UpdateServingEnvironment,
		},
		"UpsertModelVersionArtifact": Route{
			strings.ToUpper("Post"),
			"/api/model_registry/v1alpha3/model_versions/{modelversionId}/artifacts",
			c.UpsertModelVersionArtifact,
		},
	}
}

// CreateEnvironmentInferenceService - Create a InferenceService in ServingEnvironment
func (c *ModelRegistryServiceAPIController) CreateEnvironmentInferenceService(w http.ResponseWriter, r *http.Request) {
	servingenvironmentIdParam := chi.URLParam(r, "servingenvironmentId")
	inferenceServiceCreateParam := model.InferenceServiceCreate{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&inferenceServiceCreateParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertInferenceServiceCreateRequired(inferenceServiceCreateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertInferenceServiceCreateConstraints(inferenceServiceCreateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateEnvironmentInferenceService(r.Context(), servingenvironmentIdParam, inferenceServiceCreateParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// CreateInferenceService - Create a InferenceService
func (c *ModelRegistryServiceAPIController) CreateInferenceService(w http.ResponseWriter, r *http.Request) {
	inferenceServiceCreateParam := model.InferenceServiceCreate{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&inferenceServiceCreateParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertInferenceServiceCreateRequired(inferenceServiceCreateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertInferenceServiceCreateConstraints(inferenceServiceCreateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateInferenceService(r.Context(), inferenceServiceCreateParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// CreateInferenceServiceServe - Create a ServeModel action in a InferenceService
func (c *ModelRegistryServiceAPIController) CreateInferenceServiceServe(w http.ResponseWriter, r *http.Request) {
	inferenceserviceIdParam := chi.URLParam(r, "inferenceserviceId")
	serveModelCreateParam := model.ServeModelCreate{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&serveModelCreateParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertServeModelCreateRequired(serveModelCreateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertServeModelCreateConstraints(serveModelCreateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateInferenceServiceServe(r.Context(), inferenceserviceIdParam, serveModelCreateParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// CreateModelArtifact - Create a ModelArtifact
func (c *ModelRegistryServiceAPIController) CreateModelArtifact(w http.ResponseWriter, r *http.Request) {
	modelArtifactCreateParam := model.ModelArtifactCreate{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&modelArtifactCreateParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertModelArtifactCreateRequired(modelArtifactCreateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertModelArtifactCreateConstraints(modelArtifactCreateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateModelArtifact(r.Context(), modelArtifactCreateParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// CreateModelVersion - Create a ModelVersion
func (c *ModelRegistryServiceAPIController) CreateModelVersion(w http.ResponseWriter, r *http.Request) {
	modelVersionCreateParam := model.ModelVersionCreate{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&modelVersionCreateParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertModelVersionCreateRequired(modelVersionCreateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertModelVersionCreateConstraints(modelVersionCreateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateModelVersion(r.Context(), modelVersionCreateParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// CreateRegisteredModel - Create a RegisteredModel
func (c *ModelRegistryServiceAPIController) CreateRegisteredModel(w http.ResponseWriter, r *http.Request) {
	registeredModelCreateParam := model.RegisteredModelCreate{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&registeredModelCreateParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertRegisteredModelCreateRequired(registeredModelCreateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertRegisteredModelCreateConstraints(registeredModelCreateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateRegisteredModel(r.Context(), registeredModelCreateParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// CreateRegisteredModelVersion - Create a ModelVersion in RegisteredModel
func (c *ModelRegistryServiceAPIController) CreateRegisteredModelVersion(w http.ResponseWriter, r *http.Request) {
	registeredmodelIdParam := chi.URLParam(r, "registeredmodelId")
	modelVersionParam := model.ModelVersion{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&modelVersionParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertModelVersionRequired(modelVersionParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertModelVersionConstraints(modelVersionParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateRegisteredModelVersion(r.Context(), registeredmodelIdParam, modelVersionParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// CreateServingEnvironment - Create a ServingEnvironment
func (c *ModelRegistryServiceAPIController) CreateServingEnvironment(w http.ResponseWriter, r *http.Request) {
	servingEnvironmentCreateParam := model.ServingEnvironmentCreate{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&servingEnvironmentCreateParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertServingEnvironmentCreateRequired(servingEnvironmentCreateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertServingEnvironmentCreateConstraints(servingEnvironmentCreateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateServingEnvironment(r.Context(), servingEnvironmentCreateParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// FindInferenceService - Get an InferenceServices that matches search parameters.
func (c *ModelRegistryServiceAPIController) FindInferenceService(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	nameParam := query.Get("name")
	externalIdParam := query.Get("externalId")
	parentResourceIdParam := query.Get("parentResourceId")
	result, err := c.service.FindInferenceService(r.Context(), nameParam, externalIdParam, parentResourceIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// FindModelArtifact - Get a ModelArtifact that matches search parameters.
func (c *ModelRegistryServiceAPIController) FindModelArtifact(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	nameParam := query.Get("name")
	externalIdParam := query.Get("externalId")
	parentResourceIdParam := query.Get("parentResourceId")
	result, err := c.service.FindModelArtifact(r.Context(), nameParam, externalIdParam, parentResourceIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// FindModelVersion - Get a ModelVersion that matches search parameters.
func (c *ModelRegistryServiceAPIController) FindModelVersion(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	nameParam := query.Get("name")
	externalIdParam := query.Get("externalId")
	parentResourceIdParam := query.Get("parentResourceId")
	result, err := c.service.FindModelVersion(r.Context(), nameParam, externalIdParam, parentResourceIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// FindRegisteredModel - Get a RegisteredModel that matches search parameters.
func (c *ModelRegistryServiceAPIController) FindRegisteredModel(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	nameParam := query.Get("name")
	externalIdParam := query.Get("externalId")
	result, err := c.service.FindRegisteredModel(r.Context(), nameParam, externalIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// FindServingEnvironment - Find ServingEnvironment
func (c *ModelRegistryServiceAPIController) FindServingEnvironment(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	nameParam := query.Get("name")
	externalIdParam := query.Get("externalId")
	result, err := c.service.FindServingEnvironment(r.Context(), nameParam, externalIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetEnvironmentInferenceServices - List All ServingEnvironment's InferenceServices
func (c *ModelRegistryServiceAPIController) GetEnvironmentInferenceServices(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	servingenvironmentIdParam := chi.URLParam(r, "servingenvironmentId")
	nameParam := query.Get("name")
	externalIdParam := query.Get("externalId")
	pageSizeParam := query.Get("pageSize")
	orderByParam := query.Get("orderBy")
	sortOrderParam := query.Get("sortOrder")
	nextPageTokenParam := query.Get("nextPageToken")
	result, err := c.service.GetEnvironmentInferenceServices(r.Context(), servingenvironmentIdParam, nameParam, externalIdParam, pageSizeParam, model.OrderByField(orderByParam), model.SortOrder(sortOrderParam), nextPageTokenParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetInferenceService - Get a InferenceService
func (c *ModelRegistryServiceAPIController) GetInferenceService(w http.ResponseWriter, r *http.Request) {
	inferenceserviceIdParam := chi.URLParam(r, "inferenceserviceId")
	result, err := c.service.GetInferenceService(r.Context(), inferenceserviceIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetInferenceServiceModel - Get InferenceService's RegisteredModel
func (c *ModelRegistryServiceAPIController) GetInferenceServiceModel(w http.ResponseWriter, r *http.Request) {
	inferenceserviceIdParam := chi.URLParam(r, "inferenceserviceId")
	result, err := c.service.GetInferenceServiceModel(r.Context(), inferenceserviceIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetInferenceServiceServes - List All InferenceService's ServeModel actions
func (c *ModelRegistryServiceAPIController) GetInferenceServiceServes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	inferenceserviceIdParam := chi.URLParam(r, "inferenceserviceId")
	nameParam := query.Get("name")
	externalIdParam := query.Get("externalId")
	pageSizeParam := query.Get("pageSize")
	orderByParam := query.Get("orderBy")
	sortOrderParam := query.Get("sortOrder")
	nextPageTokenParam := query.Get("nextPageToken")
	result, err := c.service.GetInferenceServiceServes(r.Context(), inferenceserviceIdParam, nameParam, externalIdParam, pageSizeParam, model.OrderByField(orderByParam), model.SortOrder(sortOrderParam), nextPageTokenParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetInferenceServiceVersion - Get InferenceService's ModelVersion
func (c *ModelRegistryServiceAPIController) GetInferenceServiceVersion(w http.ResponseWriter, r *http.Request) {
	inferenceserviceIdParam := chi.URLParam(r, "inferenceserviceId")
	result, err := c.service.GetInferenceServiceVersion(r.Context(), inferenceserviceIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetInferenceServices - List All InferenceServices
func (c *ModelRegistryServiceAPIController) GetInferenceServices(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	pageSizeParam := query.Get("pageSize")
	orderByParam := query.Get("orderBy")
	sortOrderParam := query.Get("sortOrder")
	nextPageTokenParam := query.Get("nextPageToken")
	result, err := c.service.GetInferenceServices(r.Context(), pageSizeParam, model.OrderByField(orderByParam), model.SortOrder(sortOrderParam), nextPageTokenParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetModelArtifact - Get a ModelArtifact
func (c *ModelRegistryServiceAPIController) GetModelArtifact(w http.ResponseWriter, r *http.Request) {
	modelartifactIdParam := chi.URLParam(r, "modelartifactId")
	result, err := c.service.GetModelArtifact(r.Context(), modelartifactIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetModelArtifacts - List All ModelArtifacts
func (c *ModelRegistryServiceAPIController) GetModelArtifacts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	pageSizeParam := query.Get("pageSize")
	orderByParam := query.Get("orderBy")
	sortOrderParam := query.Get("sortOrder")
	nextPageTokenParam := query.Get("nextPageToken")
	result, err := c.service.GetModelArtifacts(r.Context(), pageSizeParam, model.OrderByField(orderByParam), model.SortOrder(sortOrderParam), nextPageTokenParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetModelVersion - Get a ModelVersion
func (c *ModelRegistryServiceAPIController) GetModelVersion(w http.ResponseWriter, r *http.Request) {
	modelversionIdParam := chi.URLParam(r, "modelversionId")
	result, err := c.service.GetModelVersion(r.Context(), modelversionIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetModelVersionArtifacts - List all artifacts associated with the `ModelVersion`
func (c *ModelRegistryServiceAPIController) GetModelVersionArtifacts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	modelversionIdParam := chi.URLParam(r, "modelversionId")
	nameParam := query.Get("name")
	externalIdParam := query.Get("externalId")
	pageSizeParam := query.Get("pageSize")
	orderByParam := query.Get("orderBy")
	sortOrderParam := query.Get("sortOrder")
	nextPageTokenParam := query.Get("nextPageToken")
	result, err := c.service.GetModelVersionArtifacts(r.Context(), modelversionIdParam, nameParam, externalIdParam, pageSizeParam, model.OrderByField(orderByParam), model.SortOrder(sortOrderParam), nextPageTokenParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetModelVersions - List All ModelVersions
func (c *ModelRegistryServiceAPIController) GetModelVersions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	pageSizeParam := query.Get("pageSize")
	orderByParam := query.Get("orderBy")
	sortOrderParam := query.Get("sortOrder")
	nextPageTokenParam := query.Get("nextPageToken")
	result, err := c.service.GetModelVersions(r.Context(), pageSizeParam, model.OrderByField(orderByParam), model.SortOrder(sortOrderParam), nextPageTokenParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetRegisteredModel - Get a RegisteredModel
func (c *ModelRegistryServiceAPIController) GetRegisteredModel(w http.ResponseWriter, r *http.Request) {
	registeredmodelIdParam := chi.URLParam(r, "registeredmodelId")
	result, err := c.service.GetRegisteredModel(r.Context(), registeredmodelIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetRegisteredModelVersions - List All RegisteredModel's ModelVersions
func (c *ModelRegistryServiceAPIController) GetRegisteredModelVersions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	registeredmodelIdParam := chi.URLParam(r, "registeredmodelId")
	nameParam := query.Get("name")
	externalIdParam := query.Get("externalId")
	pageSizeParam := query.Get("pageSize")
	orderByParam := query.Get("orderBy")
	sortOrderParam := query.Get("sortOrder")
	nextPageTokenParam := query.Get("nextPageToken")
	result, err := c.service.GetRegisteredModelVersions(r.Context(), registeredmodelIdParam, nameParam, externalIdParam, pageSizeParam, model.OrderByField(orderByParam), model.SortOrder(sortOrderParam), nextPageTokenParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetRegisteredModels - List All RegisteredModels
func (c *ModelRegistryServiceAPIController) GetRegisteredModels(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	pageSizeParam := query.Get("pageSize")
	orderByParam := query.Get("orderBy")
	sortOrderParam := query.Get("sortOrder")
	nextPageTokenParam := query.Get("nextPageToken")
	result, err := c.service.GetRegisteredModels(r.Context(), pageSizeParam, model.OrderByField(orderByParam), model.SortOrder(sortOrderParam), nextPageTokenParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetServingEnvironment - Get a ServingEnvironment
func (c *ModelRegistryServiceAPIController) GetServingEnvironment(w http.ResponseWriter, r *http.Request) {
	servingenvironmentIdParam := chi.URLParam(r, "servingenvironmentId")
	result, err := c.service.GetServingEnvironment(r.Context(), servingenvironmentIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetServingEnvironments - List All ServingEnvironments
func (c *ModelRegistryServiceAPIController) GetServingEnvironments(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	pageSizeParam := query.Get("pageSize")
	orderByParam := query.Get("orderBy")
	sortOrderParam := query.Get("sortOrder")
	nextPageTokenParam := query.Get("nextPageToken")
	result, err := c.service.GetServingEnvironments(r.Context(), pageSizeParam, model.OrderByField(orderByParam), model.SortOrder(sortOrderParam), nextPageTokenParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// UpdateInferenceService - Update a InferenceService
func (c *ModelRegistryServiceAPIController) UpdateInferenceService(w http.ResponseWriter, r *http.Request) {
	inferenceserviceIdParam := chi.URLParam(r, "inferenceserviceId")
	inferenceServiceUpdateParam := model.InferenceServiceUpdate{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&inferenceServiceUpdateParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertInferenceServiceUpdateRequired(inferenceServiceUpdateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertInferenceServiceUpdateConstraints(inferenceServiceUpdateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateInferenceService(r.Context(), inferenceserviceIdParam, inferenceServiceUpdateParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// UpdateModelArtifact - Update a ModelArtifact
func (c *ModelRegistryServiceAPIController) UpdateModelArtifact(w http.ResponseWriter, r *http.Request) {
	modelartifactIdParam := chi.URLParam(r, "modelartifactId")
	modelArtifactUpdateParam := model.ModelArtifactUpdate{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&modelArtifactUpdateParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertModelArtifactUpdateRequired(modelArtifactUpdateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertModelArtifactUpdateConstraints(modelArtifactUpdateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateModelArtifact(r.Context(), modelartifactIdParam, modelArtifactUpdateParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// UpdateModelVersion - Update a ModelVersion
func (c *ModelRegistryServiceAPIController) UpdateModelVersion(w http.ResponseWriter, r *http.Request) {
	modelversionIdParam := chi.URLParam(r, "modelversionId")
	modelVersionUpdateParam := model.ModelVersionUpdate{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&modelVersionUpdateParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertModelVersionUpdateRequired(modelVersionUpdateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertModelVersionUpdateConstraints(modelVersionUpdateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateModelVersion(r.Context(), modelversionIdParam, modelVersionUpdateParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// UpdateRegisteredModel - Update a RegisteredModel
func (c *ModelRegistryServiceAPIController) UpdateRegisteredModel(w http.ResponseWriter, r *http.Request) {
	registeredmodelIdParam := chi.URLParam(r, "registeredmodelId")
	registeredModelUpdateParam := model.RegisteredModelUpdate{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&registeredModelUpdateParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertRegisteredModelUpdateRequired(registeredModelUpdateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertRegisteredModelUpdateConstraints(registeredModelUpdateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateRegisteredModel(r.Context(), registeredmodelIdParam, registeredModelUpdateParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// UpdateServingEnvironment - Update a ServingEnvironment
func (c *ModelRegistryServiceAPIController) UpdateServingEnvironment(w http.ResponseWriter, r *http.Request) {
	servingenvironmentIdParam := chi.URLParam(r, "servingenvironmentId")
	servingEnvironmentUpdateParam := model.ServingEnvironmentUpdate{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&servingEnvironmentUpdateParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertServingEnvironmentUpdateRequired(servingEnvironmentUpdateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertServingEnvironmentUpdateConstraints(servingEnvironmentUpdateParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateServingEnvironment(r.Context(), servingenvironmentIdParam, servingEnvironmentUpdateParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// UpsertModelVersionArtifact - Upsert an Artifact in a ModelVersion
func (c *ModelRegistryServiceAPIController) UpsertModelVersionArtifact(w http.ResponseWriter, r *http.Request) {
	modelversionIdParam := chi.URLParam(r, "modelversionId")
	artifactParam := model.Artifact{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&artifactParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertArtifactRequired(artifactParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertArtifactConstraints(artifactParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpsertModelVersionArtifact(r.Context(), modelversionIdParam, artifactParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

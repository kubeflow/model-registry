/*
 * Model Registry REST API
 *
 * REST API for Model Registry to create and manage ML model metadata
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"context"
	"net/http"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/converter/generated"
	"github.com/kubeflow/model-registry/pkg/api"
	model "github.com/kubeflow/model-registry/pkg/openapi"
)

// ModelRegistryServiceAPIService is a service that implements the logic for the ModelRegistryServiceAPIServicer
// This service should implement the business logic for every endpoint for the ModelRegistryServiceAPI s.coreApi.
// Include any external packages or services that will be required by this service.
type ModelRegistryServiceAPIService struct {
	coreApi    api.ModelRegistryApi
	converter  converter.OpenAPIConverter
	reconciler converter.OpenAPIReconciler
}

// NewModelRegistryServiceAPIService creates a default api service
func NewModelRegistryServiceAPIService(coreApi api.ModelRegistryApi) ModelRegistryServiceAPIServicer {
	return &ModelRegistryServiceAPIService{
		coreApi:    coreApi,
		converter:  &generated.OpenAPIConverterImpl{},
		reconciler: &generated.OpenAPIReconcilerImpl{},
	}
}

// CreateEnvironmentInferenceService - Create a InferenceService in ServingEnvironment
func (s *ModelRegistryServiceAPIService) CreateEnvironmentInferenceService(ctx context.Context, servingenvironmentId string, inferenceServiceCreate model.InferenceServiceCreate) (ImplResponse, error) {
	inferenceServiceCreate.ServingEnvironmentId = servingenvironmentId
	return s.CreateInferenceService(ctx, inferenceServiceCreate)
	// TODO: return Response(http.StatusNotFound, Error{}), nil
}

// CreateInferenceService - Create a InferenceService
func (s *ModelRegistryServiceAPIService) CreateInferenceService(ctx context.Context, inferenceServiceCreate model.InferenceServiceCreate) (ImplResponse, error) {
	entity, err := s.converter.ConvertInferenceServiceCreate(&inferenceServiceCreate)
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}

	result, err := s.coreApi.UpsertInferenceService(entity)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusCreated, result), nil
	// TODO: return Response(http.StatusUnauthorized, Error{}), nil
}

// CreateInferenceServiceServe - Create a ServeModel action in a InferenceService
func (s *ModelRegistryServiceAPIService) CreateInferenceServiceServe(ctx context.Context, inferenceserviceId string, serveModelCreate model.ServeModelCreate) (ImplResponse, error) {
	entity, err := s.converter.ConvertServeModelCreate(&serveModelCreate)
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}

	result, err := s.coreApi.UpsertServeModel(entity, &inferenceserviceId)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusCreated, result), nil
	// TODO: return Response(http.StatusUnauthorized, Error{}), nil
}

// CreateArtifact - Create an Artifact
func (s *ModelRegistryServiceAPIService) CreateArtifact(ctx context.Context, artifactCreate model.ArtifactCreate) (ImplResponse, error) {
	entity, err := s.converter.ConvertArtifactCreate(&artifactCreate)
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}

	result, err := s.coreApi.UpsertArtifact(entity)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusCreated, result), nil
	// TODO: return Response(http.StatusUnauthorized, Error{}), nil
}

// CreateModelArtifact - Create a ModelArtifact
func (s *ModelRegistryServiceAPIService) CreateModelArtifact(ctx context.Context, modelArtifactCreate model.ModelArtifactCreate) (ImplResponse, error) {
	entity, err := s.converter.ConvertModelArtifactCreate(&modelArtifactCreate)
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}

	result, err := s.coreApi.UpsertModelArtifact(entity)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusCreated, result), nil
	// TODO: return Response(http.StatusUnauthorized, Error{}), nil
}

// CreateModelVersion - Create a ModelVersion
func (s *ModelRegistryServiceAPIService) CreateModelVersion(ctx context.Context, modelVersionCreate model.ModelVersionCreate) (ImplResponse, error) {
	modelVersion, err := s.converter.ConvertModelVersionCreate(&modelVersionCreate)
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}

	result, err := s.coreApi.UpsertModelVersion(modelVersion, &modelVersionCreate.RegisteredModelId)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusCreated, result), nil
	// TODO: return Response(http.StatusUnauthorized, Error{}), nil
}

// CreateModelVersionArtifact - Create an Artifact in a ModelVersion
func (s *ModelRegistryServiceAPIService) UpsertModelVersionArtifact(ctx context.Context, modelversionId string, artifact model.Artifact) (ImplResponse, error) {
	creating := (artifact.DocArtifact != nil && artifact.DocArtifact.Id == nil) || (artifact.ModelArtifact != nil && artifact.ModelArtifact.Id == nil)

	result, err := s.coreApi.UpsertModelVersionArtifact(&artifact, modelversionId)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	if creating {
		return Response(http.StatusCreated, result), nil
	}
	return Response(http.StatusOK, result), nil
	// return Response(http.StatusNotImplemented, nil), errors.New("unsupported artifactType")
	// TODO return Response(http.StatusOK, Artifact{}), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

// CreateRegisteredModel - Create a RegisteredModel
func (s *ModelRegistryServiceAPIService) CreateRegisteredModel(ctx context.Context, registeredModelCreate model.RegisteredModelCreate) (ImplResponse, error) {
	registeredModel, err := s.converter.ConvertRegisteredModelCreate(&registeredModelCreate)
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}

	result, err := s.coreApi.UpsertRegisteredModel(registeredModel)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusCreated, result), nil
	// TODO: return Response(http.StatusUnauthorized, Error{}), nil
}

// CreateRegisteredModelVersion - Create a ModelVersion in RegisteredModel
func (s *ModelRegistryServiceAPIService) CreateRegisteredModelVersion(ctx context.Context, registeredmodelId string, modelVersion model.ModelVersion) (ImplResponse, error) {
	result, err := s.coreApi.UpsertModelVersion(&modelVersion, apiutils.StrPtr(registeredmodelId))
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusCreated, result), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

// CreateServingEnvironment - Create a ServingEnvironment
func (s *ModelRegistryServiceAPIService) CreateServingEnvironment(ctx context.Context, servingEnvironmentCreate model.ServingEnvironmentCreate) (ImplResponse, error) {
	entity, err := s.converter.ConvertServingEnvironmentCreate(&servingEnvironmentCreate)
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}

	result, err := s.coreApi.UpsertServingEnvironment(entity)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusCreated, result), nil
	// TODO: return Response(http.StatusUnauthorized, Error{}), nil
}

// FindInferenceService - Get an InferenceServices that matches search parameters.
func (s *ModelRegistryServiceAPIService) FindInferenceService(ctx context.Context, name string, externalId string, parentResourceId string) (ImplResponse, error) {
	result, err := s.coreApi.GetInferenceServiceByParams(apiutils.StrPtr(name), apiutils.StrPtr(parentResourceId), apiutils.StrPtr(externalId))
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

// FindArtifact - Get an Artifact that matches search parameters.
func (s *ModelRegistryServiceAPIService) FindArtifact(ctx context.Context, name string, externalId string, parentResourceId string) (ImplResponse, error) {
	result, err := s.coreApi.GetArtifactByParams(apiutils.StrPtr(name), apiutils.StrPtr(parentResourceId), apiutils.StrPtr(externalId))
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

// FindModelArtifact - Get a ModelArtifact that matches search parameters.
func (s *ModelRegistryServiceAPIService) FindModelArtifact(ctx context.Context, name string, externalId string, parentResourceId string) (ImplResponse, error) {
	result, err := s.coreApi.GetModelArtifactByParams(apiutils.StrPtr(name), apiutils.StrPtr(parentResourceId), apiutils.StrPtr(externalId))
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

// FindModelVersion - Get a ModelVersion that matches search parameters.
func (s *ModelRegistryServiceAPIService) FindModelVersion(ctx context.Context, name string, externalId string, registeredModelId string) (ImplResponse, error) {
	result, err := s.coreApi.GetModelVersionByParams(apiutils.StrPtr(name), apiutils.StrPtr(registeredModelId), apiutils.StrPtr(externalId))
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

// FindRegisteredModel - Get a RegisteredModel that matches search parameters.
func (s *ModelRegistryServiceAPIService) FindRegisteredModel(ctx context.Context, name string, externalID string) (ImplResponse, error) {
	result, err := s.coreApi.GetRegisteredModelByParams(apiutils.StrPtr(name), apiutils.StrPtr(externalID))
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

// FindServingEnvironment - Find ServingEnvironment
func (s *ModelRegistryServiceAPIService) FindServingEnvironment(ctx context.Context, name string, externalID string) (ImplResponse, error) {
	result, err := s.coreApi.GetServingEnvironmentByParams(apiutils.StrPtr(name), apiutils.StrPtr(externalID))
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

// GetEnvironmentInferenceServices - List All ServingEnvironment&#39;s InferenceServices
func (s *ModelRegistryServiceAPIService) GetEnvironmentInferenceServices(ctx context.Context, servingenvironmentId string, name string, externalID string, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	listOpts, err := apiutils.BuildListOption(pageSize, orderBy, sortOrder, nextPageToken)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	result, err := s.coreApi.GetInferenceServices(listOpts, apiutils.StrPtr(servingenvironmentId), nil)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

// GetInferenceService - Get a InferenceService
func (s *ModelRegistryServiceAPIService) GetInferenceService(ctx context.Context, inferenceserviceId string) (ImplResponse, error) {
	result, err := s.coreApi.GetInferenceServiceById(inferenceserviceId)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO: return Response(http.StatusUnauthorized, Error{}), nil
}

// GetInferenceServiceModel - Get InferenceService&#39;s RegisteredModel
func (s *ModelRegistryServiceAPIService) GetInferenceServiceModel(ctx context.Context, inferenceserviceId string) (ImplResponse, error) {
	result, err := s.coreApi.GetRegisteredModelByInferenceService(inferenceserviceId)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO: return Response(http.StatusUnauthorized, Error{}), nil
}

// GetInferenceServiceServes - List All InferenceService&#39;s ServeModel actions
func (s *ModelRegistryServiceAPIService) GetInferenceServiceServes(ctx context.Context, inferenceserviceId string, name string, externalID string, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	listOpts, err := apiutils.BuildListOption(pageSize, orderBy, sortOrder, nextPageToken)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	result, err := s.coreApi.GetServeModels(listOpts, apiutils.StrPtr(inferenceserviceId))
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

// GetInferenceServiceVersion - Get InferenceService&#39;s ModelVersion
func (s *ModelRegistryServiceAPIService) GetInferenceServiceVersion(ctx context.Context, inferenceserviceId string) (ImplResponse, error) {
	result, err := s.coreApi.GetModelVersionByInferenceService(inferenceserviceId)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO: return Response(http.StatusUnauthorized, Error{}), nil
}

// GetInferenceServices - List All InferenceServices
func (s *ModelRegistryServiceAPIService) GetInferenceServices(ctx context.Context, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	listOpts, err := apiutils.BuildListOption(pageSize, orderBy, sortOrder, nextPageToken)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	result, err := s.coreApi.GetInferenceServices(listOpts, nil, nil)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

// GetArtifact - Get a Artifact
func (s *ModelRegistryServiceAPIService) GetArtifact(ctx context.Context, artifactId string) (ImplResponse, error) {
	result, err := s.coreApi.GetArtifactById(artifactId)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO: return Response(http.StatusUnauthorized, Error{}), nil
}

// GetArtifacts - List All Artifacts
func (s *ModelRegistryServiceAPIService) GetArtifacts(ctx context.Context, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	listOpts, err := apiutils.BuildListOption(pageSize, orderBy, sortOrder, nextPageToken)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	result, err := s.coreApi.GetArtifacts(listOpts, nil)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

// GetModelArtifact - Get a ModelArtifact
func (s *ModelRegistryServiceAPIService) GetModelArtifact(ctx context.Context, modelartifactId string) (ImplResponse, error) {
	result, err := s.coreApi.GetModelArtifactById(modelartifactId)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO: return Response(http.StatusUnauthorized, Error{}), nil
}

// GetModelArtifacts - List All ModelArtifacts
func (s *ModelRegistryServiceAPIService) GetModelArtifacts(ctx context.Context, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	listOpts, err := apiutils.BuildListOption(pageSize, orderBy, sortOrder, nextPageToken)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	result, err := s.coreApi.GetModelArtifacts(listOpts, nil)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

// GetModelVersion - Get a ModelVersion
func (s *ModelRegistryServiceAPIService) GetModelVersion(ctx context.Context, modelversionId string) (ImplResponse, error) {
	result, err := s.coreApi.GetModelVersionById(modelversionId)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO: return Response(http.StatusUnauthorized, Error{}), nil
}

// GetModelVersionArtifacts - List All ModelVersion&#39;s artifacts
func (s *ModelRegistryServiceAPIService) GetModelVersionArtifacts(ctx context.Context, modelversionId string, name string, externalID string, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	// TODO name unused
	// TODO externalID unused
	listOpts, err := apiutils.BuildListOption(pageSize, orderBy, sortOrder, nextPageToken)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	result, err := s.coreApi.GetArtifacts(listOpts, apiutils.StrPtr(modelversionId))
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

// GetModelVersions - List All ModelVersions
func (s *ModelRegistryServiceAPIService) GetModelVersions(ctx context.Context, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	listOpts, err := apiutils.BuildListOption(pageSize, orderBy, sortOrder, nextPageToken)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	result, err := s.coreApi.GetModelVersions(listOpts, nil)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

// GetRegisteredModel - Get a RegisteredModel
func (s *ModelRegistryServiceAPIService) GetRegisteredModel(ctx context.Context, registeredmodelId string) (ImplResponse, error) {
	result, err := s.coreApi.GetRegisteredModelById(registeredmodelId)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO: return Response(http.StatusUnauthorized, Error{}), nil
}

// GetRegisteredModelVersions - List All RegisteredModel&#39;s ModelVersions
func (s *ModelRegistryServiceAPIService) GetRegisteredModelVersions(ctx context.Context, registeredmodelId string, name string, externalID string, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	// TODO name unused
	// TODO externalID unused
	listOpts, err := apiutils.BuildListOption(pageSize, orderBy, sortOrder, nextPageToken)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	result, err := s.coreApi.GetModelVersions(listOpts, apiutils.StrPtr(registeredmodelId))
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

// GetRegisteredModels - List All RegisteredModels
func (s *ModelRegistryServiceAPIService) GetRegisteredModels(ctx context.Context, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	listOpts, err := apiutils.BuildListOption(pageSize, orderBy, sortOrder, nextPageToken)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	result, err := s.coreApi.GetRegisteredModels(listOpts)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

// GetServingEnvironment - Get a ServingEnvironment
func (s *ModelRegistryServiceAPIService) GetServingEnvironment(ctx context.Context, servingenvironmentId string) (ImplResponse, error) {
	result, err := s.coreApi.GetServingEnvironmentById(servingenvironmentId)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO: return Response(http.StatusUnauthorized, Error{}), nil
}

// GetServingEnvironments - List All ServingEnvironments
func (s *ModelRegistryServiceAPIService) GetServingEnvironments(ctx context.Context, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	listOpts, err := apiutils.BuildListOption(pageSize, orderBy, sortOrder, nextPageToken)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	result, err := s.coreApi.GetServingEnvironments(listOpts)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

// UpdateInferenceService - Update a InferenceService
func (s *ModelRegistryServiceAPIService) UpdateInferenceService(ctx context.Context, inferenceserviceId string, inferenceServiceUpdate model.InferenceServiceUpdate) (ImplResponse, error) {
	entity, err := s.converter.ConvertInferenceServiceUpdate(&inferenceServiceUpdate)
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}
	entity.Id = &inferenceserviceId
	existing, err := s.coreApi.GetInferenceServiceById(inferenceserviceId)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	update, err := s.reconciler.UpdateExistingInferenceService(converter.NewOpenapiUpdateWrapper(existing, entity))
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}
	result, err := s.coreApi.UpsertInferenceService(&update)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

// UpdateArtifact - Update a Artifact
func (s *ModelRegistryServiceAPIService) UpdateArtifact(ctx context.Context, artifactId string, artifactUpdate model.ArtifactUpdate) (ImplResponse, error) {
	entity, err := s.converter.ConvertArtifactUpdate(&artifactUpdate)
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}
	if artifactUpdate.DocArtifactUpdate != nil {
		entity.DocArtifact.Id = &artifactId
	} else {
		entity.ModelArtifact.Id = &artifactId
	}
	existing, err := s.coreApi.GetArtifactById(artifactId)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	update, err := s.reconciler.UpdateExistingArtifact(converter.NewOpenapiUpdateWrapper(existing, entity))
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}
	result, err := s.coreApi.UpsertArtifact(&update)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

// UpdateModelArtifact - Update a ModelArtifact
func (s *ModelRegistryServiceAPIService) UpdateModelArtifact(ctx context.Context, modelartifactId string, modelArtifactUpdate model.ModelArtifactUpdate) (ImplResponse, error) {
	modelArtifact, err := s.converter.ConvertModelArtifactUpdate(&modelArtifactUpdate)
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}
	modelArtifact.Id = &modelartifactId
	existing, err := s.coreApi.GetModelArtifactById(modelartifactId)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	update, err := s.reconciler.UpdateExistingModelArtifact(converter.NewOpenapiUpdateWrapper(existing, modelArtifact))
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}
	result, err := s.coreApi.UpsertModelArtifact(&update)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

// UpdateModelVersion - Update a ModelVersion
func (s *ModelRegistryServiceAPIService) UpdateModelVersion(ctx context.Context, modelversionId string, modelVersionUpdate model.ModelVersionUpdate) (ImplResponse, error) {
	modelVersion, err := s.converter.ConvertModelVersionUpdate(&modelVersionUpdate)
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}
	modelVersion.Id = &modelversionId
	existing, err := s.coreApi.GetModelVersionById(modelversionId)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	update, err := s.reconciler.UpdateExistingModelVersion(converter.NewOpenapiUpdateWrapper(existing, modelVersion))
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}
	result, err := s.coreApi.UpsertModelVersion(&update, nil)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

// UpdateRegisteredModel - Update a RegisteredModel
func (s *ModelRegistryServiceAPIService) UpdateRegisteredModel(ctx context.Context, registeredmodelId string, registeredModelUpdate model.RegisteredModelUpdate) (ImplResponse, error) {
	registeredModel, err := s.converter.ConvertRegisteredModelUpdate(&registeredModelUpdate)
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}
	registeredModel.Id = &registeredmodelId
	existing, err := s.coreApi.GetRegisteredModelById(registeredmodelId)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	update, err := s.reconciler.UpdateExistingRegisteredModel(converter.NewOpenapiUpdateWrapper(existing, registeredModel))
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}
	result, err := s.coreApi.UpsertRegisteredModel(&update)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

// UpdateServingEnvironment - Update a ServingEnvironment
func (s *ModelRegistryServiceAPIService) UpdateServingEnvironment(ctx context.Context, servingenvironmentId string, servingEnvironmentUpdate model.ServingEnvironmentUpdate) (ImplResponse, error) {
	entity, err := s.converter.ConvertServingEnvironmentUpdate(&servingEnvironmentUpdate)
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}
	entity.Id = &servingenvironmentId
	existing, err := s.coreApi.GetServingEnvironmentById(servingenvironmentId)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	update, err := s.reconciler.UpdateExistingServingEnvironment(converter.NewOpenapiUpdateWrapper(existing, entity))
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}
	result, err := s.coreApi.UpsertServingEnvironment(&update)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}
	return Response(http.StatusOK, result), nil
	// TODO return Response(http.StatusUnauthorized, Error{}), nil
}

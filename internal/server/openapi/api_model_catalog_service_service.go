/*
 * Model Registry REST API
 *
 * REST API for Model Registry to create and manage ML model metadata
 *
 * API version: v1alpha3
 * Generated initially by: OpenAPI Generator (https://openapi-generator.tech).
 */

package openapi

import (
	"context"
	"errors"
	"net/http"

	model "github.com/kubeflow/model-registry/pkg/openapi"
)

// ModelCatalogServiceAPIService is a service that implements the logic for the ModelCatalogServiceAPIServicer
// This service should implement the business logic for every endpoint for the ModelCatalogServiceAPI API.
// Include any external packages or services that will be required by this service.
type ModelCatalogServiceAPIService struct {
}

// NewModelCatalogServiceAPIService creates a default api service
func NewModelCatalogServiceAPIService() ModelCatalogServiceAPIServicer {
	return &ModelCatalogServiceAPIService{}
}

// ApiModelCatalogV1alpha3SourcesSourceIdModelsModelIdReadmeGet -
func (s *ModelCatalogServiceAPIService) ApiModelCatalogV1alpha3SourcesSourceIdModelsModelIdReadmeGet(ctx context.Context, sourceId string, modelId string) (ImplResponse, error) {
	// TODO - update ApiModelCatalogV1alpha3SourcesSourceIdModelsModelIdReadmeGet with the required logic for this service method.
	// Add api_model_catalog_service_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	// return Response(200, nil),nil

	// TODO: Uncomment the next line to return response Response(401, Error{}) or use other options such as http.Ok ...
	// return Response(401, Error{}), nil

	// TODO: Uncomment the next line to return response Response(404, Error{}) or use other options such as http.Ok ...
	// return Response(404, Error{}), nil

	// TODO: Uncomment the next line to return response Response(500, Error{}) or use other options such as http.Ok ...
	// return Response(500, Error{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("ApiModelCatalogV1alpha3SourcesSourceIdModelsModelIdReadmeGet method not implemented")
}

// GetAllCatalogModels -
func (s *ModelCatalogServiceAPIService) GetAllCatalogModels(ctx context.Context, source string, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, offset string) (ImplResponse, error) {
	// TODO - update GetAllCatalogModels with the required logic for this service method.
	// Add api_model_catalog_service_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, CatalogModelList{}) or use other options such as http.Ok ...
	// return Response(200, CatalogModelList{}), nil

	// TODO: Uncomment the next line to return response Response(400, Error{}) or use other options such as http.Ok ...
	// return Response(400, Error{}), nil

	// TODO: Uncomment the next line to return response Response(401, Error{}) or use other options such as http.Ok ...
	// return Response(401, Error{}), nil

	// TODO: Uncomment the next line to return response Response(404, Error{}) or use other options such as http.Ok ...
	// return Response(404, Error{}), nil

	// TODO: Uncomment the next line to return response Response(500, Error{}) or use other options such as http.Ok ...
	// return Response(500, Error{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("GetAllCatalogModels method not implemented")
}

// GetCatalogModel -
func (s *ModelCatalogServiceAPIService) GetCatalogModel(ctx context.Context, sourceId string, modelId string) (ImplResponse, error) {
	// TODO - update GetCatalogModel with the required logic for this service method.
	// Add api_model_catalog_service_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, CatalogModel{}) or use other options such as http.Ok ...
	// return Response(200, CatalogModel{}), nil

	// TODO: Uncomment the next line to return response Response(401, Error{}) or use other options such as http.Ok ...
	// return Response(401, Error{}), nil

	// TODO: Uncomment the next line to return response Response(404, Error{}) or use other options such as http.Ok ...
	// return Response(404, Error{}), nil

	// TODO: Uncomment the next line to return response Response(500, Error{}) or use other options such as http.Ok ...
	// return Response(500, Error{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("GetCatalogModel method not implemented")
}

// GetCatalogSources - List All CatalogSources
func (s *ModelCatalogServiceAPIService) GetCatalogSources(ctx context.Context, name string, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, offset string) (ImplResponse, error) {
	// TODO - update GetCatalogSources with the required logic for this service method.
	// Add api_model_catalog_service_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, CatalogSourceList{}) or use other options such as http.Ok ...
	// return Response(200, CatalogSourceList{}), nil

	// TODO: Uncomment the next line to return response Response(400, Error{}) or use other options such as http.Ok ...
	// return Response(400, Error{}), nil

	// TODO: Uncomment the next line to return response Response(401, Error{}) or use other options such as http.Ok ...
	// return Response(401, Error{}), nil

	// TODO: Uncomment the next line to return response Response(404, Error{}) or use other options such as http.Ok ...
	// return Response(404, Error{}), nil

	// TODO: Uncomment the next line to return response Response(500, Error{}) or use other options such as http.Ok ...
	// return Response(500, Error{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("GetCatalogSources method not implemented")
}

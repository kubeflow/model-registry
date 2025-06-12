package core

import (
	"fmt"
	"strconv"

	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

func (b *ModelRegistryService) UpsertServingEnvironment(servingEnvironment *openapi.ServingEnvironment) (*openapi.ServingEnvironment, error) {
	if servingEnvironment == nil {
		return nil, fmt.Errorf("invalid serving environment pointer, cannot be nil: %w", api.ErrBadRequest)
	}

	servEnv, err := b.mapper.MapFromServingEnvironment(servingEnvironment)
	if err != nil {
		return nil, err
	}

	savedServEnv, err := b.servingEnvironmentRepository.Save(servEnv)
	if err != nil {
		return nil, err
	}

	return b.mapper.MapToServingEnvironment(savedServEnv)
}

func (b *ModelRegistryService) GetServingEnvironmentById(id string) (*openapi.ServingEnvironment, error) {
	convertedId, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	model, err := b.servingEnvironmentRepository.GetByID(int32(convertedId))
	if err != nil {
		return nil, err
	}

	return b.mapper.MapToServingEnvironment(model)
}

func (b *ModelRegistryService) GetServingEnvironmentByParams(name *string, externalId *string) (*openapi.ServingEnvironment, error) {
	servEnvsList, err := b.servingEnvironmentRepository.List(models.ServingEnvironmentListOptions{
		Name:       name,
		ExternalID: externalId,
	})
	if err != nil {
		return nil, err
	}

	if len(servEnvsList.Items) == 0 {
		return nil, fmt.Errorf("no serving environment found")
	}

	return b.mapper.MapToServingEnvironment(servEnvsList.Items[0])
}

func (b *ModelRegistryService) GetServingEnvironments(listOptions api.ListOptions) (*openapi.ServingEnvironmentList, error) {
	servEnvsList, err := b.servingEnvironmentRepository.List(models.ServingEnvironmentListOptions{
		Pagination: models.Pagination{
			PageSize:      listOptions.PageSize,
			OrderBy:       listOptions.OrderBy,
			SortOrder:     listOptions.SortOrder,
			NextPageToken: listOptions.NextPageToken,
		},
	})
	if err != nil {
		return nil, err
	}

	servingEnvironmentList := &openapi.ServingEnvironmentList{
		Items: []openapi.ServingEnvironment{},
	}

	for _, servEnv := range servEnvsList.Items {
		servingEnvironment, err := b.mapper.MapToServingEnvironment(servEnv)
		if err != nil {
			return nil, err
		}
		servingEnvironmentList.Items = append(servingEnvironmentList.Items, *servingEnvironment)
	}

	servingEnvironmentList.NextPageToken = servEnvsList.NextPageToken
	servingEnvironmentList.PageSize = servEnvsList.PageSize
	servingEnvironmentList.Size = int32(servEnvsList.Size)

	return servingEnvironmentList, nil
}

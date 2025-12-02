package core

import (
	"errors"
	"fmt"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"gorm.io/gorm"
)

func (b *ModelRegistryService) UpsertExperiment(experiment *openapi.Experiment) (*openapi.Experiment, error) {
	if experiment == nil {
		return nil, fmt.Errorf("invalid experiment pointer, can't upsert nil: %w", api.ErrBadRequest)
	}

	if experiment.Id != nil {
		existing, err := b.GetExperimentById(*experiment.Id)
		if err != nil {
			return nil, err
		}

		// Use OpenAPIReconciler for proper merging instead of incomplete OverrideNotEditableForExperiment
		withNotEditable, err := b.mapper.UpdateExistingExperiment(converter.NewOpenapiUpdateWrapper(existing, experiment))
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		// Handle CustomProperties preservation for partial updates
		// If the update didn't specify CustomProperties (nil), preserve existing ones
		if experiment.CustomProperties == nil && existing.CustomProperties != nil {
			withNotEditable.CustomProperties = existing.CustomProperties
		}

		experiment = &withNotEditable
	}

	experimentEntity, err := b.mapper.MapFromExperiment(experiment)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	experimentEntity, err = b.experimentRepository.Save(experimentEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, fmt.Errorf("experiment with name %s already exists: %w", experiment.Name, api.ErrConflict)
		}

		return nil, err
	}

	toReturn, err := b.mapper.MapToExperiment(experimentEntity)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return toReturn, nil
}

func (b *ModelRegistryService) GetExperimentById(id string) (*openapi.Experiment, error) {
	convertedId, err := apiutils.ValidateIDAsInt32(id, "experiment")
	if err != nil {
		return nil, err
	}

	experiment, err := b.experimentRepository.GetByID(convertedId)
	if err != nil {
		return nil, fmt.Errorf("no experiment found for id %s: %w", id, api.ErrNotFound)
	}

	toReturn, err := b.mapper.MapToExperiment(experiment)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return toReturn, nil
}

func (b *ModelRegistryService) GetExperimentByParams(name *string, externalId *string) (*openapi.Experiment, error) {
	if name == nil && externalId == nil {
		return nil, fmt.Errorf("invalid parameters call, supply either name or externalId: %w", api.ErrBadRequest)
	}

	experiments, err := b.experimentRepository.List(models.ExperimentListOptions{
		Name:       name,
		ExternalID: externalId,
	})
	if err != nil {
		return nil, err
	}

	if len(experiments.Items) == 0 {
		return nil, fmt.Errorf("no experiments found for name=%v, externalId=%v: %w", apiutils.ZeroIfNil(name), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	if len(experiments.Items) > 1 {
		return nil, fmt.Errorf("multiple experiments found for name=%v, externalId=%v: %w", apiutils.ZeroIfNil(name), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	toReturn, err := b.mapper.MapToExperiment(experiments.Items[0])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return toReturn, nil
}

func (b *ModelRegistryService) GetExperiments(listOptions api.ListOptions) (*openapi.ExperimentList, error) {
	experiments, err := b.experimentRepository.List(models.ExperimentListOptions{
		Pagination: models.Pagination{
			PageSize:      listOptions.PageSize,
			OrderBy:       listOptions.OrderBy,
			SortOrder:     listOptions.SortOrder,
			NextPageToken: listOptions.NextPageToken,
			FilterQuery:   listOptions.FilterQuery,
		},
	})
	if err != nil {
		return nil, err
	}

	experimentList := &openapi.ExperimentList{
		Items: []openapi.Experiment{},
	}

	for _, experiment := range experiments.Items {
		experiment, err := b.mapper.MapToExperiment(experiment)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		experimentList.Items = append(experimentList.Items, *experiment)
	}

	experimentList.NextPageToken = experiments.NextPageToken
	experimentList.PageSize = experiments.PageSize
	experimentList.Size = int32(experiments.Size)

	return experimentList, nil
}

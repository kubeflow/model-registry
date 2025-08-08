package core

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/defaults"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"gorm.io/gorm"
)

func (b *ModelRegistryService) UpsertExperimentRun(experimentRun *openapi.ExperimentRun, experimentId *string) (*openapi.ExperimentRun, error) {
	if experimentRun == nil {
		return nil, fmt.Errorf("invalid experiment run pointer, can't upsert nil: %w", api.ErrBadRequest)
	}

	if experimentId == nil {
		return nil, fmt.Errorf("experiment ID is required: %w", api.ErrBadRequest)
	}

	// Convert experiment ID to int32
	experimentIDPtr, err := apiutils.ValidateIDAsInt32(*experimentId, "experiment")
	if err != nil {
		return nil, err
	}

	// Validate that the experiment exists
	_, err = b.GetExperimentById(*experimentId)
	if err != nil {
		return nil, fmt.Errorf("experiment not found: %w", err)
	}

	// Set the ExperimentId field on the experimentRun object (required for mapper)
	experimentRun.ExperimentId = *experimentId

	if experimentRun.Id != nil {
		// Update existing experiment run
		existing, err := b.GetExperimentRunById(*experimentRun.Id)
		if err != nil {
			return nil, err
		}

		withNotEditable, err := b.mapper.UpdateExistingExperimentRun(converter.NewOpenapiUpdateWrapper(existing, experimentRun))
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		experimentRun = &withNotEditable
	}

	// Validate that EndTimeSinceEpoch is not less than StartTimeSinceEpoch
	if experimentRun.EndTimeSinceEpoch != nil && experimentRun.StartTimeSinceEpoch != nil {
		endTime, err := strconv.ParseInt(*experimentRun.EndTimeSinceEpoch, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid EndTimeSinceEpoch value: %v: %w", err, api.ErrBadRequest)
		}

		startTime, err := strconv.ParseInt(*experimentRun.StartTimeSinceEpoch, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid StartTimeSinceEpoch value: %v: %w", err, api.ErrBadRequest)
		}

		if endTime < startTime {
			return nil, fmt.Errorf("EndTimeSinceEpoch (%d) cannot be less than StartTimeSinceEpoch (%d): %w", endTime, startTime, api.ErrBadRequest)
		}
	}

	experimentRunEntity, err := b.mapper.MapFromExperimentRun(experimentRun, experimentId)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	experimentRunEntity, err = b.experimentRunRepository.Save(experimentRunEntity, &experimentIDPtr)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, fmt.Errorf("experiment run with name %s already exists: %w", *experimentRun.Name, api.ErrConflict)
		}

		return nil, err
	}

	toReturn, err := b.mapper.MapToExperimentRun(experimentRunEntity)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return toReturn, nil
}

func (b *ModelRegistryService) GetExperimentRunById(id string) (*openapi.ExperimentRun, error) {
	convertedId, err := apiutils.ValidateIDAsInt32(id, "experiment run")
	if err != nil {
		return nil, err
	}

	experimentRun, err := b.experimentRunRepository.GetByID(convertedId)
	if err != nil {
		return nil, fmt.Errorf("no experiment run found for id %s: %w", id, api.ErrNotFound)
	}

	toReturn, err := b.mapper.MapToExperimentRun(experimentRun)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return toReturn, nil
}

func (b *ModelRegistryService) GetExperimentRunByParams(name *string, experimentId *string, externalId *string) (*openapi.ExperimentRun, error) {
	if (name == nil || experimentId == nil) && externalId == nil {
		return nil, fmt.Errorf("invalid parameters call, supply either (name and experimentId), or externalId: %w", api.ErrBadRequest)
	}

	var experimentIDPtr *int32
	if experimentId != nil {
		var err error
		experimentIDPtr, err = apiutils.ValidateIDAsInt32Ptr(experimentId, "experiment")
		if err != nil {
			return nil, err
		}
	}

	experimentRuns, err := b.experimentRunRepository.List(models.ExperimentRunListOptions{
		Name:         name,
		ExternalID:   externalId,
		ExperimentID: experimentIDPtr,
	})
	if err != nil {
		return nil, err
	}

	if len(experimentRuns.Items) == 0 {
		return nil, fmt.Errorf("no experiment runs found for name=%v, experimentId=%v, externalId=%v: %w", apiutils.ZeroIfNil(name), apiutils.ZeroIfNil(experimentId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	if len(experimentRuns.Items) > 1 {
		return nil, fmt.Errorf("multiple experiment runs found for name=%v, experimentId=%v, externalId=%v: %w", apiutils.ZeroIfNil(name), apiutils.ZeroIfNil(experimentId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	toReturn, err := b.mapper.MapToExperimentRun(experimentRuns.Items[0])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return toReturn, nil
}

func (b *ModelRegistryService) GetExperimentRuns(listOptions api.ListOptions, experimentId *string) (*openapi.ExperimentRunList, error) {
	var experimentIDPtr *int32
	if experimentId != nil {
		var err error
		experimentIDPtr, err = apiutils.ValidateIDAsInt32Ptr(experimentId, "experiment")
		if err != nil {
			return nil, err
		}

		// Validate that the experiment exists
		_, err = b.GetExperimentById(*experimentId)
		if err != nil {
			return nil, err
		}
	}

	experimentRuns, err := b.experimentRunRepository.List(models.ExperimentRunListOptions{
		Pagination: models.Pagination{
			PageSize:      listOptions.PageSize,
			OrderBy:       listOptions.OrderBy,
			SortOrder:     listOptions.SortOrder,
			NextPageToken: listOptions.NextPageToken,
			FilterQuery:   listOptions.FilterQuery,
		},
		ExperimentID: experimentIDPtr,
	})
	if err != nil {
		return nil, err
	}

	experimentRunList := &openapi.ExperimentRunList{
		Items: []openapi.ExperimentRun{},
	}

	for _, experimentRun := range experimentRuns.Items {
		experimentRun, err := b.mapper.MapToExperimentRun(experimentRun)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		experimentRunList.Items = append(experimentRunList.Items, *experimentRun)
	}

	experimentRunList.NextPageToken = experimentRuns.NextPageToken
	experimentRunList.PageSize = experimentRuns.PageSize
	experimentRunList.Size = int32(experimentRuns.Size)

	return experimentRunList, nil
}

func (b *ModelRegistryService) UpsertExperimentRunArtifact(artifact *openapi.Artifact, experimentRunId string) (*openapi.Artifact, error) {
	result, err := b.upsertArtifact(artifact, &experimentRunId)
	if err != nil {
		return nil, err
	}

	// Only create metric history for experiment runs, and only if it's a metric
	if result.Metric != nil {
		err = b.InsertMetricHistory(result.Metric, experimentRunId)
		if err != nil {
			return nil, fmt.Errorf("failed to insert metric history for metric %s: %w", *result.Metric.Name, err)
		}
	}

	return result, nil
}

func (b *ModelRegistryService) GetExperimentRunArtifacts(artifactType openapi.ArtifactTypeQueryParam, listOptions api.ListOptions, experimentRunId *string) (*openapi.ArtifactList, error) {
	// Note: artifactType parameter is not used in the EmbedMD implementation
	// This matches the pattern used by other artifact methods in the EmbedMD service
	return b.GetArtifacts(artifactType, listOptions, experimentRunId)
}

func (b *ModelRegistryService) GetExperimentRunMetricHistory(name *string, stepIds *string, listOptions api.ListOptions, experimentRunId *string) (*openapi.MetricList, error) {

	// Validate experiment run exists
	if experimentRunId == nil {
		return nil, fmt.Errorf("experiment run ID is required: %w", api.ErrBadRequest)
	}

	// Validate experiment run exists
	_, err := b.GetExperimentRunById(*experimentRunId)
	if err != nil {
		return nil, fmt.Errorf("experiment run not found: %w", err)
	}

	// Convert experiment run ID to int32 for repository queries
	experimentRunIdInt32, err := apiutils.ValidateIDAsInt32(*experimentRunId, "experiment run")
	if err != nil {
		return nil, err
	}

	listOptsCopy := models.MetricHistoryListOptions{
		Pagination: models.Pagination{
			PageSize:      listOptions.PageSize,
			OrderBy:       listOptions.OrderBy,
			SortOrder:     listOptions.SortOrder,
			NextPageToken: listOptions.NextPageToken,
		},
		ExperimentRunID: &experimentRunIdInt32,
	}

	// Add name filter if provided
	if name != nil && *name != "" {
		listOptsCopy.Name = name
	}

	// Add step IDs filter if provided
	if stepIds != nil && *stepIds != "" {
		listOptsCopy.StepIds = stepIds
	}

	// Query metric history repository
	metricHistories, err := b.metricHistoryRepository.List(listOptsCopy)
	if err != nil {
		return nil, err
	}

	// Convert metric history to OpenAPI metrics
	results := []openapi.Metric{}
	for _, metricHistory := range metricHistories.Items {

		mapped, err := b.mapper.MapToMetric(metricHistory)
		if err != nil {
			glog.Warningf("Failed to map metric history artifact %v: %v", metricHistory.GetID(), err)
			continue
		}

		// Remove timestamp suffix from name after the last __
		if mapped.Name != nil {
			lastIndex := strings.LastIndex(*mapped.Name, "__")
			if lastIndex != -1 {
				mapped.Name = apiutils.StrPtr((*mapped.Name)[:lastIndex])
			}
		}

		results = append(results, *mapped)
	}

	// Build response
	toReturn := openapi.MetricList{
		NextPageToken: metricHistories.NextPageToken,
		PageSize:      metricHistories.PageSize,
		Size:          int32(len(results)),
		Items:         results,
	}

	return &toReturn, nil
}

// InsertMetricHistory inserts a metric history record for an experiment run
func (b *ModelRegistryService) InsertMetricHistory(metric *openapi.Metric, experimentRunId string) error {

	if metric == nil {
		return fmt.Errorf("metric cannot be nil: %w", api.ErrBadRequest)
	}

	if experimentRunId == "" {
		return fmt.Errorf("experiment run ID is required: %w", api.ErrBadRequest)
	}

	// Validate that the experiment run exists
	_, err := b.GetExperimentRunById(experimentRunId)
	if err != nil {
		return fmt.Errorf("experiment run not found: %w", err)
	}

	// Convert experiment run ID to int32
	experimentRunIdInt32, err := apiutils.ValidateIDAsInt32(experimentRunId, "experiment run")
	if err != nil {
		return err
	}

	// Create a copy of the metric for history
	metricHistory := *metric

	// Generate a unique name with the metric last update time since epoch to avoid duplicates
	timestamp := metric.LastUpdateTimeSinceEpoch
	if timestamp == nil {
		// fallback to the current time if the last update time since epoch is not set
		timestamp = apiutils.Of(strconv.FormatInt(time.Now().UnixMilli(), 10))
	}
	metricHistory.Name = apiutils.StrPtr(fmt.Sprintf("%s:%s__%s", experimentRunId, *metric.Name, *timestamp))

	// Set the ID to nil to create a new artifact
	metricHistory.Id = nil
	// Clear the external ID to avoid duplicate key error
	metricHistory.ExternalId = nil

	// Create the MetricHistory entity with the correct TypeID
	// Get the metric history type ID from the types map
	metricHistoryTypeID, exists := b.typesMap[defaults.MetricHistoryTypeName]
	if !exists {
		return fmt.Errorf("metric history type not found in types map")
	}

	metricHistoryEntity := &models.MetricHistoryImpl{
		TypeID: apiutils.Of(int32(metricHistoryTypeID)),
		Attributes: &models.MetricHistoryAttributes{
			Name:         metricHistory.Name,
			URI:          nil, // Metric doesn't have URI field
			State:        (*string)(metricHistory.State),
			ArtifactType: apiutils.StrPtr(models.MetricHistoryType),
			ExternalID:   metricHistory.ExternalId,
		},
	}

	// Map properties from metric to metric history using the converter
	metricProperties, err := converter.MapMetricPropertiesEmbedMD(metric)
	if err != nil {
		return fmt.Errorf("failed to map metric properties: %w", err)
	}
	metricHistoryEntity.Properties = metricProperties

	// Handle custom properties using the converter
	if metric.CustomProperties != nil {
		customProps, err := converter.MapOpenAPICustomPropertiesEmbedMD(metric.CustomProperties)
		if err != nil {
			return fmt.Errorf("failed to map custom properties: %w", err)
		}
		metricHistoryEntity.CustomProperties = customProps
	}

	// Save the metric history
	_, err = b.metricHistoryRepository.Save(metricHistoryEntity, &experimentRunIdInt32)
	if err != nil {
		return fmt.Errorf("failed to insert metric history: %w", err)
	}

	glog.Infof("Successfully inserted metric history for metric %s in experiment run %s", *metric.Name, experimentRunId)
	return nil
}

package core

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/ml_metadata/proto"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

// GetExperimentRunMetricHistory retrieves metric history for an experiment run
func (serv *ModelRegistryService) GetExperimentRunMetricHistory(name *string, stepIds *string, listOptions api.ListOptions, experimentRunId *string) (*openapi.ArtifactList, error) {
	// Validate experiment run exists
	if experimentRunId == nil {
		return nil, fmt.Errorf("experiment run ID is required: %w", api.ErrBadRequest)
	}
	// Validate experiment run exists
	_, err := serv.GetExperimentRunById(*experimentRunId)
	if err != nil {
		return nil, fmt.Errorf("experiment run not found: %w", err)
	}

	// Build list operation options
	listOperationOptions, err := apiutils.BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	// Build filter query for metrics
	filterQuery := fmt.Sprintf("contexts_a.id = %s", *experimentRunId)

	// Add name filter if provided
	if name != nil && *name != "" {
		// search for metric name with a wildcard for timestamp suffix for metric history
		nameFilter := fmt.Sprintf("name LIKE \"%s:%s__%%\"", *experimentRunId, *name)
		filterQuery = fmt.Sprintf("(%s) AND (%s)", filterQuery, nameFilter)
	}

	// Add step IDs filter if provided
	if stepIds != nil && *stepIds != "" {
		stepFilter := fmt.Sprintf("properties.step.int_value IN (%s)", *stepIds)
		filterQuery = fmt.Sprintf("(%s) AND (%s)", filterQuery, stepFilter)
	}

	// Apply the filter query
	if listOperationOptions.FilterQuery != nil {
		existingFilter := *listOperationOptions.FilterQuery
		filterQuery = fmt.Sprintf("(%s) AND (%s)", existingFilter, filterQuery)
	}
	listOperationOptions.FilterQuery = &filterQuery

	// Query MLMD for metrics
	resp, err := serv.mlmdClient.GetArtifactsByType(context.Background(), &proto.GetArtifactsByTypeRequest{
		TypeName: &serv.nameConfig.MetricHistoryTypeName,
		Options:  listOperationOptions,
	})
	if err != nil {
		return nil, err
	}

	// Convert MLMD artifacts to OpenAPI artifacts
	results := []openapi.Artifact{}
	for _, artifact := range resp.Artifacts {
		// Check if this is a metric history artifact
		if *artifact.Type == serv.nameConfig.MetricHistoryTypeName {
			mapped, err := serv.mapper.MapToMetricHistory(artifact)
			if err != nil {
				glog.Warningf("Failed to map metric artifact %d: %v", *artifact.Id, err)
				continue
			}
			// remove timestamp suffix from name after the last __
			lastIndex := strings.LastIndex(*mapped.Name, "__")
			mapped.Name = apiutils.StrPtr((*mapped.Name)[:lastIndex])
			results = append(results, openapi.Artifact{
				Metric: mapped,
			})
		}
	}

	// Build response
	toReturn := openapi.ArtifactList{
		NextPageToken: apiutils.ZeroIfNil(resp.NextPageToken),
		PageSize:      apiutils.ZeroIfNil(listOptions.PageSize),
		Size:          int32(len(results)),
		Items:         results,
	}

	return &toReturn, nil
}

// InsertMetricHistory inserts a metric history record for an experiment run
func (serv *ModelRegistryService) InsertMetricHistory(metric *openapi.Metric, experimentRunId string) error {
	if metric == nil {
		return fmt.Errorf("metric cannot be nil: %w", api.ErrBadRequest)
	}

	if experimentRunId == "" {
		return fmt.Errorf("experiment run ID is required: %w", api.ErrBadRequest)
	}

	// Validate that the experiment run exists
	_, err := serv.GetExperimentRunById(experimentRunId)
	if err != nil {
		return fmt.Errorf("experiment run not found: %w", err)
	}

	// Create a copy of the metric for history
	metricHistory := *metric

	// Generate a unique name with the metric last update time since epoch to avoid duplicates
	timestamp := metric.LastUpdateTimeSinceEpoch
	if timestamp == nil {
		// fallback to the current time if the last update time since epoch is not set
		timestamp = apiutils.Of(strconv.FormatInt(time.Now().UnixMilli(), 10))
	}
	metricHistory.Name = apiutils.StrPtr(fmt.Sprintf("%s__%s", *metric.Name, *timestamp))

	// Set the ID to nil to create a new artifact
	metricHistory.Id = nil

	// Map the metric to MLMD artifact
	mlmdArtifact, err := serv.mapper.MapFromMetricArtifact(&metricHistory, &experimentRunId)
	if err != nil {
		return fmt.Errorf("failed to map metric to MLMD artifact: %w", err)
	}

	// Set the type and typeId to MetricHistory
	mlmdArtifact.Type = &serv.nameConfig.MetricHistoryTypeName
	mlmdArtifact.TypeId = apiutils.Of(serv.typesMap[serv.nameConfig.MetricHistoryTypeName])

	// Insert the artifact into MLMD
	resp, err := serv.mlmdClient.PutArtifacts(context.Background(), &proto.PutArtifactsRequest{
		Artifacts: []*proto.Artifact{mlmdArtifact},
	})
	if err != nil {
		return fmt.Errorf("failed to insert metric history into MLMD: %w", err)
	}

	if len(resp.ArtifactIds) == 0 {
		return fmt.Errorf("no artifact ID returned from MLMD")
	}

	// Create association between the metric history artifact and the experiment run
	experimentRunIdInt, err := converter.StringToInt64(&experimentRunId)
	if err != nil {
		return fmt.Errorf("failed to convert experiment run ID to int64: %w", err)
	}

	_, err = serv.mlmdClient.PutAttributionsAndAssociations(context.Background(), &proto.PutAttributionsAndAssociationsRequest{
		Attributions: []*proto.Attribution{{
			ArtifactId: &resp.ArtifactIds[0],
			ContextId:  experimentRunIdInt,
		}},
	})
	if err != nil {
		return fmt.Errorf("failed to create attribution between metric history and experiment run: %w", err)
	}

	glog.Infof("Successfully inserted metric history for metric %s in experiment run %s", *metric.Name, experimentRunId)
	return nil
}

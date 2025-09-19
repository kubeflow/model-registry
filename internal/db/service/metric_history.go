package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/utils"
	"gorm.io/gorm"
)

var ErrMetricHistoryNotFound = errors.New("metric history by id not found")

// parseStepIds parses a comma-separated string of step IDs into a slice of integers
// Validation should have been done at the API layer, but we return error for defensive programming
func parseStepIds(stepIds string) ([]int32, error) {
	if stepIds == "" {
		return nil, nil
	}

	parts := strings.Split(stepIds, ",")
	result := make([]int32, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue // skip empty parts
		}

		// Parse with error checking for defensive programming
		stepID, err := strconv.ParseInt(part, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid step ID '%s' in repository layer (validation should have caught this): %w", part, err)
		}
		result = append(result, int32(stepID))
	}

	return result, nil
}

type MetricHistoryRepositoryImpl struct {
	*GenericRepository[models.MetricHistory, schema.Artifact, schema.ArtifactProperty, *models.MetricHistoryListOptions]
}

func NewMetricHistoryRepository(db *gorm.DB, typeID int64) models.MetricHistoryRepository {
	config := GenericRepositoryConfig[models.MetricHistory, schema.Artifact, schema.ArtifactProperty, *models.MetricHistoryListOptions]{
		DB:                  db,
		TypeID:              typeID,
		EntityToSchema:      mapMetricHistoryToArtifact,
		SchemaToEntity:      mapDataLayerToMetricHistory,
		EntityToProperties:  mapMetricHistoryToArtifactProperties,
		NotFoundError:       ErrMetricHistoryNotFound,
		EntityName:          "metric history",
		PropertyFieldName:   "artifact_id",
		ApplyListFilters:    applyMetricHistoryListFilters,
		IsNewEntity:         func(entity models.MetricHistory) bool { return entity.GetID() == nil },
		HasCustomProperties: func(entity models.MetricHistory) bool { return entity.GetCustomProperties() != nil },
	}

	return &MetricHistoryRepositoryImpl{
		GenericRepository: NewGenericRepository(config),
	}
}

func (r *MetricHistoryRepositoryImpl) List(listOptions models.MetricHistoryListOptions) (*models.ListWrapper[models.MetricHistory], error) {
	return r.GenericRepository.List(&listOptions)
}

func applyMetricHistoryListFilters(query *gorm.DB, listOptions *models.MetricHistoryListOptions) *gorm.DB {
	if listOptions.Name != nil {
		query = query.Where(utils.GetTableName(query, &schema.Artifact{})+".name LIKE ?", fmt.Sprintf("%%%s%%", *listOptions.Name))
	} else if listOptions.ExternalID != nil {
		query = query.Where(utils.GetTableName(query, &schema.Artifact{})+".external_id = ?", listOptions.ExternalID)
	}

	// Add step IDs filter if provided - use unique alias to avoid conflicts with filterQuery joins
	if listOptions.StepIds != nil && *listOptions.StepIds != "" {
		// Parse step IDs (validation should have been done at API layer)
		stepIds, err := parseStepIds(*listOptions.StepIds)
		if err != nil {
			// Log error but don't fail the query - this indicates a validation bug
			// that should be caught at the API layer
			glog.Errorf("Failed to parse stepIds in repository layer: %v", err)
			return query
		}

		if len(stepIds) > 0 {
			// Proper GORM JOIN: Use properly quoted table and column names
			stepPropsTable := utils.GetTableName(query, &schema.ArtifactProperty{}) + " AS step_props"
			artifactTable := utils.GetTableName(query, &schema.Artifact{})
			query = query.Joins(fmt.Sprintf("JOIN %s ON step_props.artifact_id = %s.id", stepPropsTable, artifactTable)).
				Where("step_props.name = ? AND step_props.int_value IN (?)",
					"step", stepIds)
		}
	}

	// Join with Attribution table only when filtering by experiment run ID
	if listOptions.ExperimentRunID != nil {
		// Proper GORM JOIN: Use helper that respects naming strategy
		query = query.Joins(utils.BuildAttributionJoin(query)).
			Where(utils.GetColumnRef(query, &schema.Attribution{}, "context_id")+" = ?", listOptions.ExperimentRunID)
	}

	return query
}

func mapMetricHistoryToArtifact(metricHistory models.MetricHistory) schema.Artifact {
	if metricHistory == nil {
		return schema.Artifact{}
	}

	artifact := schema.Artifact{
		ID:     apiutils.ZeroIfNil(metricHistory.GetID()),
		TypeID: apiutils.ZeroIfNil(metricHistory.GetTypeID()),
	}

	if metricHistory.GetAttributes() != nil {
		attrs := metricHistory.GetAttributes()
		artifact.Name = attrs.Name
		artifact.ExternalID = attrs.ExternalID
		artifact.URI = attrs.URI
		// Handle State conversion - MetricHistory uses string, schema.Artifact uses int32
		if attrs.State != nil {
			stateValue := models.Artifact_State_value[*attrs.State]
			artifact.State = &stateValue
		}
		artifact.CreateTimeSinceEpoch = apiutils.ZeroIfNil(attrs.CreateTimeSinceEpoch)
		artifact.LastUpdateTimeSinceEpoch = apiutils.ZeroIfNil(attrs.LastUpdateTimeSinceEpoch)
	}

	return artifact
}

func mapMetricHistoryToArtifactProperties(metricHistory models.MetricHistory, artifactID int32) []schema.ArtifactProperty {
	var properties []schema.ArtifactProperty

	if metricHistory.GetProperties() != nil {
		for _, prop := range *metricHistory.GetProperties() {
			properties = append(properties, MapPropertiesToArtifactProperty(prop, artifactID, false))
		}
	}

	if metricHistory.GetCustomProperties() != nil {
		for _, prop := range *metricHistory.GetCustomProperties() {
			properties = append(properties, MapPropertiesToArtifactProperty(prop, artifactID, true))
		}
	}

	return properties
}

func mapDataLayerToMetricHistory(metricHistory schema.Artifact, artProperties []schema.ArtifactProperty) models.MetricHistory {
	var state *string
	if metricHistory.State != nil {
		metricState := models.Artifact_State_name[*metricHistory.State]
		state = &metricState
	}

	metricHistoryModel := &models.BaseEntity[models.MetricHistoryAttributes]{
		ID:     &metricHistory.ID,
		TypeID: &metricHistory.TypeID,
		Attributes: &models.MetricHistoryAttributes{
			Name:                     metricHistory.Name,
			ExternalID:               metricHistory.ExternalID,
			URI:                      metricHistory.URI,
			State:                    state,
			ArtifactType:             apiutils.StrPtr(models.MetricHistoryType),
			CreateTimeSinceEpoch:     &metricHistory.CreateTimeSinceEpoch,
			LastUpdateTimeSinceEpoch: &metricHistory.LastUpdateTimeSinceEpoch,
		},
	}

	properties := []models.Properties{}
	customProperties := []models.Properties{}

	for _, prop := range artProperties {
		mappedProperty := MapArtifactPropertyToProperties(prop)

		if prop.IsCustomProperty {
			customProperties = append(customProperties, mappedProperty)
		} else {
			properties = append(properties, mappedProperty)
		}
	}

	// Always set Properties and CustomProperties, even if empty
	metricHistoryModel.Properties = &properties
	metricHistoryModel.CustomProperties = &customProperties

	return metricHistoryModel
}

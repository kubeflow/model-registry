package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/utils"
	"gorm.io/gorm"
)

var ErrMetricHistoryNotFound = errors.New("metric history by id not found")

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
		// Proper GORM JOIN: Use properly quoted table and column names
		stepPropsTable := utils.GetTableName(query, &schema.ArtifactProperty{}) + " AS step_props"
		artifactTable := utils.GetTableName(query, &schema.Artifact{})
		query = query.Joins(fmt.Sprintf("JOIN %s ON step_props.artifact_id = %s.id", stepPropsTable, artifactTable)).
			Where("step_props.name = ? AND step_props.int_value IN (?)",
				"step", strings.Split(*listOptions.StepIds, ","))
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

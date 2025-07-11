package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/scopes"
	"gorm.io/gorm"
)

var ErrMetricHistoryNotFound = errors.New("metric history by id not found")

type MetricHistoryRepositoryImpl struct {
	db     *gorm.DB
	typeID int64
}

func NewMetricHistoryRepository(db *gorm.DB, typeID int64) models.MetricHistoryRepository {
	return &MetricHistoryRepositoryImpl{db: db, typeID: typeID}
}

func (r *MetricHistoryRepositoryImpl) GetByID(id int32) (models.MetricHistory, error) {
	metricHistory := &schema.Artifact{}
	properties := []schema.ArtifactProperty{}

	if err := r.db.Where("id = ? AND type_id = ?", id, r.typeID).First(metricHistory).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %v", ErrMetricHistoryNotFound, err)
		}

		return nil, fmt.Errorf("error getting metric history by id: %w", err)
	}

	if err := r.db.Where("artifact_id = ?", metricHistory.ID).Find(&properties).Error; err != nil {
		return nil, fmt.Errorf("error getting properties by metric history id: %w", err)
	}

	return mapDataLayerToMetricHistory(*metricHistory, properties), nil
}

func (r *MetricHistoryRepositoryImpl) Save(metricHistory models.MetricHistory, experimentRunID *int32) (models.MetricHistory, error) {
	now := time.Now().UnixMilli()

	metricHistoryArt := mapMetricHistoryToArtifact(metricHistory)
	propertiesArt := []schema.ArtifactProperty{}

	metricHistoryArt.LastUpdateTimeSinceEpoch = now

	if metricHistory.GetID() == nil {
		metricHistoryArt.CreateTimeSinceEpoch = now
	}

	hasCustomProperties := metricHistory.GetCustomProperties() != nil

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&metricHistoryArt).Error; err != nil {
			return fmt.Errorf("error saving metric history: %w", err)
		}

		propertiesArt = mapMetricHistoryToArtifactProperties(metricHistory, metricHistoryArt.ID)
		existingCustomPropertiesArt := []schema.ArtifactProperty{}

		if err := tx.Where("artifact_id = ? AND is_custom_property = ?", metricHistoryArt.ID, true).Find(&existingCustomPropertiesArt).Error; err != nil {
			return fmt.Errorf("error getting existing custom properties by metric history id: %w", err)
		}

		if hasCustomProperties {
			for _, existingProp := range existingCustomPropertiesArt {
				found := false
				for _, prop := range propertiesArt {
					if prop.Name == existingProp.Name && prop.ArtifactID == existingProp.ArtifactID && prop.IsCustomProperty == existingProp.IsCustomProperty {
						found = true
						break
					}
				}

				if !found {
					if err := tx.Delete(&existingProp).Error; err != nil {
						return fmt.Errorf("error deleting metric history property: %w", err)
					}
				}
			}
		}

		for _, prop := range propertiesArt {
			var existingProp schema.ArtifactProperty
			result := tx.Where("artifact_id = ? AND name = ? AND is_custom_property = ?",
				prop.ArtifactID, prop.Name, prop.IsCustomProperty).First(&existingProp)

			switch result.Error {
			case nil:
				prop.ArtifactID = existingProp.ArtifactID
				prop.Name = existingProp.Name
				prop.IsCustomProperty = existingProp.IsCustomProperty
				if err := tx.Model(&existingProp).Updates(prop).Error; err != nil {
					return fmt.Errorf("error updating metric history property: %w", err)
				}
			case gorm.ErrRecordNotFound:
				if err := tx.Create(&prop).Error; err != nil {
					return fmt.Errorf("error creating metric history property: %w", err)
				}
			default:
				return fmt.Errorf("error checking existing property: %w", result.Error)
			}
		}

		if experimentRunID != nil {
			// Check if attribution already exists to avoid duplicate key errors
			var existingAttribution schema.Attribution
			result := tx.Where("context_id = ? AND artifact_id = ?", *experimentRunID, metricHistoryArt.ID).First(&existingAttribution)

			if result.Error == gorm.ErrRecordNotFound {
				// Attribution doesn't exist, create it
				attribution := schema.Attribution{
					ContextID:  *experimentRunID,
					ArtifactID: metricHistoryArt.ID,
				}

				if err := tx.Create(&attribution).Error; err != nil {
					return fmt.Errorf("error creating attribution: %w", err)
				}
			} else if result.Error != nil {
				return fmt.Errorf("error checking existing attribution: %w", result.Error)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return mapDataLayerToMetricHistory(metricHistoryArt, propertiesArt), nil
}

func (r *MetricHistoryRepositoryImpl) List(listOptions models.MetricHistoryListOptions) (*models.ListWrapper[models.MetricHistory], error) {
	list := models.ListWrapper[models.MetricHistory]{
		PageSize: listOptions.GetPageSize(),
	}

	metricHistories := []models.MetricHistory{}
	metricHistoriesArt := []schema.Artifact{}

	query := r.db.Model(&schema.Artifact{}).Where("Artifact.type_id = ?", r.typeID)

	if listOptions.Name != nil {
		query = query.Where("Artifact.name LIKE ?", fmt.Sprintf("%%%s%%", *listOptions.Name))
	} else if listOptions.ExternalID != nil {
		query = query.Where("Artifact.external_id = ?", listOptions.ExternalID)
	}

	// Add step IDs filter if provided
	if listOptions.StepIds != nil && *listOptions.StepIds != "" {
		query = query.Joins("JOIN ArtifactProperty ON ArtifactProperty.artifact_id = Artifact.id").
			Where("ArtifactProperty.name = ? AND ArtifactProperty.int_value IN (?)",
				"step", strings.Split(*listOptions.StepIds, ","))
	}

	if listOptions.ExperimentRunID != nil {
		query = query.Joins("JOIN Attribution ON Attribution.artifact_id = Artifact.id").
			Where("Attribution.context_id = ?", listOptions.ExperimentRunID)
		query = query.Scopes(scopes.PaginateWithTablePrefix(metricHistories, &listOptions.Pagination, r.db, "Artifact"))
	} else {
		query = query.Scopes(scopes.Paginate(metricHistories, &listOptions.Pagination, r.db))
	}

	if err := query.Find(&metricHistoriesArt).Error; err != nil {
		return nil, fmt.Errorf("error listing metric histories: %w", err)
	}

	hasMore := false
	pageSize := listOptions.GetPageSize()
	if pageSize > 0 {
		hasMore = len(metricHistoriesArt) > int(pageSize)
		if hasMore {
			metricHistoriesArt = metricHistoriesArt[:len(metricHistoriesArt)-1]
		}
	}

	for _, metricHistoryArt := range metricHistoriesArt {
		properties := []schema.ArtifactProperty{}
		if err := r.db.Where("artifact_id = ?", metricHistoryArt.ID).Find(&properties).Error; err != nil {
			return nil, fmt.Errorf("error getting properties by metric history id: %w", err)
		}

		metricHistory := mapDataLayerToMetricHistory(metricHistoryArt, properties)
		metricHistories = append(metricHistories, metricHistory)
	}

	if hasMore && len(metricHistoriesArt) > 0 {
		lastModel := metricHistoriesArt[len(metricHistoriesArt)-1]
		orderBy := listOptions.GetOrderBy()
		value := ""
		if orderBy != "" {
			switch orderBy {
			case "ID":
				value = fmt.Sprintf("%d", lastModel.ID)
			case "CREATE_TIME":
				value = fmt.Sprintf("%d", lastModel.CreateTimeSinceEpoch)
			case "LAST_UPDATE_TIME":
				value = fmt.Sprintf("%d", lastModel.LastUpdateTimeSinceEpoch)
			default:
				value = fmt.Sprintf("%d", lastModel.ID)
			}
		}
		nextToken := scopes.CreateNextPageToken(lastModel.ID, value)
		listOptions.NextPageToken = &nextToken
	} else {
		listOptions.NextPageToken = nil
	}

	list.Items = metricHistories
	list.NextPageToken = listOptions.GetNextPageToken()
	list.PageSize = listOptions.GetPageSize()
	list.Size = int32(len(metricHistories))

	return &list, nil
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
		artifact.Name = metricHistory.GetAttributes().Name
		artifact.URI = metricHistory.GetAttributes().URI
		artifact.ExternalID = metricHistory.GetAttributes().ExternalID
		if metricHistory.GetAttributes().State != nil {
			stateValue := models.Artifact_State_value[*metricHistory.GetAttributes().State]
			artifact.State = &stateValue
		}
		artifact.CreateTimeSinceEpoch = apiutils.ZeroIfNil(metricHistory.GetAttributes().CreateTimeSinceEpoch)
		artifact.LastUpdateTimeSinceEpoch = apiutils.ZeroIfNil(metricHistory.GetAttributes().LastUpdateTimeSinceEpoch)
	}

	return artifact
}

func mapMetricHistoryToArtifactProperties(metricHistory models.MetricHistory, artifactID int32) []schema.ArtifactProperty {
	if metricHistory == nil {
		return []schema.ArtifactProperty{}
	}

	properties := []schema.ArtifactProperty{}

	if metricHistory.GetProperties() != nil {
		for _, prop := range *metricHistory.GetProperties() {
			properties = append(properties, mapPropertiesToArtifactProperty(prop, artifactID, false))
		}
	}

	if metricHistory.GetCustomProperties() != nil {
		for _, prop := range *metricHistory.GetCustomProperties() {
			properties = append(properties, mapPropertiesToArtifactProperty(prop, artifactID, true))
		}
	}

	return properties
}

func mapDataLayerToMetricHistory(metricHistory schema.Artifact, artProperties []schema.ArtifactProperty) models.MetricHistory {
	var state *string
	metricHistoryType := models.MetricHistoryType

	if metricHistory.State != nil {
		metricHistoryState := models.Artifact_State_name[*metricHistory.State]
		state = &metricHistoryState
	}

	metricHistoryArt := models.MetricHistoryImpl{
		ID:     &metricHistory.ID,
		TypeID: &metricHistory.TypeID,
		Attributes: &models.MetricHistoryAttributes{
			Name:                     metricHistory.Name,
			URI:                      metricHistory.URI,
			State:                    state,
			ArtifactType:             &metricHistoryType,
			ExternalID:               metricHistory.ExternalID,
			CreateTimeSinceEpoch:     &metricHistory.CreateTimeSinceEpoch,
			LastUpdateTimeSinceEpoch: &metricHistory.LastUpdateTimeSinceEpoch,
		},
	}

	customProperties := []models.Properties{}
	properties := []models.Properties{}

	for _, prop := range artProperties {
		if prop.IsCustomProperty {
			customProperties = append(customProperties, mapDataLayerToArtifactProperties(prop))
		} else {
			properties = append(properties, mapDataLayerToArtifactProperties(prop))
		}
	}

	metricHistoryArt.CustomProperties = &customProperties
	metricHistoryArt.Properties = &properties

	return &metricHistoryArt
}

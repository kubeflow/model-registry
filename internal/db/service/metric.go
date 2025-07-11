package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/scopes"
	"gorm.io/gorm"
)

var ErrMetricNotFound = errors.New("metric by id not found")

type MetricRepositoryImpl struct {
	db     *gorm.DB
	typeID int64
}

func NewMetricRepository(db *gorm.DB, typeID int64) models.MetricRepository {
	return &MetricRepositoryImpl{db: db, typeID: typeID}
}

func (r *MetricRepositoryImpl) GetByID(id int32) (models.Metric, error) {
	metric := &schema.Artifact{}
	properties := []schema.ArtifactProperty{}

	if err := r.db.Where("id = ? AND type_id = ?", id, r.typeID).First(metric).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %v", ErrMetricNotFound, err)
		}

		return nil, fmt.Errorf("error getting metric by id: %w", err)
	}

	if err := r.db.Where("artifact_id = ?", metric.ID).Find(&properties).Error; err != nil {
		return nil, fmt.Errorf("error getting properties by metric id: %w", err)
	}

	return mapDataLayerToMetric(*metric, properties), nil
}

func (r *MetricRepositoryImpl) Save(metric models.Metric, parentResourceID *int32) (models.Metric, error) {
	now := time.Now().UnixMilli()

	metricArt := mapMetricToArtifact(metric)
	propertiesArt := []schema.ArtifactProperty{}

	metricArt.LastUpdateTimeSinceEpoch = now

	if metric.GetID() == nil {
		metricArt.CreateTimeSinceEpoch = now
	}

	hasCustomProperties := metric.GetCustomProperties() != nil

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&metricArt).Error; err != nil {
			return fmt.Errorf("error saving metric: %w", err)
		}

		propertiesArt = mapMetricToArtifactProperties(metric, metricArt.ID)
		existingCustomPropertiesArt := []schema.ArtifactProperty{}

		if err := tx.Where("artifact_id = ? AND is_custom_property = ?", metricArt.ID, true).Find(&existingCustomPropertiesArt).Error; err != nil {
			return fmt.Errorf("error getting existing custom properties by metric id: %w", err)
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
						return fmt.Errorf("error deleting metric property: %w", err)
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
					return fmt.Errorf("error updating metric property: %w", err)
				}
			case gorm.ErrRecordNotFound:
				if err := tx.Create(&prop).Error; err != nil {
					return fmt.Errorf("error creating metric property: %w", err)
				}
			default:
				return fmt.Errorf("error checking existing property: %w", result.Error)
			}
		}

		if parentResourceID != nil {
			// Check if attribution already exists to avoid duplicate key errors
			var existingAttribution schema.Attribution
			result := tx.Where("context_id = ? AND artifact_id = ?", *parentResourceID, metricArt.ID).First(&existingAttribution)

			if result.Error == gorm.ErrRecordNotFound {
				// Attribution doesn't exist, create it
				attribution := schema.Attribution{
					ContextID:  *parentResourceID,
					ArtifactID: metricArt.ID,
				}

				if err := tx.Create(&attribution).Error; err != nil {
					return fmt.Errorf("error creating attribution: %w", err)
				}
			} else if result.Error != nil {
				return fmt.Errorf("error checking existing attribution: %w", result.Error)
			}
		}

		// Get all final properties for the return object
		if err := tx.Where("artifact_id = ?", metricArt.ID).Find(&propertiesArt).Error; err != nil {
			return fmt.Errorf("error getting final properties by metric id: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return mapDataLayerToMetric(metricArt, propertiesArt), nil
}

func (r *MetricRepositoryImpl) List(listOptions models.MetricListOptions) (*models.ListWrapper[models.Metric], error) {
	list := models.ListWrapper[models.Metric]{
		PageSize: listOptions.GetPageSize(),
	}

	metrics := []models.Metric{}
	metricsArt := []schema.Artifact{}

	query := r.db.Model(&schema.Artifact{}).Where("type_id = ?", r.typeID)

	if listOptions.Name != nil {
		query = query.Where("name = ?", listOptions.Name)
	} else if listOptions.ExternalID != nil {
		query = query.Where("external_id = ?", listOptions.ExternalID)
	}

	if listOptions.ParentResourceID != nil {
		query = query.Joins("JOIN Attribution ON Attribution.artifact_id = Artifact.id").
			Where("Attribution.context_id = ?", listOptions.ParentResourceID)
		query = query.Scopes(scopes.PaginateWithTablePrefix(metrics, &listOptions.Pagination, r.db, "Artifact"))
	} else {
		query = query.Scopes(scopes.Paginate(metrics, &listOptions.Pagination, r.db))
	}

	if err := query.Find(&metricsArt).Error; err != nil {
		return nil, fmt.Errorf("error listing metrics: %w", err)
	}

	hasMore := false
	pageSize := listOptions.GetPageSize()
	if pageSize > 0 {
		hasMore = len(metricsArt) > int(pageSize)
		if hasMore {
			metricsArt = metricsArt[:len(metricsArt)-1]
		}
	}

	for _, metricArt := range metricsArt {
		properties := []schema.ArtifactProperty{}
		if err := r.db.Where("artifact_id = ?", metricArt.ID).Find(&properties).Error; err != nil {
			return nil, fmt.Errorf("error getting properties by metric id: %w", err)
		}

		metric := mapDataLayerToMetric(metricArt, properties)
		metrics = append(metrics, metric)
	}

	if hasMore && len(metricsArt) > 0 {
		lastModel := metricsArt[len(metricsArt)-1]
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

	list.Items = metrics
	list.NextPageToken = listOptions.GetNextPageToken()
	list.PageSize = listOptions.GetPageSize()
	list.Size = int32(len(metrics))

	return &list, nil
}

func mapMetricToArtifact(metric models.Metric) schema.Artifact {
	if metric == nil {
		return schema.Artifact{}
	}

	artifact := schema.Artifact{
		ID:     apiutils.ZeroIfNil(metric.GetID()),
		TypeID: apiutils.ZeroIfNil(metric.GetTypeID()),
	}

	if metric.GetAttributes() != nil {
		artifact.Name = metric.GetAttributes().Name
		artifact.URI = metric.GetAttributes().URI
		artifact.ExternalID = metric.GetAttributes().ExternalID
		if metric.GetAttributes().State != nil {
			stateValue := models.Artifact_State_value[*metric.GetAttributes().State]
			artifact.State = &stateValue
		}
		artifact.CreateTimeSinceEpoch = apiutils.ZeroIfNil(metric.GetAttributes().CreateTimeSinceEpoch)
		artifact.LastUpdateTimeSinceEpoch = apiutils.ZeroIfNil(metric.GetAttributes().LastUpdateTimeSinceEpoch)
	}

	return artifact
}

func mapMetricToArtifactProperties(metric models.Metric, artifactID int32) []schema.ArtifactProperty {
	if metric == nil {
		return []schema.ArtifactProperty{}
	}

	properties := []schema.ArtifactProperty{}

	if metric.GetProperties() != nil {
		for _, prop := range *metric.GetProperties() {
			properties = append(properties, mapPropertiesToArtifactProperty(prop, artifactID, false))
		}
	}

	if metric.GetCustomProperties() != nil {
		for _, prop := range *metric.GetCustomProperties() {
			properties = append(properties, mapPropertiesToArtifactProperty(prop, artifactID, true))
		}
	}

	return properties
}

func mapDataLayerToMetric(metric schema.Artifact, artProperties []schema.ArtifactProperty) models.Metric {
	var state *string
	metricType := models.MetricType

	if metric.State != nil {
		metricState := models.Artifact_State_name[*metric.State]
		state = &metricState
	}

	metricArt := models.MetricImpl{
		ID:     &metric.ID,
		TypeID: &metric.TypeID,
		Attributes: &models.MetricAttributes{
			Name:                     metric.Name,
			URI:                      metric.URI,
			State:                    state,
			ArtifactType:             &metricType,
			ExternalID:               metric.ExternalID,
			CreateTimeSinceEpoch:     &metric.CreateTimeSinceEpoch,
			LastUpdateTimeSinceEpoch: &metric.LastUpdateTimeSinceEpoch,
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

	metricArt.CustomProperties = &customProperties
	metricArt.Properties = &properties

	return &metricArt
}

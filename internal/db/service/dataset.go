package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/scopes"
	"gorm.io/gorm"
)

var ErrDataSetNotFound = errors.New("dataset by id not found")

type DataSetRepositoryImpl struct {
	db     *gorm.DB
	typeID int64
}

func NewDataSetRepository(db *gorm.DB, typeID int64) models.DataSetRepository {
	return &DataSetRepositoryImpl{db: db, typeID: typeID}
}

func (r *DataSetRepositoryImpl) GetByID(id int32) (models.DataSet, error) {
	dataSet := &schema.Artifact{}
	properties := []schema.ArtifactProperty{}

	if err := r.db.Where("id = ? AND type_id = ?", id, r.typeID).First(dataSet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %v", ErrDataSetNotFound, err)
		}
		return nil, fmt.Errorf("error getting dataset by id: %w", err)
	}

	if err := r.db.Where("artifact_id = ?", dataSet.ID).Find(&properties).Error; err != nil {
		return nil, fmt.Errorf("error getting properties by dataset id: %w", err)
	}

	return mapDataLayerToDataSet(*dataSet, properties), nil
}

func (r *DataSetRepositoryImpl) Save(dataSet models.DataSet, parentResourceID *int32) (models.DataSet, error) {
	now := time.Now().UnixMilli()

	dataSetArt := mapDataSetToArtifact(dataSet)
	properties := mapDataSetToArtifactProperties(dataSet, dataSetArt.ID)

	dataSetArt.LastUpdateTimeSinceEpoch = now

	if dataSet.GetID() == nil {
		glog.Info("Creating new DataSet")
		dataSetArt.CreateTimeSinceEpoch = now
	} else {
		glog.Infof("Updating DataSet %d", *dataSet.GetID())
	}

	hasCustomProperties := dataSet.GetCustomProperties() != nil

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&dataSetArt).Error; err != nil {
			return fmt.Errorf("error saving dataset: %w", err)
		}

		properties = mapDataSetToArtifactProperties(dataSet, dataSetArt.ID)
		existingCustomProperties := []schema.ArtifactProperty{}

		if err := tx.Where("artifact_id = ? AND is_custom_property = ?", dataSetArt.ID, true).Find(&existingCustomProperties).Error; err != nil {
			return fmt.Errorf("error getting existing custom properties by dataset id: %w", err)
		}

		if hasCustomProperties {
			for _, existingProp := range existingCustomProperties {
				found := false
				for _, prop := range properties {
					if prop.Name == existingProp.Name && prop.ArtifactID == existingProp.ArtifactID && prop.IsCustomProperty == existingProp.IsCustomProperty {
						found = true
						break
					}
				}

				if !found {
					if err := tx.Delete(&existingProp).Error; err != nil {
						return fmt.Errorf("error deleting dataset property: %w", err)
					}
				}
			}
		}

		for _, prop := range properties {
			var existingProp schema.ArtifactProperty
			result := tx.Where("artifact_id = ? AND name = ? AND is_custom_property = ?",
				prop.ArtifactID, prop.Name, prop.IsCustomProperty).First(&existingProp)

			switch result.Error {
			case nil:
				prop.ArtifactID = existingProp.ArtifactID
				prop.Name = existingProp.Name
				prop.IsCustomProperty = existingProp.IsCustomProperty
				if err := tx.Model(&existingProp).Updates(prop).Error; err != nil {
					return fmt.Errorf("error updating dataset property: %w", err)
				}
			case gorm.ErrRecordNotFound:
				if err := tx.Create(&prop).Error; err != nil {
					return fmt.Errorf("error creating dataset property: %w", err)
				}
			default:
				return fmt.Errorf("error checking existing property: %w", result.Error)
			}
		}

		if parentResourceID != nil {
			// Check if attribution already exists to avoid duplicate key errors
			var existingAttribution schema.Attribution
			result := tx.Where("context_id = ? AND artifact_id = ?", *parentResourceID, dataSetArt.ID).First(&existingAttribution)

			if result.Error == gorm.ErrRecordNotFound {
				// Attribution doesn't exist, create it
				attribution := schema.Attribution{
					ContextID:  *parentResourceID,
					ArtifactID: dataSetArt.ID,
				}

				if err := tx.Create(&attribution).Error; err != nil {
					return fmt.Errorf("error creating attribution: %w", err)
				}
			} else if result.Error != nil {
				return fmt.Errorf("error checking existing attribution: %w", result.Error)
			}
			// If attribution already exists, do nothing
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return mapDataLayerToDataSet(dataSetArt, properties), nil
}

func (r *DataSetRepositoryImpl) List(listOptions models.DataSetListOptions) (*models.ListWrapper[models.DataSet], error) {
	list := models.ListWrapper[models.DataSet]{
		PageSize: listOptions.GetPageSize(),
	}

	dataSets := []models.DataSet{}
	dataSetsArt := []schema.Artifact{}

	query := r.db.Model(&schema.Artifact{}).Where("type_id = ?", r.typeID)

	if listOptions.Name != nil {
		query = query.Where("name = ?", listOptions.Name)
	} else if listOptions.ExternalID != nil {
		query = query.Where("external_id = ?", listOptions.ExternalID)
	}

	if listOptions.ParentResourceID != nil {
		query = query.Joins("JOIN Attribution ON Attribution.artifact_id = Artifact.id").
			Where("Attribution.context_id = ?", listOptions.ParentResourceID)
		// Use table-prefixed pagination to avoid column ambiguity
		query = query.Scopes(scopes.PaginateWithTablePrefix(dataSets, &listOptions.Pagination, r.db, "Artifact"))
	} else {
		query = query.Scopes(scopes.Paginate(dataSets, &listOptions.Pagination, r.db))
	}

	if err := query.Find(&dataSetsArt).Error; err != nil {
		return nil, fmt.Errorf("error listing datasets: %w", err)
	}

	hasMore := false
	pageSize := listOptions.GetPageSize()
	if pageSize > 0 {
		hasMore = len(dataSetsArt) > int(pageSize)
		if hasMore {
			dataSetsArt = dataSetsArt[:len(dataSetsArt)-1]
		}
	}

	for _, dataSetArt := range dataSetsArt {
		properties := []schema.ArtifactProperty{}
		if err := r.db.Where("artifact_id = ?", dataSetArt.ID).Find(&properties).Error; err != nil {
			return nil, fmt.Errorf("error getting properties by dataset id: %w", err)
		}

		dataSet := mapDataLayerToDataSet(dataSetArt, properties)
		dataSets = append(dataSets, dataSet)
	}

	if hasMore && len(dataSetsArt) > 0 {
		lastDataSet := dataSetsArt[len(dataSetsArt)-1]
		orderBy := listOptions.GetOrderBy()
		value := ""
		if orderBy != "" {
			switch orderBy {
			case "ID":
				value = fmt.Sprintf("%d", lastDataSet.ID)
			case "CREATE_TIME":
				value = fmt.Sprintf("%d", lastDataSet.CreateTimeSinceEpoch)
			case "LAST_UPDATE_TIME":
				value = fmt.Sprintf("%d", lastDataSet.LastUpdateTimeSinceEpoch)
			default:
				value = fmt.Sprintf("%d", lastDataSet.ID)
			}
		}
		nextToken := scopes.CreateNextPageToken(lastDataSet.ID, value)
		listOptions.NextPageToken = &nextToken
	} else {
		listOptions.NextPageToken = nil
	}

	list.Items = dataSets
	list.NextPageToken = listOptions.GetNextPageToken()
	list.PageSize = listOptions.GetPageSize()
	list.Size = int32(len(dataSets))

	return &list, nil
}

func mapDataSetToArtifact(dataSet models.DataSet) schema.Artifact {
	if dataSet == nil {
		return schema.Artifact{}
	}

	artifact := schema.Artifact{
		ID:     apiutils.ZeroIfNil(dataSet.GetID()),
		TypeID: apiutils.ZeroIfNil(dataSet.GetTypeID()),
	}

	if dataSet.GetAttributes() != nil {
		artifact.Name = dataSet.GetAttributes().Name
		artifact.URI = dataSet.GetAttributes().URI
		artifact.ExternalID = dataSet.GetAttributes().ExternalID
		if dataSet.GetAttributes().State != nil {
			stateValue := models.Artifact_State_value[*dataSet.GetAttributes().State]
			artifact.State = &stateValue
		}
		artifact.CreateTimeSinceEpoch = apiutils.ZeroIfNil(dataSet.GetAttributes().CreateTimeSinceEpoch)
		artifact.LastUpdateTimeSinceEpoch = apiutils.ZeroIfNil(dataSet.GetAttributes().LastUpdateTimeSinceEpoch)
	}

	return artifact
}

func mapDataSetToArtifactProperties(dataSet models.DataSet, artifactID int32) []schema.ArtifactProperty {
	if dataSet == nil {
		return []schema.ArtifactProperty{}
	}

	properties := []schema.ArtifactProperty{}

	if dataSet.GetProperties() != nil {
		for _, prop := range *dataSet.GetProperties() {
			properties = append(properties, mapPropertiesToArtifactProperty(prop, artifactID, false))
		}
	}

	if dataSet.GetCustomProperties() != nil {
		for _, prop := range *dataSet.GetCustomProperties() {
			properties = append(properties, mapPropertiesToArtifactProperty(prop, artifactID, true))
		}
	}

	return properties
}

func mapDataLayerToDataSet(dataSet schema.Artifact, artProperties []schema.ArtifactProperty) models.DataSet {
	var state *string
	dataSetType := models.DataSetType

	if dataSet.State != nil {
		dsState := models.Artifact_State_name[*dataSet.State]
		state = &dsState
	}

	dataSetArt := models.DataSetImpl{
		ID:     &dataSet.ID,
		TypeID: &dataSet.TypeID,
		Attributes: &models.DataSetAttributes{
			Name:                     dataSet.Name,
			URI:                      dataSet.URI,
			State:                    state,
			ArtifactType:             &dataSetType,
			ExternalID:               dataSet.ExternalID,
			CreateTimeSinceEpoch:     &dataSet.CreateTimeSinceEpoch,
			LastUpdateTimeSinceEpoch: &dataSet.LastUpdateTimeSinceEpoch,
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

	dataSetArt.CustomProperties = &customProperties
	dataSetArt.Properties = &properties

	return &dataSetArt
}

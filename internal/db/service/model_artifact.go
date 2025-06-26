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

var ErrModelArtifactNotFound = errors.New("model artifact by id not found")

type ModelArtifactRepositoryImpl struct {
	db     *gorm.DB
	typeID int64
}

func NewModelArtifactRepository(db *gorm.DB, typeID int64) models.ModelArtifactRepository {
	return &ModelArtifactRepositoryImpl{db: db, typeID: typeID}
}

func (r *ModelArtifactRepositoryImpl) GetByID(id int32) (models.ModelArtifact, error) {
	modelArtifact := &schema.Artifact{}
	properties := []schema.ArtifactProperty{}

	if err := r.db.Where("id = ? AND type_id = ?", id, r.typeID).First(modelArtifact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %v", ErrModelArtifactNotFound, err)
		}

		return nil, fmt.Errorf("error getting model artifact by id: %w", err)
	}

	if err := r.db.Where("artifact_id = ?", modelArtifact.ID).Find(&properties).Error; err != nil {
		return nil, fmt.Errorf("error getting properties by model artifact id: %w", err)
	}

	return mapDataLayerToModelArtifact(*modelArtifact, properties), nil
}

func (r *ModelArtifactRepositoryImpl) Save(modelArtifact models.ModelArtifact, modelVersionID *int32) (models.ModelArtifact, error) {
	now := time.Now().UnixMilli()

	modelArtifactArt := mapModelArtifactToArtifact(modelArtifact)
	properties := mapModelArtifactToArtifactProperties(modelArtifact, modelArtifactArt.ID)

	modelArtifactArt.LastUpdateTimeSinceEpoch = now

	if modelArtifact.GetID() == nil {
		glog.Info("Creating new ModelArtifact")

		modelArtifactArt.CreateTimeSinceEpoch = now
	} else {
		glog.Infof("Updating ModelArtifact %d", *modelArtifact.GetID())
	}

	hasCustomProperties := modelArtifact.GetCustomProperties() != nil

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&modelArtifactArt).Error; err != nil {
			return fmt.Errorf("error saving model artifact: %w", err)
		}

		properties = mapModelArtifactToArtifactProperties(modelArtifact, modelArtifactArt.ID)
		existingCustomProperties := []schema.ArtifactProperty{}

		if err := tx.Where("artifact_id = ? AND is_custom_property = ?", modelArtifactArt.ID, true).Find(&existingCustomProperties).Error; err != nil {
			return fmt.Errorf("error getting existing custom properties by model artifact id: %w", err)
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
						return fmt.Errorf("error deleting model artifact property: %w", err)
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
					return fmt.Errorf("error updating model artifact property: %w", err)
				}
			case gorm.ErrRecordNotFound:
				if err := tx.Create(&prop).Error; err != nil {
					return fmt.Errorf("error creating model artifact property: %w", err)
				}
			default:
				return fmt.Errorf("error checking existing property: %w", result.Error)
			}
		}

		if modelVersionID != nil {
			// Check if attribution already exists to avoid duplicate key errors
			var existingAttribution schema.Attribution
			result := tx.Where("context_id = ? AND artifact_id = ?", *modelVersionID, modelArtifactArt.ID).First(&existingAttribution)

			if result.Error == gorm.ErrRecordNotFound {
				// Attribution doesn't exist, create it
				attribution := schema.Attribution{
					ContextID:  *modelVersionID,
					ArtifactID: modelArtifactArt.ID,
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

	return mapDataLayerToModelArtifact(modelArtifactArt, properties), nil
}

func (r *ModelArtifactRepositoryImpl) List(listOptions models.ModelArtifactListOptions) (*models.ListWrapper[models.ModelArtifact], error) {
	list := models.ListWrapper[models.ModelArtifact]{
		PageSize: listOptions.GetPageSize(),
	}

	modelArtifacts := []models.ModelArtifact{}
	modelArtifactsArt := []schema.Artifact{}

	query := r.db.Model(&schema.Artifact{}).Where("type_id = ?", r.typeID)

	if listOptions.Name != nil {
		query = query.Where("name = ?", listOptions.Name)
	} else if listOptions.ExternalID != nil {
		query = query.Where("external_id = ?", listOptions.ExternalID)
	}

	if listOptions.ModelVersionID != nil {
		query = query.Joins("JOIN Attribution ON Attribution.artifact_id = Artifact.id").
			Where("Attribution.context_id = ?", listOptions.ModelVersionID)
		// Use table-prefixed pagination to avoid column ambiguity
		query = query.Scopes(scopes.PaginateWithTablePrefix(modelArtifacts, &listOptions.Pagination, r.db, "Artifact"))
	} else {
		query = query.Scopes(scopes.Paginate(modelArtifacts, &listOptions.Pagination, r.db))
	}

	if err := query.Find(&modelArtifactsArt).Error; err != nil {
		return nil, fmt.Errorf("error listing model artifacts: %w", err)
	}

	hasMore := false
	pageSize := listOptions.GetPageSize()
	if pageSize > 0 {
		hasMore = len(modelArtifactsArt) > int(pageSize)
		if hasMore {
			modelArtifactsArt = modelArtifactsArt[:len(modelArtifactsArt)-1]
		}
	}

	for _, modelArtifactArt := range modelArtifactsArt {
		properties := []schema.ArtifactProperty{}
		if err := r.db.Where("artifact_id = ?", modelArtifactArt.ID).Find(&properties).Error; err != nil {
			return nil, fmt.Errorf("error getting properties by model artifact id: %w", err)
		}

		modelArtifact := mapDataLayerToModelArtifact(modelArtifactArt, properties)
		modelArtifacts = append(modelArtifacts, modelArtifact)
	}

	if hasMore && len(modelArtifactsArt) > 0 {
		lastModel := modelArtifactsArt[len(modelArtifactsArt)-1]
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

	list.Items = modelArtifacts
	list.NextPageToken = listOptions.GetNextPageToken()
	list.PageSize = listOptions.GetPageSize()
	list.Size = int32(len(modelArtifacts))

	return &list, nil
}

func mapModelArtifactToArtifact(modelArtifact models.ModelArtifact) schema.Artifact {
	if modelArtifact == nil {
		return schema.Artifact{}
	}

	artifact := schema.Artifact{
		ID:     apiutils.ZeroIfNil(modelArtifact.GetID()),
		TypeID: apiutils.ZeroIfNil(modelArtifact.GetTypeID()),
	}

	if modelArtifact.GetAttributes() != nil {
		artifact.Name = modelArtifact.GetAttributes().Name
		artifact.URI = modelArtifact.GetAttributes().URI
		artifact.ExternalID = modelArtifact.GetAttributes().ExternalID
		if modelArtifact.GetAttributes().State != nil {
			stateValue := models.Artifact_State_value[*modelArtifact.GetAttributes().State]
			artifact.State = &stateValue
		}
		artifact.CreateTimeSinceEpoch = apiutils.ZeroIfNil(modelArtifact.GetAttributes().CreateTimeSinceEpoch)
		artifact.LastUpdateTimeSinceEpoch = apiutils.ZeroIfNil(modelArtifact.GetAttributes().LastUpdateTimeSinceEpoch)
	}

	return artifact
}

func mapModelArtifactToArtifactProperties(modelArtifact models.ModelArtifact, artifactID int32) []schema.ArtifactProperty {
	if modelArtifact == nil {
		return []schema.ArtifactProperty{}
	}

	properties := []schema.ArtifactProperty{}

	if modelArtifact.GetProperties() != nil {
		for _, prop := range *modelArtifact.GetProperties() {
			properties = append(properties, mapPropertiesToArtifactProperty(prop, artifactID, false))
		}
	}

	if modelArtifact.GetCustomProperties() != nil {
		for _, prop := range *modelArtifact.GetCustomProperties() {
			properties = append(properties, mapPropertiesToArtifactProperty(prop, artifactID, true))
		}
	}

	return properties
}

func mapDataLayerToModelArtifact(modelArtifact schema.Artifact, artProperties []schema.ArtifactProperty) models.ModelArtifact {
	var state *string
	modelArtifactType := models.ModelArtifactType

	if modelArtifact.State != nil {
		docState := models.Artifact_State_name[*modelArtifact.State]
		state = &docState
	}

	modelArtifactArt := models.ModelArtifactImpl{
		ID:     &modelArtifact.ID,
		TypeID: &modelArtifact.TypeID,
		Attributes: &models.ModelArtifactAttributes{
			Name:                     modelArtifact.Name,
			URI:                      modelArtifact.URI,
			State:                    state,
			ArtifactType:             &modelArtifactType,
			ExternalID:               modelArtifact.ExternalID,
			CreateTimeSinceEpoch:     &modelArtifact.CreateTimeSinceEpoch,
			LastUpdateTimeSinceEpoch: &modelArtifact.LastUpdateTimeSinceEpoch,
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

	modelArtifactArt.CustomProperties = &customProperties
	modelArtifactArt.Properties = &properties

	return &modelArtifactArt
}

func mapPropertiesToArtifactProperty(prop models.Properties, artifactID int32, isCustomProperty bool) schema.ArtifactProperty {
	artProp := schema.ArtifactProperty{
		ArtifactID:       artifactID,
		Name:             prop.Name,
		IsCustomProperty: isCustomProperty,
		IntValue:         prop.IntValue,
		DoubleValue:      prop.DoubleValue,
		StringValue:      prop.StringValue,
		BoolValue:        prop.BoolValue,
		ByteValue:        prop.ByteValue,
		ProtoValue:       prop.ProtoValue,
	}

	return artProp
}

func mapDataLayerToArtifactProperties(artProperty schema.ArtifactProperty) models.Properties {
	prop := models.Properties{
		Name:             artProperty.Name,
		IsCustomProperty: artProperty.IsCustomProperty,
		IntValue:         artProperty.IntValue,
		DoubleValue:      artProperty.DoubleValue,
		StringValue:      artProperty.StringValue,
		BoolValue:        artProperty.BoolValue,
		ByteValue:        artProperty.ByteValue,
		ProtoValue:       artProperty.ProtoValue,
	}

	return prop
}

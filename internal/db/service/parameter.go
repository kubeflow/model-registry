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

var ErrParameterNotFound = errors.New("parameter by id not found")

type ParameterRepositoryImpl struct {
	db     *gorm.DB
	typeID int64
}

func NewParameterRepository(db *gorm.DB, typeID int64) models.ParameterRepository {
	return &ParameterRepositoryImpl{db: db, typeID: typeID}
}

func (r *ParameterRepositoryImpl) GetByID(id int32) (models.Parameter, error) {
	parameter := &schema.Artifact{}
	properties := []schema.ArtifactProperty{}

	if err := r.db.Where("id = ? AND type_id = ?", id, r.typeID).First(parameter).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %v", ErrParameterNotFound, err)
		}

		return nil, fmt.Errorf("error getting parameter by id: %w", err)
	}

	if err := r.db.Where("artifact_id = ?", parameter.ID).Find(&properties).Error; err != nil {
		return nil, fmt.Errorf("error getting properties by parameter id: %w", err)
	}

	return mapDataLayerToParameter(*parameter, properties), nil
}

func (r *ParameterRepositoryImpl) Save(parameter models.Parameter, parentResourceID *int32) (models.Parameter, error) {
	now := time.Now().UnixMilli()

	parameterArt := mapParameterToArtifact(parameter)
	propertiesArt := []schema.ArtifactProperty{}

	parameterArt.LastUpdateTimeSinceEpoch = now

	if parameter.GetID() == nil {
		glog.Info("Creating new Parameter")
		parameterArt.CreateTimeSinceEpoch = now
	} else {
		glog.Infof("Updating Parameter %d", *parameter.GetID())
	}

	hasCustomProperties := parameter.GetCustomProperties() != nil

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&parameterArt).Error; err != nil {
			return fmt.Errorf("error saving parameter: %w", err)
		}

		propertiesArt = mapParameterToArtifactProperties(parameter, parameterArt.ID)
		existingCustomPropertiesArt := []schema.ArtifactProperty{}

		if err := tx.Where("artifact_id = ? AND is_custom_property = ?", parameterArt.ID, true).Find(&existingCustomPropertiesArt).Error; err != nil {
			return fmt.Errorf("error getting existing custom properties by parameter id: %w", err)
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
						return fmt.Errorf("error deleting parameter property: %w", err)
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
					return fmt.Errorf("error updating parameter property: %w", err)
				}
			case gorm.ErrRecordNotFound:
				if err := tx.Create(&prop).Error; err != nil {
					return fmt.Errorf("error creating parameter property: %w", err)
				}
			default:
				return fmt.Errorf("error checking existing property: %w", result.Error)
			}
		}

		if parentResourceID != nil {
			// Check if attribution already exists to avoid duplicate key errors
			var existingAttribution schema.Attribution
			result := tx.Where("context_id = ? AND artifact_id = ?", *parentResourceID, parameterArt.ID).First(&existingAttribution)

			if result.Error == gorm.ErrRecordNotFound {
				// Attribution doesn't exist, create it
				attribution := schema.Attribution{
					ContextID:  *parentResourceID,
					ArtifactID: parameterArt.ID,
				}

				if err := tx.Create(&attribution).Error; err != nil {
					return fmt.Errorf("error creating attribution: %w", err)
				}
			} else if result.Error != nil {
				return fmt.Errorf("error checking existing attribution: %w", result.Error)
			}

		}

		// Get all final properties for the return object
		if err := tx.Where("artifact_id = ?", parameterArt.ID).Find(&propertiesArt).Error; err != nil {
			return fmt.Errorf("error getting final properties by parameter id: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return mapDataLayerToParameter(parameterArt, propertiesArt), nil
}

func (r *ParameterRepositoryImpl) List(listOptions models.ParameterListOptions) (*models.ListWrapper[models.Parameter], error) {
	list := models.ListWrapper[models.Parameter]{
		PageSize: listOptions.GetPageSize(),
	}

	parameters := []models.Parameter{}
	parametersArt := []schema.Artifact{}

	query := r.db.Model(&schema.Artifact{}).Where("type_id = ?", r.typeID)

	if listOptions.Name != nil {
		query = query.Where("name = ?", listOptions.Name)
	} else if listOptions.ExternalID != nil {
		query = query.Where("external_id = ?", listOptions.ExternalID)
	}

	if listOptions.ParentResourceID != nil {
		query = query.Joins("JOIN Attribution ON Attribution.artifact_id = Artifact.id").
			Where("Attribution.context_id = ?", listOptions.ParentResourceID)
		query = query.Scopes(scopes.PaginateWithTablePrefix(parameters, &listOptions.Pagination, r.db, "Artifact"))
	} else {
		query = query.Scopes(scopes.Paginate(parameters, &listOptions.Pagination, r.db))
	}

	if err := query.Find(&parametersArt).Error; err != nil {
		return nil, fmt.Errorf("error listing parameters: %w", err)
	}

	hasMore := false
	pageSize := listOptions.GetPageSize()
	if pageSize > 0 {
		hasMore = len(parametersArt) > int(pageSize)
		if hasMore {
			parametersArt = parametersArt[:len(parametersArt)-1]
		}
	}

	for _, parameterArt := range parametersArt {
		properties := []schema.ArtifactProperty{}
		if err := r.db.Where("artifact_id = ?", parameterArt.ID).Find(&properties).Error; err != nil {
			return nil, fmt.Errorf("error getting properties by parameter id: %w", err)
		}

		parameter := mapDataLayerToParameter(parameterArt, properties)
		parameters = append(parameters, parameter)
	}

	if hasMore && len(parametersArt) > 0 {
		lastModel := parametersArt[len(parametersArt)-1]
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

	list.Items = parameters
	list.NextPageToken = listOptions.GetNextPageToken()
	list.PageSize = listOptions.GetPageSize()
	list.Size = int32(len(parameters))

	return &list, nil
}

func mapParameterToArtifact(parameter models.Parameter) schema.Artifact {
	if parameter == nil {
		return schema.Artifact{}
	}

	artifact := schema.Artifact{
		ID:     apiutils.ZeroIfNil(parameter.GetID()),
		TypeID: apiutils.ZeroIfNil(parameter.GetTypeID()),
	}

	if parameter.GetAttributes() != nil {
		artifact.Name = parameter.GetAttributes().Name
		artifact.URI = parameter.GetAttributes().URI
		artifact.ExternalID = parameter.GetAttributes().ExternalID
		if parameter.GetAttributes().State != nil {
			stateValue := models.Artifact_State_value[*parameter.GetAttributes().State]
			artifact.State = &stateValue
		}
		artifact.CreateTimeSinceEpoch = apiutils.ZeroIfNil(parameter.GetAttributes().CreateTimeSinceEpoch)
		artifact.LastUpdateTimeSinceEpoch = apiutils.ZeroIfNil(parameter.GetAttributes().LastUpdateTimeSinceEpoch)
	}

	return artifact
}

func mapParameterToArtifactProperties(parameter models.Parameter, artifactID int32) []schema.ArtifactProperty {
	if parameter == nil {
		return []schema.ArtifactProperty{}
	}

	properties := []schema.ArtifactProperty{}

	if parameter.GetProperties() != nil {
		for _, prop := range *parameter.GetProperties() {
			properties = append(properties, mapPropertiesToArtifactProperty(prop, artifactID, false))
		}
	}

	if parameter.GetCustomProperties() != nil {
		for _, prop := range *parameter.GetCustomProperties() {
			properties = append(properties, mapPropertiesToArtifactProperty(prop, artifactID, true))
		}
	}

	return properties
}

func mapDataLayerToParameter(parameter schema.Artifact, artProperties []schema.ArtifactProperty) models.Parameter {
	var state *string
	parameterType := models.ParameterType

	if parameter.State != nil {
		parameterState := models.Artifact_State_name[*parameter.State]
		state = &parameterState
	}

	parameterArt := models.ParameterImpl{
		ID:     &parameter.ID,
		TypeID: &parameter.TypeID,
		Attributes: &models.ParameterAttributes{
			Name:                     parameter.Name,
			URI:                      parameter.URI,
			State:                    state,
			ArtifactType:             &parameterType,
			ExternalID:               parameter.ExternalID,
			CreateTimeSinceEpoch:     &parameter.CreateTimeSinceEpoch,
			LastUpdateTimeSinceEpoch: &parameter.LastUpdateTimeSinceEpoch,
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

	parameterArt.CustomProperties = &customProperties
	parameterArt.Properties = &properties

	return &parameterArt
}

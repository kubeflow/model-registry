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

var ErrServeModelNotFound = errors.New("serve model by id not found")

type ServeModelRepositoryImpl struct {
	db     *gorm.DB
	typeID int64
}

func NewServeModelRepository(db *gorm.DB, typeID int64) models.ServeModelRepository {
	return &ServeModelRepositoryImpl{db: db, typeID: typeID}
}

func (r *ServeModelRepositoryImpl) GetByID(id int32) (models.ServeModel, error) {
	serveModel := &schema.Execution{}
	properties := []schema.ExecutionProperty{}

	if err := r.db.Where("id = ? AND type_id = ?", id, r.typeID).First(serveModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %v", ErrServeModelNotFound, err)
		}

		return nil, fmt.Errorf("error getting serve model by id: %w", err)
	}

	if err := r.db.Where("execution_id = ?", serveModel.ID).Find(&properties).Error; err != nil {
		return nil, fmt.Errorf("error getting properties by serve model id: %w", err)
	}

	return mapDataLayerToServeModel(*serveModel, properties), nil
}

func (r *ServeModelRepositoryImpl) Save(serveModel models.ServeModel, inferenceServiceID *int32) (models.ServeModel, error) {
	now := time.Now().UnixMilli()

	serveModelExec := mapServeModelToExecution(serveModel)
	properties := mapServeModelToExecutionProperties(serveModel, serveModelExec.ID)

	serveModelExec.LastUpdateTimeSinceEpoch = now

	if serveModel.GetID() == nil {
		glog.Info("Creating new ServeModel")
		serveModelExec.CreateTimeSinceEpoch = now
	} else {
		glog.Infof("Updating ServeModel %d", *serveModel.GetID())
	}

	hasCustomProperties := serveModel.GetCustomProperties() != nil

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&serveModelExec).Error; err != nil {
			return fmt.Errorf("error saving serve model: %w", err)
		}

		properties = mapServeModelToExecutionProperties(serveModel, serveModelExec.ID)
		existingCustomProperties := []schema.ExecutionProperty{}

		if err := tx.Where("execution_id = ? AND is_custom_property = ?", serveModelExec.ID, true).Find(&existingCustomProperties).Error; err != nil {
			return fmt.Errorf("error getting existing custom properties by serve model id: %w", err)
		}

		if hasCustomProperties {
			for _, existingProp := range existingCustomProperties {
				found := false
				for _, prop := range properties {
					if prop.Name == existingProp.Name && prop.ExecutionID == existingProp.ExecutionID && prop.IsCustomProperty == existingProp.IsCustomProperty {
						found = true
						break
					}
				}

				if !found {
					if err := tx.Delete(&existingProp).Error; err != nil {
						return fmt.Errorf("error deleting serve model property: %w", err)
					}
				}
			}
		}

		for _, prop := range properties {
			var existingProp schema.ExecutionProperty
			result := tx.Where("execution_id = ? AND name = ? AND is_custom_property = ?",
				prop.ExecutionID, prop.Name, prop.IsCustomProperty).First(&existingProp)

			switch result.Error {
			case nil:
				prop.ExecutionID = existingProp.ExecutionID
				prop.Name = existingProp.Name
				prop.IsCustomProperty = existingProp.IsCustomProperty
				if err := tx.Model(&existingProp).Updates(prop).Error; err != nil {
					return fmt.Errorf("error updating serve model property: %w", err)
				}
			case gorm.ErrRecordNotFound:
				if err := tx.Create(&prop).Error; err != nil {
					return fmt.Errorf("error creating serve model property: %w", err)
				}
			default:
				return fmt.Errorf("error checking existing property: %w", result.Error)
			}
		}

		if inferenceServiceID == nil {
			var associations []schema.Association
			tx.Where("execution_id = ?", serveModelExec.ID).Find(&associations)

			if len(associations) > 1 {
				return fmt.Errorf("multiple InferenceService found for ServeModel %d: %w", serveModelExec.ID, ErrModelArtifactNotFound)
			}

			if len(associations) == 0 {
				return fmt.Errorf("no InferenceService found for ServeModel %d: %w", serveModelExec.ID, ErrModelArtifactNotFound)
			}

			inferenceServiceID = &associations[0].ContextID
		}

		var existingAssociation schema.Association
		result := tx.Where("context_id = ? AND execution_id = ?", inferenceServiceID, serveModelExec.ID).First(&existingAssociation)

		if result.Error == gorm.ErrRecordNotFound {
			association := schema.Association{
				ContextID:   *inferenceServiceID,
				ExecutionID: serveModelExec.ID,
			}

			if err := tx.Create(&association).Error; err != nil {
				return fmt.Errorf("error creating association: %w", err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("error checking existing association: %w", result.Error)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return mapDataLayerToServeModel(serveModelExec, properties), nil
}

func (r *ServeModelRepositoryImpl) List(listOptions models.ServeModelListOptions) (*models.ListWrapper[models.ServeModel], error) {
	list := models.ListWrapper[models.ServeModel]{
		PageSize: listOptions.GetPageSize(),
	}

	serveModels := []models.ServeModel{}
	serveModelsExec := []schema.Execution{}

	query := r.db.Model(&schema.Execution{}).Where("type_id = ?", r.typeID)

	if listOptions.Name != nil {
		query = query.Where("name = ?", listOptions.Name)
	} else if listOptions.ExternalID != nil {
		query = query.Where("external_id = ?", listOptions.ExternalID)
	}

	if listOptions.InferenceServiceID != nil {
		query = query.Joins("JOIN Association ON Association.execution_id = Execution.id").
			Where("Association.context_id = ?", listOptions.InferenceServiceID)
		query = query.Scopes(scopes.PaginateWithTablePrefix(serveModels, &listOptions.Pagination, r.db, "Execution"))
	} else {
		query = query.Scopes(scopes.Paginate(serveModels, &listOptions.Pagination, r.db))
	}

	if err := query.Find(&serveModelsExec).Error; err != nil {
		return nil, fmt.Errorf("error listing serve models: %w", err)
	}

	hasMore := false
	pageSize := listOptions.GetPageSize()
	if pageSize > 0 {
		hasMore = len(serveModelsExec) > int(pageSize)
		if hasMore {
			serveModelsExec = serveModelsExec[:len(serveModelsExec)-1]
		}
	}

	for _, serveModelExec := range serveModelsExec {
		properties := []schema.ExecutionProperty{}
		if err := r.db.Where("execution_id = ?", serveModelExec.ID).Find(&properties).Error; err != nil {
			return nil, fmt.Errorf("error getting properties by serve model id: %w", err)
		}

		serveModel := mapDataLayerToServeModel(serveModelExec, properties)
		serveModels = append(serveModels, serveModel)
	}

	if hasMore && len(serveModelsExec) > 0 {
		lastModel := serveModelsExec[len(serveModelsExec)-1]
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

	list.Items = serveModels
	list.NextPageToken = listOptions.GetNextPageToken()
	list.PageSize = listOptions.GetPageSize()
	list.Size = int32(len(serveModels))

	return &list, nil
}

func mapServeModelToExecution(serveModel models.ServeModel) schema.Execution {
	if serveModel == nil {
		return schema.Execution{}
	}

	glog.Infof("Mapping ServeModel to Execution: %+v", serveModel)

	serveModelExec := schema.Execution{
		ID:     apiutils.ZeroIfNil(serveModel.GetID()),
		TypeID: apiutils.ZeroIfNil(serveModel.GetTypeID()),
	}

	if serveModel.GetAttributes() != nil {
		serveModelExec.Name = serveModel.GetAttributes().Name
		serveModelExec.ExternalID = serveModel.GetAttributes().ExternalID
		if serveModel.GetAttributes().LastKnownState != nil {
			lastKnownState := models.Execution_State_value[*serveModel.GetAttributes().LastKnownState]
			serveModelExec.LastKnownState = &lastKnownState
		}
		serveModelExec.CreateTimeSinceEpoch = apiutils.ZeroIfNil(serveModel.GetAttributes().CreateTimeSinceEpoch)
		serveModelExec.LastUpdateTimeSinceEpoch = apiutils.ZeroIfNil(serveModel.GetAttributes().LastUpdateTimeSinceEpoch)
	}

	return serveModelExec
}

func mapServeModelToExecutionProperties(serveModel models.ServeModel, executionID int32) []schema.ExecutionProperty {
	if serveModel == nil {
		return []schema.ExecutionProperty{}
	}

	properties := []schema.ExecutionProperty{}

	if serveModel.GetProperties() != nil {
		for _, prop := range *serveModel.GetProperties() {
			properties = append(properties, mapPropertiesToExecutionProperty(prop, executionID))
		}
	}

	if serveModel.GetCustomProperties() != nil {
		for _, prop := range *serveModel.GetCustomProperties() {
			properties = append(properties, mapPropertiesToExecutionProperty(prop, executionID))
		}
	}

	return properties
}

func mapPropertiesToExecutionProperty(prop models.Properties, executionID int32) schema.ExecutionProperty {
	execProp := schema.ExecutionProperty{
		ExecutionID:      executionID,
		Name:             prop.Name,
		IsCustomProperty: prop.IsCustomProperty,
		IntValue:         prop.IntValue,
		DoubleValue:      prop.DoubleValue,
		StringValue:      prop.StringValue,
		BoolValue:        prop.BoolValue,
		ByteValue:        prop.ByteValue,
		ProtoValue:       prop.ProtoValue,
	}

	return execProp
}

func mapDataLayerToServeModel(serveModel schema.Execution, properties []schema.ExecutionProperty) models.ServeModel {
	var lastKnownState *string
	if serveModel.LastKnownState != nil {
		docState := models.Execution_State_name[*serveModel.LastKnownState]
		lastKnownState = &docState
	}

	serveModelExec := models.ServeModelImpl{
		ID:     &serveModel.ID,
		TypeID: &serveModel.TypeID,
		Attributes: &models.ServeModelAttributes{
			Name:                     serveModel.Name,
			ExternalID:               serveModel.ExternalID,
			LastKnownState:           lastKnownState,
			CreateTimeSinceEpoch:     &serveModel.CreateTimeSinceEpoch,
			LastUpdateTimeSinceEpoch: &serveModel.LastUpdateTimeSinceEpoch,
		},
	}

	customProperties := []models.Properties{}
	execProperties := []models.Properties{}

	for _, prop := range properties {
		if prop.IsCustomProperty {
			customProperties = append(customProperties, mapDataLayerToExecutionProperties(prop))
		} else {
			execProperties = append(execProperties, mapDataLayerToExecutionProperties(prop))
		}
	}

	serveModelExec.CustomProperties = &customProperties
	serveModelExec.Properties = &execProperties

	return &serveModelExec
}

func mapDataLayerToExecutionProperties(prop schema.ExecutionProperty) models.Properties {
	execProp := models.Properties{
		Name:             prop.Name,
		IsCustomProperty: prop.IsCustomProperty,
		IntValue:         prop.IntValue,
		DoubleValue:      prop.DoubleValue,
		StringValue:      prop.StringValue,
		BoolValue:        prop.BoolValue,
		ByteValue:        prop.ByteValue,
		ProtoValue:       prop.ProtoValue,
	}

	return execProp
}

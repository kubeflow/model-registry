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

var ErrRegisteredModelNotFound = errors.New("registered model by id not found")

type RegisteredModelRepositoryImpl struct {
	db     *gorm.DB
	typeID int64
}

func NewRegisteredModelRepository(db *gorm.DB, typeID int64) models.RegisteredModelRepository {
	return &RegisteredModelRepositoryImpl{db: db, typeID: typeID}
}

func (r *RegisteredModelRepositoryImpl) GetByID(id int32) (models.RegisteredModel, error) {
	modelCtx := &schema.Context{}
	propertiesCtx := []schema.ContextProperty{}

	if err := r.db.Where("id = ? AND type_id = ?", id, r.typeID).First(modelCtx).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %v", ErrRegisteredModelNotFound, err)
		}

		return nil, fmt.Errorf("error getting model by id: %w", err)
	}

	if err := r.db.Where("context_id = ?", modelCtx.ID).Find(&propertiesCtx).Error; err != nil {
		return nil, fmt.Errorf("error getting properties by model id: %w", err)
	}

	return mapDataLayerToRegisteredModel(*modelCtx, propertiesCtx), nil
}

func (r *RegisteredModelRepositoryImpl) Save(model models.RegisteredModel) (models.RegisteredModel, error) {
	now := time.Now().UnixMilli()

	modelCtx := mapRegisteredModelToContext(model)
	propertiesCtx := []schema.ContextProperty{}

	modelCtx.LastUpdateTimeSinceEpoch = now

	if model.GetID() == nil {
		modelCtx.CreateTimeSinceEpoch = now
	}

	hasCustomProperties := model.GetCustomProperties() != nil

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&modelCtx).Error; err != nil {
			return fmt.Errorf("error saving model context: %w", err)
		}

		propertiesCtx = mapRegisteredModelToContextProperties(model, modelCtx.ID)
		existingCustomPropertiesCtx := []schema.ContextProperty{}

		if err := tx.Where("context_id = ? AND is_custom_property = ?", modelCtx.ID, true).Find(&existingCustomPropertiesCtx).Error; err != nil {
			return fmt.Errorf("error getting existing custom properties by model id: %w", err)
		}

		if hasCustomProperties {
			for _, existingProp := range existingCustomPropertiesCtx {
				found := false
				for _, prop := range propertiesCtx {
					if prop.Name == existingProp.Name && prop.ContextID == existingProp.ContextID && prop.IsCustomProperty == existingProp.IsCustomProperty {
						found = true
						break
					}
				}

				if !found {
					if err := tx.Delete(&existingProp).Error; err != nil {
						return fmt.Errorf("error deleting model context property: %w", err)
					}
				}
			}
		}

		for _, prop := range propertiesCtx {
			var existingProp schema.ContextProperty
			result := tx.Where("context_id = ? AND name = ? AND is_custom_property = ?",
				prop.ContextID, prop.Name, prop.IsCustomProperty).First(&existingProp)

			switch result.Error {
			case nil:
				prop.ContextID = existingProp.ContextID
				prop.Name = existingProp.Name
				prop.IsCustomProperty = existingProp.IsCustomProperty
				if err := tx.Model(&existingProp).Updates(prop).Error; err != nil {
					return fmt.Errorf("error updating model context property: %w", err)
				}
			case gorm.ErrRecordNotFound:
				if err := tx.Create(&prop).Error; err != nil {
					return fmt.Errorf("error creating model context property: %w", err)
				}
			default:
				return fmt.Errorf("error checking existing property: %w", result.Error)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return mapDataLayerToRegisteredModel(modelCtx, propertiesCtx), nil
}

func (r *RegisteredModelRepositoryImpl) List(listOptions models.RegisteredModelListOptions) (*models.ListWrapper[models.RegisteredModel], error) {
	list := models.ListWrapper[models.RegisteredModel]{
		PageSize: listOptions.GetPageSize(),
	}

	models := []models.RegisteredModel{}
	modelsCtx := []schema.Context{}

	query := r.db.Model(&schema.Context{}).Where("type_id = ?", r.typeID)
	if listOptions.Name != nil {
		query = query.Where("name = ?", listOptions.Name)
	} else if listOptions.ExternalID != nil {
		query = query.Where("external_id = ?", listOptions.ExternalID)
	}

	query = query.Scopes(scopes.Paginate(models, &listOptions.Pagination, r.db))

	if err := query.Find(&modelsCtx).Error; err != nil {
		return nil, fmt.Errorf("error listing models: %w", err)
	}

	hasMore := false
	pageSize := listOptions.GetPageSize()
	if pageSize > 0 {
		hasMore = len(modelsCtx) > int(pageSize)
		if hasMore {
			modelsCtx = modelsCtx[:len(modelsCtx)-1]
		}
	}

	for _, modelCtx := range modelsCtx {
		propertiesCtx := []schema.ContextProperty{}
		if err := r.db.Where("context_id = ?", modelCtx.ID).Find(&propertiesCtx).Error; err != nil {
			return nil, fmt.Errorf("error getting properties for model %d: %w", modelCtx.ID, err)
		}
		model := mapDataLayerToRegisteredModel(modelCtx, propertiesCtx)
		models = append(models, model)
	}

	if hasMore && len(modelsCtx) > 0 {
		lastModel := modelsCtx[len(modelsCtx)-1]
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

	list.Items = models
	list.NextPageToken = listOptions.GetNextPageToken()
	list.PageSize = listOptions.GetPageSize()
	list.Size = int32(len(models))

	return &list, nil
}

func mapRegisteredModelToContext(model models.RegisteredModel) schema.Context {
	if model == nil {
		return schema.Context{}
	}

	modelCtx := schema.Context{
		ID:     apiutils.ZeroIfNil(model.GetID()),
		TypeID: apiutils.ZeroIfNil(model.GetTypeID()),
	}

	if model.GetAttributes() != nil {
		modelCtx.Name = apiutils.ZeroIfNil(model.GetAttributes().Name)
		modelCtx.ExternalID = model.GetAttributes().ExternalID
		modelCtx.CreateTimeSinceEpoch = apiutils.ZeroIfNil(model.GetAttributes().CreateTimeSinceEpoch)
		modelCtx.LastUpdateTimeSinceEpoch = apiutils.ZeroIfNil(model.GetAttributes().LastUpdateTimeSinceEpoch)
	}

	return modelCtx
}

func mapRegisteredModelToContextProperties(model models.RegisteredModel, modelId int32) []schema.ContextProperty {
	if model == nil {
		return []schema.ContextProperty{}
	}

	propertiesCtx := []schema.ContextProperty{}

	if model.GetProperties() != nil {
		for _, prop := range *model.GetProperties() {
			propCtx := mapPropertiesToContextProperty(prop, modelId, false)
			propertiesCtx = append(propertiesCtx, propCtx)
		}
	}

	if model.GetCustomProperties() != nil {
		for _, prop := range *model.GetCustomProperties() {
			propCtx := mapPropertiesToContextProperty(prop, modelId, true)
			propertiesCtx = append(propertiesCtx, propCtx)
		}
	}

	return propertiesCtx
}

func mapPropertiesToContextProperty(prop models.Properties, contextID int32, isCustomProperty bool) schema.ContextProperty {
	propCtx := schema.ContextProperty{
		ContextID:        contextID,
		Name:             prop.Name,
		IsCustomProperty: isCustomProperty,
		IntValue:         prop.IntValue,
		DoubleValue:      prop.DoubleValue,
		StringValue:      prop.StringValue,
		BoolValue:        prop.BoolValue,
		ByteValue:        prop.ByteValue,
		ProtoValue:       prop.ProtoValue,
	}

	return propCtx
}

func mapDataLayerToRegisteredModel(modelCtx schema.Context, propertiesCtx []schema.ContextProperty) models.RegisteredModel {
	model := models.RegisteredModelImpl{
		ID:     &modelCtx.ID,
		TypeID: &modelCtx.TypeID,
		Attributes: &models.RegisteredModelAttributes{
			Name:                     &modelCtx.Name,
			ExternalID:               modelCtx.ExternalID,
			CreateTimeSinceEpoch:     &modelCtx.CreateTimeSinceEpoch,
			LastUpdateTimeSinceEpoch: &modelCtx.LastUpdateTimeSinceEpoch,
		},
	}

	customProperties := []models.Properties{}
	properties := []models.Properties{}

	for _, prop := range propertiesCtx {
		if prop.IsCustomProperty {
			customProperties = append(customProperties, mapDataLayerToProperties(prop))
		} else {
			properties = append(properties, mapDataLayerToProperties(prop))
		}
	}

	model.CustomProperties = &customProperties
	model.Properties = &properties

	return &model
}

func mapDataLayerToProperties(propCtx schema.ContextProperty) models.Properties {
	prop := models.Properties{
		Name:             propCtx.Name,
		IsCustomProperty: propCtx.IsCustomProperty,
		IntValue:         propCtx.IntValue,
		DoubleValue:      propCtx.DoubleValue,
		StringValue:      propCtx.StringValue,
		BoolValue:        propCtx.BoolValue,
		ByteValue:        propCtx.ByteValue,
		ProtoValue:       propCtx.ProtoValue,
	}

	return prop
}

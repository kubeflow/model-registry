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

var ErrModelVersionNotFound = errors.New("model version by id not found")

type ModelVersionRepositoryImpl struct {
	db     *gorm.DB
	typeID int64
}

func NewModelVersionRepository(db *gorm.DB, typeID int64) models.ModelVersionRepository {
	return &ModelVersionRepositoryImpl{
		db:     db,
		typeID: typeID,
	}
}

func (r *ModelVersionRepositoryImpl) GetByID(id int32) (models.ModelVersion, error) {
	modelVersionCtx := &schema.Context{}
	propertiesCtx := []schema.ContextProperty{}

	if err := r.db.Where("id = ? AND type_id = ?", id, r.typeID).First(&modelVersionCtx).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %v", ErrModelVersionNotFound, err)
		}

		return nil, fmt.Errorf("failed to get model version by id: %w", err)
	}

	if err := r.db.Where("context_id = ?", modelVersionCtx.ID).Find(&propertiesCtx).Error; err != nil {
		return nil, fmt.Errorf("failed to get model version properties: %w", err)
	}

	return mapDataLayerToModelVersion(*modelVersionCtx, propertiesCtx), nil
}

func (r *ModelVersionRepositoryImpl) Save(modelVersion models.ModelVersion) (models.ModelVersion, error) {
	now := time.Now().UnixMilli()

	modelVersionCtx := mapModelVersionToContext(modelVersion)
	propertiesCtx := []schema.ContextProperty{}

	modelVersionCtx.LastUpdateTimeSinceEpoch = now

	if modelVersion.GetID() == nil {
		glog.Info("Creating new ModelVersion")

		modelVersionCtx.CreateTimeSinceEpoch = now
	} else {
		glog.Infof("Updating ModelVersion %d", *modelVersion.GetID())
	}

	hasCustomProperties := modelVersion.GetCustomProperties() != nil

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&modelVersionCtx).Error; err != nil {
			return fmt.Errorf("error saving model version context: %w", err)
		}

		propertiesCtx = mapModelVersionToContextProperties(modelVersion, modelVersionCtx.ID)
		existingCustomPropertiesCtx := []schema.ContextProperty{}

		if err := tx.Where("context_id = ? AND is_custom_property = ?", modelVersionCtx.ID, true).Find(&existingCustomPropertiesCtx).Error; err != nil {
			return fmt.Errorf("error getting existing custom properties by model version id: %w", err)
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
						return fmt.Errorf("error deleting model version context property: %w", err)
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
					return fmt.Errorf("error updating model version context property: %w", err)
				}
			case gorm.ErrRecordNotFound:
				if err := tx.Create(&prop).Error; err != nil {
					return fmt.Errorf("error creating model version context property: %w", err)
				}
			default:
				return fmt.Errorf("error checking existing property: %w", result.Error)
			}
		}

		registeredModelID := int32(0)
		for _, prop := range *modelVersion.GetProperties() {
			if prop.Name == "registered_model_id" {
				registeredModelID = *prop.IntValue
				break
			}
		}

		parentsCtx := schema.ParentContext{
			ContextID:       modelVersionCtx.ID,
			ParentContextID: registeredModelID,
		}

		if err := tx.Save(&parentsCtx).Error; err != nil {
			return fmt.Errorf("error saving model version parent context: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return mapDataLayerToModelVersion(modelVersionCtx, propertiesCtx), nil
}

func (r *ModelVersionRepositoryImpl) List(listOptions models.ModelVersionListOptions) (*models.ListWrapper[models.ModelVersion], error) {
	list := models.ListWrapper[models.ModelVersion]{
		PageSize: listOptions.GetPageSize(),
	}

	modelVersions := []models.ModelVersion{}
	modelVersionsCtx := []schema.Context{}

	query := r.db.Model(&schema.Context{}).Where("type_id = ?", r.typeID)

	if listOptions.Name != nil {
		query = query.Where("name = ?", listOptions.Name)
	} else if listOptions.ExternalID != nil {
		query = query.Where("external_id = ?", listOptions.ExternalID)
	}

	if listOptions.ParentResourceID != nil {
		query = query.Joins("JOIN ParentContext ON ParentContext.context_id = Context.id").
			Where("ParentContext.parent_context_id = ?", listOptions.ParentResourceID)
	}

	query = query.Scopes(scopes.Paginate(modelVersionsCtx, &listOptions.Pagination, r.db))

	if err := query.Find(&modelVersionsCtx).Error; err != nil {
		return nil, fmt.Errorf("error listing model versions: %w", err)
	}

	hasMore := false
	pageSize := listOptions.GetPageSize()
	if pageSize > 0 {
		hasMore = len(modelVersionsCtx) > int(pageSize)
		if hasMore {
			modelVersionsCtx = modelVersionsCtx[:len(modelVersionsCtx)-1]
		}
	}

	for _, modelCtx := range modelVersionsCtx {
		propertiesCtx := []schema.ContextProperty{}
		if err := r.db.Where("context_id = ?", modelCtx.ID).Find(&propertiesCtx).Error; err != nil {
			return nil, fmt.Errorf("error getting properties for model version with id %d: %w", modelCtx.ID, err)
		}
		modelVersion := mapDataLayerToModelVersion(modelCtx, propertiesCtx)
		modelVersions = append(modelVersions, modelVersion)
	}

	if hasMore && len(modelVersionsCtx) > 0 {
		lastModel := modelVersionsCtx[len(modelVersionsCtx)-1]
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

	list.Items = modelVersions
	list.NextPageToken = listOptions.GetNextPageToken()
	list.PageSize = listOptions.GetPageSize()
	list.Size = int32(len(modelVersions))

	return &list, nil
}

func mapModelVersionToContext(modelVersion models.ModelVersion) schema.Context {
	if modelVersion == nil {
		return schema.Context{}
	}

	modelVersionCtx := schema.Context{
		ID:     apiutils.ZeroIfNil(modelVersion.GetID()),
		TypeID: apiutils.ZeroIfNil(modelVersion.GetTypeID()),
	}

	if modelVersion.GetAttributes() != nil {
		modelVersionCtx.Name = apiutils.ZeroIfNil(modelVersion.GetAttributes().Name)
		modelVersionCtx.ExternalID = modelVersion.GetAttributes().ExternalID
		modelVersionCtx.CreateTimeSinceEpoch = apiutils.ZeroIfNil(modelVersion.GetAttributes().CreateTimeSinceEpoch)
		modelVersionCtx.LastUpdateTimeSinceEpoch = apiutils.ZeroIfNil(modelVersion.GetAttributes().LastUpdateTimeSinceEpoch)
	}

	return modelVersionCtx
}

func mapModelVersionToContextProperties(modelVersion models.ModelVersion, contextID int32) []schema.ContextProperty {
	if modelVersion == nil {
		return []schema.ContextProperty{}
	}

	propertiesCtx := []schema.ContextProperty{}

	if modelVersion.GetProperties() != nil {
		for _, prop := range *modelVersion.GetProperties() {
			propertiesCtx = append(propertiesCtx, mapPropertiesToContextProperty(prop, contextID, false))
		}
	}

	if modelVersion.GetCustomProperties() != nil {
		for _, prop := range *modelVersion.GetCustomProperties() {
			propertiesCtx = append(propertiesCtx, mapPropertiesToContextProperty(prop, contextID, true))
		}
	}

	return propertiesCtx
}

func mapDataLayerToModelVersion(modelVersionCtx schema.Context, propertiesCtx []schema.ContextProperty) models.ModelVersion {
	modelVersion := models.ModelVersionImpl{
		ID:     &modelVersionCtx.ID,
		TypeID: &modelVersionCtx.TypeID,
		Attributes: &models.ModelVersionAttributes{
			Name:                     &modelVersionCtx.Name,
			ExternalID:               modelVersionCtx.ExternalID,
			CreateTimeSinceEpoch:     &modelVersionCtx.CreateTimeSinceEpoch,
			LastUpdateTimeSinceEpoch: &modelVersionCtx.LastUpdateTimeSinceEpoch,
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

	modelVersion.CustomProperties = &customProperties
	modelVersion.Properties = &properties

	return &modelVersion
}

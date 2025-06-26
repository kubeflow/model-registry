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

var ErrInferenceServiceNotFound = errors.New("inference service by id not found")

type InferenceServiceRepositoryImpl struct {
	db     *gorm.DB
	typeID int64
}

func NewInferenceServiceRepository(db *gorm.DB, typeID int64) models.InferenceServiceRepository {
	return &InferenceServiceRepositoryImpl{db: db, typeID: typeID}
}

func (r *InferenceServiceRepositoryImpl) GetByID(id int32) (models.InferenceService, error) {
	infSvcCtx := &schema.Context{}
	propertiesCtx := []schema.ContextProperty{}

	if err := r.db.Where("id = ? AND type_id = ?", id, r.typeID).First(infSvcCtx).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %v", ErrInferenceServiceNotFound, err)
		}

		return nil, fmt.Errorf("error getting inference service by id: %w", err)
	}

	if err := r.db.Where("context_id = ?", infSvcCtx.ID).Find(&propertiesCtx).Error; err != nil {
		return nil, fmt.Errorf("error getting properties by inference service id: %w", err)
	}

	return mapDataLayerToInferenceService(*infSvcCtx, propertiesCtx), nil
}

func (r *InferenceServiceRepositoryImpl) Save(inferenceService models.InferenceService) (models.InferenceService, error) {
	now := time.Now().UnixMilli()

	infSvcCtx := mapInferenceServiceToContext(inferenceService)
	propertiesCtx := []schema.ContextProperty{}

	infSvcCtx.LastUpdateTimeSinceEpoch = now

	if inferenceService.GetID() == nil {
		glog.Info("Creating new InferenceService")

		infSvcCtx.CreateTimeSinceEpoch = now
	} else {
		glog.Infof("Updating InferenceService %d", *inferenceService.GetID())
	}

	hasCustomProperties := inferenceService.GetCustomProperties() != nil

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&infSvcCtx).Error; err != nil {
			return fmt.Errorf("error saving inference service context: %w", err)
		}

		propertiesCtx = mapInferenceServiceToContextProperties(inferenceService, infSvcCtx.ID)
		existingCustomPropertiesCtx := []schema.ContextProperty{}

		if err := tx.Where("context_id = ? AND is_custom_property = ?", infSvcCtx.ID, true).Find(&existingCustomPropertiesCtx).Error; err != nil {
			return fmt.Errorf("error getting existing custom properties by inference service id: %w", err)
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
						return fmt.Errorf("error deleting inference service context property: %w", err)
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
					return fmt.Errorf("error updating inference service context property: %w", err)
				}
			case gorm.ErrRecordNotFound:
				if err := tx.Create(&prop).Error; err != nil {
					return fmt.Errorf("error creating inference service context property: %w", err)
				}
			default:
				return fmt.Errorf("error checking existing property: %w", result.Error)
			}
		}

		servEnvID := int32(0)
		for _, prop := range *inferenceService.GetProperties() {
			if prop.Name == "serving_environment_id" {
				servEnvID = *prop.IntValue
				break
			}
		}

		parentsCtx := schema.ParentContext{
			ContextID:       infSvcCtx.ID,
			ParentContextID: servEnvID,
		}

		if err := tx.Save(&parentsCtx).Error; err != nil {
			return fmt.Errorf("error saving inference service parent context: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return mapDataLayerToInferenceService(infSvcCtx, propertiesCtx), nil
}

func (r *InferenceServiceRepositoryImpl) List(listOptions models.InferenceServiceListOptions) (*models.ListWrapper[models.InferenceService], error) {
	list := models.ListWrapper[models.InferenceService]{
		PageSize: listOptions.GetPageSize(),
	}

	infSvcs := []models.InferenceService{}
	infSvcCtx := []schema.Context{}

	query := r.db.Model(&schema.Context{}).Where("type_id = ?", r.typeID)

	if listOptions.Name != nil {
		query = query.Where("name = ?", listOptions.Name)
	} else if listOptions.ExternalID != nil {
		query = query.Where("external_id = ?", listOptions.ExternalID)
	}

	needsTablePrefix := false
	if listOptions.ParentResourceID != nil {
		query = query.Joins("JOIN ParentContext ON ParentContext.context_id = Context.id").
			Where("ParentContext.parent_context_id = ?", listOptions.ParentResourceID)
		needsTablePrefix = true
	}

	if listOptions.Runtime != nil {
		query = query.Joins("JOIN ContextProperty ON ContextProperty.context_id = Context.id AND ContextProperty.name = 'runtime'").
			Where("ContextProperty.string_value = ?", listOptions.Runtime)
		needsTablePrefix = true
	}

	if needsTablePrefix {
		query = query.Scopes(scopes.PaginateWithTablePrefix(infSvcs, &listOptions.Pagination, r.db, "Context"))
	} else {
		query = query.Scopes(scopes.Paginate(infSvcs, &listOptions.Pagination, r.db))
	}

	if err := query.Find(&infSvcCtx).Error; err != nil {
		return nil, fmt.Errorf("error listing inference services: %w", err)
	}

	hasMore := false
	pageSize := listOptions.GetPageSize()
	if pageSize > 0 {
		hasMore = len(infSvcCtx) > int(pageSize)
		if hasMore {
			infSvcCtx = infSvcCtx[:len(infSvcCtx)-1]
		}
	}

	for _, modelCtx := range infSvcCtx {
		propertiesCtx := []schema.ContextProperty{}
		if err := r.db.Where("context_id = ?", modelCtx.ID).Find(&propertiesCtx).Error; err != nil {
			return nil, fmt.Errorf("error getting properties for inference service with id %d: %w", modelCtx.ID, err)
		}
		infSvc := mapDataLayerToInferenceService(modelCtx, propertiesCtx)
		infSvcs = append(infSvcs, infSvc)
	}

	if hasMore && len(infSvcCtx) > 0 {
		lastModel := infSvcCtx[len(infSvcCtx)-1]
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

	list.Items = infSvcs
	list.NextPageToken = listOptions.GetNextPageToken()
	list.PageSize = listOptions.GetPageSize()
	list.Size = int32(len(infSvcs))

	return &list, nil
}

func mapInferenceServiceToContext(inferenceService models.InferenceService) schema.Context {
	if inferenceService == nil {
		return schema.Context{}
	}

	infSvcCtx := schema.Context{
		ID:     apiutils.ZeroIfNil(inferenceService.GetID()),
		TypeID: apiutils.ZeroIfNil(inferenceService.GetTypeID()),
	}

	if inferenceService.GetAttributes() != nil {
		infSvcCtx.Name = apiutils.ZeroIfNil(inferenceService.GetAttributes().Name)
		infSvcCtx.ExternalID = inferenceService.GetAttributes().ExternalID
		infSvcCtx.CreateTimeSinceEpoch = apiutils.ZeroIfNil(inferenceService.GetAttributes().CreateTimeSinceEpoch)
		infSvcCtx.LastUpdateTimeSinceEpoch = apiutils.ZeroIfNil(inferenceService.GetAttributes().LastUpdateTimeSinceEpoch)
	}

	return infSvcCtx
}

func mapInferenceServiceToContextProperties(inferenceService models.InferenceService, contextID int32) []schema.ContextProperty {
	if inferenceService == nil {
		return []schema.ContextProperty{}
	}

	propertiesCtx := []schema.ContextProperty{}

	if inferenceService.GetProperties() != nil {
		for _, prop := range *inferenceService.GetProperties() {
			propertiesCtx = append(propertiesCtx, mapPropertiesToContextProperty(prop, contextID, false))
		}
	}

	if inferenceService.GetCustomProperties() != nil {
		for _, prop := range *inferenceService.GetCustomProperties() {
			propertiesCtx = append(propertiesCtx, mapPropertiesToContextProperty(prop, contextID, true))
		}
	}

	return propertiesCtx
}

func mapDataLayerToInferenceService(infSvcCtx schema.Context, propertiesCtx []schema.ContextProperty) models.InferenceService {
	infSvc := models.InferenceServiceImpl{
		ID:     &infSvcCtx.ID,
		TypeID: &infSvcCtx.TypeID,
		Attributes: &models.InferenceServiceAttributes{
			Name:                     &infSvcCtx.Name,
			ExternalID:               infSvcCtx.ExternalID,
			CreateTimeSinceEpoch:     &infSvcCtx.CreateTimeSinceEpoch,
			LastUpdateTimeSinceEpoch: &infSvcCtx.LastUpdateTimeSinceEpoch,
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

	infSvc.CustomProperties = &customProperties
	infSvc.Properties = &properties

	return &infSvc
}

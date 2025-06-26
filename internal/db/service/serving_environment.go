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

var ErrServingEnvironmentNotFound = errors.New("serving environment by id not found")

type ServingEnvironmentRepositoryImpl struct {
	db     *gorm.DB
	typeID int64
}

func NewServingEnvironmentRepository(db *gorm.DB, typeID int64) models.ServingEnvironmentRepository {
	return &ServingEnvironmentRepositoryImpl{db: db, typeID: typeID}
}

func (r *ServingEnvironmentRepositoryImpl) GetByID(id int32) (models.ServingEnvironment, error) {
	modelCtx := &schema.Context{}
	propertiesCtx := []schema.ContextProperty{}

	if err := r.db.Where("id = ? AND type_id = ?", id, r.typeID).First(modelCtx).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %v", ErrServingEnvironmentNotFound, err)
		}

		return nil, fmt.Errorf("error getting serving environment by id: %w", err)
	}

	if err := r.db.Where("context_id = ?", modelCtx.ID).Find(&propertiesCtx).Error; err != nil {
		return nil, fmt.Errorf("error getting properties by serving environment id: %w", err)
	}

	return mapDataLayerToServingEnvironment(*modelCtx, propertiesCtx), nil
}

func (r *ServingEnvironmentRepositoryImpl) Save(servEnv models.ServingEnvironment) (models.ServingEnvironment, error) {
	now := time.Now().UnixMilli()

	servEnvCtx := mapServingEnvironmentToContext(servEnv)
	propertiesCtx := []schema.ContextProperty{}

	servEnvCtx.LastUpdateTimeSinceEpoch = now

	if servEnv.GetID() == nil {
		servEnvCtx.CreateTimeSinceEpoch = now
	}

	hasCustomProperties := servEnv.GetCustomProperties() != nil

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&servEnvCtx).Error; err != nil {
			return fmt.Errorf("error saving serving environment context: %w", err)
		}

		propertiesCtx = mapServingEnvironmentToContextProperties(servEnv, servEnvCtx.ID)
		existingCustomPropertiesCtx := []schema.ContextProperty{}

		if err := tx.Where("context_id = ? AND is_custom_property = ?", servEnvCtx.ID, true).Find(&existingCustomPropertiesCtx).Error; err != nil {
			return fmt.Errorf("error getting existing custom properties by serving environment id: %w", err)
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
						return fmt.Errorf("error deleting serving environment context property: %w", err)
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
					return fmt.Errorf("error updating serving environment context property: %w", err)
				}
			case gorm.ErrRecordNotFound:
				if err := tx.Create(&prop).Error; err != nil {
					return fmt.Errorf("error creating serving environment context property: %w", err)
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

	return mapDataLayerToServingEnvironment(servEnvCtx, propertiesCtx), nil
}

func (r *ServingEnvironmentRepositoryImpl) List(listOptions models.ServingEnvironmentListOptions) (*models.ListWrapper[models.ServingEnvironment], error) {
	list := models.ListWrapper[models.ServingEnvironment]{
		PageSize: listOptions.GetPageSize(),
	}

	servEnvs := []models.ServingEnvironment{}
	servEnvsCtx := []schema.Context{}

	query := r.db.Model(&schema.Context{}).Where("type_id = ?", r.typeID)
	if listOptions.Name != nil {
		query = query.Where("name = ?", listOptions.Name)
	} else if listOptions.ExternalID != nil {
		query = query.Where("external_id = ?", listOptions.ExternalID)
	}

	query = query.Scopes(scopes.Paginate(servEnvs, &listOptions.Pagination, r.db))

	if err := query.Find(&servEnvsCtx).Error; err != nil {
		return nil, fmt.Errorf("error listing serving environments: %w", err)
	}

	hasMore := false
	pageSize := listOptions.GetPageSize()
	if pageSize > 0 {
		hasMore = len(servEnvsCtx) > int(pageSize)
		if hasMore {
			servEnvsCtx = servEnvsCtx[:len(servEnvsCtx)-1]
		}
	}

	for _, modelCtx := range servEnvsCtx {
		propertiesCtx := []schema.ContextProperty{}
		if err := r.db.Where("context_id = ?", modelCtx.ID).Find(&propertiesCtx).Error; err != nil {
			return nil, fmt.Errorf("error getting properties for model %d: %w", modelCtx.ID, err)
		}
		servEnv := mapDataLayerToServingEnvironment(modelCtx, propertiesCtx)
		servEnvs = append(servEnvs, servEnv)
	}

	if hasMore && len(servEnvsCtx) > 0 {
		lastModel := servEnvsCtx[len(servEnvsCtx)-1]
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

	list.Items = servEnvs
	list.NextPageToken = listOptions.GetNextPageToken()
	list.PageSize = listOptions.GetPageSize()
	list.Size = int32(len(servEnvs))

	return &list, nil
}

func mapServingEnvironmentToContext(servEnv models.ServingEnvironment) schema.Context {
	if servEnv == nil {
		return schema.Context{}
	}

	servEnvCtx := schema.Context{
		ID:     apiutils.ZeroIfNil(servEnv.GetID()),
		TypeID: apiutils.ZeroIfNil(servEnv.GetTypeID()),
	}

	if servEnv.GetAttributes() != nil {
		servEnvCtx.Name = apiutils.ZeroIfNil(servEnv.GetAttributes().Name)
		servEnvCtx.ExternalID = servEnv.GetAttributes().ExternalID
		servEnvCtx.CreateTimeSinceEpoch = apiutils.ZeroIfNil(servEnv.GetAttributes().CreateTimeSinceEpoch)
		servEnvCtx.LastUpdateTimeSinceEpoch = apiutils.ZeroIfNil(servEnv.GetAttributes().LastUpdateTimeSinceEpoch)
	}

	return servEnvCtx
}

func mapServingEnvironmentToContextProperties(servEnv models.ServingEnvironment, servEnvId int32) []schema.ContextProperty {
	if servEnv == nil {
		return []schema.ContextProperty{}
	}

	propertiesCtx := []schema.ContextProperty{}

	if servEnv.GetProperties() != nil {
		for _, prop := range *servEnv.GetProperties() {
			propCtx := mapPropertiesToContextProperty(prop, servEnvId, false)
			propertiesCtx = append(propertiesCtx, propCtx)
		}
	}

	if servEnv.GetCustomProperties() != nil {
		for _, prop := range *servEnv.GetCustomProperties() {
			propCtx := mapPropertiesToContextProperty(prop, servEnvId, true)
			propertiesCtx = append(propertiesCtx, propCtx)
		}
	}

	return propertiesCtx
}

func mapDataLayerToServingEnvironment(modelCtx schema.Context, propertiesCtx []schema.ContextProperty) models.ServingEnvironment {
	servingEnv := models.ServingEnvironmentImpl{
		ID:     &modelCtx.ID,
		TypeID: &modelCtx.TypeID,
		Attributes: &models.ServingEnvironmentAttributes{
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

	servingEnv.CustomProperties = &customProperties
	servingEnv.Properties = &properties

	return &servingEnv
}

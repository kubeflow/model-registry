package service

import (
	"errors"

	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	dbmodels "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/service"
	"gorm.io/gorm"
)

var ErrCatalogModelNotFound = errors.New("catalog model by id not found")

type CatalogModelRepositoryImpl struct {
	*service.GenericRepository[models.CatalogModel, schema.Context, schema.ContextProperty, *models.CatalogModelListOptions]
}

func NewCatalogModelRepository(db *gorm.DB, typeID int64) models.CatalogModelRepository {
	config := service.GenericRepositoryConfig[models.CatalogModel, schema.Context, schema.ContextProperty, *models.CatalogModelListOptions]{
		DB:                  db,
		TypeID:              typeID,
		EntityToSchema:      mapCatalogModelToContext,
		SchemaToEntity:      mapDataLayerToCatalogModel,
		EntityToProperties:  mapCatalogModelToContextProperties,
		NotFoundError:       ErrCatalogModelNotFound,
		EntityName:          "catalog model",
		PropertyFieldName:   "context_id",
		ApplyListFilters:    applyCatalogModelListFilters,
		IsNewEntity:         func(entity models.CatalogModel) bool { return entity.GetID() == nil },
		HasCustomProperties: func(entity models.CatalogModel) bool { return entity.GetCustomProperties() != nil },
	}

	return &CatalogModelRepositoryImpl{
		GenericRepository: service.NewGenericRepository(config),
	}
}

func (r *CatalogModelRepositoryImpl) Save(model models.CatalogModel) (models.CatalogModel, error) {
	return r.GenericRepository.Save(model, nil)
}

func (r *CatalogModelRepositoryImpl) List(listOptions models.CatalogModelListOptions) (*dbmodels.ListWrapper[models.CatalogModel], error) {
	return r.GenericRepository.List(&listOptions)
}

func applyCatalogModelListFilters(query *gorm.DB, listOptions *models.CatalogModelListOptions) *gorm.DB {
	if listOptions.Name != nil {
		query = query.Where("name LIKE ?", listOptions.Name)
	} else if listOptions.ExternalID != nil {
		query = query.Where("external_id = ?", listOptions.ExternalID)
	}
	return query
}

func mapCatalogModelToContext(model models.CatalogModel) schema.Context {
	attrs := model.GetAttributes()
	context := schema.Context{
		TypeID: *model.GetTypeID(),
	}

	if model.GetID() != nil {
		context.ID = *model.GetID()
	}

	if attrs != nil {
		if attrs.Name != nil {
			context.Name = *attrs.Name
		}
		context.ExternalID = attrs.ExternalID
		if attrs.CreateTimeSinceEpoch != nil {
			context.CreateTimeSinceEpoch = *attrs.CreateTimeSinceEpoch
		}
		if attrs.LastUpdateTimeSinceEpoch != nil {
			context.LastUpdateTimeSinceEpoch = *attrs.LastUpdateTimeSinceEpoch
		}
	}

	return context
}

func mapCatalogModelToContextProperties(model models.CatalogModel, contextID int32) []schema.ContextProperty {
	var properties []schema.ContextProperty

	if model.GetProperties() != nil {
		for _, prop := range *model.GetProperties() {
			properties = append(properties, service.MapPropertiesToContextProperty(prop, contextID, false))
		}
	}

	if model.GetCustomProperties() != nil {
		for _, prop := range *model.GetCustomProperties() {
			properties = append(properties, service.MapPropertiesToContextProperty(prop, contextID, true))
		}
	}

	return properties
}

func mapDataLayerToCatalogModel(modelCtx schema.Context, propertiesCtx []schema.ContextProperty) models.CatalogModel {
	catalogModel := &models.CatalogModelImpl{
		ID:     &modelCtx.ID,
		TypeID: &modelCtx.TypeID,
		Attributes: &models.CatalogModelAttributes{
			Name:                     &modelCtx.Name,
			ExternalID:               modelCtx.ExternalID,
			CreateTimeSinceEpoch:     &modelCtx.CreateTimeSinceEpoch,
			LastUpdateTimeSinceEpoch: &modelCtx.LastUpdateTimeSinceEpoch,
		},
	}

	properties := []dbmodels.Properties{}
	customProperties := []dbmodels.Properties{}

	for _, prop := range propertiesCtx {
		mappedProperty := service.MapContextPropertyToProperties(prop)

		if prop.IsCustomProperty {
			customProperties = append(customProperties, mappedProperty)
		} else {
			properties = append(properties, mappedProperty)
		}
	}

	catalogModel.Properties = &properties
	catalogModel.CustomProperties = &customProperties

	return catalogModel
}
package service

import (
	"errors"
	"fmt"

	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/internal/apiutils"
	dbmodels "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/db/utils"
	"gorm.io/gorm"
)

var ErrCatalogModelArtifactNotFound = errors.New("catalog model artifact by id not found")

type CatalogModelArtifactRepositoryImpl struct {
	*service.GenericRepository[models.CatalogModelArtifact, schema.Artifact, schema.ArtifactProperty, *models.CatalogModelArtifactListOptions]
}

func NewCatalogModelArtifactRepository(db *gorm.DB, typeID int32) models.CatalogModelArtifactRepository {
	config := service.GenericRepositoryConfig[models.CatalogModelArtifact, schema.Artifact, schema.ArtifactProperty, *models.CatalogModelArtifactListOptions]{
		DB:                  db,
		TypeID:              typeID,
		EntityToSchema:      mapCatalogModelArtifactToArtifact,
		SchemaToEntity:      mapDataLayerToCatalogModelArtifact,
		EntityToProperties:  mapCatalogModelArtifactToArtifactProperties,
		NotFoundError:       ErrCatalogModelArtifactNotFound,
		EntityName:          "catalog model artifact",
		PropertyFieldName:   "artifact_id",
		ApplyListFilters:    applyCatalogModelArtifactListFilters,
		IsNewEntity:         func(entity models.CatalogModelArtifact) bool { return entity.GetID() == nil },
		HasCustomProperties: func(entity models.CatalogModelArtifact) bool { return entity.GetCustomProperties() != nil },
	}

	return &CatalogModelArtifactRepositoryImpl{
		GenericRepository: service.NewGenericRepository(config),
	}
}

func (r *CatalogModelArtifactRepositoryImpl) Save(modelArtifact models.CatalogModelArtifact, parentResourceID *int32) (models.CatalogModelArtifact, error) {
	config := r.GetConfig()
	if modelArtifact.GetTypeID() == nil {
		if config.TypeID > 0 {
			modelArtifact.SetTypeID(config.TypeID)
		}
	}

	attr := modelArtifact.GetAttributes()
	if modelArtifact.GetID() == nil && attr != nil && attr.Name != nil {
		existing, err := r.lookupModelArtifactByName(*attr.Name)
		if err != nil {
			if !errors.Is(err, ErrCatalogModelArtifactNotFound) {
				return nil, fmt.Errorf("error finding existing model artifact named %s: %w", *attr.Name, err)
			}
		} else {
			modelArtifact.SetID(existing.ID)
		}
	}

	return r.GenericRepository.Save(modelArtifact, parentResourceID)
}

func (r *CatalogModelArtifactRepositoryImpl) List(listOptions models.CatalogModelArtifactListOptions) (*dbmodels.ListWrapper[models.CatalogModelArtifact], error) {
	return r.GenericRepository.List(&listOptions)
}

func (r *CatalogModelArtifactRepositoryImpl) lookupModelArtifactByName(name string) (*schema.Artifact, error) {
	var entity schema.Artifact

	config := r.GetConfig()

	if err := config.DB.Where("name = ? AND type_id = ?", name, config.TypeID).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %v", config.NotFoundError, err)
		}
		return nil, fmt.Errorf("error getting %s by name: %w", config.EntityName, err)
	}

	return &entity, nil
}

func applyCatalogModelArtifactListFilters(query *gorm.DB, listOptions *models.CatalogModelArtifactListOptions) *gorm.DB {
	if listOptions.Name != nil {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%:%s", *listOptions.Name))
	} else if listOptions.ExternalID != nil {
		query = query.Where("external_id = ?", listOptions.ExternalID)
	}

	if listOptions.ParentResourceID != nil {
		query = query.Joins(utils.BuildAttributionJoin(query)).
			Where(utils.GetColumnRef(query, &schema.Attribution{}, "context_id")+" = ?", listOptions.ParentResourceID)
	}

	return query
}

func mapCatalogModelArtifactToArtifact(catalogModelArtifact models.CatalogModelArtifact) schema.Artifact {
	if catalogModelArtifact == nil {
		return schema.Artifact{}
	}

	artifact := schema.Artifact{
		ID:     apiutils.ZeroIfNil(catalogModelArtifact.GetID()),
		TypeID: apiutils.ZeroIfNil(catalogModelArtifact.GetTypeID()),
	}

	if catalogModelArtifact.GetAttributes() != nil {
		artifact.Name = catalogModelArtifact.GetAttributes().Name
		artifact.URI = catalogModelArtifact.GetAttributes().URI
		artifact.ExternalID = catalogModelArtifact.GetAttributes().ExternalID
		artifact.CreateTimeSinceEpoch = apiutils.ZeroIfNil(catalogModelArtifact.GetAttributes().CreateTimeSinceEpoch)
		artifact.LastUpdateTimeSinceEpoch = apiutils.ZeroIfNil(catalogModelArtifact.GetAttributes().LastUpdateTimeSinceEpoch)
	}

	return artifact
}

func mapCatalogModelArtifactToArtifactProperties(catalogModelArtifact models.CatalogModelArtifact, artifactID int32) []schema.ArtifactProperty {
	if catalogModelArtifact == nil {
		return []schema.ArtifactProperty{}
	}

	properties := []schema.ArtifactProperty{}

	if catalogModelArtifact.GetProperties() != nil {
		for _, prop := range *catalogModelArtifact.GetProperties() {
			properties = append(properties, service.MapPropertiesToArtifactProperty(prop, artifactID, false))
		}
	}

	if catalogModelArtifact.GetCustomProperties() != nil {
		for _, prop := range *catalogModelArtifact.GetCustomProperties() {
			properties = append(properties, service.MapPropertiesToArtifactProperty(prop, artifactID, true))
		}
	}

	return properties
}

func mapDataLayerToCatalogModelArtifact(artifact schema.Artifact, artProperties []schema.ArtifactProperty) models.CatalogModelArtifact {
	catalogModelArtifact := models.CatalogModelArtifactImpl{
		ID:     &artifact.ID,
		TypeID: &artifact.TypeID,
		Attributes: &models.CatalogModelArtifactAttributes{
			Name:                     artifact.Name,
			URI:                      artifact.URI,
			ArtifactType:             apiutils.Of(models.CatalogModelArtifactType),
			ExternalID:               artifact.ExternalID,
			CreateTimeSinceEpoch:     &artifact.CreateTimeSinceEpoch,
			LastUpdateTimeSinceEpoch: &artifact.LastUpdateTimeSinceEpoch,
		},
	}

	customProperties := []dbmodels.Properties{}
	properties := []dbmodels.Properties{}

	for _, prop := range artProperties {
		if prop.IsCustomProperty {
			customProperties = append(customProperties, service.MapArtifactPropertyToProperties(prop))
		} else {
			properties = append(properties, service.MapArtifactPropertyToProperties(prop))
		}
	}

	catalogModelArtifact.CustomProperties = &customProperties
	catalogModelArtifact.Properties = &properties

	return &catalogModelArtifact
}

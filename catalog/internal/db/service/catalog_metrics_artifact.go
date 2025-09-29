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

var ErrCatalogMetricsArtifactNotFound = errors.New("catalog metrics artifact by id not found")

type CatalogMetricsArtifactRepositoryImpl struct {
	*service.GenericRepository[models.CatalogMetricsArtifact, schema.Artifact, schema.ArtifactProperty, *models.CatalogMetricsArtifactListOptions]
}

func NewCatalogMetricsArtifactRepository(db *gorm.DB, typeID int64) models.CatalogMetricsArtifactRepository {
	config := service.GenericRepositoryConfig[models.CatalogMetricsArtifact, schema.Artifact, schema.ArtifactProperty, *models.CatalogMetricsArtifactListOptions]{
		DB:                  db,
		TypeID:              typeID,
		EntityToSchema:      mapCatalogMetricsArtifactToArtifact,
		SchemaToEntity:      mapDataLayerToCatalogMetricsArtifact,
		EntityToProperties:  mapCatalogMetricsArtifactToArtifactProperties,
		NotFoundError:       ErrCatalogMetricsArtifactNotFound,
		EntityName:          "catalog metrics artifact",
		PropertyFieldName:   "artifact_id",
		ApplyListFilters:    applyCatalogMetricsArtifactListFilters,
		IsNewEntity:         func(entity models.CatalogMetricsArtifact) bool { return entity.GetID() == nil },
		HasCustomProperties: func(entity models.CatalogMetricsArtifact) bool { return entity.GetCustomProperties() != nil },
	}

	return &CatalogMetricsArtifactRepositoryImpl{
		GenericRepository: service.NewGenericRepository(config),
	}
}

func (r *CatalogMetricsArtifactRepositoryImpl) List(listOptions models.CatalogMetricsArtifactListOptions) (*dbmodels.ListWrapper[models.CatalogMetricsArtifact], error) {
	return r.GenericRepository.List(&listOptions)
}

func (r *CatalogMetricsArtifactRepositoryImpl) Save(ma models.CatalogMetricsArtifact, parentResourceID *int32) (models.CatalogMetricsArtifact, error) {
	attr := ma.GetAttributes()
	if attr == nil {
		return ma, fmt.Errorf("invalid artifact: nil attributes")
	}

	switch attr.MetricsType {
	case models.MetricsTypeAccuracy, models.MetricsTypePerformance:
		// OK
	default:
		return ma, fmt.Errorf("invalid artifact: unknown metrics type: %s", attr.MetricsType)
	}

	return r.GenericRepository.Save(ma, parentResourceID)
}

func applyCatalogMetricsArtifactListFilters(query *gorm.DB, listOptions *models.CatalogMetricsArtifactListOptions) *gorm.DB {
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

func mapCatalogMetricsArtifactToArtifact(catalogMetricsArtifact models.CatalogMetricsArtifact) schema.Artifact {
	if catalogMetricsArtifact == nil {
		return schema.Artifact{}
	}

	artifact := schema.Artifact{
		ID:     apiutils.ZeroIfNil(catalogMetricsArtifact.GetID()),
		TypeID: apiutils.ZeroIfNil(catalogMetricsArtifact.GetTypeID()),
	}

	if catalogMetricsArtifact.GetAttributes() != nil {
		artifact.Name = catalogMetricsArtifact.GetAttributes().Name
		artifact.ExternalID = catalogMetricsArtifact.GetAttributes().ExternalID
		artifact.CreateTimeSinceEpoch = apiutils.ZeroIfNil(catalogMetricsArtifact.GetAttributes().CreateTimeSinceEpoch)
		artifact.LastUpdateTimeSinceEpoch = apiutils.ZeroIfNil(catalogMetricsArtifact.GetAttributes().LastUpdateTimeSinceEpoch)
	}

	return artifact
}

func mapCatalogMetricsArtifactToArtifactProperties(catalogMetricsArtifact models.CatalogMetricsArtifact, artifactID int32) []schema.ArtifactProperty {
	if catalogMetricsArtifact == nil {
		return []schema.ArtifactProperty{}
	}

	properties := []schema.ArtifactProperty{}

	// Add the metricsType as a property
	if catalogMetricsArtifact.GetAttributes() != nil {
		metricsTypeProp := dbmodels.Properties{
			Name:        "metricsType",
			StringValue: apiutils.Of(string(catalogMetricsArtifact.GetAttributes().MetricsType)),
		}
		properties = append(properties, service.MapPropertiesToArtifactProperty(metricsTypeProp, artifactID, false))
	}

	if catalogMetricsArtifact.GetProperties() != nil {
		for _, prop := range *catalogMetricsArtifact.GetProperties() {
			properties = append(properties, service.MapPropertiesToArtifactProperty(prop, artifactID, false))
		}
	}

	if catalogMetricsArtifact.GetCustomProperties() != nil {
		for _, prop := range *catalogMetricsArtifact.GetCustomProperties() {
			properties = append(properties, service.MapPropertiesToArtifactProperty(prop, artifactID, true))
		}
	}

	return properties
}

func mapDataLayerToCatalogMetricsArtifact(artifact schema.Artifact, artProperties []schema.ArtifactProperty) models.CatalogMetricsArtifact {
	catalogMetricsArtifact := models.CatalogMetricsArtifactImpl{
		ID:     &artifact.ID,
		TypeID: &artifact.TypeID,
		Attributes: &models.CatalogMetricsArtifactAttributes{
			Name:                     artifact.Name,
			ArtifactType:             apiutils.Of(models.CatalogMetricsArtifactType),
			ExternalID:               artifact.ExternalID,
			CreateTimeSinceEpoch:     &artifact.CreateTimeSinceEpoch,
			LastUpdateTimeSinceEpoch: &artifact.LastUpdateTimeSinceEpoch,
		},
	}

	customProperties := []dbmodels.Properties{}
	properties := []dbmodels.Properties{}

	for _, prop := range artProperties {
		mappedProperty := service.MapArtifactPropertyToProperties(prop)

		// Extract metricsType from properties and set it as an attribute
		if mappedProperty.Name == "metricsType" && !prop.IsCustomProperty {
			if mappedProperty.StringValue != nil {
				catalogMetricsArtifact.Attributes.MetricsType = models.MetricsType(*mappedProperty.StringValue)
			}
		} else if prop.IsCustomProperty {
			customProperties = append(customProperties, mappedProperty)
		} else {
			properties = append(properties, mappedProperty)
		}
	}

	catalogMetricsArtifact.CustomProperties = &customProperties
	catalogMetricsArtifact.Properties = &properties

	return &catalogMetricsArtifact
}

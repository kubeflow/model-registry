package service

import (
	"errors"
	"fmt"

	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/internal/datastore"
	dbmodels "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/utils"
	"gorm.io/gorm"
)

var ErrCatalogArtifactNotFound = errors.New("catalog artifact by id not found")

type CatalogArtifactRepositoryImpl struct {
	db       *gorm.DB
	idToName map[int64]string
	nameToID datastore.ArtifactTypeMap
}

func NewCatalogArtifactRepository(db *gorm.DB, artifactTypes datastore.ArtifactTypeMap) models.CatalogArtifactRepository {
	idToName := make(map[int64]string, len(artifactTypes))
	for name, id := range artifactTypes {
		idToName[id] = name
	}

	return &CatalogArtifactRepositoryImpl{
		db:       db,
		nameToID: artifactTypes,
		idToName: idToName,
	}
}

func (r *CatalogArtifactRepositoryImpl) GetByID(id int32) (models.CatalogArtifact, error) {
	artifact := &schema.Artifact{}
	properties := []schema.ArtifactProperty{}

	if err := r.db.Where("id = ?", id).First(artifact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.CatalogArtifact{}, fmt.Errorf("%w: %v", ErrCatalogArtifactNotFound, err)
		}
		return models.CatalogArtifact{}, fmt.Errorf("error getting catalog artifact by id: %w", err)
	}

	if err := r.db.Where("artifact_id = ?", artifact.ID).Find(&properties).Error; err != nil {
		return models.CatalogArtifact{}, fmt.Errorf("error getting properties by artifact id: %w", err)
	}

	// Use the same logic as mapDataLayerToCatalogArtifact to handle artifact types
	mappedArtifact, err := r.mapDataLayerToCatalogArtifact(*artifact, properties)
	if err != nil {
		return models.CatalogArtifact{}, fmt.Errorf("error mapping catalog artifact: %w", err)
	}

	return mappedArtifact, nil
}

func (r *CatalogArtifactRepositoryImpl) List(listOptions models.CatalogArtifactListOptions) (*dbmodels.ListWrapper[models.CatalogArtifact], error) {
	list := dbmodels.ListWrapper[models.CatalogArtifact]{
		PageSize: listOptions.GetPageSize(),
	}

	artifacts := []models.CatalogArtifact{}
	artifactsArt := []schema.Artifact{}

	query := r.db.Model(&schema.Artifact{})

	// Apply filters similar to the internal artifact service
	if listOptions.Name != nil {
		// Name is not prefixed with the parent resource id to allow for filtering by name only
		// Parent resource Id is used later to filter by Attribution.context_id
		query = query.Where("name LIKE ?", fmt.Sprintf("%%:%s", *listOptions.Name))
	} else if listOptions.ExternalID != nil {
		query = query.Where("external_id = ?", listOptions.ExternalID)
	}

	// Filter by artifact type if specified
	if listOptions.ArtifactType != nil {
		typeID, err := r.getTypeIDFromArtifactType(*listOptions.ArtifactType)
		if err != nil {
			return nil, fmt.Errorf("invalid catalog artifact type %s: %w", *listOptions.ArtifactType, err)
		}
		query = query.Where("type_id = ?", typeID)
	} else {
		// Only include catalog artifact types
		catalogTypeIDs := []int64{}
		for _, typeID := range r.nameToID {
			catalogTypeIDs = append(catalogTypeIDs, typeID)
		}
		query = query.Where("type_id IN ?", catalogTypeIDs)
	}

	// Apply parent resource filtering if specified
	if listOptions.ParentResourceID != nil {
		// Proper GORM JOIN: Use helper that respects naming strategy
		query = query.Joins(utils.BuildAttributionJoin(query)).
			Where(utils.GetColumnRef(query, &schema.Attribution{}, "context_id")+" = ?", listOptions.ParentResourceID).
			Select(utils.GetTableName(query, &schema.Artifact{}) + ".*") // Explicitly select from Artifact table to avoid ambiguity
	}

	// Apply pagination
	pageSize := listOptions.GetPageSize()
	if pageSize > 0 {
		query = query.Limit(int(pageSize + 1)) // +1 to check if there are more results
	}

	if err := query.Find(&artifactsArt).Error; err != nil {
		return nil, fmt.Errorf("error listing catalog artifacts: %w", err)
	}

	// Check for pagination
	hasMore := false
	if pageSize > 0 {
		hasMore = len(artifactsArt) > int(pageSize)
		if hasMore {
			artifactsArt = artifactsArt[:len(artifactsArt)-1]
		}
	}

	// Map each artifact with its properties
	for _, artifactArt := range artifactsArt {
		properties := []schema.ArtifactProperty{}
		if err := r.db.Where("artifact_id = ?", artifactArt.ID).Find(&properties).Error; err != nil {
			return nil, fmt.Errorf("error getting properties by artifact id: %w", err)
		}

		artifact, err := r.mapDataLayerToCatalogArtifact(artifactArt, properties)
		if err != nil {
			return nil, fmt.Errorf("error mapping catalog artifact: %w", err)
		}
		artifacts = append(artifacts, artifact)
	}

	// Handle pagination token
	if hasMore && len(artifactsArt) > 0 {
		// For now, use a simple pagination approach
		listOptions.NextPageToken = nil // Implementation can be enhanced later
	} else {
		listOptions.NextPageToken = nil
	}

	list.Items = artifacts
	list.NextPageToken = listOptions.GetNextPageToken()
	list.Size = int32(len(artifacts))

	return &list, nil
}

// getTypeIDFromArtifactType maps catalog artifact type strings to their corresponding type IDs
func (r *CatalogArtifactRepositoryImpl) getTypeIDFromArtifactType(artifactType string) (int64, error) {
	switch artifactType {
	case "model-artifact":
		return r.nameToID[CatalogModelArtifactTypeName], nil
	case "metrics-artifact":
		return r.nameToID[CatalogMetricsArtifactTypeName], nil
	default:
		return 0, fmt.Errorf("unsupported catalog artifact type: %s", artifactType)
	}
}

func (r *CatalogArtifactRepositoryImpl) mapDataLayerToCatalogArtifact(artifact schema.Artifact, properties []schema.ArtifactProperty) (models.CatalogArtifact, error) {
	artToReturn := models.CatalogArtifact{}

	typeName := r.idToName[int64(artifact.TypeID)]

	switch typeName {
	case CatalogModelArtifactTypeName:
		modelArtifact := mapDataLayerToCatalogModelArtifact(artifact, properties)
		artToReturn.CatalogModelArtifact = &modelArtifact
	case CatalogMetricsArtifactTypeName:
		metricsArtifact := mapDataLayerToCatalogMetricsArtifact(artifact, properties)
		artToReturn.CatalogMetricsArtifact = &metricsArtifact
	default:
		return models.CatalogArtifact{}, fmt.Errorf("invalid catalog artifact type: %s=%d (expected: %v)", typeName, artifact.TypeID, r.idToName)
	}

	return artToReturn, nil
}

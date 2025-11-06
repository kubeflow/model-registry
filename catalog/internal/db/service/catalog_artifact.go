package service

import (
	"errors"
	"fmt"

	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/internal/datastore"
	dbmodels "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/scopes"
	"github.com/kubeflow/model-registry/internal/db/utils"
	"github.com/kubeflow/model-registry/pkg/api"
	"gorm.io/gorm"
)

var ErrCatalogArtifactNotFound = errors.New("catalog artifact by id not found")

type CatalogArtifactRepositoryImpl struct {
	db       *gorm.DB
	idToName map[int32]string
	nameToID datastore.ArtifactTypeMap
}

func NewCatalogArtifactRepository(db *gorm.DB, artifactTypes datastore.ArtifactTypeMap) models.CatalogArtifactRepository {
	idToName := make(map[int32]string, len(artifactTypes))
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

	// Filter by artifact type(s) if specified
	if len(listOptions.ArtifactTypesFilter) > 0 {
		// Handle multiple artifact types
		typeIDs := []int32{}
		for _, artifactType := range listOptions.ArtifactTypesFilter {
			// Handle "null" string as invalid artifact type
			if artifactType == "null" || artifactType == "" {
				return nil, fmt.Errorf("invalid artifact type: empty or null value provided: %w", api.ErrBadRequest)
			}
			typeID, err := r.getTypeIDFromArtifactType(artifactType)
			if err != nil {
				return nil, fmt.Errorf("invalid catalog artifact type %s: %w", artifactType, err)
			}
			typeIDs = append(typeIDs, typeID)
		}
		query = query.Where("type_id IN ?", typeIDs)
	} else if listOptions.ArtifactType != nil {
		// Handle single artifact type for backward compatibility
		// Handle "null" string as invalid artifact type
		if *listOptions.ArtifactType == "null" || *listOptions.ArtifactType == "" {
			return nil, fmt.Errorf("invalid artifact type: empty or null value provided: %w", api.ErrBadRequest)
		}
		typeID, err := r.getTypeIDFromArtifactType(*listOptions.ArtifactType)
		if err != nil {
			return nil, fmt.Errorf("invalid catalog artifact type %s: %w", *listOptions.ArtifactType, err)
		}
		query = query.Where("type_id = ?", typeID)
	} else {
		// Only include catalog artifact types
		catalogTypeIDs := []int32{}
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

	orderBy := listOptions.GetOrderBy()
	sortOrder := listOptions.GetSortOrder()
	nextPageToken := listOptions.GetNextPageToken()
	pageSize := listOptions.GetPageSize()

	// Handle NAME ordering specially (catalog-specific) to avoid string-to-integer cast issues
	if orderBy == "NAME" {
		artifactTable := utils.GetTableName(query, &schema.Artifact{})
		query = ApplyNameOrdering(query, artifactTable, sortOrder, nextPageToken, pageSize)
	} else {
		// For non-NAME ordering, use standard pagination
		pagination := &dbmodels.Pagination{
			PageSize:      &pageSize,
			OrderBy:       &orderBy,
			SortOrder:     &sortOrder,
			NextPageToken: &nextPageToken,
		}

		// Use catalog-specific allowed columns (includes NAME)
		query = query.Scopes(scopes.PaginateWithOptions(artifactsArt, pagination, r.db, "Artifact", CatalogOrderByColumns))
	}

	if err := query.Find(&artifactsArt).Error; err != nil {
		return nil, fmt.Errorf("error listing catalog artifacts: %w", err)
	}

	hasMore := false
	if pageSize > 0 {
		hasMore = len(artifactsArt) > int(pageSize)
		if hasMore {
			artifactsArt = artifactsArt[:len(artifactsArt)-1] // Remove the extra item used for hasMore detection
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

	// Handle pagination token - generate token when there are more pages
	if hasMore && len(artifactsArt) > 0 {
		// Use the last artifact to generate pagination token
		lastArtifact := artifactsArt[len(artifactsArt)-1]
		nextToken := r.createPaginationToken(lastArtifact, listOptions)
		listOptions.NextPageToken = &nextToken
	} else {
		listOptions.NextPageToken = nil
	}

	list.Items = artifacts
	list.NextPageToken = listOptions.GetNextPageToken()
	list.Size = int32(len(artifacts))

	return &list, nil
}

// getTypeIDFromArtifactType maps catalog artifact type strings to their corresponding type IDs
func (r *CatalogArtifactRepositoryImpl) getTypeIDFromArtifactType(artifactType string) (int32, error) {
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

	typeName := r.idToName[artifact.TypeID]

	switch typeName {
	case CatalogModelArtifactTypeName:
		artToReturn.CatalogModelArtifact = mapDataLayerToCatalogModelArtifact(artifact, properties)
	case CatalogMetricsArtifactTypeName:
		artToReturn.CatalogMetricsArtifact = mapDataLayerToCatalogMetricsArtifact(artifact, properties)
	default:
		return models.CatalogArtifact{}, fmt.Errorf("invalid catalog artifact type: %s=%d (expected: %v)", typeName, artifact.TypeID, r.idToName)
	}

	return artToReturn, nil
}

// createPaginationToken generates a pagination token based on the last artifact and ordering
func (r *CatalogArtifactRepositoryImpl) createPaginationToken(artifact schema.Artifact, listOptions models.CatalogArtifactListOptions) string {
	orderBy := listOptions.GetOrderBy()
	value := ""

	// Generate token value based on ordering field
	switch orderBy {
	case "ID":
		value = fmt.Sprintf("%d", artifact.ID)
	case "CREATE_TIME":
		value = fmt.Sprintf("%d", artifact.CreateTimeSinceEpoch)
	case "LAST_UPDATE_TIME":
		value = fmt.Sprintf("%d", artifact.LastUpdateTimeSinceEpoch)
	case "NAME":
		return CreateNamePaginationToken(artifact.ID, artifact.Name)
	default:
		// Default to ID ordering
		value = fmt.Sprintf("%d", artifact.ID)
	}

	return scopes.CreateNextPageToken(artifact.ID, value)
}

func (r *CatalogArtifactRepositoryImpl) DeleteByParentID(artifactTypeName string, parentResourceID int32) error {
	typeID, ok := r.nameToID[artifactTypeName]
	if !ok {
		return fmt.Errorf("unknown artifact type name: %s", artifactTypeName)
	}

	return r.db.Exec(`DELETE FROM "Artifact" WHERE id IN (SELECT artifact_id from "Attribution" INNER JOIN "Artifact" artifact ON artifact.id=artifact_id where context_id=? and type_id=?)`, parentResourceID, typeID).Error
}

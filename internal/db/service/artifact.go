package service

import (
	"errors"
	"fmt"

	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/scopes"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"gorm.io/gorm"
)

var ErrArtifactNotFound = errors.New("artifact by id not found")

type ArtifactRepositoryImpl struct {
	db                  *gorm.DB
	modelArtifactTypeID int64
	docArtifactTypeID   int64
	dataSetTypeID       int64
	metricTypeID        int64
	parameterTypeID     int64
	metricHistoryTypeID int64
}

func NewArtifactRepository(db *gorm.DB, modelArtifactTypeID int64, docArtifactTypeID int64, dataSetTypeID int64, metricTypeID int64, parameterTypeID int64, metricHistoryTypeID int64) models.ArtifactRepository {
	return &ArtifactRepositoryImpl{
		db:                  db,
		modelArtifactTypeID: modelArtifactTypeID,
		docArtifactTypeID:   docArtifactTypeID,
		dataSetTypeID:       dataSetTypeID,
		metricTypeID:        metricTypeID,
		parameterTypeID:     parameterTypeID,
		metricHistoryTypeID: metricHistoryTypeID,
	}
}

func (r *ArtifactRepositoryImpl) GetByID(id int32) (models.Artifact, error) {
	var artifactWithParent artifactWithParentID
	properties := []schema.ArtifactProperty{}

	// Use LEFT JOIN to get both artifact and parent resource ID in one query
	if err := r.db.Model(&schema.Artifact{}).
		Select("Artifact.*, Attribution.context_id").
		Joins("LEFT JOIN Attribution ON Attribution.artifact_id = Artifact.id").
		Where("Artifact.id = ?", id).
		First(&artifactWithParent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Artifact{}, fmt.Errorf("%w: %v", ErrArtifactNotFound, err)
		}

		return models.Artifact{}, fmt.Errorf("error getting artifact by id: %w", err)
	}

	if err := r.db.Where("artifact_id = ?", artifactWithParent.ID).Find(&properties).Error; err != nil {
		return models.Artifact{}, fmt.Errorf("error getting properties by artifact id: %w", err)
	}

	// Use the same logic as mapDataLayerToArtifact to handle all artifact types
	mappedArtifact, err := r.mapDataLayerToArtifact(artifactWithParent.Artifact, properties)
	if err != nil {
		return models.Artifact{}, fmt.Errorf("error mapping artifact: %w", err)
	}

	// Set parent resource ID from the JOIN result
	if artifactWithParent.ParentResourceID != nil {
		parentResourceIDStr := fmt.Sprintf("%d", *artifactWithParent.ParentResourceID)
		mappedArtifact.ParentResourceID = &parentResourceIDStr
	}

	return mappedArtifact, nil
}

// artifactWithParentID represents an artifact with its parent resource ID from Attribution table
type artifactWithParentID struct {
	schema.Artifact
	ParentResourceID *int32 `gorm:"column:context_id"`
}

func (r *ArtifactRepositoryImpl) List(listOptions models.ArtifactListOptions) (*models.ListWrapper[models.Artifact], error) {
	list := models.ListWrapper[models.Artifact]{
		PageSize: listOptions.GetPageSize(),
	}

	artifacts := []models.Artifact{}
	artifactsWithParent := []artifactWithParentID{}

	// Always LEFT JOIN with Attribution to get parent resource ID in a single query
	query := r.db.Model(&schema.Artifact{}).
		Select("Artifact.*, Attribution.context_id").
		Joins("LEFT JOIN Attribution ON Attribution.artifact_id = Artifact.id")

	// Exclude metric history records - they should only be returned via metric history endpoints
	query = query.Where("type_id != ?", r.metricHistoryTypeID)

	if listOptions.Name != nil {
		// Name is not prefixed with the parent resource id to allow for filtering by name only
		// Parent resource Id is used later to filter by Attribution.context_id
		query = query.Where("name LIKE ?", fmt.Sprintf("%%:%s", *listOptions.Name))
	} else if listOptions.ExternalID != nil {
		query = query.Where("external_id = ?", listOptions.ExternalID)
	}

	// Filter by artifact type if specified
	if listOptions.ArtifactType != nil {
		// Handle "null" string as invalid artifact type
		if *listOptions.ArtifactType == "null" || *listOptions.ArtifactType == "" {
			return nil, fmt.Errorf("invalid artifact type: empty or null value provided: %w", api.ErrBadRequest)
		}
		typeID, err := r.getTypeIDFromArtifactType(*listOptions.ArtifactType)
		if err != nil {
			return nil, fmt.Errorf("invalid artifact type %s: %w", *listOptions.ArtifactType, api.ErrBadRequest)
		}
		query = query.Where("Artifact.type_id = ?", typeID)
	}

	// Filter by parent resource ID if specified
	if listOptions.ParentResourceID != nil {
		query = query.Where("Attribution.context_id = ?", listOptions.ParentResourceID)
	}

	// Use table-prefixed pagination to handle JOINs properly
	query = query.Scopes(scopes.PaginateWithTablePrefix(artifacts, &listOptions.Pagination, r.db, "Artifact"))

	if err := query.Find(&artifactsWithParent).Error; err != nil {
		return nil, fmt.Errorf("error listing artifacts: %w", err)
	}

	hasMore := false
	pageSize := listOptions.GetPageSize()
	if pageSize > 0 {
		hasMore = len(artifactsWithParent) > int(pageSize)
		if hasMore {
			artifactsWithParent = artifactsWithParent[:len(artifactsWithParent)-1]
		}
	}

	for _, artifactWithParent := range artifactsWithParent {
		properties := []schema.ArtifactProperty{}
		if err := r.db.Where("artifact_id = ?", artifactWithParent.ID).Find(&properties).Error; err != nil {
			return nil, fmt.Errorf("error getting properties by artifact id: %w", err)
		}

		artifact, err := r.mapDataLayerToArtifact(artifactWithParent.Artifact, properties)
		if err != nil {
			return nil, fmt.Errorf("error mapping artifact: %w", err)
		}

		// Set parent resource ID from the JOIN result
		if artifactWithParent.ParentResourceID != nil {
			parentResourceIDStr := fmt.Sprintf("%d", *artifactWithParent.ParentResourceID)
			artifact.ParentResourceID = &parentResourceIDStr
		}

		artifacts = append(artifacts, artifact)
	}

	if hasMore && len(artifactsWithParent) > 0 {
		lastArtifact := artifactsWithParent[len(artifactsWithParent)-1]
		orderBy := listOptions.GetOrderBy()
		value := ""
		if orderBy != "" {
			switch orderBy {
			case "ID":
				value = fmt.Sprintf("%d", lastArtifact.ID)
			case "CREATE_TIME":
				value = fmt.Sprintf("%d", lastArtifact.CreateTimeSinceEpoch)
			case "LAST_UPDATE_TIME":
				value = fmt.Sprintf("%d", lastArtifact.LastUpdateTimeSinceEpoch)
			default:
				value = fmt.Sprintf("%d", lastArtifact.ID)
			}
		}
		nextToken := scopes.CreateNextPageToken(lastArtifact.ID, value)
		listOptions.NextPageToken = &nextToken
	} else {
		listOptions.NextPageToken = nil
	}

	list.Items = artifacts
	list.NextPageToken = listOptions.GetNextPageToken()
	list.PageSize = listOptions.GetPageSize()
	list.Size = int32(len(artifacts))

	return &list, nil
}

// getTypeIDFromArtifactType maps artifact type strings to their corresponding type IDs
func (r *ArtifactRepositoryImpl) getTypeIDFromArtifactType(artifactType string) (int64, error) {
	switch artifactType {
	case string(openapi.ARTIFACTTYPEQUERYPARAM_MODEL_ARTIFACT):
		return r.modelArtifactTypeID, nil
	case string(openapi.ARTIFACTTYPEQUERYPARAM_DOC_ARTIFACT):
		return r.docArtifactTypeID, nil
	case string(openapi.ARTIFACTTYPEQUERYPARAM_DATASET_ARTIFACT):
		return r.dataSetTypeID, nil
	case string(openapi.ARTIFACTTYPEQUERYPARAM_METRIC):
		return r.metricTypeID, nil
	case string(openapi.ARTIFACTTYPEQUERYPARAM_PARAMETER):
		return r.parameterTypeID, nil
	default:
		return 0, fmt.Errorf("unsupported artifact type: %s: %w", artifactType, api.ErrBadRequest)
	}
}

func (r *ArtifactRepositoryImpl) mapDataLayerToArtifact(artifact schema.Artifact, properties []schema.ArtifactProperty) (models.Artifact, error) {
	artToReturn := models.Artifact{}

	switch artifact.TypeID {
	case int32(r.modelArtifactTypeID):
		modelArtifact := mapDataLayerToModelArtifact(artifact, properties)
		artToReturn.ModelArtifact = &modelArtifact
	case int32(r.docArtifactTypeID):
		docArtifact := mapDataLayerToDocArtifact(artifact, properties)
		artToReturn.DocArtifact = &docArtifact
	case int32(r.dataSetTypeID):
		dataSet := mapDataLayerToDataSet(artifact, properties)
		artToReturn.DataSet = &dataSet
	case int32(r.metricTypeID):
		metric := mapDataLayerToMetric(artifact, properties)
		artToReturn.Metric = &metric
	case int32(r.parameterTypeID):
		parameter := mapDataLayerToParameter(artifact, properties)
		artToReturn.Parameter = &parameter
	default:
		return models.Artifact{}, fmt.Errorf("invalid artifact type: %d (expected: modelArtifact=%d, docArtifact=%d, dataSet=%d, metric=%d, parameter=%d, metricHistory=%d [filtered])",
			artifact.TypeID, r.modelArtifactTypeID, r.docArtifactTypeID, r.dataSetTypeID, r.metricTypeID, r.parameterTypeID, r.metricHistoryTypeID)
	}

	return artToReturn, nil
}

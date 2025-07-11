package service

import (
	"errors"
	"fmt"

	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/scopes"
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
}

func NewArtifactRepository(db *gorm.DB, modelArtifactTypeID int64, docArtifactTypeID int64, dataSetTypeID int64, metricTypeID int64, parameterTypeID int64) models.ArtifactRepository {
	return &ArtifactRepositoryImpl{
		db:                  db,
		modelArtifactTypeID: modelArtifactTypeID,
		docArtifactTypeID:   docArtifactTypeID,
		dataSetTypeID:       dataSetTypeID,
		metricTypeID:        metricTypeID,
		parameterTypeID:     parameterTypeID,
	}
}

func (r *ArtifactRepositoryImpl) GetByID(id int32) (models.Artifact, error) {
	artifact := &schema.Artifact{}
	properties := []schema.ArtifactProperty{}

	if err := r.db.Where("id = ?", id).First(artifact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Artifact{}, fmt.Errorf("%w: %v", ErrArtifactNotFound, err)
		}

		return models.Artifact{}, fmt.Errorf("error getting artifact by id: %w", err)
	}

	if err := r.db.Where("artifact_id = ?", artifact.ID).Find(&properties).Error; err != nil {
		return models.Artifact{}, fmt.Errorf("error getting properties by artifact id: %w", err)
	}

	// Use the same logic as mapDataLayerToArtifact to handle all artifact types
	mappedArtifact, err := r.mapDataLayerToArtifact(*artifact, properties)
	if err != nil {
		return models.Artifact{}, fmt.Errorf("error mapping artifact: %w", err)
	}

	return mappedArtifact, nil
}

func (r *ArtifactRepositoryImpl) List(listOptions models.ArtifactListOptions) (*models.ListWrapper[models.Artifact], error) {
	list := models.ListWrapper[models.Artifact]{
		PageSize: listOptions.GetPageSize(),
	}

	artifacts := []models.Artifact{}
	artifactsArt := []schema.Artifact{}

	query := r.db.Model(&schema.Artifact{})

	if listOptions.Name != nil {
		query = query.Where("name = ?", listOptions.Name)
	} else if listOptions.ExternalID != nil {
		query = query.Where("external_id = ?", listOptions.ExternalID)
	}

	if listOptions.ParentResourceID != nil {
		query = query.Joins("JOIN Attribution ON Attribution.artifact_id = Artifact.id").
			Where("Attribution.context_id = ?", listOptions.ParentResourceID).
			Select("Artifact.*") // Explicitly select from Artifact table to avoid ambiguity
		// Use table-prefixed pagination to avoid column ambiguity
		query = query.Scopes(scopes.PaginateWithTablePrefix(artifacts, &listOptions.Pagination, r.db, "Artifact"))
	} else {
		query = query.Scopes(scopes.Paginate(artifacts, &listOptions.Pagination, r.db))
	}

	if err := query.Find(&artifactsArt).Error; err != nil {
		return nil, fmt.Errorf("error listing artifacts: %w", err)
	}

	hasMore := false
	pageSize := listOptions.GetPageSize()
	if pageSize > 0 {
		hasMore = len(artifactsArt) > int(pageSize)
		if hasMore {
			artifactsArt = artifactsArt[:len(artifactsArt)-1]
		}
	}

	for _, artifactArt := range artifactsArt {
		properties := []schema.ArtifactProperty{}
		if err := r.db.Where("artifact_id = ?", artifactArt.ID).Find(&properties).Error; err != nil {
			return nil, fmt.Errorf("error getting properties by artifact id: %w", err)
		}

		artifact, err := r.mapDataLayerToArtifact(artifactArt, properties)
		if err != nil {
			return nil, fmt.Errorf("error mapping artifact: %w", err)
		}
		artifacts = append(artifacts, artifact)
	}

	if hasMore && len(artifactsArt) > 0 {
		lastArtifact := artifactsArt[len(artifactsArt)-1]
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
		return models.Artifact{}, fmt.Errorf("invalid artifact type: %d", artifact.TypeID)
	}

	return artToReturn, nil
}

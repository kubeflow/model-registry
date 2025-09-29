package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	mr_models "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/pkg/api"
)

type dbCatalogImpl struct {
	catalogModelRepository           models.CatalogModelRepository
	catalogModelArtifactRepository   models.CatalogModelArtifactRepository
	catalogMetricsArtifactRepository models.CatalogMetricsArtifactRepository
}

func NewDBCatalog(
	catalogModelRepository models.CatalogModelRepository,
	catalogModelArtifactRepository models.CatalogModelArtifactRepository,
	catalogMetricsArtifactRepository models.CatalogMetricsArtifactRepository,
) CatalogSourceProvider {
	return &dbCatalogImpl{
		catalogModelRepository:           catalogModelRepository,
		catalogModelArtifactRepository:   catalogModelArtifactRepository,
		catalogMetricsArtifactRepository: catalogMetricsArtifactRepository,
	}
}

func (d *dbCatalogImpl) GetModel(ctx context.Context, modelName string, sourceID string) (*model.CatalogModel, error) {
	modelsList, err := d.catalogModelRepository.List(models.CatalogModelListOptions{
		Name:      &modelName,
		SourceIDs: &[]string{sourceID},
	})
	if err != nil {
		return nil, err
	}

	if len(modelsList.Items) == 0 {
		return nil, fmt.Errorf("no models found for name=%v: %w", modelName, api.ErrNotFound)
	}

	if len(modelsList.Items) > 1 {
		return nil, fmt.Errorf("multiple models found for name=%v: %w", modelName, api.ErrNotFound)
	}

	model := mapCatalogModelToCatalogModel(modelsList.Items[0])

	return &model, nil
}

func (d *dbCatalogImpl) ListModels(ctx context.Context, params ListModelsParams) (model.CatalogModelList, error) {
	pageSize := int32(params.PageSize)
	orderBy := string(params.OrderBy)
	sortOrder := string(params.SortOrder)
	nextPageToken := params.NextPageToken

	modelsList, err := d.catalogModelRepository.List(models.CatalogModelListOptions{
		SourceIDs: &params.SourceIDs,
		Pagination: mr_models.Pagination{
			PageSize:      &pageSize,
			OrderBy:       &orderBy,
			SortOrder:     &sortOrder,
			NextPageToken: nextPageToken,
		},
	})
	if err != nil {
		return model.CatalogModelList{}, err
	}

	modelList := &model.CatalogModelList{
		Items: make([]model.CatalogModel, 0),
	}

	for _, model := range modelsList.Items {
		modelList.Items = append(modelList.Items, mapCatalogModelToCatalogModel(model))
	}

	modelList.NextPageToken = *nextPageToken
	modelList.PageSize = pageSize
	modelList.Size = int32(len(modelsList.Items))

	return *modelList, nil
}

func (d *dbCatalogImpl) GetArtifacts(ctx context.Context, modelName string, sourceID string, params ListArtifactsParams) (model.CatalogArtifactList, error) {
	pageSize := int32(params.PageSize)
	orderBy := string(params.OrderBy)
	sortOrder := string(params.SortOrder)
	nextPageToken := params.NextPageToken

	m, err := d.GetModel(ctx, modelName, sourceID)
	if err != nil {
		return model.CatalogArtifactList{}, err
	}

	parentResourceID, err := strconv.ParseInt(*m.Id, 10, 32)
	if err != nil {
		return model.CatalogArtifactList{}, err
	}

	parentResourceID32 := int32(parentResourceID)

	artifactsList, err := d.catalogModelArtifactRepository.List(models.CatalogModelArtifactListOptions{
		ParentResourceID: &parentResourceID32,
		Pagination: mr_models.Pagination{
			PageSize:      &pageSize,
			OrderBy:       &orderBy,
			SortOrder:     &sortOrder,
			NextPageToken: nextPageToken,
		},
	})
	if err != nil {
		return model.CatalogArtifactList{}, err
	}

	artifactList := &model.CatalogArtifactList{
		Items: make([]model.CatalogArtifact, 0),
	}

	for _, artifact := range artifactsList.Items {
		mappedArtifact := mapCatalogModelArtifactToCatalogArtifact(artifact)
		artifactList.Items = append(artifactList.Items, mappedArtifact)
	}

	artifactList.NextPageToken = *nextPageToken
	artifactList.PageSize = pageSize
	artifactList.Size = int32(len(artifactsList.Items))

	return *artifactList, nil
}

func mapCatalogModelToCatalogModel(m models.CatalogModel) model.CatalogModel {
	res := model.CatalogModel{}

	id := strconv.FormatInt(int64(*m.GetID()), 10)
	res.Id = &id

	if m.GetAttributes() != nil {
		res.Name = *m.GetAttributes().Name
		res.ExternalId = m.GetAttributes().ExternalID

		if m.GetAttributes().CreateTimeSinceEpoch != nil {
			createTimeSinceEpoch := strconv.FormatInt(*m.GetAttributes().CreateTimeSinceEpoch, 10)
			res.CreateTimeSinceEpoch = &createTimeSinceEpoch
		}
		if m.GetAttributes().LastUpdateTimeSinceEpoch != nil {
			lastUpdateTimeSinceEpoch := strconv.FormatInt(*m.GetAttributes().LastUpdateTimeSinceEpoch, 10)
			res.LastUpdateTimeSinceEpoch = &lastUpdateTimeSinceEpoch
		}
	}

	if m.GetProperties() != nil {
		for _, prop := range *m.GetProperties() {
			switch prop.Name {
			case "source_id":
				if prop.StringValue != nil {
					res.SourceId = prop.StringValue
				}
			case "description":
				if prop.StringValue != nil {
					res.Description = prop.StringValue
				}
			case "library_name":
				if prop.StringValue != nil {
					res.LibraryName = prop.StringValue
				}
			case "license_link":
				if prop.StringValue != nil {
					res.LicenseLink = prop.StringValue
				}
			case "license":
				if prop.StringValue != nil {
					res.License = prop.StringValue
				}
			case "logo":
				if prop.StringValue != nil {
					res.Logo = prop.StringValue
				}
			case "maturity":
				if prop.StringValue != nil {
					res.Maturity = prop.StringValue
				}
			case "provider":
				if prop.StringValue != nil {
					res.Provider = prop.StringValue
				}
			case "readme":
				if prop.StringValue != nil {
					res.Readme = prop.StringValue
				}
			case "language":
				if prop.StringValue != nil {
					var languages []string
					if err := json.Unmarshal([]byte(*prop.StringValue), &languages); err == nil {
						res.Language = languages
					}
				}
			case "tasks":
				if prop.StringValue != nil {
					var tasks []string
					if err := json.Unmarshal([]byte(*prop.StringValue), &tasks); err == nil {
						res.Tasks = tasks
					}
				}
			}
		}
	}

	return res
}

func mapCatalogModelArtifactToCatalogArtifact(a models.CatalogModelArtifact) model.CatalogArtifact {
	// Create CatalogModelArtifact
	catalogModelArtifact := &model.CatalogModelArtifact{
		ArtifactType: "model-artifact", // Default artifact type
		Uri:          "",               // Will be set from attributes below
	}

	// Map basic fields
	if a.GetID() != nil {
		id := strconv.FormatInt(int64(*a.GetID()), 10)
		catalogModelArtifact.Id = &id
	}

	if a.GetAttributes() != nil {
		attrs := a.GetAttributes()

		catalogModelArtifact.Name = attrs.Name
		catalogModelArtifact.ExternalId = attrs.ExternalID

		if attrs.URI != nil {
			catalogModelArtifact.Uri = *attrs.URI
		}

		if attrs.CreateTimeSinceEpoch != nil {
			createTime := strconv.FormatInt(*attrs.CreateTimeSinceEpoch, 10)
			catalogModelArtifact.CreateTimeSinceEpoch = &createTime
		}

		if attrs.LastUpdateTimeSinceEpoch != nil {
			updateTime := strconv.FormatInt(*attrs.LastUpdateTimeSinceEpoch, 10)
			catalogModelArtifact.LastUpdateTimeSinceEpoch = &updateTime
		}
	}

	return model.CatalogArtifact{
		CatalogModelArtifact: catalogModelArtifact,
	}
}

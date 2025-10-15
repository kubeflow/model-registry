package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	dbmodels "github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	apimodels "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/internal/converter"
	mrmodels "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

type dbCatalogImpl struct {
	catalogModelRepository    dbmodels.CatalogModelRepository
	catalogArtifactRepository dbmodels.CatalogArtifactRepository
	sources                   *SourceCollection
}

func NewDBCatalog(services service.Services, sources *SourceCollection) APIProvider {
	return &dbCatalogImpl{
		catalogArtifactRepository: services.CatalogArtifactRepository,
		catalogModelRepository:    services.CatalogModelRepository,
		sources:                   sources,
	}
}

func (d *dbCatalogImpl) GetModel(ctx context.Context, modelName string, sourceID string) (*apimodels.CatalogModel, error) {
	modelsList, err := d.catalogModelRepository.List(dbmodels.CatalogModelListOptions{
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

	model := mapDBModelToAPIModel(modelsList.Items[0])

	return &model, nil
}

func (d *dbCatalogImpl) ListModels(ctx context.Context, params ListModelsParams) (apimodels.CatalogModelList, error) {
	pageSize := int32(params.PageSize)

	// Use consistent defaults to match pagination logic
	orderBy := string(params.OrderBy)
	if orderBy == "" {
		orderBy = mrmodels.DefaultOrderBy
	}

	sortOrder := string(params.SortOrder)
	if sortOrder == "" {
		sortOrder = mrmodels.DefaultSortOrder
	}

	nextPageToken := params.NextPageToken

	var queryPtr *string
	if params.Query != "" {
		queryPtr = &params.Query
	}

	sourceIDs := params.SourceIDs
	if len(sourceIDs) == 0 && len(params.SourceLabels) > 0 {
		sources := d.sources.ByLabel(params.SourceLabels)
		if len(sources) == 0 {
			// No matching sources, so no matching models.
			return apimodels.CatalogModelList{
				Items:    make([]apimodels.CatalogModel, 0),
				PageSize: pageSize,
			}, nil
		}

		sourceIDs = make([]string, len(sources))
		for i, source := range sources {
			sourceIDs[i] = source.Id
		}
	}

	modelsList, err := d.catalogModelRepository.List(dbmodels.CatalogModelListOptions{
		SourceIDs: &sourceIDs,
		Query:     queryPtr,
		Pagination: mrmodels.Pagination{
			FilterQuery:   &params.FilterQuery,
			PageSize:      &pageSize,
			OrderBy:       &orderBy,
			SortOrder:     &sortOrder,
			NextPageToken: nextPageToken,
		},
	})
	if err != nil {
		return apimodels.CatalogModelList{}, err
	}

	modelList := &apimodels.CatalogModelList{
		Items: make([]apimodels.CatalogModel, 0, len(modelsList.Items)),
	}

	for _, model := range modelsList.Items {
		modelList.Items = append(modelList.Items, mapDBModelToAPIModel(model))
	}

	modelList.NextPageToken = modelsList.NextPageToken
	modelList.PageSize = pageSize
	modelList.Size = int32(len(modelsList.Items))

	return *modelList, nil
}

func (d *dbCatalogImpl) GetArtifacts(ctx context.Context, modelName string, sourceID string, params ListArtifactsParams) (apimodels.CatalogArtifactList, error) {
	pageSize := int32(params.PageSize)

	// Use consistent defaults to match pagination logic
	orderBy := string(params.OrderBy)
	if orderBy == "" {
		orderBy = mrmodels.DefaultOrderBy
	}

	sortOrder := string(params.SortOrder)
	if sortOrder == "" {
		sortOrder = mrmodels.DefaultSortOrder
	}

	nextPageToken := params.NextPageToken

	m, err := d.GetModel(ctx, modelName, sourceID)
	if err != nil {
		if errors.Is(err, api.ErrNotFound) {
			return apimodels.CatalogArtifactList{}, fmt.Errorf("invalid model name '%s' for source '%s': %w", modelName, sourceID, api.ErrBadRequest)
		}
		return apimodels.CatalogArtifactList{}, err
	}

	parentResourceID, err := strconv.ParseInt(*m.Id, 10, 32)
	if err != nil {
		return apimodels.CatalogArtifactList{}, err
	}

	parentResourceID32 := int32(parentResourceID)

	artifactsList, err := d.catalogArtifactRepository.List(dbmodels.CatalogArtifactListOptions{
		ParentResourceID: &parentResourceID32,
		ArtifactType:     params.ArtifactType,
		Pagination: mrmodels.Pagination{
			PageSize:      &pageSize,
			OrderBy:       &orderBy,
			SortOrder:     &sortOrder,
			NextPageToken: nextPageToken,
		},
	})
	if err != nil {
		return apimodels.CatalogArtifactList{}, err
	}

	artifactList := &apimodels.CatalogArtifactList{
		Items: make([]apimodels.CatalogArtifact, 0),
	}

	for _, artifact := range artifactsList.Items {
		mappedArtifact, err := mapDBArtifactToAPIArtifact(artifact)
		if err != nil {
			return apimodels.CatalogArtifactList{}, err
		}
		artifactList.Items = append(artifactList.Items, mappedArtifact)
	}

	artifactList.NextPageToken = artifactsList.NextPageToken
	artifactList.PageSize = pageSize
	artifactList.Size = int32(len(artifactList.Items))

	return *artifactList, nil
}

func (d *dbCatalogImpl) GetFilterOptions(ctx context.Context) (*apimodels.FilterOptionsList, error) {
	// Max length threshold for filter values (excludes verbose fields like readme, description)
	const maxFilterValueLength = 100

	// Query database for filterable properties
	filterableProps, err := d.catalogModelRepository.GetFilterableProperties(maxFilterValueLength)
	if err != nil {
		return nil, err
	}

	// Build FilterOptionsList
	options := make(map[string]apimodels.FilterOption)

	// Process each property and its values
	for fieldName, values := range filterableProps {
		// Skip internal/technical fields that shouldn't be exposed as filters
		if fieldName == "source_id" || fieldName == "logo" || fieldName == "license_link" {
			continue
		}

		// Deduplicate values
		uniqueValues := make(map[string]bool)

		// Parse JSON arrays for fields like language and tasks
		for _, value := range values {
			var arrayValues []string
			if err := json.Unmarshal([]byte(value), &arrayValues); err == nil {
				// Successfully parsed as array, add individual values
				for _, v := range arrayValues {
					uniqueValues[v] = true
				}
			} else {
				// Not a JSON array
				uniqueValues[value] = true
			}
		}

		if len(uniqueValues) > 0 {
			expandedValues := make([]interface{}, 0, len(uniqueValues))
			for v := range uniqueValues {
				expandedValues = append(expandedValues, v)
			}

			options[fieldName] = apimodels.FilterOption{
				Type:   "string",
				Values: expandedValues,
			}
		}
	}

	return &apimodels.FilterOptionsList{
		Filters: &options,
	}, nil
}

func mapDBModelToAPIModel(m dbmodels.CatalogModel) apimodels.CatalogModel {
	res := apimodels.CatalogModel{}

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

	// Map custom properties
	if m.GetCustomProperties() != nil && len(*m.GetCustomProperties()) > 0 {
		customProps := make(map[string]apimodels.MetadataValue, len(*m.GetCustomProperties()))
		for _, prop := range *m.GetCustomProperties() {
			if prop.StringValue != nil {
				customProps[prop.Name] = apimodels.MetadataStringValueAsMetadataValue(
					apimodels.NewMetadataStringValue(*prop.StringValue, "MetadataStringValue"),
				)
			}
		}
		if len(customProps) > 0 {
			res.CustomProperties = &customProps
		}
	}

	return res
}

func mapDBArtifactToAPIArtifact(a dbmodels.CatalogArtifact) (apimodels.CatalogArtifact, error) {
	if a.CatalogModelArtifact != nil {
		return mapToModelArtifact(a.CatalogModelArtifact)
	} else if a.CatalogMetricsArtifact != nil {
		metricsTypeValue := string(a.CatalogMetricsArtifact.GetAttributes().MetricsType)
		return mapToMetricsArtifact(a.CatalogMetricsArtifact, metricsTypeValue)
	}

	return apimodels.CatalogArtifact{}, fmt.Errorf("invalid catalog artifact type: %v", a)
}

func mapToModelArtifact(a dbmodels.CatalogModelArtifact) (apimodels.CatalogArtifact, error) {
	catalogModelArtifact := &apimodels.CatalogModelArtifact{
		ArtifactType: dbmodels.CatalogModelArtifactType,
	}

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

	if a.GetProperties() != nil {
		for _, prop := range *a.GetProperties() {
			switch prop.Name {
			case "description":
				if prop.StringValue != nil {
					catalogModelArtifact.Description = prop.StringValue
				}
			case "artifactType":
				if prop.StringValue != nil {
					catalogModelArtifact.ArtifactType = *prop.StringValue
				}
			}
		}
	}

	// Map custom properties
	if a.GetCustomProperties() != nil && len(*a.GetCustomProperties()) > 0 {
		customPropsMap, err := converter.MapEmbedMDCustomProperties(*a.GetCustomProperties())
		if err != nil {
			return apimodels.CatalogArtifact{}, fmt.Errorf("error mapping custom properties: %w", err)
		}

		catalogCustomProps := convertMetadataValueMap(customPropsMap)
		catalogModelArtifact.CustomProperties = &catalogCustomProps
	}

	return apimodels.CatalogArtifact{
		CatalogModelArtifact: catalogModelArtifact,
	}, nil
}

func mapToMetricsArtifact(a dbmodels.CatalogMetricsArtifact, metricsType string) (apimodels.CatalogArtifact, error) {
	catalogMetricsArtifact := &apimodels.CatalogMetricsArtifact{
		ArtifactType: dbmodels.CatalogMetricsArtifactType,
		MetricsType:  metricsType,
	}

	if a.GetID() != nil {
		id := strconv.FormatInt(int64(*a.GetID()), 10)
		catalogMetricsArtifact.Id = &id
	}

	if a.GetAttributes() != nil {
		attrs := a.GetAttributes()

		catalogMetricsArtifact.Name = attrs.Name
		catalogMetricsArtifact.ExternalId = attrs.ExternalID

		if attrs.CreateTimeSinceEpoch != nil {
			createTime := strconv.FormatInt(*attrs.CreateTimeSinceEpoch, 10)
			catalogMetricsArtifact.CreateTimeSinceEpoch = &createTime
		}

		if attrs.LastUpdateTimeSinceEpoch != nil {
			updateTime := strconv.FormatInt(*attrs.LastUpdateTimeSinceEpoch, 10)
			catalogMetricsArtifact.LastUpdateTimeSinceEpoch = &updateTime
		}
	}

	if a.GetProperties() != nil {
		for _, prop := range *a.GetProperties() {
			switch prop.Name {
			case "description":
				if prop.StringValue != nil {
					catalogMetricsArtifact.Description = prop.StringValue
				}
			}
		}
	}

	// Map custom properties
	if a.GetCustomProperties() != nil && len(*a.GetCustomProperties()) > 0 {
		customPropsMap, err := converter.MapEmbedMDCustomProperties(*a.GetCustomProperties())
		if err != nil {
			return apimodels.CatalogArtifact{}, fmt.Errorf("error mapping custom properties: %w", err)
		}

		catalogCustomProps := convertMetadataValueMap(customPropsMap)
		catalogMetricsArtifact.CustomProperties = &catalogCustomProps

	}

	return apimodels.CatalogArtifact{
		CatalogMetricsArtifact: catalogMetricsArtifact,
	}, nil
}

// convertMetadataValueMap converts from pkg/openapi.MetadataValue to catalog/pkg/openapi.MetadataValue
func convertMetadataValueMap(source map[string]openapi.MetadataValue) map[string]apimodels.MetadataValue {
	result := make(map[string]apimodels.MetadataValue)

	for key, value := range source {
		catalogValue := apimodels.MetadataValue{}

		if value.MetadataStringValue != nil {
			catalogValue.MetadataStringValue = &apimodels.MetadataStringValue{
				StringValue:  value.MetadataStringValue.StringValue,
				MetadataType: value.MetadataStringValue.MetadataType,
			}
		} else if value.MetadataIntValue != nil {
			catalogValue.MetadataIntValue = &apimodels.MetadataIntValue{
				IntValue:     value.MetadataIntValue.IntValue,
				MetadataType: value.MetadataIntValue.MetadataType,
			}
		} else if value.MetadataDoubleValue != nil {
			catalogValue.MetadataDoubleValue = &apimodels.MetadataDoubleValue{
				DoubleValue:  value.MetadataDoubleValue.DoubleValue,
				MetadataType: value.MetadataDoubleValue.MetadataType,
			}
		} else if value.MetadataBoolValue != nil {
			catalogValue.MetadataBoolValue = &apimodels.MetadataBoolValue{
				BoolValue:    value.MetadataBoolValue.BoolValue,
				MetadataType: value.MetadataBoolValue.MetadataType,
			}
		} else if value.MetadataStructValue != nil {
			catalogValue.MetadataStructValue = &apimodels.MetadataStructValue{
				StructValue:  value.MetadataStructValue.StructValue,
				MetadataType: value.MetadataStructValue.MetadataType,
			}
		}

		result[key] = catalogValue
	}

	return result
}

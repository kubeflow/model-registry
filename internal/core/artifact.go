package core

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"gorm.io/gorm"
)

func (b *ModelRegistryService) upsertArtifact(artifact *openapi.Artifact, parentResourceId *string) (*openapi.Artifact, error) {
	var parentResourceIDPtr *int32

	artToReturn := &openapi.Artifact{}

	if artifact == nil {
		return nil, fmt.Errorf("invalid artifact pointer, cannot be nil: %w", api.ErrBadRequest)
	}

	// If parentResourceId is nil, we need to generate a new UUID for the parent resource ID for the name prefixing,
	// the id will still be nil when invoking the Save method so that the attribution is created only when needed.
	if parentResourceId == nil {
		uuid := uuid.New().String()
		parentResourceId = &uuid
	} else {
		convertedId, err := strconv.ParseInt(*parentResourceId, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		convertedIdInt32 := int32(convertedId)
		parentResourceIDPtr = &convertedIdInt32
	}

	if ma := artifact.ModelArtifact; ma != nil {
		if ma.Id != nil {
			existing, err := b.getArtifact(*ma.Id, true)
			if err != nil {
				return nil, fmt.Errorf("mismatched types, artifact with id %s is not a model artifact: %w", *ma.Id, api.ErrBadRequest)
			}

			if existing.ModelArtifact == nil {
				return nil, fmt.Errorf("artifact with id %s is not a model artifact: %w", *ma.Id, api.ErrBadRequest)
			}

			withNotEditable, err := b.mapper.UpdateExistingModelArtifact(converter.NewOpenapiUpdateWrapper(existing.ModelArtifact, ma))
			if err != nil {
				return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
			}

			ma = &withNotEditable
		} else {
			name := ""
			if ma.Name != nil {
				name = *ma.Name
			}

			prefixedName := converter.PrefixWhenOwned(parentResourceId, name)
			ma.Name = &prefixedName
		}

		modelArtifact, err := b.mapper.MapFromModelArtifact(ma)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		modelArtifact, err = b.modelArtifactRepository.Save(modelArtifact, parentResourceIDPtr)
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return nil, fmt.Errorf("model artifact with name %s already exists: %w", *ma.Name, api.ErrConflict)
			}
			return nil, err
		}

		toReturn, err := b.mapper.MapToModelArtifact(modelArtifact)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		artToReturn.ModelArtifact = toReturn
		return artToReturn, nil
	} else if da := artifact.DocArtifact; da != nil {
		if da.Id != nil {
			existing, err := b.getArtifact(*da.Id, true)
			if err != nil {
				return nil, fmt.Errorf("mismatched types, artifact with id %s is not a doc artifact: %w", *da.Id, api.ErrBadRequest)
			}

			if existing.DocArtifact == nil {
				return nil, fmt.Errorf("artifact with id %s is not a doc artifact: %w", *da.Id, api.ErrBadRequest)
			}

			withNotEditable, err := b.mapper.UpdateExistingDocArtifact(converter.NewOpenapiUpdateWrapper(existing.DocArtifact, da))
			if err != nil {
				return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
			}

			da = &withNotEditable
		} else {
			name := ""
			if da.Name != nil {
				name = *da.Name
			}

			prefixedName := converter.PrefixWhenOwned(parentResourceId, name)
			da.Name = &prefixedName
		}

		docArtifact, err := b.mapper.MapFromDocArtifact(da)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		docArtifact, err = b.docArtifactRepository.Save(docArtifact, parentResourceIDPtr)
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return nil, fmt.Errorf("doc artifact with name %s already exists: %w", *da.Name, api.ErrConflict)
			}
			return nil, err
		}

		toReturn, err := b.mapper.MapToDocArtifact(docArtifact)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		artToReturn.DocArtifact = toReturn
		return artToReturn, nil
	} else if ds := artifact.DataSet; ds != nil {
		// Handle DataSet artifacts using embedmd converters
		if ds.Id != nil {
			existing, err := b.getArtifact(*ds.Id, true)
			if err != nil {
				return nil, fmt.Errorf("mismatched types, artifact with id %s is not a dataset artifact: %w", *ds.Id, api.ErrBadRequest)
			}

			if existing.DataSet == nil {
				return nil, fmt.Errorf("artifact with id %s is not a dataset artifact: %w", *ds.Id, api.ErrBadRequest)
			}

			withNotEditable, err := b.mapper.UpdateExistingDataSet(converter.NewOpenapiUpdateWrapper(existing.DataSet, ds))
			if err != nil {
				return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
			}

			ds = &withNotEditable
		} else {
			name := ""
			if ds.Name != nil {
				name = *ds.Name
			}

			prefixedName := converter.PrefixWhenOwned(parentResourceId, name)
			ds.Name = &prefixedName
		}

		dataSetEntity, err := b.mapper.MapFromDataSet(ds)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		dataSetEntity, err = b.dataSetRepository.Save(dataSetEntity, parentResourceIDPtr)
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return nil, fmt.Errorf("dataset artifact with name %s already exists: %w", *ds.Name, api.ErrConflict)
			}
			return nil, err
		}

		toReturn, err := b.mapper.MapToDataSet(dataSetEntity)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		artToReturn.DataSet = toReturn
		return artToReturn, nil
	} else if me := artifact.Metric; me != nil {
		// Handle metric artifacts using embedmd converters
		if me.Id != nil {
			existing, err := b.getArtifact(*me.Id, true)
			if err != nil {
				return nil, fmt.Errorf("mismatched types, artifact with id %s is not a metric artifact: %w", *me.Id, api.ErrBadRequest)
			}

			if existing.Metric == nil {
				return nil, fmt.Errorf("artifact with id %s is not a metric artifact: %w", *me.Id, api.ErrBadRequest)
			}

			withNotEditable, err := b.mapper.UpdateExistingMetric(converter.NewOpenapiUpdateWrapper(existing.Metric, me))
			if err != nil {
				return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
			}

			me = &withNotEditable
		} else {
			name := ""
			if me.Name != nil {
				name = *me.Name
			}

			// Only prefix the name if we're in a parent resource context (model version or experiment run)
			if parentResourceId != nil && parentResourceIDPtr != nil {
				prefixedName := converter.PrefixWhenOwned(parentResourceId, name)
				me.Name = &prefixedName
			}
		}

		metricEntity, err := b.mapper.MapFromMetric(me)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		metricEntity, err = b.metricRepository.Save(metricEntity, parentResourceIDPtr)
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return nil, fmt.Errorf("metric artifact with name %s already exists: %w", *me.Name, api.ErrConflict)
			}
			return nil, err
		}

		toReturn, err := b.mapper.MapToMetricFromMetric(metricEntity)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		artToReturn.Metric = toReturn
		return artToReturn, nil
	} else if pa := artifact.Parameter; pa != nil {
		// ADD PARAMETER SUPPORT using embedmd converters
		if pa.Id != nil {
			existing, err := b.getArtifact(*pa.Id, true)
			if err != nil {
				return nil, fmt.Errorf("mismatched types, artifact with id %s is not a parameter artifact: %w", *pa.Id, api.ErrBadRequest)
			}

			if existing.Parameter == nil {
				return nil, fmt.Errorf("artifact with id %s is not a parameter artifact: %w", *pa.Id, api.ErrBadRequest)
			}

			withNotEditable, err := b.mapper.UpdateExistingParameter(converter.NewOpenapiUpdateWrapper(existing.Parameter, pa))
			if err != nil {
				return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
			}

			pa = &withNotEditable
		} else {
			name := ""
			if pa.Name != nil {
				name = *pa.Name
			}

			// Only prefix the name if we're in a parent resource context (model version or experiment run)
			if parentResourceId != nil && parentResourceIDPtr != nil {
				prefixedName := converter.PrefixWhenOwned(parentResourceId, name)
				pa.Name = &prefixedName
			}
		}

		parameterEntity, err := b.mapper.MapFromParameter(pa)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		parameterEntity, err = b.parameterRepository.Save(parameterEntity, parentResourceIDPtr)
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return nil, fmt.Errorf("parameter artifact with name %s already exists: %w", *pa.Name, api.ErrConflict)
			}
			return nil, err
		}

		toReturn, err := b.mapper.MapToParameter(parameterEntity)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		artToReturn.Parameter = toReturn
		return artToReturn, nil
	}

	return nil, fmt.Errorf("invalid artifact type, must be either ModelArtifact, DocArtifact, DataSet, Metric, or Parameter: %w", api.ErrBadRequest)
}

func (b *ModelRegistryService) UpsertModelVersionArtifact(artifact *openapi.Artifact, parentResourceId string) (*openapi.Artifact, error) {
	return b.upsertArtifact(artifact, &parentResourceId)
}

func (b *ModelRegistryService) UpsertArtifact(artifact *openapi.Artifact) (*openapi.Artifact, error) {
	return b.upsertArtifact(artifact, nil)
}

func (b *ModelRegistryService) getArtifact(id string, preserveName bool) (*openapi.Artifact, error) {
	artToReturn := &openapi.Artifact{}
	convertedId, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	artifact, err := b.artifactRepository.GetByID(int32(convertedId))
	if err != nil {
		return nil, fmt.Errorf("no artifact found for id %s: %w", id, api.ErrNotFound)
	}

	if artifact.ModelArtifact != nil {
		toReturn, err := b.mapper.MapToModelArtifact(*artifact.ModelArtifact)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		artToReturn.ModelArtifact = toReturn
	} else if artifact.DocArtifact != nil {
		toReturn, err := b.mapper.MapToDocArtifact(*artifact.DocArtifact)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		artToReturn.DocArtifact = toReturn
	} else if artifact.DataSet != nil {
		toReturn, err := b.mapper.MapToDataSet(*artifact.DataSet)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		artToReturn.DataSet = toReturn
	} else if artifact.Metric != nil {
		toReturn, err := b.mapper.MapToMetricFromMetric(*artifact.Metric)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		artToReturn.Metric = toReturn
	} else if artifact.Parameter != nil {
		toReturn, err := b.mapper.MapToParameter(*artifact.Parameter)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		artToReturn.Parameter = toReturn
	}

	if preserveName {
		if artifact.ModelArtifact != nil {
			artToReturn.ModelArtifact.Name = (*artifact.ModelArtifact).GetAttributes().Name
		}
		if artifact.DocArtifact != nil {
			artToReturn.DocArtifact.Name = (*artifact.DocArtifact).GetAttributes().Name
		}
		if artifact.DataSet != nil {
			artToReturn.DataSet.Name = (*artifact.DataSet).GetAttributes().Name
		}
		if artifact.Metric != nil {
			artToReturn.Metric.Name = (*artifact.Metric).GetAttributes().Name
		}
		if artifact.Parameter != nil {
			artToReturn.Parameter.Name = (*artifact.Parameter).GetAttributes().Name
		}
	}

	return artToReturn, nil
}

func (b *ModelRegistryService) GetArtifactById(id string) (*openapi.Artifact, error) {
	return b.getArtifact(id, false)
}

func (b *ModelRegistryService) getArtifactsByParams(artifactName *string, parentResourceId *string, externalId *string, artifactType string) (*openapi.Artifact, error) {
	artToReturn := &openapi.Artifact{}

	if artifactName == nil && parentResourceId == nil && externalId == nil {
		return nil, fmt.Errorf("invalid parameters call, supply either (artifactName and parentResourceId), or externalId: %w", api.ErrBadRequest)
	}

	if artifactName != nil && parentResourceId != nil {
		combinedName := converter.PrefixWhenOwned(parentResourceId, *artifactName)
		artifactName = &combinedName
	}

	artifacts, err := b.artifactRepository.List(models.ArtifactListOptions{
		Name:       artifactName,
		ExternalID: externalId,
	})
	if err != nil {
		return nil, err
	}

	if artifactType != "" {
		artifactType += " "
	}

	if len(artifacts.Items) == 0 {
		return nil, fmt.Errorf("no %sartifacts found for name=%v, parentResourceId=%v, externalId=%v: %w", artifactType, apiutils.ZeroIfNil(artifactName), apiutils.ZeroIfNil(parentResourceId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	if len(artifacts.Items) > 1 {
		return nil, fmt.Errorf("multiple %sartifacts found for name=%v, parentResourceId=%v, externalId=%v: %w", artifactType, apiutils.ZeroIfNil(artifactName), apiutils.ZeroIfNil(parentResourceId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	if artifacts.Items[0].ModelArtifact != nil {
		modelArtifact, err := b.mapper.MapToModelArtifact(*artifacts.Items[0].ModelArtifact)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		artToReturn.ModelArtifact = modelArtifact
	} else {
		docArtifact, err := b.mapper.MapToDocArtifact(*artifacts.Items[0].DocArtifact)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		artToReturn.DocArtifact = docArtifact
	}

	return artToReturn, nil
}

func (b *ModelRegistryService) GetArtifactByParams(artifactName *string, parentResourceId *string, externalId *string) (*openapi.Artifact, error) {
	return b.getArtifactsByParams(artifactName, parentResourceId, externalId, "")
}

func (b *ModelRegistryService) GetArtifacts(artifactType openapi.ArtifactTypeQueryParam, listOptions api.ListOptions, parentResourceId *string) (*openapi.ArtifactList, error) {
	var parentResourceIDPtr *int32

	if parentResourceId != nil {
		convertedId, err := strconv.ParseInt(*parentResourceId, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		convertedIdInt32 := int32(convertedId)

		parentResourceIDPtr = &convertedIdInt32
	}

	artifacts, err := b.artifactRepository.List(models.ArtifactListOptions{
		Pagination: models.Pagination{
			PageSize:      listOptions.PageSize,
			OrderBy:       listOptions.OrderBy,
			SortOrder:     listOptions.SortOrder,
			NextPageToken: listOptions.NextPageToken,
		},
		ParentResourceID: parentResourceIDPtr,
	})
	if err != nil {
		return nil, err
	}

	artifactsList := &openapi.ArtifactList{
		Items: []openapi.Artifact{},
	}

	for _, artifact := range artifacts.Items {
		if artifact.ModelArtifact != nil {
			modelArtifact, err := b.mapper.MapToModelArtifact(*artifact.ModelArtifact)
			if err != nil {
				return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
			}
			artifactsList.Items = append(artifactsList.Items, openapi.Artifact{ModelArtifact: modelArtifact})
		} else {
			docArtifact, err := b.mapper.MapToDocArtifact(*artifact.DocArtifact)
			if err != nil {
				return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
			}
			artifactsList.Items = append(artifactsList.Items, openapi.Artifact{DocArtifact: docArtifact})
		}
	}

	artifactsList.NextPageToken = artifacts.NextPageToken
	artifactsList.PageSize = artifacts.PageSize
	artifactsList.Size = int32(artifacts.Size)

	return artifactsList, nil
}

func (b *ModelRegistryService) UpsertModelArtifact(modelArtifact *openapi.ModelArtifact) (*openapi.ModelArtifact, error) {
	if modelArtifact == nil {
		return nil, fmt.Errorf("invalid model artifact pointer, can't upsert nil: %w", api.ErrBadRequest)
	}

	art, err := b.UpsertArtifact(&openapi.Artifact{
		ModelArtifact: modelArtifact,
	})
	if err != nil {
		return nil, err
	}

	return art.ModelArtifact, nil
}

func (b *ModelRegistryService) GetModelArtifactById(id string) (*openapi.ModelArtifact, error) {
	art, err := b.GetArtifactById(id)
	if err != nil {
		return nil, err
	}

	if art.ModelArtifact == nil {
		return nil, fmt.Errorf("artifact with id %s is not a model artifact: %w", id, api.ErrNotFound)
	}

	return art.ModelArtifact, nil
}

func (b *ModelRegistryService) GetModelArtifactByInferenceService(inferenceServiceId string) (*openapi.ModelArtifact, error) {
	mv, err := b.GetModelVersionByInferenceService(inferenceServiceId)
	if err != nil {
		return nil, err
	}

	artifactList, err := b.GetArtifacts(openapi.ARTIFACTTYPEQUERYPARAM_MODEL_ARTIFACT, api.ListOptions{}, mv.Id)
	if err != nil {
		return nil, err
	}

	if len(artifactList.Items) == 0 {
		return nil, fmt.Errorf("no artifacts found for model version %s: %w", *mv.Id, api.ErrNotFound)
	}

	if artifactList.Items[0].ModelArtifact == nil {
		return nil, fmt.Errorf("no artifacts found for model version %s: %w", *mv.Id, api.ErrNotFound)
	}

	return artifactList.Items[0].ModelArtifact, nil
}

func (b *ModelRegistryService) GetModelArtifactByParams(artifactName *string, parentResourceId *string, externalId *string) (*openapi.ModelArtifact, error) {
	art, err := b.getArtifactsByParams(artifactName, parentResourceId, externalId, "model")
	if err != nil {
		return nil, err
	}

	if art.ModelArtifact == nil {
		return nil, fmt.Errorf("artifact with name=%v, parentResourceId=%v, externalId=%v is not a model artifact: %w", apiutils.ZeroIfNil(artifactName), apiutils.ZeroIfNil(parentResourceId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	return art.ModelArtifact, nil
}

func (b *ModelRegistryService) GetModelArtifacts(listOptions api.ListOptions, parentResourceId *string) (*openapi.ModelArtifactList, error) {
	var parentResourceIDPtr *int32

	if parentResourceId != nil {
		convertedId, err := strconv.ParseInt(*parentResourceId, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		convertedIdInt32 := int32(convertedId)

		parentResourceIDPtr = &convertedIdInt32
	}

	modelArtifacts, err := b.modelArtifactRepository.List(models.ModelArtifactListOptions{
		Pagination: models.Pagination{
			PageSize:      listOptions.PageSize,
			OrderBy:       listOptions.OrderBy,
			SortOrder:     listOptions.SortOrder,
			NextPageToken: listOptions.NextPageToken,
		},
		ParentResourceID: parentResourceIDPtr,
	})
	if err != nil {
		return nil, err
	}

	modelArtifactList := &openapi.ModelArtifactList{
		Items: []openapi.ModelArtifact{},
	}

	for _, modelArtifact := range modelArtifacts.Items {
		modelArtifact, err := b.mapper.MapToModelArtifact(modelArtifact)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		modelArtifactList.Items = append(modelArtifactList.Items, *modelArtifact)
	}

	modelArtifactList.NextPageToken = modelArtifacts.NextPageToken
	modelArtifactList.PageSize = modelArtifacts.PageSize
	modelArtifactList.Size = int32(modelArtifacts.Size)

	return modelArtifactList, nil
}

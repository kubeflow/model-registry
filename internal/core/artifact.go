package core

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/mapper"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"gorm.io/gorm"
)

type ModelRegistryService struct {
	artifactRepository           models.ArtifactRepository
	modelArtifactRepository      models.ModelArtifactRepository
	docArtifactRepository        models.DocArtifactRepository
	registeredModelRepository    models.RegisteredModelRepository
	modelVersionRepository       models.ModelVersionRepository
	servingEnvironmentRepository models.ServingEnvironmentRepository
	inferenceServiceRepository   models.InferenceServiceRepository
	serveModelRepository         models.ServeModelRepository
	mapper                       mapper.EmbedMDMapper
}

func NewModelRegistryService(
	artifactRepository models.ArtifactRepository,
	modelArtifactRepository models.ModelArtifactRepository,
	docArtifactRepository models.DocArtifactRepository,
	registeredModelRepository models.RegisteredModelRepository,
	modelVersionRepository models.ModelVersionRepository,
	servingEnvironmentRepository models.ServingEnvironmentRepository,
	inferenceServiceRepository models.InferenceServiceRepository,
	serveModelRepository models.ServeModelRepository,
	typesMap map[string]int64) *ModelRegistryService {
	return &ModelRegistryService{
		artifactRepository:           artifactRepository,
		modelArtifactRepository:      modelArtifactRepository,
		docArtifactRepository:        docArtifactRepository,
		registeredModelRepository:    registeredModelRepository,
		modelVersionRepository:       modelVersionRepository,
		servingEnvironmentRepository: servingEnvironmentRepository,
		inferenceServiceRepository:   inferenceServiceRepository,
		serveModelRepository:         serveModelRepository,
		mapper:                       *mapper.NewEmbedMDMapper(typesMap),
	}
}

func (b *ModelRegistryService) upsertArtifact(artifact *openapi.Artifact, modelVersionId *string) (*openapi.Artifact, error) {
	var modelVersionIDPtr *int32

	artToReturn := &openapi.Artifact{}

	if artifact == nil {
		return nil, fmt.Errorf("invalid artifact pointer, cannot be nil: %w", api.ErrBadRequest)
	}

	// If modelVersionId is nil, we need to generate a new UUID for the model version ID for the name prefixing,
	// the id will still be nil when invoking the Save method so that the attribution is created only when needed.
	if modelVersionId == nil {
		uuid := uuid.New().String()

		modelVersionId = &uuid
	} else {
		convertedId, err := strconv.ParseInt(*modelVersionId, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		convertedIdInt32 := int32(convertedId)

		modelVersionIDPtr = &convertedIdInt32
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

			withNotEditable, err := b.mapper.OverrideNotEditableForModelArtifact(converter.NewOpenapiUpdateWrapper(existing.ModelArtifact, ma))
			if err != nil {
				return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
			}

			ma = &withNotEditable
		} else {
			name := ""

			if ma.Name != nil {
				name = *ma.Name
			}

			prefixedName := converter.PrefixWhenOwned(modelVersionId, name)
			ma.Name = &prefixedName
		}

		modelArtifact, err := b.mapper.MapFromModelArtifact(ma)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		modelArtifact, err = b.modelArtifactRepository.Save(modelArtifact, modelVersionIDPtr)
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

			withNotEditable, err := b.mapper.OverrideNotEditableForDocArtifact(converter.NewOpenapiUpdateWrapper(existing.DocArtifact, da))
			if err != nil {
				return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
			}

			da = &withNotEditable
		} else {
			name := ""

			if da.Name != nil {
				name = *da.Name
			}

			prefixedName := converter.PrefixWhenOwned(modelVersionId, name)
			da.Name = &prefixedName
		}

		docArtifact, err := b.mapper.MapFromDocArtifact(da)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		docArtifact, err = b.docArtifactRepository.Save(docArtifact, modelVersionIDPtr)
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
	}

	return nil, fmt.Errorf("invalid artifact type, must be either ModelArtifact or DocArtifact: %w", api.ErrBadRequest)
}

func (b *ModelRegistryService) UpsertModelVersionArtifact(artifact *openapi.Artifact, modelVersionId string) (*openapi.Artifact, error) {
	return b.upsertArtifact(artifact, &modelVersionId)
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
	} else {
		toReturn, err := b.mapper.MapToDocArtifact(*artifact.DocArtifact)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		artToReturn.DocArtifact = toReturn
	}

	if preserveName {
		if artifact.ModelArtifact != nil {
			artToReturn.ModelArtifact.Name = (*artifact.ModelArtifact).GetAttributes().Name
		}
		if artifact.DocArtifact != nil {
			artToReturn.DocArtifact.Name = (*artifact.DocArtifact).GetAttributes().Name
		}
	}

	return artToReturn, nil
}

func (b *ModelRegistryService) GetArtifactById(id string) (*openapi.Artifact, error) {
	return b.getArtifact(id, false)
}

func (b *ModelRegistryService) getArtifactsByParams(artifactName *string, modelVersionId *string, externalId *string, artifactType string) (*openapi.Artifact, error) {
	artToReturn := &openapi.Artifact{}

	if artifactName == nil && modelVersionId == nil && externalId == nil {
		return nil, fmt.Errorf("invalid parameters call, supply either (artifactName and modelVersionId), or externalId: %w", api.ErrBadRequest)
	}

	if artifactName != nil && modelVersionId != nil {
		combinedName := converter.PrefixWhenOwned(modelVersionId, *artifactName)
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
		return nil, fmt.Errorf("no %sartifacts found for name=%v, modelVersionId=%v, externalId=%v: %w", artifactType, apiutils.ZeroIfNil(artifactName), apiutils.ZeroIfNil(modelVersionId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	if len(artifacts.Items) > 1 {
		return nil, fmt.Errorf("multiple %sartifacts found for name=%v, modelVersionId=%v, externalId=%v: %w", artifactType, apiutils.ZeroIfNil(artifactName), apiutils.ZeroIfNil(modelVersionId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
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

func (b *ModelRegistryService) GetArtifactByParams(artifactName *string, modelVersionId *string, externalId *string) (*openapi.Artifact, error) {
	return b.getArtifactsByParams(artifactName, modelVersionId, externalId, "")
}

func (b *ModelRegistryService) GetArtifacts(listOptions api.ListOptions, modelVersionId *string) (*openapi.ArtifactList, error) {
	var modelVersionIDPtr *int32

	if modelVersionId != nil {
		convertedId, err := strconv.ParseInt(*modelVersionId, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		convertedIdInt32 := int32(convertedId)

		modelVersionIDPtr = &convertedIdInt32
	}

	artifacts, err := b.artifactRepository.List(models.ArtifactListOptions{
		Pagination: models.Pagination{
			PageSize:      listOptions.PageSize,
			OrderBy:       listOptions.OrderBy,
			SortOrder:     listOptions.SortOrder,
			NextPageToken: listOptions.NextPageToken,
		},
		ModelVersionID: modelVersionIDPtr,
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

	artifactList, err := b.GetArtifacts(api.ListOptions{}, mv.Id)
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

func (b *ModelRegistryService) GetModelArtifactByParams(artifactName *string, modelVersionId *string, externalId *string) (*openapi.ModelArtifact, error) {
	art, err := b.getArtifactsByParams(artifactName, modelVersionId, externalId, "model")
	if err != nil {
		return nil, err
	}

	if art.ModelArtifact == nil {
		return nil, fmt.Errorf("artifact with name=%v, modelVersionId=%v, externalId=%v is not a model artifact: %w", apiutils.ZeroIfNil(artifactName), apiutils.ZeroIfNil(modelVersionId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	return art.ModelArtifact, nil
}

func (b *ModelRegistryService) GetModelArtifacts(listOptions api.ListOptions, modelVersionId *string) (*openapi.ModelArtifactList, error) {
	var modelVersionIDPtr *int32

	if modelVersionId != nil {
		convertedId, err := strconv.ParseInt(*modelVersionId, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		convertedIdInt32 := int32(convertedId)

		modelVersionIDPtr = &convertedIdInt32
	}

	modelArtifacts, err := b.modelArtifactRepository.List(models.ModelArtifactListOptions{
		Pagination: models.Pagination{
			PageSize:      listOptions.PageSize,
			OrderBy:       listOptions.OrderBy,
			SortOrder:     listOptions.SortOrder,
			NextPageToken: listOptions.NextPageToken,
		},
		ModelVersionID: modelVersionIDPtr,
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

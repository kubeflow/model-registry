package core

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/mapper"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
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
		convertedId, err := strconv.Atoi(*modelVersionId)
		if err != nil {
			return nil, fmt.Errorf("invalid model version id: %w", err)
		}

		convertedIdInt32 := int32(convertedId)

		modelVersionIDPtr = &convertedIdInt32
	}

	if artifact.ModelArtifact != nil {
		modelArtifact, err := b.mapper.MapFromModelArtifact(artifact.ModelArtifact)
		if err != nil {
			return nil, err
		}

		prefixedName := converter.PrefixWhenOwned(modelVersionId, *modelArtifact.GetAttributes().Name)
		modelArtifact.GetAttributes().Name = &prefixedName

		modelArtifact, err = b.modelArtifactRepository.Save(modelArtifact, modelVersionIDPtr)
		if err != nil {
			return nil, err
		}

		toReturn, err := b.mapper.MapToModelArtifact(modelArtifact)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		artToReturn.ModelArtifact = toReturn

		return artToReturn, nil
	}

	docArtifact, err := b.mapper.MapFromDocArtifact(artifact.DocArtifact)
	if err != nil {
		return nil, err
	}

	prefixedName := converter.PrefixWhenOwned(modelVersionId, *docArtifact.GetAttributes().Name)
	docArtifact.GetAttributes().Name = &prefixedName

	docArtifact, err = b.docArtifactRepository.Save(docArtifact, modelVersionIDPtr)
	if err != nil {
		return nil, err
	}

	toReturn, err := b.mapper.MapToDocArtifact(docArtifact)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	artToReturn.DocArtifact = toReturn

	return artToReturn, nil
}

func (b *ModelRegistryService) UpsertModelVersionArtifact(artifact *openapi.Artifact, modelVersionId string) (*openapi.Artifact, error) {
	return b.upsertArtifact(artifact, &modelVersionId)
}

func (b *ModelRegistryService) UpsertArtifact(artifact *openapi.Artifact) (*openapi.Artifact, error) {
	return b.upsertArtifact(artifact, nil)
}

func (b *ModelRegistryService) GetArtifactById(id string) (*openapi.Artifact, error) {
	artToReturn := &openapi.Artifact{}
	convertedId, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	artifact, err := b.artifactRepository.GetByID(int32(convertedId))
	if err != nil {
		return nil, err
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

	return artToReturn, nil
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
		convertedId, err := strconv.Atoi(*modelVersionId)
		if err != nil {
			return nil, fmt.Errorf("invalid model version id: %w", err)
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
		convertedId, err := strconv.Atoi(*modelVersionId)
		if err != nil {
			return nil, fmt.Errorf("invalid model version id: %w", err)
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

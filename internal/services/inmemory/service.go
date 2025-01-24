package services

import (
	"fmt"
	"os"
	"strconv"

	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"gopkg.in/yaml.v3"
)

type InMemory struct {
	RegisteredModels    map[string]*openapi.RegisteredModel
	ModelVersions       map[string]*openapi.ModelVersion
	Artifacts           map[string]*openapi.Artifact
	ModelArtifacts      map[string]*openapi.ModelArtifact
	ServingEnvironments map[string]*openapi.ServingEnvironment
	InferenceServices   map[string]*openapi.InferenceService
	ServeModels         map[string]*openapi.ServeModel
}

func NewInMemory() *InMemory {
	return &InMemory{
		RegisteredModels:    make(map[string]*openapi.RegisteredModel),
		ModelVersions:       make(map[string]*openapi.ModelVersion),
		Artifacts:           make(map[string]*openapi.Artifact),
		ModelArtifacts:      make(map[string]*openapi.ModelArtifact),
		ServingEnvironments: make(map[string]*openapi.ServingEnvironment),
		InferenceServices:   make(map[string]*openapi.InferenceService),
		ServeModels:         make(map[string]*openapi.ServeModel),
	}
}

func (s *InMemory) Seed(path string) error {
	seedData := InMemorySeedData{}

	seedFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open seed data file: %w", err)
	}

	decoder := yaml.NewDecoder(seedFile)

	// decoder.KnownFields(true)

	err = decoder.Decode(&seedData)
	if err != nil {
		return fmt.Errorf("failed to decode seed data: %w", err)
	}

	for i, rm := range seedData.RegisteredModels {
		if rm.Id == nil {
			newID := strconv.Itoa(i + 1)
			rm.Id = &newID
		}

		s.RegisteredModels[*rm.Id] = rm
	}

	for i, mv := range seedData.ModelVersions {
		if mv.Id == nil {
			newID := strconv.Itoa(i + 1)
			mv.Id = &newID
		}

		s.ModelVersions[*mv.Id] = mv
	}

	for i, a := range seedData.Artifacts {
		if a.DocArtifact != nil {
			if a.DocArtifact.Id == nil {
				newID := strconv.Itoa(i*2 + 1)
				a.DocArtifact.Id = &newID
			}

			s.Artifacts[*a.DocArtifact.Id] = a
		}

		if a.ModelArtifact != nil {
			if a.ModelArtifact.Id == nil {
				newID := strconv.Itoa(i*2 + 1)

				a.ModelArtifact.Id = &newID
			}

			s.Artifacts[*a.ModelArtifact.Id] = a
		}
	}

	for i, ma := range seedData.ModelArtifacts {
		if ma.Id == nil {
			newID := strconv.Itoa(i + 1)
			ma.Id = &newID
		}

		s.ModelArtifacts[*ma.Id] = ma
	}

	for i, se := range seedData.ServingEnvironments {
		if se.Id == nil {
			newID := strconv.Itoa(i + 1)
			se.Id = &newID
		}

		s.ServingEnvironments[*se.Id] = se
	}

	for i, is := range seedData.InferenceServices {
		if is.Id == nil {
			newID := strconv.Itoa(i + 1)
			is.Id = &newID
		}

		s.InferenceServices[*is.Id] = is
	}

	for i, sm := range seedData.ServeModels {
		if sm.Id == nil {
			newID := strconv.Itoa(i + 1)
			sm.Id = &newID
		}

		s.ServeModels[*sm.Id] = sm
	}

	return nil
}

func (s *InMemory) UpsertRegisteredModel(registeredModel *openapi.RegisteredModel) (*openapi.RegisteredModel, error) {
	if registeredModel.Id == nil {
		newID := strconv.Itoa(len(s.RegisteredModels) + 1)
		registeredModel.Id = &newID

		s.RegisteredModels[newID] = registeredModel

		return registeredModel, nil
	}

	s.RegisteredModels[*registeredModel.Id] = registeredModel

	return registeredModel, nil
}

func (s *InMemory) GetRegisteredModelById(id string) (*openapi.RegisteredModel, error) {
	return s.RegisteredModels[id], nil
}

func (s *InMemory) GetRegisteredModelByInferenceService(inferenceServiceId string) (*openapi.RegisteredModel, error) {
	ifsvc, err := s.GetInferenceServiceById(inferenceServiceId)
	if err != nil {
		return nil, err
	}

	return s.RegisteredModels[ifsvc.RegisteredModelId], nil
}

func (s *InMemory) GetRegisteredModelByParams(name *string, externalId *string) (*openapi.RegisteredModel, error) {
	for _, rm := range s.RegisteredModels {
		if name != nil && rm.Name == *name {
			return rm, nil
		}

		if externalId != nil && rm.ExternalId != nil && *rm.ExternalId == *externalId {
			return rm, nil
		}
	}

	return nil, nil
}

func (s *InMemory) GetRegisteredModels(listOptions api.ListOptions) (*openapi.RegisteredModelList, error) {
	modelList := openapi.RegisteredModelList{}

	for _, rm := range s.RegisteredModels {
		modelList.Items = append(modelList.Items, *rm)
	}

	modelList.PageSize = 1
	modelList.Size = int32(len(modelList.Items))

	return &modelList, nil
}

func (s *InMemory) UpsertModelVersion(modelVersion *openapi.ModelVersion, registeredModelId *string) (*openapi.ModelVersion, error) {
	if modelVersion.Id == nil {
		newID := strconv.Itoa(len(s.ModelVersions) + 1)
		modelVersion.Id = &newID

		s.ModelVersions[newID] = modelVersion

		return modelVersion, nil
	}

	s.ModelVersions[*modelVersion.Id] = modelVersion

	return modelVersion, nil
}

func (s *InMemory) GetModelVersionById(id string) (*openapi.ModelVersion, error) {
	return s.ModelVersions[id], nil
}

func (s *InMemory) GetModelVersionByInferenceService(inferenceServiceId string) (*openapi.ModelVersion, error) {
	ifsvc, err := s.GetInferenceServiceById(inferenceServiceId)
	if err != nil {
		return nil, err
	}

	if ifsvc.ModelVersionId == nil {
		return nil, nil
	}

	return s.ModelVersions[*ifsvc.ModelVersionId], nil
}

func (s *InMemory) GetModelVersionByParams(versionName *string, registeredModelId *string, externalId *string) (*openapi.ModelVersion, error) {
	for _, mv := range s.ModelVersions {
		if versionName != nil && registeredModelId != nil {
			searchName := converter.PrefixWhenOwned(registeredModelId, *versionName)

			if mv.Name == searchName {
				return mv, nil
			}
		}

		if externalId != nil && mv.ExternalId != nil && *mv.ExternalId == *externalId {
			return mv, nil
		}
	}

	return nil, nil
}

func (s *InMemory) GetModelVersions(listOptions api.ListOptions, registeredModelId *string) (*openapi.ModelVersionList, error) {
	modelVersionList := openapi.ModelVersionList{}

	for _, mv := range s.ModelVersions {
		modelVersionList.Items = append(modelVersionList.Items, *mv)
	}

	modelVersionList.PageSize = 1
	modelVersionList.Size = int32(len(modelVersionList.Items))

	return &modelVersionList, nil
}

func (s *InMemory) UpsertModelVersionArtifact(artifact *openapi.Artifact, modelVersionId string) (*openapi.Artifact, error) {
	return nil, nil
}

func (s *InMemory) UpsertArtifact(artifact *openapi.Artifact) (*openapi.Artifact, error) {
	if artifact.DocArtifact != nil {
		if artifact.DocArtifact.Id == nil {
			newID := strconv.Itoa(len(s.Artifacts)*2 + 1)
			artifact.DocArtifact.Id = &newID

			s.Artifacts[newID] = artifact

			return artifact, nil
		}

		s.Artifacts[*artifact.DocArtifact.Id] = artifact

		return artifact, nil
	}

	if artifact.ModelArtifact != nil {
		if artifact.ModelArtifact.Id == nil {
			newID := strconv.Itoa(len(s.Artifacts)*2 + 1)
			artifact.ModelArtifact.Id = &newID

			s.Artifacts[newID] = artifact

			return artifact, nil
		}

		s.Artifacts[*artifact.ModelArtifact.Id] = artifact

		return artifact, nil
	}

	return nil, nil
}

func (s *InMemory) GetArtifactById(id string) (*openapi.Artifact, error) {
	return s.Artifacts[id], nil
}

func (s *InMemory) GetArtifactByParams(artifactName *string, modelVersionId *string, externalId *string) (*openapi.Artifact, error) {
	for _, a := range s.Artifacts {
		if artifactName != nil && modelVersionId != nil {
			searchName := converter.PrefixWhenOwned(modelVersionId, *artifactName)

			if a.DocArtifact != nil && *a.DocArtifact.Name == searchName {
				return a, nil
			}

			if a.ModelArtifact != nil && *a.ModelArtifact.Name == searchName {
				return a, nil
			}
		}

		if externalId != nil {
			if a.DocArtifact != nil && a.DocArtifact.ExternalId != nil && *a.DocArtifact.ExternalId == *externalId {
				return a, nil
			}

			if a.ModelArtifact != nil && a.ModelArtifact.ExternalId != nil && *a.ModelArtifact.ExternalId == *externalId {
				return a, nil
			}
		}
	}

	return nil, nil
}

func (s *InMemory) GetArtifacts(listOptions api.ListOptions, modelVersionId *string) (*openapi.ArtifactList, error) {
	artifactList := openapi.ArtifactList{}

	for _, a := range s.Artifacts {
		artifactList.Items = append(artifactList.Items, *a)
	}

	artifactList.PageSize = 1
	artifactList.Size = int32(len(artifactList.Items))

	return &artifactList, nil
}

func (s *InMemory) UpsertModelArtifact(modelArtifact *openapi.ModelArtifact) (*openapi.ModelArtifact, error) {
	if modelArtifact.Id == nil {
		newID := strconv.Itoa(len(s.ModelArtifacts) + 1)
		modelArtifact.Id = &newID

		s.ModelArtifacts[newID] = modelArtifact

		return modelArtifact, nil
	}

	s.ModelArtifacts[*modelArtifact.Id] = modelArtifact

	return modelArtifact, nil
}

func (s *InMemory) GetModelArtifactById(id string) (*openapi.ModelArtifact, error) {
	return s.ModelArtifacts[id], nil
}

func (s *InMemory) GetModelArtifactByInferenceService(inferenceServiceId string) (*openapi.ModelArtifact, error) {
	ifsvc, err := s.GetInferenceServiceById(inferenceServiceId)
	if err != nil {
		return nil, err
	}

	artifacts, err := s.GetModelArtifacts(api.ListOptions{}, ifsvc.ModelVersionId)
	if err != nil {
		return nil, err
	}

	if len(artifacts.Items) == 0 {
		return nil, fmt.Errorf("no model artifact found for inference service %s", inferenceServiceId)
	}

	return &artifacts.Items[0], nil
}

func (s *InMemory) GetModelArtifactByParams(artifactName *string, modelVersionId *string, externalId *string) (*openapi.ModelArtifact, error) {
	for _, ma := range s.ModelArtifacts {
		if artifactName != nil && modelVersionId != nil {
			searchName := converter.PrefixWhenOwned(modelVersionId, *artifactName)

			if ma.Name != nil && *ma.Name == searchName {
				return ma, nil
			}
		}

		if externalId != nil && ma.ExternalId != nil && *ma.ExternalId == *externalId {
			return ma, nil
		}
	}

	return nil, nil
}

func (s *InMemory) GetModelArtifacts(listOptions api.ListOptions, modelVersionId *string) (*openapi.ModelArtifactList, error) {
	modelArtifactList := openapi.ModelArtifactList{}

	for _, ma := range s.ModelArtifacts {
		modelArtifactList.Items = append(modelArtifactList.Items, *ma)
	}

	modelArtifactList.PageSize = 1
	modelArtifactList.Size = int32(len(modelArtifactList.Items))

	return &modelArtifactList, nil
}

func (s *InMemory) UpsertServingEnvironment(registeredModel *openapi.ServingEnvironment) (*openapi.ServingEnvironment, error) {
	if registeredModel.Id == nil {
		newID := strconv.Itoa(len(s.ServingEnvironments) + 1)
		registeredModel.Id = &newID

		s.ServingEnvironments[newID] = registeredModel

		return registeredModel, nil
	}

	s.ServingEnvironments[*registeredModel.Id] = registeredModel

	return registeredModel, nil
}

func (s *InMemory) GetServingEnvironmentById(id string) (*openapi.ServingEnvironment, error) {
	return s.ServingEnvironments[id], nil
}

func (s *InMemory) GetServingEnvironmentByParams(name *string, externalId *string) (*openapi.ServingEnvironment, error) {
	for _, se := range s.ServingEnvironments {
		if name != nil && se.Name != nil && *se.Name == *name {
			return se, nil
		}

		if externalId != nil && se.ExternalId != nil && *se.ExternalId == *externalId {
			return se, nil
		}
	}

	return nil, nil
}

func (s *InMemory) GetServingEnvironments(listOptions api.ListOptions) (*openapi.ServingEnvironmentList, error) {
	servingEnvironmentList := openapi.ServingEnvironmentList{}

	for _, se := range s.ServingEnvironments {
		servingEnvironmentList.Items = append(servingEnvironmentList.Items, *se)
	}

	servingEnvironmentList.PageSize = 1
	servingEnvironmentList.Size = int32(len(servingEnvironmentList.Items))

	return &servingEnvironmentList, nil
}

func (s *InMemory) UpsertInferenceService(inferenceService *openapi.InferenceService) (*openapi.InferenceService, error) {
	if inferenceService.Id == nil {
		newID := strconv.Itoa(len(s.InferenceServices) + 1)
		inferenceService.Id = &newID

		s.InferenceServices[newID] = inferenceService

		return inferenceService, nil
	}

	s.InferenceServices[*inferenceService.Id] = inferenceService

	return inferenceService, nil
}

func (s *InMemory) GetInferenceServiceById(id string) (*openapi.InferenceService, error) {
	return s.InferenceServices[id], nil
}

func (s *InMemory) GetInferenceServiceByParams(name *string, parentResourceId *string, externalId *string) (*openapi.InferenceService, error) {
	for _, is := range s.InferenceServices {
		if name != nil && parentResourceId != nil {
			searchName := converter.PrefixWhenOwned(parentResourceId, *name)

			if is.Name != nil && *is.Name == searchName {
				return is, nil
			}
		}

		if externalId != nil && is.ExternalId != nil && *is.ExternalId == *externalId {
			return is, nil
		}
	}

	return nil, nil
}

func (s *InMemory) GetInferenceServices(listOptions api.ListOptions, servingEnvironmentId *string, runtime *string) (*openapi.InferenceServiceList, error) {
	inferenceServiceList := openapi.InferenceServiceList{}

	items := []openapi.InferenceService{}

	for _, is := range s.InferenceServices {
		if servingEnvironmentId != nil && is.ServingEnvironmentId != *servingEnvironmentId {
			continue
		}

		if runtime != nil && is.Runtime != nil && *is.Runtime != *runtime {
			continue
		}

		items = append(items, *is)
	}

	inferenceServiceList.Items = items
	inferenceServiceList.PageSize = 1
	inferenceServiceList.Size = int32(len(items))

	return &inferenceServiceList, nil
}

func (s *InMemory) UpsertServeModel(serveModel *openapi.ServeModel, inferenceServiceId *string) (*openapi.ServeModel, error) {
	// TODO: Implement
	return nil, nil
}

func (s *InMemory) GetServeModelById(id string) (*openapi.ServeModel, error) {
	return s.ServeModels[id], nil
}

func (s *InMemory) GetServeModels(listOptions api.ListOptions, inferenceServiceId *string) (*openapi.ServeModelList, error) {
	// TODO: Implement
	return nil, nil
}

package services

import "github.com/kubeflow/model-registry/pkg/openapi"

type InMemorySeedData struct {
	RegisteredModels    []*openapi.RegisteredModel    `yaml:"registeredModels"`
	ModelVersions       []*openapi.ModelVersion       `yaml:"modelVersions"`
	Artifacts           []*openapi.Artifact           `yaml:"artifacts"`
	ModelArtifacts      []*openapi.ModelArtifact      `yaml:"modelArtifacts"`
	ServingEnvironments []*openapi.ServingEnvironment `yaml:"servingEnvironments"`
	InferenceServices   []*openapi.InferenceService   `yaml:"inferenceServices"`
	ServeModels         []*openapi.ServeModel         `yaml:"serveModels"`
}

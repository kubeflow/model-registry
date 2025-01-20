package repositories

import (
	"log/slog"
)

type ModelRegistryClientInterface interface {
	RegisteredModelInterface
	ModelVersionInterface
	ArtifactInterface
}

type ModelRegistryClient struct {
	logger *slog.Logger
	RegisteredModel
	ModelVersion
	Artifact
}

func NewModelRegistryClient(logger *slog.Logger) (ModelRegistryClientInterface, error) {
	return &ModelRegistryClient{logger: logger}, nil
}

package data

import (
	"log/slog"
)

type ModelRegistryClientInterface interface {
	RegisteredModelInterface
	ModelVersionInterface
}

type ModelRegistryClient struct {
	logger *slog.Logger
	RegisteredModel
	ModelVersion
}

func NewModelRegistryClient(logger *slog.Logger) (ModelRegistryClientInterface, error) {
	return &ModelRegistryClient{logger: logger}, nil
}

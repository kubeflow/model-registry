package data

import (
	"log/slog"
)

type ModelRegistryClientInterface interface {
	RegisteredModelInterface
}

type ModelRegistryClient struct {
	logger *slog.Logger
	RegisteredModel
}

func NewModelRegistryClient(logger *slog.Logger) (ModelRegistryClientInterface, error) {
	return &ModelRegistryClient{logger: logger}, nil
}

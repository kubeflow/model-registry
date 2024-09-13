package data

import (
	"encoding/json"
	"fmt"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/integrations"
	"net/url"
)

const modelVersionPath = "/model_versions"

type ModelVersionInterface interface {
	GetModelVersion(client integrations.HTTPClientInterface, id string) (*openapi.ModelVersion, error)
}

type ModelVersion struct {
	ModelVersionInterface
}

func (v ModelVersion) GetModelVersion(client integrations.HTTPClientInterface, id string) (*openapi.ModelVersion, error) {
	path, err := url.JoinPath(modelVersionPath, id)
	if err != nil {
		return nil, err
	}

	response, err := client.GET(path)

	if err != nil {
		return nil, fmt.Errorf("error fetching model version: %w", err)
	}

	var model openapi.ModelVersion
	if err := json.Unmarshal(response, &model); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &model, nil
}

package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations"
)

const modelArtifactPath = "/model_artifacts"

type ModelArtifactInterface interface {
	UpdateModelArtifact(client integrations.HTTPClientInterface, id string, jsonData []byte) (*openapi.ModelArtifact, error)
}

type ModelArtifact struct {
	ModelArtifactInterface
}

func (a ModelArtifact) UpdateModelArtifact(client integrations.HTTPClientInterface, id string, jsonData []byte) (*openapi.ModelArtifact, error) {
	path, err := url.JoinPath(modelArtifactPath, id)
	if err != nil {
		return nil, err
	}

	responseData, err := client.PATCH(path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error patching registered model: %w", err)
	}

	var modelArtifact openapi.ModelArtifact
	if err := json.Unmarshal(responseData, &modelArtifact); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &modelArtifact, nil
}

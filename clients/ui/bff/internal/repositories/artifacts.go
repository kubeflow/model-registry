package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations"
)

const artifactPath = "/artifacts"

type ArtifactInterface interface {
	GetAllArtifacts(client integrations.HTTPClientInterface, pageValues url.Values) (*openapi.ArtifactList, error)
	GetArtifact(client integrations.HTTPClientInterface, id string) (*openapi.Artifact, error)
	CreateArtifact(client integrations.HTTPClientInterface, jsonData []byte) (*openapi.Artifact, error)
	UpdateArtifact(client integrations.HTTPClientInterface, id string, jsonData []byte) (*openapi.Artifact, error)
}

type Artifact struct {
	ArtifactInterface
}

func (a Artifact) GetAllArtifacts(client integrations.HTTPClientInterface, pageValues url.Values) (*openapi.ArtifactList, error) {
	responseData, err := client.GET(UrlWithPageParams(artifactPath, pageValues))
	if err != nil {
		return nil, fmt.Errorf("error fetching artifacts: %w", err)
	}

	var artifacts openapi.ArtifactList
	if err := json.Unmarshal(responseData, &artifacts); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &artifacts, nil
}

func (a Artifact) GetArtifact(client integrations.HTTPClientInterface, id string) (*openapi.Artifact, error) {
	path, err := url.JoinPath(artifactPath, id)
	if err != nil {
		return nil, err
	}

	responseData, err := client.GET(path)
	if err != nil {
		return nil, fmt.Errorf("error fetching artifacts: %w", err)
	}

	var artifact openapi.Artifact
	if err := json.Unmarshal(responseData, &artifact); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &artifact, nil
}

func (a Artifact) CreateArtifact(client integrations.HTTPClientInterface, jsonData []byte) (*openapi.Artifact, error) {
	responseData, err := client.POST(artifactPath, bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, fmt.Errorf("error creating artifact: %w", err)
	}

	var artifact openapi.Artifact
	if err := json.Unmarshal(responseData, &artifact); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &artifact, nil
}

func (a Artifact) UpdateArtifact(client integrations.HTTPClientInterface, id string, jsonData []byte) (*openapi.Artifact, error) {
	path, err := url.JoinPath(artifactPath, id)
	if err != nil {
		return nil, err
	}

	responseData, err := client.PATCH(path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error patching registered model: %w", err)
	}

	var artifact openapi.Artifact
	if err := json.Unmarshal(responseData, &artifact); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &artifact, nil
}

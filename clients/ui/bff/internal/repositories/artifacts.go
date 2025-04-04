package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/mrserver"
	"net/url"

	"github.com/kubeflow/model-registry/pkg/openapi"
)

const artifactPath = "/artifacts"

type ArtifactInterface interface {
	GetAllArtifacts(client mrserver.HTTPClientInterface, pageValues url.Values) (*openapi.ArtifactList, error)
	GetArtifact(client mrserver.HTTPClientInterface, id string) (*openapi.Artifact, error)
	CreateArtifact(client mrserver.HTTPClientInterface, jsonData []byte) (*openapi.Artifact, error)
}

type Artifact struct {
	ArtifactInterface
}

func (a Artifact) GetAllArtifacts(client mrserver.HTTPClientInterface, pageValues url.Values) (*openapi.ArtifactList, error) {
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

func (a Artifact) GetArtifact(client mrserver.HTTPClientInterface, id string) (*openapi.Artifact, error) {
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

func (a Artifact) CreateArtifact(client mrserver.HTTPClientInterface, jsonData []byte) (*openapi.Artifact, error) {
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

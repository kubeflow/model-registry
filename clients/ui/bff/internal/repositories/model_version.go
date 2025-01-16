package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations"
	"net/url"
)

const modelVersionPath = "/model_versions"
const artifactsByModelVersionPath = "/artifacts"

type ModelVersionInterface interface {
	GetAllModelVersions(client integrations.HTTPClientInterface) (*openapi.ModelVersionList, error)
	GetModelVersion(client integrations.HTTPClientInterface, id string) (*openapi.ModelVersion, error)
	CreateModelVersion(client integrations.HTTPClientInterface, jsonData []byte) (*openapi.ModelVersion, error)
	UpdateModelVersion(client integrations.HTTPClientInterface, id string, jsonData []byte) (*openapi.ModelVersion, error)
	GetModelArtifactsByModelVersion(client integrations.HTTPClientInterface, id string, pageValues url.Values) (*openapi.ModelArtifactList, error)
	CreateModelArtifactByModelVersion(client integrations.HTTPClientInterface, id string, jsonData []byte) (*openapi.ModelArtifact, error)
}

type ModelVersion struct {
	ModelVersionInterface
}

func (v ModelVersion) GetAllModelVersions(client integrations.HTTPClientInterface) (*openapi.ModelVersionList, error) {
	response, err := client.GET(modelVersionPath)

	if err != nil {
		return nil, fmt.Errorf("error fetching model versions: %w", err)
	}

	var models openapi.ModelVersionList
	if err := json.Unmarshal(response, &models); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &models, nil
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

func (v ModelVersion) CreateModelVersion(client integrations.HTTPClientInterface, jsonData []byte) (*openapi.ModelVersion, error) {
	responseData, err := client.POST(modelVersionPath, bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, fmt.Errorf("error posting registered model: %w", err)
	}

	var model openapi.ModelVersion
	if err := json.Unmarshal(responseData, &model); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &model, nil
}

func (v ModelVersion) UpdateModelVersion(client integrations.HTTPClientInterface, id string, jsonData []byte) (*openapi.ModelVersion, error) {
	path, err := url.JoinPath(modelVersionPath, id)

	if err != nil {
		return nil, err
	}

	responseData, err := client.PATCH(path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error patching ModelVersion: %w", err)
	}

	var model openapi.ModelVersion
	if err := json.Unmarshal(responseData, &model); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &model, nil
}

func (v ModelVersion) GetModelArtifactsByModelVersion(client integrations.HTTPClientInterface, id string, pageValues url.Values) (*openapi.ModelArtifactList, error) {
	path, err := url.JoinPath(modelVersionPath, id, artifactsByModelVersionPath)

	if err != nil {
		return nil, err
	}

	responseData, err := client.GET(UrlWithPageParams(path, pageValues))
	if err != nil {
		return nil, fmt.Errorf("error fetching model version artifacts: %w", err)
	}

	var model openapi.ModelArtifactList
	if err := json.Unmarshal(responseData, &model); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &model, nil
}

func (v ModelVersion) CreateModelArtifactByModelVersion(client integrations.HTTPClientInterface, id string, jsonData []byte) (*openapi.ModelArtifact, error) {
	path, err := url.JoinPath(modelVersionPath, id, artifactsByModelVersionPath)
	if err != nil {
		return nil, err
	}

	responseData, err := client.POST(path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error posting model artifact: %w", err)
	}

	var model openapi.ModelArtifact
	if err := json.Unmarshal(responseData, &model); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &model, nil
}

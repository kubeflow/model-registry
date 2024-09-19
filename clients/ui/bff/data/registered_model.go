package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/integrations"
	"net/url"
)

const registeredModelPath = "/registered_models"
const versionsPath = "/versions"

type RegisteredModelInterface interface {
	GetAllRegisteredModels(client integrations.HTTPClientInterface) (*openapi.RegisteredModelList, error)
	CreateRegisteredModel(client integrations.HTTPClientInterface, jsonData []byte) (*openapi.RegisteredModel, error)
	GetRegisteredModel(client integrations.HTTPClientInterface, id string) (*openapi.RegisteredModel, error)
	UpdateRegisteredModel(client integrations.HTTPClientInterface, id string, jsonData []byte) (*openapi.RegisteredModel, error)
	GetAllModelVersions(client integrations.HTTPClientInterface, id string) (*openapi.ModelVersionList, error)
	CreateModelVersionForRegisteredModel(client integrations.HTTPClientInterface, id string, jsonData []byte) (*openapi.ModelVersion, error)
}

type RegisteredModel struct {
	RegisteredModelInterface
}

func (m RegisteredModel) GetAllRegisteredModels(client integrations.HTTPClientInterface) (*openapi.RegisteredModelList, error) {

	responseData, err := client.GET(registeredModelPath)
	if err != nil {
		return nil, fmt.Errorf("error fetching registered models: %w", err)
	}

	var modelList openapi.RegisteredModelList
	if err := json.Unmarshal(responseData, &modelList); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &modelList, nil
}

func (m RegisteredModel) CreateRegisteredModel(client integrations.HTTPClientInterface, jsonData []byte) (*openapi.RegisteredModel, error) {
	responseData, err := client.POST(registeredModelPath, bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, fmt.Errorf("error posting registered model: %w", err)
	}

	var model openapi.RegisteredModel
	if err := json.Unmarshal(responseData, &model); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &model, nil
}

func (m RegisteredModel) GetRegisteredModel(client integrations.HTTPClientInterface, id string) (*openapi.RegisteredModel, error) {
	path, err := url.JoinPath(registeredModelPath, id)
	if err != nil {
		return nil, err
	}
	responseData, err := client.GET(path)

	if err != nil {
		return nil, fmt.Errorf("error fetching registered model: %w", err)
	}

	var model openapi.RegisteredModel
	if err := json.Unmarshal(responseData, &model); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &model, nil
}

func (m RegisteredModel) UpdateRegisteredModel(client integrations.HTTPClientInterface, id string, jsonData []byte) (*openapi.RegisteredModel, error) {
	path, err := url.JoinPath(registeredModelPath, id)

	if err != nil {
		return nil, err
	}

	responseData, err := client.PATCH(path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error patching registered model: %w", err)
	}

	var model openapi.RegisteredModel
	if err := json.Unmarshal(responseData, &model); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &model, nil
}

func (m RegisteredModel) GetAllModelVersions(client integrations.HTTPClientInterface, id string) (*openapi.ModelVersionList, error) {
	path, err := url.JoinPath(registeredModelPath, id, versionsPath)

	if err != nil {
		return nil, err
	}

	responseData, err := client.GET(path)

	if err != nil {
		return nil, fmt.Errorf("error fetching model versions: %w", err)
	}

	var model openapi.ModelVersionList
	if err := json.Unmarshal(responseData, &model); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &model, nil
}

func (m RegisteredModel) CreateModelVersionForRegisteredModel(client integrations.HTTPClientInterface, id string, jsonData []byte) (*openapi.ModelVersion, error) {
	path, err := url.JoinPath(registeredModelPath, id, versionsPath)

	if err != nil {
		return nil, err
	}

	responseData, err := client.POST(path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error posting model version: %w", err)
	}

	var model openapi.ModelVersion
	if err := json.Unmarshal(responseData, &model); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &model, nil
}

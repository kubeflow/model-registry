package proxy

import (
	"fmt"

	platformproxy "github.com/kubeflow/hub/internal/platform/proxy"
	"github.com/kubeflow/hub/pkg/api"
)

// Unexported constants for test compatibility
const (
	responseOK   = "OK"
	responseFail = "FAIL"
)

// Model-registry-specific constants
const (
	HealthCheckModelRegistry = "model-registry"

	detailDatastoreType                 = "datastore_type"
	detailRegisteredModelsAccessible    = "registered_models_accessible"
	detailRegisteredModelsCount         = "registered_models_count"
	detailArtifactsAccessible           = "artifacts_accessible"
	detailArtifactsCount                = "artifacts_count"
	detailModelVersionsAccessible       = "model_versions_accessible"
	detailModelVersionsCount            = "model_versions_count"
	detailServingEnvironmentsAccessible = "serving_environments_accessible"
	detailServingEnvironmentsCount      = "serving_environments_count"
	detailInferenceServicesAccessible    = "inference_services_accessible"
	detailInferenceServicesCount        = "inference_services_count"
	detailExperimentsAccessible         = "experiments_accessible"
	detailExperimentsCount              = "experiments_count"
	detailExperimentRunsAccessible      = "experiment_runs_accessible"
	detailExperimentRunsCount           = "experiment_runs_count"
	detailTotalResourcesChecked         = "total_resources_checked"
)

// ModelRegistryHealthChecker checks model registry service health
type ModelRegistryHealthChecker struct {
	service api.ModelRegistryApi
}

func NewModelRegistryHealthChecker(service api.ModelRegistryApi) *ModelRegistryHealthChecker {
	return &ModelRegistryHealthChecker{
		service: service,
	}
}

func (m *ModelRegistryHealthChecker) Check() platformproxy.HealthCheck {
	check := platformproxy.HealthCheck{
		Name:    HealthCheckModelRegistry,
		Details: make(map[string]any),
	}

	if m.service == nil {
		check.Status = platformproxy.StatusFail
		check.Message = "model registry service not available"
		return check
	}

	listOptions := api.ListOptions{
		PageSize: func() *int32 { i := int32(1); return &i }(),
	}

	models, err := m.service.GetRegisteredModels(listOptions)
	if err != nil {
		check.Status = platformproxy.StatusFail
		check.Message = fmt.Sprintf("failed to list registered models: %v", err)
		return check
	}
	check.Details[detailRegisteredModelsAccessible] = true
	check.Details[detailRegisteredModelsCount] = models.Size

	artifacts, err := m.service.GetArtifacts("", listOptions, nil)
	if err != nil {
		check.Status = platformproxy.StatusFail
		check.Message = fmt.Sprintf("failed to list artifacts: %v", err)
		return check
	}
	check.Details[detailArtifactsAccessible] = true
	check.Details[detailArtifactsCount] = artifacts.Size

	versions, err := m.service.GetModelVersions(listOptions, nil)
	if err != nil {
		check.Status = platformproxy.StatusFail
		check.Message = fmt.Sprintf("failed to list model versions: %v", err)
		return check
	}
	check.Details[detailModelVersionsAccessible] = true
	check.Details[detailModelVersionsCount] = versions.Size

	servingEnvs, err := m.service.GetServingEnvironments(listOptions)
	if err != nil {
		check.Status = platformproxy.StatusFail
		check.Message = fmt.Sprintf("failed to list serving environments: %v", err)
		return check
	}
	check.Details[detailServingEnvironmentsAccessible] = true
	check.Details[detailServingEnvironmentsCount] = servingEnvs.Size

	inferenceServices, err := m.service.GetInferenceServices(listOptions, nil, nil)
	if err != nil {
		check.Status = platformproxy.StatusFail
		check.Message = fmt.Sprintf("failed to list inference services: %v", err)
		return check
	}
	check.Details[detailInferenceServicesAccessible] = true
	check.Details[detailInferenceServicesCount] = inferenceServices.Size

	experiments, err := m.service.GetExperiments(listOptions)
	if err != nil {
		check.Status = platformproxy.StatusFail
		check.Message = fmt.Sprintf("failed to list experiments: %v", err)
		return check
	}
	check.Details[detailExperimentsAccessible] = true
	check.Details[detailExperimentsCount] = experiments.Size

	experimentRuns, err := m.service.GetExperimentRuns(listOptions, nil)
	if err != nil {
		check.Status = platformproxy.StatusFail
		check.Message = fmt.Sprintf("failed to list experiment runs: %v", err)
		return check
	}
	check.Details[detailExperimentRunsAccessible] = true
	check.Details[detailExperimentRunsCount] = experimentRuns.Size

	check.Status = platformproxy.StatusPass
	check.Message = "model registry service is healthy"
	check.Details[detailTotalResourcesChecked] = 5

	return check
}

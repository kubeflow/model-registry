package repositories

import (
	"context"
	"fmt"

	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	helper "github.com/kubeflow/model-registry/ui/bff/internal/helpers"
	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ModelRegistryRepository struct {
}

func NewModelRegistryRepository() *ModelRegistryRepository {
	return &ModelRegistryRepository{}
}

func (m *ModelRegistryRepository) GetAllModelRegistries(sessionCtx context.Context, client k8s.KubernetesClientInterface, namespace string) ([]models.ModelRegistryModel, error) {
	// Default to non-federated mode for backward compatibility
	return m.GetAllModelRegistriesWithMode(sessionCtx, client, namespace, false)
}

// GetAllModelRegistriesWithMode fetches all model registries with support for federated mode
func (m *ModelRegistryRepository) GetAllModelRegistriesWithMode(sessionCtx context.Context, client k8s.KubernetesClientInterface, namespace string, isFederatedMode bool) ([]models.ModelRegistryModel, error) {
	logger := helper.GetContextLogger(sessionCtx)
	logger.Debug("GetAllModelRegistriesWithMode called", "namespace", namespace, "isFederatedMode", isFederatedMode)

	var resources []k8s.ServiceDetails
	var err error

	// Check if we have authorization context from the middleware
	if authCtx, ok := sessionCtx.Value(constants.ServiceAuthorizationContextKey).(*models.ServiceAuthorizationContext); ok {
		if authCtx.AllowList {
			logger.Debug("User can get list services ")
			resources, err = client.GetServiceDetails(sessionCtx, namespace)
		} else {
			logger.Debug("User has limited access - we need use Rule base access",
				"serviceCount", len(authCtx.AllowedServiceNames),
				"services", authCtx.AllowedServiceNames)
			resources, err = m.getSpecificServiceDetails(sessionCtx, client, namespace, authCtx.AllowedServiceNames)
		}
	} else {
		logger.Warn("No authorization context found - using fallback behavior")
		resources, err = client.GetServiceDetails(sessionCtx, namespace)
	}

	if err != nil {
		logger.Error("Error fetching service details", "error", err, "namespace", namespace)
		return nil, fmt.Errorf("error fetching model registries: %w", err)
	}

	var registries = []models.ModelRegistryModel{}
	for _, s := range resources {
		serverAddress := m.ResolveServerAddress(s.ClusterIP, s.HTTPPort, s.IsHTTPS, s.ExternalAddressRest, isFederatedMode)
		registry := models.ModelRegistryModel{
			Name:          s.Name,
			Description:   s.Description,
			DisplayName:   s.DisplayName,
			ServerAddress: serverAddress,
			IsHTTPS:       s.IsHTTPS,
		}
		registries = append(registries, registry)
	}

	return registries, nil
}

// getSpecificServiceDetails fetches details for specific services by name
func (m *ModelRegistryRepository) getSpecificServiceDetails(sessionCtx context.Context, client k8s.KubernetesClientInterface, namespace string, serviceNames []string) ([]k8s.ServiceDetails, error) {
	logger := helper.GetContextLogger(sessionCtx)
	logger.Debug("getSpecificServiceDetails called", "namespace", namespace, "serviceNames", serviceNames)

	var resources []k8s.ServiceDetails

	for _, serviceName := range serviceNames {
		logger.Debug("Fetching service details", "serviceName", serviceName, "namespace", namespace)
		// Validate if service is a model registry service by passing the component label value
		serviceDetail, err := client.GetServiceDetailsByName(sessionCtx, namespace, serviceName, k8s.ComponentLabelValue)
		if err != nil {
			logger.Warn("Failed to get service details, skipping",
				"serviceName", serviceName,
				"namespace", namespace,
				"error", err)
			// Log the error but continue with other services
			continue
		}
		logger.Debug("Service details retrieved successfully", "serviceName", serviceName)
		resources = append(resources, serviceDetail)
	}
	return resources, nil
}

func (m *ModelRegistryRepository) GetModelRegistry(sessionCtx context.Context, client k8s.KubernetesClientInterface, namespace string, modelRegistryID string) (models.ModelRegistryModel, error) {
	// Default to non-federated mode for backward compatibility
	return m.GetModelRegistryWithMode(sessionCtx, client, namespace, modelRegistryID, false)
}

// GetModelRegistryWithMode fetches a specific model registry with support for federated mode
func (m *ModelRegistryRepository) GetModelRegistryWithMode(sessionCtx context.Context, client k8s.KubernetesClientInterface, namespace string, modelRegistryID string, isFederatedMode bool) (models.ModelRegistryModel, error) {

	s, err := client.GetServiceDetailsByName(sessionCtx, namespace, modelRegistryID, k8s.ComponentLabelValue)
	if err != nil {
		return models.ModelRegistryModel{}, fmt.Errorf("error fetching model registry: %w", err)
	}

	modelRegistry := models.ModelRegistryModel{
		Name:          s.Name,
		Description:   s.Description,
		DisplayName:   s.DisplayName,
		ServerAddress: m.ResolveServerAddress(s.ClusterIP, s.HTTPPort, s.IsHTTPS, s.ExternalAddressRest, isFederatedMode),
		IsHTTPS:       s.IsHTTPS,
	}

	return modelRegistry, nil
}

func (m *ModelRegistryRepository) ResolveServerAddress(clusterIP string, httpPort int32, isHTTPS bool, externalAddressRest string, isFederatedMode bool) string {
	// Default behavior - use cluster IP and port
	protocol := "http"
	if isHTTPS {
		protocol = "https"
	}
	// In federated mode, if external address is available, use it
	if isFederatedMode && externalAddressRest != "" {
		// External address is assumed to be HTTPS
		url := fmt.Sprintf("%s://%s/api/model_registry/v1alpha3", protocol, externalAddressRest)
		return url
	}

	url := fmt.Sprintf("%s://%s:%d/api/model_registry/v1alpha3", protocol, clusterIP, httpPort)
	return url
}

// GetAllModelTransferJobs returns just one mock sample to unblock the UI work and the rest of the logic will be added in followup PR
// TODO: Replace with actual implementation for all the methods
func (m *ModelRegistryRepository) GetAllModelTransferJobs(ctx context.Context, client k8s.KubernetesClientInterface, namespace string) (*models.ModelTransferJobList, error) {
	jobList, err := client.GetAllModelTransferJobs(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch model transfer jobs: %w", err)
	}

	var transferJobs []models.ModelTransferJob
	for _, job := range jobList.Items {
		transferJobs = append(transferJobs, convertK8sJobToModel(&job))
	}

	return &models.ModelTransferJobList{
		Items:    transferJobs,
		Size:     len(transferJobs),
		PageSize: len(transferJobs),
	}, nil

}

func (m *ModelRegistryRepository) CreateModelTransferJob(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, payload models.ModelTransferJob) error {
	job := convertModelToK8sJob(payload, namespace)

	err := client.CreateModelTransferJob(ctx, namespace, job)
	if err != nil {
		return fmt.Errorf("failed to create model transfer job: %w", err)
	}
	return nil
}

func (m *ModelRegistryRepository) UpdateModelTransferJob(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, jobId string, updates map[string]string) error {
	err := client.UpdateModelTransferJob(ctx, namespace, jobId, updates)
	if err != nil {
		return fmt.Errorf("failed to update model transfer job %s: %w", jobId, err)
	}
	return nil
}

func (m *ModelRegistryRepository) DeleteModelTransferJob(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, jobId string) error {
	err := client.DeleteModelTransferJob(ctx, namespace, jobId)
	if err != nil {
		return fmt.Errorf("failed to delete model transfer job %s: %w", jobId, err)
	}
	return nil
}

// TODO: These functions convert the minimum required fields for now. Improve these to convert all the necessary fields
func convertModelToK8sJob(payload models.ModelTransferJob, namespace string) *batchv1.Job {
	backoffLimit := int32(3)

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      payload.Name,
			Namespace: namespace,
			Labels: map[string]string{
				"modelregistry.kubeflow.org/job-type": "async-upload",
				"modelregistry.kubeflow.org/job-id":   payload.Id,
			},
			Annotations: map[string]string{
				"modelregistry.kubeflow.org/source-type":   string(payload.Source.Type),
				"modelregistry.kubeflow.org/dest-type":     string(payload.Destination.Type),
				"modelregistry.kubeflow.org/model-name":    payload.RegisteredModelName,
				"modelregistry.kubeflow.org/upload-intent": string(payload.UploadIntent),
			},
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: &backoffLimit,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:  "async-upload",
							Image: "quay.io/opendatahub/model-registry-job-async-upload:latest",
						},
					},
				},
			},
		},
	}
}

func convertK8sJobToModel(job *batchv1.Job) models.ModelTransferJob {
	annotations := job.Annotations
	if annotations == nil {
		annotations = map[string]string{}
	}
	labels := job.Labels
	if labels == nil {
		labels = map[string]string{}
	}

	status := models.ModelTransferJobStatusPending
	if job.Status.Succeeded > 0 {
		status = models.ModelTransferJobStatusCompleted
	} else if job.Status.Failed > 0 {
		status = models.ModelTransferJobStatusFailed
	} else if job.Status.Active > 0 {
		status = models.ModelTransferJobStatusRunning
	}

	return models.ModelTransferJob{
		Id:                   labels["modelregistry.kubeflow.org/job-id"],
		Name:                 job.Name,
		Description:          annotations["modelregistry.kubeflow.org/description"],
		RegisteredModelName:  annotations["modelregistry.kubeflow.org/model-name"],
		ModelVersionName:     annotations["modelregistry.kubeflow.org/version-name"],
		RegisteredModelId:    annotations["modelregistry.kubeflow.org/registered-model-id"],
		ModelVersionId:       annotations["modelregistry.kubeflow.org/model-version-id"],
		ModelArtifactId:      annotations["modelregistry.kubeflow.org/model-artifact-id"],
		Author:               annotations["modelregistry.kubeflow.org/author"],
		Status:               status,
		CreateTimeSinceEpoch: fmt.Sprintf("%d", job.CreationTimestamp.UnixMilli()),
		Namespace:            job.Namespace,
	}
}

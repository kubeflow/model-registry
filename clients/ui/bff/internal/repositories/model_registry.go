package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	helper "github.com/kubeflow/model-registry/ui/bff/internal/helpers"
	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ErrModelTransferJobNotFound is returned when the model transfer job does not exist.
var ErrModelTransferJobNotFound = errors.New("model transfer job not found")

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

func (m *ModelRegistryRepository) GetAllModelTransferJobs(ctx context.Context, client k8s.KubernetesClientInterface, namespace string) (*models.ModelTransferJobList, error) {
	logger := helper.GetContextLogger(ctx)

	jobList, err := client.GetAllModelTransferJobs(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch model transfer jobs: %w", err)
	}

	jobNames := make([]string, 0, len(jobList.Items))
	for _, job := range jobList.Items {
		jobNames = append(jobNames, job.Name)
	}

	podsByJob := map[string][]corev1.Pod{}
	podNameToJobName := map[string]string{}
	if len(jobNames) > 0 {
		podList, podErr := client.GetTransferJobPods(ctx, namespace, jobNames)
		if podErr != nil {
			logger.Warn("failed to fetch transfer job pods, continuing without pod data", slog.Any("error", podErr))
		} else {
			for i := range podList.Items {
				pod := podList.Items[i]
				jobName := pod.Labels["job-name"]
				if jobName != "" {
					podsByJob[jobName] = append(podsByJob[jobName], pod)
					podNameToJobName[pod.Name] = jobName
				}
			}
		}
	}

	eventsByJob := map[string][]models.ModelTransferJobEvent{}
	if len(podNameToJobName) > 0 {
		podNames := make([]string, 0, len(podNameToJobName))
		for podName := range podNameToJobName {
			podNames = append(podNames, podName)
		}
		eventList, eventErr := client.GetEventsForPods(ctx, namespace, podNames)
		if eventErr != nil {
			logger.Warn("failed to fetch events for transfer job pods, continuing without events", slog.Any("error", eventErr))
		} else {
			for _, event := range eventList.Items {
				jobName := podNameToJobName[event.InvolvedObject.Name]
				if jobName == "" {
					continue
				}
				ts := event.LastTimestamp.UTC().Format("2006-01-02T15:04:05Z")
				if event.LastTimestamp.IsZero() && !event.EventTime.IsZero() {
					ts = event.EventTime.UTC().Format("2006-01-02T15:04:05Z")
				}
				eventsByJob[jobName] = append(eventsByJob[jobName], models.ModelTransferJobEvent{
					Timestamp: ts,
					Type:      event.Type,
					Reason:    event.Reason,
					Message:   event.Message,
				})
			}
		}
	}

	var transferJobs []models.ModelTransferJob
	for _, job := range jobList.Items {
		transferJob := convertK8sJobToModel(&job)
		if pods, ok := podsByJob[job.Name]; ok {
			enrichJobFromPods(&transferJob, pods, logger)
		}
		if events, ok := eventsByJob[job.Name]; ok {
			transferJob.Events = events
		}
		if transferJob.Events == nil {
			transferJob.Events = []models.ModelTransferJobEvent{}
		}
		transferJobs = append(transferJobs, transferJob)
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

func (m *ModelRegistryRepository) DeleteModelTransferJob(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, jobName string) error {
	err := client.DeleteModelTransferJob(ctx, namespace, jobName)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return fmt.Errorf("%w: %s", ErrModelTransferJobNotFound, jobName)
		}
		return fmt.Errorf("failed to delete model transfer job %s: %w", jobName, err)
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
	var errorMessage, errorDescription string
	if job.Status.Succeeded > 0 {
		status = models.ModelTransferJobStatusCompleted
	} else if job.Status.Failed > 0 {
		status = models.ModelTransferJobStatusFailed
		for _, condition := range job.Status.Conditions {
			if condition.Type == batchv1.JobFailed && condition.Status == corev1.ConditionTrue {
				errorMessage = condition.Reason
				errorDescription = condition.Message
				break
			}
		}
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
		ErrorMessage:         errorMessage,
		ErrorDescription:     errorDescription,
		CreateTimeSinceEpoch: fmt.Sprintf("%d", job.CreationTimestamp.UnixMilli()),
		Namespace:            job.Namespace,
	}
}

type terminationMessageIDs struct {
	RegisteredModel struct {
		ID string `json:"id"`
	} `json:"RegisteredModel"`
	ModelVersion struct {
		ID string `json:"id"`
	} `json:"ModelVersion"`
	ModelArtifact struct {
		ID string `json:"id"`
	} `json:"ModelArtifact"`
}

func enrichJobFromPods(job *models.ModelTransferJob, pods []corev1.Pod, logger *slog.Logger) {
	for _, pod := range pods {
		for _, cs := range pod.Status.ContainerStatuses {
			if cs.State.Terminated == nil {
				continue
			}
			msg := cs.State.Terminated.Message
			if msg == "" {
				continue
			}

			if cs.State.Terminated.ExitCode == 0 {
				var ids terminationMessageIDs
				if err := json.Unmarshal([]byte(msg), &ids); err != nil {
					logger.Debug("failed to parse termination message as JSON",
						slog.String("pod", pod.Name),
						slog.Any("error", err))
					continue
				}
				if ids.RegisteredModel.ID != "" {
					job.RegisteredModelId = ids.RegisteredModel.ID
				}
				if ids.ModelVersion.ID != "" {
					job.ModelVersionId = ids.ModelVersion.ID
				}
				if ids.ModelArtifact.ID != "" {
					job.ModelArtifactId = ids.ModelArtifact.ID
				}
				return
			}

			if job.ErrorMessage == "" {
				job.ErrorMessage = msg
			}
			return
		}
	}
}

package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
	helper "github.com/kubeflow/model-registry/ui/bff/internal/helpers"
	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation"
)

var (
	ErrJobNotFound           = errors.New("model transfer job not found")
	ErrJobValidationFailed   = errors.New("validation failed")
	ErrModelRegistryNotFound = errors.New("model registry not found in the selected namespace")
)

func (m *ModelRegistryRepository) GetAllModelTransferJobs(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, modelRegistryID string) (*models.ModelTransferJobList, error) {
	if modelRegistryID == "" {
		return &models.ModelTransferJobList{Items: []models.ModelTransferJob{}, Size: 0, PageSize: 0}, nil
	}

	jobList, err := client.GetAllModelTransferJobs(ctx, namespace, modelRegistryID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch model transfer jobs: %w", err)
	}

	transferJobs := make([]models.ModelTransferJob, 0, len(jobList.Items))
	for _, job := range jobList.Items {
		if job.DeletionTimestamp != nil {
			continue
		}
		transferJobs = append(transferJobs, convertK8sJobToModel(&job))
	}

	return &models.ModelTransferJobList{
		Items:    transferJobs,
		Size:     len(transferJobs),
		PageSize: len(transferJobs),
	}, nil
}

func (m *ModelRegistryRepository) CreateModelTransferJob(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, payload models.ModelTransferJob, modelRegistryID string) (*models.ModelTransferJob, error) {
	return m.createModelTransferJobResources(ctx, client, namespace, payload, modelRegistryID, "")
}

func (m *ModelRegistryRepository) createModelTransferJobResources(
	ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	payload models.ModelTransferJob,
	modelRegistryID string,
	existingDestSecretName string,
) (*models.ModelTransferJob, error) {
	payload.Source.Bucket = strings.TrimSpace(payload.Source.Bucket)
	payload.Source.Key = strings.TrimSpace(payload.Source.Key)
	payload.Source.URI = strings.TrimSpace(payload.Source.URI)
	payload.Destination.URI = strings.TrimSpace(payload.Destination.URI)
	payload.Destination.Registry = strings.TrimSpace(payload.Destination.Registry)

	skipDestCredsValidation := existingDestSecretName != ""
	if err := validateCreatePayload(payload, skipDestCredsValidation); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	if payload.Destination.Registry == "" && payload.Destination.URI != "" {
		payload.Destination.Registry, _ = extractRegistryFromURI(payload.Destination.URI)
	}

	logger := helper.GetContextLogger(ctx)

	modelRegistryAddress, err := m.getModelRegistryAddress(ctx, client, namespace, modelRegistryID)
	if err != nil {
		return nil, err
	}

	jobID := uuid.NewString()
	jobName := payload.Name

	configMapName := fmt.Sprintf("%s-metadata-configmap", jobName)
	destSecretName := fmt.Sprintf("%s-dest-creds", jobName)
	if existingDestSecretName != "" {
		destSecretName = existingDestSecretName
	}

	configMap := buildModelMetadataConfigMap(configMapName, namespace, payload, jobID, jobName)
	if err := client.CreateConfigMap(ctx, namespace, configMap); err != nil {
		return nil, fmt.Errorf("failed to create metadata configmap: %w", err)
	}

	var sourceSecretName string
	if payload.Source.Type == models.ModelTransferJobSourceTypeS3 {
		sourceSecretName = fmt.Sprintf("%s-source-creds", jobName)
		sourceSecret := buildSourceSecret(sourceSecretName, namespace, payload, jobID)
		if err := client.CreateSecret(ctx, namespace, sourceSecret); err != nil {
			cleanupCreatedResources(ctx, client, namespace, configMapName, "", "")
			return nil, fmt.Errorf("failed to create source secret: %w", err)
		}
	}

	if existingDestSecretName == "" {
		destSecret, err := buildDestinationSecret(destSecretName, namespace, payload, jobID)
		if err != nil {
			cleanupCreatedResources(ctx, client, namespace, configMapName, sourceSecretName, "")
			return nil, fmt.Errorf("failed to build destination secret: %w", err)
		}
		if err := client.CreateSecret(ctx, namespace, destSecret); err != nil {
			cleanupCreatedResources(ctx, client, namespace, configMapName, sourceSecretName, "")
			return nil, fmt.Errorf("failed to create destination secret: %w", err)
		}
	}

	job := buildK8sJob(jobName, namespace, jobID, payload, configMapName, sourceSecretName, destSecretName, modelRegistryAddress, modelRegistryID)
	createdJob, err := client.CreateModelTransferJob(ctx, namespace, job)
	if err != nil {
		cleanupCreatedResources(ctx, client, namespace, configMapName, sourceSecretName, destSecretName)
		if apierrors.IsAlreadyExists(err) {
			return nil, fmt.Errorf("%w: job '%s' already exists", ErrJobValidationFailed, jobName)
		}
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	if createdJob == nil {
		logger.Error("created job is nil - unexpected K8s client behavior")
		cleanupCreatedResources(ctx, client, namespace, configMapName, sourceSecretName, destSecretName)
		if err := client.DeleteModelTransferJob(ctx, namespace, jobName); err != nil && !apierrors.IsNotFound(err) {
			logger.Warn("failed to cleanup job after nil response", "jobName", jobName, "error", err)
		}
		return nil, fmt.Errorf("unexpected Kubernetes API behavior: created job object was nil")
	}

	ownerRef := metav1.OwnerReference{
		APIVersion: "batch/v1",
		Kind:       "Job",
		Name:       createdJob.Name,
		UID:        createdJob.UID,
	}

	if err := client.PatchConfigMapOwnerReference(ctx, namespace, configMapName, ownerRef); err != nil {
		logger.Warn("failed to set ownerReference on configmap", "error", err)
	}
	if sourceSecretName != "" {
		if err := client.PatchSecretOwnerReference(ctx, namespace, sourceSecretName, ownerRef); err != nil {
			logger.Warn("failed to set ownerReference on source secret", "error", err)
		}
	}
	if err := client.PatchSecretOwnerReference(ctx, namespace, destSecretName, ownerRef); err != nil {
		logger.Warn("failed to set ownerReference on destination secret", "error", err)
	}

	result := convertK8sJobToModel(createdJob)
	return &result, nil
}

func (m *ModelRegistryRepository) UpdateModelTransferJob(
	ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	oldJobName string,
	newPayload models.ModelTransferJob,
	deleteOldJob bool,
	modelRegistryID string,
) (*models.ModelTransferJob, error) {

	logger := helper.GetContextLogger(ctx)

	newJobName := newPayload.Name
	if newJobName == "" {
		return nil, fmt.Errorf("%w: new job name is required", ErrJobValidationFailed)
	}
	if len(newJobName) > 50 {
		return nil, fmt.Errorf("%w: job name must be 50 characters or less", ErrJobValidationFailed)
	}
	if errs := validation.IsDNS1123Subdomain(newJobName); len(errs) > 0 {
		return nil, fmt.Errorf("%w: invalid job name: %s", ErrJobValidationFailed, strings.Join(errs, ", "))
	}
	if newJobName == oldJobName {
		return nil, fmt.Errorf("%w: new job name must be different from old job name", ErrJobValidationFailed)
	}

	oldJob, err := client.GetModelTransferJob(ctx, namespace, oldJobName)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("%w: %s", ErrJobNotFound, oldJobName)
		}
		return nil, fmt.Errorf("failed to get existing job: %w", err)
	}

	oldAnnotations := oldJob.Annotations
	if oldAnnotations == nil {
		oldAnnotations = map[string]string{}
	}

	jobRegistry := oldJob.Labels["modelregistry.kubeflow.org/model-registry-name"]
	if jobRegistry != modelRegistryID {
		return nil, fmt.Errorf("%w: %s", ErrJobNotFound, oldJobName)
	}

	oldConfigMapName := oldAnnotations["modelregistry.kubeflow.org/configmap-name"]
	if oldConfigMapName == "" {
		return nil, fmt.Errorf("old job missing required annotation: configmap-name (job may not have been created via this API)")
	}

	oldDestSecretName := oldAnnotations["modelregistry.kubeflow.org/dest-secret"]
	if oldDestSecretName == "" {
		return nil, fmt.Errorf("old job missing required annotation: dest-secret")
	}

	oldSourceSecretName := oldAnnotations["modelregistry.kubeflow.org/source-secret"]

	if newPayload.Source.Type == "" {
		if sourceType := oldAnnotations["modelregistry.kubeflow.org/source-type"]; sourceType != "" {
			newPayload.Source.Type = models.ModelTransferJobSourceType(sourceType)
		}
	}
	if newPayload.Source.Bucket == "" {
		newPayload.Source.Bucket = oldAnnotations["modelregistry.kubeflow.org/source-bucket"]
	}
	if newPayload.Source.Key == "" {
		newPayload.Source.Key = oldAnnotations["modelregistry.kubeflow.org/source-key"]
	}
	if newPayload.Source.URI == "" {
		newPayload.Source.URI = oldAnnotations["modelregistry.kubeflow.org/source-uri"]
	}
	if newPayload.Destination.Type == "" {
		if destType := oldAnnotations["modelregistry.kubeflow.org/dest-type"]; destType != "" {
			newPayload.Destination.Type = models.ModelTransferJobDestinationType(destType)
		}
	}
	if newPayload.Destination.Registry == "" {
		newPayload.Destination.Registry = oldAnnotations["modelregistry.kubeflow.org/dest-registry"]
	}
	if newPayload.Destination.URI == "" {
		newPayload.Destination.URI = oldAnnotations["modelregistry.kubeflow.org/dest-uri"]
	}
	if newPayload.UploadIntent == "" {
		newPayload.UploadIntent = models.ModelTransferJobUploadIntent(oldAnnotations["modelregistry.kubeflow.org/upload-intent"])
	}
	if newPayload.RegisteredModelName == "" {
		newPayload.RegisteredModelName = oldAnnotations["modelregistry.kubeflow.org/model-name"]
	}
	if newPayload.ModelVersionName == "" {
		newPayload.ModelVersionName = oldAnnotations["modelregistry.kubeflow.org/version-name"]
	}
	if newPayload.RegisteredModelId == "" {
		newPayload.RegisteredModelId = oldAnnotations["modelregistry.kubeflow.org/registered-model-id"]
	}
	if newPayload.ModelVersionId == "" {
		newPayload.ModelVersionId = oldAnnotations["modelregistry.kubeflow.org/model-version-id"]
	}
	if newPayload.ModelArtifactId == "" {
		newPayload.ModelArtifactId = oldAnnotations["modelregistry.kubeflow.org/model-artifact-id"]
	}

	oldConfigMap, err := client.GetConfigMap(ctx, namespace, oldConfigMapName)
	if err != nil {
		logger.Warn("failed to get old configmap", "name", oldConfigMapName, "error", err)
	}

	var oldSourceSecret *corev1.Secret
	if oldSourceSecretName != "" {
		oldSourceSecret, err = client.GetSecret(ctx, namespace, oldSourceSecretName)
		if err != nil {
			logger.Warn("failed to get old source secret", "name", oldSourceSecretName, "error", err)
		}
	}

	oldDestSecret, err := client.GetSecret(ctx, namespace, oldDestSecretName)
	if err != nil {
		return nil, fmt.Errorf("failed to get old dest secret: %w", err)
	}

	if newPayload.Source.Type == models.ModelTransferJobSourceTypeS3 {
		if newPayload.Source.AwsAccessKeyId == "" && oldSourceSecret != nil && oldSourceSecret.Data != nil {
			if val, ok := oldSourceSecret.Data["AWS_ACCESS_KEY_ID"]; ok {
				newPayload.Source.AwsAccessKeyId = string(val)
			}

			if val, ok := oldSourceSecret.Data["AWS_SECRET_ACCESS_KEY"]; ok {
				newPayload.Source.AwsSecretAccessKey = string(val)
			}

			if newPayload.Source.Region == "" {
				if val, ok := oldSourceSecret.Data["AWS_DEFAULT_REGION"]; ok {
					newPayload.Source.Region = string(val)
				}
			}
			if newPayload.Source.Endpoint == "" {
				if val, ok := oldSourceSecret.Data["AWS_S3_ENDPOINT"]; ok {
					newPayload.Source.Endpoint = string(val)
				}
			}
		}
	}

	reuseDestCreds := (newPayload.Destination.Username == "" && newPayload.Destination.Password == "") &&
		oldDestSecret != nil && len(oldDestSecret.Data[".dockerconfigjson"]) > 0

	var existingDestSecretName string
	if reuseDestCreds {
		jobID := uuid.NewString()
		newDestSecretName := fmt.Sprintf("%s-dest-creds", newPayload.Name)
		clonedSecret := cloneDestSecretFromExisting(newDestSecretName, namespace, jobID, oldDestSecret)
		if clonedSecret == nil {
			return nil, fmt.Errorf("could not clone destination secret for reuse")
		}
		if err := client.CreateSecret(ctx, namespace, clonedSecret); err != nil {
			return nil, fmt.Errorf("failed to create cloned destination secret: %w", err)
		}
		existingDestSecretName = newDestSecretName
	}

	if newPayload.ModelArtifactName == "" && oldConfigMap != nil && oldConfigMap.Data != nil {
		if val, ok := oldConfigMap.Data["ModelArtifact.name"]; ok {
			newPayload.ModelArtifactName = val
		}
	}

	if err := validateCreatePayload(newPayload, reuseDestCreds); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	result, err := m.createModelTransferJobResources(ctx, client, namespace, newPayload, modelRegistryID, existingDestSecretName)
	if err != nil {
		if reuseDestCreds && existingDestSecretName != "" {
			if delErr := client.DeleteSecret(ctx, namespace, existingDestSecretName); delErr != nil {
				logger.Warn("failed to cleanup cloned destination secret after create failure", "name", existingDestSecretName, "error", delErr)
			}
		}
		return nil, err
	}

	if deleteOldJob {
		if err := client.DeleteModelTransferJob(ctx, namespace, oldJobName); err != nil {
			logger.Warn("failed to delete old job", "name", oldJobName, "error", err)
		}
	}
	return result, nil
}

func (m *ModelRegistryRepository) DeleteModelTransferJob(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, jobName string, modelRegistryID string) (*models.ModelTransferJob, error) {
	if modelRegistryID == "" {
		return nil, fmt.Errorf("%w: model registry name is required", ErrJobValidationFailed)
	}
	job, err := client.GetModelTransferJob(ctx, namespace, jobName)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("%w: %s", ErrJobNotFound, jobName)
		}
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	jobRegistry := job.Labels["modelregistry.kubeflow.org/model-registry-name"]
	if jobRegistry != modelRegistryID {
		return nil, fmt.Errorf("%w: %s", ErrJobNotFound, jobName)
	}

	if err := client.DeleteModelTransferJob(ctx, namespace, jobName); err != nil {
		if apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("%w: %s", ErrJobNotFound, jobName)
		}
		return nil, fmt.Errorf("failed to delete model transfer job %s: %w", jobName, err)
	}

	result := convertK8sJobToModel(job)
	return &result, nil
}

func (m *ModelRegistryRepository) getModelRegistryAddress(ctx context.Context, client k8s.KubernetesClientInterface, namespace, modelRegistryID string) (string, error) {
	modelRegistry, err := m.GetModelRegistry(ctx, client, namespace, modelRegistryID)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return "", ErrModelRegistryNotFound
		}
		return "", fmt.Errorf("failed to get model registry: %w", err)
	}
	return modelRegistry.ServerAddress, nil
}

func buildK8sJob(jobName, namespace, jobID string, payload models.ModelTransferJob,
	configMapName, sourceSecretName, destSecretName, modelRegistryAddress, modelRegistryID string) *batchv1.Job {

	backoffLimit := int32(3)

	volumes := []corev1.Volume{
		{
			Name: "destination-credentials",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: destSecretName,
					Items: []corev1.KeyToPath{
						{Key: ".dockerconfigjson", Path: ".dockerconfigjson"},
					},
				},
			},
		},
		{
			Name: "model-metadata",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: configMapName,
					},
				},
			},
		},
	}

	volumeMounts := []corev1.VolumeMount{
		{Name: "destination-credentials", MountPath: "/opt/creds/destination", ReadOnly: true},
		{Name: "model-metadata", MountPath: "/etc/model-metadata", ReadOnly: true},
	}
	registryPort, registrySecure := parseRegistryServerAddress(modelRegistryAddress)
	envVars := []corev1.EnvVar{
		{Name: "MODEL_SYNC_SOURCE_TYPE", Value: string(payload.Source.Type)},
		{Name: "MODEL_SYNC_DESTINATION_TYPE", Value: string(payload.Destination.Type)},
		{Name: "MODEL_SYNC_DESTINATION_OCI_URI", Value: payload.Destination.URI},
		{Name: "MODEL_SYNC_DESTINATION_OCI_REGISTRY", Value: payload.Destination.Registry},
		{Name: "MODEL_SYNC_DESTINATION_OCI_CREDENTIALS_PATH", Value: "/opt/creds/destination/.dockerconfigjson"},
		{Name: "MODEL_SYNC_REGISTRY_SERVER_ADDRESS", Value: registryOriginOnly(modelRegistryAddress)},
		{Name: "MODEL_SYNC_REGISTRY_PORT", Value: registryPort},
		{Name: "MODEL_SYNC_REGISTRY_IS_SECURE", Value: strconv.FormatBool(registrySecure)},
		{Name: "MODEL_SYNC_METADATA_CONFIGMAP_PATH", Value: "/etc/model-metadata"},
		{Name: "MODEL_SYNC_MODEL_UPLOAD_INTENT", Value: string(payload.UploadIntent)},
	}
	if payload.Destination.Type == models.ModelTransferJobDestinationTypeOCI && payload.Destination.Registry == "quay.io" {
		envVars = append(envVars, corev1.EnvVar{Name: "MODEL_SYNC_DESTINATION_OCI_BASE_IMAGE", Value: "quay.io/quay/busybox:latest"})
	}

	annotations := map[string]string{
		"modelregistry.kubeflow.org/source-type":         string(payload.Source.Type),
		"modelregistry.kubeflow.org/dest-type":           string(payload.Destination.Type),
		"modelregistry.kubeflow.org/dest-uri":            payload.Destination.URI,
		"modelregistry.kubeflow.org/dest-registry":       payload.Destination.Registry,
		"modelregistry.kubeflow.org/upload-intent":       string(payload.UploadIntent),
		"modelregistry.kubeflow.org/model-name":          payload.RegisteredModelName,
		"modelregistry.kubeflow.org/version-name":        payload.ModelVersionName,
		"modelregistry.kubeflow.org/artifact-name":       payload.ModelVersionName,
		"modelregistry.kubeflow.org/description":         payload.Description,
		"modelregistry.kubeflow.org/author":              payload.Author,
		"modelregistry.kubeflow.org/configmap-name":      configMapName,
		"modelregistry.kubeflow.org/dest-secret":         destSecretName,
		"modelregistry.kubeflow.org/registered-model-id": payload.RegisteredModelId,
		"modelregistry.kubeflow.org/model-version-id":    payload.ModelVersionId,
		"modelregistry.kubeflow.org/model-artifact-id":   payload.ModelArtifactId,
	}

	switch payload.Source.Type {
	case models.ModelTransferJobSourceTypeS3:
		volumes = append(volumes, corev1.Volume{
			Name: "source-credentials",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: sourceSecretName,
				},
			},
		})
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name: "source-credentials", MountPath: "/opt/creds/source", ReadOnly: true,
		})
		envVars = append(envVars,
			corev1.EnvVar{Name: "MODEL_SYNC_SOURCE_AWS_KEY", Value: payload.Source.Key},
			corev1.EnvVar{Name: "MODEL_SYNC_SOURCE_S3_CREDENTIALS_PATH", Value: "/opt/creds/source"},
		)
		annotations["modelregistry.kubeflow.org/source-bucket"] = payload.Source.Bucket
		annotations["modelregistry.kubeflow.org/source-key"] = payload.Source.Key
		annotations["modelregistry.kubeflow.org/source-secret"] = sourceSecretName
		annotations["modelregistry.kubeflow.org/model-registry-name"] = modelRegistryID

	case models.ModelTransferJobSourceTypeURI:
		envVars = append(envVars,
			corev1.EnvVar{Name: "MODEL_SYNC_SOURCE_URI", Value: payload.Source.URI},
		)
		annotations["modelregistry.kubeflow.org/source-uri"] = payload.Source.URI
		annotations["modelregistry.kubeflow.org/model-registry-name"] = modelRegistryID
	}

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
			Labels: map[string]string{
				"modelregistry.kubeflow.org/job-type":            "async-upload",
				"modelregistry.kubeflow.org/job-id":              jobID,
				"modelregistry.kubeflow.org/model-registry-name": modelRegistryID,
			},
			Annotations: annotations,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: &backoffLimit,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"modelregistry.kubeflow.org/job-type": "async-upload",
					},
				},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot: func() *bool { b := true; return &b }(),
					},
					Volumes: volumes,
					Containers: []corev1.Container{
						{
							Name:            "async-upload",
							Image:           "quay.io/opendatahub/model-registry-job-async-upload:latest",
							ImagePullPolicy: corev1.PullIfNotPresent,
							VolumeMounts:    volumeMounts,
							Env:             envVars,
							SecurityContext: &corev1.SecurityContext{
								AllowPrivilegeEscalation: func() *bool { b := false; return &b }(),
								Capabilities: &corev1.Capabilities{
									Drop: []corev1.Capability{"ALL"},
								},
								SeccompProfile: &corev1.SeccompProfile{
									Type: corev1.SeccompProfileTypeRuntimeDefault,
								},
							},
						},
					},
				},
			},
		},
	}
}

func convertK8sJobToModel(job *batchv1.Job) models.ModelTransferJob {
	if job == nil {
		return models.ModelTransferJob{}
	}

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

	errorMessage := ""
	for _, condition := range job.Status.Conditions {
		if condition.Type == batchv1.JobFailed && condition.Status == corev1.ConditionTrue {
			errorMessage = condition.Message
			break
		}
	}

	lastUpdateTime := fmt.Sprintf("%d", job.CreationTimestamp.UnixMilli())
	if job.Status.CompletionTime != nil {
		lastUpdateTime = fmt.Sprintf("%d", job.Status.CompletionTime.UnixMilli())
	}

	return models.ModelTransferJob{
		Id:          labels["modelregistry.kubeflow.org/job-id"],
		Name:        job.Name,
		Description: annotations["modelregistry.kubeflow.org/description"],
		Source: models.ModelTransferJobSource{
			Type:   models.ModelTransferJobSourceType(annotations["modelregistry.kubeflow.org/source-type"]),
			Bucket: annotations["modelregistry.kubeflow.org/source-bucket"],
			Key:    annotations["modelregistry.kubeflow.org/source-key"],
			URI:    annotations["modelregistry.kubeflow.org/source-uri"],
		},
		Destination: models.ModelTransferJobDestination{
			Type:     models.ModelTransferJobDestinationType(annotations["modelregistry.kubeflow.org/dest-type"]),
			URI:      annotations["modelregistry.kubeflow.org/dest-uri"],
			Registry: annotations["modelregistry.kubeflow.org/dest-registry"],
		},
		UploadIntent:             models.ModelTransferJobUploadIntent(annotations["modelregistry.kubeflow.org/upload-intent"]),
		RegisteredModelName:      annotations["modelregistry.kubeflow.org/model-name"],
		ModelVersionName:         annotations["modelregistry.kubeflow.org/version-name"],
		ModelArtifactName:        annotations["modelregistry.kubeflow.org/version-name"],
		RegisteredModelId:        annotations["modelregistry.kubeflow.org/registered-model-id"],
		ModelVersionId:           annotations["modelregistry.kubeflow.org/model-version-id"],
		ModelArtifactId:          annotations["modelregistry.kubeflow.org/model-artifact-id"],
		Author:                   annotations["modelregistry.kubeflow.org/author"],
		Status:                   status,
		ErrorMessage:             errorMessage,
		CreateTimeSinceEpoch:     fmt.Sprintf("%d", job.CreationTimestamp.UnixMilli()),
		LastUpdateTimeSinceEpoch: lastUpdateTime,
		Namespace:                job.Namespace,
		SourceSecretName:         annotations["modelregistry.kubeflow.org/source-secret"],
		DestSecretName:           annotations["modelregistry.kubeflow.org/dest-secret"],
	}
}

func buildModelMetadataConfigMap(name, namespace string, payload models.ModelTransferJob, jobID string, jobName string) *corev1.ConfigMap {
	data := map[string]string{
		"ModelVersion.name":   payload.ModelVersionName,
		"ModelVersion.author": payload.Author,
		"ModelArtifact.name":  payload.ModelVersionName,
	}

	if payload.VersionDescription != "" {
		data["ModelVersion.description"] = payload.VersionDescription
		data["ModelArtifact.description"] = payload.VersionDescription
	}
	if payload.SourceModelFormat != "" {
		data["ModelArtifact.modelFormatName"] = payload.SourceModelFormat
	}
	if payload.SourceModelFormatVersion != "" {
		data["ModelArtifact.modelFormatVersion"] = payload.SourceModelFormatVersion
	}
	if len(payload.ModelCustomProperties) > 0 {
		if b, err := json.Marshal(payload.ModelCustomProperties); err == nil {
			data["RegisteredModel.customProperties"] = string(b)
		}
	}
	if len(payload.VersionCustomProperties) > 0 {
		if b, err := json.Marshal(payload.VersionCustomProperties); err == nil {
			data["ModelVersion.customProperties"] = string(b)
			data["ModelArtifact.customProperties"] = string(b)
		}
	}

	switch payload.UploadIntent {
	case models.ModelTransferJobUploadIntentCreateModel:
		data["RegisteredModel.name"] = payload.RegisteredModelName
		data["RegisteredModel.description"] = payload.Description
		data["RegisteredModel.owner"] = payload.Author
	case models.ModelTransferJobUploadIntentCreateVersion:
		data["RegisteredModel.id"] = payload.RegisteredModelId
	case models.ModelTransferJobUploadIntentUpdateArtifact:
		data["RegisteredModel.id"] = payload.RegisteredModelId
		data["ModelVersion.id"] = payload.ModelVersionId
		data["ModelArtifact.id"] = payload.ModelArtifactId
	}

	data["ModelArtifact.model_source_kind"] = "Job"
	data["ModelArtifact.model_source_class"] = "async-upload"
	data["ModelArtifact.model_source_group"] = "batch/v1"
	data["ModelArtifact.model_source_name"] = jobName

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"modelregistry.kubeflow.org/job-type": "async-upload",
				"modelregistry.kubeflow.org/job-id":   jobID,
			},
		},
		Data: data,
	}
}

func buildSourceSecret(name, namespace string, payload models.ModelTransferJob, jobID string) *corev1.Secret {
	stringData := map[string]string{
		"AWS_ACCESS_KEY_ID":     payload.Source.AwsAccessKeyId,
		"AWS_SECRET_ACCESS_KEY": payload.Source.AwsSecretAccessKey,
		"AWS_S3_ENDPOINT":       payload.Source.Endpoint,
		"AWS_S3_BUCKET":         payload.Source.Bucket,
	}

	if payload.Source.Region != "" {
		stringData["AWS_DEFAULT_REGION"] = payload.Source.Region
	}
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"modelregistry.kubeflow.org/job-type": "async-upload",
				"modelregistry.kubeflow.org/job-id":   jobID,
			},
		},
		Type:       corev1.SecretTypeOpaque,
		StringData: stringData,
	}
}

func buildDestinationSecret(name, namespace string, payload models.ModelTransferJob, jobID string) (*corev1.Secret, error) {
	// NOTE: Due to async-upload bug, auth is NOT base64 encoded here
	auth := fmt.Sprintf("%s:%s", payload.Destination.Username, payload.Destination.Password)

	registry := payload.Destination.Registry
	if registry == "" {
		registry, _ = extractRegistryFromURI(payload.Destination.URI)
	}

	dockerConfig := fmt.Sprintf(`{"auths":{"%s":{"auth":"%s","email":"%s"}}}`,
		registry, auth, payload.Destination.Email)

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"modelregistry.kubeflow.org/job-type": "async-upload",
				"modelregistry.kubeflow.org/job-id":   jobID,
			},
		},
		Type: corev1.SecretTypeDockerConfigJson,
		StringData: map[string]string{
			".dockerconfigjson": dockerConfig,
		},
	}, nil
}

func validateCreatePayload(payload models.ModelTransferJob, skipDestCredsValidation bool) error {
	if payload.Name == "" {
		return fmt.Errorf("%w: job name is required", ErrJobValidationFailed)
	}
	if len(payload.Name) > 50 {
		return fmt.Errorf("%w: job name must be 50 characters or less", ErrJobValidationFailed)
	}

	if errs := validation.IsDNS1123Subdomain(payload.Name); len(errs) > 0 {
		return fmt.Errorf("%w: %s", ErrJobValidationFailed, strings.Join(errs, ", "))
	}

	if payload.Source.Type == "" {
		return fmt.Errorf("%w: source type is required", ErrJobValidationFailed)
	}

	switch payload.Source.Type {
	case models.ModelTransferJobSourceTypeS3:
		if payload.Source.Bucket == "" {
			return fmt.Errorf("%w: source bucket is required for S3 source type", ErrJobValidationFailed)
		}
		if payload.Source.Key == "" {
			return fmt.Errorf("%w: source key is required for S3 source type", ErrJobValidationFailed)
		}
		if payload.Source.AwsAccessKeyId == "" {
			return fmt.Errorf("%w: AWS access key ID is required for S3 source type", ErrJobValidationFailed)
		}
		if payload.Source.AwsSecretAccessKey == "" {
			return fmt.Errorf("%w: AWS secret access key is required for S3 source type", ErrJobValidationFailed)
		}
	case models.ModelTransferJobSourceTypeURI:
		if payload.Source.URI == "" {
			return fmt.Errorf("%w: source URI is required for URI source type", ErrJobValidationFailed)
		}
	default:
		return fmt.Errorf("%w: invalid source type: %s", ErrJobValidationFailed, payload.Source.Type)
	}

	if payload.Destination.Type == "" {
		return fmt.Errorf("%w: destination type is required", ErrJobValidationFailed)
	}
	if payload.Destination.Type == models.ModelTransferJobDestinationTypeOCI {
		if payload.Destination.URI == "" {
			return fmt.Errorf("%w: destination URI is required for OCI destination type", ErrJobValidationFailed)
		}

		registry := payload.Destination.Registry
		if registry == "" {
			if _, err := extractRegistryFromURI(payload.Destination.URI); err != nil {
				return fmt.Errorf("%w: cannot extract registry from destination URI: %w", ErrJobValidationFailed, err)
			}
		}

		if !skipDestCredsValidation {
			if payload.Destination.Username == "" {
				return fmt.Errorf("%w: destination username is required for OCI destination type", ErrJobValidationFailed)
			}
			if payload.Destination.Password == "" {
				return fmt.Errorf("%w: destination password is required for OCI destination type", ErrJobValidationFailed)
			}
		}
	} else {
		return fmt.Errorf("%w: invalid destination type: %s", ErrJobValidationFailed, payload.Destination.Type)
	}

	if payload.UploadIntent == "" {
		return fmt.Errorf("%w: upload intent is required", ErrJobValidationFailed)
	}
	validIntents := map[models.ModelTransferJobUploadIntent]bool{
		models.ModelTransferJobUploadIntentCreateModel:    true,
		models.ModelTransferJobUploadIntentCreateVersion:  true,
		models.ModelTransferJobUploadIntentUpdateArtifact: true,
	}
	if !validIntents[payload.UploadIntent] {
		return fmt.Errorf("%w: invalid upload intent: %s", ErrJobValidationFailed, payload.UploadIntent)
	}

	switch payload.UploadIntent {
	case models.ModelTransferJobUploadIntentCreateModel:
		if payload.RegisteredModelName == "" {
			return fmt.Errorf("%w: registered model name is required for create_model intent", ErrJobValidationFailed)
		}
		if payload.ModelVersionName == "" {
			return fmt.Errorf("%w: model version name is required for create_model intent", ErrJobValidationFailed)
		}
	case models.ModelTransferJobUploadIntentCreateVersion:
		if payload.RegisteredModelId == "" {
			return fmt.Errorf("%w: registered model ID is required for create_version intent", ErrJobValidationFailed)
		}
		if payload.ModelVersionName == "" {
			return fmt.Errorf("%w: model version name is required for create_version intent", ErrJobValidationFailed)
		}
	case models.ModelTransferJobUploadIntentUpdateArtifact:
		if payload.ModelArtifactId == "" {
			return fmt.Errorf("%w: model artifact ID is required for update_artifact intent", ErrJobValidationFailed)
		}
	}

	return nil
}

func cleanupCreatedResources(ctx context.Context, client k8s.KubernetesClientInterface, namespace, configMapName, sourceSecretName, destSecretName string) {
	logger := helper.GetContextLogger(ctx)

	if configMapName != "" {
		if err := client.DeleteConfigMap(ctx, namespace, configMapName); err != nil {
			logger.Warn("failed to cleanup configmap", "name", configMapName, "error", err)
		}
	}
	if sourceSecretName != "" {
		if err := client.DeleteSecret(ctx, namespace, sourceSecretName); err != nil {
			logger.Warn("failed to cleanup source secret", "name", sourceSecretName, "error", err)
		}
	}
	if destSecretName != "" {
		if err := client.DeleteSecret(ctx, namespace, destSecretName); err != nil {
			logger.Warn("failed to cleanup dest secret", "name", destSecretName, "error", err)
		}
	}
}

func extractRegistryFromURI(uri string) (string, error) {
	parts := strings.Split(uri, "/")
	if len(parts) > 0 && parts[0] != "" {
		return parts[0], nil
	}
	return "", fmt.Errorf("cannot extract registry from URI: %s", uri)
}

func parseRegistryServerAddress(serverAddress string) (port string, isSecure bool) {
	u, err := url.Parse(serverAddress)
	if err != nil {
		return "8080", false
	}
	port = u.Port()
	if port == "" {
		if strings.ToLower(u.Scheme) == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}
	isSecure = strings.ToLower(u.Scheme) == "https"
	return port, isSecure
}

// registryOriginOnly returns only scheme + host (no port, no path). The async-upload image
// adds the port from MODEL_SYNC_REGISTRY_PORT and the API path itself; including port here
// would produce host:8080:8080 and break.
func registryOriginOnly(serverAddress string) string {
	u, err := url.Parse(serverAddress)
	if err != nil {
		return serverAddress
	}
	host := u.Hostname()
	if host == "" {
		return serverAddress
	}
	scheme := u.Scheme
	if scheme == "" {
		scheme = "http"
	}
	return scheme + "://" + host
}

func cloneDestSecretFromExisting(newName, namespace, jobID string, oldSecret *corev1.Secret) *corev1.Secret {
	if oldSecret == nil {
		return nil
	}
	data := make(map[string][]byte)
	for k, v := range oldSecret.Data {
		data[k] = append([]byte(nil), v...)
	}
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      newName,
			Namespace: namespace,
			Labels: map[string]string{
				"modelregistry.kubeflow.org/job-type": "async-upload",
				"modelregistry.kubeflow.org/job-id":   jobID,
			},
		},
		Type: corev1.SecretTypeDockerConfigJson,
		Data: data,
	}
}

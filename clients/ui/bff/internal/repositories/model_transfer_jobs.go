package repositories

import (
	"context"
	"errors"
	"fmt"
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
	ErrJobNotFound          = errors.New("model transfer job not found")
	ErrJobValidationFailed  = errors.New("validation failed")
	ErrJobNameRequired      = errors.New("job name is required")
	ErrJobNameTooLong       = errors.New("job name exceeds maximum length")
	ErrJobNameInvalid       = errors.New("job name contains invalid characters")
	ErrSourceTypeRequired   = errors.New("source type is required")
	ErrSourceTypeInvalid    = errors.New("invalid source type")
	ErrDestTypeRequired     = errors.New("destination type is required")
	ErrDestTypeInvalid      = errors.New("invalid destination type")
	ErrUploadIntentRequired = errors.New("upload intent is required")
	ErrUploadIntentInvalid  = errors.New("invalid upload intent")
)

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

func (m *ModelRegistryRepository) CreateModelTransferJob(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, payload models.ModelTransferJob, modelRegistryID string) (*models.ModelTransferJob, error) {
	if err := validateCreatePayload(payload); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	logger := helper.GetContextLogger(ctx)

	modelRegistryAddress, err := m.getModelRegistryAddress(ctx, client, namespace, modelRegistryID)
	if err != nil {
		return nil, err
	}

	jobId := uuid.NewString()
	jobName := payload.Name

	configMapName := fmt.Sprintf("%s-metadata", jobName)
	destSecretName := fmt.Sprintf("%s-dest-creds", jobName)

	configMap := buildModelMetadataConfigMap(configMapName, namespace, payload, jobId)
	if err := client.CreateConfigMap(ctx, namespace, configMap); err != nil {
		return nil, fmt.Errorf("failed to create metadata configmap: %w", err)
	}

	var sourceSecretName string
	if payload.Source.Type == models.ModelTransferJobSourceTypeS3 {
		sourceSecretName = fmt.Sprintf("%s-source-creds", jobName)
		sourceSecret := buildSourceSecret(sourceSecretName, namespace, payload, jobId)
		if err := client.CreateSecret(ctx, namespace, sourceSecret); err != nil {
			if delErr := client.DeleteConfigMap(ctx, namespace, configMapName); delErr != nil {
				logger.Error("failed to cleanup configmap after error", "name", configMapName, "error", delErr)
			}
			return nil, fmt.Errorf("failed to create source secret: %w", err)
		}

	}

	destSecret, err := buildDestinationSecret(destSecretName, namespace, payload, jobId)
	if err := client.CreateSecret(ctx, namespace, destSecret); err != nil {

		client.DeleteConfigMap(ctx, namespace, configMapName)

		if sourceSecretName != "" {
			client.DeleteSecret(ctx, namespace, sourceSecretName)
		}
		return nil, fmt.Errorf("failed to create destination secret: %w", err)
	}

	job := buildFullK8sJob(jobName, namespace, jobId, payload, configMapName, sourceSecretName, destSecretName, modelRegistryAddress)
	if err := client.CreateModelTransferJob(ctx, namespace, job); err != nil {

		client.DeleteConfigMap(ctx, namespace, configMapName)
		if sourceSecretName != "" {
			client.DeleteSecret(ctx, namespace, sourceSecretName)
		}
		client.DeleteSecret(ctx, namespace, destSecretName)
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	createdJob, err := client.GetModelTransferJob(ctx, namespace, jobName)
	if err != nil {
		logger.Warn("job created but failed to retrieve", "error", err)
		return &models.ModelTransferJob{
			Id:   jobId,
			Name: jobName,
		}, nil
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
	if newPayload.Name == "" {
		return nil, fmt.Errorf("%w: new job name is required", ErrJobValidationFailed)
	}
	if len(newPayload.Name) > 50 {
		return nil, fmt.Errorf("%w: job name must be 50 characters or less", ErrJobValidationFailed)
	}
	if errs := validation.IsDNS1123Subdomain(newPayload.Name); len(errs) > 0 {
		return nil, fmt.Errorf("%w: %s", ErrJobValidationFailed, strings.Join(errs, ", "))
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

	oldConfigMapName := oldAnnotations["modelregistry.kubeflow.org/configmap-name"]
	if oldConfigMapName == "" {
		return nil, fmt.Errorf("old job is missing required annotation: modelregistry.kubeflow.org/configmap-name")
	}

	oldDestSecretName := oldAnnotations["modelregistry.kubeflow.org/dest-secret"]
	if oldDestSecretName == "" {
		return nil, fmt.Errorf("old job is missing required annotation: modelregistry.kubeflow.org/dest-secret")
	}

	oldSourceSecretName := oldAnnotations["modelregistry.kubeflow.org/source-secret"]

	if deleteOldJob {
		oldSourceSecret, _ := client.GetSecret(ctx, namespace, oldSourceSecretName)
		oldDestSecret, _ := client.GetSecret(ctx, namespace, oldDestSecretName)
		oldConfigMap, _ := client.GetConfigMap(ctx, namespace, oldConfigMapName)

		if err := client.DeleteModelTransferJob(ctx, namespace, oldJobName); err != nil {
			logger.Warn("failed to delete old job", "name", oldJobName, "error", err)
		}
		if oldConfigMapName != "" {
			if err := client.DeleteConfigMap(ctx, namespace, oldConfigMapName); err != nil {
				logger.Warn("failed to delete old configmap", "name", oldConfigMapName, "error", err)
			}
		}
		if oldSourceSecretName != "" {
			if err := client.DeleteSecret(ctx, namespace, oldSourceSecretName); err != nil {
				logger.Warn("failed to delete old source secret", "name", oldSourceSecretName, "error", err)
			}
		}
		if oldDestSecretName != "" {
			if err := client.DeleteSecret(ctx, namespace, oldDestSecretName); err != nil {
				logger.Warn("failed to delete old dest secret", "name", oldDestSecretName, "error", err)
			}
		}

		mergedPayload := mergePayloadWithOldData(newPayload, oldJob, oldSourceSecret, oldDestSecret, oldConfigMap)

		return m.CreateModelTransferJob(ctx, client, namespace, mergedPayload, modelRegistryID)
	} else {

		newJobName := newPayload.Name

		if newJobName == oldJobName {
			return nil, fmt.Errorf("new job name must be different from old job name (%s)", oldJobName)
		}

		newJobId := uuid.NewString()

		modelRegistryAddress, err := m.getModelRegistryAddress(ctx, client, namespace, modelRegistryID)
		if err != nil {
			return nil, err
		}

		job := buildFullK8sJob(newJobName, namespace, newJobId, newPayload,
			oldConfigMapName, oldSourceSecretName, oldDestSecretName, modelRegistryAddress)

		if err := client.CreateModelTransferJob(ctx, namespace, job); err != nil {
			return nil, fmt.Errorf("failed to create new job: %w", err)
		}

		createdJob, err := client.GetModelTransferJob(ctx, namespace, newJobName)
		if err != nil {
			return nil, fmt.Errorf("failed to get created job: %w", err)
		}

		result := convertK8sJobToModel(createdJob)
		return &result, nil
	}
}

func (m *ModelRegistryRepository) DeleteModelTransferJob(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, jobName string) (*models.ModelTransferJob, error) {
	var errs []string

	job, err := client.GetModelTransferJob(ctx, namespace, jobName)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("%w: %s", ErrJobNotFound, jobName)
		}
		errs = append(errs, fmt.Sprintf("could not get job for cleanup: %v", err))
	}

	var configMapName, sourceSecretName, destSecretName string
	if job != nil && job.Annotations != nil {
		configMapName = job.Annotations["modelregistry.kubeflow.org/configmap-name"]
		sourceSecretName = job.Annotations["modelregistry.kubeflow.org/source-secret"]
		destSecretName = job.Annotations["modelregistry.kubeflow.org/dest-secret"]
	}

	if err := client.DeleteModelTransferJob(ctx, namespace, jobName); err != nil {
		return nil, fmt.Errorf("failed to delete model transfer job %s: %w", jobName, err)
	}

	if configMapName != "" {
		if err := client.DeleteConfigMap(ctx, namespace, configMapName); err != nil {
			errs = append(errs, fmt.Sprintf("failed to delete configmap %s: %v", configMapName, err))
		}
	}
	if sourceSecretName != "" {
		if err := client.DeleteSecret(ctx, namespace, sourceSecretName); err != nil {
			errs = append(errs, fmt.Sprintf("failed to delete source secret %s: %v", sourceSecretName, err))
		}
	}
	if destSecretName != "" {
		if err := client.DeleteSecret(ctx, namespace, destSecretName); err != nil {
			errs = append(errs, fmt.Sprintf("failed to delete destination secret %s: %v", destSecretName, err))
		}
	}

	if len(errs) > 0 {
		result := convertK8sJobToModel(job)

		return &result, fmt.Errorf("job deleted but cleanup had errors: %s", strings.Join(errs, "; "))
	}

	result := convertK8sJobToModel(job)
	return &result, nil
}

func (m *ModelRegistryRepository) getModelRegistryAddress(ctx context.Context, client k8s.KubernetesClientInterface, namespace, modelRegistryID string) (string, error) {
	modelRegistry, err := m.GetModelRegistry(ctx, client, namespace, modelRegistryID)
	if err != nil {
		return "", fmt.Errorf("failed to get model registry: %w", err)
	}
	return modelRegistry.ServerAddress, nil
}

func buildFullK8sJob(jobName, namespace, jobId string, payload models.ModelTransferJob,
	configMapName, sourceSecretName, destSecretName, modelRegistryAddress string) *batchv1.Job {

	backoffLimit := int32(3)
	baseImage := models.DefaultOCIBaseImage

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
	envVars := []corev1.EnvVar{
		{Name: "MODEL_SYNC_SOURCE_TYPE", Value: string(payload.Source.Type)},
		{Name: "MODEL_SYNC_DESTINATION_TYPE", Value: string(payload.Destination.Type)},
		{Name: "MODEL_SYNC_DESTINATION_OCI_URI", Value: payload.Destination.URI},
		{Name: "MODEL_SYNC_DESTINATION_OCI_REGISTRY", Value: payload.Destination.Registry},
		{Name: "MODEL_SYNC_DESTINATION_OCI_BASE_IMAGE", Value: baseImage},
		{Name: "MODEL_SYNC_DESTINATION_OCI_CREDENTIALS_PATH", Value: "/opt/creds/destination/.dockerconfigjson"},
		{Name: "MODEL_SYNC_REGISTRY_SERVER_ADDRESS", Value: modelRegistryAddress},
		{Name: "MODEL_SYNC_METADATA_CONFIGMAP_PATH", Value: "/etc/model-metadata"},
		{Name: "MODEL_SYNC_MODEL_UPLOAD_INTENT", Value: string(payload.UploadIntent)},
	}

	annotations := map[string]string{
		"modelregistry.kubeflow.org/source-type":    string(payload.Source.Type),
		"modelregistry.kubeflow.org/dest-type":      string(payload.Destination.Type),
		"modelregistry.kubeflow.org/dest-uri":       payload.Destination.URI,
		"modelregistry.kubeflow.org/dest-registry":  payload.Destination.Registry,
		"modelregistry.kubeflow.org/upload-intent":  string(payload.UploadIntent),
		"modelregistry.kubeflow.org/model-name":     payload.RegisteredModelName,
		"modelregistry.kubeflow.org/version-name":   payload.ModelVersionName,
		"modelregistry.kubeflow.org/description":    payload.Description,
		"modelregistry.kubeflow.org/author":         payload.Author,
		"modelregistry.kubeflow.org/configmap-name": configMapName,
		"modelregistry.kubeflow.org/dest-secret":    destSecretName,
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

	case models.ModelTransferJobSourceTypeURI:
		envVars = append(envVars,
			corev1.EnvVar{Name: "MODEL_SYNC_SOURCE_URI", Value: payload.Source.URI},
		)
		annotations["modelregistry.kubeflow.org/source-uri"] = payload.Source.URI
	}

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
			Labels: map[string]string{
				"modelregistry.kubeflow.org/job-type": "async-upload",
				"modelregistry.kubeflow.org/job-id":   jobId,
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
		RegisteredModelId:        annotations["modelregistry.kubeflow.org/registered-model-id"],
		ModelVersionId:           annotations["modelregistry.kubeflow.org/model-version-id"],
		ModelArtifactId:          annotations["modelregistry.kubeflow.org/model-artifact-id"],
		Author:                   annotations["modelregistry.kubeflow.org/author"],
		Status:                   status,
		ErrorMessage:             errorMessage,
		CreateTimeSinceEpoch:     fmt.Sprintf("%d", job.CreationTimestamp.UnixMilli()),
		LastUpdateTimeSinceEpoch: lastUpdateTime,
		Namespace:                job.Namespace,
	}
}

func buildModelMetadataConfigMap(name, namespace string, payload models.ModelTransferJob, jobId string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"modelregistry.kubeflow.org/job-type": "async-upload",
				"modelregistry.kubeflow.org/job-id":   jobId,
			},
		},
		Data: map[string]string{
			"RegisteredModel.name":        payload.RegisteredModelName,
			"RegisteredModel.description": payload.Description,
			"RegisteredModel.owner":       payload.Author,
			"ModelVersion.name":           payload.ModelVersionName,
			"ModelVersion.author":         payload.Author,
			"ModelArtifact.name":          payload.ModelArtifactName,
		},
	}
}

func buildSourceSecret(name, namespace string, payload models.ModelTransferJob, jobId string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"modelregistry.kubeflow.org/job-type": "async-upload",
				"modelregistry.kubeflow.org/job-id":   jobId,
			},
		},
		Type: corev1.SecretTypeOpaque,
		StringData: map[string]string{
			"AWS_ACCESS_KEY_ID":     payload.Source.AwsAccessKeyId,
			"AWS_SECRET_ACCESS_KEY": payload.Source.AwsSecretAccessKey,
			"AWS_DEFAULT_REGION":    payload.Source.Region,
			"AWS_S3_ENDPOINT":       payload.Source.Endpoint,
			"AWS_S3_BUCKET":         payload.Source.Bucket,
		},
	}
}

func buildDestinationSecret(name, namespace string, payload models.ModelTransferJob, jobId string) (*corev1.Secret, error) {
	// Build docker config JSON for OCI authentication
	// NOTE: Due to async-upload bug, auth is NOT base64 encoded here
	auth := fmt.Sprintf("%s:%s", payload.Destination.Username, payload.Destination.Password)

	registry := payload.Destination.Registry
	if registry == "" {
		parts := strings.Split(payload.Destination.URI, "/")
		if len(parts) > 0 && parts[0] != "" {
			registry = parts[0]
		} else {
			return nil, fmt.Errorf("cannot determine registry from destination URI: %s", payload.Destination.URI)
		}
	}

	dockerConfig := fmt.Sprintf(`{"auths":{"%s":{"auth":"%s","email":"%s"}}}`,
		registry, auth, payload.Destination.Email)

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"modelregistry.kubeflow.org/job-type": "async-upload",
				"modelregistry.kubeflow.org/job-id":   jobId,
			},
		},
		Type: corev1.SecretTypeDockerConfigJson,
		StringData: map[string]string{
			".dockerconfigjson": dockerConfig,
		},
	}, nil
}

func mergePayloadWithOldData(newPayload models.ModelTransferJob,
	oldJob *batchv1.Job,
	oldSourceSecret *corev1.Secret, oldDestSecret *corev1.Secret,
	oldConfigMap *corev1.ConfigMap) models.ModelTransferJob {

	oldAnnotations := map[string]string{}
	if oldJob != nil && oldJob.Annotations != nil {
		oldAnnotations = oldJob.Annotations
	}

	if newPayload.Source.Type == "" {
		newPayload.Source.Type = models.ModelTransferJobSourceType(
			oldAnnotations["modelregistry.kubeflow.org/source-type"])
	}

	if newPayload.Source.Key == "" {
		newPayload.Source.Key = oldAnnotations["modelregistry.kubeflow.org/source-key"]
	}

	if newPayload.Source.URI == "" {
		newPayload.Source.URI = oldAnnotations["modelregistry.kubeflow.org/source-uri"]
	}

	if newPayload.Destination.URI == "" {
		newPayload.Destination.URI = oldAnnotations["modelregistry.kubeflow.org/dest-uri"]
	}

	if newPayload.Destination.Type == "" {
		newPayload.Destination.Type = models.ModelTransferJobDestinationType(
			oldAnnotations["modelregistry.kubeflow.org/dest-type"])
	}

	if newPayload.UploadIntent == "" {
		newPayload.UploadIntent = models.ModelTransferJobUploadIntent(
			oldAnnotations["modelregistry.kubeflow.org/upload-intent"])
	}

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
		if newPayload.Source.Bucket == "" {
			if val, ok := oldSourceSecret.Data["AWS_S3_BUCKET"]; ok {
				newPayload.Source.Bucket = string(val)
			}
		}
	}

	if newPayload.Destination.Username == "" && oldDestSecret != nil && oldDestSecret.Data != nil {
		if val, ok := oldDestSecret.Data["username"]; ok {
			newPayload.Destination.Username = string(val)
		}
		if val, ok := oldDestSecret.Data["password"]; ok {
			newPayload.Destination.Password = string(val)
		}
		if val, ok := oldDestSecret.Data["email"]; ok {
			newPayload.Destination.Email = string(val)
		}
		if newPayload.Destination.Registry == "" {
			if val, ok := oldDestSecret.Data["registry"]; ok {
				newPayload.Destination.Registry = string(val)
			}
		}
	}

	if oldConfigMap != nil && oldConfigMap.Data != nil {
		if newPayload.RegisteredModelName == "" {
			if val, ok := oldConfigMap.Data["RegisteredModel.name"]; ok {
				newPayload.RegisteredModelName = val
			}
		}
		if newPayload.ModelVersionName == "" {
			if val, ok := oldConfigMap.Data["ModelVersion.name"]; ok {
				newPayload.ModelVersionName = val
			}
		}
		if newPayload.Description == "" {
			if val, ok := oldConfigMap.Data["RegisteredModel.description"]; ok {
				newPayload.Description = val
			}
		}
		if newPayload.Author == "" {
			if val, ok := oldConfigMap.Data["ModelVersion.author"]; ok {
				newPayload.Author = val
			}
		}
	}

	return newPayload
}

func validateCreatePayload(payload models.ModelTransferJob) error {
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
			parts := strings.Split(payload.Destination.URI, "/")
			if len(parts) == 0 || parts[0] == "" {
				return fmt.Errorf("invalid destination URI: cannot extract registry from %s", payload.Destination.URI)
			}
		}
		if payload.Destination.Username == "" {
			return fmt.Errorf("%w: destination username is required for OCI destination type", ErrJobValidationFailed)
		}
		if payload.Destination.Password == "" {
			return fmt.Errorf("%w: destination password is required for OCI destination type", ErrJobValidationFailed)
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

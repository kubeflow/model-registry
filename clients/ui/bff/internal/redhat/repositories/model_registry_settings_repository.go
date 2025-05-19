package repositories

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	helper "github.com/kubeflow/model-registry/ui/bff/internal/helpers"
	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	databasePasswordKey = "database-password"
	caCertificateKey    = "ca.crt"
)

var modelRegistryGVR = schema.GroupVersionResource{
	Group:    "modelregistry.opendatahub.io",
	Version:  "v1beta1",
	Resource: "modelregistries",
}

// ModelRegistrySettingsRepository contains the Kubernetes-aware implementation of the Model Registry settings flows.
// TODO(upstream): Relocate this repository upstream once the core CRUD logic is accepted for all distros.
type ModelRegistrySettingsRepository struct {
	logger *slog.Logger
}

func NewModelRegistrySettingsRepository(logger *slog.Logger) *ModelRegistrySettingsRepository {
	return &ModelRegistrySettingsRepository{logger: logger}
}

type kubeClients struct {
	dynamic dynamic.Interface
	core    kubernetes.Interface
}

// TODO(upstream): The generic client construction belongs in the shared Kubernetes factory helpers.
func (r *ModelRegistrySettingsRepository) buildClients(client k8s.KubernetesClientInterface) (*kubeClients, error) {
	cfg, err := restConfigForClient(client)
	if err != nil {
		return nil, err
	}

	dyn, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic Kubernetes client: %w", err)
	}

	core, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create typed Kubernetes client: %w", err)
	}

	return &kubeClients{dynamic: dyn, core: core}, nil
}

type restConfigProvider interface {
	RESTConfig() *rest.Config
}

// TODO(upstream): Fold this helper into the upstream factory so downstream code does not duplicate it.
func restConfigForClient(client k8s.KubernetesClientInterface) (*rest.Config, error) {
	if provider, ok := client.(restConfigProvider); ok {
		if cfg := provider.RESTConfig(); cfg != nil {
			return rest.CopyConfig(cfg), nil
		}
	}

	cfg, err := helper.GetKubeconfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}
	return rest.CopyConfig(cfg), nil
}

// TODO(upstream): Move list/get/create/patch/delete implementations upstream once API contracts settle.
func (r *ModelRegistrySettingsRepository) List(
	ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	labelSelector string,
) ([]models.ModelRegistryKind, error) {
	if namespace == "" {
		return nil, errors.New("namespace is required")
	}

	clients, err := r.buildClients(client)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{}
	if labelSelector != "" {
		opts.LabelSelector = labelSelector
	}

	list, err := clients.dynamic.Resource(modelRegistryGVR).Namespace(namespace).List(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list ModelRegistry resources: %w", err)
	}

	registries := make([]models.ModelRegistryKind, 0, len(list.Items))
	for _, item := range list.Items {
		typed, err := convertUnstructuredToModelRegistryKind(item, false, ctx, clients.core)
		if err != nil {
			return nil, err
		}
		registries = append(registries, typed)
	}

	return registries, nil
}

func (r *ModelRegistrySettingsRepository) Get(
	ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	name string,
) (models.ModelRegistryKind, *string, error) {
	if namespace == "" || name == "" {
		return models.ModelRegistryKind{}, nil, errors.New("namespace and name are required")
	}

	clients, err := r.buildClients(client)
	if err != nil {
		return models.ModelRegistryKind{}, nil, err
	}

	item, err := clients.dynamic.Resource(modelRegistryGVR).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return models.ModelRegistryKind{}, nil, fmt.Errorf("failed to fetch ModelRegistry %q: %w", name, err)
	}

	typed, err := convertUnstructuredToModelRegistryKind(*item, true, ctx, clients.core)
	if err != nil {
		return models.ModelRegistryKind{}, nil, err
	}

	var password *string
	if val := typed.Spec.DatabaseConfig.PasswordSecret.Value; val != "" {
		password = &val
	}

	return typed, password, nil
}

func (r *ModelRegistrySettingsRepository) Create(
	ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	payload models.ModelRegistrySettingsPayload,
	dryRunOnly bool,
) (models.ModelRegistryKind, error) {
	if namespace == "" {
		return models.ModelRegistryKind{}, errors.New("namespace is required")
	}

	clients, err := r.buildClients(client)
	if err != nil {
		return models.ModelRegistryKind{}, err
	}

	model := payload.ModelRegistry
	if model.Metadata.Namespace == "" {
		model.Metadata.Namespace = namespace
	}
	if model.APIVersion == "" {
		model.APIVersion = fmt.Sprintf("%s/%s", modelRegistryGVR.Group, modelRegistryGVR.Version)
	}
	if model.Kind == "" {
		model.Kind = "ModelRegistry"
	}

	if payload.NewDatabaseCACertificate != nil && *payload.NewDatabaseCACertificate != "" {
		if model.Spec.DatabaseConfig.SSLRootCertificateConfigMap == nil || model.Spec.DatabaseConfig.SSLRootCertificateConfigMap.Name == "" {
			return models.ModelRegistryKind{}, errors.New("sslRootCertificateConfigMap.name is required when providing a new CA certificate")
		}
		entry, err := r.createCACertificateConfigMap(ctx, clients, namespace, model.Spec.DatabaseConfig.SSLRootCertificateConfigMap, *payload.NewDatabaseCACertificate)
		if err != nil {
			return models.ModelRegistryKind{}, err
		}
		model.Spec.DatabaseConfig.SSLRootCertificateConfigMap = entry
	}

	createOnce := func(dryRun bool) (models.ModelRegistryKind, error) {
		modelCopy := model

		if payload.DatabasePassword != nil && *payload.DatabasePassword != "" {
			secret, err := r.createDatabaseSecret(ctx, clients, namespace, modelCopy, *payload.DatabasePassword, dryRun)
			if err != nil {
				return models.ModelRegistryKind{}, err
			}
			modelCopy.Spec.DatabaseConfig.PasswordSecret = models.PasswordSecret{
				Name: secret.Name,
				Key:  databasePasswordKey,
			}
		}

		unstructuredModel, err := convertModelToUnstructured(modelCopy)
		if err != nil {
			return models.ModelRegistryKind{}, err
		}

		if unstructuredModel.GetNamespace() == "" {
			unstructuredModel.SetNamespace(namespace)
		}

		opts := metav1.CreateOptions{}
		if dryRun {
			opts.DryRun = []string{metav1.DryRunAll}
		}

		created, err := clients.dynamic.Resource(modelRegistryGVR).Namespace(namespace).Create(ctx, &unstructuredModel, opts)
		if err != nil {
			return models.ModelRegistryKind{}, fmt.Errorf("failed to create ModelRegistry resource: %w", err)
		}

		typed, err := convertUnstructuredToModelRegistryKind(*created, false, ctx, clients.core)
		if err != nil {
			return models.ModelRegistryKind{}, err
		}
		return typed, nil
	}

	dryRunResult, err := createOnce(true)
	if err != nil {
		return models.ModelRegistryKind{}, err
	}
	if dryRunOnly {
		return dryRunResult, nil
	}
	return createOnce(false)
}

func (r *ModelRegistrySettingsRepository) Patch(
	ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	name string,
	patch map[string]interface{},
	databasePassword *string,
	newCACertificate *string,
	dryRunOnly bool,
) (models.ModelRegistryKind, *string, error) {
	if namespace == "" || name == "" {
		return models.ModelRegistryKind{}, nil, errors.New("namespace and name are required")
	}

	clients, err := r.buildClients(client)
	if err != nil {
		return models.ModelRegistryKind{}, nil, err
	}

	if newCACertificate != nil && *newCACertificate != "" {
		spec, _ := patch["spec"].(map[string]interface{})
		if spec == nil {
			return models.ModelRegistryKind{}, nil, errors.New("spec patch is required when providing a new CA certificate")
		}

		var dbKey string
		for _, candidate := range []string{"mysql", "postgres"} {
			if _, ok := spec[candidate]; ok {
				dbKey = candidate
				break
			}
		}

		if dbKey == "" {
			return models.ModelRegistryKind{}, nil, errors.New("database configuration must be included in the patch when updating the CA certificate")
		}

		dbSpec, _ := spec[dbKey].(map[string]interface{})
		if dbSpec == nil {
			return models.ModelRegistryKind{}, nil, errors.New("invalid database spec in patch")
		}

		entryMap, _ := dbSpec["sslRootCertificateConfigMap"].(map[string]interface{})
		if entryMap == nil {
			return models.ModelRegistryKind{}, nil, errors.New("sslRootCertificateConfigMap is required when uploading a new certificate")
		}

		entry := &models.Entry{
			Name: getString(entryMap, "name"),
			Key:  getString(entryMap, "key"),
		}
		if entry.Name == "" {
			return models.ModelRegistryKind{}, nil, errors.New("sslRootCertificateConfigMap.name is required")
		}

		resolvedEntry, err := r.createCACertificateConfigMap(ctx, clients, namespace, entry, *newCACertificate)
		if err != nil {
			return models.ModelRegistryKind{}, nil, err
		}

		dbSpec["sslRootCertificateConfigMap"] = map[string]interface{}{
			"name": resolvedEntry.Name,
			"key":  resolvedEntry.Key,
		}
	}

	patchBytes, err := json.Marshal(patch)
	if err != nil {
		return models.ModelRegistryKind{}, nil, fmt.Errorf("failed to encode patch body: %w", err)
	}

	patchOnce := func(dryRun bool) (models.ModelRegistryKind, *string, error) {
		opts := metav1.PatchOptions{}
		if dryRun {
			opts.DryRun = []string{metav1.DryRunAll}
		}

		patched, err := clients.dynamic.Resource(modelRegistryGVR).Namespace(namespace).Patch(ctx, name, types.MergePatchType, patchBytes, opts)
		if err != nil {
			return models.ModelRegistryKind{}, nil, fmt.Errorf("failed to patch ModelRegistry %q: %w", name, err)
		}

		typed, err := convertUnstructuredToModelRegistryKind(*patched, true, ctx, clients.core)
		if err != nil {
			return models.ModelRegistryKind{}, nil, err
		}

		var passwordResp *string
		if databasePassword != nil {
			secretRef := typed.Spec.DatabaseConfig.PasswordSecret
			if secretRef.Name != "" {
				if *databasePassword == "" {
					if err := r.deleteDatabasePasswordSecret(ctx, clients, namespace, secretRef, dryRun); err != nil && !apierrors.IsNotFound(err) {
						return models.ModelRegistryKind{}, nil, err
					}
				} else {
					if err := r.updateDatabasePassword(ctx, clients, namespace, secretRef, *databasePassword, dryRun); err != nil {
						return models.ModelRegistryKind{}, nil, err
					}
					passwordResp = databasePassword
				}
			}
		}

		return typed, passwordResp, nil
	}

	dryRunResult, dryRunPassword, err := patchOnce(true)
	if err != nil {
		return models.ModelRegistryKind{}, nil, err
	}
	if dryRunOnly {
		return dryRunResult, dryRunPassword, nil
	}

	return patchOnce(false)
}

func (r *ModelRegistrySettingsRepository) Delete(
	ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	name string,
	dryRunOnly bool,
) (*metav1.Status, error) {
	if namespace == "" || name == "" {
		return nil, errors.New("namespace and name are required")
	}

	clients, err := r.buildClients(client)
	if err != nil {
		return nil, err
	}

	existing, err := clients.dynamic.Resource(modelRegistryGVR).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ModelRegistry %q prior to deletion: %w", name, err)
	}

	typed, err := convertUnstructuredToModelRegistryKind(*existing, false, ctx, clients.core)
	if err != nil {
		return nil, err
	}

	secretRef := typed.Spec.DatabaseConfig.PasswordSecret

	deleteOnce := func(dryRun bool) (*metav1.Status, error) {
		opts := metav1.DeleteOptions{}
		if dryRun {
			opts.DryRun = []string{metav1.DryRunAll}
		}

		if err := clients.dynamic.Resource(modelRegistryGVR).Namespace(namespace).Delete(ctx, name, opts); err != nil {
			return nil, fmt.Errorf("failed to delete ModelRegistry %q: %w", name, err)
		}

		if secretRef.Name != "" {
			if err := r.deleteDatabasePasswordSecret(ctx, clients, namespace, secretRef, dryRun); err != nil && !apierrors.IsNotFound(err) {
				return nil, err
			}
		}

		status := &metav1.Status{Status: "Success"}
		if dryRun {
			status.Message = "DryRun"
		}
		return status, nil
	}

	dryRunResult, err := deleteOnce(true)
	if err != nil {
		return nil, err
	}
	if dryRunOnly {
		return dryRunResult, nil
	}

	return deleteOnce(false)
}

// TODO(downstream): Keep template-specific annotations downstream; upstream should rely on neutral secrets.
func (r *ModelRegistrySettingsRepository) createDatabaseSecret(
	ctx context.Context,
	clients *kubeClients,
	namespace string,
	model models.ModelRegistryKind,
	password string,
	dryRun bool,
) (*corev1.Secret, error) {
	if password == "" {
		return nil, errors.New("database password is required")
	}
	if model.Metadata.Name == "" {
		return nil, errors.New("model registry name is required")
	}
	if model.Spec.DatabaseConfig.Database == "" || model.Spec.DatabaseConfig.Username == "" {
		return nil, errors.New("database name and username are required to create credentials")
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-db-", model.Metadata.Name),
			Namespace:    namespace,
			Annotations: map[string]string{
				"template.openshift.io/expose-database_name": "{.data['database-name']}",
				"template.openshift.io/expose-username":      "{.data['database-user']}",
				"template.openshift.io/expose-password":      "{.data['database-password']}",
			},
			Labels: map[string]string{
				"modelregistry.opendatahub.io/managed": "true",
				"modelregistry.opendatahub.io/name":    model.Metadata.Name,
			},
		},
		StringData: map[string]string{
			"database-name":     model.Spec.DatabaseConfig.Database,
			"database-user":     model.Spec.DatabaseConfig.Username,
			databasePasswordKey: password,
		},
		Type: corev1.SecretTypeOpaque,
	}

	opts := metav1.CreateOptions{}
	if dryRun {
		opts.DryRun = []string{metav1.DryRunAll}
	}

	created, err := clients.core.CoreV1().Secrets(namespace).Create(ctx, secret, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create database secret: %w", err)
	}

	return created, nil
}

// TODO(upstream): Consider moving CA ConfigMap creation upstream if the API becomes shared across vendors.
func (r *ModelRegistrySettingsRepository) createCACertificateConfigMap(
	ctx context.Context,
	clients *kubeClients,
	namespace string,
	entry *models.Entry,
	certificate string,
) (*models.Entry, error) {
	key := entry.Key
	if key == "" {
		key = caCertificateKey
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      entry.Name,
			Namespace: namespace,
		},
		Data: map[string]string{
			key: certificate,
		},
	}

	created, err := clients.core.CoreV1().ConfigMaps(namespace).Create(ctx, cm, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create ConfigMap %q: %w", entry.Name, err)
	}

	resolvedKey := key
	for dataKey := range created.Data {
		resolvedKey = dataKey
		break
	}

	return &models.Entry{Name: created.Name, Key: resolvedKey}, nil
}

// TODO(upstream): Shared password rotation logic should live in the default repository implementation.
func (r *ModelRegistrySettingsRepository) updateDatabasePassword(
	ctx context.Context,
	clients *kubeClients,
	namespace string,
	secretRef models.PasswordSecret,
	password string,
	dryRun bool,
) error {
	if secretRef.Name == "" || secretRef.Key == "" {
		return nil
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(password))
	payload := map[string]map[string]string{
		"data": {
			secretRef.Key: encoded,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to encode secret patch: %w", err)
	}

	opts := metav1.PatchOptions{}
	if dryRun {
		opts.DryRun = []string{metav1.DryRunAll}
	}

	if _, err := clients.core.CoreV1().Secrets(namespace).Patch(ctx, secretRef.Name, types.MergePatchType, body, opts); err != nil {
		return fmt.Errorf("failed to update database secret: %w", err)
	}

	return nil
}

// TODO(upstream): Shared delete logic should move upstream together with the rest of the repository.
func (r *ModelRegistrySettingsRepository) deleteDatabasePasswordSecret(
	ctx context.Context,
	clients *kubeClients,
	namespace string,
	secretRef models.PasswordSecret,
	dryRun bool,
) error {
	if secretRef.Name == "" {
		return nil
	}

	opts := metav1.DeleteOptions{}
	if dryRun {
		opts.DryRun = []string{metav1.DryRunAll}
	}

	if err := clients.core.CoreV1().Secrets(namespace).Delete(ctx, secretRef.Name, opts); err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("failed to delete database secret: %w", err)
	}

	return nil
}

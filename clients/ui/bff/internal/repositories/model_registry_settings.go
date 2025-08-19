package repositories

import (
	"context"
	"fmt"

	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ModelRegistrySettingsRepository struct {
}

func NewModelRegistrySettingsRepository() *ModelRegistrySettingsRepository {
	return &ModelRegistrySettingsRepository{}
}

func (r *ModelRegistrySettingsRepository) GetGroups(ctx context.Context, client k8s.KubernetesClientInterface) ([]models.Group, error) {
	groupNames, err := client.GetGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching groups: %w", err)
	}

	var groups []models.Group
	for _, name := range groupNames {
		// Create mock users for each group to make the data more realistic
		var users []string
		switch name {
		case "dora-group-mock":
			users = []string{"dora-user@example.com", "dora-admin@example.com"}
		case "bella-group-mock":
			users = []string{"bella-user@example.com", "bella-maintainer@example.com"}
		default:
			users = []string{fmt.Sprintf("%s-user@example.com", name)}
		}

		groups = append(groups, models.NewGroup(name, users))
	}

	return groups, nil
}

func (m *ModelRegistrySettingsRepository) GetModelRegistrySettings(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, name string) (models.ModelRegistryKind, error) {

	unstructuredObj, err := client.GetModelRegistrySettingsByName(ctx, namespace, name)
	if err != nil {
		return models.ModelRegistryKind{}, fmt.Errorf("failed to get model registry setting %q: %w", name, err)
	}

	typed, err := ConvertUnstructuredToModelRegistryKind(unstructuredObj, true, ctx, client)
	if err != nil {
		return models.ModelRegistryKind{}, err
	}

	return typed, nil
}

func (m *ModelRegistrySettingsRepository) GetAllModelRegistriesSettings(sessionCtx context.Context, client k8s.KubernetesClientInterface, namespace string, labelSelector string) ([]models.ModelRegistryKind, error) {
	modelRegistries, err := client.GetModelRegistrySettings(sessionCtx, namespace, labelSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to get model registry settings: %w", err)
	}

	var modelRegistriesModels = []models.ModelRegistryKind{}
	for _, u := range modelRegistries {
		typed, err := ConvertUnstructuredToModelRegistryKind(u, false, sessionCtx, client) // skip secret extraction
		if err != nil {
			return nil, err
		}
		modelRegistriesModels = append(modelRegistriesModels, typed)
	}

	return modelRegistriesModels, nil
}

// ConvertUnstructuredToModelRegistryKind converts a Kubernetes unstructured object into a typed ModelRegistryKind struct.
// If extractSecrets is true, the method attempts to extract PasswordSecret and SSL-related secrets/config maps.
// Otherwise, it skips them (e.g., for list all operations).
func ConvertUnstructuredToModelRegistryKind(obj unstructured.Unstructured, extractSecrets bool, ctx context.Context, client k8s.KubernetesClientInterface) (models.ModelRegistryKind, error) {
	var typed models.ModelRegistryKind

	// Base conversion
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, &typed)
	if err != nil {
		return typed, fmt.Errorf("failed to convert model registry %q: %w", obj.GetName(), err)
	}

	// Extract optional database block manually (that is flattened in our ModelRegistryKind)
	spec, ok := obj.Object["spec"].(map[string]interface{})
	if !ok {
		return typed, fmt.Errorf("model registry %q missing spec field", obj.GetName())
	}

	switch {
	case spec["mysql"] != nil:
		db, _ := spec["mysql"].(map[string]interface{})
		typed.Spec.DatabaseConfig.DatabaseType = models.MySQL
		typed.Spec.DatabaseConfig.Database = getString(db, "database")
		typed.Spec.DatabaseConfig.Host = getString(db, "host")
		typed.Spec.DatabaseConfig.Port = getInt(db, "port")
		typed.Spec.DatabaseConfig.Username = getString(db, "username")
		typed.Spec.DatabaseConfig.SkipDBCreation = getBool(db, "skipDBCreation")

		if extractSecrets {
			if s, ok := db["passwordSecret"].(map[string]interface{}); ok {
				key := getString(s, "key")
				name := getString(s, "name")
				typed.Spec.DatabaseConfig.PasswordSecret.Key = key
				typed.Spec.DatabaseConfig.PasswordSecret.Name = name
				val, err := client.GetDatabaseSecretValue(ctx, typed.Metadata.Namespace, name, key)
				if err == nil {
					typed.Spec.DatabaseConfig.PasswordSecret.Value = val
				}

			}
			// TODO: implement this when we have the requirements from UI
			// typed.Spec.DatabaseConfig.SSLRootCertificateConfigMap = getMapKey(db, "sslRootCertificateConfigMap")
			// typed.Spec.DatabaseConfig.SSLRootCertificateSecret = getMapKey(db, "sslRootCertificateSecret")
		}

	case spec["postgres"] != nil:
		db, _ := spec["postgres"].(map[string]interface{})
		typed.Spec.DatabaseConfig.DatabaseType = models.Postgres
		typed.Spec.DatabaseConfig.Database = getString(db, "database")
		typed.Spec.DatabaseConfig.Host = getString(db, "host")
		typed.Spec.DatabaseConfig.Port = getInt(db, "port")
		typed.Spec.DatabaseConfig.Username = getString(db, "username")
		typed.Spec.DatabaseConfig.SkipDBCreation = getBool(db, "skipDBCreation")

		if extractSecrets {
			if s, ok := db["passwordSecret"].(map[string]interface{}); ok {
				key := getString(s, "key")
				name := getString(s, "name")
				typed.Spec.DatabaseConfig.PasswordSecret.Key = key
				typed.Spec.DatabaseConfig.PasswordSecret.Name = name
				val, err := client.GetDatabaseSecretValue(ctx, typed.Metadata.Namespace, name, key)
				if err == nil {
					typed.Spec.DatabaseConfig.PasswordSecret.Value = val
				}
			}
			// TODO: implement this when we have the requirements from UI
			// typed.Spec.DatabaseConfig.SSLRootCertificateConfigMap = getMapKey(db, "sslRootCertificateConfigMap")
			// typed.Spec.DatabaseConfig.SSLRootCertificateSecret = getMapKey(db, "sslRootCertificateSecret")
		}
	}

	return typed, nil
}

// CreateModelRegistryKindWithSecret creates the DB Secret (if needed) and the ModelRegistryKind .
// If dryRun==true, it first performs a dry-run of both Secret and CR; if that succeeds, it then does the real create.
// If dryRun==false, it skips the dry-run pre-check and just does the real create.
func (m *ModelRegistrySettingsRepository) CreateModelRegistryKindWithSecret(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, model models.ModelRegistryKind, databasePassword *string, dryRun bool) (models.ModelRegistryKind, error) {

	//TODO: We need to support creation of newCACertificateConfigMap

	// 1) Pre-flight dry-run (only if requested)
	if dryRun {
		if _, err := m.createModelRegistryKindAndSecret(ctx, client, namespace, model, databasePassword, true); err != nil {
			return models.ModelRegistryKind{}, err
		}
	}
	// 2) Actual create (secret + CR)
	finalMRKind, err := m.createModelRegistryKindAndSecret(ctx, client, namespace, model, databasePassword, false)
	if err != nil {
		return models.ModelRegistryKind{}, err
	}

	// 3) Convert the returned Unstructured into our typed model
	return ConvertUnstructuredToModelRegistryKind(finalMRKind, false, ctx, client)
}

// createModelRegistryKindAndSecret handles one “both secrets & CR” operation under a single dryRun flag.
func (m *ModelRegistrySettingsRepository) createModelRegistryKindAndSecret(
	ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	model models.ModelRegistryKind,
	databasePassword *string,
	dryRun bool,
) (unstructured.Unstructured, error) {

	if err := m.injectSecretIfRequired(ctx, client, namespace, &model, databasePassword, dryRun); err != nil {
		return unstructured.Unstructured{}, err
	}

	u, err := m.convertModelToUnstructured(model)
	if err != nil {
		return unstructured.Unstructured{}, err
	}

	return client.CreateModelRegistryKind(ctx, namespace, u, dryRun)
}

func (m *ModelRegistrySettingsRepository) injectSecretIfRequired(
	ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	model *models.ModelRegistryKind,
	databasePassword *string,
	dryRun bool,
) error {
	if databasePassword == nil || *databasePassword == "" {
		return nil
	}
	secret, err := client.CreateDatabaseSecret(
		ctx,
		model.Metadata.Name,
		namespace,
		model.Spec.DatabaseConfig.Database,
		model.Spec.DatabaseConfig.Username,
		*databasePassword,
		dryRun,
	)
	if err != nil {
		return fmt.Errorf("failed to create database secret: %w", err)
	}
	model.Spec.DatabaseConfig.PasswordSecret = models.PasswordSecret{
		Name: secret.Name,
		Key:  "database-password",
	}
	return nil
}

func (m *ModelRegistrySettingsRepository) convertModelToUnstructured(
	model models.ModelRegistryKind,
) (unstructured.Unstructured, error) {

	objMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&model)
	if err != nil {
		return unstructured.Unstructured{}, fmt.Errorf("failed to convert model to unstructured: %w", err)
	}

	// Clean up unnecessary fields
	delete(objMap, "status")
	if metadata, ok := objMap["metadata"].(map[string]interface{}); ok {
		delete(metadata, "creationTimestamp") // drop system timestamp
	}

	spec, ok := objMap["spec"].(map[string]interface{})
	if !ok {
		return unstructured.Unstructured{}, fmt.Errorf("invalid spec format")
	}
	delete(spec, "databaseConfig")

	// Build the nested db block (mysql or postgres)
	dbType := string(model.Spec.DatabaseConfig.DatabaseType)
	dbSpec := map[string]interface{}{
		"database":       model.Spec.DatabaseConfig.Database,
		"host":           model.Spec.DatabaseConfig.Host,
		"port":           model.Spec.DatabaseConfig.Port,
		"username":       model.Spec.DatabaseConfig.Username,
		"skipDBCreation": model.Spec.DatabaseConfig.SkipDBCreation,
	}
	if model.Spec.DatabaseConfig.PasswordSecret.Name != "" {
		dbSpec["passwordSecret"] = map[string]interface{}{
			"name": model.Spec.DatabaseConfig.PasswordSecret.Name,
			"key":  model.Spec.DatabaseConfig.PasswordSecret.Key,
		}
	}
	spec[dbType] = dbSpec

	// Construct the Unstructured value
	var u unstructured.Unstructured
	u.Object = objMap
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "modelregistry.opendatahub.io",
		Version: "v1alpha1",
		Kind:    "ModelRegistry",
	})

	return u, nil
}

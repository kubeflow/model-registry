package repositories

import (
	"context"
	"fmt"

	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

// convertUnstructuredToModelRegistryKind converts a dynamic ModelRegistry object into the strongly typed model.
// TODO(upstream): Move the conversion helpers next to the upstream CRD types so every distro reuses them.
func convertUnstructuredToModelRegistryKind(
	obj unstructured.Unstructured,
	extractSecrets bool,
	ctx context.Context,
	coreClient kubernetes.Interface,
) (models.ModelRegistryKind, error) {
	var typed models.ModelRegistryKind

	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, &typed); err != nil {
		return typed, fmt.Errorf("failed to convert model registry %q: %w", obj.GetName(), err)
	}

	spec, ok := obj.Object["spec"].(map[string]interface{})
	if !ok {
		return typed, fmt.Errorf("model registry %q missing spec field", obj.GetName())
	}

	switch {
	case spec["mysql"] != nil:
		if err := populateDatabaseConfig(spec["mysql"], models.MySQL, &typed, extractSecrets, ctx, coreClient); err != nil {
			return typed, err
		}
	case spec["postgres"] != nil:
		if err := populateDatabaseConfig(spec["postgres"], models.Postgres, &typed, extractSecrets, ctx, coreClient); err != nil {
			return typed, err
		}
	default:
		return typed, fmt.Errorf("model registry %q missing supported database configuration", obj.GetName())
	}

	return typed, nil
}

// TODO(upstream): This parser should be shared once upstream repositories adopt the same schema.
func populateDatabaseConfig(
	dbSpec interface{},
	dbType models.DatabaseType,
	typed *models.ModelRegistryKind,
	extractSecrets bool,
	ctx context.Context,
	coreClient kubernetes.Interface,
) error {
	db, ok := dbSpec.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid database spec format")
	}

	typed.Spec.DatabaseConfig.DatabaseType = dbType
	typed.Spec.DatabaseConfig.Database = getString(db, "database")
	typed.Spec.DatabaseConfig.Host = getString(db, "host")
	typed.Spec.DatabaseConfig.Port = getInt(db, "port")
	typed.Spec.DatabaseConfig.Username = getString(db, "username")
	typed.Spec.DatabaseConfig.SkipDBCreation = getBool(db, "skipDBCreation")

	if cm, ok := db["sslRootCertificateConfigMap"].(map[string]interface{}); ok {
		typed.Spec.DatabaseConfig.SSLRootCertificateConfigMap = &models.Entry{
			Name: getString(cm, "name"),
			Key:  getString(cm, "key"),
		}
	}

	if secretEntry, ok := db["sslRootCertificateSecret"].(map[string]interface{}); ok {
		typed.Spec.DatabaseConfig.SSLRootCertificateSecret = &models.Entry{
			Name: getString(secretEntry, "name"),
			Key:  getString(secretEntry, "key"),
		}
	}

	if secret, ok := db["passwordSecret"].(map[string]interface{}); ok {
		typed.Spec.DatabaseConfig.PasswordSecret = models.PasswordSecret{
			Name: getString(secret, "name"),
			Key:  getString(secret, "key"),
		}
		if extractSecrets && typed.Spec.DatabaseConfig.PasswordSecret.Name != "" {
			val, err := getDatabaseSecretValue(
				ctx,
				coreClient,
				typed.Metadata.Namespace,
				typed.Spec.DatabaseConfig.PasswordSecret.Name,
				typed.Spec.DatabaseConfig.PasswordSecret.Key,
			)
			if err == nil {
				typed.Spec.DatabaseConfig.PasswordSecret.Value = val
			}
		}
	}

	return nil
}

// getDatabaseSecretValue fetches a secret value from Kubernetes.
func getDatabaseSecretValue(ctx context.Context, coreClient kubernetes.Interface, namespace, secretName, key string) (string, error) {
	secret, err := coreClient.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get secret %q: %w", secretName, err)
	}

	val, ok := secret.Data[key]
	if !ok {
		return "", fmt.Errorf("key %q not found in secret %q", key, secretName)
	}

	return string(val), nil
}

// convertModelToUnstructured transforms our model into the CR form expected by Kubernetes.
// TODO(upstream): Relocate serialization helpers to the upstream API package when ready.
func convertModelToUnstructured(model models.ModelRegistryKind) (unstructured.Unstructured, error) {
	objMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&model)
	if err != nil {
		return unstructured.Unstructured{}, fmt.Errorf("failed to convert model to unstructured: %w", err)
	}

	delete(objMap, "status")
	if metadata, ok := objMap["metadata"].(map[string]interface{}); ok {
		delete(metadata, "creationTimestamp")
	}

	spec, ok := objMap["spec"].(map[string]interface{})
	if !ok {
		return unstructured.Unstructured{}, fmt.Errorf("invalid spec format")
	}

	delete(spec, "databaseConfig")
	dbType := string(model.Spec.DatabaseConfig.DatabaseType)
	if dbType == "" {
		return unstructured.Unstructured{}, fmt.Errorf("database type is required")
	}

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

	if model.Spec.DatabaseConfig.SSLRootCertificateConfigMap != nil {
		dbSpec["sslRootCertificateConfigMap"] = map[string]interface{}{
			"name": model.Spec.DatabaseConfig.SSLRootCertificateConfigMap.Name,
			"key":  model.Spec.DatabaseConfig.SSLRootCertificateConfigMap.Key,
		}
	}

	if model.Spec.DatabaseConfig.SSLRootCertificateSecret != nil {
		dbSpec["sslRootCertificateSecret"] = map[string]interface{}{
			"name": model.Spec.DatabaseConfig.SSLRootCertificateSecret.Name,
			"key":  model.Spec.DatabaseConfig.SSLRootCertificateSecret.Key,
		}
	}

	spec[dbType] = dbSpec

	var u unstructured.Unstructured
	u.Object = objMap
	return u, nil
}

// BuildModelRegistryPatch serializes the provided model into a merge-patch-friendly map.
// TODO(upstream): Expose this helper with the upstream API once ModelRegistry settings move upstream.
func BuildModelRegistryPatch(model models.ModelRegistryKind) (map[string]interface{}, error) {
	u, err := convertModelToUnstructured(model)
	if err != nil {
		return nil, err
	}
	return u.Object, nil
}

func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

func getInt(m map[string]interface{}, key string) int {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case int64:
			return int(v)
		case float64:
			return int(v)
		}
	}
	return 0
}

func getBool(m map[string]interface{}, key string) bool {
	if val, ok := m[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

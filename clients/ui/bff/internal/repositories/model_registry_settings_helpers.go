package repositories

import (
	"fmt"

	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	SecretKeySuffix    = "-database-password"
	SecretPasswordKey  = "database-password"
	ModelRegistryGroup = "modelregistry.opendatahub.io"
	ModelRegistryKind  = "ModelRegistry"
)

// BuildSecret creates a Kubernetes Secret for database password
func BuildSecret(namespace, registryName, password string) *corev1.Secret {
	secretName := registryName + SecretKeySuffix
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		Type: corev1.SecretTypeOpaque,
		StringData: map[string]string{
			SecretPasswordKey: password,
		},
	}
}

// addPasswordSecret adds password secret configuration to a database config map if the secret is not nil
func addPasswordSecret(config map[string]interface{}, secret *models.PasswordSecret) {
	if secret != nil {
		config["passwordSecret"] = map[string]interface{}{
			"name": secret.Name,
			"key":  secret.Key,
		}
	}
}

// addSSLConfig adds SSL certificate configuration (ConfigMap and Secret) to a database config map
func addSSLConfig(config map[string]interface{}, configMap, secret *models.Entry) {
	if configMap != nil {
		config["sslRootCertificateConfigMap"] = map[string]interface{}{
			"name": configMap.Name,
			"key":  configMap.Key,
		}
	}

	if secret != nil {
		config["sslRootCertificateSecret"] = map[string]interface{}{
			"name": secret.Name,
			"key":  secret.Key,
		}
	}
}

// ModelRegistryKindToUnstructured converts a ModelRegistryKind to an unstructured object
// for use with the dynamic client
func ModelRegistryKindToUnstructured(mr models.ModelRegistryKind) (*unstructured.Unstructured, error) {
	obj := &unstructured.Unstructured{
		Object: make(map[string]interface{}),
	}

	obj.SetAPIVersion(ModelRegistryGroup + "/v1beta1")
	obj.SetKind(ModelRegistryKind)
	obj.SetName(mr.Metadata.Name)
	obj.SetNamespace(mr.Metadata.Namespace)

	// Set annotations
	if len(mr.Metadata.Annotations) > 0 {
		obj.SetAnnotations(mr.Metadata.Annotations)
	}

	// Build spec
	spec := make(map[string]interface{})

	// Add grpc (empty object)
	spec["grpc"] = map[string]interface{}{}

	// Add rest (empty object)
	spec["rest"] = map[string]interface{}{}

	// Add istio configuration
	spec["istio"] = map[string]interface{}{
		"gateway": map[string]interface{}{
			"grpc": map[string]interface{}{
				"tls": map[string]interface{}{},
			},
			"rest": map[string]interface{}{
				"tls": map[string]interface{}{},
			},
		},
	}

	// Add database configuration
	if mr.Spec.MySQL != nil {
		mysqlConfig := make(map[string]interface{})
		mysqlConfig["host"] = mr.Spec.MySQL.Host
		mysqlConfig["database"] = mr.Spec.MySQL.Database
		mysqlConfig["username"] = mr.Spec.MySQL.Username

		if mr.Spec.MySQL.Port != nil {
			mysqlConfig["port"] = int64(*mr.Spec.MySQL.Port)
		}

		addPasswordSecret(mysqlConfig, mr.Spec.MySQL.PasswordSecret)

		if mr.Spec.MySQL.SkipDBCreation != nil {
			mysqlConfig["skipDBCreation"] = *mr.Spec.MySQL.SkipDBCreation
		}

		addSSLConfig(mysqlConfig, mr.Spec.MySQL.SSLRootCertificateConfigMap, mr.Spec.MySQL.SSLRootCertificateSecret)

		spec["mysql"] = mysqlConfig
	}

	if mr.Spec.Postgres != nil {
		postgresConfig := make(map[string]interface{})
		postgresConfig["database"] = mr.Spec.Postgres.Database

		// Only add host/username for external databases
		if mr.Spec.Postgres.Host != "" {
			postgresConfig["host"] = mr.Spec.Postgres.Host
		}

		if mr.Spec.Postgres.Username != "" {
			postgresConfig["username"] = mr.Spec.Postgres.Username
		}

		if mr.Spec.Postgres.Port != nil {
			postgresConfig["port"] = int64(*mr.Spec.Postgres.Port)
		}

		addPasswordSecret(postgresConfig, mr.Spec.Postgres.PasswordSecret)

		if mr.Spec.Postgres.GenerateDeployment != nil {
			postgresConfig["generateDeployment"] = *mr.Spec.Postgres.GenerateDeployment
		}

		if mr.Spec.Postgres.SkipDBCreation != nil {
			postgresConfig["skipDBCreation"] = *mr.Spec.Postgres.SkipDBCreation
		}

		if mr.Spec.Postgres.SSLMode != "" {
			postgresConfig["sslMode"] = mr.Spec.Postgres.SSLMode
		}

		addSSLConfig(postgresConfig, mr.Spec.Postgres.SSLRootCertificateConfigMap, mr.Spec.Postgres.SSLRootCertificateSecret)

		spec["postgres"] = postgresConfig
	}

	if err := unstructured.SetNestedMap(obj.Object, spec, "spec"); err != nil {
		return nil, fmt.Errorf("failed to set spec: %w", err)
	}

	return obj, nil
}

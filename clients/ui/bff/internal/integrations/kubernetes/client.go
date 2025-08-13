package kubernetes

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const ComponentLabelValue = "model-registry"

type KubernetesClientInterface interface {
	// Service discovery
	GetServiceNames(ctx context.Context, namespace string) ([]string, error)
	GetServiceDetailsByName(ctx context.Context, namespace, serviceName string) (ServiceDetails, error)
	GetServiceDetails(ctx context.Context, namespace string) ([]ServiceDetails, error)

	// Model Registry Settings
	GetModelRegistrySettings(ctx context.Context, namespace string, labelSelector string) ([]unstructured.Unstructured, error)
	GetModelRegistrySettingsByName(ctx context.Context, namespace string, name string) (unstructured.Unstructured, error)
	CreateModelRegistryKind(ctx context.Context, namespace string, modelRegistryKind unstructured.Unstructured, dryRun bool) (unstructured.Unstructured, error)

	// Namespace access
	GetNamespaces(ctx context.Context, identity *RequestIdentity) ([]corev1.Namespace, error)

	// Database Secret access
	GetDatabaseSecretValue(ctx context.Context, namespace, secretName, key string) (string, error)
	CreateDatabaseSecret(ctx context.Context, name string, namespace string, database string, databaseUsername string, databasePassword string, dryRun bool) (*corev1.Secret, error)

	// Permission checks (abstracted SAR/SelfSAR)
	CanListServicesInNamespace(ctx context.Context, identity *RequestIdentity, namespace string) (bool, error)
	CanAccessServiceInNamespace(ctx context.Context, identity *RequestIdentity, namespace, serviceName string) (bool, error)

	// Meta
	IsClusterAdmin(identity *RequestIdentity) (bool, error)
	BearerToken() (string, error)
	GetUser(identity *RequestIdentity) (string, error)

	// Model Registry Settings
	GetGroups(ctx context.Context) ([]string, error)
}

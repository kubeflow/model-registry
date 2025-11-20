package kubernetes

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const ComponentLabelValue = "model-registry"
const ComponentLabelValueCatalog = "model-catalog"

type KubernetesClientInterface interface {
	// Service discovery
	GetServiceNames(ctx context.Context, namespace string) ([]string, error)
	GetServiceDetailsByName(ctx context.Context, namespace, serviceName string, serviceType string) (ServiceDetails, error)
	GetServiceDetails(ctx context.Context, namespace string) ([]ServiceDetails, error)

	// Namespace access
	GetNamespaces(ctx context.Context, identity *RequestIdentity) ([]corev1.Namespace, error)

	// Permission checks (abstracted SAR/SelfSAR)
	CanListServicesInNamespace(ctx context.Context, identity *RequestIdentity, namespace string) (bool, error)
	CanAccessServiceInNamespace(ctx context.Context, identity *RequestIdentity, namespace, serviceName string) (bool, error)
	GetSelfSubjectRulesReview(ctx context.Context, identity *RequestIdentity, namespace string) ([]string, error)

	// Meta
	IsClusterAdmin(identity *RequestIdentity) (bool, error)
	BearerToken() (string, error)
	GetUser(identity *RequestIdentity) (string, error)

	// Model Registry Settings
	GetGroups(ctx context.Context) ([]string, error)
	CreateModelRegistry(ctx context.Context, namespace string, obj *unstructured.Unstructured) (*unstructured.Unstructured, error)
	CreateSecret(ctx context.Context, namespace string, secret *corev1.Secret) (*corev1.Secret, error)
}

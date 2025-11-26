package kubernetes

import (
	"context"
	corev1 "k8s.io/api/core/v1"
)

const ComponentLabelValue = "model-registry"
const ComponentLabelValueCatalog = "model-catalog"

// TODO ppadti double check if the config map key is indeed sources.yaml
const CatalogSourceKey = "sources.yaml"
const CatalogSourceDefaultConfigMapName = "model-catalog-source-config"
const CatalogSourceUserConfigMapName = "model-catalog-sources"

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

	//Model Catalog Settings
	GetAllCatalogSourceConfigs(ctx context.Context, namespace string) (corev1.ConfigMap, corev1.ConfigMap, error)
	//TODO ppadti add other methods here
}

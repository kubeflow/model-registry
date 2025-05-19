package kubernetes

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

var modelRegistryGVR = schema.GroupVersionResource{
	Group:    "modelregistry.opendatahub.io",
	Version:  "v1alpha1",
	Resource: "modelregistries",
}

// OpenShift Groups GVR
var groupGVR = schema.GroupVersionResource{
	Group:    "user.openshift.io",
	Version:  "v1",
	Resource: "groups",
}

// RHOAITokenKubernetesClient wraps the default TokenKubernetesClient but overrides selected methods.
// It also adds a dynamic client to the base client.
type RHOAITokenKubernetesClient struct {
	*kubernetes.TokenKubernetesClient
	DynClient dynamic.Interface
}

func NewRHOAIKubernetesClient(token string, logger *slog.Logger) (kubernetes.KubernetesClientInterface, error) {
	baseClient, err := kubernetes.NewTokenKubernetesClient(token, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create base token client: %w", err)
	}

	typed := baseClient.(*kubernetes.TokenKubernetesClient)

	dynClient, err := dynamic.NewForConfig(typed.RESTConfig())

	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	return &RHOAITokenKubernetesClient{
		TokenKubernetesClient: baseClient.(*kubernetes.TokenKubernetesClient),
		DynClient:             dynClient,
	}, nil
}

// GetGroups overrides the default stub logic with RHOAI-specific logic to fetch OpenShift groups.
// This method queries the OpenShift user.openshift.io/v1 Group API to retrieve all available groups
// that can be used for RBAC and permissions management in RHOAI.
func (c *RHOAITokenKubernetesClient) GetGroups(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	c.Logger.Info("Fetching OpenShift groups for RHOAI")

	// List all OpenShift groups using the dynamic client
	// Note: Groups are cluster-scoped resources, so no namespace is specified
	list, err := c.DynClient.
		Resource(groupGVR).
		List(ctx, v1.ListOptions{})
	if err != nil {
		c.Logger.Error("Failed to list OpenShift Groups", "error", err)
		// Check if this is a permission error
		if err.Error() != "" {
			c.Logger.Warn("User may not have permissions to list groups cluster-wide")
		}
		return nil, fmt.Errorf("failed to list OpenShift groups: %w", err)
	}

	var groupNames []string
	for _, item := range list.Items {
		name, found, err := unstructured.NestedString(item.Object, "metadata", "name")
		if err != nil {
			c.Logger.Warn("Failed to extract group name", "error", err, "group", item.GetName())
			continue
		}
		if !found || name == "" {
			c.Logger.Warn("Group name not found or empty in metadata", "group", item.GetName())
			continue
		}
		groupNames = append(groupNames, name)
	}

	c.Logger.Info("Successfully fetched OpenShift groups", "count", len(groupNames), "groups", groupNames)
	return groupNames, nil
}

func (kc *RHOAITokenKubernetesClient) CreateModelRegistryKind(ctx context.Context, namespace string, modelRegistryKind unstructured.Unstructured, dryRun bool) (unstructured.Unstructured, error) {

	opts := v1.CreateOptions{}
	if dryRun {
		opts.DryRun = []string{v1.DryRunAll}
	}

	created, err := kc.DynClient.
		Resource(modelRegistryGVR).
		Namespace(namespace).
		Create(ctx, &modelRegistryKind, opts)

	if err != nil {
		return unstructured.Unstructured{}, fmt.Errorf("failed to create ModelRegistryKind: %w", err)
	}

	return *created, nil
}

// GetModelRegistrySettings overrides the default stub logic with RHOAI-specific logic.
func (c *RHOAITokenKubernetesClient) GetModelRegistrySettings(ctx context.Context, namespace string, labelSelector string) ([]unstructured.Unstructured, error) {

	listOptions := v1.ListOptions{}
	if labelSelector != "" {
		listOptions.LabelSelector = labelSelector
	}

	list, err := c.DynClient.
		Resource(modelRegistryGVR).
		Namespace(namespace).
		List(ctx, listOptions)
	if err != nil {
		c.Logger.Error("Failed to list ModelRegistry CRs", "error", err)
		return nil, fmt.Errorf("failed to list ModelRegistry CRs in namespace %q: %w", namespace, err)
	}

	return list.Items, nil
}

func (c *RHOAITokenKubernetesClient) GetModelRegistrySettingsByName(ctx context.Context, namespace string, name string) (unstructured.Unstructured, error) {

	result, err := c.DynClient.
		Resource(modelRegistryGVR).
		Namespace(namespace).
		Get(ctx, name, v1.GetOptions{})
	if err != nil {
		c.Logger.Error("Failed to get ModelRegistry CR", "error", err, "name", name, "namespace", namespace)
		return unstructured.Unstructured{}, fmt.Errorf("failed to get ModelRegistry %q in namespace %q: %w", name, namespace, err)
	}

	return *result, nil
}

func (c *RHOAITokenKubernetesClient) CreateDatabaseSecret(ctx context.Context, name string, namespace string, database string, databaseUsername string, databasePassword string, dryRun bool) (*corev1.Secret, error) {

	if database == "" || databaseUsername == "" {
		return nil, fmt.Errorf("invalid database config: database name or username missing")
	}

	secret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			GenerateName: name + "-db-",
			Namespace:    namespace,
			Annotations: map[string]string{
				"template.openshift.io/expose-database_name": "{.data['database-name']}",
				"template.openshift.io/expose-username":      "{.data['database-user']}",
				"template.openshift.io/expose-password":      "{.data['database-password']}",
			},
			Labels: map[string]string{
				"modelregistry.opendatahub.io/managed": "true",
				"modelregistry.opendatahub.io/name":    name,
			},
		},
		StringData: map[string]string{
			"database-name":     database,
			"database-user":     databaseUsername,
			"database-password": databasePassword,
		},
		Type: corev1.SecretTypeOpaque,
	}

	options := v1.CreateOptions{}
	if dryRun {
		options.DryRun = []string{v1.DryRunAll}
	}

	created, err := c.Client.CoreV1().Secrets(namespace).Create(ctx, secret, options)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret: %w", err)
	}

	return created, nil
}

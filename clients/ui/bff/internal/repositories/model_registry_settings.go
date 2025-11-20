package repositories

import (
	"context"
	"fmt"

	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
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

// CreateModelRegistry creates a ModelRegistry custom resource and optionally a Secret for the database password
func (r *ModelRegistrySettingsRepository) CreateModelRegistry(
	ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	payload models.ModelRegistrySettingsPayload,
) (*models.ModelRegistryKind, error) {
	registryName := payload.ModelRegistry.Metadata.Name
	if registryName == "" {
		return nil, fmt.Errorf("model registry name is required")
	}

	// Determine if we need to create a Secret for the database password
	needsSecret := false
	if payload.DatabasePassword != nil && *payload.DatabasePassword != "" {
		// External databases require a Secret
		if payload.ModelRegistry.Spec.MySQL != nil ||
			(payload.ModelRegistry.Spec.Postgres != nil && payload.ModelRegistry.Spec.Postgres.GenerateDeployment == nil) {
			needsSecret = true
		}
	}

	// Create Secret if needed
	if needsSecret {
		secret := BuildSecret(namespace, registryName, *payload.DatabasePassword)
		_, err := client.CreateSecret(ctx, namespace, secret)
		if err != nil {
			return nil, fmt.Errorf("failed to create database password secret: %w", err)
		}

		// Update the ModelRegistry spec to reference the Secret
		secretName := registryName + SecretKeySuffix
		passwordSecret := &models.PasswordSecret{
			Name: secretName,
			Key:  SecretPasswordKey,
		}

		if payload.ModelRegistry.Spec.MySQL != nil {
			payload.ModelRegistry.Spec.MySQL.PasswordSecret = passwordSecret
		} else if payload.ModelRegistry.Spec.Postgres != nil {
			payload.ModelRegistry.Spec.Postgres.PasswordSecret = passwordSecret
		}
	}

	// Convert ModelRegistryKind to unstructured
	unstructuredObj, err := ModelRegistryKindToUnstructured(payload.ModelRegistry)
	if err != nil {
		return nil, fmt.Errorf("failed to convert ModelRegistry to unstructured: %w", err)
	}

	// Create ModelRegistry CR
	_, err = client.CreateModelRegistry(ctx, namespace, unstructuredObj)
	if err != nil {
		return nil, fmt.Errorf("failed to create ModelRegistry CR: %w", err)
	}

	// Return the created ModelRegistry
	return &payload.ModelRegistry, nil
}

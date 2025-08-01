package models

import corev1 "k8s.io/api/core/v1"

type NamespaceModel struct {
	Name        string  `json:"name"`
	DisplayName *string `json:"displayName,omitempty"`
}

func NewNamespaceModelFromNamespace(name string) NamespaceModel {
	displayName := name // For now, use name as display name, but this can be customized later
	return NamespaceModel{
		Name:        name,
		DisplayName: &displayName,
	}
}

// NewNamespaceModelFromK8sNamespace creates a NamespaceModel from a Kubernetes namespace object
func NewNamespaceModelFromK8sNamespace(namespace corev1.Namespace) NamespaceModel {
	var displayName *string
	
	// Try to get display name from OpenShift annotation
	if annotations := namespace.GetAnnotations(); annotations != nil {
		if osDisplayName, exists := annotations["openshift.io/display-name"]; exists && osDisplayName != "" {
			displayName = &osDisplayName
		}
	}
	
	// Fallback to namespace name if no display name annotation exists
	if displayName == nil {
		fallback := namespace.Name
		displayName = &fallback
	}
	
	return NamespaceModel{
		Name:        namespace.Name,
		DisplayName: displayName,
	}
}

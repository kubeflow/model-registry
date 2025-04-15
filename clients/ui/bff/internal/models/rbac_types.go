package models

import (
	rbacv1 "k8s.io/api/rbac/v1"
)

// CertificateItem represents a Secret or ConfigMap suitable for certificate usage.
type CertificateItem struct {
	Name string   `json:"name"`
	Keys []string `json:"keys"`
}

// CertificateList holds lists of Secrets and ConfigMaps.
type CertificateList struct {
	Secrets    []CertificateItem `json:"secrets"`
	ConfigMaps []CertificateItem `json:"configMaps"`
}

// RoleBinding represents a Kubernetes RoleBinding (simplified for API).
// Using the actual k8s type for consistency might be better if full fidelity is needed.
type RoleBinding rbacv1.RoleBinding

// RoleBindingList represents a list of Kubernetes RoleBindings.
type RoleBindingList struct {
	Items []RoleBinding `json:"items"`
	// Add other list metadata if needed (e.g., from metav1.ListMeta)
}

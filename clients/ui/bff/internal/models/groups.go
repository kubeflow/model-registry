package models

import "time"

// K8sObjectMeta represents the metadata of a Kubernetes object
type K8sObjectMeta struct {
	Name              string            `json:"name"`
	Namespace         *string           `json:"namespace,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	CreationTimestamp *time.Time        `json:"creationTimestamp,omitempty"`
	ResourceVersion   *string           `json:"resourceVersion,omitempty"`
	UID               *string           `json:"uid,omitempty"`
}

// Group represents an OpenShift/Kubernetes Group object
type Group struct {
	APIVersion *string       `json:"apiVersion,omitempty"`
	Kind       *string       `json:"kind,omitempty"`
	Metadata   K8sObjectMeta `json:"metadata"`
	Users      []string      `json:"users"`
}

// Legacy GroupModel for backward compatibility (if needed)
type GroupModel struct {
	Name string `json:"name"`
}

// NewGroup creates a new Group with the specified name and users
func NewGroup(name string, users []string) Group {
	apiVersion := "user.openshift.io/v1"
	kind := "Group"
	now := time.Now()

	return Group{
		APIVersion: &apiVersion,
		Kind:       &kind,
		Metadata: K8sObjectMeta{
			Name:              name,
			CreationTimestamp: &now,
		},
		Users: users,
	}
}

func NewGroupModel(name string) GroupModel {
	return GroupModel{
		Name: name,
	}
}

package models

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

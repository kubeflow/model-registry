package models

type NamespaceModel struct {
	Name        string `json:"name"`
	DisplayName string `json:"display-name"`
}

func NewNamespaceModelFromNamespace(name string) NamespaceModel {
	return NamespaceModel{
		Name:        name,
		DisplayName: name, // For now, use name as display name, but this can be customized later
	}
}

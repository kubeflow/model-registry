package models

type NamespaceModel struct {
	Name string `json:"name"`
}

func NewNamespaceModelFromNamespace(name string) NamespaceModel {
	return NamespaceModel{
		Name: name,
	}
}

package models

type ModelRegistryModel struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

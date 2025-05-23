package models

type GroupModel struct {
	Name string `json:"name"`
}

func NewGroupModel(name string) GroupModel {
	return GroupModel{
		Name: name,
	}
}

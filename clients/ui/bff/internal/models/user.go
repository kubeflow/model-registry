package models

type User struct {
	UserID       string `json:"userId"`
	ClusterAdmin bool   `json:"clusterAdmin"`
}

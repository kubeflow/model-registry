package models

type User struct {
	UserID       string `json:"user-id"`
	ClusterAdmin bool   `json:"cluster-admin"`
}

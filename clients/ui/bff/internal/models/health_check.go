package models

type SystemInfo struct {
	Version string `json:"version"`
}

type HealthCheckModel struct {
	Status     string     `json:"status"`
	SystemInfo SystemInfo `json:"system_info"`
	UserID     string     `json:"userId"`
}

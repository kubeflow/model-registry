package models

type ModelRegistryModel struct {
	Name          string `json:"name"`
	DisplayName   string `json:"displayName"`
	Description   string `json:"description"`
	ServerAddress string `json:"serverAddress"`
	IsHTTPS       bool   `json:"isHttps"`
}

// ServiceAuthorizationContext holds the authorization decision context
type ServiceAuthorizationContext struct {
	AllowList           bool
	AllowedServiceNames []string
	Namespace           string
}

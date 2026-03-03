package models

type ModelRegistryModel struct {
	Name          string `json:"name"`
	DisplayName   string `json:"displayName"`
	Description   string `json:"description"`
	ServerAddress string `json:"serverAddress"`
	IsHTTPS       bool   `json:"isHttps"`
	IsAvailable   bool   `json:"isAvailable"` // true if Endpoints has ready addresses
}

// ServiceAuthorizationContext holds the authorization decision context
type ServiceAuthorizationContext struct {
	AllowList           bool
	AllowedServiceNames []string
	Namespace           string
}

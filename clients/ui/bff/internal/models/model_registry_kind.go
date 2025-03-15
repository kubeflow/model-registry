package models

import (
	"time"
)

type ModelRegistryKind struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   Metadata          `json:"metadata"`
	Spec       ModelRegistrySpec `json:"spec"`
	Status     Status            `json:"status"`
}

type Metadata struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	CreationTimestamp time.Time         `json:"creationTimestamp"`
	Annotations       map[string]string `json:"annotations,omitempty"`
}

// EmptyObject represents an empty object at create time, properties here aren't used by the UI
type EmptyObject struct{}

type ModelRegistrySpec struct {
	GRPC           EmptyObject    `json:"grpc"` // Empty object at create time, properties here aren't used by the UI
	REST           EmptyObject    `json:"rest"` // Empty object at create time, properties here aren't used by the UI
	Istio          IstioConfig    `json:"istio"`
	DatabaseConfig DatabaseConfig `json:"databaseConfig"`
}

type IstioConfig struct {
	Gateway GatewayConfig `json:"gateway"`
}

type GatewayConfig struct {
	GRPC GRPCConfig `json:"grpc"`
	REST RESTConfig `json:"rest"`
}

type GRPCConfig struct {
	TLS EmptyObject `json:"tls"` // Empty object at create time, properties here aren't used by the UI
}

type RESTConfig struct {
	TLS EmptyObject `json:"tls"` // Empty object at create time, properties here aren't used by the UI
}

type DatabaseType string

const (
	MySQL    DatabaseType = "mysql"
	Postgres DatabaseType = "postgres"
)

type DatabaseConfig struct {
	DatabaseType                DatabaseType   `json:"databaseType"`
	Database                    string         `json:"database"`
	Host                        string         `json:"host"`
	PasswordSecret              PasswordSecret `json:"passwordSecret,omitempty"`
	Port                        int            `json:"port"`
	SkipDBCreation              bool           `json:"skipDBCreation"`
	Username                    string         `json:"username"`
	SSLRootCertificateConfigMap string         `json:"sslRootCertificateConfigMap"`
	SSLRootCertificateSecret    string         `json:"sslRootCertificateSecret"`
}

type PasswordSecret struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type Status struct {
	Conditions []Condition `json:"conditions"`
}

type Condition struct {
	LastTransitionTime time.Time `json:"lastTransitionTime"`
	Message            string    `json:"message"`
	Reason             string    `json:"reason"`
	Status             string    `json:"status"`
	Type               string    `json:"type"`
}

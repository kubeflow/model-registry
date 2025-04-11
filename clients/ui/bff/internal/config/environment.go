package config

import "log/slog"

const (
	// AuthMethodInternal uses the credentials of the running backend.
	// If running inside the cluster, it uses the pod's service account.
	// If running locally (e.g. for development), it uses the current user's kubeconfig context.
	// This is the default authentication method.
	// This uses kubeflow-userid header to carry the user identity.
	AuthMethodInternal = "internal"

	// AuthMethodUser uses a user-provided Bearer token for authentication.
	AuthMethodUser = "user_token"
)

type EnvConfig struct {
	Port            int
	MockK8Client    bool
	MockMRClient    bool
	DevMode         bool
	StandaloneMode  bool
	DevModePort     int
	StaticAssetsDir string
	LogLevel        slog.Level
	AllowedOrigins  []string
	// Either AuthMethodInternal or AuthMethodUser
	AuthMethod string
}

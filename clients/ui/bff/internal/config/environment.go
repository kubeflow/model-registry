package config

import "log/slog"

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
}

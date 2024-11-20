package config

type EnvConfig struct {
	Port         int
	MockK8Client bool
	MockMRClient bool
	DevMode      bool
	DevModePort  int
}

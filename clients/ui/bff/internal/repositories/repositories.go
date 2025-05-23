package repositories

// Repositories struct is a single convenient container to hold and represent all our repositories.
type Repositories struct {
	HealthCheck           *HealthCheckRepository
	ModelRegistry         *ModelRegistryRepository
	ModelRegistrySettings *ModelRegistrySettingsRepository
	ModelRegistryClient   ModelRegistryClientInterface
	User                  *UserRepository
	Namespace             *NamespaceRepository
}

func NewRepositories(modelRegistryClient ModelRegistryClientInterface) *Repositories {
	return &Repositories{
		HealthCheck:           NewHealthCheckRepository(),
		ModelRegistry:         NewModelRegistryRepository(),
		ModelRegistrySettings: NewModelRegistrySettingsRepository(),
		ModelRegistryClient:   modelRegistryClient,
		User:                  NewUserRepository(),
		Namespace:             NewNamespaceRepository(),
	}
}

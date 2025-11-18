package repositories

// Repositories struct is a single convenient container to hold and represent all our repositories.
type Repositories struct {
	HealthCheck                    *HealthCheckRepository
	ModelRegistry                  *ModelRegistryRepository
	ModelCatalog                   *ModelCatalogRepository
	ModelRegistrySettings          *ModelRegistrySettingsRepository
	ModelRegistryClient            ModelRegistryClientInterface
	ModelCatalogClient             ModelCatalogClientInterface
	ModelCatalogSettingsRepository ModelCatalogSettingsRepositoryInterface
	User                           *UserRepository
	Namespace                      *NamespaceRepository
}

func NewRepositories(modelRegistryClient ModelRegistryClientInterface, modelCatalogClient ModelCatalogClientInterface, modelCatalogSettingsRepository ModelCatalogSettingsRepositoryInterface) *Repositories {
	return &Repositories{
		HealthCheck:                    NewHealthCheckRepository(),
		ModelRegistry:                  NewModelRegistryRepository(),
		ModelCatalog:                   NewCatalogRepository(),
		ModelCatalogClient:             modelCatalogClient,
		ModelRegistrySettings:          NewModelRegistrySettingsRepository(),
		ModelRegistryClient:            modelRegistryClient,
		ModelCatalogSettingsRepository: modelCatalogSettingsRepository,
		User:                           NewUserRepository(),
		Namespace:                      NewNamespaceRepository(),
	}
}

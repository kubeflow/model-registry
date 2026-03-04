package repositories

// Repositories struct is a single convenient container to hold and represent all our repositories.
type Repositories struct {
	HealthCheck                    *HealthCheckRepository
	ModelRegistry                  *ModelRegistryRepository
	ModelCatalog                   *ModelCatalogRepository
	ModelRegistrySettings          *ModelRegistrySettingsRepository
	ModelRegistryClient            ModelRegistryClientInterface
	ModelCatalogClient             ModelCatalogClientInterface
	ModelCatalogSettingsRepository *ModelCatalogSettingsRepository
	User                           *UserRepository
	Namespace                      *NamespaceRepository
}

func NewRepositories(modelRegistryClient ModelRegistryClientInterface, modelCatalogClient ModelCatalogClientInterface, isFederatedMode bool, podNamespace string) *Repositories {
	return &Repositories{
		HealthCheck:                    NewHealthCheckRepository(),
		ModelRegistry:                  NewModelRegistryRepository(isFederatedMode, podNamespace),
		ModelCatalog:                   NewCatalogRepository(),
		ModelCatalogClient:             modelCatalogClient,
		ModelRegistrySettings:          NewModelRegistrySettingsRepository(),
		ModelRegistryClient:            modelRegistryClient,
		ModelCatalogSettingsRepository: NewModelCatalogSettingsRepository(),
		User:                           NewUserRepository(),
		Namespace:                      NewNamespaceRepository(),
	}
}

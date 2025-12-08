import {
  CatalogSourceConfigList,
  CatalogSourceType,
  HuggingFaceCatalogSourceConfig,
  YamlCatalogSourceConfig,
} from '~/app/modelCatalogTypes';

export const mockYamlCatalogSourceConfig = (
  partial?: Partial<YamlCatalogSourceConfig>,
): YamlCatalogSourceConfig => ({
  id: 'sample_source_1',
  name: 'Source 1',
  type: CatalogSourceType.YAML,
  enabled: true,
  includedModels: [],
  excludedModels: [],
  isDefault: true,
  yaml: 'models:\n  - name: model1',
  ...partial,
});

export const mockHuggingFaceCatalogSourceConfig = (
  partial?: Partial<HuggingFaceCatalogSourceConfig>,
): HuggingFaceCatalogSourceConfig => ({
  id: 'source_2',
  name: 'Huggingface source 2',
  type: CatalogSourceType.HUGGING_FACE,
  enabled: true,
  includedModels: [],
  excludedModels: [],
  isDefault: false,
  allowedOrganization: 'org1',
  apiKey: 'apikey',
  ...partial,
});

export const mockCatalogSourceConfigList = (
  partial?: Partial<CatalogSourceConfigList>,
): CatalogSourceConfigList => ({
  catalogs: [
    mockYamlCatalogSourceConfig({
      id: 'sample_source_1',
      name: 'Sample source 1',
      isDefault: true,
      includedModels: [],
      excludedModels: [],
    }),
    mockYamlCatalogSourceConfig({
      id: 'source_2',
      name: 'Source 2',
      isDefault: false,
      includedModels: ['model1', 'model2'],
      excludedModels: ['model3'],
      enabled: false,
    }),
    mockHuggingFaceCatalogSourceConfig({
      id: 'huggingface_source_3',
      name: 'Huggingface source 3',
      allowedOrganization: 'org1',
      isDefault: false,
    }),
    mockYamlCatalogSourceConfig({
      id: 'sample_source_4',
      name: 'Sample source 4',
      isDefault: false,
      includedModels: ['model1', 'model2'],
      excludedModels: ['model3'],
    }),
  ],
  ...partial,
});

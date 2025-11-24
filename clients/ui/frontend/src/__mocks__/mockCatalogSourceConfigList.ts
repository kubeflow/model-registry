import {
  CatalogSourceConfig,
  CatalogSourceConfigList,
  YamlCatalogSourceConfig,
  HuggingFaceCatalogSourceConfig,
  CatalogSourceType,
} from '~/app/modelCatalogTypes';

export const mockYamlCatalogSourceConfig = (
  partial?: Partial<YamlCatalogSourceConfig>,
): YamlCatalogSourceConfig => ({
  id: 'yaml-source-1',
  name: 'Red Hat AI',
  type: CatalogSourceType.YAML,
  enabled: true,
  labels: ['Red Hat AI'],
  includedModels: [],
  excludedModels: [],
  isDefault: true,
  yaml: 'version: 1.0\nmodels:\n  - name: example-model',
  ...partial,
});

export const mockHuggingFaceCatalogSourceConfig = (
  partial?: Partial<HuggingFaceCatalogSourceConfig>,
): HuggingFaceCatalogSourceConfig => ({
  id: 'huggingface-source-1',
  name: 'Huggingface_Admin_1',
  type: CatalogSourceType.HUGGING_FACE,
  enabled: true,
  labels: ['Hugging Face'],
  includedModels: [],
  excludedModels: [],
  isDefault: false,
  allowedOrganization: 'Google',
  apiKey: undefined,
  ...partial,
});

export const mockCatalogSourceConfig = (
  partial?: Partial<CatalogSourceConfig>,
): CatalogSourceConfig => {
  if (partial?.type === CatalogSourceType.HUGGING_FACE) {
    return mockHuggingFaceCatalogSourceConfig(partial as Partial<HuggingFaceCatalogSourceConfig>);
  }
  return mockYamlCatalogSourceConfig(partial as Partial<YamlCatalogSourceConfig>);
};

export const mockCatalogSourceConfigList = (
  partial?: Partial<CatalogSourceConfigList>,
): CatalogSourceConfigList => ({
  catalogs: [
    mockYamlCatalogSourceConfig({ id: 'red-hat-ai', name: 'Red Hat AI', isDefault: true }),
    mockYamlCatalogSourceConfig({
      id: 'red-hat-ai-validated',
      name: 'Red Hat AI validated',
      isDefault: true,
    }),
    mockHuggingFaceCatalogSourceConfig({
      id: 'huggingface-admin-1',
      name: 'Huggingface_Admin_1',
      allowedOrganization: 'Google',
      isDefault: false,
    }),
    mockYamlCatalogSourceConfig({
      id: 'yaml-amdimport-1',
      name: 'YAMLAmdImport_1',
      isDefault: false,
      includedModels: ['model1', 'model2'],
      excludedModels: ['model3'],
    }),
  ],
  ...partial,
});

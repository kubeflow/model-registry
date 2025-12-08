import { mockHuggingFaceCatalogSourceConfig, mockYamlCatalogSourceConfig } from '~/__mocks__';
import { CatalogSourceType } from '~/app/modelCatalogTypes';
import {
  catalogSourceConfigToFormData,
  generateSourceIdFromName,
  transformFormDataToPayload,
} from '~/app/pages/modelCatalogSettings/utils/modelCatalogSettingsUtils';
import { ManageSourceFormData } from '~/app/pages/modelCatalogSettings/useManageSourceData';

const catalogSourceDeafultConfigYAMLMock = mockYamlCatalogSourceConfig({});
const catalogSourceConfigYAMLMock = mockYamlCatalogSourceConfig({ isDefault: false });
const catalogSourceConfigHFMock = mockHuggingFaceCatalogSourceConfig({});

const yamlFormData: ManageSourceFormData = {
  accessToken: '',
  allowedModels: '',
  enabled: true,
  excludedModels: '',
  id: 'sample_source_1',
  isDefault: false,
  name: 'Source 1',
  organization: '',
  sourceType: CatalogSourceType.YAML,
  yamlContent: 'models:\n  - name: model1',
};
const yamlDefaultFormData: ManageSourceFormData = {
  accessToken: '',
  allowedModels: '',
  enabled: true,
  excludedModels: '',
  id: 'sample_source_1',
  isDefault: true,
  name: 'Source 1',
  organization: '',
  sourceType: CatalogSourceType.YAML,
  yamlContent: 'models:\n  - name: model1',
};
const hfFormData: ManageSourceFormData = {
  accessToken: 'apikey',
  allowedModels: '',
  enabled: true,
  excludedModels: '',
  id: 'source_2',
  isDefault: false,
  name: 'Huggingface source 2',
  organization: 'org1',
  sourceType: CatalogSourceType.HUGGING_FACE,
  yamlContent: '',
};

describe('generateSourceIdFromName', () => {
  it('should trim extra spaces', () => {
    expect(generateSourceIdFromName('  testname')).toBe('testname');
  });

  it('should replace - with _', () => {
    expect(generateSourceIdFromName('test-name')).toBe('test_name');
  });

  it('should Remove anything that is NOT alphanumeric and NOT underscore and replace it with _', () => {
    expect(generateSourceIdFromName('Test-Name!')).toBe('test_name');
  });

  it('should convert upper case to lower case', () => {
    expect(generateSourceIdFromName('TestName')).toBe('testname');
  });
});

describe('catalogSourceConfigToFormData', () => {
  it('should convert the data from catalogSourceConfig to formData', () => {
    expect(catalogSourceConfigToFormData(catalogSourceDeafultConfigYAMLMock)).toEqual(
      yamlDefaultFormData,
    );
    expect(catalogSourceConfigToFormData(catalogSourceConfigYAMLMock)).toEqual(yamlFormData);
    expect(catalogSourceConfigToFormData(catalogSourceConfigHFMock)).toEqual(hfFormData);
  });
});

describe('transformFormDataToPayload', () => {
  it('should transform the form data to payload format', () => {
    expect(transformFormDataToPayload(yamlFormData, false)).toEqual({
      enabled: true,
      id: 'sample_source_1',
      name: 'Source 1',
      yaml: 'models:\n  - name: model1',
      isDefault: false,
      type: CatalogSourceType.YAML,
      excludedModels: [],
      includedModels: [],
    });

    expect(transformFormDataToPayload(yamlFormData, true)).toEqual({
      enabled: true,
      name: 'Source 1',
      isDefault: false,
      type: CatalogSourceType.YAML,
      yaml: 'models:\n  - name: model1',
      excludedModels: [],
      includedModels: [],
    });

    expect(transformFormDataToPayload(yamlDefaultFormData, true)).toEqual({
      enabled: true,
      excludedModels: [],
      includedModels: [],
    });

    expect(transformFormDataToPayload(yamlDefaultFormData, false)).toEqual({
      enabled: true,
      excludedModels: [],
      includedModels: [],
    });

    expect(transformFormDataToPayload(hfFormData, false)).toEqual({
      allowedOrganization: 'org1',
      apiKey: 'apikey',
      id: 'source_2',
      enabled: true,
      excludedModels: [],
      includedModels: [],
      isDefault: false,
      name: 'Huggingface source 2',
      type: 'huggingface',
    });
  });
});

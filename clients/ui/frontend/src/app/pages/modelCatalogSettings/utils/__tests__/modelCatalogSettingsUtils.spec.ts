import { mockHuggingFaceCatalogSourceConfig, mockYamlCatalogSourceConfig } from '~/__mocks__';
import { CatalogSourceType } from '~/app/modelCatalogTypes';
import {
  catalogSourceConfigToFormData,
  generateSourceIdFromName,
  transformFormDataToPayload,
} from '~/app/pages/modelCatalogSettings/utils/modelCatalogSettingsUtils';
import { ManageSourceFormData } from '~/app/pages/modelCatalogSettings/useManageSourceData';

const catalogSourceConfigYAMLMock = mockYamlCatalogSourceConfig({});
const catalogSourceConfigHFMock = mockHuggingFaceCatalogSourceConfig({});

const yamlFormData: ManageSourceFormData = {
  accessToken: '',
  allowedModels: '',
  enabled: true,
  excludedModels: '',
  id: 'sample-source-1',
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
  id: 'source-2',
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
    expect(catalogSourceConfigToFormData(catalogSourceConfigYAMLMock)).toEqual(yamlFormData);

    expect(catalogSourceConfigToFormData(catalogSourceConfigHFMock)).toEqual(hfFormData);
  });
});

describe('transformFormDataToPayload', () => {
  it('should transform the form data to payload format', () => {
    expect(transformFormDataToPayload(yamlFormData)).toEqual(mockYamlCatalogSourceConfig({}));
    expect(transformFormDataToPayload(hfFormData)).toEqual(mockHuggingFaceCatalogSourceConfig({}));
  });
});

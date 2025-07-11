/* eslint-disable camelcase */
import {
  ModelRegistryCustomProperties,
  ModelRegistryMetadataType,
  ModelRegistryStringCustomProperties,
} from '~/app/types';
import { getProperties } from '~/app/pages/modelRegistry/screens/utils';

describe('getProperties', () => {
  it('should return an empty object when customProperties is empty', () => {
    const customProperties: ModelRegistryCustomProperties = {};
    const result = getProperties(customProperties);
    expect(result).toEqual({});
  });

  it('should return a filtered object including only string properties with a non-empty value', () => {
    const customProperties: ModelRegistryCustomProperties = {
      property1: { metadataType: ModelRegistryMetadataType.STRING, string_value: 'non-empty' },
      property2: {
        metadataType: ModelRegistryMetadataType.STRING,
        string_value: 'another-non-empty',
      },
      label1: { metadataType: ModelRegistryMetadataType.STRING, string_value: '' },
      label2: { metadataType: ModelRegistryMetadataType.STRING, string_value: '' },
      int1: { metadataType: ModelRegistryMetadataType.INT, int_value: '1' },
      int2: { metadataType: ModelRegistryMetadataType.INT, int_value: '2' },
    };
    const result = getProperties(customProperties);
    expect(result).toEqual({
      property1: { metadataType: ModelRegistryMetadataType.STRING, string_value: 'non-empty' },
      property2: {
        metadataType: ModelRegistryMetadataType.STRING,
        string_value: 'another-non-empty',
      },
    } satisfies ModelRegistryStringCustomProperties);
  });

  it('should return an empty object when all values in customProperties are empty strings or non-string values', () => {
    const customProperties: ModelRegistryCustomProperties = {
      label1: { metadataType: ModelRegistryMetadataType.STRING, string_value: '' },
      label2: { metadataType: ModelRegistryMetadataType.STRING, string_value: '' },
      int1: { metadataType: ModelRegistryMetadataType.INT, int_value: '1' },
      int2: { metadataType: ModelRegistryMetadataType.INT, int_value: '2' },
    };
    const result = getProperties(customProperties);
    expect(result).toEqual({});
  });

  it('should return with _lastModified and _registeredFrom props filtered out', () => {
    const customProperties: ModelRegistryCustomProperties = {
      property1: { metadataType: ModelRegistryMetadataType.STRING, string_value: 'non-empty' },
      _lastModified: { metadataType: ModelRegistryMetadataType.STRING, string_value: 'non-empty' },
      _registeredFromSomething: {
        metadataType: ModelRegistryMetadataType.STRING,
        string_value: 'non-empty',
      },
    };
    const result = getProperties(customProperties);
    expect(result).toEqual({
      property1: { metadataType: ModelRegistryMetadataType.STRING, string_value: 'non-empty' },
    });
  });
});

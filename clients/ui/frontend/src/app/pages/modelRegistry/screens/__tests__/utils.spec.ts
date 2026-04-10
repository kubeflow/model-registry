/* eslint-disable camelcase */
import { mockModelVersion } from '~/__mocks__/mockModelVersion';
import { mockRegisteredModel } from '~/__mocks__/mockRegisteredModel';
import {
  ModelRegistryCustomProperties,
  ModelRegistryCustomPropertyDouble,
  ModelRegistryCustomPropertyInt,
  ModelRegistryCustomPropertyString,
  ModelRegistryMetadataType,
  RegisteredModel,
  ModelVersion,
  ModelState,
} from '~/app/types';
import { mockModelArtifact } from '~/__mocks__/mockModelArtifact';
import { ModelSourceKind } from '~/concepts/modelRegistry/types';
import { modelSourcePropertiesToCatalogParams } from '~/concepts/modelRegistry/utils';
import {
  filterModelVersions,
  getLabels,
  getProperties,
  getPropertyValue,
  mergeUpdatedProperty,
  mergeUpdatedLabels,
  filterRegisteredModels,
  sortModelVersionsByCreateTime,
  isValidHttpUrl,
  getCustomPropString,
  isCompanyUri,
} from '~/app/pages/modelRegistry/screens/utils';
import { COMPANY_URI } from '~/app/utilities/const';
import {
  ModelRegistryFilterDataType,
  ModelRegistryVersionsFilterDataType,
} from '~/app/pages/modelRegistry/screens/const';

describe('getLabels', () => {
  it('should return an empty array when customProperties is empty', () => {
    const customProperties: ModelRegistryCustomProperties = {};
    const result = getLabels(customProperties);
    expect(result).toEqual([]);
  });

  it('should return an array of keys with empty string values in customProperties', () => {
    const customProperties: ModelRegistryCustomProperties = {
      label1: { metadataType: ModelRegistryMetadataType.STRING, string_value: '' },
      label2: { metadataType: ModelRegistryMetadataType.STRING, string_value: 'non-empty' },
      label3: { metadataType: ModelRegistryMetadataType.STRING, string_value: '' },
    };
    const result = getLabels(customProperties);
    expect(result).toEqual(['label1', 'label3']);
  });

  it('should return an empty array when all values in customProperties are non-empty strings', () => {
    const customProperties: ModelRegistryCustomProperties = {
      label1: { metadataType: ModelRegistryMetadataType.STRING, string_value: 'non-empty' },
      label2: { metadataType: ModelRegistryMetadataType.STRING, string_value: 'another-non-empty' },
    };
    const result = getLabels(customProperties);
    expect(result).toEqual([]);
  });
});

describe('mergeUpdatedLabels', () => {
  it('should return an empty object when customProperties and updatedLabels are empty', () => {
    const customProperties: ModelRegistryCustomProperties = {};
    const result = mergeUpdatedLabels(customProperties, []);
    expect(result).toEqual({});
  });

  it('should return an unmodified object if updatedLabels match existing labels', () => {
    const customProperties: ModelRegistryCustomProperties = {
      someUnrelatedProp: { string_value: 'foo', metadataType: ModelRegistryMetadataType.STRING },
      label1: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
    };
    const result = mergeUpdatedLabels(customProperties, ['label1']);
    expect(result).toEqual(customProperties);
  });

  it('should return an object with labels added', () => {
    const customProperties: ModelRegistryCustomProperties = {};
    const result = mergeUpdatedLabels(customProperties, ['label1', 'label2']);
    expect(result).toEqual({
      label1: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
      label2: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
    } satisfies ModelRegistryCustomProperties);
  });

  it('should return an object with labels removed', () => {
    const customProperties: ModelRegistryCustomProperties = {
      label1: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
      label2: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
      label3: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
      label4: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
    };
    const result = mergeUpdatedLabels(customProperties, ['label2', 'label4']);
    expect(result).toEqual({
      label2: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
      label4: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
    } satisfies ModelRegistryCustomProperties);
  });

  it('should return an object with labels both added and removed', () => {
    const customProperties: ModelRegistryCustomProperties = {
      label1: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
      label2: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
      label3: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
    };
    const result = mergeUpdatedLabels(customProperties, ['label1', 'label3', 'label4']);
    expect(result).toEqual({
      label1: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
      label3: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
      label4: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
    } satisfies ModelRegistryCustomProperties);
  });

  it('should not affect non-label properties on the object', () => {
    const customProperties: ModelRegistryCustomProperties = {
      someUnrelatedStrProp: { string_value: 'foo', metadataType: ModelRegistryMetadataType.STRING },
      someUnrelatedIntProp: { int_value: '3', metadataType: ModelRegistryMetadataType.INT },
      label1: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
      label2: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
    };
    const result = mergeUpdatedLabels(customProperties, ['label2', 'label3']);
    expect(result).toEqual({
      someUnrelatedStrProp: { string_value: 'foo', metadataType: ModelRegistryMetadataType.STRING },
      someUnrelatedIntProp: { int_value: '3', metadataType: ModelRegistryMetadataType.INT },
      label2: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
      label3: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
    } satisfies ModelRegistryCustomProperties);
  });
});

describe('getProperties', () => {
  it('should return an empty object when customProperties is empty', () => {
    const customProperties: ModelRegistryCustomProperties = {};
    const result = getProperties(customProperties);
    expect(result).toEqual({});
  });

  it('should return a filtered object including string, int, and double properties', () => {
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
      double1: { metadataType: ModelRegistryMetadataType.DOUBLE, double_value: 3.14 },
    };
    const result = getProperties(customProperties);
    expect(result).toEqual({
      property1: { metadataType: ModelRegistryMetadataType.STRING, string_value: 'non-empty' },
      property2: {
        metadataType: ModelRegistryMetadataType.STRING,
        string_value: 'another-non-empty',
      },
      int1: { metadataType: ModelRegistryMetadataType.INT, int_value: '1' },
      int2: { metadataType: ModelRegistryMetadataType.INT, int_value: '2' },
      double1: { metadataType: ModelRegistryMetadataType.DOUBLE, double_value: 3.14 },
    });
  });

  it('should return numeric properties when only numeric values exist', () => {
    const customProperties: ModelRegistryCustomProperties = {
      label1: { metadataType: ModelRegistryMetadataType.STRING, string_value: '' },
      label2: { metadataType: ModelRegistryMetadataType.STRING, string_value: '' },
      int1: { metadataType: ModelRegistryMetadataType.INT, int_value: '1' },
      int2: { metadataType: ModelRegistryMetadataType.INT, int_value: '2' },
    };
    const result = getProperties(customProperties);
    expect(result).toEqual({
      int1: { metadataType: ModelRegistryMetadataType.INT, int_value: '1' },
      int2: { metadataType: ModelRegistryMetadataType.INT, int_value: '2' },
    });
  });

  it('should return with _lastModified, _registeredFrom, and model_type props filtered out', () => {
    const customProperties: ModelRegistryCustomProperties = {
      property1: { metadataType: ModelRegistryMetadataType.STRING, string_value: 'non-empty' },
      _lastModified: { metadataType: ModelRegistryMetadataType.STRING, string_value: 'non-empty' },
      model_type: { metadataType: ModelRegistryMetadataType.STRING, string_value: 'predictive' },
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

describe('getCustomPropString', () => {
  it('should return the string value of a custom property', () => {
    const customProperties: ModelRegistryCustomProperties = {
      property1: { metadataType: ModelRegistryMetadataType.STRING, string_value: 'prop1' },
      property2: { metadataType: ModelRegistryMetadataType.STRING, string_value: 'prop2' },
    };
    const prop1Result = getCustomPropString(customProperties, 'property1');
    const prop2Result = getCustomPropString(customProperties, 'property2');
    expect(prop1Result).toEqual('prop1');
    expect(prop2Result).toEqual('prop2');
  });
});

describe('getPropertyValue', () => {
  it('should extract string value from STRING property', () => {
    const prop = {
      metadataType: ModelRegistryMetadataType.STRING,
      string_value: 'test value',
    } satisfies ModelRegistryCustomPropertyString;
    expect(getPropertyValue(prop)).toBe('test value');
  });

  it('should extract int value from INT property', () => {
    const prop = {
      metadataType: ModelRegistryMetadataType.INT,
      int_value: '42',
    } satisfies ModelRegistryCustomPropertyInt;
    expect(getPropertyValue(prop)).toBe('42');
  });

  it('should extract and convert double value from DOUBLE property', () => {
    const prop = {
      metadataType: ModelRegistryMetadataType.DOUBLE,
      double_value: 3.14159,
    } satisfies ModelRegistryCustomPropertyDouble;
    expect(getPropertyValue(prop)).toBe('3.14159');
  });

  it('should handle negative integer values', () => {
    const prop = {
      metadataType: ModelRegistryMetadataType.INT,
      int_value: '-100',
    } satisfies ModelRegistryCustomPropertyInt;
    expect(getPropertyValue(prop)).toBe('-100');
  });

  it('should handle negative double values', () => {
    const prop = {
      metadataType: ModelRegistryMetadataType.DOUBLE,
      double_value: -0.5,
    } satisfies ModelRegistryCustomPropertyDouble;
    expect(getPropertyValue(prop)).toBe('-0.5');
  });

  it('should handle empty string value', () => {
    const prop = {
      metadataType: ModelRegistryMetadataType.STRING,
      string_value: '',
    } satisfies ModelRegistryCustomPropertyString;
    expect(getPropertyValue(prop)).toBe('');
  });
});

describe('getCatalogModelDetailsProps', () => {
  it('should return a CatalogModelDetailsParams object from top-level properties when available', () => {
    const modelArtifact = mockModelArtifact({
      modelSourceKind: ModelSourceKind.CATALOG,
      modelSourceClass: 'sourceClass',
      modelSourceName: 'sourceName',
    });
    const result = modelSourcePropertiesToCatalogParams(modelArtifact);
    expect(result).toEqual({
      sourceId: 'sourceClass',
      modelName: 'sourceName',
    });
  });
});

describe('mergeUpdatedProperty', () => {
  it('should handle the create operation', () => {
    const customProperties: ModelRegistryCustomProperties = {
      label1: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
      prop1: { string_value: 'val1', metadataType: ModelRegistryMetadataType.STRING },
    };
    const result = mergeUpdatedProperty({
      customProperties,
      op: 'create',
      newPair: { key: 'prop2', value: 'val2' },
    });
    expect(result).toEqual({
      label1: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
      prop1: { string_value: 'val1', metadataType: ModelRegistryMetadataType.STRING },
      prop2: { string_value: 'val2', metadataType: ModelRegistryMetadataType.STRING },
    } satisfies ModelRegistryCustomProperties);
  });

  it('should handle the update operation without a key change', () => {
    const customProperties: ModelRegistryCustomProperties = {
      label1: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
      prop1: { string_value: 'val1', metadataType: ModelRegistryMetadataType.STRING },
    };
    const result = mergeUpdatedProperty({
      customProperties,
      op: 'update',
      oldKey: 'prop1',
      newPair: { key: 'prop1', value: 'updatedVal1' },
    });
    expect(result).toEqual({
      label1: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
      prop1: { string_value: 'updatedVal1', metadataType: ModelRegistryMetadataType.STRING },
    } satisfies ModelRegistryCustomProperties);
  });

  it('should handle the update operation with a key change', () => {
    const customProperties: ModelRegistryCustomProperties = {
      label1: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
      prop1: { string_value: 'val1', metadataType: ModelRegistryMetadataType.STRING },
    };
    const result = mergeUpdatedProperty({
      customProperties,
      op: 'update',
      oldKey: 'prop1',
      newPair: { key: 'prop2', value: 'val2' },
    });
    expect(result).toEqual({
      label1: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
      prop2: { string_value: 'val2', metadataType: ModelRegistryMetadataType.STRING },
    } satisfies ModelRegistryCustomProperties);
  });

  it('should perform a create if using the update operation with an invalid oldKey', () => {
    const customProperties: ModelRegistryCustomProperties = {
      label1: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
      prop1: { string_value: 'val1', metadataType: ModelRegistryMetadataType.STRING },
    };
    const result = mergeUpdatedProperty({
      customProperties,
      op: 'update',
      oldKey: 'prop2',
      newPair: { key: 'prop3', value: 'val3' },
    });
    expect(result).toEqual({
      label1: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
      prop1: { string_value: 'val1', metadataType: ModelRegistryMetadataType.STRING },
      prop3: { string_value: 'val3', metadataType: ModelRegistryMetadataType.STRING },
    } satisfies ModelRegistryCustomProperties);
  });

  it('should handle the delete operation', () => {
    const customProperties: ModelRegistryCustomProperties = {
      label1: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
      prop1: { string_value: 'val1', metadataType: ModelRegistryMetadataType.STRING },
      prop2: { string_value: 'val2', metadataType: ModelRegistryMetadataType.STRING },
    };
    const result = mergeUpdatedProperty({
      customProperties,
      op: 'delete',
      oldKey: 'prop2',
    });
    expect(result).toEqual({
      label1: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
      prop1: { string_value: 'val1', metadataType: ModelRegistryMetadataType.STRING },
    } satisfies ModelRegistryCustomProperties);
  });

  it('should do nothing if using the delete operation with an invalid oldKey', () => {
    const customProperties: ModelRegistryCustomProperties = {
      label1: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
      prop1: { string_value: 'val1', metadataType: ModelRegistryMetadataType.STRING },
    };
    const result = mergeUpdatedProperty({
      customProperties,
      op: 'delete',
      oldKey: 'prop2',
    });
    expect(result).toEqual({
      label1: { string_value: '', metadataType: ModelRegistryMetadataType.STRING },
      prop1: { string_value: 'val1', metadataType: ModelRegistryMetadataType.STRING },
    } satisfies ModelRegistryCustomProperties);
  });

  describe('type detection', () => {
    it('should detect and create INT property for integer values', () => {
      const customProperties: ModelRegistryCustomProperties = {};
      const result = mergeUpdatedProperty({
        customProperties,
        op: 'create',
        newPair: { key: 'count', value: '42' },
      });
      expect(result).toEqual({
        count: { int_value: '42', metadataType: ModelRegistryMetadataType.INT },
      } satisfies ModelRegistryCustomProperties);
    });

    it('should detect and create INT property for negative integers', () => {
      const customProperties: ModelRegistryCustomProperties = {};
      const result = mergeUpdatedProperty({
        customProperties,
        op: 'create',
        newPair: { key: 'offset', value: '-100' },
      });
      expect(result).toEqual({
        offset: { int_value: '-100', metadataType: ModelRegistryMetadataType.INT },
      } satisfies ModelRegistryCustomProperties);
    });

    it('should detect and create DOUBLE property for decimal values', () => {
      const customProperties: ModelRegistryCustomProperties = {};
      const result = mergeUpdatedProperty({
        customProperties,
        op: 'create',
        newPair: { key: 'accuracy', value: '0.95' },
      });
      expect(result).toEqual({
        accuracy: { double_value: 0.95, metadataType: ModelRegistryMetadataType.DOUBLE },
      } satisfies ModelRegistryCustomProperties);
    });

    it('should detect and create DOUBLE property for negative decimals', () => {
      const customProperties: ModelRegistryCustomProperties = {};
      const result = mergeUpdatedProperty({
        customProperties,
        op: 'create',
        newPair: { key: 'loss', value: '-0.5' },
      });
      expect(result).toEqual({
        loss: { double_value: -0.5, metadataType: ModelRegistryMetadataType.DOUBLE },
      } satisfies ModelRegistryCustomProperties);
    });

    it('should create STRING property for non-numeric values', () => {
      const customProperties: ModelRegistryCustomProperties = {};
      const result = mergeUpdatedProperty({
        customProperties,
        op: 'create',
        newPair: { key: 'model_name', value: 'bert' },
      });
      expect(result).toEqual({
        model_name: { string_value: 'bert', metadataType: ModelRegistryMetadataType.STRING },
      } satisfies ModelRegistryCustomProperties);
    });

    it('should create STRING property for scientific notation', () => {
      const customProperties: ModelRegistryCustomProperties = {};
      const result = mergeUpdatedProperty({
        customProperties,
        op: 'create',
        newPair: { key: 'learning_rate', value: '1e-5' },
      });
      expect(result).toEqual({
        learning_rate: { string_value: '1e-5', metadataType: ModelRegistryMetadataType.STRING },
      } satisfies ModelRegistryCustomProperties);
    });

    it('should create STRING property for values with leading zeros', () => {
      const customProperties: ModelRegistryCustomProperties = {};
      const result = mergeUpdatedProperty({
        customProperties,
        op: 'create',
        newPair: { key: 'version', value: '007' },
      });
      expect(result).toEqual({
        version: { string_value: '007', metadataType: ModelRegistryMetadataType.STRING },
      } satisfies ModelRegistryCustomProperties);
    });
  });

  describe('type transitions', () => {
    it('should update INT to DOUBLE when value changes to decimal', () => {
      const customProperties: ModelRegistryCustomProperties = {
        metric: { int_value: '42', metadataType: ModelRegistryMetadataType.INT },
      };
      const result = mergeUpdatedProperty({
        customProperties,
        op: 'update',
        oldKey: 'metric',
        newPair: { key: 'metric', value: '42.5' },
      });
      expect(result).toEqual({
        metric: { double_value: 42.5, metadataType: ModelRegistryMetadataType.DOUBLE },
      } satisfies ModelRegistryCustomProperties);
    });

    it('should update DOUBLE to INT when value changes to integer', () => {
      const customProperties: ModelRegistryCustomProperties = {
        metric: { double_value: 42.5, metadataType: ModelRegistryMetadataType.DOUBLE },
      };
      const result = mergeUpdatedProperty({
        customProperties,
        op: 'update',
        oldKey: 'metric',
        newPair: { key: 'metric', value: '100' },
      });
      expect(result).toEqual({
        metric: { int_value: '100', metadataType: ModelRegistryMetadataType.INT },
      } satisfies ModelRegistryCustomProperties);
    });

    it('should update INT to STRING when value changes to non-numeric', () => {
      const customProperties: ModelRegistryCustomProperties = {
        field: { int_value: '42', metadataType: ModelRegistryMetadataType.INT },
      };
      const result = mergeUpdatedProperty({
        customProperties,
        op: 'update',
        oldKey: 'field',
        newPair: { key: 'field', value: 'abc' },
      });
      expect(result).toEqual({
        field: { string_value: 'abc', metadataType: ModelRegistryMetadataType.STRING },
      } satisfies ModelRegistryCustomProperties);
    });

    it('should update STRING to INT when value changes to integer', () => {
      const customProperties: ModelRegistryCustomProperties = {
        field: { string_value: 'abc', metadataType: ModelRegistryMetadataType.STRING },
      };
      const result = mergeUpdatedProperty({
        customProperties,
        op: 'update',
        oldKey: 'field',
        newPair: { key: 'field', value: '42' },
      });
      expect(result).toEqual({
        field: { int_value: '42', metadataType: ModelRegistryMetadataType.INT },
      } satisfies ModelRegistryCustomProperties);
    });
  });
});

describe('filterModelVersions', () => {
  const modelVersions: ModelVersion[] = [
    mockModelVersion({ name: 'Test 1', state: ModelState.ARCHIVED }),
    mockModelVersion({
      name: 'Test 2',
      description: 'Description2',
    }),
    mockModelVersion({ name: 'Test 3', author: 'Author3', state: ModelState.ARCHIVED }),
    mockModelVersion({ name: 'Test 4', state: ModelState.ARCHIVED }),
    mockModelVersion({ name: 'Test 5' }),
  ];

  test('filters by name', () => {
    const filtered = filterModelVersions(modelVersions, {
      Keyword: 'Test 1',
      Author: '',
    } satisfies ModelRegistryVersionsFilterDataType);
    expect(filtered).toEqual([modelVersions[0]]);
  });

  test('filters by description', () => {
    const filtered = filterModelVersions(modelVersions, {
      Keyword: 'Description2',
      Author: '',
    } satisfies ModelRegistryVersionsFilterDataType);
    expect(filtered).toEqual([modelVersions[1]]);
  });

  test('filters by author', () => {
    const filtered = filterModelVersions(modelVersions, {
      Keyword: '',
      Author: 'Author3',
    } satisfies ModelRegistryVersionsFilterDataType);
    expect(filtered).toEqual([modelVersions[2]]);
  });

  test('filters by keyword and author', () => {
    const filtered = filterModelVersions(modelVersions, {
      Keyword: 'Test 3',
      Author: 'Author3',
    } satisfies ModelRegistryVersionsFilterDataType);
    expect(filtered).toEqual([modelVersions[2]]);
  });

  test('does not filter when search is empty', () => {
    const filtered = filterModelVersions(modelVersions, {
      Keyword: '',
      Author: '',
    } satisfies ModelRegistryVersionsFilterDataType);
    expect(filtered).toEqual(modelVersions);
  });

  test('does not filter when keyword is correct but author is incorrect', () => {
    const filtered = filterModelVersions(modelVersions, {
      Keyword: 'Test 1',
      Author: 'Bob',
    } satisfies ModelRegistryVersionsFilterDataType);
    expect(filtered).toEqual([]);
  });

  test('does not filter when keyword is incorrect but author is correct', () => {
    const filtered = filterModelVersions(modelVersions, {
      Keyword: 'Test 6',
      Author: 'Author3',
    } satisfies ModelRegistryVersionsFilterDataType);
    expect(filtered).toEqual([]);
  });
});

describe('filterRegisteredModels', () => {
  const registeredModels: RegisteredModel[] = [
    mockRegisteredModel({ name: 'Test 1', state: ModelState.ARCHIVED, owner: 'Alice' }),
    mockRegisteredModel({
      name: 'Test 2',
      description: 'Description2',
      owner: 'Bob',
    }),
    mockRegisteredModel({ name: 'Test 3', state: ModelState.ARCHIVED, owner: 'Charlie' }),
    mockRegisteredModel({ name: 'Test 4', state: ModelState.ARCHIVED, owner: 'Alice' }),
    mockRegisteredModel({ name: 'Test 5', owner: 'Bob' }),
  ];

  test('filters by name', () => {
    const filtered = filterRegisteredModels(registeredModels, [], {
      Keyword: 'Test 1',
      Owner: '',
    } satisfies ModelRegistryFilterDataType);
    expect(filtered).toEqual([registeredModels[0]]);
  });

  test('filters by description', () => {
    const filtered = filterRegisteredModels(registeredModels, [], {
      Keyword: 'Description2',
      Owner: '',
    } satisfies ModelRegistryFilterDataType);
    expect(filtered).toEqual([registeredModels[1]]);
  });

  test('filters by owner', () => {
    const filtered = filterRegisteredModels(registeredModels, [], {
      Keyword: '',
      Owner: 'Alice',
    } satisfies ModelRegistryFilterDataType);
    expect(filtered).toEqual([registeredModels[0], registeredModels[3]]);
  });

  test('filters by keyword and owner', () => {
    const filtered = filterRegisteredModels(registeredModels, [], {
      Keyword: 'Test 1',
      Owner: 'Alice',
    } satisfies ModelRegistryFilterDataType);
    expect(filtered).toEqual([registeredModels[0]]);
  });

  test('does not filter when search is empty', () => {
    const filtered = filterRegisteredModels(registeredModels, [], {
      Keyword: '',
      Owner: '',
    } satisfies ModelRegistryFilterDataType);
    expect(filtered).toEqual(registeredModels);
  });

  test('does not filter when keyword is correct but owner is incorrect', () => {
    const filtered = filterRegisteredModels(registeredModels, [], {
      Keyword: 'Test 1',
      Owner: 'Bob',
    } satisfies ModelRegistryFilterDataType);
    expect(filtered).toEqual([]);
  });

  test('does not filter when keyword is incorrect but owner is correct', () => {
    const filtered = filterRegisteredModels(registeredModels, [], {
      Keyword: 'Test 6',
      Owner: 'Alice',
    } satisfies ModelRegistryFilterDataType);
    expect(filtered).toEqual([]);
  });
});

describe('sortModelVersionsByCreateTime', () => {
  it('should return list of sorted modelVersions by create time', () => {
    const modelVersions: ModelVersion[] = [
      mockModelVersion({
        name: 'model version 1',
        author: 'Author 1',
        id: '1',
        createTimeSinceEpoch: '1725018764650',
        lastUpdateTimeSinceEpoch: '1725030215299',
      }),
      mockModelVersion({
        name: 'model version 1',
        author: 'Author 1',
        id: '1',
        createTimeSinceEpoch: '1725028468207',
        lastUpdateTimeSinceEpoch: '1725030142332',
      }),
    ];

    const result = sortModelVersionsByCreateTime(modelVersions);
    expect(result).toEqual([modelVersions[1], modelVersions[0]]);
  });
});

describe('isValidHttpUrl', () => {
  it('should return true for a valid HTTPS URL', () => {
    expect(isValidHttpUrl('https://example.com')).toBe(true);
  });

  it('should return true for a valid HTTP URL', () => {
    expect(isValidHttpUrl('http://example.com')).toBe(true);
  });

  it('should return false for a URL with an unsupported protocol', () => {
    expect(isValidHttpUrl('ftp://example.com')).toBe(false);
  });

  it('should return false for an invalid URL string', () => {
    expect(isValidHttpUrl('random text')).toBe(false);
  });

  it('should return false for an empty string', () => {
    expect(isValidHttpUrl('')).toBe(false);
  });

  it('should return false for a string without a protocol', () => {
    expect(isValidHttpUrl('www.example.com')).toBe(false);
  });

  it('should return false for null input', () => {
    expect(isValidHttpUrl(null as unknown as string)).toBe(false);
  });

  it('should return false for undefined input', () => {
    expect(isValidHttpUrl(undefined as unknown as string)).toBe(false);
  });

  it('should return false for a number input', () => {
    expect(isValidHttpUrl(12345 as unknown as string)).toBe(false);
  });

  it('should return false for a URL with a missing domain', () => {
    expect(isValidHttpUrl('http://')).toBe(false);
  });

  it('should return false for a URL with an invalid format', () => {
    expect(isValidHttpUrl('http://example..com')).toBe(false);
  });
});

describe('isCompanyUri', () => {
  it('should return true for company registry URI', () => {
    expect(isCompanyUri(`${COMPANY_URI}/test/test`)).toBe(true);
  });

  it('should return false for non-company registry URI', () => {
    expect(isCompanyUri(`${COMPANY_URI}1/test/test`)).toBe(false);
  });
});

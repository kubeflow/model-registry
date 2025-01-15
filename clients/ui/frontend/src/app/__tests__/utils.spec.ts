import { mockModelVersion } from '~/__mocks__/mockModelVersion';
import { mockRegisteredModel } from '~/__mocks__/mockRegisteredModel';
import {
  filterArchiveModels,
  filterArchiveVersions,
  filterLiveModels,
  filterLiveVersions,
  getLastCreatedItem,
  ObjectStorageFields,
  objectStorageFieldsToUri,
  uriToObjectStorageFields,
} from '~/app/utils';
import { RegisteredModel, ModelState, ModelVersion } from '~/app/types';

describe('objectStorageFieldsToUri', () => {
  it('converts fields to URI with all fields present', () => {
    const uri = objectStorageFieldsToUri({
      endpoint: 'http://s3.amazonaws.com/',
      bucket: 'test-bucket',
      region: 'us-east-1',
      path: 'demo-models/flan-t5-small-caikit',
    });
    expect(uri).toEqual(
      's3://test-bucket/demo-models/flan-t5-small-caikit?endpoint=http%3A%2F%2Fs3.amazonaws.com%2F&defaultRegion=us-east-1',
    );
  });

  it('converts fields to URI with region missing', () => {
    const uri = objectStorageFieldsToUri({
      endpoint: 'http://s3.amazonaws.com/',
      bucket: 'test-bucket',
      path: 'demo-models/flan-t5-small-caikit',
    });
    expect(uri).toEqual(
      's3://test-bucket/demo-models/flan-t5-small-caikit?endpoint=http%3A%2F%2Fs3.amazonaws.com%2F',
    );
  });

  it('converts fields to URI with region empty', () => {
    const uri = objectStorageFieldsToUri({
      endpoint: 'http://s3.amazonaws.com/',
      bucket: 'test-bucket',
      region: '',
      path: 'demo-models/flan-t5-small-caikit',
    });
    expect(uri).toEqual(
      's3://test-bucket/demo-models/flan-t5-small-caikit?endpoint=http%3A%2F%2Fs3.amazonaws.com%2F',
    );
  });

  it('falls back to null if endpoint is empty', () => {
    const uri = objectStorageFieldsToUri({
      endpoint: '',
      bucket: 'test-bucket',
      region: 'us-east-1',
      path: 'demo-models/flan-t5-small-caikit',
    });
    expect(uri).toEqual(null);
  });

  it('falls back to null if bucket is empty', () => {
    const uri = objectStorageFieldsToUri({
      endpoint: 'http://s3.amazonaws.com/',
      bucket: '',
      region: 'us-east-1',
      path: 'demo-models/flan-t5-small-caikit',
    });
    expect(uri).toEqual(null);
  });

  it('falls back to null if path is empty', () => {
    const uri = objectStorageFieldsToUri({
      endpoint: 'http://s3.amazonaws.com/',
      bucket: 'test-bucket',
      region: 'us-east-1',
      path: '',
    });
    expect(uri).toEqual(null);
  });
});

describe('uriToObjectStorageFields', () => {
  it('converts URI to fields with all params present', () => {
    const fields = uriToObjectStorageFields(
      's3://test-bucket/demo-models/flan-t5-small-caikit?endpoint=http%3A%2F%2Fs3.amazonaws.com%2F&defaultRegion=us-east-1',
    );
    expect(fields).toEqual({
      endpoint: 'http://s3.amazonaws.com/',
      bucket: 'test-bucket',
      region: 'us-east-1',
      path: 'demo-models/flan-t5-small-caikit',
    } satisfies ObjectStorageFields);
  });

  it('converts URI to fields with region missing', () => {
    const fields = uriToObjectStorageFields(
      's3://test-bucket/demo-models/flan-t5-small-caikit?endpoint=http%3A%2F%2Fs3.amazonaws.com%2F',
    );
    expect(fields).toEqual({
      endpoint: 'http://s3.amazonaws.com/',
      bucket: 'test-bucket',
      path: 'demo-models/flan-t5-small-caikit',
      region: undefined,
    } satisfies ObjectStorageFields);
  });

  it('falls back to null if endpoint is missing', () => {
    const fields = uriToObjectStorageFields('s3://test-bucket/demo-models/flan-t5-small-caikit');
    expect(fields).toBeNull();
  });

  it('falls back to null if path is missing', () => {
    const fields = uriToObjectStorageFields(
      's3://test-bucket/?endpoint=http%3A%2F%2Fs3.amazonaws.com%2F&defaultRegion=us-east-1',
    );
    expect(fields).toBeNull();
  });

  it('falls back to null if bucket is missing', () => {
    const fields = uriToObjectStorageFields(
      's3://?endpoint=http%3A%2F%2Fs3.amazonaws.com%2F&defaultRegion=us-east-1',
    );
    expect(fields).toBeNull();
  });

  it('falls back to null if the URI is malformed', () => {
    const fields = uriToObjectStorageFields('test-bucket/demo-models/flan-t5-small-caikit');
    expect(fields).toBeNull();
  });
});

describe('getLastCreatedItem', () => {
  it('returns the latest item correctly', () => {
    const items = [
      {
        foo: 'a',
        createTimeSinceEpoch: '1712234877179', // Apr 04 2024
      },
      {
        foo: 'b',
        createTimeSinceEpoch: '1723659611927', // Aug 14 2024
      },
    ];
    expect(getLastCreatedItem(items)).toBe(items[1]);
  });

  it('returns first item if items have no createTimeSinceEpoch', () => {
    const items = [
      { foo: 'a', createTimeSinceEpoch: undefined },
      { foo: 'b', createTimeSinceEpoch: undefined },
    ];
    expect(getLastCreatedItem(items)).toBe(items[0]);
  });
});

describe('Filter model state', () => {
  const models: RegisteredModel[] = [
    mockRegisteredModel({ name: 'Test 1', state: ModelState.ARCHIVED }),
    mockRegisteredModel({
      name: 'Test 2',
      state: ModelState.LIVE,
      description: 'Description2',
    }),
    mockRegisteredModel({ name: 'Test 3', state: ModelState.ARCHIVED }),
    mockRegisteredModel({ name: 'Test 4', state: ModelState.ARCHIVED }),
    mockRegisteredModel({ name: 'Test 5', state: ModelState.LIVE }),
  ];

  describe('filterArchiveModels', () => {
    it('should filter out only the archived versions', () => {
      const archivedModels = filterArchiveModels(models);
      expect(archivedModels).toEqual([models[0], models[2], models[3]]);
    });

    it('should return an empty array if the input array is empty', () => {
      const result = filterArchiveModels([]);
      expect(result).toEqual([]);
    });
  });

  describe('filterLiveModels', () => {
    it('should filter out only the live models', () => {
      const liveModels = filterLiveModels(models);
      expect(liveModels).toEqual([models[1], models[4]]);
    });

    it('should return an empty array if the input array is empty', () => {
      const result = filterLiveModels([]);
      expect(result).toEqual([]);
    });
  });
});

describe('Filter model version state', () => {
  const modelVersions: ModelVersion[] = [
    mockModelVersion({ name: 'Test 1', state: ModelState.ARCHIVED }),
    mockModelVersion({
      name: 'Test 2',
      state: ModelState.LIVE,
      description: 'Description2',
    }),
    mockModelVersion({ name: 'Test 3', author: 'Author3', state: ModelState.ARCHIVED }),
    mockModelVersion({ name: 'Test 4', state: ModelState.ARCHIVED }),
    mockModelVersion({ name: 'Test 5', state: ModelState.LIVE }),
  ];

  describe('filterArchiveVersions', () => {
    it('should filter out only the archived versions', () => {
      const archivedVersions = filterArchiveVersions(modelVersions);
      expect(archivedVersions).toEqual([modelVersions[0], modelVersions[2], modelVersions[3]]);
    });

    it('should return an empty array if the input array is empty', () => {
      const result = filterArchiveVersions([]);
      expect(result).toEqual([]);
    });
  });

  describe('filterLiveVersions', () => {
    it('should filter out only the live versions', () => {
      const liveVersions = filterLiveVersions(modelVersions);
      expect(liveVersions).toEqual([modelVersions[1], modelVersions[4]]);
    });

    it('should return an empty array if the input array is empty', () => {
      const result = filterLiveVersions([]);
      expect(result).toEqual([]);
    });
  });
});

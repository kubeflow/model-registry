import { restCREATE, handleRestFailures } from 'mod-arch-core';
import { previewCatalogSource } from '~/app/api/modelCatalogSettings/service';

const mockRestPromise = Promise.resolve({ data: {} });

jest.mock('mod-arch-core', () => ({
  restCREATE: jest.fn(() => mockRestPromise),
  assembleModArchBody: jest.fn((data) => data),
  isModArchResponse: jest.fn(() => true),
  handleRestFailures: jest.fn(() => mockRestPromise),
}));

const handleRestFailuresMock = jest.mocked(handleRestFailures);
const restCREATEMock = jest.mocked(restCREATE);

const APIOptionsMock = {};

describe('previewCatalogSource', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should call restCREATE with base query params when no additional params provided', async () => {
    const mockData = {
      type: 'yaml' as const,
      includedModels: ['*'],
      excludedModels: [],
      properties: { yaml: 'models:\n  - name: test' },
    };

    await previewCatalogSource('/api/v1/settings/model_catalog', { namespace: 'kubeflow' })(
      APIOptionsMock,
      mockData,
    );

    expect(restCREATEMock).toHaveBeenCalledTimes(1);
    expect(restCREATEMock).toHaveBeenCalledWith(
      '/api/v1/settings/model_catalog',
      '/source_preview',
      mockData,
      { namespace: 'kubeflow' },
      APIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
  });

  it('should merge additional query params with base query params', async () => {
    const mockData = {
      type: 'yaml' as const,
      includedModels: ['*'],
      excludedModels: [],
      properties: { yaml: 'models:\n  - name: test' },
    };

    await previewCatalogSource('/api/v1/settings/model_catalog', { namespace: 'kubeflow' })(
      APIOptionsMock,
      mockData,
      { filterStatus: 'included', pageSize: 20 },
    );

    expect(restCREATEMock).toHaveBeenCalledTimes(1);
    expect(restCREATEMock).toHaveBeenCalledWith(
      '/api/v1/settings/model_catalog',
      '/source_preview',
      mockData,
      { namespace: 'kubeflow', filterStatus: 'included', pageSize: 20 },
      APIOptionsMock,
    );
  });

  it('should include nextPageToken in query params when provided', async () => {
    const mockData = {
      type: 'yaml' as const,
      includedModels: ['*'],
      excludedModels: [],
      properties: { yaml: 'models:\n  - name: test' },
    };

    await previewCatalogSource('/api/v1/settings/model_catalog', { namespace: 'kubeflow' })(
      APIOptionsMock,
      mockData,
      { filterStatus: 'excluded', pageSize: 10, nextPageToken: 'abc123' },
    );

    expect(restCREATEMock).toHaveBeenCalledWith(
      '/api/v1/settings/model_catalog',
      '/source_preview',
      mockData,
      { namespace: 'kubeflow', filterStatus: 'excluded', pageSize: 10, nextPageToken: 'abc123' },
      APIOptionsMock,
    );
  });

  it('should override base query params with additional params if same key exists', async () => {
    const mockData = {
      type: 'yaml' as const,
      includedModels: [],
      excludedModels: [],
      properties: {},
    };

    // Base has someKey: 'base', additional has someKey: 'override'
    await previewCatalogSource('/api/v1/settings/model_catalog', { someKey: 'base' })(
      APIOptionsMock,
      mockData,
      { filterStatus: 'all' },
    );

    expect(restCREATEMock).toHaveBeenCalledWith(
      '/api/v1/settings/model_catalog',
      '/source_preview',
      mockData,
      { someKey: 'base', filterStatus: 'all' },
      APIOptionsMock,
    );
  });
});

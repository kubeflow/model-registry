import {
  isModArchResponse,
  restCREATE,
  restGET,
  restPATCH,
  handleRestFailures,
} from 'mod-arch-shared';
import { ModelState, ModelArtifactState } from '~/app/types';
import {
  createRegisteredModel,
  getRegisteredModel,
  getModelVersion,
  getListModelVersions,
  getListModelArtifacts,
  getModelVersionsByRegisteredModel,
  getListRegisteredModels,
  patchModelArtifact,
  patchModelVersion,
  patchRegisteredModel,
  getModelArtifactsByModelVersion,
  createModelVersionForRegisteredModel,
  createModelArtifactForModelVersion,
} from '~/app/api/service';
import { mockRegisteredModel } from '~/__mocks__';
import { BFF_API_VERSION } from '~/app/utilities/const';

const mockRestPromise = Promise.resolve({ data: {} });
const mockRestResponse = {};

jest.mock('mod-arch-shared', () => ({
  restCREATE: jest.fn(() => mockRestPromise),
  restGET: jest.fn(() => mockRestPromise),
  restPATCH: jest.fn(() => mockRestPromise),
  assembleModArchBody: jest.fn(() => ({})),
  isModArchResponse: jest.fn(() => true),
  handleRestFailures: jest.fn(() => mockRestPromise),
  asEnumMember: jest.fn((value, enumObj) =>
    value && Object.values(enumObj).includes(value) ? value : undefined,
  ),
  Theme: {
    Default: 'default',
    MUI: 'mui',
  },
  DeploymentMode: {
    Integrated: 'integrated',
    Default: 'default',
  },
  PlatformMode: {
    Kubeflow: 'kubeflow',
    Default: 'default',
  },
}));

const handleRestFailuresMock = jest.mocked(handleRestFailures);
const restCREATEMock = jest.mocked(restCREATE);
const restGETMock = jest.mocked(restGET);
const restPATCHMock = jest.mocked(restPATCH);
const isModelRegistryResponseMock = jest.mocked(isModArchResponse);

const APIOptionsMock = {};

describe('createRegisteredModel', () => {
  it('should call restCREATE and handleRestFailures to create registered model', async () => {
    const mockData = {
      description: 'test',
      externalID: '1',
      name: 'test new registered model',
      state: ModelState.LIVE,
      customProperties: {},
    };
    const response = await createRegisteredModel(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
    )(APIOptionsMock, mockData);
    expect(response).toEqual(mockRestResponse);
    expect(restCREATEMock).toHaveBeenCalledTimes(1);
    expect(isModelRegistryResponseMock).toHaveBeenCalledTimes(1);
    expect(restCREATEMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/registered_models`,
      {},
      {},
      APIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockRestPromise);
  });
});

describe('createModelVersionForRegisteredModel', () => {
  it('should call restCREATE and handleRestFailures to create model version for a model', async () => {
    const mockData = {
      description: 'test',
      externalID: '1',
      author: 'test author',
      registeredModelId: '1',
      name: 'test new model version',
      state: ModelState.LIVE,
      customProperties: {},
    };
    const response = await createModelVersionForRegisteredModel(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
    )(APIOptionsMock, '1', mockData, mockRegisteredModel({}), true);
    expect(response).toEqual(mockRestResponse);
    expect(restCREATEMock).toHaveBeenCalledTimes(1);
    expect(restCREATEMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/registered_models/1/versions`,
      {},
      {},
      APIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockRestPromise);
  });
});

describe('createModelArtifactForModelVersion', () => {
  it('should call restCREATE and handleRestFailures to create model artifact for version', async () => {
    const mockData = {
      description: 'test',
      externalID: 'test',
      uri: 'test-uri',
      state: ModelArtifactState.LIVE,
      name: 'test-name',
      modelFormatName: 'test-modelformatname',
      storageKey: 'teststoragekey',
      storagePath: 'teststoragePath',
      modelFormatVersion: 'testmodelFormatVersion',
      serviceAccountName: 'testserviceAccountname',
      customProperties: {},
      artifactType: 'model-artifact',
    };
    const response = await createModelArtifactForModelVersion(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
    )(APIOptionsMock, '2', mockData);
    expect(response).toEqual(mockRestResponse);
    expect(restCREATEMock).toHaveBeenCalledTimes(1);
    expect(restCREATEMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/model_versions/2/artifacts`,
      {},
      {},
      APIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockRestPromise);
  });
});

describe('getRegisteredModel', () => {
  it('should call restGET and handleRestFailures to fetch registered model', async () => {
    const response = await getRegisteredModel(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
    )(APIOptionsMock, '1');
    expect(response).toEqual(mockRestResponse);
    expect(restGETMock).toHaveBeenCalledTimes(1);
    expect(restGETMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/registered_models/1`,
      {},
      APIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockRestPromise);
  });
});

describe('getModelVersion', () => {
  it('should call restGET and handleRestFailures to fetch model version', async () => {
    const response = await getModelVersion(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
    )(APIOptionsMock, '1');
    expect(response).toEqual(mockRestResponse);
    expect(restGETMock).toHaveBeenCalledTimes(1);
    expect(restGETMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/model_versions/1`,
      {},
      APIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockRestPromise);
  });
});

describe('getListRegisteredModels', () => {
  it('should call restGET and handleRestFailures to list registered models', async () => {
    const response = await getListRegisteredModels(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
    )({});
    expect(response).toEqual(mockRestResponse);
    expect(restGETMock).toHaveBeenCalledTimes(1);
    expect(restGETMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/registered_models`,
      {},
      APIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockRestPromise);
  });
});

describe('getListModelArtifacts', () => {
  it('should call restGET and handleRestFailures to list models artifacts', async () => {
    const response = await getListModelArtifacts(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
    )({});
    expect(response).toEqual(mockRestResponse);
    expect(restGETMock).toHaveBeenCalledTimes(1);
    expect(restGETMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/model_artifacts`,
      {},
      APIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockRestPromise);
  });
});

describe('getListModelVersions', () => {
  it('should call restGET and handleRestFailures to list models versions', async () => {
    const response = await getListModelVersions(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
    )({});
    expect(response).toEqual(mockRestResponse);
    expect(restGETMock).toHaveBeenCalledTimes(1);
    expect(restGETMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/model_versions`,
      {},
      APIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockRestPromise);
  });
});

describe('getModelVersionsByRegisteredModel', () => {
  it('should call restGET and handleRestFailures to list models versions by registered model', async () => {
    const response = await getModelVersionsByRegisteredModel(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
    )({}, '1');
    expect(response).toEqual(mockRestResponse);
    expect(restGETMock).toHaveBeenCalledTimes(1);
    expect(restGETMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/registered_models/1/versions`,
      {},
      APIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockRestPromise);
  });
});

describe('getModelArtifactsByModelVersion', () => {
  it('should call restGET and handleRestFailures to list models artifacts by model version', async () => {
    const response = await getModelArtifactsByModelVersion(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
    )({}, '1');
    expect(response).toEqual(mockRestResponse);
    expect(restGETMock).toHaveBeenCalledTimes(1);
    expect(restGETMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/model_versions/1/artifacts`,
      {},
      APIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockRestPromise);
  });
});

describe('patchRegisteredModel', () => {
  it('should call restPATCH and handleRestFailures to update registered model', async () => {
    const mockData = { description: 'new test' };
    const response = await patchRegisteredModel(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      {},
    )(APIOptionsMock, mockData, '1');
    expect(response).toEqual(mockRestResponse);
    expect(restPATCHMock).toHaveBeenCalledTimes(1);
    expect(restPATCHMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/registered_models/1`,
      {},
      {},
      APIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockRestPromise);
  });
});

describe('patchModelVersion', () => {
  it('should call restPATCH and handleRestFailures to update model version', async () => {
    const mockData = { description: 'new test' };
    const response = await patchModelVersion(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      {},
    )(APIOptionsMock, mockData, '1');
    expect(response).toEqual(mockRestResponse);
    expect(restPATCHMock).toHaveBeenCalledTimes(1);
    expect(restPATCHMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/model_versions/1`,
      {},
      {},
      APIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockRestPromise);
  });
});

describe('patchModelArtifact', () => {
  it('should call restPATCH and handleRestFailures to update model artifact', async () => {
    const mockData = { description: 'new test' };
    const response = await patchModelArtifact(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      {},
    )(APIOptionsMock, mockData, '1');
    expect(response).toEqual(mockRestResponse);
    expect(restPATCHMock).toHaveBeenCalledTimes(1);
    expect(restPATCHMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/model_artifacts/1`,
      {},
      {},
      APIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockRestPromise);
  });
});

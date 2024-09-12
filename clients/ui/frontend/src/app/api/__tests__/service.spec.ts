import { restCREATE, restGET, restPATCH } from '~/app/api/apiUtils';
import { handleRestFailures } from '~/app/api/errorUtils';
import { ModelState, ModelArtifactState } from '~/app/types';
import {
  createModelArtifact,
  createModelVersion,
  createRegisteredModel,
  getRegisteredModel,
  getModelVersion,
  getModelArtifact,
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
import { BFF_API_VERSION } from '~/app/const';

const mockProxyPromise = Promise.resolve();

jest.mock('~/app/api/apiUtils', () => ({
  restCREATE: jest.fn(() => mockProxyPromise),
  restGET: jest.fn(() => mockProxyPromise),
  restPATCH: jest.fn(() => mockProxyPromise),
}));

const mockResultPromise = Promise.resolve();

jest.mock('~/app/api/errorUtils', () => ({
  handleRestFailures: jest.fn(() => mockResultPromise),
}));

const handleRestFailuresMock = jest.mocked(handleRestFailures);
const restCREATEMock = jest.mocked(restCREATE);
const restGETMock = jest.mocked(restGET);
const restPATCHMock = jest.mocked(restPATCH);

const K8sAPIOptionsMock = {};

describe('createRegisteredModel', () => {
  it('should call restCREATE and handleRestFailures to create registered model', () => {
    expect(
      createRegisteredModel(`/api/${BFF_API_VERSION}/model_registry/model-registry-1/`)(
        K8sAPIOptionsMock,
        {
          description: 'test',
          externalID: '1',
          name: 'test new registered model',
          state: ModelState.LIVE,
          customProperties: {},
        },
      ),
    ).toBe(mockResultPromise);
    expect(restCREATEMock).toHaveBeenCalledTimes(1);
    expect(restCREATEMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/registered_models`,
      {
        description: 'test',
        externalID: '1',
        name: 'test new registered model',
        state: ModelState.LIVE,
        customProperties: {},
      },
      {},
      K8sAPIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockProxyPromise);
  });
});

describe('createModelVersion', () => {
  it('should call restCREATE and handleRestFailures to create model version', () => {
    expect(
      createModelVersion(`/api/${BFF_API_VERSION}/model_registry/model-registry-1/`)(
        K8sAPIOptionsMock,
        {
          description: 'test',
          externalID: '1',
          author: 'test author',
          registeredModelId: '1',
          name: 'test new model version',
          state: ModelState.LIVE,
          customProperties: {},
        },
      ),
    ).toBe(mockResultPromise);
    expect(restCREATEMock).toHaveBeenCalledTimes(1);
    expect(restCREATEMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/model_versions`,
      {
        description: 'test',
        externalID: '1',
        author: 'test author',
        registeredModelId: '1',
        name: 'test new model version',
        state: ModelState.LIVE,
        customProperties: {},
      },
      {},
      K8sAPIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockProxyPromise);
  });
});

describe('createModelVersionForRegisteredModel', () => {
  it('should call restCREATE and handleRestFailures to create model version for a model', () => {
    expect(
      createModelVersionForRegisteredModel(
        `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      )(K8sAPIOptionsMock, '1', {
        description: 'test',
        externalID: '1',
        author: 'test author',
        registeredModelId: '1',
        name: 'test new model version',
        state: ModelState.LIVE,
        customProperties: {},
      }),
    ).toBe(mockResultPromise);
    expect(restCREATEMock).toHaveBeenCalledTimes(1);
    expect(restCREATEMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/registered_models/1/versions`,
      {
        description: 'test',
        externalID: '1',
        author: 'test author',
        registeredModelId: '1',
        name: 'test new model version',
        state: ModelState.LIVE,
        customProperties: {},
      },
      {},
      K8sAPIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockProxyPromise);
  });
});

describe('createModelArtifact', () => {
  it('should call restCREATE and handleRestFailures to create model artifact', () => {
    expect(
      createModelArtifact(`/api/${BFF_API_VERSION}/model_registry/model-registry-1/`)(
        K8sAPIOptionsMock,
        {
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
        },
      ),
    ).toBe(mockResultPromise);
    expect(restCREATEMock).toHaveBeenCalledTimes(1);
    expect(restCREATEMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/model_artifacts`,
      {
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
      },
      {},
      K8sAPIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockProxyPromise);
  });
});

describe('createModelArtifactForModelVersion', () => {
  it('should call restCREATE and handleRestFailures to create model artifact for version', () => {
    expect(
      createModelArtifactForModelVersion(
        `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      )(K8sAPIOptionsMock, '2', {
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
      }),
    ).toBe(mockResultPromise);
    expect(restCREATEMock).toHaveBeenCalledTimes(1);
    expect(restCREATEMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/model_versions/2/artifacts`,
      {
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
      },
      {},
      K8sAPIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockProxyPromise);
  });
});

describe('getRegisteredModel', () => {
  it('should call restGET and handleRestFailures to fetch registered model', () => {
    expect(
      getRegisteredModel(`/api/${BFF_API_VERSION}/model_registry/model-registry-1/`)(
        K8sAPIOptionsMock,
        '1',
      ),
    ).toBe(mockResultPromise);
    expect(restGETMock).toHaveBeenCalledTimes(1);
    expect(restGETMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/registered_models/1`,
      {},
      K8sAPIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockProxyPromise);
  });
});

describe('getModelVersion', () => {
  it('should call restGET and handleRestFailures to fetch model version', () => {
    expect(
      getModelVersion(`/api/${BFF_API_VERSION}/model_registry/model-registry-1/`)(
        K8sAPIOptionsMock,
        '1',
      ),
    ).toBe(mockResultPromise);
    expect(restGETMock).toHaveBeenCalledTimes(1);
    expect(restGETMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/model_versions/1`,
      {},
      K8sAPIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockProxyPromise);
  });
});

describe('getModelArtifact', () => {
  it('should call restGET and handleRestFailures to fetch model version', () => {
    expect(
      getModelArtifact(`/api/${BFF_API_VERSION}/model_registry/model-registry-1/`)(
        K8sAPIOptionsMock,
        '1',
      ),
    ).toBe(mockResultPromise);
    expect(restGETMock).toHaveBeenCalledTimes(1);
    expect(restGETMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/model_artifacts/1`,
      {},
      K8sAPIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockProxyPromise);
  });
});

describe('getListRegisteredModels', () => {
  it('should call restGET and handleRestFailures to list registered models', () => {
    expect(
      getListRegisteredModels(`/api/${BFF_API_VERSION}/model_registry/model-registry-1/`)({}),
    ).toBe(mockResultPromise);
    expect(restGETMock).toHaveBeenCalledTimes(1);
    expect(restGETMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/registered_models`,
      {},
      K8sAPIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockProxyPromise);
  });
});

describe('getListModelArtifacts', () => {
  it('should call restGET and handleRestFailures to list models artifacts', () => {
    expect(
      getListModelArtifacts(`/api/${BFF_API_VERSION}/model_registry/model-registry-1/`)({}),
    ).toBe(mockResultPromise);
    expect(restGETMock).toHaveBeenCalledTimes(1);
    expect(restGETMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/model_artifacts`,
      {},
      K8sAPIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockProxyPromise);
  });
});

describe('getListModelVersions', () => {
  it('should call restGET and handleRestFailures to list models versions', () => {
    expect(
      getListModelVersions(`/api/${BFF_API_VERSION}/model_registry/model-registry-1/`)({}),
    ).toBe(mockResultPromise);
    expect(restGETMock).toHaveBeenCalledTimes(1);
    expect(restGETMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/model_versions`,
      {},
      K8sAPIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockProxyPromise);
  });
});

describe('getModelVersionsByRegisteredModel', () => {
  it('should call restGET and handleRestFailures to list models versions by registered model', () => {
    expect(
      getModelVersionsByRegisteredModel(`/api/${BFF_API_VERSION}/model_registry/model-registry-1/`)(
        {},
        '1',
      ),
    ).toBe(mockResultPromise);
    expect(restGETMock).toHaveBeenCalledTimes(1);
    expect(restGETMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/registered_models/1/versions`,
      {},
      K8sAPIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockProxyPromise);
  });
});

describe('getModelArtifactsByModelVersion', () => {
  it('should call restGET and handleRestFailures to list models artifacts by model version', () => {
    expect(
      getModelArtifactsByModelVersion(`/api/${BFF_API_VERSION}/model_registry/model-registry-1/`)(
        {},
        '1',
      ),
    ).toBe(mockResultPromise);
    expect(restGETMock).toHaveBeenCalledTimes(1);
    expect(restGETMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/model_versions/1/artifacts`,
      {},
      K8sAPIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockProxyPromise);
  });
});

describe('patchRegisteredModel', () => {
  it('should call restPATCH and handleRestFailures to update registered model', () => {
    expect(
      patchRegisteredModel(`/api/${BFF_API_VERSION}/model_registry/model-registry-1/`)(
        K8sAPIOptionsMock,
        { description: 'new test' },
        '1',
      ),
    ).toBe(mockResultPromise);
    expect(restPATCHMock).toHaveBeenCalledTimes(1);
    expect(restPATCHMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/registered_models/1`,
      { description: 'new test' },
      K8sAPIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockProxyPromise);
  });
});

describe('patchModelVersion', () => {
  it('should call restPATCH and handleRestFailures to update model version', () => {
    expect(
      patchModelVersion(`/api/${BFF_API_VERSION}/model_registry/model-registry-1/`)(
        K8sAPIOptionsMock,
        { description: 'new test' },
        '1',
      ),
    ).toBe(mockResultPromise);
    expect(restPATCHMock).toHaveBeenCalledTimes(1);
    expect(restPATCHMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/model_versions/1`,
      { description: 'new test' },
      K8sAPIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockProxyPromise);
  });
});

describe('patchModelArtifact', () => {
  it('should call restPATCH and handleRestFailures to update model artifact', () => {
    expect(
      patchModelArtifact(`/api/${BFF_API_VERSION}/model_registry/model-registry-1/`)(
        K8sAPIOptionsMock,
        { description: 'new test' },
        '1',
      ),
    ).toBe(mockResultPromise);
    expect(restPATCHMock).toHaveBeenCalledTimes(1);
    expect(restPATCHMock).toHaveBeenCalledWith(
      `/api/${BFF_API_VERSION}/model_registry/model-registry-1/`,
      `/model_artifacts/1`,
      { description: 'new test' },
      K8sAPIOptionsMock,
    );
    expect(handleRestFailuresMock).toHaveBeenCalledTimes(1);
    expect(handleRestFailuresMock).toHaveBeenCalledWith(mockProxyPromise);
  });
});

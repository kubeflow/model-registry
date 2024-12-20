import React from 'react';
import { APIState } from '~/shared/api/types';
import { ModelRegistryAPIs } from '~/app/types';
import {
  createModelArtifact,
  createModelArtifactForModelVersion,
  createModelVersion,
  createModelVersionForRegisteredModel,
  createRegisteredModel,
  getListModelArtifacts,
  getListModelVersions,
  getListRegisteredModels,
  getModelArtifact,
  getModelArtifactsByModelVersion,
  getModelVersion,
  getModelVersionsByRegisteredModel,
  getRegisteredModel,
  patchModelArtifact,
  patchModelVersion,
  patchRegisteredModel,
} from '~/shared/api/service';
import useAPIState from '~/shared/api/useAPIState';

export type ModelRegistryAPIState = APIState<ModelRegistryAPIs>;

const useModelRegistryAPIState = (
  hostPath: string | null,
  queryParameters?: Record<string, unknown>,
): [apiState: ModelRegistryAPIState, refreshAPIState: () => void] => {
  const createAPI = React.useCallback(
    (path: string) => ({
      createRegisteredModel: createRegisteredModel(path, queryParameters),
      createModelVersion: createModelVersion(path, queryParameters),
      createModelVersionForRegisteredModel: createModelVersionForRegisteredModel(
        path,
        queryParameters,
      ),
      createModelArtifact: createModelArtifact(path, queryParameters),
      createModelArtifactForModelVersion: createModelArtifactForModelVersion(path, queryParameters),
      getRegisteredModel: getRegisteredModel(path, queryParameters),
      getModelVersion: getModelVersion(path, queryParameters),
      getModelArtifact: getModelArtifact(path, queryParameters),
      listModelArtifacts: getListModelArtifacts(path, queryParameters),
      listModelVersions: getListModelVersions(path, queryParameters),
      listRegisteredModels: getListRegisteredModels(path, queryParameters),
      getModelVersionsByRegisteredModel: getModelVersionsByRegisteredModel(path, queryParameters),
      getModelArtifactsByModelVersion: getModelArtifactsByModelVersion(path, queryParameters),
      patchRegisteredModel: patchRegisteredModel(path, queryParameters),
      patchModelVersion: patchModelVersion(path, queryParameters),
      patchModelArtifact: patchModelArtifact(path, queryParameters),
    }),
    [queryParameters],
  );

  return useAPIState(hostPath, createAPI);
};

export default useModelRegistryAPIState;

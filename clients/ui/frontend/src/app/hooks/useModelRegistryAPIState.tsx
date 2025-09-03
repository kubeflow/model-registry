import React from 'react';
import { APIState, useAPIState } from 'mod-arch-core';
import { ModelRegistryAPIs } from '~/app/types';
import {
  createModelArtifactForModelVersion,
  createModelVersionForRegisteredModel,
  createRegisteredModel,
  getListModelArtifacts,
  getListModelVersions,
  getListRegisteredModels,
  getModelArtifactsByModelVersion,
  getModelVersion,
  getModelVersionsByRegisteredModel,
  getRegisteredModel,
  patchModelArtifact,
  patchModelVersion,
  patchRegisteredModel,
} from '~/app/api/service';

export type ModelRegistryAPIState = APIState<ModelRegistryAPIs>;

const useModelRegistryAPIState = (
  hostPath: string | null,
  queryParameters?: Record<string, unknown>,
): [apiState: ModelRegistryAPIState, refreshAPIState: () => void] => {
  const createAPI = React.useCallback(
    (path: string) => ({
      createRegisteredModel: createRegisteredModel(path, queryParameters),
      createModelVersionForRegisteredModel: createModelVersionForRegisteredModel(
        path,
        queryParameters,
      ),
      createModelArtifactForModelVersion: createModelArtifactForModelVersion(path, queryParameters),
      getRegisteredModel: getRegisteredModel(path, queryParameters),
      getModelVersion: getModelVersion(path, queryParameters),
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

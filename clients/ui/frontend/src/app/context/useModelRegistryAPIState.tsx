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
): [apiState: ModelRegistryAPIState, refreshAPIState: () => void] => {
  const createAPI = React.useCallback(
    (path: string) => ({
      createRegisteredModel: createRegisteredModel(path),
      createModelVersion: createModelVersion(path),
      createModelVersionForRegisteredModel: createModelVersionForRegisteredModel(path),
      createModelArtifact: createModelArtifact(path),
      createModelArtifactForModelVersion: createModelArtifactForModelVersion(path),
      getRegisteredModel: getRegisteredModel(path),
      getModelVersion: getModelVersion(path),
      getModelArtifact: getModelArtifact(path),
      listModelArtifacts: getListModelArtifacts(path),
      listModelVersions: getListModelVersions(path),
      listRegisteredModels: getListRegisteredModels(path),
      getModelVersionsByRegisteredModel: getModelVersionsByRegisteredModel(path),
      getModelArtifactsByModelVersion: getModelArtifactsByModelVersion(path),
      patchRegisteredModel: patchRegisteredModel(path),
      patchModelVersion: patchModelVersion(path),
      patchModelArtifact: patchModelArtifact(path),
    }),
    [],
  );

  return useAPIState(hostPath, createAPI);
};

export default useModelRegistryAPIState;

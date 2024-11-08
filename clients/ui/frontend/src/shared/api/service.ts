import {
  CreateModelArtifactData,
  CreateModelVersionData,
  CreateRegisteredModelData,
  ModelArtifact,
  ModelArtifactList,
  ModelVersionList,
  ModelVersion,
  RegisteredModelList,
  RegisteredModel,
} from '~/app/types';
import {
  assembleModelRegistryBody,
  isModelRegistryResponse,
  restCREATE,
  restGET,
  restPATCH,
} from '~/shared/api/apiUtils';
import { APIOptions } from '~/shared/api/types';
import { handleRestFailures } from '~/shared/api/errorUtils';

export const createRegisteredModel =
  (hostPath: string) =>
  (opts: APIOptions, data: CreateRegisteredModelData): Promise<RegisteredModel> =>
    handleRestFailures(
      restCREATE(hostPath, `/registered_models`, assembleModelRegistryBody(data), {}, opts),
    ).then((response) => {
      if (isModelRegistryResponse<RegisteredModel>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const createModelVersion =
  (hostPath: string) =>
  (opts: APIOptions, data: CreateModelVersionData): Promise<ModelVersion> =>
    handleRestFailures(
      restCREATE(hostPath, `/model_versions`, assembleModelRegistryBody(data), {}, opts),
    ).then((response) => {
      if (isModelRegistryResponse<ModelVersion>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const createModelVersionForRegisteredModel =
  (hostPath: string) =>
  (
    opts: APIOptions,
    registeredModelId: string,
    data: CreateModelVersionData,
  ): Promise<ModelVersion> =>
    handleRestFailures(
      restCREATE(
        hostPath,
        `/registered_models/${registeredModelId}/versions`,
        assembleModelRegistryBody(data),
        {},
        opts,
      ),
    ).then((response) => {
      if (isModelRegistryResponse<ModelVersion>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const createModelArtifact =
  (hostPath: string) =>
  (opts: APIOptions, data: CreateModelArtifactData): Promise<ModelArtifact> =>
    handleRestFailures(
      restCREATE(hostPath, `/model_artifacts`, assembleModelRegistryBody(data), {}, opts),
    ).then((response) => {
      if (isModelRegistryResponse<ModelArtifact>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const createModelArtifactForModelVersion =
  (hostPath: string) =>
  (
    opts: APIOptions,
    modelVersionId: string,
    data: CreateModelArtifactData,
  ): Promise<ModelArtifact> =>
    handleRestFailures(
      restCREATE(
        hostPath,
        `/model_versions/${modelVersionId}/artifacts`,
        assembleModelRegistryBody(data),
        {},
        opts,
      ),
    ).then((response) => {
      if (isModelRegistryResponse<ModelArtifact>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getRegisteredModel =
  (hostPath: string) =>
  (opts: APIOptions, registeredModelId: string): Promise<RegisteredModel> =>
    handleRestFailures(restGET(hostPath, `/registered_models/${registeredModelId}`, {}, opts)).then(
      (response) => {
        if (isModelRegistryResponse<RegisteredModel>(response)) {
          return response.data;
        }
        throw new Error('Invalid response format');
      },
    );

export const getModelVersion =
  (hostPath: string) =>
  (opts: APIOptions, modelversionId: string): Promise<ModelVersion> =>
    handleRestFailures(restGET(hostPath, `/model_versions/${modelversionId}`, {}, opts)).then(
      (response) => {
        if (isModelRegistryResponse<ModelVersion>(response)) {
          return response.data;
        }
        throw new Error('Invalid response format');
      },
    );

export const getModelArtifact =
  (hostPath: string) =>
  (opts: APIOptions, modelArtifactId: string): Promise<ModelArtifact> =>
    handleRestFailures(restGET(hostPath, `/model_artifacts/${modelArtifactId}`, {}, opts)).then(
      (response) => {
        if (isModelRegistryResponse<ModelArtifact>(response)) {
          return response.data;
        }
        throw new Error('Invalid response format');
      },
    );

export const getListModelArtifacts =
  (hostPath: string) =>
  (opts: APIOptions): Promise<ModelArtifactList> =>
    handleRestFailures(restGET(hostPath, `/model_artifacts`, {}, opts)).then((response) => {
      if (isModelRegistryResponse<ModelArtifactList>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getListModelVersions =
  (hostPath: string) =>
  (opts: APIOptions): Promise<ModelVersionList> =>
    handleRestFailures(restGET(hostPath, `/model_versions`, {}, opts)).then((response) => {
      if (isModelRegistryResponse<ModelVersionList>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getListRegisteredModels =
  (hostPath: string) =>
  (opts: APIOptions): Promise<RegisteredModelList> =>
    handleRestFailures(restGET(hostPath, `/registered_models`, {}, opts)).then((response) => {
      if (isModelRegistryResponse<RegisteredModelList>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getModelVersionsByRegisteredModel =
  (hostPath: string) =>
  (opts: APIOptions, registeredmodelId: string): Promise<ModelVersionList> =>
    handleRestFailures(
      restGET(hostPath, `/registered_models/${registeredmodelId}/versions`, {}, opts),
    ).then((response) => {
      if (isModelRegistryResponse<ModelVersionList>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getModelArtifactsByModelVersion =
  (hostPath: string) =>
  (opts: APIOptions, modelVersionId: string): Promise<ModelArtifactList> =>
    handleRestFailures(
      restGET(hostPath, `/model_versions/${modelVersionId}/artifacts`, {}, opts),
    ).then((response) => {
      if (isModelRegistryResponse<ModelArtifactList>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const patchRegisteredModel =
  (hostPath: string) =>
  (
    opts: APIOptions,
    data: Partial<RegisteredModel>,
    registeredModelId: string,
  ): Promise<RegisteredModel> =>
    handleRestFailures(
      restPATCH(
        hostPath,
        `/registered_models/${registeredModelId}`,
        assembleModelRegistryBody(data),
        opts,
      ),
    ).then((response) => {
      if (isModelRegistryResponse<RegisteredModel>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const patchModelVersion =
  (hostPath: string) =>
  (opts: APIOptions, data: Partial<ModelVersion>, modelversionId: string): Promise<ModelVersion> =>
    handleRestFailures(
      restPATCH(
        hostPath,
        `/model_versions/${modelversionId}`,
        assembleModelRegistryBody(data),
        opts,
      ),
    ).then((response) => {
      if (isModelRegistryResponse<ModelVersion>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const patchModelArtifact =
  (hostPath: string) =>
  (
    opts: APIOptions,
    data: Partial<ModelArtifact>,
    modelartifactId: string,
  ): Promise<ModelArtifact> =>
    handleRestFailures(
      restPATCH(
        hostPath,
        `/model_artifacts/${modelartifactId}`,
        assembleModelRegistryBody(data),
        opts,
      ),
    ).then((response) => {
      if (isModelRegistryResponse<ModelArtifact>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

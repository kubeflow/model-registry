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
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, data: CreateRegisteredModelData): Promise<RegisteredModel> =>
    handleRestFailures(
      restCREATE(
        hostPath,
        `/registered_models`,
        assembleModelRegistryBody(data),
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModelRegistryResponse<RegisteredModel>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const createModelVersion =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, data: CreateModelVersionData): Promise<ModelVersion> =>
    handleRestFailures(
      restCREATE(hostPath, `/model_versions`, assembleModelRegistryBody(data), queryParams, opts),
    ).then((response) => {
      if (isModelRegistryResponse<ModelVersion>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const createModelVersionForRegisteredModel =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
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
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModelRegistryResponse<ModelVersion>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const createModelArtifact =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, data: CreateModelArtifactData): Promise<ModelArtifact> =>
    handleRestFailures(
      restCREATE(hostPath, `/model_artifacts`, assembleModelRegistryBody(data), queryParams, opts),
    ).then((response) => {
      if (isModelRegistryResponse<ModelArtifact>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const createModelArtifactForModelVersion =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
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
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModelRegistryResponse<ModelArtifact>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getRegisteredModel =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, registeredModelId: string): Promise<RegisteredModel> =>
    handleRestFailures(
      restGET(hostPath, `/registered_models/${registeredModelId}`, queryParams, opts),
    ).then((response) => {
      if (isModelRegistryResponse<RegisteredModel>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getModelVersion =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, modelversionId: string): Promise<ModelVersion> =>
    handleRestFailures(
      restGET(hostPath, `/model_versions/${modelversionId}`, queryParams, opts),
    ).then((response) => {
      if (isModelRegistryResponse<ModelVersion>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getModelArtifact =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, modelArtifactId: string): Promise<ModelArtifact> =>
    handleRestFailures(
      restGET(hostPath, `/model_artifacts/${modelArtifactId}`, queryParams, opts),
    ).then((response) => {
      if (isModelRegistryResponse<ModelArtifact>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getListModelArtifacts =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions): Promise<ModelArtifactList> =>
    handleRestFailures(restGET(hostPath, `/model_artifacts`, queryParams, opts)).then(
      (response) => {
        if (isModelRegistryResponse<ModelArtifactList>(response)) {
          return response.data;
        }
        throw new Error('Invalid response format');
      },
    );

export const getListModelVersions =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions): Promise<ModelVersionList> =>
    handleRestFailures(restGET(hostPath, `/model_versions`, queryParams, opts)).then((response) => {
      if (isModelRegistryResponse<ModelVersionList>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getListRegisteredModels =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions): Promise<RegisteredModelList> =>
    handleRestFailures(restGET(hostPath, `/registered_models`, queryParams, opts)).then(
      (response) => {
        if (isModelRegistryResponse<RegisteredModelList>(response)) {
          return response.data;
        }
        throw new Error('Invalid response format');
      },
    );

export const getModelVersionsByRegisteredModel =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, registeredmodelId: string): Promise<ModelVersionList> =>
    handleRestFailures(
      restGET(hostPath, `/registered_models/${registeredmodelId}/versions`, queryParams, opts),
    ).then((response) => {
      if (isModelRegistryResponse<ModelVersionList>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getModelArtifactsByModelVersion =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, modelVersionId: string): Promise<ModelArtifactList> =>
    handleRestFailures(
      restGET(hostPath, `/model_versions/${modelVersionId}/artifacts`, queryParams, opts),
    ).then((response) => {
      if (isModelRegistryResponse<ModelArtifactList>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const patchRegisteredModel =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
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
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModelRegistryResponse<RegisteredModel>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const patchModelVersion =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, data: Partial<ModelVersion>, modelversionId: string): Promise<ModelVersion> =>
    handleRestFailures(
      restPATCH(
        hostPath,
        `/model_versions/${modelversionId}`,
        assembleModelRegistryBody(data),
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModelRegistryResponse<ModelVersion>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const patchModelArtifact =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
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
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModelRegistryResponse<ModelArtifact>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

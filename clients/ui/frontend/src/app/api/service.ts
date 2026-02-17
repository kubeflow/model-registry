import {
  APIOptions,
  assembleModArchBody,
  isModArchResponse,
  restCREATE,
  restGET,
  restPATCH,
  handleRestFailures,
  restDELETE,
} from 'mod-arch-core';
import {
  CreateModelArtifactData,
  CreateModelVersionData,
  CreateRegisteredModelData,
  CreateModelTransferJobData,
  ModelArtifact,
  ModelArtifactList,
  ModelVersionList,
  ModelVersion,
  RegisteredModelList,
  RegisteredModel,
  ModelTransferJobList,
  ModelTransferJob,
} from '~/app/types';
import { bumpRegisteredModelTimestamp } from '~/app/api/updateTimestamps';

export const createRegisteredModel =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, data: CreateRegisteredModelData): Promise<RegisteredModel> =>
    handleRestFailures(
      restCREATE(hostPath, `/registered_models`, assembleModArchBody(data), queryParams, opts),
    ).then((response) => {
      if (isModArchResponse<RegisteredModel>(response)) {
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
    registeredModel: RegisteredModel,
    isFirstVersion?: boolean,
  ): Promise<ModelVersion> =>
    handleRestFailures(
      restCREATE(
        hostPath,
        `/registered_models/${registeredModelId}/versions`,
        assembleModArchBody(data),
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModArchResponse<ModelVersion>(response)) {
        const newVersion = response.data;

        if (!isFirstVersion) {
          return bumpRegisteredModelTimestamp(
            { patchRegisteredModel: patchRegisteredModel(hostPath, queryParams) },
            registeredModel,
          ).then(() => newVersion);
        }
        return newVersion;
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
        assembleModArchBody(data),
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModArchResponse<ModelArtifact>(response)) {
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
      if (isModArchResponse<RegisteredModel>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getModelVersion =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, modelVersionId: string): Promise<ModelVersion> =>
    handleRestFailures(
      restGET(hostPath, `/model_versions/${modelVersionId}`, queryParams, opts),
    ).then((response) => {
      if (isModArchResponse<ModelVersion>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getListModelArtifacts =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions): Promise<ModelArtifactList> =>
    handleRestFailures(restGET(hostPath, `/model_artifacts`, queryParams, opts)).then(
      (response) => {
        if (isModArchResponse<ModelArtifactList>(response)) {
          return response.data;
        }
        throw new Error('Invalid response format');
      },
    );

export const getListModelVersions =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions): Promise<ModelVersionList> =>
    handleRestFailures(restGET(hostPath, `/model_versions`, queryParams, opts)).then((response) => {
      if (isModArchResponse<ModelVersionList>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getListRegisteredModels =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions): Promise<RegisteredModelList> =>
    handleRestFailures(restGET(hostPath, `/registered_models`, queryParams, opts)).then(
      (response) => {
        if (isModArchResponse<RegisteredModelList>(response)) {
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
      if (isModArchResponse<ModelVersionList>(response)) {
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
      if (isModArchResponse<ModelArtifactList>(response)) {
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
        assembleModArchBody(data),
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModArchResponse<RegisteredModel>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const patchModelVersion =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, data: Partial<ModelVersion>, modelVersionId: string): Promise<ModelVersion> =>
    handleRestFailures(
      restPATCH(
        hostPath,
        `/model_versions/${modelVersionId}`,
        assembleModArchBody(data),
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModArchResponse<ModelVersion>(response)) {
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
        assembleModArchBody(data),
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModArchResponse<ModelArtifact>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getListModelTransferJobs =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions): Promise<ModelTransferJobList> =>
    handleRestFailures(restGET(hostPath, `/model_transfer_jobs`, queryParams, opts)).then(
      (response) => {
        if (isModArchResponse<ModelTransferJobList>(response)) {
          return response.data;
        }
        throw new Error('Invalid response format');
      },
    );

export const createModelTransferJob =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, data: CreateModelTransferJobData): Promise<ModelTransferJob> =>
    handleRestFailures(
      restCREATE(hostPath, '/model_transfer_jobs', assembleModArchBody(data), queryParams, opts),
    ).then((response) => {
      if (isModArchResponse<ModelTransferJob>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const updateModelTransferJob =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, jobId: string, data: Partial<ModelTransferJob>): Promise<ModelTransferJob> =>
    handleRestFailures(
      restPATCH(
        hostPath,
        `/model_transfer_jobs/${jobId}`,
        assembleModArchBody(data),
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModArchResponse<ModelTransferJob>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const deleteModelTransferJob =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, jobName: string): Promise<void> =>
    handleRestFailures(
      restDELETE(
        hostPath,
        `/model_transfer_jobs/${encodeURIComponent(jobName)}`,
        {},
        queryParams,
        opts,
      ),
    );

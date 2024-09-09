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
import { restCREATE, restGET, restPATCH } from '~/app/api/apiUtils';
import { APIOptions } from '~/types';
import { handleRestFailures } from '~/app/api/errorUtils';
import { BFF_API_VERSION } from '~/app/const';

export const createRegisteredModel =
  (hostPath: string, mrName: string) =>
  (opts: APIOptions, data: CreateRegisteredModelData): Promise<RegisteredModel> =>
    handleRestFailures(
      restCREATE(
        hostPath,
        `/api/${BFF_API_VERSION}/model_registry/${mrName}/registered_models`,
        data,
        {},
        opts,
      ),
    );

export const createModelVersion =
  (hostPath: string, mrName: string) =>
  (opts: APIOptions, data: CreateModelVersionData): Promise<ModelVersion> =>
    handleRestFailures(
      restCREATE(
        hostPath,
        `/api/${BFF_API_VERSION}/model_registry/${mrName}/model_versions`,
        data,
        {},
        opts,
      ),
    );
export const createModelVersionForRegisteredModel =
  (hostPath: string, mrName: string) =>
  (
    opts: APIOptions,
    registeredModelId: string,
    data: CreateModelVersionData,
  ): Promise<ModelVersion> =>
    handleRestFailures(
      restCREATE(
        hostPath,
        `/api/${BFF_API_VERSION}/model_registry/${mrName}/registered_models/${registeredModelId}/versions`,
        data,
        {},
        opts,
      ),
    );

export const createModelArtifact =
  (hostPath: string, mrName: string) =>
  (opts: APIOptions, data: CreateModelArtifactData): Promise<ModelArtifact> =>
    handleRestFailures(
      restCREATE(
        hostPath,
        `/api/${BFF_API_VERSION}/model_registry/${mrName}/model_artifacts`,
        data,
        {},
        opts,
      ),
    );

export const createModelArtifactForModelVersion =
  (hostPath: string, mrName: string) =>
  (
    opts: APIOptions,
    modelVersionId: string,
    data: CreateModelArtifactData,
  ): Promise<ModelArtifact> =>
    handleRestFailures(
      restCREATE(
        hostPath,
        `/api/${BFF_API_VERSION}/model_registry/${mrName}/model_versions/${modelVersionId}/artifacts`,
        data,
        {},
        opts,
      ),
    );

export const getRegisteredModel =
  (hostPath: string, mrName: string) =>
  (opts: APIOptions, registeredModelId: string): Promise<RegisteredModel> =>
    handleRestFailures(
      restGET(
        hostPath,
        `/api/${BFF_API_VERSION}/model_registry/${mrName}/registered_models/${registeredModelId}`,
        {},
        opts,
      ),
    );

export const getModelVersion =
  (hostPath: string, mrName: string) =>
  (opts: APIOptions, modelversionId: string): Promise<ModelVersion> =>
    handleRestFailures(
      restGET(
        hostPath,
        `/api/${BFF_API_VERSION}/model_registry/${mrName}/model_versions/${modelversionId}`,
        {},
        opts,
      ),
    );

export const getModelArtifact =
  (hostPath: string, mrName: string) =>
  (opts: APIOptions, modelArtifactId: string): Promise<ModelArtifact> =>
    handleRestFailures(
      restGET(
        hostPath,
        `/api/${BFF_API_VERSION}/model_registry/${mrName}/model_artifacts/${modelArtifactId}`,
        {},
        opts,
      ),
    );

export const getListModelArtifacts =
  (hostPath: string, mrName: string) =>
  (opts: APIOptions): Promise<ModelArtifactList> =>
    handleRestFailures(
      restGET(
        hostPath,
        `/api/${BFF_API_VERSION}/model_registry/${mrName}/model_artifacts`,
        {},
        opts,
      ),
    );

export const getListModelVersions =
  (hostPath: string, mrName: string) =>
  (opts: APIOptions): Promise<ModelVersionList> =>
    handleRestFailures(
      restGET(
        hostPath,
        `/api/${BFF_API_VERSION}/model_registry/${mrName}/model_versions`,
        {},
        opts,
      ),
    );

export const getListRegisteredModels =
  (hostPath: string, mrName: string) =>
  (opts: APIOptions): Promise<RegisteredModelList> =>
    handleRestFailures(
      restGET(
        hostPath,
        `/api/${BFF_API_VERSION}/model_registry/${mrName}/registered_models`,
        {},
        opts,
      ),
    );

export const getModelVersionsByRegisteredModel =
  (hostPath: string, mrName: string) =>
  (opts: APIOptions, registeredmodelId: string): Promise<ModelVersionList> =>
    handleRestFailures(
      restGET(
        hostPath,
        `/api/${BFF_API_VERSION}/model_registry/${mrName}/registered_models/${registeredmodelId}/versions`,
        {},
        opts,
      ),
    );

export const getModelArtifactsByModelVersion =
  (hostPath: string, mrName: string) =>
  (opts: APIOptions, modelVersionId: string): Promise<ModelArtifactList> =>
    handleRestFailures(
      restGET(
        hostPath,
        `/api/${BFF_API_VERSION}/model_registry/${mrName}/model_versions/${modelVersionId}/artifacts`,
        {},
        opts,
      ),
    );

export const patchRegisteredModel =
  (hostPath: string, mrName: string) =>
  (
    opts: APIOptions,
    data: Partial<RegisteredModel>,
    registeredModelId: string,
  ): Promise<RegisteredModel> =>
    handleRestFailures(
      restPATCH(
        hostPath,
        `/api/${BFF_API_VERSION}/model_registry/${mrName}/registered_models/${registeredModelId}`,
        data,
        opts,
      ),
    );

export const patchModelVersion =
  (hostPath: string, mrName: string) =>
  (opts: APIOptions, data: Partial<ModelVersion>, modelversionId: string): Promise<ModelVersion> =>
    handleRestFailures(
      restPATCH(
        hostPath,
        `/api/${BFF_API_VERSION}/model_registry/${mrName}/model_versions/${modelversionId}`,
        data,
        opts,
      ),
    );

export const patchModelArtifact =
  (hostPath: string, mrName: string) =>
  (
    opts: APIOptions,
    data: Partial<ModelArtifact>,
    modelartifactId: string,
  ): Promise<ModelArtifact> =>
    handleRestFailures(
      restPATCH(
        hostPath,
        `/api/${BFF_API_VERSION}/model_registry/${mrName}/model_artifacts/${modelartifactId}`,
        data,
        opts,
      ),
    );

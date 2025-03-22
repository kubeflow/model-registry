import {
  ModelArtifact,
  ModelArtifactState,
  ModelState,
  ModelVersion,
  RegisteredModel,
  RegisteredModelList,
} from '~/app/types';
import { ModelRegistryAPIState } from '~/app/hooks/useModelRegistryAPIState';
import { objectStorageFieldsToUri } from '~/app/utils';
import {
  ModelLocationType,
  RegisterCatalogModelFormData,
  RegisterModelFormData,
  RegisterVersionFormData,
  RegistrationCommonFormData,
} from './useRegisterModelData';
import { RegistrationErrorType, MR_CHARACTER_LIMIT } from './const';

export type RegisterModelCreatedResources = RegisterVersionCreatedResources & {
  registeredModel?: RegisteredModel;
};

export type RegisterVersionCreatedResources = {
  modelVersion?: ModelVersion;
  modelArtifact?: ModelArtifact;
};

export const registerModel = async (
  apiState: ModelRegistryAPIState,
  formData: RegisterModelFormData,
  author: string,
): Promise<{
  data: RegisterModelCreatedResources;
  errors: { [key: string]: Error | undefined };
}> => {
  let registeredModel;
  const error: { [key: string]: Error | undefined } = {};
  try {
    registeredModel = await apiState.api.createRegisteredModel(
      {},
      {
        name: formData.modelName,
        description: formData.modelDescription,
        customProperties: formData.modelCustomProperties || {},
        owner: author,
        state: ModelState.LIVE,
      },
    );
  } catch (e) {
    if (e instanceof Error) {
      error[RegistrationErrorType.REGISTERED_MODEL] = e;
    }
    return { data: { registeredModel }, errors: error };
  }
  const {
    data: { modelVersion, modelArtifact },
    errors,
  } = await registerVersion(apiState, registeredModel, formData, author, true);

  return {
    data: { registeredModel, modelVersion, modelArtifact },
    errors,
  };
};

export const registerVersion = async (
  apiState: ModelRegistryAPIState,
  registeredModel: RegisteredModel,
  formData: Omit<RegisterVersionFormData, 'registeredModelId'>,
  author: string,
  isFirstVersion?: boolean,
): Promise<{
  data: RegisterVersionCreatedResources;
  errors: { [key: string]: Error | undefined };
}> => {
  let modelVersion;
  let modelArtifact;
  const errors: { [key: string]: Error | undefined } = {};
  try {
    modelVersion = await apiState.api.createModelVersionForRegisteredModel(
      {},
      registeredModel.id,
      {
        name: formData.versionName,
        description: formData.versionDescription,
        customProperties: formData.versionCustomProperties || {},
        state: ModelState.LIVE,
        author,
        registeredModelId: registeredModel.id,
      },
      registeredModel,
      isFirstVersion,
    );
  } catch (e) {
    if (e instanceof Error) {
      errors[RegistrationErrorType.MODEL_VERSION] = e;
    }
    return { data: { modelVersion, modelArtifact }, errors };
  }

  try {
    modelArtifact = await apiState.api.createModelArtifactForModelVersion({}, modelVersion.id, {
      name: `${formData.versionName}`,
      description: formData.versionDescription,
      customProperties: {},
      state: ModelArtifactState.LIVE,
      author,
      modelFormatName: formData.sourceModelFormat,
      modelFormatVersion: formData.sourceModelFormatVersion,
      // storageKey: 'TODO',
      uri:
        formData.modelLocationType === ModelLocationType.ObjectStorage
          ? objectStorageFieldsToUri({
              endpoint: formData.modelLocationEndpoint,
              bucket: formData.modelLocationBucket,
              region: formData.modelLocationRegion,
              path: formData.modelLocationPath,
            }) || '' // We'll only hit this case if required fields are empty strings, so form validation should catch it.
          : formData.modelLocationURI,
      artifactType: 'model-artifact',
    });
  } catch (e) {
    if (e instanceof Error) {
      errors[RegistrationErrorType.MODEL_ARTIFACT] = e;
    }
  }

  return { data: { modelVersion, modelArtifact }, errors };
};

const isSubmitDisabledForCommonFields = (formData: RegistrationCommonFormData): boolean => {
  const {
    versionName,
    modelLocationType,
    modelLocationURI,
    modelLocationBucket,
    modelLocationEndpoint,
    modelLocationPath,
  } = formData;
  return (
    !versionName ||
    (modelLocationType === ModelLocationType.URI && !modelLocationURI) ||
    (modelLocationType === ModelLocationType.ObjectStorage &&
      (!modelLocationBucket || !modelLocationEndpoint || !modelLocationPath)) ||
    !isNameValid(versionName)
  );
};

export const isRegisterModelSubmitDisabled = (
  formData: RegisterModelFormData,
  registeredModels: RegisteredModelList,
): boolean =>
  !formData.modelName ||
  isSubmitDisabledForCommonFields(formData) ||
  !isNameValid(formData.modelName) ||
  isModelNameExisting(formData.modelName, registeredModels);

export const isRegisterVersionSubmitDisabled = (formData: RegisterVersionFormData): boolean =>
  !formData.registeredModelId || isSubmitDisabledForCommonFields(formData);

export const isRegisterCatalogModelSubmitDisabled = (
  formData: RegisterCatalogModelFormData,
  registeredModels: RegisteredModelList,
): boolean => isRegisterModelSubmitDisabled(formData, registeredModels) || !formData.modelRegistry;

export const isNameValid = (name: string): boolean => name.length <= MR_CHARACTER_LIMIT;

export const isModelNameExisting = (name: string, registeredModels: RegisteredModelList): boolean =>
  registeredModels.items.some((model) => model.name === name);

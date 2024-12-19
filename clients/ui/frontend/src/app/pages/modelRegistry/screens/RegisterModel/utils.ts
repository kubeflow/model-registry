import {
  ModelArtifact,
  ModelArtifactState,
  ModelState,
  ModelVersion,
  RegisteredModel,
} from '~/app/types';
import { ModelRegistryAPIState } from '~/app/hooks/useModelRegistryAPIState';
import { objectStorageFieldsToUri } from '~/app/pages/modelRegistry/screens/utils';
import {
  ModelLocationType,
  RegisterModelFormData,
  RegisterVersionFormData,
  RegistrationCommonFormData,
} from './useRegisterModelData';

export type RegisterModelCreatedResources = RegisterVersionCreatedResources & {
  registeredModel: RegisteredModel;
};

export type RegisterVersionCreatedResources = {
  modelVersion: ModelVersion;
  modelArtifact: ModelArtifact;
};

export const registerModel = async (
  apiState: ModelRegistryAPIState,
  formData: RegisterModelFormData,
  author: string,
): Promise<RegisterModelCreatedResources> => {
  const registeredModel = await apiState.api.createRegisteredModel(
    {},
    {
      name: formData.modelName,
      description: formData.modelDescription,
      customProperties: {},
      owner: author,
      state: ModelState.LIVE,
    },
  );
  const { modelVersion, modelArtifact } = await registerVersion(
    apiState,
    registeredModel,
    formData,
    author,
  );
  return { registeredModel, modelVersion, modelArtifact };
};

export const registerVersion = async (
  apiState: ModelRegistryAPIState,
  registeredModel: RegisteredModel,
  formData: Omit<RegisterVersionFormData, 'registeredModelId'>,
  author: string,
): Promise<RegisterVersionCreatedResources> => {
  const modelVersion = await apiState.api.createModelVersionForRegisteredModel(
    {},
    registeredModel.id,
    {
      name: formData.versionName,
      description: formData.versionDescription,
      customProperties: {},
      state: ModelState.LIVE,
      author,
      registeredModelId: registeredModel.id,
    },
  );
  const modelArtifact = await apiState.api.createModelArtifactForModelVersion({}, modelVersion.id, {
    name: `${registeredModel.name}-${formData.versionName}-artifact`,
    description: formData.versionDescription,
    customProperties: {},
    state: ModelArtifactState.LIVE,
    author,
    modelFormatName: formData.sourceModelFormat,
    modelFormatVersion: formData.sourceModelFormatVersion,
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
  return { modelVersion, modelArtifact };
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
      (!modelLocationBucket || !modelLocationEndpoint || !modelLocationPath))
  );
};

export const isRegisterModelSubmitDisabled = (formData: RegisterModelFormData): boolean =>
  !formData.modelName || isSubmitDisabledForCommonFields(formData);

export const isRegisterVersionSubmitDisabled = (formData: RegisterVersionFormData): boolean =>
  !formData.registeredModelId || isSubmitDisabledForCommonFields(formData);

import {
  ModelArtifact,
  ModelArtifactState,
  ModelState,
  ModelVersion,
  RegisteredModel,
  RegisteredModelList,
  ModelTransferJob,
  CreateModelTransferJobData,
  ModelTransferJobSource,
  ModelTransferJobSourceType,
  ModelTransferJobDestinationType,
  ModelTransferJobOCIDestination,
  ModelTransferJobUploadIntent,
  ModelTransferJobStatus,
} from '~/app/types';
import { ModelRegistryAPIState } from '~/app/hooks/useModelRegistryAPIState';
import { objectStorageFieldsToUri } from '~/app/utils';
import { RegistrationMode } from '~/app/pages/modelRegistry/screens/const';
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
      ...formData.additionalArtifactProperties,
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
    modelLocationS3AccessKeyId,
    modelLocationS3SecretAccessKey,
    registrationMode,
    namespace,
    destinationOciRegistry,
    destinationOciUri,
    destinationOciUsername,
    destinationOciPassword,
    jobResourceName,
  } = formData;

  // RegisterAndStore mode validation - require destination fields, credentials, and job name
  if (registrationMode === RegistrationMode.RegisterAndStore) {
    // Base requirements for register-and-store mode
    if (!namespace || !destinationOciRegistry || !destinationOciUri || !jobResourceName) {
      return true;
    }
    // Destination credentials are required
    if (!destinationOciUsername || !destinationOciPassword) {
      return true;
    }
    // Source credentials are required for S3/ObjectStorage
    if (
      modelLocationType === ModelLocationType.ObjectStorage &&
      (!modelLocationS3AccessKeyId || !modelLocationS3SecretAccessKey)
    ) {
      return true;
    }
  }

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

// Helper function to build ModelTransferJob payload from form data
// TODO: When ModelTransferJob API is extended, add support for:
//   - Credentials: formData.modelLocationS3AccessKeyId, formData.modelLocationS3SecretAccessKey
//                  formData.destinationOciUsername, formData.destinationOciPassword, formData.destinationOciEmail
//   - Model metadata: formData.modelDescription (for CREATE_MODEL), formData.versionDescription
//   - Model format: formData.sourceModelFormat, formData.sourceModelFormatVersion
//   - Custom properties: formData.modelCustomProperties, formData.versionCustomProperties
export const buildModelTransferJobPayload = (
  formData: RegisterModelFormData | RegisterVersionFormData,
  author: string,
  uploadIntent: ModelTransferJobUploadIntent,
  registeredModelId?: string,
  registeredModelName?: string,
): CreateModelTransferJobData => {
  // Build source based on modelLocationType
  const source: ModelTransferJobSource =
    formData.modelLocationType === ModelLocationType.ObjectStorage
      ? {
          type: ModelTransferJobSourceType.S3,
          bucket: formData.modelLocationBucket,
          key: formData.modelLocationPath,
          region: formData.modelLocationRegion || undefined,
          endpoint: formData.modelLocationEndpoint || undefined,
        }
      : {
          type: ModelTransferJobSourceType.URI,
          uri: formData.modelLocationURI,
        };

  // Build OCI destination
  const destination: ModelTransferJobOCIDestination = {
    type: ModelTransferJobDestinationType.OCI,
    uri: formData.destinationOciUri,
    registry: formData.destinationOciRegistry || undefined,
  };

  // RegisterModelFormData has modelName (user-provided for new model).
  // RegisterVersionFormData omits it since the model already exists; we use registeredModelName instead.
  const modelName = 'modelName' in formData ? formData.modelName : registeredModelName;

  return {
    name: formData.jobResourceName,
    source,
    destination,
    uploadIntent,
    registeredModelId,
    registeredModelName: modelName,
    modelVersionName: formData.versionName,
    namespace: formData.namespace,
    author,
    status: ModelTransferJobStatus.PENDING,
  };
};

// Result type for transfer job creation
export type RegisterViaTransferJobResult = {
  transferJob?: ModelTransferJob;
  error?: Error;
};

// Options for registerViaTransferJob based on intent
type RegisterViaTransferJobOptions =
  | {
      intent: typeof ModelTransferJobUploadIntent.CREATE_MODEL;
      formData: RegisterModelFormData;
    }
  | {
      intent: typeof ModelTransferJobUploadIntent.CREATE_VERSION;
      formData: RegisterVersionFormData;
      registeredModel: RegisteredModel;
    };

// Create transfer job for async model registration (handles both new models and new versions)
export const registerViaTransferJob = async (
  apiState: ModelRegistryAPIState,
  author: string,
  options: RegisterViaTransferJobOptions,
): Promise<RegisterViaTransferJobResult> => {
  try {
    const payload =
      options.intent === ModelTransferJobUploadIntent.CREATE_MODEL
        ? buildModelTransferJobPayload(options.formData, author, options.intent)
        : buildModelTransferJobPayload(
            options.formData,
            author,
            options.intent,
            options.registeredModel.id,
            options.registeredModel.name,
          );
    const transferJob = await apiState.api.createModelTransferJob({}, payload);
    return { transferJob };
  } catch (e) {
    return { error: e instanceof Error ? e : new Error('Failed to create transfer job') };
  }
};

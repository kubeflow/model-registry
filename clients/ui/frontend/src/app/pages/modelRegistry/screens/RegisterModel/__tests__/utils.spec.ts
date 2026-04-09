import {
  RegisteredModelList,
  ModelRegistryMetadataType,
  ModelTransferJobSourceType,
  ModelTransferJobDestinationType,
  ModelTransferJobUploadIntent,
  ModelTransferJobStatus,
} from '~/app/types';
import {
  isModelNameExisting,
  isNameValid,
  buildModelTransferJobPayload,
  isRegisterModelSubmitDisabled,
  isRegisterCatalogModelSubmitDisabled,
} from '~/app/pages/modelRegistry/screens/RegisterModel/utils';
import { MR_CHARACTER_LIMIT } from '~/app/pages/modelRegistry/screens/RegisterModel/const';
import {
  ModelLocationType,
  RegisterModelFormData,
} from '~/app/pages/modelRegistry/screens/RegisterModel/useRegisterModelData';
import { RegistrationMode } from '~/app/pages/modelRegistry/screens/const';
import { CatalogModelCustomPropertyKey, ModelType } from '~/concepts/modelCatalog/const';

/** Shared fields for registration utils tests (transfer job + submit-disabled). */
const registrationFormTestBase = {
  versionName: 'v1.0.0',
  versionDescription: 'Test version',
  sourceModelFormat: 'onnx',
  sourceModelFormatVersion: '1.0',
  modelLocationType: ModelLocationType.ObjectStorage,
  modelLocationEndpoint: 'https://s3.amazonaws.com',
  modelLocationBucket: 'test-bucket',
  modelLocationRegion: 'us-east-1',
  modelLocationPath: 'models/test',
  modelLocationURI: '',
  modelLocationS3AccessKeyId: '',
  modelLocationS3SecretAccessKey: '',
  registrationMode: RegistrationMode.RegisterAndStore,
  namespace: 'test-namespace',
  destinationOciRegistry: 'quay.io',
  destinationOciUsername: '',
  destinationOciPassword: '',
  destinationOciUri: 'quay.io/org/model:v1',
  jobName: 'test-job',
  jobResourceName: 'test-job-resource',
  versionCustomProperties: {},
  additionalArtifactProperties: {},
};

describe('RegisterModel utils', () => {
  const emptyRegisteredModelList = {
    items: [],
    size: 0,
    pageSize: 20,
    nextPageToken: '',
  } as RegisteredModelList;

  /** Register + URI path: only model type should gate MR submit when required. */
  const mrRegisterForm = (
    modelCustomProperties: RegisterModelFormData['modelCustomProperties'],
  ): RegisterModelFormData => ({
    ...registrationFormTestBase,
    modelName: 'unique-new-model',
    modelDescription: '',
    registrationMode: RegistrationMode.Register,
    modelLocationType: ModelLocationType.URI,
    modelLocationURI: 'https://example.com/model.onnx',
    modelLocationEndpoint: '',
    modelLocationBucket: '',
    modelLocationRegion: '',
    modelLocationPath: '',
    namespace: '',
    destinationOciRegistry: '',
    destinationOciUsername: '',
    destinationOciPassword: '',
    destinationOciUri: '',
    jobName: '',
    jobResourceName: '',
    modelCustomProperties,
    versionCustomProperties: {},
  });

  describe('isRegisterModelSubmitDisabled (model type)', () => {
    it('disables submit until model type is selected when requireModelType is true', () => {
      expect(
        isRegisterModelSubmitDisabled(mrRegisterForm({}), emptyRegisteredModelList, undefined, undefined, {
          requireModelType: true,
        }),
      ).toBe(true);
    });

    it('allows submit once model type is set when requireModelType is true', () => {
      expect(
        isRegisterModelSubmitDisabled(
          mrRegisterForm({
            [CatalogModelCustomPropertyKey.MODEL_TYPE]: {
              metadataType: ModelRegistryMetadataType.STRING,
              // eslint-disable-next-line camelcase
              string_value: ModelType.GENERATIVE,
            },
          }),
          emptyRegisteredModelList,
          undefined,
          undefined,
          { requireModelType: true },
        ),
      ).toBe(false);
    });

    it('does not require model type by default', () => {
      expect(isRegisterModelSubmitDisabled(mrRegisterForm({}), emptyRegisteredModelList)).toBe(
        false,
      );
    });
  });

  describe('isRegisterCatalogModelSubmitDisabled', () => {
    it('allows submit without model type when registry is selected', () => {
      expect(
        isRegisterCatalogModelSubmitDisabled(
          { ...mrRegisterForm({}), modelRegistry: 'test-mr' },
          emptyRegisteredModelList,
        ),
      ).toBe(false);
    });
  });

  describe('isModelNameExisting', () => {
    const existingModelName = 'model2';
    const newModelName = 'model4';
    const modelList = {
      items: [{ name: 'model1' }, { name: existingModelName }, { name: 'model3' }],
    } as RegisteredModelList;
    it('should return true if model name exists in list', () => {
      expect(isModelNameExisting(existingModelName, modelList)).toBe(true);
    });

    it('should return false if model name does not exist in list', () => {
      expect(isModelNameExisting(newModelName, modelList)).toBe(false);
    });
  });

  describe('isNameValid', () => {
    it('should return true for valid model names (currently only limited by character count)', () => {
      expect(isNameValid('x'.repeat(MR_CHARACTER_LIMIT))).toBe(true);
      expect(isNameValid('')).toBe(true); //will be caught by form 'required' validation
    });
    it('should return false for names that are too long', () => {
      expect(isNameValid('x'.repeat(MR_CHARACTER_LIMIT + 1))).toBe(false);
    });
  });

  describe('buildModelTransferJobPayload', () => {
    it('should build payload with S3 source for ObjectStorage location type', () => {
      const formData = {
        ...registrationFormTestBase,
        modelName: 'Test Model',
        modelDescription: '',
      };
      const payload = buildModelTransferJobPayload(
        formData,
        'test-author',
        ModelTransferJobUploadIntent.CREATE_MODEL,
      );

      expect(payload.source.type).toBe(ModelTransferJobSourceType.S3);
      expect(payload.source).toMatchObject({
        bucket: 'test-bucket',
        key: 'models/test',
        region: 'us-east-1',
      });
    });

    it('should build payload with URI source for URI location type', () => {
      const formData = {
        ...registrationFormTestBase,
        modelName: 'Test Model',
        modelDescription: '',
        modelLocationType: ModelLocationType.URI,
        modelLocationURI: 'https://example.com/model.onnx',
      };
      const payload = buildModelTransferJobPayload(
        formData,
        'test-author',
        ModelTransferJobUploadIntent.CREATE_MODEL,
      );

      expect(payload.source.type).toBe(ModelTransferJobSourceType.URI);
      expect(payload.source).toMatchObject({ uri: 'https://example.com/model.onnx' });
    });

    it('should build OCI destination correctly', () => {
      const formData = {
        ...registrationFormTestBase,
        modelName: 'Test Model',
        modelDescription: '',
      };
      const payload = buildModelTransferJobPayload(
        formData,
        'test-author',
        ModelTransferJobUploadIntent.CREATE_MODEL,
      );

      expect(payload.destination.type).toBe(ModelTransferJobDestinationType.OCI);
      expect(payload.destination).toMatchObject({
        uri: 'quay.io/org/model:v1',
        registry: 'quay.io',
      });
    });

    it('should set CREATE_MODEL intent and include model name', () => {
      const formData = {
        ...registrationFormTestBase,
        modelName: 'My New Model',
        modelDescription: '',
      };
      const payload = buildModelTransferJobPayload(
        formData,
        'test-author',
        ModelTransferJobUploadIntent.CREATE_MODEL,
      );

      expect(payload.uploadIntent).toBe(ModelTransferJobUploadIntent.CREATE_MODEL);
      expect(payload.registeredModelName).toBe('My New Model');
    });

    it('should set CREATE_VERSION intent with registeredModelId', () => {
      const formData = { ...registrationFormTestBase, registeredModelId: 'existing-model-123' };
      const payload = buildModelTransferJobPayload(
        formData,
        'test-author',
        ModelTransferJobUploadIntent.CREATE_VERSION,
        'existing-model-123',
        'Existing Model Name',
      );

      expect(payload.uploadIntent).toBe(ModelTransferJobUploadIntent.CREATE_VERSION);
      expect(payload.registeredModelId).toBe('existing-model-123');
      expect(payload.registeredModelName).toBe('Existing Model Name');
    });

    it('should include namespace, author, and job resource name', () => {
      const formData = {
        ...registrationFormTestBase,
        modelName: 'Test Model',
        modelDescription: '',
      };
      const payload = buildModelTransferJobPayload(
        formData,
        'test-author',
        ModelTransferJobUploadIntent.CREATE_MODEL,
      );

      expect(payload.namespace).toBe('test-namespace');
      expect(payload.author).toBe('test-author');
      expect(payload.name).toBe('test-job-resource');
    });

    it('should set PENDING status and omit server-generated fields', () => {
      const formData = {
        ...registrationFormTestBase,
        modelName: 'Test Model',
        modelDescription: '',
      };
      const payload = buildModelTransferJobPayload(
        formData,
        'test-author',
        ModelTransferJobUploadIntent.CREATE_MODEL,
      );

      expect(payload.status).toBe(ModelTransferJobStatus.PENDING);
      // id, createTimeSinceEpoch, lastUpdateTimeSinceEpoch are omitted by CreateModelTransferJobData type
      expect('id' in payload).toBe(false);
      expect('createTimeSinceEpoch' in payload).toBe(false);
      expect('lastUpdateTimeSinceEpoch' in payload).toBe(false);
    });
  });
});

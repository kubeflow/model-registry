import { Form, FormGroup, PageSection, Stack, StackItem } from '@patternfly/react-core';
import React from 'react';
import { useNavigate } from 'react-router-dom';
import spacing from '@patternfly/react-styles/css/utilities/Spacing/spacing';
import {
  ModelLocationType,
  RegisterCatalogModelFormData,
  useRegisterCatalogModelData,
} from '~/app/pages/modelRegistry/screens/RegisterModel/useRegisterModelData';
import RegistrationCommonFormSections from '~/app/pages/modelRegistry/screens/RegisterModel/RegistrationCommonFormSections';
import {
  isModelNameExisting,
  isNameValid,
  isRegisterCatalogModelSubmitDisabled,
  registerModel,
} from '~/app/pages/modelRegistry/screens/RegisterModel/utils';
import { SubmitLabel } from '~/app/pages/modelRegistry/screens/RegisterModel/const';
import RegisterModelDetailsFormSection from '~/app/pages/modelRegistry/screens/RegisterModel/RegisterModelDetailsFormSection';
import RegistrationFormFooter from '~/app/pages/modelRegistry/screens/RegisterModel/RegistrationFormFooter';
import { ModelCatalogItem } from '~/app/modelCatalogTypes';
import { ModelRegistry, ModelRegistryMetadataType } from '~/app/types';
import { ModelRegistryContext } from '~/app/context/ModelRegistryContext';
import useRegisteredModels from '~/app/hooks/useRegisteredModels';
import { extractVersionTag } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import ModelRegistrySelector from '~/app/pages/modelRegistry/screens/ModelRegistrySelector';

interface RegisterCatalogModelFormProps {
  model: ModelCatalogItem;
  modelId?: string;
  preferredModelRegistry: ModelRegistry;
}

const RegisterCatalogModelForm: React.FC<RegisterCatalogModelFormProps> = ({
  model,
  modelId,
  preferredModelRegistry,
}) => {
  const navigate = useNavigate();
  const { apiState } = React.useContext(ModelRegistryContext);
  const [registeredModels, registeredModelsLoaded] = useRegisteredModels();

  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [submitError, setSubmitError] = React.useState<Error | undefined>(undefined);

  const versionTag = extractVersionTag(model.tags);

  // Extract framework information from model tags or description
  const getFrameworkFromModel = (catalogModel: ModelCatalogItem): string => {
    if (catalogModel.framework) {
      return catalogModel.framework;
    }
    // Try to extract from tags or description
    const tags = catalogModel.tags || [];
    const description = catalogModel.description || '';

    // Common framework patterns
    const frameworkPatterns = [
      'pytorch',
      'tensorflow',
      'onnx',
      'sklearn',
      'xgboost',
      'lightgbm',
      'huggingface',
      'transformers',
      'torch',
      'tf',
      'keras',
    ];

    for (const pattern of frameworkPatterns) {
      if (
        tags.some((tag) => tag.toLowerCase().includes(pattern)) ||
        description.toLowerCase().includes(pattern)
      ) {
        return pattern.charAt(0).toUpperCase() + pattern.slice(1);
      }
    }

    return 'PyTorch'; // Default fallback
  };

  const getVersionFromModel = (catalogModel: ModelCatalogItem): string => {
    if (catalogModel.framework) {
      // Try to extract version from tags or use default
      const tags = catalogModel.tags || [];
      const versionPattern = /(\d+\.\d+\.\d+|\d+\.\d+|\d+)/;

      for (const tag of tags) {
        const match = tag.match(versionPattern);
        if (match) {
          return match[1];
        }
      }
    }
    return '1.0.0'; // Default fallback
  };

  const initialFormData: RegisterCatalogModelFormData = {
    modelName: `${model.name}-${versionTag || ''}`,
    modelDescription: model.description || '',
    versionName: 'Version 1',
    versionDescription: '',
    sourceModelFormat: getFrameworkFromModel(model),
    sourceModelFormatVersion: getVersionFromModel(model),
    modelLocationType: ModelLocationType.URI,
    modelLocationEndpoint: '',
    modelLocationBucket: '',
    modelLocationRegion: '',
    modelLocationPath: '',
    modelLocationURI: model.url || '',
    modelRegistry: preferredModelRegistry.name,
    modelCustomProperties: {},
    versionCustomProperties: {
      License: {
        // eslint-disable-next-line camelcase
        string_value: model.license || '',
        metadataType: ModelRegistryMetadataType.STRING,
      },
      Provider: {
        // eslint-disable-next-line camelcase
        string_value: model.provider || '',
        metadataType: ModelRegistryMetadataType.STRING,
      },
      'Registered from': {
        // eslint-disable-next-line camelcase
        string_value: 'Model catalog',
        metadataType: ModelRegistryMetadataType.STRING,
      },
      'Source model': {
        // eslint-disable-next-line camelcase
        string_value: model.name,
        metadataType: ModelRegistryMetadataType.STRING,
      },
      'Source model version': {
        // eslint-disable-next-line camelcase
        string_value: versionTag || '',
        metadataType: ModelRegistryMetadataType.STRING,
      },
      'Source model id': {
        // eslint-disable-next-line camelcase
        string_value: model.id || '',
        metadataType: ModelRegistryMetadataType.STRING,
      },
      Framework: {
        // eslint-disable-next-line camelcase
        string_value: model.framework || '',
        metadataType: ModelRegistryMetadataType.STRING,
      },
      Task: {
        // eslint-disable-next-line camelcase
        string_value: model.task || '',
        metadataType: ModelRegistryMetadataType.STRING,
      },
    },
  };

  const [formData, setData] = useRegisterCatalogModelData();

  // Initialize form data from the catalog model on mount/update
  React.useEffect(() => {
    setData('modelName', initialFormData.modelName);
    setData('modelDescription', initialFormData.modelDescription);
    setData('versionName', initialFormData.versionName);
    setData('versionDescription', initialFormData.versionDescription);
    setData('sourceModelFormat', initialFormData.sourceModelFormat);
    setData('sourceModelFormatVersion', initialFormData.sourceModelFormatVersion);
    setData('modelLocationType', initialFormData.modelLocationType);
    setData('modelLocationEndpoint', initialFormData.modelLocationEndpoint);
    setData('modelLocationBucket', initialFormData.modelLocationBucket);
    setData('modelLocationRegion', initialFormData.modelLocationRegion);
    setData('modelLocationPath', initialFormData.modelLocationPath);
    setData('modelLocationURI', initialFormData.modelLocationURI);
    setData('modelRegistry', initialFormData.modelRegistry);
    setData('modelCustomProperties', initialFormData.modelCustomProperties);
    setData('versionCustomProperties', initialFormData.versionCustomProperties);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [model.id, preferredModelRegistry.name, versionTag]);

  const [submittedRegisteredModelName, setSubmittedRegisteredModelName] =
    React.useState<string>('');
  const [submittedVersionName, setSubmittedVersionName] = React.useState<string>('');
  const [registrationErrorType, setRegistrationErrorType] = React.useState<string | undefined>(
    undefined,
  );

  const isModelNameValid = isNameValid(formData.modelName);
  const isModelNameDuplicate = registeredModelsLoaded
    ? isModelNameExisting(formData.modelName, registeredModels)
    : false;
  const hasModelNameError = !isModelNameValid || isModelNameDuplicate;

  const isSubmitDisabled =
    isSubmitting || isRegisterCatalogModelSubmitDisabled(formData, registeredModels);

  const handleSubmit = async () => {
    setIsSubmitting(true);
    setSubmitError(undefined);

    // Additional validation before submission
    if (!formData.sourceModelFormat || formData.sourceModelFormat.trim() === '') {
      setSubmitError(new Error('Source model format is required'));
      setIsSubmitting(false);
      return;
    }

    if (!formData.sourceModelFormatVersion || formData.sourceModelFormatVersion.trim() === '') {
      setSubmitError(new Error('Source model format version is required'));
      setIsSubmitting(false);
      return;
    }

    if (!formData.modelLocationURI || formData.modelLocationURI.trim() === '') {
      setSubmitError(new Error('Model location URI is required'));
      setIsSubmitting(false);
      return;
    }

    // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
    if (!apiState.api) {
      setSubmitError(new Error('Model registry API is not available'));
      setIsSubmitting(false);
      return;
    }

    try {
      const {
        data: { registeredModel, modelVersion, modelArtifact },
        errors,
      } = await registerModel(apiState, formData, 'user'); // TODO: Get actual user

      if (registeredModel && modelVersion && modelArtifact) {
        const navigationPath = `/model-registry/${encodeURIComponent(preferredModelRegistry.name)}/registeredModels/${registeredModel.id}`;
        navigate(navigationPath);
      } else if (Object.keys(errors).length > 0) {
        setIsSubmitting(false);
        setSubmittedRegisteredModelName(formData.modelName);
        setSubmittedVersionName(formData.versionName);
        const resourceName = Object.keys(errors)[0];
        setRegistrationErrorType(resourceName);
        setSubmitError(errors[resourceName]);
      }
    } catch (error) {
      setSubmitError(error instanceof Error ? error : new Error('Registration failed'));
      setIsSubmitting(false);
    }
  };

  const onCancel = () => {
    navigate(`/model-catalog/${modelId || model.id}`);
  };

  return (
    <>
      <PageSection hasBodyWrapper={false} isFilled>
        <Form isWidthLimited>
          <Stack hasGutter>
            <StackItem className={spacing.mbLg}>
              <FormGroup
                id="model-registry-container"
                label="Model registry"
                isRequired
                fieldId="model-registry-name"
              >
                <ModelRegistrySelector
                  modelRegistry={formData.modelRegistry}
                  onSelection={(mr) => setData('modelRegistry', mr)}
                  primary
                  isFullWidth
                  hasError={false}
                />
              </FormGroup>
            </StackItem>
            <StackItem>
              <RegisterModelDetailsFormSection
                formData={formData}
                setData={setData}
                hasModelNameError={hasModelNameError}
                isModelNameDuplicate={isModelNameDuplicate}
                isCatalogModel
              />
              <RegistrationCommonFormSections
                formData={formData}
                setData={setData}
                isFirstVersion={false}
              />
            </StackItem>
          </Stack>
        </Form>
      </PageSection>
      <RegistrationFormFooter
        submitLabel={SubmitLabel.REGISTER_MODEL}
        submitError={submitError}
        isSubmitDisabled={isSubmitDisabled}
        isSubmitting={isSubmitting}
        onSubmit={handleSubmit}
        onCancel={onCancel}
        registrationErrorType={registrationErrorType}
        versionName={submittedVersionName}
        modelName={submittedRegisteredModelName}
      />
    </>
  );
};

export default RegisterCatalogModelForm;

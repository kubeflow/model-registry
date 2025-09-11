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
import useUser from '~/app/hooks/useUser';
import { extractVersionTag } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import ModelRegistrySelector from '~/app/pages/modelRegistry/screens/ModelRegistrySelector';
import { registeredModelUrl } from '~/app/pages/modelRegistry/screens/routeUtils';
import { getCatalogModelDetailsRoute } from '~/app/routes/modelCatalog/catalogModelDetails';
import {
  catalogParamsToModelSourceProperties,
  createCustomPropertiesFromModel,
} from '~/concepts/modelRegistry/utils';

interface RegisterCatalogModelFormProps {
  model: ModelCatalogItem;
  preferredModelRegistry: ModelRegistry;
}

const RegisterCatalogModelForm: React.FC<RegisterCatalogModelFormProps> = ({
  model,
  preferredModelRegistry,
}) => {
  const navigate = useNavigate();
  const { apiState } = React.useContext(ModelRegistryContext);
  const [registeredModels, registeredModelsLoaded] = useRegisteredModels();
  const user = useUser();

  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [submitError, setSubmitError] = React.useState<Error | undefined>(undefined);

  const versionTag = extractVersionTag(model.tags);

  const catalogParams = {
    modelName: model.name,
    tag: versionTag || '',
    sourceName: model.provider || '',
    repositoryName: model.id || '',
  };
  const sourceProperties = catalogParamsToModelSourceProperties(catalogParams);

  const initialFormData: RegisterCatalogModelFormData = {
    modelName: `${model.name}-${versionTag || ''}`,
    modelDescription: model.description || '',
    versionName: 'Version 1',
    versionDescription: '',
    sourceModelFormat: '',
    sourceModelFormatVersion: '',
    modelLocationType: ModelLocationType.URI,
    modelLocationEndpoint: '',
    modelLocationBucket: '',
    modelLocationRegion: '',
    modelLocationPath: '',
    modelLocationURI: model.url || '',
    modelRegistry: preferredModelRegistry.name,
    modelCustomProperties: createCustomPropertiesFromModel(model),
    versionCustomProperties: {
      License: {
        // eslint-disable-next-line camelcase
        string_value: model.license || '',
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
    additionalArtifactProperties: sourceProperties,
  };

  const [formData, setData] = useRegisterCatalogModelData(initialFormData);

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
      } = await registerModel(apiState, formData, user.userId || 'user');

      if (registeredModel && modelVersion && modelArtifact) {
        const navigationPath = registeredModelUrl(registeredModel.id, preferredModelRegistry.name);
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
    navigate(getCatalogModelDetailsRoute({ modelName: model.name, tag: '' }));
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
              />
              <RegistrationCommonFormSections
                formData={formData}
                setData={setData}
                isFirstVersion={false}
                isCatalogModel
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

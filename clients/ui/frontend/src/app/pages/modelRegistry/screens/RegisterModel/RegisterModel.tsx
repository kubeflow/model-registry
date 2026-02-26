import React from 'react';
import {
  Breadcrumb,
  BreadcrumbItem,
  Form,
  PageSection,
  Stack,
  StackItem,
} from '@patternfly/react-core';
import { useParams, useNavigate } from 'react-router';
import { Link } from 'react-router-dom';
import { ApplicationsPage, FormSection } from 'mod-arch-shared';
import { useThemeContext } from 'mod-arch-kubeflow';
import { useCheckNamespaceRegistryAccess } from '~/app/hooks/useCheckNamespaceRegistryAccess';
import { useModelRegistryNamespace } from '~/app/hooks/useModelRegistryNamespace';
import { modelRegistryUrl, modelVersionUrl } from '~/app/pages/modelRegistry/screens/routeUtils';
import { RegistrationMode } from '~/app/pages/modelRegistry/screens/const';
import { ModelTransferJobUploadIntent } from '~/app/types';
import { ModelRegistryContext } from '~/app/context/ModelRegistryContext';
import { AppContext } from '~/app/context/AppContext';
import useRegisteredModels from '~/app/hooks/useRegisteredModels';
import { useRegisterModelData } from './useRegisterModelData';
import {
  isModelNameExisting,
  isNameValid,
  isRegisterModelSubmitDisabled,
  registerModel,
  registerViaTransferJob,
} from './utils';
import RegistrationCommonFormSections from './RegistrationCommonFormSections';
import RegistrationFormFooter from './RegistrationFormFooter';
import { SubmitLabel, RegistrationErrorType } from './const';
import type { RegistrationInlineAlert } from './RegistrationFormFooter';
import PrefilledModelRegistryField from './PrefilledModelRegistryField';
import RegisterModelDetailsFormSection from './RegisterModelDetailsFormSection';
import { useRegistrationNotification } from './useRegistrationNotification';

const RegisterModel: React.FC = () => {
  const { modelRegistry: mrName } = useParams();
  const registryNamespace = useModelRegistryNamespace();
  const navigate = useNavigate();
  const { apiState } = React.useContext(ModelRegistryContext);
  const { user } = React.useContext(AppContext);
  const { isMUITheme } = useThemeContext();
  const author = user.userId || '';
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [submitError, setSubmitError] = React.useState<Error | undefined>(undefined);
  const [formData, setData] = useRegisterModelData();
  const [submittedRegisteredModelName, setSubmittedRegisteredModelName] =
    React.useState<string>('');
  const [submittedVersionName, setSubmittedVersionName] = React.useState<string>('');
  const [registrationErrorType, setRegistrationErrorType] = React.useState<string | undefined>(
    undefined,
  );
  const [inlineAlert, setInlineAlert] = React.useState<RegistrationInlineAlert | undefined>(
    undefined,
  );
  const registrationNotification = useRegistrationNotification(setInlineAlert);
  const [registeredModels, registeredModelsLoaded, registeredModelsLoadError] =
    useRegisteredModels();
  const {
    hasAccess: namespaceHasAccess,
    isLoading: isNamespaceAccessLoading,
    error: namespaceAccessError,
  } = useCheckNamespaceRegistryAccess(mrName, registryNamespace, formData.namespace ?? '');

  const isModelNameValid = isNameValid(formData.modelName);
  const isModelNameDuplicate = isModelNameExisting(formData.modelName, registeredModels);
  const hasModelNameError = !isModelNameValid || isModelNameDuplicate;
  const isSubmitDisabled =
    isSubmitting ||
    isRegisterModelSubmitDisabled(
      formData,
      registeredModels,
      namespaceHasAccess,
      isNamespaceAccessLoading,
    );

  const handleSubmit = async () => {
    setIsSubmitting(true);
    setSubmitError(undefined);
    setInlineAlert(undefined);

    const versionModelName = `${formData.modelName} / ${formData.versionName}`;
    const toastParams = { versionModelName, mrName: mrName ?? '' };

    // Branch based on registration mode
    if (formData.registrationMode === RegistrationMode.RegisterAndStore) {
      registrationNotification.showRegisterAndStoreSubmitting(toastParams);
      const { transferJob, error } = await registerViaTransferJob(apiState, author, {
        intent: ModelTransferJobUploadIntent.CREATE_MODEL,
        formData,
      });

      if (transferJob) {
        registrationNotification.showRegisterAndStoreSuccess(toastParams);
        navigate(modelRegistryUrl(mrName));
      } else if (error) {
        setIsSubmitting(false);
        setRegistrationErrorType(RegistrationErrorType.TRANSFER_JOB);
        setSubmitError(error);
        registrationNotification.showRegisterAndStoreError(toastParams);
      }
    } else {
      // Register mode: Existing synchronous registration flow
      const {
        data: { registeredModel, modelVersion, modelArtifact },
        errors,
      } = await registerModel(apiState, formData, author);

      if (registeredModel && modelVersion && modelArtifact) {
        navigate(modelVersionUrl(modelVersion.id, registeredModel.id, mrName));
      } else if (Object.keys(errors).length > 0) {
        setIsSubmitting(false);
        setSubmittedRegisteredModelName(formData.modelName);
        setSubmittedVersionName(formData.versionName);
        const resourceName = Object.keys(errors)[0];
        setRegistrationErrorType(resourceName);
        setSubmitError(errors[resourceName]);
      }
    }
  };
  const onCancel = () => {
    navigate(modelRegistryUrl(mrName));
  };

  return (
    <ApplicationsPage
      title="Register model"
      description="Create and register the first version of a new model."
      breadcrumb={
        <Breadcrumb>
          <BreadcrumbItem
            render={() => <Link to={modelRegistryUrl(mrName)}>Model registry - {mrName}</Link>}
          />
          <BreadcrumbItem>Register model</BreadcrumbItem>
        </Breadcrumb>
      }
      loaded={registeredModelsLoaded}
      loadError={registeredModelsLoadError}
      empty={false}
    >
      <PageSection hasBodyWrapper={false} isFilled>
        <Form isWidthLimited>
          <Stack hasGutter>
            <FormSection className="pf-v6-u-pb-xl">
              <PrefilledModelRegistryField mrName={mrName} />
            </FormSection>
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
                isFirstVersion
                namespaceHasAccess={namespaceHasAccess}
                isNamespaceAccessLoading={isNamespaceAccessLoading}
                namespaceAccessError={namespaceAccessError}
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
        inlineAlert={!isMUITheme ? inlineAlert : undefined}
      />
    </ApplicationsPage>
  );
};

export default RegisterModel;

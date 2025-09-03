import React from 'react';
import {
  Breadcrumb,
  BreadcrumbItem,
  Form,
  PageSection,
  Stack,
  StackItem,
} from '@patternfly/react-core';
import spacing from '@patternfly/react-styles/css/utilities/Spacing/spacing';
import { useParams, useNavigate } from 'react-router';
import { Link } from 'react-router-dom';
import { ApplicationsPage } from 'mod-arch-shared';
import { modelRegistryUrl, modelVersionUrl } from '~/app/pages/modelRegistry/screens/routeUtils';
import { ModelRegistryContext } from '~/app/context/ModelRegistryContext';
import { AppContext } from '~/app/context/AppContext';
import useRegisteredModels from '~/app/hooks/useRegisteredModels';
import { useRegisterModelData } from './useRegisterModelData';
import {
  isModelNameExisting,
  isNameValid,
  isRegisterModelSubmitDisabled,
  registerModel,
} from './utils';
import RegistrationCommonFormSections from './RegistrationCommonFormSections';
import RegistrationFormFooter from './RegistrationFormFooter';
import { SubmitLabel } from './const';
import PrefilledModelRegistryField from './PrefilledModelRegistryField';
import RegisterModelDetailsFormSection from './RegisterModelDetailsFormSection';

const RegisterModel: React.FC = () => {
  const { modelRegistry: mrName } = useParams();
  const navigate = useNavigate();
  const { apiState } = React.useContext(ModelRegistryContext);
  const { user } = React.useContext(AppContext);
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
  const [registeredModels, registeredModelsLoaded, registeredModelsLoadError] =
    useRegisteredModels();

  const isModelNameValid = isNameValid(formData.modelName);
  const isModelNameDuplicate = isModelNameExisting(formData.modelName, registeredModels);
  const hasModelNameError = !isModelNameValid || isModelNameDuplicate;
  const isSubmitDisabled =
    isSubmitting || isRegisterModelSubmitDisabled(formData, registeredModels);

  const handleSubmit = async () => {
    setIsSubmitting(true);
    setSubmitError(undefined);

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
            <StackItem className={spacing.mbLg}>
              <PrefilledModelRegistryField mrName={mrName} />
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
                isFirstVersion
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
    </ApplicationsPage>
  );
};

export default RegisterModel;

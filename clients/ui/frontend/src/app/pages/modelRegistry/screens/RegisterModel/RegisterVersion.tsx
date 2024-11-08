import React from 'react';
import {
  Breadcrumb,
  BreadcrumbItem,
  Form,
  FormGroup,
  PageSection,
  Spinner,
  Stack,
  StackItem,
} from '@patternfly/react-core';
import spacing from '@patternfly/react-styles/css/utilities/Spacing/spacing';
import { useParams, useNavigate } from 'react-router';
import { Link } from 'react-router-dom';
import ApplicationsPage from '~/shared/components/ApplicationsPage';
import { modelRegistryUrl, registeredModelUrl } from '~/app/pages/modelRegistry/screens/routeUtils';
import useRegisteredModels from '~/app/hooks/useRegisteredModels';
import { ValueOf } from '~/shared/typeHelpers';
import { filterLiveModels } from '~/app/pages/modelRegistry/screens/utils';
import { RegistrationCommonFormData, useRegisterVersionData } from './useRegisterModelData';
import { isRegisterVersionSubmitDisabled, registerVersion } from './utils';
import RegistrationCommonFormSections from './RegistrationCommonFormSections';
import { useRegistrationCommonState } from './useRegistrationCommonState';
import PrefilledModelRegistryField from './PrefilledModelRegistryField';
import RegistrationFormFooter from './RegistrationFormFooter';
import RegisteredModelSelector from './RegisteredModelSelector';
import { usePrefillRegisterVersionFields } from './usePrefillRegisterVersionFields';

const RegisterVersion: React.FC = () => {
  const { modelRegistry: mrName, registeredModelId: prefilledRegisteredModelId } = useParams();

  const navigate = useNavigate();

  const { isSubmitting, submitError, setSubmitError, handleSubmit, apiState, author } =
    useRegistrationCommonState();

  const [formData, setData] = useRegisterVersionData(prefilledRegisteredModelId);
  const isSubmitDisabled = isSubmitting || isRegisterVersionSubmitDisabled(formData);
  const { registeredModelId } = formData;

  const [allRegisteredModels, loadedRegisteredModels, loadRegisteredModelsError] =
    useRegisteredModels();
  const liveRegisteredModels = filterLiveModels(allRegisteredModels.items);
  const registeredModel = liveRegisteredModels.find(({ id }) => id === registeredModelId);

  const { loadedPrefillData, loadPrefillDataError, latestVersion } =
    usePrefillRegisterVersionFields({
      registeredModel,
      setData,
    });

  const onSubmit = () => {
    if (!registeredModel) {
      return; // We shouldn't be able to hit this due to form validation
    }
    handleSubmit(async () => {
      await registerVersion(apiState, registeredModel, formData, author);
      navigate(registeredModelUrl(registeredModel.id, mrName));
    });
  };
  const onCancel = () =>
    navigate(
      prefilledRegisteredModelId && registeredModel
        ? registeredModelUrl(registeredModel.id, mrName)
        : modelRegistryUrl(mrName),
    );

  return (
    <ApplicationsPage
      title="Register new version"
      description="Register a latest version to the model you selected below."
      breadcrumb={
        <Breadcrumb>
          <BreadcrumbItem
            render={() => <Link to={modelRegistryUrl(mrName)}>Model registry - {mrName}</Link>}
          />
          {prefilledRegisteredModelId && registeredModel && (
            <BreadcrumbItem
              render={() => (
                <Link to={registeredModelUrl(registeredModel.id, mrName)}>
                  {registeredModel.name}
                </Link>
              )}
            />
          )}
          <BreadcrumbItem>Register new version</BreadcrumbItem>
        </Breadcrumb>
      }
      loadError={loadRegisteredModelsError || loadPrefillDataError}
      // Data for prefilling is refetched when the model selection changes, so we don't handle its loaded state here.
      // Instead we show a spinner in RegisteredModelSelector after that selection changes.
      loaded={loadedRegisteredModels}
      empty={false}
    >
      <PageSection hasBodyWrapper={false} isFilled>
        <Form isWidthLimited>
          <Stack hasGutter>
            <StackItem>
              <PrefilledModelRegistryField mrName={mrName} />
            </StackItem>
            <StackItem className={spacing.mbLg}>
              <FormGroup
                id="registered-model-container"
                isRequired
                fieldId="model-name"
                labelHelp={
                  !loadedPrefillData ? <Spinner size="sm" className={spacing.mlMd} /> : undefined
                }
              >
                <RegisteredModelSelector
                  registeredModels={liveRegisteredModels}
                  registeredModelId={registeredModelId}
                  setRegisteredModelId={(id) => setData('registeredModelId', id)}
                  isDisabled={!!prefilledRegisteredModelId}
                />
              </FormGroup>
            </StackItem>
            <StackItem>
              <RegistrationCommonFormSections
                formData={formData}
                setData={(
                  propKey: keyof RegistrationCommonFormData,
                  propValue: ValueOf<RegistrationCommonFormData>,
                ) => setData(propKey, propValue)}
                isFirstVersion={false}
                latestVersion={latestVersion}
              />
            </StackItem>
          </Stack>
        </Form>
      </PageSection>
      <RegistrationFormFooter
        submitLabel="Register new version"
        submitError={submitError}
        setSubmitError={setSubmitError}
        isSubmitDisabled={isSubmitDisabled}
        isSubmitting={isSubmitting}
        onSubmit={onSubmit}
        onCancel={onCancel}
      />
    </ApplicationsPage>
  );
};

export default RegisterVersion;

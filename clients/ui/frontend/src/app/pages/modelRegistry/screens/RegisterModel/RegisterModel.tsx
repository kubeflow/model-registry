import React from 'react';
import {
  Breadcrumb,
  BreadcrumbItem,
  Form,
  FormGroup,
  PageSection,
  Stack,
  StackItem,
  TextArea,
  TextInput,
} from '@patternfly/react-core';
import spacing from '@patternfly/react-styles/css/utilities/Spacing/spacing';
import { useParams, useNavigate } from 'react-router';
import { Link } from 'react-router-dom';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import FormSection from '~/shared/components/pf-overrides/FormSection';
import ApplicationsPage from '~/shared/components/ApplicationsPage';
import { modelRegistryUrl, registeredModelUrl } from '~/app/pages/modelRegistry/screens/routeUtils';
import { ValueOf } from '~/shared/typeHelpers';
import { isMUITheme } from '~/shared/utilities/const';
import { useRegisterModelData, RegistrationCommonFormData } from './useRegisterModelData';
import { isRegisterModelSubmitDisabled, registerModel } from './utils';
import { useRegistrationCommonState } from './useRegistrationCommonState';
import RegistrationCommonFormSections from './RegistrationCommonFormSections';
import RegistrationFormFooter from './RegistrationFormFooter';

const RegisterModel: React.FC = () => {
  const { modelRegistry: mrName } = useParams();
  const navigate = useNavigate();

  const { isSubmitting, submitError, setSubmitError, handleSubmit, apiState, author } =
    useRegistrationCommonState();

  const [formData, setData] = useRegisterModelData();
  const isSubmitDisabled = isSubmitting || isRegisterModelSubmitDisabled(formData);
  const { modelName, modelDescription } = formData;

  const onSubmit = () =>
    handleSubmit(async () => {
      const { registeredModel } = await registerModel(apiState, formData, author);
      navigate(registeredModelUrl(registeredModel.id, mrName));
    });
  const onCancel = () => navigate(modelRegistryUrl(mrName));

  const modelRegistryInput = (
    <TextInput isDisabled isRequired type="text" id="mr-name" name="mr-name" value={mrName} />
  );

  const modelNameInput = (
    <TextInput
      isRequired
      type="text"
      id="model-name"
      name="model-name"
      value={modelName}
      onChange={(_e, value) => setData('modelName', value)}
    />
  );

  const modelDescriptionInput = (
    <TextArea
      type="text"
      id="model-description"
      name="model-description"
      value={modelDescription}
      onChange={(_e, value) => setData('modelDescription', value)}
    />
  );

  return (
    <ApplicationsPage
      title="Register model"
      description="Create a new model and register the first version of your new model."
      breadcrumb={
        <Breadcrumb>
          <BreadcrumbItem
            render={() => <Link to={modelRegistryUrl(mrName)}>Model registry - {mrName}</Link>}
          />
          <BreadcrumbItem>Register model</BreadcrumbItem>
        </Breadcrumb>
      }
      loaded
      empty={false}
    >
      <PageSection hasBodyWrapper={false} isFilled>
        <Form isWidthLimited>
          <Stack hasGutter>
            <StackItem className={spacing.mbLg}>
              <FormGroup
                className="form-group-disabled"
                label="Model registry"
                isRequired
                fieldId="mr-name"
              >
                {isMUITheme() ? (
                  <FormFieldset component={modelRegistryInput} field="Model Registry" />
                ) : (
                  modelRegistryInput
                )}
              </FormGroup>
            </StackItem>
            <StackItem>
              <FormSection
                title="Model details"
                description="Provide general details that apply to all versions of this model."
              >
                <FormGroup label="Model name" isRequired fieldId="model-name">
                  {isMUITheme() ? (
                    <FormFieldset component={modelNameInput} field="Model Name" />
                  ) : (
                    modelNameInput
                  )}
                </FormGroup>
                <FormGroup
                  className="model-description"
                  label="Model description"
                  fieldId="model-description"
                >
                  {isMUITheme() ? (
                    <FormFieldset component={modelDescriptionInput} field="Model Description" />
                  ) : (
                    modelDescriptionInput
                  )}
                </FormGroup>
              </FormSection>
              <RegistrationCommonFormSections
                formData={formData}
                setData={(
                  propKey: keyof RegistrationCommonFormData,
                  propValue: ValueOf<RegistrationCommonFormData>,
                ) => setData(propKey, propValue)}
                isFirstVersion
              />
            </StackItem>
          </Stack>
        </Form>
      </PageSection>
      <RegistrationFormFooter
        submitLabel="Register model"
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

export default RegisterModel;

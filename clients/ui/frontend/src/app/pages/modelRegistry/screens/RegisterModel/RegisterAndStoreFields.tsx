import React from 'react';
import { HelperText, HelperTextItem } from '@patternfly/react-core';
import { UpdateObjectAtPropAndValue } from 'mod-arch-shared';
import FormSection from '~/app/pages/modelRegistry/components/pf-overrides/FormSection';
import K8sNameDescriptionField from '~/concepts/k8s/K8sNameDescriptionField/K8sNameDescriptionField';
import { RegistrationCommonFormData } from './useRegisterModelData';
import RegistrationModelLocationFields from './RegistrationModelLocationFields';

type RegisterAndStoreFieldsProps<D extends RegistrationCommonFormData> = {
  formData: D;
  setData: UpdateObjectAtPropAndValue<D>;
  isCatalogModel?: boolean;
};

const RegisterAndStoreFields = <D extends RegistrationCommonFormData>({
  formData,
  setData,
  isCatalogModel,
}: RegisterAndStoreFieldsProps<D>): React.ReactNode => (
  <>
    <K8sNameDescriptionField
      dataTestId="model-transfer-job"
      nameLabel="Model transfer job name"
      data={{
        name: formData.jobName,
        description: '',
      }}
      onDataChange={(data) => {
        setData('jobName', data.name);
      }}
      resourceNameHelperText={
        <>
          <HelperText>
            <HelperTextItem>Cannot exceed 30 characters.</HelperTextItem>
            <HelperTextItem>
              Must start and end with a letter or number. Valid characters include lowercase
              letters, numbers, and hyphens (-).
            </HelperTextItem>
            <HelperTextItem>
              Auto generated value will be used as resource name if field is blank.
            </HelperTextItem>
          </HelperText>
        </>
      }
      hideDescription
    />
    {/*
      TODO add a namespace selector - don't replicate the ODH notion of "projects", we will start with a simple k8s namespace selector.
      Needs to list all namespaces the user can see, which is something we already have in the app header here, look how that was done.

      TODO hide the rest of this section until a namespace is selected, once that's implemented.
    */}
    <FormSection
      title="Model origin location"
      description="Specify the location that is currently being used to store the model."
    >
      <RegistrationModelLocationFields
        formData={formData}
        setData={setData}
        isCatalogModel={isCatalogModel}
      />
    </FormSection>
    <FormSection
      title="Model destination location"
      description="Specify the location that will be used to store the registered model."
    >
      TODO destination location fields here
      {/* 
        TODO this will need the inputs we require to create the secret required by the job,
        but in a generic way - don't recreate the whole ODH create/select connection flow
        here. start with bare minimum simple inputs.
        Jobs are documented here: https://github.com/kubeflow/model-registry/tree/main/jobs/async-upload
      */}
    </FormSection>
  </>
);

export default RegisterAndStoreFields;

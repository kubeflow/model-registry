import React from 'react';
import { UpdateObjectAtPropAndValue } from 'mod-arch-shared';
import FormSection from '~/app/pages/modelRegistry/components/pf-overrides/FormSection';
import { RegistrationCommonFormData } from './useRegisterModelData';
import RegistrationModelLocationFields from './RegistrationModelLocationFields';
import RegistrationDestinationLocationFields from './RegistrationDestinationLocationFields';

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
    TODO job name and namespace fields here
    {/*
      TODO use the K8sNameResourceField component here for the job name.

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
        includeCredentialFields
      />
    </FormSection>
    <FormSection
      title="Model destination location"
      description="Specify the OCI registry location that will be used to store the registered model."
    >
      <RegistrationDestinationLocationFields formData={formData} setData={setData} />
    </FormSection>
  </>
);

export default RegisterAndStoreFields;

import React from 'react';
import { UpdateObjectAtPropAndValue } from 'mod-arch-shared';
import { FormGroup } from '@patternfly/react-core';

import { useThemeContext } from 'mod-arch-kubeflow';
import FormSection from '~/app/pages/modelRegistry/components/pf-overrides/FormSection';
import NamespaceSelector from '~/app/standalone/NamespaceSelector';
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
}: RegisterAndStoreFieldsProps<D>): React.ReactNode => {
  const { isMUITheme } = useThemeContext();

  const handleNamespaceSelect = (namespace: string) => {
    setData('namespace', namespace);
  };

  const namespaceSelectorElement = (
    <NamespaceSelector
      placeholderText="Select a namespace"
      onSelect={handleNamespaceSelect}
      selectedNamespace={formData.namespace}
      isFullWidth
      ignoreMandatoryNamespace
    />
  );

  return (
    <>
      TODO job name field here
      {/*
      TODO use the K8sNameResourceField component here for the job name.

      */}
      {isMUITheme ? (
        <FormSection title="Namespace">{namespaceSelectorElement}</FormSection>
      ) : (
        <FormGroup label="Project" data-testid="namespace-form-group" isRequired>
          {namespaceSelectorElement}
        </FormGroup>
      )}
      {formData.namespace && (
        <>
          <FormSection
            data-testid="model-origin-location-section"
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
            data-testid="model-destination-location-section"
            description="Specify the OCI registry location that will be used to store the registered model."
          >
            <RegistrationDestinationLocationFields formData={formData} setData={setData} />
          </FormSection>
        </>
      )}
    </>
  );
};

export default RegisterAndStoreFields;

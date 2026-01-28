import React from 'react';
import { FormGroup, TextInput } from '@patternfly/react-core';
import { UpdateObjectAtPropAndValue } from 'mod-arch-shared';
import PasswordInput from '~/app/shared/components/PasswordInput';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import { RegistrationCommonFormData } from './useRegisterModelData';

type RegistrationDestinationLocationFieldsProps<D extends RegistrationCommonFormData> = {
  formData: D;
  setData: UpdateObjectAtPropAndValue<D>;
};

const RegistrationDestinationLocationFields = <D extends RegistrationCommonFormData>({
  formData,
  setData,
}: RegistrationDestinationLocationFieldsProps<D>): React.ReactNode => {
  const {
    destinationOciRegistry,
    destinationOciUsername,
    destinationOciPassword,
    destinationOciUri,
    destinationOciEmail,
  } = formData;

  // OCI fields
  const ociRegistryInput = (
    <TextInput
      isRequired
      type="text"
      id="destination-oci-registry"
      name="destination-oci-registry"
      value={destinationOciRegistry}
      onChange={(_e, value) => setData('destinationOciRegistry', value)}
    />
  );

  const ociUsernameInput = (
    <TextInput
      isRequired
      type="text"
      id="destination-oci-username"
      name="destination-oci-username"
      value={destinationOciUsername}
      onChange={(_e, value) => setData('destinationOciUsername', value)}
    />
  );

  const ociPasswordInput = (
    <PasswordInput
      isRequired
      id="destination-oci-password"
      name="destination-oci-password"
      value={destinationOciPassword}
      onChange={(_e, value) => setData('destinationOciPassword', value)}
    />
  );

  const ociUriInput = (
    <TextInput
      isRequired
      type="text"
      id="destination-oci-uri"
      name="destination-oci-uri"
      value={destinationOciUri}
      onChange={(_e, value) => setData('destinationOciUri', value)}
    />
  );

  const ociEmailInput = (
    <TextInput
      type="email"
      id="destination-oci-email"
      name="destination-oci-email"
      value={destinationOciEmail}
      onChange={(_e, value) => setData('destinationOciEmail', value)}
    />
  );

  return (
    <>
      <FormGroup label="Registry" isRequired fieldId="destination-oci-registry">
        <FormFieldset component={ociRegistryInput} field="Registry" />
      </FormGroup>
      <FormGroup label="URI" isRequired fieldId="destination-oci-uri">
        <FormFieldset component={ociUriInput} field="URI" />
      </FormGroup>
      <FormGroup label="Username" isRequired fieldId="destination-oci-username">
        <FormFieldset component={ociUsernameInput} field="Username" />
      </FormGroup>
      <FormGroup label="Email" fieldId="destination-oci-email">
        <FormFieldset component={ociEmailInput} field="Email" />
      </FormGroup>
      <FormGroup label="Password" isRequired fieldId="destination-oci-password">
        <FormFieldset component={ociPasswordInput} field="Password" />
      </FormGroup>
    </>
  );
};

export default RegistrationDestinationLocationFields;

import React from 'react';
import {
  FormGroup,
  TextInput,
  Radio,
  HelperText,
  HelperTextItem,
  TextInputGroupMain,
  TextInputGroup,
} from '@patternfly/react-core';
import spacing from '@patternfly/react-styles/css/utilities/Spacing/spacing';
import { UpdateObjectAtPropAndValue } from 'mod-arch-shared';
import PasswordInput from '~/app/shared/components/PasswordInput';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import { DestinationStorageType, RegistrationCommonFormData } from './useRegisterModelData';

type RegistrationDestinationLocationFieldsProps<D extends RegistrationCommonFormData> = {
  formData: D;
  setData: UpdateObjectAtPropAndValue<D>;
};

const RegistrationDestinationLocationFields = <D extends RegistrationCommonFormData>({
  formData,
  setData,
}: RegistrationDestinationLocationFieldsProps<D>): React.ReactNode => {
  const {
    destinationStorageType,
    destinationS3AccessKeyId,
    destinationS3SecretAccessKey,
    destinationS3Endpoint,
    destinationS3Bucket,
    destinationS3Region,
    destinationS3Path,
    destinationOciRegistry,
    destinationOciUsername,
    destinationOciPassword,
    destinationOciUri,
    destinationOciEmail,
  } = formData;

  // S3 fields
  const s3AccessKeyIdInput = (
    <TextInput
      isRequired
      type="text"
      id="destination-s3-access-key-id"
      name="destination-s3-access-key-id"
      value={destinationS3AccessKeyId}
      onChange={(_e, value) => setData('destinationS3AccessKeyId', value)}
    />
  );

  const s3SecretAccessKeyInput = (
    <PasswordInput
      isRequired
      id="destination-s3-secret-access-key"
      name="destination-s3-secret-access-key"
      value={destinationS3SecretAccessKey}
      onChange={(_e, value) => setData('destinationS3SecretAccessKey', value)}
    />
  );

  const s3EndpointInput = (
    <TextInput
      type="text"
      id="destination-s3-endpoint"
      name="destination-s3-endpoint"
      value={destinationS3Endpoint}
      onChange={(_e, value) => setData('destinationS3Endpoint', value)}
    />
  );

  const s3BucketInput = (
    <TextInput
      type="text"
      id="destination-s3-bucket"
      name="destination-s3-bucket"
      value={destinationS3Bucket}
      onChange={(_e, value) => setData('destinationS3Bucket', value)}
    />
  );

  const s3RegionInput = (
    <TextInput
      type="text"
      id="destination-s3-region"
      name="destination-s3-region"
      value={destinationS3Region}
      onChange={(_e, value) => setData('destinationS3Region', value)}
    />
  );

  const s3PathInput = (
    <TextInputGroup>
      <TextInputGroupMain
        icon="/"
        type="text"
        id="destination-s3-path"
        name="destination-s3-path"
        value={destinationS3Path}
        onChange={(_e, value) => setData('destinationS3Path', value)}
      />
    </TextInputGroup>
  );

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
      <Radio
        isChecked={destinationStorageType === DestinationStorageType.S3}
        name="destination-storage-type-s3"
        onChange={() => {
          setData('destinationStorageType', DestinationStorageType.S3);
        }}
        label="S3"
        id="destination-storage-type-s3"
      />
      {destinationStorageType === DestinationStorageType.S3 && (
        <>
          <FormGroup
            className={spacing.mlLg}
            label="Endpoint"
            isRequired
            fieldId="destination-s3-endpoint"
          >
            <FormFieldset component={s3EndpointInput} field="Endpoint" />
          </FormGroup>
          <FormGroup
            className={spacing.mlLg}
            label="Bucket"
            isRequired
            fieldId="destination-s3-bucket"
          >
            <FormFieldset component={s3BucketInput} field="Bucket" />
          </FormGroup>
          <FormGroup className={spacing.mlLg} label="Region" fieldId="destination-s3-region">
            <FormFieldset component={s3RegionInput} field="Region" />
          </FormGroup>
          <FormGroup
            className={`destination-s3-path ${spacing.mlLg}`}
            label="Path"
            isRequired
            fieldId="destination-s3-path"
          >
            <FormFieldset component={s3PathInput} field="Path" />
            <HelperText>
              <HelperTextItem>
                Enter a path to a model or folder. This path cannot point to a root folder.
              </HelperTextItem>
            </HelperText>
          </FormGroup>
          <FormGroup
            className={spacing.mlLg}
            label="Access Key ID"
            isRequired
            fieldId="destination-s3-access-key-id"
          >
            <FormFieldset component={s3AccessKeyIdInput} field="Access Key ID" />
          </FormGroup>
          <FormGroup
            className={spacing.mlLg}
            label="Secret Access Key"
            isRequired
            fieldId="destination-s3-secret-access-key"
          >
            <FormFieldset component={s3SecretAccessKeyInput} field="Secret Access Key" />
          </FormGroup>
        </>
      )}
      <Radio
        isChecked={destinationStorageType === DestinationStorageType.OCI}
        name="destination-storage-type-oci"
        onChange={() => {
          setData('destinationStorageType', DestinationStorageType.OCI);
        }}
        label="OCI Registry"
        id="destination-storage-type-oci"
      />
      {destinationStorageType === DestinationStorageType.OCI && (
        <>
          <FormGroup
            className={spacing.mlLg}
            label="Registry"
            isRequired
            fieldId="destination-oci-registry"
          >
            <FormFieldset component={ociRegistryInput} field="Registry" />
          </FormGroup>
          <FormGroup className={spacing.mlLg} label="URI" isRequired fieldId="destination-oci-uri">
            <FormFieldset component={ociUriInput} field="URI" />
          </FormGroup>
          <FormGroup
            className={spacing.mlLg}
            label="Username"
            isRequired
            fieldId="destination-oci-username"
          >
            <FormFieldset component={ociUsernameInput} field="Username" />
          </FormGroup>
          <FormGroup
            className={spacing.mlLg}
            label="Email"
            isRequired
            fieldId="destination-oci-email"
          >
            <FormFieldset component={ociEmailInput} field="Email" />
          </FormGroup>
          <FormGroup
            className={spacing.mlLg}
            label="Password"
            isRequired
            fieldId="destination-oci-password"
          >
            <FormFieldset component={ociPasswordInput} field="Password" />
          </FormGroup>
        </>
      )}
    </>
  );
};

export default RegistrationDestinationLocationFields;

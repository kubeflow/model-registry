import React from 'react';
import { FormGroup, TextInput, Radio } from '@patternfly/react-core';
import spacing from '@patternfly/react-styles/css/utilities/Spacing/spacing';
import { UpdateObjectAtPropAndValue } from 'mod-arch-shared';
import PasswordInput from '~/app/shared/components/PasswordInput';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import {
  DestinationStorageType,
  ModelLocationType,
  RegistrationCommonFormData,
} from './useRegisterModelData';
import RegistrationModelLocationFields from './RegistrationModelLocationFields';

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
    destinationOciRegistry,
    destinationOciUsername,
    destinationOciPassword,
    destinationOciUri,
    destinationOciEmail,
  } = formData;

  // S3 common fields (endpoint, bucket, region, path) from RegistrationModelLocationFields
  const destinationFormData: D = {
    ...formData,
    modelLocationType: ModelLocationType.ObjectStorage,
    modelLocationEndpoint: formData.destinationS3Endpoint,
    modelLocationBucket: formData.destinationS3Bucket,
    modelLocationRegion: formData.destinationS3Region,
    modelLocationPath: formData.destinationS3Path,
    modelLocationURI: formData.modelLocationURI,
  };

  const destinationSetData: UpdateObjectAtPropAndValue<D> = (key: keyof D, value: D[keyof D]) => {
    if (key === 'modelLocationEndpoint') {
      setData('destinationS3Endpoint', String(value));
    } else if (key === 'modelLocationBucket') {
      setData('destinationS3Bucket', String(value));
    } else if (key === 'modelLocationRegion') {
      setData('destinationS3Region', String(value));
    } else if (key === 'modelLocationPath') {
      setData('destinationS3Path', String(value));
    } else {
      setData(key, value);
    }
  };

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
          <RegistrationModelLocationFields
            formData={destinationFormData}
            setData={destinationSetData}
            hideRadioButtons
          />
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

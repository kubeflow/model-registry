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
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import { ModelLocationType, RegistrationCommonFormData } from './useRegisterModelData';

type RegistrationModelLocationFieldsProps<D extends RegistrationCommonFormData> = {
  formData: D;
  setData: UpdateObjectAtPropAndValue<D>;
  isCatalogModel?: boolean;
};

const RegistrationModelLocationFields = <D extends RegistrationCommonFormData>({
  formData,
  setData,
  isCatalogModel,
}: RegistrationModelLocationFieldsProps<D>): React.ReactNode => {
  const {
    modelLocationType,
    modelLocationEndpoint,
    modelLocationBucket,
    modelLocationRegion,
    modelLocationPath,
    modelLocationURI,
  } = formData;

  const endpointInput = (
    <TextInput
      isRequired
      type="text"
      id="location-endpoint"
      name="location-endpoint"
      value={modelLocationEndpoint}
      onChange={(_e, value) => setData('modelLocationEndpoint', value)}
    />
  );

  const bucketInput = (
    <TextInput
      isRequired
      type="text"
      id="location-bucket"
      name="location-bucket"
      value={modelLocationBucket}
      onChange={(_e, value) => setData('modelLocationBucket', value)}
    />
  );

  const regionInput = (
    <TextInput
      type="text"
      id="location-region"
      name="location-region"
      value={modelLocationRegion}
      onChange={(_e, value) => setData('modelLocationRegion', value)}
    />
  );

  const pathInput = (
    <TextInputGroup>
      <TextInputGroupMain
        icon="/"
        type="text"
        id="location-path"
        name="location-path"
        value={modelLocationPath}
        onChange={(_e, value) => setData('modelLocationPath', value)}
      />
    </TextInputGroup>
  );

  const uriInput = (
    <TextInput
      isRequired
      type="text"
      id="location-uri"
      name="location-uri"
      value={modelLocationURI}
      onChange={(_e, value) => setData('modelLocationURI', value)}
      isDisabled={isCatalogModel}
    />
  );

  return (
    <>
      <Radio
        isChecked={modelLocationType === ModelLocationType.ObjectStorage}
        name="location-type-object-storage"
        isDisabled={isCatalogModel}
        onChange={() => {
          setData('modelLocationType', ModelLocationType.ObjectStorage);
        }}
        label="Object storage"
        id="location-type-object-storage"
      />
      {modelLocationType === ModelLocationType.ObjectStorage && (
        <>
          <FormGroup
            className={spacing.mlLg}
            label="Endpoint"
            isRequired
            fieldId="location-endpoint"
          >
            <FormFieldset component={endpointInput} field="Endpoint" />
          </FormGroup>
          <FormGroup className={spacing.mlLg} label="Bucket" isRequired fieldId="location-bucket">
            <FormFieldset component={bucketInput} field="Bucket" />
          </FormGroup>
          <FormGroup className={spacing.mlLg} label="Region" fieldId="location-region">
            <FormFieldset component={regionInput} field="Region" />
          </FormGroup>
          <FormGroup
            className={`location-path` + ` ${spacing.mlLg}`}
            label="Path"
            isRequired
            fieldId="location-path"
          >
            <FormFieldset component={pathInput} field="Path" />
            <HelperText>
              <HelperTextItem>
                Enter a path to a model or folder. This path cannot point to a root folder.
              </HelperTextItem>
            </HelperText>
          </FormGroup>
        </>
      )}
      <Radio
        isChecked={modelLocationType === ModelLocationType.URI}
        name="location-type-uri"
        onChange={() => {
          setData('modelLocationType', ModelLocationType.URI);
        }}
        label="URI"
        id="location-type-uri"
      />
      {modelLocationType === ModelLocationType.URI && (
        <>
          <FormGroup className={spacing.mlLg} label="URI" isRequired fieldId="location-uri">
            <FormFieldset component={uriInput} field="URI" />
          </FormGroup>
        </>
      )}
    </>
  );
};

export default RegistrationModelLocationFields;

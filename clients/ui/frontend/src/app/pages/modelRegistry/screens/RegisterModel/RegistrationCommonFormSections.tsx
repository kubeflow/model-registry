import React from 'react';
import {
  FormGroup,
  TextInput,
  TextArea,
  Radio,
  HelperText,
  HelperTextItem,
  FormHelperText,
  TextInputGroupMain,
  TextInputGroup,
} from '@patternfly/react-core';
import spacing from '@patternfly/react-styles/css/utilities/Spacing/spacing';
import { FormSection, UpdateObjectAtPropAndValue } from 'mod-arch-shared';
// import { DataConnection, UpdateObjectAtPropAndValue } from '~/pages/projects/types';
// import { convertAWSSecretData } from '~/pages/projects/screens/detail/data-connections/utils';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import { ModelVersion } from '~/app/types';
import { ModelLocationType, RegistrationCommonFormData } from './useRegisterModelData';
import { isNameValid } from './utils';
import { MR_CHARACTER_LIMIT } from './const';
// import { ConnectionModal } from './ConnectionModal';

type RegistrationCommonFormSectionsProps<D extends RegistrationCommonFormData> = {
  formData: D;
  setData: UpdateObjectAtPropAndValue<D>;
  isFirstVersion: boolean;
  latestVersion?: ModelVersion;
  isCatalogModel?: boolean;
};

const RegistrationCommonFormSections = <D extends RegistrationCommonFormData>({
  formData,
  setData,
  isFirstVersion,
  latestVersion,
  isCatalogModel,
}: RegistrationCommonFormSectionsProps<D>): React.ReactNode => {
  // const [isAutofillModalOpen, setAutofillModalOpen] = React.useState(false);
  const isVersionNameValid = isNameValid(formData.versionName);

  // const connectionDataMap: Record<
  //   string,
  //   keyof Pick<
  //     RegistrationCommonFormData,
  //     'modelLocationEndpoint' | 'modelLocationBucket' | 'modelLocationRegion'
  //   >
  // > = {
  //   AWS_S3_ENDPOINT: 'modelLocationEndpoint',
  //   AWS_S3_BUCKET: 'modelLocationBucket',
  //   AWS_DEFAULT_REGION: 'modelLocationRegion',
  // };

  // const fillObjectStorageByConnection = (connection: DataConnection) => {
  //   convertAWSSecretData(connection).forEach((dataItem) => {
  //     setData(connectionDataMap[dataItem.key], dataItem.value);
  //   });
  // };

  const {
    versionName,
    versionDescription,
    sourceModelFormat,
    sourceModelFormatVersion,
    modelLocationType,
    modelLocationEndpoint,
    modelLocationBucket,
    modelLocationRegion,
    modelLocationPath,
    modelLocationURI,
  } = formData;

  const versionNameInput = (
    <TextInput
      isRequired
      type="text"
      id="version-name"
      name="version-name"
      value={versionName}
      onChange={(_e, value) => setData('versionName', value)}
      validated={isVersionNameValid ? 'default' : 'error'}
    />
  );

  const versionDescriptionInput = (
    <TextArea
      type="text"
      id="version-description"
      name="version-description"
      value={versionDescription}
      onChange={(_e, value) => setData('versionDescription', value)}
    />
  );

  const sourceModelFormatInput = (
    <TextInput
      type="text"
      placeholder="Example, tensorflow"
      id="source-model-format"
      name="source-model-format"
      value={sourceModelFormat}
      onChange={(_e, value) => setData('sourceModelFormat', value)}
    />
  );

  const sourceModelFormatVersionInput = (
    <TextInput
      type="text"
      placeholder="Example, 1"
      id="source-model-format-version"
      name="source-model-format-version"
      value={sourceModelFormatVersion}
      onChange={(_e, value) => setData('sourceModelFormatVersion', value)}
    />
  );

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
    />
  );

  return (
    <>
      <FormSection
        title="Version details"
        description={
          isFirstVersion
            ? 'Provide model details that apply to every version of this model.'
            : 'Configure details for the version of this model.'
        }
      >
        <FormGroup label="Version name" isRequired fieldId="version-name">
          <FormFieldset component={versionNameInput} field="Version Name" />
          <FormHelperText>
            {latestVersion && (
              <HelperText>
                <HelperTextItem>Current version is {latestVersion.name}</HelperTextItem>
              </HelperText>
            )}
            {!isVersionNameValid && (
              <HelperText>
                <HelperTextItem variant="error">
                  Cannot exceed {MR_CHARACTER_LIMIT} characters
                </HelperTextItem>
              </HelperText>
            )}
          </FormHelperText>
        </FormGroup>
        <FormGroup label="Version description" fieldId="version-description">
          <FormFieldset component={versionDescriptionInput} field="Version Description" />
        </FormGroup>
        <FormGroup label="Source model format" fieldId="source-model-format">
          <FormFieldset component={sourceModelFormatInput} field="Source Model Format" />
        </FormGroup>
        <FormGroup label="Source model format version" fieldId="source-model-format-version">
          <FormFieldset
            component={sourceModelFormatVersionInput}
            field="Source Model Format Version"
          />
        </FormGroup>
      </FormSection>
      <FormSection
        title="Model location"
        description="Specify the model location by providing either the object storage details or the URI."
      >
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
        {/* {modelLocationType === ModelLocationType.ObjectStorage && (
            <SplitItem>
              <Button
                data-testid="object-storage-autofill-button"
                variant="link"
                isInline
                icon={<OptimizeIcon />}
                onClick={() => setAutofillModalOpen(true)}
              >
                Autofill from data connection
              </Button>
            </SplitItem>
          )} */}
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
      </FormSection>
      {/* {isAutofillModalOpen ? (
        <ConnectionModal
          onClose={() => setAutofillModalOpen(false)}
          onSubmit={(connection) => {
            fillObjectStorageByConnection(connection);
            setAutofillModalOpen(false);
          }}
        />
      ) : null} */}
    </>
  );
};

export default RegistrationCommonFormSections;

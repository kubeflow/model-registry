import React from 'react';
import {
  FormGroup,
  TextInput,
  TextArea,
  Radio,
  Split,
  SplitItem,
  HelperText,
  HelperTextItem,
  FormHelperText,
  TextInputGroup,
  TextInputGroupMain,
} from '@patternfly/react-core';
import spacing from '@patternfly/react-styles/css/utilities/Spacing/spacing';
import { UpdateObjectAtPropAndValue } from '~/shared/types';
// import { DataConnection, UpdateObjectAtPropAndValue } from '~/pages/projects/types';
// import { convertAWSSecretData } from '~/pages/projects/screens/detail/data-connections/utils';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import FormSection from '~/shared/components/pf-overrides/FormSection';
import { ModelVersion } from '~/app/types';
import { isMUITheme } from '~/shared/utilities/const';
import { ModelLocationType, RegistrationCommonFormData } from './useRegisterModelData';
// import { ConnectionModal } from './ConnectionModal';

type RegistrationCommonFormSectionsProps = {
  formData: RegistrationCommonFormData;
  setData: UpdateObjectAtPropAndValue<RegistrationCommonFormData>;
  isFirstVersion: boolean;
  latestVersion?: ModelVersion;
};

const RegistrationCommonFormSections: React.FC<RegistrationCommonFormSectionsProps> = ({
  formData,
  setData,
  isFirstVersion,
  latestVersion,
}) => {
  // TODO: [Data connections] Check wether we should use data connections
  // const [isAutofillModalOpen, setAutofillModalOpen] = React.useState(false);

  // const connectionDataMap: Record<string, keyof RegistrationCommonFormData> = {
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
            ? 'Configure details for the first version of this model.'
            : 'Configure details for the version of this model.'
        }
      >
        <FormGroup label="Version name" isRequired fieldId="version-name">
          {isMUITheme() ? (
            <FormFieldset component={versionNameInput} field="Version Name" />
          ) : (
            versionNameInput
          )}
        </FormGroup>
        {latestVersion && (
          <FormHelperText>
            <HelperText>
              <HelperTextItem>Current version is {latestVersion.name}</HelperTextItem>
            </HelperText>
          </FormHelperText>
        )}
        <FormGroup
          className="version-description"
          label="Version description"
          fieldId="version-description"
        >
          {isMUITheme() ? (
            <FormFieldset component={versionDescriptionInput} field="Version Description" />
          ) : (
            versionDescriptionInput
          )}
        </FormGroup>
        <FormGroup label="Source model format" fieldId="source-model-format">
          {isMUITheme() ? (
            <FormFieldset component={sourceModelFormatInput} field="Source Model Format" />
          ) : (
            sourceModelFormatInput
          )}
        </FormGroup>
        <FormGroup label="Source model format version" fieldId="source-model-format-version">
          {isMUITheme() ? (
            <FormFieldset
              component={sourceModelFormatVersionInput}
              field="Source Model Format Version"
            />
          ) : (
            sourceModelFormatVersionInput
          )}
        </FormGroup>
      </FormSection>
      <FormSection
        title="Model location"
        description="Specify the model location by providing either the object storage details or the URI."
      >
        <Split>
          <SplitItem isFilled>
            <Radio
              isChecked={modelLocationType === ModelLocationType.ObjectStorage}
              name="location-type-object-storage"
              onChange={() => {
                setData('modelLocationType', ModelLocationType.ObjectStorage);
              }}
              label="Object storage"
              id="location-type-object-storage"
            />
          </SplitItem>
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
        </Split>
        {modelLocationType === ModelLocationType.ObjectStorage && (
          <>
            <FormGroup
              className={spacing.mlLg}
              label="Endpoint"
              isRequired
              fieldId="location-endpoint"
            >
              {isMUITheme() ? (
                <FormFieldset component={endpointInput} field="Endpoint" />
              ) : (
                endpointInput
              )}
            </FormGroup>
            <FormGroup className={spacing.mlLg} label="Bucket" isRequired fieldId="location-bucket">
              {isMUITheme() ? <FormFieldset component={bucketInput} field="Bucket" /> : bucketInput}
            </FormGroup>
            <FormGroup className={spacing.mlLg} label="Region" fieldId="location-region">
              {isMUITheme() ? <FormFieldset component={regionInput} field="Region" /> : regionInput}
            </FormGroup>
            <FormGroup className={spacing.mlLg} label="Path" isRequired fieldId="location-path">
              {isMUITheme() ? <FormFieldset component={pathInput} field="Path" /> : pathInput}
            </FormGroup>
            <FormHelperText className="path-helper-text">
              <HelperText>
                <HelperTextItem>
                  Enter a path to a model or folder. This path cannot point to a root folder.
                </HelperTextItem>
              </HelperText>
            </FormHelperText>
          </>
        )}
        <Split>
          <SplitItem isFilled>
            <Radio
              isChecked={modelLocationType === ModelLocationType.URI}
              name="location-type-uri"
              onChange={() => {
                setData('modelLocationType', ModelLocationType.URI);
              }}
              label="URI"
              id="location-type-uri"
            />
          </SplitItem>
        </Split>
        {modelLocationType === ModelLocationType.URI && (
          <>
            <FormGroup className={spacing.mlLg} label="URI" isRequired fieldId="location-uri">
              {isMUITheme() ? <FormFieldset component={uriInput} field="URI" /> : uriInput}
            </FormGroup>
          </>
        )}
      </FormSection>
      {/* <ConnectionModal
        isOpen={isAutofillModalOpen}
        onClose={() => setAutofillModalOpen(false)}
        onSubmit={(connection) => {
          fillObjectStorageByConnection(connection);
          setAutofillModalOpen(false);
        }}
      /> */}
    </>
  );
};

export default RegistrationCommonFormSections;

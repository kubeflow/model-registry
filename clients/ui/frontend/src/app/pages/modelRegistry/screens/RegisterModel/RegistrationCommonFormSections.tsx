import React from 'react';
import {
  FormGroup,
  TextInput,
  TextArea,
  Radio,
  Split,
  SplitItem,
  InputGroupText,
  InputGroupItem,
  HelperText,
  HelperTextItem,
  FormHelperText,
} from '@patternfly/react-core';
import spacing from '@patternfly/react-styles/css/utilities/Spacing/spacing';
import { UpdateObjectAtPropAndValue } from '~/types';
// import { DataConnection, UpdateObjectAtPropAndValue } from '~/pages/projects/types';
// import { convertAWSSecretData } from '~/pages/projects/screens/detail/data-connections/utils';
import FormSection from '~/app/components/pf-overrides/FormSection';
import { ModelVersion } from '~/app/types';
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
  // const [isAutofillModalOpen, setAutofillModalOpen] = React.useState(false); TODO: Check wether we should use data connections

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
          <TextInput
            isRequired
            type="text"
            id="version-name"
            name="version-name"
            value={versionName}
            onChange={(_e, value) => setData('versionName', value)}
          />
          {latestVersion && (
            <FormHelperText>
              <HelperText>
                <HelperTextItem>Current version is {latestVersion.name}</HelperTextItem>
              </HelperText>
            </FormHelperText>
          )}
        </FormGroup>
        <FormGroup label="Version description" fieldId="version-description">
          <TextArea
            type="text"
            id="version-description"
            name="version-description"
            value={versionDescription}
            onChange={(_e, value) => setData('versionDescription', value)}
          />
        </FormGroup>
        <FormGroup label="Source model format" fieldId="source-model-format">
          <TextInput
            type="text"
            placeholder="Example, tensorflow"
            id="source-model-format"
            name="source-model-format"
            value={sourceModelFormat}
            onChange={(_e, value) => setData('sourceModelFormat', value)}
          />
        </FormGroup>
        <FormGroup label="Source model format version" fieldId="source-model-format-version">
          <TextInput
            type="text"
            placeholder="Example, 1"
            id="source-model-format-version"
            name="source-model-format-version"
            value={sourceModelFormatVersion}
            onChange={(_e, value) => setData('sourceModelFormatVersion', value)}
          />
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
              <TextInput
                isRequired
                type="text"
                id="location-endpoint"
                name="location-endpoint"
                value={modelLocationEndpoint}
                onChange={(_e, value) => setData('modelLocationEndpoint', value)}
              />
            </FormGroup>
            <FormGroup className={spacing.mlLg} label="Bucket" isRequired fieldId="location-bucket">
              <TextInput
                isRequired
                type="text"
                id="location-bucket"
                name="location-bucket"
                value={modelLocationBucket}
                onChange={(_e, value) => setData('modelLocationBucket', value)}
              />
            </FormGroup>
            <FormGroup className={spacing.mlLg} label="Region" fieldId="location-region">
              <TextInput
                type="text"
                id="location-region"
                name="location-region"
                value={modelLocationRegion}
                onChange={(_e, value) => setData('modelLocationRegion', value)}
              />
            </FormGroup>
            <FormGroup className={spacing.mlLg} label="Path" isRequired fieldId="location-path">
              <Split hasGutter>
                <SplitItem>
                  <InputGroupText isPlain>/</InputGroupText>
                </SplitItem>
                <SplitItem isFilled>
                  <InputGroupItem>
                    <TextInput
                      isRequired
                      type="text"
                      id="location-path"
                      name="location-path"
                      value={modelLocationPath}
                      onChange={(_e, value) => setData('modelLocationPath', value)}
                    />
                  </InputGroupItem>
                </SplitItem>
              </Split>
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
          body={
            modelLocationType === ModelLocationType.URI && (
              <FormGroup label="URI" isRequired fieldId="location-uri">
                <TextInput
                  isRequired
                  type="text"
                  id="location-uri"
                  name="location-uri"
                  value={modelLocationURI}
                  onChange={(_e, value) => setData('modelLocationURI', value)}
                />
              </FormGroup>
            )
          }
        />
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

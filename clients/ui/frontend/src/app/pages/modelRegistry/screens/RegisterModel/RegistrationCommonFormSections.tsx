import React from 'react';
import {
  FormGroup,
  TextInput,
  TextArea,
  HelperText,
  HelperTextItem,
  FormHelperText,
  ToggleGroup,
  ToggleGroupItem,
} from '@patternfly/react-core';
import spacing from '@patternfly/react-styles/css/utilities/Spacing/spacing';
import { UpdateObjectAtPropAndValue } from 'mod-arch-shared';
// import { DataConnection, UpdateObjectAtPropAndValue } from '~/pages/projects/types';
// import { convertAWSSecretData } from '~/pages/projects/screens/detail/data-connections/utils';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import { ModelVersion } from '~/app/types';
import FormSection from '~/app/pages/modelRegistry/components/pf-overrides/FormSection';
import { useTempDevFeatureAvailable, TempDevFeature } from '~/app/hooks/useTempDevFeatureAvailable';
import { RegistrationMode } from '~/app/pages/modelRegistry/screens/const';
import { RegistrationCommonFormData } from './useRegisterModelData';
import RegistrationModelLocationFields from './RegistrationModelLocationFields';
import RegisterAndStoreFields from './RegisterAndStoreFields';
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
  const isVersionNameValid = isNameValid(formData.versionName);
  const isRegistryStorageFeatureAvailable = useTempDevFeatureAvailable(
    TempDevFeature.RegistryStorage,
  );
  const registrationMode = formData.registrationMode || RegistrationMode.Register;

  const { versionName, versionDescription, sourceModelFormat, sourceModelFormatVersion } = formData;

  const handleModeChange = (mode: RegistrationMode) => {
    setData('registrationMode', mode);

    if (mode === RegistrationMode.Register) {
      setData('namespace', '');
    }
  };

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
      autoResize
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
          {latestVersion && (
            <FormHelperText>
              <HelperText>
                <HelperTextItem>Current version is {latestVersion.name}</HelperTextItem>
              </HelperText>
              {!isVersionNameValid && (
                <HelperText>
                  <HelperTextItem variant="error">
                    Cannot exceed {MR_CHARACTER_LIMIT} characters
                  </HelperTextItem>
                </HelperText>
              )}
            </FormHelperText>
          )}
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
        title={
          // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
          isRegistryStorageFeatureAvailable ? 'Model location and storage' : 'Model location'
        }
        description={
          // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
          isRegistryStorageFeatureAvailable ? (
            <>
              Choose <strong>Register</strong> to use the model&apos;s original storage location for
              artifact storage, or <strong>Register and store</strong> to specify a different
              artifact storage location.
            </>
          ) : (
            'Specify the model location by providing either the object storage details or the URI.'
          )
        }
      >
        {/* eslint-disable-next-line @typescript-eslint/no-unnecessary-condition */}
        {isRegistryStorageFeatureAvailable && (
          <ToggleGroup
            aria-label="Registration mode"
            className={spacing.myMd}
            data-testid="registration-mode-toggle-group"
          >
            <ToggleGroupItem
              text="Register"
              isSelected={registrationMode === RegistrationMode.Register}
              data-testid="registration-mode-register"
              onChange={() => handleModeChange(RegistrationMode.Register)}
            />
            <ToggleGroupItem
              text="Register and store"
              isSelected={registrationMode === RegistrationMode.RegisterAndStore}
              onChange={() => handleModeChange(RegistrationMode.RegisterAndStore)}
              data-testid="registration-mode-register-and-store"
            />
          </ToggleGroup>
        )}
        {registrationMode === RegistrationMode.Register ? (
          <RegistrationModelLocationFields
            formData={formData}
            setData={setData}
            isCatalogModel={isCatalogModel}
          />
        ) : (
          <RegisterAndStoreFields
            formData={formData}
            setData={setData}
            isCatalogModel={isCatalogModel}
          />
        )}
      </FormSection>
    </>
  );
};

export default RegistrationCommonFormSections;

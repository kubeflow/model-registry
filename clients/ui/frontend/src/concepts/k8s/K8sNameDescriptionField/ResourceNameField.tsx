import * as React from 'react';
import {
  FormGroup,
  HelperText,
  HelperTextItem,
  TextInput,
  ValidatedOptions,
} from '@patternfly/react-core';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import { K8sNameDescriptionFieldData, K8sNameDescriptionFieldUpdateFunction } from './types';

type ResourceNameFieldProps = {
  allowEdit: boolean;
  dataTestId: string;
  k8sName: K8sNameDescriptionFieldData['k8sName'];
  onDataChange?: K8sNameDescriptionFieldUpdateFunction;
};

/** Sub-resource; not for public consumption */
const ResourceNameField: React.FC<ResourceNameFieldProps> = ({
  allowEdit,
  dataTestId,
  k8sName,
  onDataChange,
}) => {
  if (k8sName.state.immutable) {
    return (
      <FormGroup label="Resource name" fieldId={`${dataTestId}-resource-name`}>
        <FormFieldset component={<div>{k8sName.value}</div>} field="Resource name" />
      </FormGroup>
    );
  }

  if (!allowEdit) {
    return null;
  }

  let validated: ValidatedOptions = ValidatedOptions.default;
  if (k8sName.state.invalidLength || k8sName.state.invalidCharacters) {
    validated = ValidatedOptions.error;
  } else if (k8sName.value.length > 0) {
    validated = ValidatedOptions.success;
  }

  const textInput = (
    <TextInput
      aria-readonly={!onDataChange}
      id={`${dataTestId}-resourceName`}
      data-testid={`${dataTestId}-resourceName`}
      name={`${dataTestId}-resourceName`}
      type="text"
      isRequired
      value={k8sName.value}
      onChange={(_event, value) => onDataChange?.('k8sName', value)}
      validated={validated}
    />
  );

  return (
    <>
      <FormGroup
        label="Resource name"
        className="resource-name"
        isRequired
        fieldId={`${dataTestId}-resource-name`}
      >
        <FormFieldset component={textInput} field="Resource name" />
      </FormGroup>
      <HelperText>
        {k8sName.state.invalidLength && (
          <HelperTextItem variant="error">
            Cannot exceed {k8sName.state.maxLength} characters
          </HelperTextItem>
        )}
        {k8sName.state.invalidCharacters && (
          <HelperTextItem variant="error">
            Must start and end with a lowercase letter or number. Valid characters include lowercase
            letters, numbers, and hyphens (-).
          </HelperTextItem>
        )}
        {!k8sName.state.invalidLength && !k8sName.state.invalidCharacters && (
          <HelperTextItem>
            The resource name is used to identify your resource, and is generated based on the name
            you enter. The resource name cannot be edited after creation.
          </HelperTextItem>
        )}
      </HelperText>
    </>
  );
};

export default ResourceNameField;

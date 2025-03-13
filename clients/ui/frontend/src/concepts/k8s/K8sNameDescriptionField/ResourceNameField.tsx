import * as React from 'react';
import { FormGroup, HelperText, TextInput } from '@patternfly/react-core';
import ResourceNameDefinitionTooltip from '~/concepts/k8s/ResourceNameDefinitionTootip';

type ResourceNameFieldProps = {
  allowEdit: boolean;
  dataTestId: string;
  //   k8sName: K8sNameDescriptionFieldData['k8sName'];
  //   onDataChange?: K8sNameDescriptionFieldUpdateFunction;
};

/** Sub-resource; not for public consumption */
const ResourceNameField: React.FC<ResourceNameFieldProps> = ({
  allowEdit,
  dataTestId,
  //   k8sName,
  //   onDataChange,
}) => {
  const formGroupProps: React.ComponentProps<typeof FormGroup> = {
    label: 'Resource name',
    labelHelp: <ResourceNameDefinitionTooltip />,
    fieldId: `${dataTestId}-resourceName`,
  };

  // TODO: Implement this once we have the endpoint.

  //   if (k8sName.state.immutable) {
  //     return <FormGroup {...formGroupProps}>{k8sName.value}</FormGroup>;
  //   }

  if (!allowEdit) {
    return null;
  }

  // TODO: Implement this once we have the endpoint.

  //   let validated: ValidatedOptions = ValidatedOptions.default;
  //   if (k8sName.state.invalidLength || k8sName.state.invalidCharacters) {
  //     validated = ValidatedOptions.error;
  //   } else if (k8sName.value.length > 0) {
  //     validated = ValidatedOptions.success;
  //   }

  //   const usePrefix = k8sName.state.staticPrefix && !!k8sName.state.safePrefix;
  const textInput = (
    <TextInput
      id={`${dataTestId}-resourceName`}
      data-testid={`${dataTestId}-resourceName`}
      name={`${dataTestId}-resourceName`}
      isRequired
      //   value={
      //     usePrefix && k8sName.state.safePrefix
      //       ? k8sName.value.replace(new RegExp(`^${k8sName.state.safePrefix}`), '')
      //       : k8sName.value
      //   }
      //   onChange={(event, value) =>
      //     onDataChange?.(
      //       'k8sName',
      //       usePrefix && k8sName.state.safePrefix ? `${k8sName.state.safePrefix}${value}` : value,
      //     )
      //   }
      //   validated={validated}
    />
  );
  return (
    <FormGroup {...formGroupProps} isRequired>
      {/* {usePrefix ? (
        <InputGroup>
          <InputGroupText>{k8sName.state.safePrefix}</InputGroupText>
          <InputGroupItem isFill>{textInput}</InputGroupItem>
        </InputGroup>
      ) : ( */}
      {textInput}
      {/* )} */}
      <HelperText>
        {/* <HelperTextItemMaxLength k8sName={k8sName} />
        <HelperTextItemValidCharacters k8sName={k8sName} /> */}
      </HelperText>
    </FormGroup>
  );
};

export default ResourceNameField;

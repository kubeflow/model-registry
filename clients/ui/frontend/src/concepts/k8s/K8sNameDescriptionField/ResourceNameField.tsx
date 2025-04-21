import * as React from 'react';
import { FormGroup, HelperText, HelperTextItem, TextInput } from '@patternfly/react-core';
import { useThemeContext } from 'mod-arch-shared';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';

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
  // TODO: Implement this once we have the endpoint.

  //   if (k8sName.state.immutable) {
  //     return <FormGroup {...formGroupProps}>{k8sName.value}</FormGroup>;
  //   }

  const { isMUITheme } = useThemeContext();
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
      type="text"
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

  const resourceNameFormGroup = (
    <>
      <FormGroup
        label="Resource name"
        className="resource-name"
        isRequired
        fieldId={`${dataTestId}-resource-name`}
      >
        <FormFieldset component={textInput} field="Host" />
      </FormGroup>
      <HelperText>
        <HelperTextItem>
          The resource name is used to identify your resource, and is generated based on the name
          you enter. The resource name cannot be edited after creation.
        </HelperTextItem>
        {/* <HelperTextItemMaxLength k8sName={k8sName} />
         <HelperTextItemValidCharacters k8sName={k8sName} /> */}
      </HelperText>
    </>
  );

  // TODO: Implement this once we have the endpoint.
  // return (
  //   <FormGroup {...formGroupProps} isRequired>
  //     {usePrefix ? (
  //       <InputGroup>
  //         <InputGroupText>{k8sName.state.safePrefix}</InputGroupText>
  //         <InputGroupItem isFill>{textInput}</InputGroupItem>
  //       </InputGroup>
  //     ) : (
  //       { textInput }
  //     )}
  //     <HelperText>
  //       <HelperTextItemMaxLength k8sName={k8sName} />
  //       <HelperTextItemValidCharacters k8sName={k8sName} />
  //     </HelperText>
  //   </FormGroup>
  // );

  return isMUITheme ? resourceNameFormGroup : textInput;
};

export default ResourceNameField;

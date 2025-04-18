import React from 'react';
import { FormGroup, TextInput } from '@patternfly/react-core';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import { useThemeContext } from '~/app/ThemeContext';

type PrefilledModelRegistryFieldProps = {
  mrName?: string;
};

const PrefilledModelRegistryField: React.FC<PrefilledModelRegistryFieldProps> = ({ mrName }) => {
  const { isMUITheme } = useThemeContext();

  const mrNameInput = (
    <TextInput isDisabled isRequired type="text" id="mr-name" name="mr-name" value={mrName} />
  );

  return (
    <FormGroup className="form-group-disabled" label="Model registry" isRequired fieldId="mr-name">
      {isMUITheme ? <FormFieldset component={mrNameInput} field="Model Registry" /> : mrNameInput}
    </FormGroup>
  );
};

export default PrefilledModelRegistryField;

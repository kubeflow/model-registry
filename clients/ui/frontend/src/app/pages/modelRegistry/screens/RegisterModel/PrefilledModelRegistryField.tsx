import React from 'react';
import { FormGroup, TextInput, Stack, StackItem } from '@patternfly/react-core';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';

type PrefilledModelRegistryFieldProps = {
  mrName?: string;
};

const PrefilledModelRegistryField: React.FC<PrefilledModelRegistryFieldProps> = ({ mrName }) => {
  const mrNameInput = (
    <TextInput isDisabled isRequired type="text" id="mr-name" name="mr-name" value={mrName} />
  );

  return (
    <Stack hasGutter>
      <StackItem>
        <FormGroup
          className="form-group-disabled"
          label="Model registry"
          isRequired
          fieldId="mr-name"
        >
          <FormFieldset component={mrNameInput} field="Model Registry" />
        </FormGroup>
      </StackItem>
    </Stack>
  );
};

export default PrefilledModelRegistryField;

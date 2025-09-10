import React from 'react';
import { FormGroup, TextInput, Alert, Stack, StackItem } from '@patternfly/react-core';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';

type PrefilledModelRegistryFieldProps = {
  mrName?: string;
  isTriggeredFromCatalog?: boolean;
};

const PrefilledModelRegistryField: React.FC<PrefilledModelRegistryFieldProps> = ({
  mrName,
  isTriggeredFromCatalog = false,
}) => {
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
      {isTriggeredFromCatalog && (
        <StackItem>
          <Alert
            variant="info"
            isInline
            isPlain
            title="Additional model metadata, such as its model card details, labels, provider, and license, will be available to view and edit after registration is complete."
          />
        </StackItem>
      )}
    </Stack>
  );
};

export default PrefilledModelRegistryField;

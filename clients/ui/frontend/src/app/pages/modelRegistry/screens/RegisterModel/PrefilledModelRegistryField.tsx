import React from 'react';
import { FormGroup, TextInput } from '@patternfly/react-core';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';

type PrefilledModelRegistryFieldProps = {
  mrName?: string;
};

const PrefilledModelRegistryField: React.FC<PrefilledModelRegistryFieldProps> = ({ mrName }) => (
  <FormGroup className="form-group-disabled" label="Model registry" isRequired fieldId="mr-name">
    <FormFieldset
      component={
        <TextInput isDisabled isRequired type="text" id="mr-name" name="mr-name" value={mrName} />
      }
      field="Model Registry"
    />
  </FormGroup>
);

export default PrefilledModelRegistryField;

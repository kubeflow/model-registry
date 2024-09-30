import React from 'react';
import { FormGroup, TextInput } from '@patternfly/react-core';

type PrefilledModelRegistryFieldProps = {
  mrName?: string;
};

const PrefilledModelRegistryField: React.FC<PrefilledModelRegistryFieldProps> = ({ mrName }) => (
  <FormGroup label="Model registry" isRequired fieldId="mr-name">
    <TextInput isDisabled isRequired type="text" id="mr-name" name="mr-name" value={mrName} />
  </FormGroup>
);

export default PrefilledModelRegistryField;

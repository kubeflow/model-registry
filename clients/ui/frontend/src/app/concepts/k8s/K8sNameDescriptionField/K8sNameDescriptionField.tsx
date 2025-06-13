import * as React from 'react';
import { FormGroup, TextField, FormHelperText } from '@mui/material';
import { K8sNameDescriptionFieldType } from './useK8sNameDescriptionField';

type K8sNameDescriptionFieldProps = {
  data: K8sNameDescriptionFieldType;
  onDataChange: (data: K8sNameDescriptionFieldType) => void;
  k8sNameIsEditable?: boolean;
  dataTestId?: string;
};

const K8sNameDescriptionField: React.FC<K8sNameDescriptionFieldProps> = ({
  data,
  onDataChange,
  k8sNameIsEditable,
  dataTestId,
}) => (
  <>
    <FormGroup>
      <TextField
        label="Name"
        value={data.name}
        onChange={(e) => onDataChange({ ...data, name: e.target.value })}
        required
        data-testid={`${dataTestId}-name-field`}
      />
    </FormGroup>
    {k8sNameIsEditable && (
      <FormGroup>
        <TextField
          label="Kubernetes Name"
          value={data.k8sName.value}
          onChange={(e) => onDataChange({ ...data, k8sName: { ...data.k8sName, value: e.target.value } })}
          required
          error={!!data.k8sName.error}
          helperText={data.k8sName.error}
          data-testid={`${dataTestId}-k8s-name-field`}
        />
      </FormGroup>
    )}
    <FormGroup>
      <TextField
        label="Description"
        multiline
        rows={4}
        value={data.description}
        onChange={(e) => onDataChange({ ...data, description: e.target.value })}
        data-testid={`${dataTestId}-description-field`}
      />
    </FormGroup>
  </>
);

export default K8sNameDescriptionField; 
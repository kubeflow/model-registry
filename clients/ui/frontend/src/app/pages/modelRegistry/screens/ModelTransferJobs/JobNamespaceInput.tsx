import React from 'react';
import {
  FormGroup,
  HelperText,
  HelperTextItem,
  TextInput,
} from '@patternfly/react-core';

type JobNamespaceInputProps = {
  value: string;
  onChange: (namespace: string) => void;
};

const DEBOUNCE_MS = 2000;

const JobNamespaceInput: React.FC<JobNamespaceInputProps> = ({ value, onChange }) => {
  const [textInputValue, setTextInputValue] = React.useState(value);
  const debounceTimerRef = React.useRef<ReturnType<typeof setTimeout>>();

  const handleChange = (_event: React.FormEvent, newValue: string) => {
    setTextInputValue(newValue);
    if (debounceTimerRef.current) {
      clearTimeout(debounceTimerRef.current);
    }
    if (newValue) {
      debounceTimerRef.current = setTimeout(() => {
        onChange(newValue);
      }, DEBOUNCE_MS);
    }
  };

  const handleBlur = () => {
    if (debounceTimerRef.current) {
      clearTimeout(debounceTimerRef.current);
    }
    if (textInputValue && textInputValue !== value) {
      onChange(textInputValue);
    }
  };

  React.useEffect(
    () => () => {
      if (debounceTimerRef.current) {
        clearTimeout(debounceTimerRef.current);
      }
    },
    [],
  );

  return (
    <FormGroup label="Transfer job namespace" fieldId="job-namespace-input">
      <TextInput
        id="job-namespace-input"
        data-testid="job-namespace-input"
        value={textInputValue}
        onChange={handleChange}
        onBlur={handleBlur}
        placeholder="Enter namespace"
      />
      <HelperText>
        <HelperTextItem variant="indeterminate">
          You do not have permission to list transfer jobs across all namespaces. Enter a namespace
          to view transfer jobs in that namespace.
        </HelperTextItem>
      </HelperText>
    </FormGroup>
  );
};

export default JobNamespaceInput;

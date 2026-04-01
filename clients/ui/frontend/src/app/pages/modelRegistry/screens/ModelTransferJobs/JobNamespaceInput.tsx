import React from 'react';
import { Bullseye, Button, FlexItem, Popover, TextInput } from '@patternfly/react-core';
import { InfoCircleIcon } from '@patternfly/react-icons';

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
    <>
      <FlexItem>
        <Bullseye>Transfer job namespace</Bullseye>
      </FlexItem>
      <FlexItem>
        <TextInput
          id="job-namespace-input"
          data-testid="job-namespace-input"
          value={textInputValue}
          onChange={handleChange}
          onBlur={handleBlur}
          placeholder="Enter namespace"
          aria-label="Transfer job namespace"
          style={{ width: '200px' }}
        />
      </FlexItem>
      <FlexItem>
        <Popover
          aria-label="Transfer job namespace info"
          bodyContent="You do not have permission to list transfer jobs across all namespaces. Enter a namespace to view transfer jobs in that namespace."
        >
          <Button
            variant="plain"
            aria-label="More info about transfer job namespace"
            data-testid="job-namespace-info"
          >
            <InfoCircleIcon />
          </Button>
        </Popover>
      </FlexItem>
    </>
  );
};

export default JobNamespaceInput;

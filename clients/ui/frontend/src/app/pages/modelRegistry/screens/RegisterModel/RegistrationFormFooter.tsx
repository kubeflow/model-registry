import React from 'react';
import {
  PageSection,
  Stack,
  StackItem,
  Alert,
  AlertActionCloseButton,
  ActionGroup,
  Button,
} from '@patternfly/react-core';

type RegistrationFormFooterProps = {
  submitLabel: string;
  submitError?: Error;
  setSubmitError: (e?: Error) => void;
  isSubmitDisabled: boolean;
  isSubmitting: boolean;
  onSubmit: () => void;
  onCancel: () => void;
};

const RegistrationFormFooter: React.FC<RegistrationFormFooterProps> = ({
  submitLabel,
  submitError,
  setSubmitError,
  isSubmitDisabled,
  isSubmitting,
  onSubmit,
  onCancel,
}) => (
  <PageSection hasBodyWrapper={false} stickyOnBreakpoint={{ default: 'bottom' }}>
    <Stack hasGutter>
      {submitError && (
        <StackItem>
          <Alert
            isInline
            variant="danger"
            title={submitError.name}
            actionClose={<AlertActionCloseButton onClose={() => setSubmitError(undefined)} />}
          >
            {submitError.message}
          </Alert>
        </StackItem>
      )}
      <StackItem>
        <ActionGroup>
          <Button
            isDisabled={isSubmitDisabled}
            variant="primary"
            id="create-button"
            data-testid="create-button"
            isLoading={isSubmitting}
            onClick={onSubmit}
          >
            {submitLabel}
          </Button>
          <Button isDisabled={isSubmitting} variant="link" id="cancel-button" onClick={onCancel}>
            Cancel
          </Button>
        </ActionGroup>
      </StackItem>
    </Stack>
  </PageSection>
);

export default RegistrationFormFooter;

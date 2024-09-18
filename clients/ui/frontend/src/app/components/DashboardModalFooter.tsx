import * as React from 'react';
import {
  ActionList,
  ActionListItem,
  Alert,
  Button,
  ButtonProps,
  Stack,
  StackItem,
} from '@patternfly/react-core';

type DashboardModalFooterProps = {
  submitLabel: string;
  submitButtonVariant?: ButtonProps['variant'];
  onSubmit: () => void;
  onCancel: () => void;
  isSubmitDisabled: boolean;
  isSubmitLoading?: boolean;
  isCancelDisabled?: boolean;
  alertTitle: string;
  error?: Error;
};

const DashboardModalFooter: React.FC<DashboardModalFooterProps> = ({
  submitLabel,
  submitButtonVariant = 'primary',
  onSubmit,
  onCancel,
  isSubmitDisabled,
  isSubmitLoading,
  isCancelDisabled,
  error,
  alertTitle,
}) => (
  // make sure alert uses the full width
  <Stack hasGutter style={{ flex: 'auto' }}>
    {error && (
      <StackItem>
        <Alert data-testid="error-message-alert" isInline variant="danger" title={alertTitle}>
          {error.message}
        </Alert>
      </StackItem>
    )}
    <StackItem>
      <ActionList>
        <ActionListItem>
          <Button
            key="submit"
            variant={submitButtonVariant}
            isDisabled={isSubmitDisabled}
            onClick={onSubmit}
            isLoading={isSubmitLoading}
            data-testid="modal-submit-button"
          >
            {submitLabel}
          </Button>
        </ActionListItem>
        <ActionListItem>
          <Button
            key="cancel"
            variant="link"
            isDisabled={isCancelDisabled}
            onClick={onCancel}
            data-testid="modal-cancel-button"
          >
            Cancel
          </Button>
        </ActionListItem>
      </ActionList>
    </StackItem>
  </Stack>
);

export default DashboardModalFooter;

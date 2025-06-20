import * as React from 'react';
import {
  ActionList,
  ActionListItem,
  ActionListGroup,
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
  isSubmitDisabled?: boolean;
  isSubmitLoading?: boolean;
  isCancelDisabled?: boolean;
  alertTitle?: string;
  error?: Error | React.ReactNode;
  alertLinks?: React.ReactNode;
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
  alertLinks,
}) => (
  // make sure alert uses the full width
  <Stack hasGutter style={{ flex: 'auto' }}>
    {error && (
      <StackItem>
        <Alert
          data-testid="error-message-alert"
          isInline
          variant="danger"
          title={alertTitle}
          actionLinks={alertLinks}
        >
          {error instanceof Error ? error.message : error}
        </Alert>
      </StackItem>
    )}
    <StackItem>
      <ActionList>
        <ActionListGroup>
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
        </ActionListGroup>
      </ActionList>
    </StackItem>
  </Stack>
);

export default DashboardModalFooter;

import * as React from 'react';
import { Alert, Button, ButtonProps } from '@patternfly/react-core';

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
  <>
    {error && (
      <Alert data-testid="error-message-alert" isInline variant="danger" title={alertTitle}>
        {error.message}
      </Alert>
    )}
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
    <Button
      key="cancel"
      variant="link"
      isDisabled={isCancelDisabled}
      onClick={onCancel}
      data-testid="modal-cancel-button"
    >
      Cancel
    </Button>
  </>
);

export default DashboardModalFooter;

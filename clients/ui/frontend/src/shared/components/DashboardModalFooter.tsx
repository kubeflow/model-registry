import * as React from 'react';
import { Button, ButtonProps, ModalFooter } from '@patternfly/react-core';

type DashboardModalFooterProps = {
  submitLabel: string;
  submitButtonVariant?: ButtonProps['variant'];
  onSubmit: () => void;
  onCancel: () => void;
  isSubmitDisabled: boolean;
  isSubmitLoading?: boolean;
  isCancelDisabled?: boolean;
};

const DashboardModalFooter: React.FC<DashboardModalFooterProps> = ({
  submitLabel,
  submitButtonVariant = 'primary',
  onSubmit,
  onCancel,
  isSubmitDisabled,
  isSubmitLoading,
  isCancelDisabled,
}) => (
  <ModalFooter>
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
  </ModalFooter>
);

export default DashboardModalFooter;

import * as React from 'react';
import { Modal } from '@patternfly/react-core/deprecated';
import DashboardModalFooter from '~/app/components/DashboardModalFooter';

interface RestoreModelVersionModalProps {
  onCancel: () => void;
  onSubmit: () => void;
  isOpen: boolean;
  modelVersionName: string;
}

export const RestoreModelVersionModal: React.FC<RestoreModelVersionModalProps> = ({
  onCancel,
  onSubmit,
  isOpen,
  modelVersionName,
}) => {
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [error, setError] = React.useState<Error>();

  const onClose = React.useCallback(() => {
    onCancel();
  }, [onCancel]);

  const onConfirm = React.useCallback(async () => {
    setIsSubmitting(true);

    try {
      await onSubmit();
      onClose();
    } catch (e) {
      if (e instanceof Error) {
        setError(e);
      }
    } finally {
      setIsSubmitting(false);
    }
  }, [onSubmit, onClose]);

  return (
    <Modal
      isOpen={isOpen}
      title="Restore version?"
      variant="small"
      onClose={onClose}
      footer={
        <DashboardModalFooter
          onCancel={onClose}
          onSubmit={onConfirm}
          submitLabel="Restore"
          isSubmitLoading={isSubmitting}
          error={error}
          alertTitle="Error"
          isSubmitDisabled={isSubmitting}
        />
      }
      data-testid="restore-model-version-modal"
    >
      <b>{modelVersionName}</b> will be restored and returned to the versions list.
    </Modal>
  );
};

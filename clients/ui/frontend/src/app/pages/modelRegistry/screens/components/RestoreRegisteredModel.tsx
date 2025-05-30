import * as React from 'react';
import { Modal } from '@patternfly/react-core/deprecated';
import { DashboardModalFooter } from 'mod-arch-shared';
import { useNotification } from '~/app/hooks/useNotification';

interface RestoreRegisteredModelModalProps {
  onCancel: () => void;
  onSubmit: () => void;
  registeredModelName: string;
}

export const RestoreRegisteredModelModal: React.FC<RestoreRegisteredModelModalProps> = ({
  onCancel,
  onSubmit,
  registeredModelName,
}) => {
  const notification = useNotification();
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
      notification.success(`${registeredModelName} and all its versions restored.`);
    } catch (e) {
      if (e instanceof Error) {
        setError(e);
      }
    } finally {
      setIsSubmitting(false);
    }
  }, [onSubmit, onClose, notification, registeredModelName]);

  return (
    <Modal
      isOpen
      title="Restore model?"
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
      data-testid="restore-registered-model-modal"
    >
      <b>{registeredModelName}</b> and all of its versions will be restored and returned to the
      registered models list.
    </Modal>
  );
};

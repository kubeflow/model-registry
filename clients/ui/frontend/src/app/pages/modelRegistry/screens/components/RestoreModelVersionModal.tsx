import * as React from 'react';
import { DashboardModalFooter } from 'mod-arch-shared';
import { ModalBody, ModalFooter, ModalHeader, Modal } from '@patternfly/react-core';
import { useNotification } from '~/app/hooks/useNotification';

interface RestoreModelVersionModalProps {
  onCancel: () => void;
  onSubmit: () => void;
  modelVersionName: string;
}

export const RestoreModelVersionModal: React.FC<RestoreModelVersionModalProps> = ({
  onCancel,
  onSubmit,
  modelVersionName,
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
      notification.success(`${modelVersionName} restored.`);
    } catch (e) {
      if (e instanceof Error) {
        setError(e);
      }
    } finally {
      setIsSubmitting(false);
    }
  }, [onSubmit, onClose, notification, modelVersionName]);

  return (
    <Modal isOpen variant="small" onClose={onClose} data-testid="restore-model-version-modal">
      <ModalHeader title="Restore model version?" titleIconVariant="warning" />
      <ModalBody>
        <b>{modelVersionName}</b> will be restored and returned to the versions list.
      </ModalBody>
      <ModalFooter>
        <DashboardModalFooter
          onCancel={onClose}
          onSubmit={onConfirm}
          submitLabel="Restore"
          isSubmitLoading={isSubmitting}
          error={error}
          alertTitle="Error"
          isSubmitDisabled={isSubmitting}
        />
      </ModalFooter>
    </Modal>
  );
};

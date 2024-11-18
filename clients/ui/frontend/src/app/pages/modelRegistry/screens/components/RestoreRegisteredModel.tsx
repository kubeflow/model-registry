import * as React from 'react';
import { Alert, Form, ModalHeader, Modal, ModalBody } from '@patternfly/react-core';
import DashboardModalFooter from '~/shared/components/DashboardModalFooter';
import { useNotification } from '~/app/hooks/useNotification';

interface RestoreRegisteredModelModalProps {
  onCancel: () => void;
  onSubmit: () => void;
  isOpen: boolean;
  registeredModelName: string;
}

export const RestoreRegisteredModelModal: React.FC<RestoreRegisteredModelModalProps> = ({
  onCancel,
  onSubmit,
  isOpen,
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
  }, [notification, registeredModelName, onSubmit, onClose]);

  const description = (
    <>
      <b>{registeredModelName}</b> and all of its versions will be restored and returned to the
      registered models list.
    </>
  );

  return (
    <Modal
      isOpen={isOpen}
      variant="small"
      onClose={onClose}
      data-testid="restore-registered-model-modal"
    >
      <ModalHeader title="Restore model?" />
      <ModalBody>
        <Form>
          {error && (
            <Alert data-testid="error-message-alert" isInline variant="danger" title="Error">
              {error.message}
            </Alert>
          )}
        </Form>
        {description}
      </ModalBody>
      <DashboardModalFooter
        onCancel={onClose}
        onSubmit={onConfirm}
        submitLabel="Restore"
        isSubmitLoading={isSubmitting}
        isSubmitDisabled={isSubmitting}
      />
    </Modal>
  );
};

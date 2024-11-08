import * as React from 'react';
import { Form, Modal, ModalHeader, ModalBody, Alert } from '@patternfly/react-core';
import DashboardModalFooter from '~/shared/components/DashboardModalFooter';
import { useNotification } from '~/app/hooks/useNotification';

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
  const notification = useNotification();

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
  }, [notification, modelVersionName, onSubmit, onClose]);

  const description = (
    <>
      <b>{modelVersionName}</b> will be restored and returned to the versions list.
    </>
  );

  return (
    <Modal
      isOpen={isOpen}
      title="Restore version?"
      variant="small"
      onClose={onClose}
      data-testid="restore-model-version-modal"
    >
      <ModalHeader title="Restore version?" />
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

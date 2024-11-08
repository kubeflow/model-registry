import * as React from 'react';
import {
  Alert,
  Form,
  FormGroup,
  Modal,
  ModalBody,
  ModalHeader,
  TextInput,
} from '@patternfly/react-core';
import DashboardModalFooter from '~/shared/components/DashboardModalFooter';
import { useNotification } from '~/app/hooks/useNotification';

interface ArchiveRegisteredModelModalProps {
  onCancel: () => void;
  onSubmit: () => void;
  isOpen: boolean;
  registeredModelName: string;
}

export const ArchiveRegisteredModelModal: React.FC<ArchiveRegisteredModelModalProps> = ({
  onCancel,
  onSubmit,
  isOpen,
  registeredModelName,
}) => {
  const notification = useNotification();
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [error, setError] = React.useState<Error>();
  const [confirmInputValue, setConfirmInputValue] = React.useState('');
  const isDisabled = confirmInputValue.trim() !== registeredModelName || isSubmitting;

  const onClose = React.useCallback(() => {
    setConfirmInputValue('');
    onCancel();
  }, [onCancel]);

  const onConfirm = React.useCallback(async () => {
    setIsSubmitting(true);

    try {
      await onSubmit();
      onClose();
      notification.success(`${registeredModelName} and all its versions archived.`);
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
      <b>{registeredModelName}</b> and all of its versions will be archived and unavailable for use
      unless it is restored.
      <br />
      <br />
      Type <strong>{registeredModelName}</strong> to confirm archiving:
    </>
  );

  return (
    <Modal
      isOpen={isOpen}
      title="Archive model?"
      variant="small"
      onClose={onClose}
      data-testid="archive-registered-model-modal"
    >
      <ModalHeader title="Archive model?" titleIconVariant="warning" />
      <ModalBody>
        <Form>
          {error && (
            <Alert data-testid="error-message-alert" isInline variant="danger" title="Error">
              {error.message}
            </Alert>
          )}
          <FormGroup>
            {description}
            <TextInput
              id="confirm-archive-input"
              data-testid="confirm-archive-input"
              aria-label="confirm archive input"
              value={confirmInputValue}
              onChange={(_e, newValue) => setConfirmInputValue(newValue)}
              onKeyDown={(event) => {
                if (event.key === 'Enter' && !isDisabled) {
                  onConfirm();
                }
              }}
            />
          </FormGroup>
        </Form>
      </ModalBody>
      <DashboardModalFooter
        onCancel={onClose}
        onSubmit={onConfirm}
        submitLabel="Archive"
        isSubmitLoading={isSubmitting}
        isSubmitDisabled={isDisabled}
      />
    </Modal>
  );
};

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

interface ArchiveModelVersionModalProps {
  onCancel: () => void;
  onSubmit: () => void;
  isOpen: boolean;
  modelVersionName: string;
}

export const ArchiveModelVersionModal: React.FC<ArchiveModelVersionModalProps> = ({
  onCancel,
  onSubmit,
  isOpen,
  modelVersionName,
}) => {
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [error, setError] = React.useState<Error>();
  const [confirmInputValue, setConfirmInputValue] = React.useState('');
  const isDisabled = confirmInputValue.trim() !== modelVersionName || isSubmitting;
  const notification = useNotification();

  const onClose = React.useCallback(() => {
    setConfirmInputValue('');
    onCancel();
  }, [onCancel]);

  const onConfirm = React.useCallback(async () => {
    setIsSubmitting(true);

    try {
      await onSubmit();
      onClose();
      notification.success(`${modelVersionName} archived.`);
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
      <b>{modelVersionName}</b> will be archived and unavailable for use unless it is restored.
      <br />
      <br />
      Type <strong>{modelVersionName}</strong> to confirm archiving:
    </>
  );

  return (
    <Modal
      isOpen={isOpen}
      variant="small"
      onClose={onClose}
      data-testid="archive-model-version-modal"
    >
      <ModalHeader title="Archive version?" titleIconVariant="warning" />
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

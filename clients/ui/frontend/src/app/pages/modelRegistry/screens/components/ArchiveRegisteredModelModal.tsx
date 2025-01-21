import * as React from 'react';
import { Flex, FlexItem, Stack, StackItem, TextInput } from '@patternfly/react-core';
import { Modal } from '@patternfly/react-core/deprecated';
import DashboardModalFooter from '~/shared/components/DashboardModalFooter';
import { useNotification } from '~/app/hooks/useNotification';

interface ArchiveRegisteredModelModalProps {
  onCancel: () => void;
  onSubmit: () => void;
  registeredModelName: string;
}

export const ArchiveRegisteredModelModal: React.FC<ArchiveRegisteredModelModalProps> = ({
  onCancel,
  onSubmit,
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
  }, [onSubmit, onClose, notification, registeredModelName]);

  return (
    <Modal
      isOpen
      title="Archive model?"
      titleIconVariant="warning"
      variant="small"
      onClose={onClose}
      footer={
        <DashboardModalFooter
          onCancel={onClose}
          onSubmit={onConfirm}
          submitLabel="Archive"
          isSubmitLoading={isSubmitting}
          isSubmitDisabled={isDisabled}
          error={error}
          alertTitle="Error"
        />
      }
      data-testid="archive-registered-model-modal"
    >
      <Stack hasGutter>
        <StackItem>
          <b>{registeredModelName}</b> and all of its versions will be archived and unavailable for
          use unless it is restored.
        </StackItem>
        <StackItem>
          <Flex direction={{ default: 'column' }} spaceItems={{ default: 'spaceItemsSm' }}>
            <FlexItem>
              Type <strong>{registeredModelName}</strong> to confirm archiving:
            </FlexItem>
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
          </Flex>
        </StackItem>
      </Stack>
    </Modal>
  );
};

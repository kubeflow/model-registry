import * as React from 'react';
import { Flex, FlexItem, Stack, StackItem, TextInput } from '@patternfly/react-core';
import { Modal } from '@patternfly/react-core/deprecated';
import { DashboardModalFooter } from 'mod-arch-shared';
import { useNotification } from '~/app/hooks/useNotification';

interface ArchiveModelVersionModalProps {
  onCancel: () => void;
  onSubmit: () => void;
  modelVersionName: string;
}

export const ArchiveModelVersionModal: React.FC<ArchiveModelVersionModalProps> = ({
  onCancel,
  onSubmit,
  modelVersionName,
}) => {
  const notification = useNotification();
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [error, setError] = React.useState<Error>();
  const [confirmInputValue, setConfirmInputValue] = React.useState('');
  const isDisabled = confirmInputValue.trim() !== modelVersionName || isSubmitting;

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
  }, [onSubmit, onClose, notification, modelVersionName]);

  return (
    <Modal
      isOpen
      title="Archive model version?"
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
      data-testid="archive-model-version-modal"
    >
      <Stack hasGutter>
        <StackItem>
          <b>{modelVersionName}</b> will be archived and unavailable for use unless it is restored.
        </StackItem>
        <StackItem>
          <Flex direction={{ default: 'column' }} spaceItems={{ default: 'spaceItemsSm' }}>
            <FlexItem>
              Type <strong>{modelVersionName}</strong> to confirm archiving:
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

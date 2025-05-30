import React from 'react';
import { Content, TextInput, Stack, StackItem } from '@patternfly/react-core';
import { Modal } from '@patternfly/react-core/deprecated';
import { DashboardModalFooter, ModelRegistryKind } from 'mod-arch-shared';

type DeleteModelRegistryModalProps = {
  modelRegistry: ModelRegistryKind;
  onClose: () => void;
  refresh: () => void;
};

const DeleteModelRegistryModal: React.FC<DeleteModelRegistryModalProps> = ({
  modelRegistry: mr,
  onClose,
  refresh,
}) => {
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [error, setError] = React.useState<Error>();
  const [confirmInputValue, setConfirmInputValue] = React.useState('');
  const isDisabled = confirmInputValue.trim() !== mr.metadata.name || isSubmitting;

  const onBeforeClose = () => {
    setConfirmInputValue('');
    setIsSubmitting(false);
    setError(undefined);
    onClose();
  };

  const onConfirm = async () => {
    setIsSubmitting(true);
    setError(undefined);
    try {
      // TODO: implement when CRD endpoint is ready
      // await deleteModelRegistryBackend(mr.metadata.name);
      await refresh();
      onBeforeClose();
    } catch (e) {
      if (e instanceof Error) {
        setError(e);
      }
      setIsSubmitting(false);
    }
  };

  return (
    <Modal
      data-testid="delete-mr-modal"
      titleIconVariant="warning"
      title="Delete model registry?"
      isOpen
      onClose={onClose}
      variant="medium"
      footer={
        <DashboardModalFooter
          submitLabel="Delete model registry"
          submitButtonVariant="danger"
          onSubmit={onConfirm}
          onCancel={onBeforeClose}
          isSubmitLoading={isSubmitting}
          isSubmitDisabled={isDisabled}
          error={error}
          alertTitle="Error deleting model registry"
        />
      }
    >
      <Stack hasGutter>
        <StackItem>
          <Content>
            <Content component="p">
              The <strong>{mr.metadata.name}</strong> model registry, its default group, and any
              permissions associated with it will be deleted. Data located in the database connected
              to the registry will be unaffected.
            </Content>
            <Content component="p">
              Type <strong>{mr.metadata.name}</strong> to confirm deletion:
            </Content>
          </Content>
        </StackItem>
        <StackItem>
          <TextInput
            id="confirm-delete-input"
            data-testid="confirm-delete-input"
            aria-label="Confirm delete input"
            value={confirmInputValue}
            onChange={(_e, newValue) => setConfirmInputValue(newValue)}
            onKeyDown={(event) => {
              if (event.key === 'Enter' && !isDisabled) {
                onConfirm();
              }
            }}
          />
        </StackItem>
      </Stack>
    </Modal>
  );
};

export default DeleteModelRegistryModal;

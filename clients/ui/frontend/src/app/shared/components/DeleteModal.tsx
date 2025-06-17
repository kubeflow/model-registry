// // TODO: Move this code to shared library once the migration completes.
import * as React from 'react';
import {
  Alert,
  Button,
  Flex,
  FlexItem,
  Stack,
  StackItem,
  TextInput,
  Modal,
  ModalBody,
  ModalHeader,
  ModalFooter,
} from '@patternfly/react-core';

type DeleteModalProps = {
  title: string;
  onClose: () => void;
  deleting: boolean;
  onDelete: () => void;
  deleteName: string;
  submitButtonLabel?: string;
  error?: Error;
  children: React.ReactNode;
  testId?: string;
  genericLabel?: boolean;
};

const DeleteModal: React.FC<DeleteModalProps> = ({
  children,
  title,
  onClose,
  deleting,
  onDelete,
  deleteName,
  error,
  submitButtonLabel = 'Delete',
  testId,
  genericLabel,
}) => {
  const [value, setValue] = React.useState('');

  const deleteNameSanitized = React.useMemo(
    () => deleteName.trim().replace(/\s+/g, ' '),
    [deleteName],
  );

  const onBeforeClose = (deleted: boolean) => {
    if (deleted) {
      onDelete();
    } else {
      onClose();
    }
  };

  return (
    <Modal
      isOpen
      onClose={() => onBeforeClose(false)}
      variant="small"
      data-testid={testId || 'delete-modal'}
    >
      <ModalHeader title={title} titleIconVariant="warning" />
      <ModalBody>
        <Stack hasGutter>
          <StackItem>{children}</StackItem>

          <StackItem>
            <Flex direction={{ default: 'column' }} spaceItems={{ default: 'spaceItemsSm' }}>
              <FlexItem>
                Type <strong>{deleteNameSanitized}</strong> to confirm
                {genericLabel ? '' : ' deletion'}:
              </FlexItem>

              <TextInput
                id="delete-modal-input"
                data-testid="delete-modal-input"
                aria-label="Delete modal input"
                value={value}
                onChange={(_e, newValue) => setValue(newValue)}
                onKeyDown={(event) => {
                  if (event.key === 'Enter' && value.trim() === deleteNameSanitized && !deleting) {
                    onDelete();
                  }
                }}
              />
            </Flex>
          </StackItem>

          {error && (
            <StackItem>
              <Alert
                data-testid="delete-model-error-message-alert"
                title={`Error deleting ${deleteNameSanitized}`}
                isInline
                variant="danger"
              >
                {error.message}
              </Alert>
            </StackItem>
          )}
        </Stack>
      </ModalBody>
      <ModalFooter>
        <Button
          key="delete-button"
          variant="danger"
          isLoading={deleting}
          isDisabled={deleting || value.trim() !== deleteNameSanitized}
          onClick={() => onBeforeClose(true)}
        >
          {submitButtonLabel}
        </Button>
        <Button key="cancel-button" variant="link" onClick={() => onBeforeClose(false)}>
          Cancel
        </Button>
      </ModalFooter>
    </Modal>
  );
};

export default DeleteModal;

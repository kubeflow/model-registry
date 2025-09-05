import {
  Modal,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Button,
  Checkbox,
  ModalVariant,
  Stack,
  StackItem,
} from '@patternfly/react-core';
import * as React from 'react';
import useDeletePropertiesModalAvailability from '~/app/hooks/useDeletePropertiesModalAvailability';

type DeletePropertiesModalProps = {
  isOpen: boolean;
  setIsOpen: (isOpen: boolean) => void;
  deleteProperty: () => Promise<unknown>;
  modelName?: string;
};

const DeletePropertiesModal: React.FC<DeletePropertiesModalProps> = ({
  modelName,
  isOpen,
  setIsOpen,
  deleteProperty,
}) => {
  const [dontShowModalValue, setDontShowModalValue] = useDeletePropertiesModalAvailability();
  return (
    <Modal
      isOpen={isOpen}
      onClose={() => setIsOpen(false)}
      variant={ModalVariant.small}
      ouiaId="DeletePropertyModal"
      aria-labelledby="delete-property-modal-title"
      aria-describedby="delete-property-modal-body"
      data-testid="delete-property-modal"
    >
      <ModalHeader
        title="Delete property from all model versions?"
        labelId="delete-property-modal-title"
      />
      <ModalBody id="delete-property-modal-body">
        <Stack hasGutter>
          <StackItem>
            Editing the model details will apply changes to all versions of the{' '}
            {modelName ? (
              <>
                <b>{modelName}</b> model.
              </>
            ) : (
              <>model.</>
            )}
          </StackItem>
          <StackItem>
            <Checkbox
              id="dont-show-again"
              label="Don't show this again"
              isChecked={dontShowModalValue}
              onChange={(_, checked) => setDontShowModalValue(checked)}
            />
          </StackItem>
        </Stack>
      </ModalBody>
      <ModalFooter>
        <Button
          key="confirm"
          variant="primary"
          onClick={() => {
            deleteProperty();
            setIsOpen(false);
          }}
          data-testid="delete-property-modal-confirm"
        >
          Confirm
        </Button>
        <Button
          key="cancel"
          variant="link"
          onClick={() => setIsOpen(false)}
          data-testid="delete-property-modal-cancel"
        >
          Cancel
        </Button>
      </ModalFooter>
    </Modal>
  );
};

export default DeletePropertiesModal;

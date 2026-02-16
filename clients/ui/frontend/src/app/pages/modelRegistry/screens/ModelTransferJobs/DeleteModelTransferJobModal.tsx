import * as React from 'react';
import {
  Alert,
  Button,
  Content,
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
} from '@patternfly/react-core';
import { ModelTransferJob } from '~/app/types';

type DeleteModelTransferJobModalProps = {
  job: ModelTransferJob;
  onClose: () => void;
  onDelete: (job: ModelTransferJob) => Promise<void>;
};

const DeleteModelTransferJobModal: React.FC<DeleteModelTransferJobModalProps> = ({
  job,
  onClose,
  onDelete,
}) => {
  const [isDeleting, setIsDeleting] = React.useState(false);
  const [error, setError] = React.useState<Error | undefined>();

  const handleDelete = React.useCallback(async () => {
    setIsDeleting(true);
    setError(undefined);
    try {
      await onDelete(job);
      onClose();
    } catch (e) {
      setError(e instanceof Error ? e : new Error(String(e)));
    } finally {
      setIsDeleting(false);
    }
  }, [job, onDelete, onClose]);

  return (
    <Modal
      isOpen
      onClose={onClose}
      variant="small"
      data-testid="delete-model-transfer-job-modal"
      aria-label="Delete model transfer job"
    >
      <ModalHeader title="Delete model transfer job?" />
      <ModalBody>
        <Content>
          <Content component="p">
            The <strong>{job.name}</strong> model transfer job will be deleted, but the storage
            location of the model will not be affected.
          </Content>
          {error && (
            <Alert
              data-testid="delete-model-transfer-job-error"
              title="Error deleting model transfer job"
              variant="danger"
            >
              {error.message}
            </Alert>
          )}
        </Content>
      </ModalBody>
      <ModalFooter>
        <Button
          key="delete"
          variant="primary"
          onClick={handleDelete}
          isDisabled={isDeleting}
          isLoading={isDeleting}
          data-testid="delete-model-transfer-job-submit"
        >
          Delete
        </Button>
        <Button
          key="cancel"
          variant="link"
          onClick={onClose}
          isDisabled={isDeleting}
          data-testid="delete-model-transfer-job-cancel"
        >
          Cancel
        </Button>
      </ModalFooter>
    </Modal>
  );
};

export default DeleteModelTransferJobModal;

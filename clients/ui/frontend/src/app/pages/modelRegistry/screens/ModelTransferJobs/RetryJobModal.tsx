import * as React from 'react';
import {
  Alert,
  Button,
  Checkbox,
  FormGroup,
  HelperText,
  HelperTextItem,
  Stack,
  StackItem,
  TextInput,
  Modal,
  ModalBody,
  ModalHeader,
  ModalFooter,
} from '@patternfly/react-core';
import { ModelTransferJob, ModelTransferJobUploadIntent } from '~/app/types';
import ResourceNameDefinitionTooltip from '~/concepts/k8s/ResourceNameDefinitionTooltip';
import { checkValidK8sName } from '~/concepts/k8s/K8sNameDescriptionField/utils';

/**
 * Generates a retry job name by appending or incrementing a numeric suffix.
 * e.g., "my-job" -> "my-job-2", "my-job-2" -> "my-job-3"
 */
const generateRetryJobName = (originalName: string): string => {
  const numericSuffixMatch = originalName.match(/^(.+)-(\d+)$/);
  if (numericSuffixMatch) {
    const [, baseName, numStr] = numericSuffixMatch;
    return `${baseName}-${parseInt(numStr, 10) + 1}`;
  }
  return `${originalName}-2`;
};

type RetryJobModalProps = {
  job: ModelTransferJob;
  onClose: () => void;
  onRetry: (newJobName: string, deleteOldJob: boolean) => Promise<void>;
};

const RetryJobModal: React.FC<RetryJobModalProps> = ({ job, onClose, onRetry }) => {
  const [newJobName, setNewJobName] = React.useState(() => generateRetryJobName(job.name));
  const [showResourceNameField, setShowResourceNameField] = React.useState(false);
  const [deleteOldJob, setDeleteOldJob] = React.useState(true);
  const [isRetrying, setIsRetrying] = React.useState(false);
  const [error, setError] = React.useState<Error | undefined>();

  const validation = React.useMemo(() => checkValidK8sName(newJobName), [newJobName]);
  const isNameValid = validation.valid && newJobName.length > 0 && newJobName.length <= 253;

  const handleRetry = async () => {
    setIsRetrying(true);
    setError(undefined);
    try {
      await onRetry(newJobName, deleteOldJob);
      onClose();
    } catch (e) {
      setError(e instanceof Error ? e : new Error(String(e)));
    } finally {
      setIsRetrying(false);
    }
  };

  return (
    <Modal isOpen onClose={onClose} variant="small" data-testid="retry-job-modal">
      <ModalHeader title="Retry model transfer job?" />
      <ModalBody>
        <Stack hasGutter>
          <StackItem>
            {job.uploadIntent === ModelTransferJobUploadIntent.CREATE_MODEL ? (
              <>
                A new transfer job will be created for the{' '}
                <strong>{job.registeredModelName}</strong> model.
              </>
            ) : (
              <>
                A new transfer job will be created for the <strong>{job.modelVersionName}</strong>{' '}
                version of <strong>{job.registeredModelName}</strong>.
              </>
            )}
          </StackItem>

          <StackItem>
            <FormGroup label="New model transfer job name" isRequired fieldId="retry-job-name">
              <TextInput
                id="retry-job-name"
                data-testid="retry-job-name-input"
                value={newJobName}
                onChange={(_e, value) => {
                  setNewJobName(value);
                  setShowResourceNameField(true);
                }}
                validated={!isNameValid && showResourceNameField ? 'error' : 'default'}
                isRequired
              />
            </FormGroup>
            <HelperText>
              {!showResourceNameField && (
                <HelperTextItem>
                  <Button
                    data-testid="retry-job-edit-resource-link"
                    variant="link"
                    isInline
                    onClick={() => setShowResourceNameField(true)}
                  >
                    Edit resource name
                  </Button>{' '}
                  <ResourceNameDefinitionTooltip />
                </HelperTextItem>
              )}
              {showResourceNameField && !validation.valid && validation.invalidCharacters && (
                <HelperTextItem variant="error">
                  Must start and end with a lowercase letter or number. Valid characters include
                  lowercase letters, numbers, and hyphens (-).
                </HelperTextItem>
              )}
              {showResourceNameField && newJobName.length > 253 && (
                <HelperTextItem variant="error">Cannot exceed 253 characters</HelperTextItem>
              )}
              {showResourceNameField && validation.valid && newJobName.length > 0 && (
                <HelperTextItem>
                  The resource name is used to identify your resource, and cannot be edited after
                  creation.
                </HelperTextItem>
              )}
            </HelperText>
          </StackItem>

          <StackItem>
            <Checkbox
              id="delete-old-job-checkbox"
              data-testid="delete-old-job-checkbox"
              label={
                <>
                  Delete the failed <strong>{job.name}</strong> transfer job
                </>
              }
              isChecked={deleteOldJob}
              onChange={(_e, checked) => setDeleteOldJob(checked)}
            />
          </StackItem>

          {error && (
            <StackItem>
              <Alert
                data-testid="retry-job-error-alert"
                title="Error retrying transfer job"
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
          key="retry-button"
          variant="primary"
          isLoading={isRetrying}
          isDisabled={isRetrying || !isNameValid}
          onClick={handleRetry}
          data-testid="retry-job-submit-button"
        >
          Retry
        </Button>
        <Button key="cancel-button" variant="link" onClick={onClose}>
          Cancel
        </Button>
      </ModalFooter>
    </Modal>
  );
};

export default RetryJobModal;

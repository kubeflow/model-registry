import * as React from 'react';
import {
  Alert,
  Button,
  Checkbox,
  Stack,
  StackItem,
  Modal,
  ModalBody,
  ModalHeader,
  ModalFooter,
} from '@patternfly/react-core';
import { ModelTransferJob, ModelTransferJobUploadIntent } from '~/app/types';
import K8sNameDescriptionField, {
  useK8sNameDescriptionFieldData,
} from '~/concepts/k8s/K8sNameDescriptionField/K8sNameDescriptionField';

/**
 * Generates a retry job name by appending or incrementing a numeric suffix.
 * Ensures the result doesn't exceed maxLength by truncating the base name if needed.
 * e.g., "my-job" -> "my-job-2", "my-job-2" -> "my-job-3"
 * @param originalName The original job name
 * @param maxLength Maximum length for the generated name (default: 63, DNS-1123 label limit)
 */
const generateRetryJobName = (originalName: string, maxLength = 63): string => {
  const numericSuffixMatch = originalName.match(/^(.+)-(\d+)$/);

  if (numericSuffixMatch) {
    const [, baseName, numStr] = numericSuffixMatch;
    const newSuffix = parseInt(numStr, 10) + 1;
    const suffixStr = `-${newSuffix}`;
    const maxBaseLength = maxLength - suffixStr.length;

    // Truncate base name to fit within limit
    const truncatedBase = baseName.slice(0, maxBaseLength);
    // Remove trailing dashes from truncation (K8s names must end with alphanumeric)
    return `${truncatedBase.replace(/-+$/, '')}${suffixStr}`;
  }

  // For first retry: truncate to maxLength - 2 to leave room for "-2"
  const maxBaseLength = maxLength - 2;
  const truncatedName = originalName.slice(0, maxBaseLength);
  return `${truncatedName.replace(/-+$/, '')}-2`;
};

type RetryJobModalProps = {
  job: ModelTransferJob;
  onClose: () => void;
  onRetry: (newJobName: string, newJobDisplayName: string, deleteOldJob: boolean) => Promise<void>;
};

const RetryJobModal: React.FC<RetryJobModalProps> = ({ job, onClose, onRetry }) => {
  const generatedName = generateRetryJobName(job.name);
  const { data: fieldData, onDataChange } = useK8sNameDescriptionFieldData({
    initialData: { name: generatedName },
    editableK8sName: true,
    maxK8sNameLength: 63,
  });
  const [deleteOldJob, setDeleteOldJob] = React.useState(true);
  const [isRetrying, setIsRetrying] = React.useState(false);
  const [error, setError] = React.useState<Error | undefined>();

  // Trigger initial k8sName generation from name on mount
  React.useEffect(() => {
    onDataChange('name', generatedName);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const { k8sName } = fieldData;
  const isNameValid =
    !k8sName.state.invalidCharacters && !k8sName.state.invalidLength && k8sName.value.length > 0;

  const handleRetry = async () => {
    setIsRetrying(true);
    setError(undefined);
    try {
      await onRetry(fieldData.k8sName.value, fieldData.name, deleteOldJob);
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
            <K8sNameDescriptionField
              data={fieldData}
              onDataChange={onDataChange}
              dataTestId="retry-job"
              nameLabel="New model transfer job name"
              hideDescription
            />
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

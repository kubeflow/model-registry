import { Button, Content, ContentVariants, Label, Truncate } from '@patternfly/react-core';
import { ActionsColumn, Td, Tr } from '@patternfly/react-table';
import {
  CheckCircleIcon,
  ExclamationCircleIcon,
  InProgressIcon,
  PendingIcon,
  BanIcon,
} from '@patternfly/react-icons';
import * as React from 'react';
import { useNavigate } from 'react-router-dom';
import ModelTimestamp from '~/app/pages/modelRegistry/screens/components/ModelTimestamp';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import { registeredModelUrl, modelVersionUrl } from '~/app/pages/modelRegistry/screens/routeUtils';
import {
  ModelTransferJob,
  ModelTransferJobStatus,
  ModelTransferJobUploadIntent,
} from '~/app/types';
import { EMPTY_CUSTOM_PROPERTY_VALUE } from '~/concepts/modelCatalog/const';
import ModelTransferJobStatusModal from './ModelTransferJobStatusModal';

type ModelTransferJobTableRowProps = {
  job: ModelTransferJob;
  onRequestDelete?: (job: ModelTransferJob) => void;
};

export const getStatusLabel = (
  status: ModelTransferJobStatus,
): {
  label: string;
  color: React.ComponentProps<typeof Label>['color'];
  icon: React.ReactNode;
} => {
  switch (status) {
    case ModelTransferJobStatus.COMPLETED:
      return { label: 'Complete', color: 'green', icon: <CheckCircleIcon /> };
    case ModelTransferJobStatus.RUNNING:
      return { label: 'Running', color: 'blue', icon: <InProgressIcon /> };
    case ModelTransferJobStatus.PENDING:
      return { label: 'Pending', color: 'grey', icon: <PendingIcon /> };
    case ModelTransferJobStatus.FAILED:
      return { label: 'Failed', color: 'red', icon: <ExclamationCircleIcon /> };
    case ModelTransferJobStatus.CANCELLED:
      return { label: 'Cancelled', color: 'grey', icon: <BanIcon /> };
    default:
      return { label: status, color: 'grey', icon: null };
  }
};

const ModelTransferJobTableRow: React.FC<ModelTransferJobTableRowProps> = ({
  job,
  onRequestDelete,
}) => {
  const navigate = useNavigate();
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);
  const [isStatusModalOpen, setIsStatusModalOpen] = React.useState(false);

  const handleModelNameClick = () => {
    if (job.registeredModelId) {
      navigate(registeredModelUrl(job.registeredModelId, preferredModelRegistry?.name));
    }
  };

  const handleVersionNameClick = () => {
    if (job.modelVersionId && job.registeredModelId) {
      navigate(
        modelVersionUrl(job.modelVersionId, job.registeredModelId, preferredModelRegistry?.name),
      );
    }
  };

  const statusInfo = getStatusLabel(job.status);

  const actions = React.useMemo(
    () => [
      {
        title: 'Delete',
        onClick: () => {
          onRequestDelete?.(job);
        },
      },
    ],
    [job, onRequestDelete],
  );

  return (
    <Tr>
      <Td dataLabel="Job name">
        <div data-testid="job-name">
          <Truncate content={job.name} />
        </div>
        {job.description && (
          <Content data-testid="job-description" component={ContentVariants.small}>
            <Truncate content={job.description} />
          </Content>
        )}
      </Td>
      <Td dataLabel="Model name">
        {job.registeredModelName ? (
          job.uploadIntent === ModelTransferJobUploadIntent.CREATE_MODEL &&
          job.status === ModelTransferJobStatus.COMPLETED ? (
            <Button variant="link" isInline onClick={handleModelNameClick}>
              <Truncate content={job.registeredModelName} />
            </Button>
          ) : (
            <Truncate content={job.registeredModelName} />
          )
        ) : (
          EMPTY_CUSTOM_PROPERTY_VALUE
        )}
      </Td>
      <Td dataLabel="Model version name">
        {job.modelVersionName ? (
          (job.uploadIntent === ModelTransferJobUploadIntent.CREATE_MODEL ||
            job.uploadIntent === ModelTransferJobUploadIntent.CREATE_VERSION) &&
          job.status === ModelTransferJobStatus.COMPLETED ? (
            <Button variant="link" isInline onClick={handleVersionNameClick}>
              <Truncate content={job.modelVersionName} />
            </Button>
          ) : (
            <Truncate content={job.modelVersionName} />
          )
        ) : (
          EMPTY_CUSTOM_PROPERTY_VALUE
        )}
      </Td>
      <Td dataLabel="Namespace">
        <Content component="p" data-testid="job-namespace">
          {job.namespace || EMPTY_CUSTOM_PROPERTY_VALUE}
        </Content>
      </Td>
      <Td dataLabel="Created">
        <ModelTimestamp timeSinceEpoch={job.createTimeSinceEpoch} />
      </Td>
      <Td dataLabel="Author">
        <Content component="p" data-testid="job-author">
          {job.author || EMPTY_CUSTOM_PROPERTY_VALUE}
        </Content>
      </Td>
      <Td dataLabel="Transfer job status">
        <Label
          color={statusInfo.color}
          icon={statusInfo.icon}
          data-testid="job-status"
          style={{ cursor: 'pointer' }}
          onClick={() => setIsStatusModalOpen(true)}
        >
          {statusInfo.label}
        </Label>
      </Td>
      <Td isActionCell>
        <ActionsColumn items={actions} />
      </Td>
      {isStatusModalOpen && (
        <ModelTransferJobStatusModal
          job={job}
          isOpen={isStatusModalOpen}
          onClose={() => setIsStatusModalOpen(false)}
        />
      )}
    </Tr>
  );
};

export default ModelTransferJobTableRow;

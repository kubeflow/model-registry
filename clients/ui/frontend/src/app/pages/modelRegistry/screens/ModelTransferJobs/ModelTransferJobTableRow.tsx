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
import { ModelTransferJob, ModelTransferJobStatus } from '~/app/types';
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

  const hasModelLink = Boolean(job.registeredModelId);
  const hasVersionLink = Boolean(job.modelVersionId && job.registeredModelId);

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
          hasModelLink ? (
            <Button
              variant="link"
              isInline
              onClick={handleModelNameClick}
              data-testid="job-model-link"
            >
              <Truncate content={job.registeredModelName} />
            </Button>
          ) : (
            <span data-testid="job-model-name">
              <Truncate content={job.registeredModelName} />
            </span>
          )
        ) : (
          EMPTY_CUSTOM_PROPERTY_VALUE
        )}
      </Td>
      <Td dataLabel="Model version name">
        {job.modelVersionName ? (
          hasVersionLink ? (
            <Button
              variant="link"
              isInline
              onClick={handleVersionNameClick}
              data-testid="job-version-link"
            >
              <Truncate content={job.modelVersionName} />
            </Button>
          ) : (
            <span data-testid="job-version-name">
              <Truncate content={job.modelVersionName} />
            </span>
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

import {
  Button,
  Content,
  ContentVariants,
  HelperText,
  HelperTextItem,
  Label,
  Truncate,
} from '@patternfly/react-core';
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

type ModelTransferJobTableRowProps = {
  job: ModelTransferJob;
};

const getStatusLabel = (
  status: ModelTransferJobStatus,
): { label: string; color: React.ComponentProps<typeof Label>['color']; icon: React.ReactNode } => {
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

const ModelTransferJobTableRow: React.FC<ModelTransferJobTableRowProps> = ({ job }) => {
  const navigate = useNavigate();
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);

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

  const actions = React.useMemo(() => {
    const items = [];

    // Show Retry action for failed jobs
    if (job.status === ModelTransferJobStatus.FAILED) {
      items.push({
        title: 'Retry',
        onClick: () => {
          // TODO: Implement retry functionality
        },
      });
    }

    // Always show Delete action
    items.push({
      title: 'Delete',
      onClick: () => {
        // TODO: Implement delete functionality
      },
    });

    return items;
  }, [job.status]);

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
          <Button variant="link" isInline onClick={handleModelNameClick}>
            <Truncate content={job.registeredModelName} />
          </Button>
        ) : (
          EMPTY_CUSTOM_PROPERTY_VALUE
        )}
      </Td>
      <Td dataLabel="Model version name">
        {job.modelVersionName ? (
          <Button variant="link" isInline onClick={handleVersionNameClick}>
            <Truncate content={job.modelVersionName} />
          </Button>
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
        <div>
          <Label color={statusInfo.color} icon={statusInfo.icon} data-testid="job-status">
            {statusInfo.label}
          </Label>
          {job.status === ModelTransferJobStatus.FAILED && job.errorMessage && (
            <HelperText data-testid="job-error-message">
              <HelperTextItem variant="error">
                <Truncate content={job.errorMessage} />
              </HelperTextItem>
            </HelperText>
          )}
        </div>
      </Td>
      <Td isActionCell>
        <ActionsColumn items={actions} />
      </Td>
    </Tr>
  );
};

export default ModelTransferJobTableRow;

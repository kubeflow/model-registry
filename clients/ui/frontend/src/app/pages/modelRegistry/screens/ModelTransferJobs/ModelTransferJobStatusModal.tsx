import * as React from 'react';
import {
  Alert,
  Button,
  Flex,
  FlexItem,
  Label,
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  Tab,
  TabContent,
  TabContentBody,
  Tabs,
  TabTitleText,
  Spinner,
} from '@patternfly/react-core';
import {
  useFetchState,
  FetchStateCallbackPromise,
  NotReadyError,
  POLL_INTERVAL,
} from 'mod-arch-core';
import {
  ModelTransferJob,
  ModelTransferJobEvent,
  ModelTransferJobStatus,
  ModelTransferJobUploadIntent,
} from '~/app/types';
import EventLog from '~/app/shared/components/EventLog';
import { useModelRegistryAPI } from '~/app/hooks/useModelRegistryAPI';
import { getStatusLabel } from './ModelTransferJobTableRow';

type ModelTransferJobStatusModalProps = {
  job: ModelTransferJob;
  isOpen: boolean;
  onClose: () => void;
};

const ModelTransferJobStatusModal: React.FC<ModelTransferJobStatusModalProps> = ({
  job,
  isOpen,
  onClose,
}) => {
  const [activeTabKey, setActiveTabKey] = React.useState(0);
  const statusInfo = getStatusLabel(job.status);
  const { api, apiAvailable } = useModelRegistryAPI();

  // Determine if we should poll for event updates (only for active jobs)
  const shouldPollEvents =
    isOpen &&
    (job.status === ModelTransferJobStatus.PENDING ||
      job.status === ModelTransferJobStatus.RUNNING);

  // Fetch events with useFetchState - memoized by job name
  const fetchEvents = React.useCallback<FetchStateCallbackPromise<ModelTransferJobEvent[]>>(
    (opts) => {
      if (!isOpen || !apiAvailable) {
        return Promise.reject(new NotReadyError('Modal is closed or API not available'));
      }
      return api.getModelTransferJobEvents(opts, job.name);
    },
    [isOpen, apiAvailable, api, job.name],
  );

  const [events, eventsLoaded, eventsLoadError] = useFetchState(fetchEvents, [], {
    initialPromisePurity: true,
    refreshRate: shouldPollEvents ? POLL_INTERVAL : undefined,
  });

  if (!isOpen) {
    return null;
  }

  const getModalTitle = (intent: ModelTransferJobUploadIntent): string => {
    switch (intent) {
      case ModelTransferJobUploadIntent.CREATE_MODEL:
        return 'Model creation status';
      case ModelTransferJobUploadIntent.CREATE_VERSION:
        return 'Model version status';
      default:
        return 'Transfer job status';
    }
  };

  const title = (
    <Flex spaceItems={{ default: 'spaceItemsSm' }} alignItems={{ default: 'alignItemsCenter' }}>
      <FlexItem>{getModalTitle(job.uploadIntent)}</FlexItem>
      <FlexItem>
        <Label color={statusInfo.color} icon={statusInfo.icon}>
          {statusInfo.label}
        </Label>
      </FlexItem>
    </Flex>
  );

  return (
    <Modal variant="medium" isOpen onClose={onClose} data-testid="transfer-job-status-modal">
      <ModalHeader title={title} />
      <ModalBody>
        {job.status === ModelTransferJobStatus.FAILED && (
          <Alert
            variant="danger"
            isInline
            title={job.errorMessage || 'Failure reason (unknown)'}
            className="pf-v6-u-mb-md"
            data-testid="transfer-job-failure-alert"
          />
        )}
        <Tabs
          activeKey={activeTabKey}
          onSelect={(_event, tabKey) => setActiveTabKey(Number(tabKey))}
          data-testid="transfer-job-status-tabs"
        >
          <Tab eventKey={0} title={<TabTitleText>Event log</TabTitleText>}>
            <TabContent id="event-log-tab" activeKey={activeTabKey} eventKey={0}>
              <TabContentBody hasPadding>
                {!eventsLoaded ? (
                  <Flex justifyContent={{ default: 'justifyContentCenter' }}>
                    <Spinner size="lg" />
                  </Flex>
                ) : eventsLoadError ? (
                  <Alert variant="danger" isInline title="Failed to load events">
                    {eventsLoadError.message}
                  </Alert>
                ) : (
                  <EventLog events={events} data-testid="transfer-job-event-log" />
                )}
              </TabContentBody>
            </TabContent>
          </Tab>
        </Tabs>
      </ModalBody>
      <ModalFooter>
        <Button variant="link" onClick={onClose}>
          Close
        </Button>
      </ModalFooter>
    </Modal>
  );
};

export default ModelTransferJobStatusModal;

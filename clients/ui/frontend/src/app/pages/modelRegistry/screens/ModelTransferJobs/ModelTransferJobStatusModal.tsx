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
  const [events, setEvents] = React.useState<ModelTransferJobEvent[]>([]);
  const [isLoadingEvents, setIsLoadingEvents] = React.useState(false);
  const [eventsError, setEventsError] = React.useState<string | null>(null);
  const statusInfo = getStatusLabel(job.status);
  const { api, apiAvailable } = useModelRegistryAPI();

  React.useEffect(() => {
    if (isOpen && apiAvailable) {
      setIsLoadingEvents(true);
      setEventsError(null);
      api
        .getModelTransferJobEvents({}, job.name)
        .then((fetchedEvents: React.SetStateAction<ModelTransferJobEvent[]>) => {
          setEvents(fetchedEvents);
          setIsLoadingEvents(false);
        })
        .catch((error: { message: string }) => {
          setEventsError(error.message || 'Failed to load events');
          setIsLoadingEvents(false);
        });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isOpen, job.name]);

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
                {isLoadingEvents ? (
                  <Flex justifyContent={{ default: 'justifyContentCenter' }}>
                    <Spinner size="lg" />
                  </Flex>
                ) : eventsError ? (
                  <Alert variant="danger" isInline title="Failed to load events">
                    {eventsError}
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

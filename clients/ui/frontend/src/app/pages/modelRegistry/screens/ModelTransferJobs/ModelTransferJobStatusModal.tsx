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
} from '@patternfly/react-core';
import { ModelTransferJob, ModelTransferJobStatus } from '~/app/types';
import EventLog from '~/app/shared/components/EventLog';
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

  if (!isOpen) {
    return null;
  }

  const title = (
    <Flex spaceItems={{ default: 'spaceItemsSm' }} alignItems={{ default: 'alignItemsCenter' }}>
      <FlexItem>Model version status</FlexItem>
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
          >
            {job.errorDescription && <p>{job.errorDescription}</p>}
          </Alert>
        )}
        <Tabs
          activeKey={activeTabKey}
          onSelect={(_event, tabKey) => setActiveTabKey(Number(tabKey))}
          data-testid="transfer-job-status-tabs"
        >
          <Tab eventKey={0} title={<TabTitleText>Event log</TabTitleText>}>
            <TabContent id="event-log-tab" activeKey={activeTabKey} eventKey={0}>
              <TabContentBody hasPadding>
                <EventLog events={job.events ?? []} data-testid="transfer-job-event-log" />
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

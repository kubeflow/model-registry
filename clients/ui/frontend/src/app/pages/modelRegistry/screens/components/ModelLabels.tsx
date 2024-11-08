import { Button, Label, LabelGroup, Popover, SearchInput, Content } from '@patternfly/react-core';
import { Modal, ModalVariant } from '@patternfly/react-core/deprecated';
import React from 'react';
import { ModelVersion, RegisteredModel } from '~/app/types';
import useDebounceCallback from '~/shared/utilities/useDebounceCallback';
import { getLabels } from '~/app/pages/modelRegistry/screens/utils';

// Threshold count to decide whether to choose modal or popover
const MODAL_THRESHOLD = 4;

type ModelLabelsProps = {
  name: string;
  customProperties: RegisteredModel['customProperties'] | ModelVersion['customProperties'];
};

const ModelLabels: React.FC<ModelLabelsProps> = ({ name, customProperties }) => {
  const [isLabelModalOpen, setIsLabelModalOpen] = React.useState(false);
  const [searchValue, setSearchValue] = React.useState('');

  const allLabels = getLabels(customProperties);
  const filteredLabels = allLabels.filter(
    (label) => label && label.toLowerCase().includes(searchValue.toLowerCase()),
  );

  const doSetSearchDebounced = useDebounceCallback(setSearchValue);

  const labelsComponent = (labels: string[], textMaxWidth?: string) =>
    labels.map((label, index) => (
      <Label color="blue" data-testid="label" textMaxWidth={textMaxWidth || '40ch'} key={index}>
        {label}
      </Label>
    ));

  const getLabelComponent = (labels: JSX.Element[]) => {
    const labelCount = labels.length;
    if (labelCount) {
      return labelCount > MODAL_THRESHOLD
        ? getLabelModal(labelCount)
        : getLabelPopover(labelCount, labels);
    }
    return null;
  };

  const getLabelPopover = (labelCount: number, labels: JSX.Element[]) => (
    <Popover
      bodyContent={
        <LabelGroup data-testid="popover-label-group" numLabels={labelCount}>
          {labels}
        </LabelGroup>
      }
    >
      <Label data-testid="popover-label-text" variant="overflow">
        {labelCount} more
      </Label>
    </Popover>
  );

  const getLabelModal = (labelCount: number) => (
    <Label
      data-testid="modal-label-text"
      variant="overflow"
      onClick={() => setIsLabelModalOpen(true)}
    >
      {labelCount} more
    </Label>
  );

  const labelModal = (
    <Modal
      variant={ModalVariant.small}
      title="Labels"
      isOpen={isLabelModalOpen}
      onClose={() => setIsLabelModalOpen(false)}
      description={
        <Content component="p">
          The following are all the labels of <strong>{name}</strong>
        </Content>
      }
      actions={[
        <Button
          data-testid="close-modal"
          key="close"
          variant="primary"
          onClick={() => setIsLabelModalOpen(false)}
        >
          Close
        </Button>,
      ]}
    >
      <SearchInput
        aria-label="Label modal search"
        data-testid="label-modal-search"
        placeholder="Find a label"
        value={searchValue}
        onChange={(_event, value) => doSetSearchDebounced(value)}
        onClear={() => setSearchValue('')}
      />
      <br />
      <LabelGroup data-testid="modal-label-group" numLabels={allLabels.length}>
        {labelsComponent(filteredLabels, '50ch')}
      </LabelGroup>
    </Modal>
  );

  if (Object.keys(customProperties).length === 0) {
    return '-';
  }

  return (
    <>
      <LabelGroup numLabels={MODAL_THRESHOLD}>
        {labelsComponent(allLabels.slice(0, 3))}
        {getLabelComponent(labelsComponent(allLabels.slice(3)))}
      </LabelGroup>
      {labelModal}
    </>
  );
};

export default ModelLabels;

import * as React from 'react';
import {
  Alert,
  Flex,
  FlexItem,
  Label,
  Modal,
  ModalBody,
  ModalHeader,
  ModalVariant,
} from '@patternfly/react-core';
import { ExclamationCircleIcon } from '@patternfly/react-icons';

type CatalogSourceStatusErrorModalProps = {
  isOpen: boolean;
  onClose: () => void;
  errorMessage: string;
};

const CatalogSourceStatusErrorModal: React.FC<CatalogSourceStatusErrorModalProps> = ({
  isOpen,
  onClose,
  errorMessage,
}) => {
  const titleWithLabel = (
    <Flex spaceItems={{ default: 'spaceItemsSm' }} alignItems={{ default: 'alignItemsCenter' }}>
      <FlexItem>Source status</FlexItem>
      <FlexItem>
        <Label color="red" icon={<ExclamationCircleIcon />}>
          Failed
        </Label>
      </FlexItem>
    </Flex>
  );

  return (
    <Modal
      variant={ModalVariant.medium}
      isOpen={isOpen}
      onClose={onClose}
      data-testid="catalog-source-status-error-modal"
    >
      <ModalHeader title={titleWithLabel} />
      <ModalBody>
        <Alert
          variant="danger"
          isInline
          title="Validation failed"
          data-testid="catalog-source-status-error-alert"
        >
          <p data-testid="catalog-source-status-error-details">
            The system cannot establish a connection to the source. Ensure that the organization and
            access token are accurate, then try again.
          </p>
          {errorMessage && <p data-testid="catalog-source-status-error-message">{errorMessage}</p>}
        </Alert>
      </ModalBody>
    </Modal>
  );
};

export default CatalogSourceStatusErrorModal;

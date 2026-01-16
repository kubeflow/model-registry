// TODO this component was copied from odh-dashboard temporarily and should be abstracted out into mod-arch-shared.

import * as React from 'react';
import { useId } from 'react';
import {
  Modal,
  ModalBody,
  ModalHeader,
  ModalFooter,
  Button,
  ButtonProps,
  ModalProps,
  ModalHeaderProps,
} from '@patternfly/react-core';
import '~/concepts/dashboard/ModalStyles.scss';

export type ButtonAction = {
  label: string;
  onClick: () => void;
  variant?: ButtonProps['variant'];
  dataTestId?: string;
};

type ContentModalProps = {
  onClose: () => void;
  contents: React.ReactNode;
  title: string | React.ReactNode;
  buttonActions?: ButtonAction[];
  description?: React.ReactNode;
  disableFocusTrap?: boolean;
  dataTestId?: string;
  bodyClassName?: string;
  variant?: ModalProps['variant'];
  bodyLabel?: string;
  titleIconVariant?: ModalHeaderProps['titleIconVariant'];
};

// all buttons are always 'on'/enabled in this modal
/**
 * Generic Modal component for better accessibility.
 * This is used to make the modal more accessible for users who use keyboard navigation.
 * and easier to use in general.
 *
 * The buttons are defined via the buttonActions prop
 *
 * originally, the 'cancel' button was meant to be activated upon enter, but that is not standard UX;
 * and the 'escape' button already closes the dialog; which is standard. (this is implemented by patternfly)
 */
const ContentModal: React.FC<ContentModalProps> = ({
  onClose,
  contents,
  title,
  buttonActions,
  description,
  disableFocusTrap,
  dataTestId = 'content-modal',
  bodyClassName = 'odh-modal__content-height',
  variant = 'medium',
  bodyLabel,
  titleIconVariant,
}) => {
  const headingId = useId(); // used for aria-labelledby (a11y)
  const descriptionId = useId(); // used for aria-describedby (a11y)

  return (
    <Modal
      data-testid={dataTestId}
      isOpen
      variant={variant}
      onClose={onClose}
      disableFocusTrap={disableFocusTrap}
      aria-labelledby={headingId}
      aria-describedby={description ? descriptionId : undefined}
    >
      <ModalHeader
        title={title}
        labelId={headingId}
        description={description ? <div id={descriptionId}>{description}</div> : undefined}
        titleIconVariant={titleIconVariant}
        data-testid="generic-modal-header"
      />
      <ModalBody className={bodyClassName} aria-label={bodyLabel}>
        {contents}
      </ModalBody>
      <ModalFooter>
        {buttonActions?.map((action, index) => (
          <Button
            key={`${action.label}-${index}`}
            variant={action.variant}
            onClick={action.onClick}
            data-testid={action.dataTestId}
          >
            {action.label}
          </Button>
        ))}
      </ModalFooter>
    </Modal>
  );
};

export default ContentModal;

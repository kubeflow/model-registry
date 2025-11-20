import * as React from 'react';
import { Button, ButtonProps } from '@patternfly/react-core';
import { BUTTON_LABELS } from '~/app/pages/modelCatalogSettings/constants';

type PreviewButtonProps = {
  onClick: () => void;
  isDisabled: boolean;
  variant?: ButtonProps['variant'];
  testId?: string;
};

const PreviewButton: React.FC<PreviewButtonProps> = ({
  onClick,
  isDisabled,
  variant = 'primary',
  testId = 'preview-button',
}) => (
  <Button variant={variant} onClick={onClick} isDisabled={isDisabled} data-testid={testId}>
    {BUTTON_LABELS.PREVIEW}
  </Button>
);

export default PreviewButton;

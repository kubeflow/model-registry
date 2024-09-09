import React from 'react';
import {
  Button,
  ButtonVariant,
  EmptyState,
  EmptyStateActions,
  EmptyStateBody,
  EmptyStateFooter,
  EmptyStateVariant,
} from '@patternfly/react-core';
import { PlusCircleIcon } from '@patternfly/react-icons';

type EmptyModelRegistryStateType = {
  testid?: string;
  title: string;
  description: string;
  primaryActionText?: string;
  primaryActionOnClick?: () => void;
  secondaryActionText?: string;
  secondaryActionOnClick?: () => void;
  headerIcon?: React.ComponentType;
  customAction?: React.ReactNode;
};

const EmptyModelRegistryState: React.FC<EmptyModelRegistryStateType> = ({
  testid,
  title,
  description,
  primaryActionText,
  secondaryActionText,
  primaryActionOnClick,
  secondaryActionOnClick,
  headerIcon,
  customAction,
}) => (
  <EmptyState
    icon={headerIcon ?? PlusCircleIcon}
    titleText={title}
    variant={EmptyStateVariant.sm}
    data-testid={testid}
  >
    <EmptyStateBody>{description}</EmptyStateBody>
    <EmptyStateFooter>
      {primaryActionText && (
        <EmptyStateActions>
          <Button
            data-testid="empty-model-registry-primary-action"
            variant={ButtonVariant.primary}
            onClick={primaryActionOnClick}
          >
            {primaryActionText}
          </Button>
        </EmptyStateActions>
      )}

      {secondaryActionText && (
        <EmptyStateActions>
          <Button
            data-testid="empty-model-registry-secondary-action"
            variant="link"
            onClick={secondaryActionOnClick}
          >
            {secondaryActionText}
          </Button>
        </EmptyStateActions>
      )}

      {customAction && <EmptyStateActions>{customAction}</EmptyStateActions>}
    </EmptyStateFooter>
  </EmptyState>
);

export default EmptyModelRegistryState;

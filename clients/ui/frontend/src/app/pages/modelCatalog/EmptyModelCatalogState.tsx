import React from 'react';
import {
  EmptyState,
  EmptyStateActions,
  EmptyStateBody,
  EmptyStateFooter,
  EmptyStateVariant,
} from '@patternfly/react-core';
import { PlusCircleIcon } from '@patternfly/react-icons';

type EmptyModelCatalogStateType = {
  testid?: string;
  className?: string;
  title: string;
  description: React.ReactNode;
  headerIcon?: React.ComponentType;
  children?: React.ReactNode;
  customAction?: React.ReactNode;
  variant?: EmptyStateVariant;
};

const EmptyModelCatalogState: React.FC<EmptyModelCatalogStateType> = ({
  testid,
  className,
  title,
  description,
  headerIcon,
  children,
  customAction,
  variant = EmptyStateVariant.sm,
}) => (
  <EmptyState
    className={className}
    icon={headerIcon ?? PlusCircleIcon}
    titleText={title}
    variant={variant}
    data-testid={testid}
  >
    <EmptyStateBody>{description}</EmptyStateBody>
    {children}
    <EmptyStateFooter>
      {customAction && <EmptyStateActions>{customAction}</EmptyStateActions>}
    </EmptyStateFooter>
  </EmptyState>
);

export default EmptyModelCatalogState;

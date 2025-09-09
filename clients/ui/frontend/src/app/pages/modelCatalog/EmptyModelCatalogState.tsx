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
  title: string;
  description: string;
  headerIcon?: React.ComponentType;
  children?: React.ReactNode;
  customAction?: React.ReactNode;
};

const EmptyModelCatalogState: React.FC<EmptyModelCatalogStateType> = ({
  testid,
  title,
  description,
  headerIcon,
  children,
  customAction,
}) => (
  <EmptyState
    icon={headerIcon ?? PlusCircleIcon}
    titleText={title}
    variant={EmptyStateVariant.sm}
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

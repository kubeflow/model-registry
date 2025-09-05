import React from 'react';
import { EmptyState, EmptyStateBody, EmptyStateVariant } from '@patternfly/react-core';
import { PlusCircleIcon } from '@patternfly/react-icons';

type EmptyModelCatalogStateType = {
  testid?: string;
  title: string;
  description: string;
  headerIcon?: React.ComponentType;
  children?: React.ReactNode;
};

const EmptyModelCatalogState: React.FC<EmptyModelCatalogStateType> = ({
  testid,
  title,
  description,
  headerIcon,
  children,
}) => (
  <EmptyState
    icon={headerIcon ?? PlusCircleIcon}
    titleText={title}
    variant={EmptyStateVariant.sm}
    data-testid={testid}
  >
    <EmptyStateBody>{description}</EmptyStateBody>
    {children}
  </EmptyState>
);

export default EmptyModelCatalogState;

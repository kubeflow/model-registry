import * as React from 'react';
import {
  EmptyState,
  EmptyStateBody,
  Stack,
  StackItem,
  EmptyStateFooter,
} from '@patternfly/react-core';
import { PathMissingIcon } from '@patternfly/react-icons';

type EmptyStateErrorMessageProps = {
  children?: React.ReactNode;
  title: string;
  bodyText: string;
};

const EmptyStateErrorMessage: React.FC<EmptyStateErrorMessageProps> = ({
  title,
  bodyText,
  children,
}) => (
  <EmptyState headingLevel="h2" icon={PathMissingIcon} titleText={title}>
    <EmptyStateFooter>
      <Stack hasGutter>
        <StackItem>
          <EmptyStateBody>{bodyText}</EmptyStateBody>
        </StackItem>
        {children && <StackItem>{children}</StackItem>}
      </Stack>
    </EmptyStateFooter>
  </EmptyState>
);

export default EmptyStateErrorMessage;

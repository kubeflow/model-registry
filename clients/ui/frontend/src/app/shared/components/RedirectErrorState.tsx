import {
  PageSection,
  EmptyState,
  EmptyStateVariant,
  EmptyStateBody,
  EmptyStateFooter,
  EmptyStateActions,
} from '@patternfly/react-core';
import { ExclamationCircleIcon } from '@patternfly/react-icons';
import React from 'react';

type RedirectErrorStateProps = {
  title?: string;
  errorMessage?: string;
  actions?: React.ReactNode | React.ReactNode[];
};

/**
 * A component that displays an error state with optional title, message and actions
 * Used for showing redirect/navigation errors with fallback options
 *
 * Props for the RedirectErrorState component
 * @property {string} [title] - Optional title text to display in the error state
 * @property {string} [errorMessage] - Optional error message to display
 * @property {React.ReactNode | React.ReactNode[]} [actions] - Custom action buttons/elements to display
 *
 *
 * @example
 * ```tsx
 * // With custom actions
 * <RedirectErrorState
 *   title="Error redirecting to pipelines"
 *   errorMessage={error.message}
 *   actions={
 *     <>
 *       <Button variant="link" onClick={() => navigate('/pipelines')}>
 *         Go to pipelines
 *       </Button>
 *       <Button variant="link" onClick={() => navigate('/experiments')}>
 *         Go to experiments
 *       </Button>
 *     </>
 *   }
 * />
 * ```
 */

const RedirectErrorState: React.FC<RedirectErrorStateProps> = ({
  title,
  errorMessage,
  actions,
}) => (
  <PageSection hasBodyWrapper={false} isFilled>
    <EmptyState
      headingLevel="h1"
      icon={ExclamationCircleIcon}
      titleText={title ?? 'Error redirecting'}
      variant={EmptyStateVariant.lg}
      data-testid="redirect-error"
    >
      {errorMessage && <EmptyStateBody>{errorMessage}</EmptyStateBody>}
      {actions && (
        <EmptyStateFooter>
          <EmptyStateActions>{actions}</EmptyStateActions>
        </EmptyStateFooter>
      )}
    </EmptyState>
  </PageSection>
);

export default RedirectErrorState;

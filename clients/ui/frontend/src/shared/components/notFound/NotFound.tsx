import * as React from 'react';
import { HomeIcon, PathMissingIcon } from '@patternfly/react-icons';
import {
  Button,
  EmptyState,
  EmptyStateBody,
  EmptyStateFooter,
  EmptyStateVariant,
  PageSection,
} from '@patternfly/react-core';

const NotFound: React.FC = () => (
  <PageSection hasBodyWrapper={false}>
    <EmptyState
      headingLevel="h2"
      icon={PathMissingIcon}
      titleText="We can‘t find that page"
      variant={EmptyStateVariant.full}
      data-testid="not-found-page"
    >
      <EmptyStateBody data-testid="not-found-page-description">
        Another page might have what you need. Return to the home page.
      </EmptyStateBody>
      <EmptyStateFooter>
        <Button
          icon={<HomeIcon />}
          data-testid="home-page-button"
          component="a"
          href="/"
          variant="primary"
        >
          Home
        </Button>
      </EmptyStateFooter>
    </EmptyState>
  </PageSection>
);

export default NotFound;

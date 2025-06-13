import React from 'react';
import { ExclamationCircleIcon, QuestionCircleIcon } from '@patternfly/react-icons';
import {
  PageSection,
  Content,
  EmptyState,
  EmptyStateVariant,
  Spinner,
  EmptyStateBody,
  PageBreadcrumb,
  StackItem,
  Stack,
  Flex,
  FlexItem,
} from '@patternfly/react-core';

type ApplicationsPageProps = {
  title?: React.ReactNode;
  breadcrumb?: React.ReactNode;
  description?: React.ReactNode;
  loaded: boolean;
  empty: boolean;
  loadError?: Error;
  loadErrorPage?: React.ReactNode;
  children?: React.ReactNode;
  errorMessage?: string;
  emptyMessage?: string;
  emptyStatePage?: React.ReactNode;
  headerAction?: React.ReactNode;
  headerContent?: React.ReactNode;
  provideChildrenPadding?: boolean;
  removeChildrenTopPadding?: boolean;
  subtext?: React.ReactNode;
  loadingContent?: React.ReactNode;
  noHeader?: boolean;
};

const ApplicationsPage: React.FC<ApplicationsPageProps> = ({
  title,
  breadcrumb,
  description,
  loaded,
  empty,
  loadError,
  loadErrorPage,
  children,
  errorMessage,
  emptyMessage,
  emptyStatePage,
  headerAction,
  headerContent,
  provideChildrenPadding,
  removeChildrenTopPadding,
  subtext,
  loadingContent,
  noHeader,
}) => {
  const renderHeader = () => (
    <PageSection hasBodyWrapper={false}>
      <Stack hasGutter>
        <StackItem>
          <Flex
            justifyContent={{ default: 'justifyContentSpaceBetween' }}
            alignItems={{ default: 'alignItemsFlexStart' }}
          >
            <FlexItem flex={{ default: 'flex_1' }}>
              <Content component="h1" data-testid="app-page-title">
                {title}
              </Content>
              <Stack hasGutter>
                {subtext && <StackItem>{subtext}</StackItem>}
                {description && <StackItem>{description}</StackItem>}
              </Stack>
            </FlexItem>
            <FlexItem>{headerAction}</FlexItem>
          </Flex>
        </StackItem>
        {headerContent && <StackItem>{headerContent}</StackItem>}
      </Stack>
    </PageSection>
  );

  const renderContents = () => {
    if (loadError) {
      return !loadErrorPage ? (
        <PageSection hasBodyWrapper={false} isFilled>
          <EmptyState
            headingLevel="h1"
            icon={ExclamationCircleIcon}
            titleText={errorMessage !== undefined ? errorMessage : 'Error loading components'}
            variant={EmptyStateVariant.lg}
            data-id="error-empty-state"
          >
            <EmptyStateBody data-testid="error-empty-state-body">
              {loadError.message}
            </EmptyStateBody>
          </EmptyState>
        </PageSection>
      ) : (
        loadErrorPage
      );
    }

    if (!loaded) {
      return (
        loadingContent || (
          <PageSection hasBodyWrapper={false} isFilled>
            <EmptyState
              headingLevel="h1"
              titleText="Loading"
              variant={EmptyStateVariant.lg}
              data-id="loading-empty-state"
            >
              <Spinner size="xl" />
            </EmptyState>
          </PageSection>
        )
      );
    }

    if (empty) {
      return !emptyStatePage ? (
        <PageSection hasBodyWrapper={false} isFilled>
          <EmptyState
            headingLevel="h1"
            icon={QuestionCircleIcon}
            titleText={emptyMessage !== undefined ? emptyMessage : 'No Components Found'}
            variant={EmptyStateVariant.lg}
            data-id="empty-empty-state"
          />
        </PageSection>
      ) : (
        emptyStatePage
      );
    }

    if (provideChildrenPadding) {
      return (
        <PageSection isFilled style={removeChildrenTopPadding ? { paddingTop: 0 } : undefined}>
          {children}
        </PageSection>
      );
    }

    return children;
  };

  return (
    // TODO: PageBreadcrumb and the PageSection items here are children of the DrawerBody not the PageMain. DrawerBody is not flex which the Page items expect the parent to be.
    <Flex
      direction={{ default: 'column' }}
      flexWrap={{ default: 'nowrap' }}
      style={{ height: '100%' }}
    >
      <FlexItem>
        {breadcrumb && <PageBreadcrumb hasBodyWrapper={false}>{breadcrumb}</PageBreadcrumb>}
      </FlexItem>
      <FlexItem>{!noHeader && renderHeader()}</FlexItem>
      <FlexItem flex={{ default: 'flex_1' }}>
        <Flex
          direction={{ default: 'column' }}
          style={{ height: '100%' }}
          flexWrap={{ default: 'nowrap' }}
        >
          {renderContents()}
        </Flex>
      </FlexItem>
    </Flex>
  );
};

export default ApplicationsPage;

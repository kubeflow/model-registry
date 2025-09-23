import * as React from 'react';
import {
  Tabs,
  Tab,
  TabTitleText,
  PageSection,
  Flex,
  FlexItem,
  Content,
  ContentVariants,
} from '@patternfly/react-core';
import { CatalogModel, CatalogModelDetailsParams } from '~/app/modelCatalogTypes';
import ModelDetailsView from './ModelDetailsView';

// Utility function to check if a model is validated
const isModelValidated = (model: CatalogModel): boolean =>
  model.tasks?.includes('validated') ?? false;

export enum ModelDetailsTab {
  OVERVIEW = 'overview',
  PERFORMANCE_INSIGHTS = 'performance-insights',
}

export enum ModelDetailsTabTitle {
  OVERVIEW = 'Overview',
  PERFORMANCE_INSIGHTS = 'Performance insights',
}

type ModelDetailsTabsProps = {
  model: CatalogModel;
  decodedParams: CatalogModelDetailsParams;
};

const PerformanceInsightsPlaceholder = () => (
  <PageSection hasBodyWrapper={false} isFilled data-testid="performance-insights-tab-content">
    <Flex
      direction={{ default: 'column' }}
      alignItems={{ default: 'alignItemsCenter' }}
      justifyContent={{ default: 'justifyContentCenter' }}
      style={{ minHeight: '400px' }}
    >
      <FlexItem>
        <Content component={ContentVariants.p}>Performance Insights - Coming Soon</Content>
      </FlexItem>
    </Flex>
  </PageSection>
);

const ModelDetailsTabs = ({ model, decodedParams }: ModelDetailsTabsProps): React.JSX.Element => {
  const [activeTabKey, setActiveTabKey] = React.useState<ModelDetailsTab>(ModelDetailsTab.OVERVIEW);
  const isValidated = isModelValidated(model);

  const handleTabClick = (
    _event: React.MouseEvent<HTMLElement, MouseEvent>,
    tabIndex: string | number,
  ) => {
    const validTab = Object.values(ModelDetailsTab).find((tab) => tab === tabIndex);
    if (validTab) {
      setActiveTabKey(validTab);
    }
  };

  // If model is not validated, just show the overview content without tabs
  if (!isValidated) {
    return (
      <PageSection hasBodyWrapper={false} isFilled data-testid="model-overview-tab-content">
        <ModelDetailsView model={model} decodedParams={decodedParams} />
      </PageSection>
    );
  }

  return (
    <Tabs
      activeKey={activeTabKey}
      onSelect={handleTabClick}
      aria-label="Model details page tabs"
      role="region"
      data-testid="model-details-page-tabs"
    >
      <Tab
        eventKey={ModelDetailsTab.OVERVIEW}
        title={<TabTitleText>{ModelDetailsTabTitle.OVERVIEW}</TabTitleText>}
        aria-label="Model overview tab"
        data-testid="model-overview-tab"
      >
        <PageSection hasBodyWrapper={false} isFilled data-testid="model-overview-tab-content">
          <ModelDetailsView model={model} decodedParams={decodedParams} />
        </PageSection>
      </Tab>
      <Tab
        eventKey={ModelDetailsTab.PERFORMANCE_INSIGHTS}
        title={<TabTitleText>{ModelDetailsTabTitle.PERFORMANCE_INSIGHTS}</TabTitleText>}
        aria-label="Performance insights tab"
        data-testid="performance-insights-tab"
      >
        <PerformanceInsightsPlaceholder />
      </Tab>
    </Tabs>
  );
};

export default ModelDetailsTabs;

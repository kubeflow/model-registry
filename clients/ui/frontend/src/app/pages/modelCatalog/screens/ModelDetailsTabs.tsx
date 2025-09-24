import * as React from 'react';
import { Tabs, Tab, TabTitleText, PageSection } from '@patternfly/react-core';
import { CatalogModel, CatalogModelDetailsParams } from '~/app/modelCatalogTypes';
import { getLabels } from '~/app/pages/modelRegistry/screens/utils';
import ModelDetailsView from './ModelDetailsView';
import PerformanceInsightsView from './PerformanceInsightsView';

// Utility function to check if a model is validated
const isModelValidated = (model: CatalogModel): boolean => {
  if (!model.customProperties) {
    return false;
  }
  const labels = getLabels(model.customProperties);
  return labels.includes('validated');
};

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
        <PageSection hasBodyWrapper={false} isFilled data-testid="performance-insights-tab-content">
          <PerformanceInsightsView />
        </PageSection>
      </Tab>
    </Tabs>
  );
};

export default ModelDetailsTabs;

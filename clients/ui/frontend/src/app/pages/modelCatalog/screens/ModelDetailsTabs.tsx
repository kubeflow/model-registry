import * as React from 'react';
import { Tabs, Tab, TabTitleText, PageSection } from '@patternfly/react-core';
import { useNavigate } from 'react-router-dom';
import { CatalogArtifactList, CatalogModel } from '~/app/modelCatalogTypes';
import { isModelValidated } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import ModelDetailsView from './ModelDetailsView';
import PerformanceInsightsView from './PerformanceInsightsView';

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
  tab: ModelDetailsTab;
  artifacts: CatalogArtifactList;
  artifactLoaded: boolean;
  artifactsLoadError: Error | undefined;
};

const ModelDetailsTabs = ({
  model,
  tab,
  artifacts,
  artifactLoaded,
  artifactsLoadError,
}: ModelDetailsTabsProps): React.JSX.Element => {
  const isValidated = isModelValidated(model);
  const navigate = useNavigate();

  const handleTabClick = (
    _event: React.MouseEvent<HTMLElement, MouseEvent>,
    tabIndex: string | number,
  ) => {
    const validTab = Object.values(ModelDetailsTab).find((t) => t === tabIndex);
    if (validTab) {
      navigate(`../${validTab}`, { relative: 'path' });
    }
  };

  // If model is not validated, just show the overview content without tabs
  if (!isValidated) {
    return (
      <PageSection hasBodyWrapper={false} isFilled data-testid="model-overview-tab-content">
        <ModelDetailsView
          model={model}
          artifacts={artifacts}
          artifactLoaded={artifactLoaded}
          artifactsLoadError={artifactsLoadError}
        />
      </PageSection>
    );
  }

  return (
    <Tabs
      activeKey={tab}
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
          <ModelDetailsView
            model={model}
            artifacts={artifacts}
            artifactLoaded={artifactLoaded}
            artifactsLoadError={artifactsLoadError}
          />
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

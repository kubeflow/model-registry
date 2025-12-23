import * as React from 'react';
import { Tabs, Tab, TabTitleText, PageSection } from '@patternfly/react-core';
import { useNavigate, useParams } from 'react-router-dom';
import {
  CatalogArtifactList,
  CatalogArtifactType,
  CatalogModel,
  CatalogModelDetailsParams,
  CatalogPerformanceMetricsArtifact,
  MetricsType,
} from '~/app/modelCatalogTypes';
import {
  shouldShowValidatedInsights,
  filterArtifactsByType,
  decodeParams,
  getActiveLatencyFieldName,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { ModelDetailsTab, ModelCatalogNumberFilterKey } from '~/concepts/modelCatalog/const';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { useCatalogPerformanceArtifacts } from '~/app/hooks/modelCatalog/useCatalogPerformanceArtifacts';
import ModelDetailsView from './ModelDetailsView';
import PerformanceInsightsView from './PerformanceInsightsView';

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
  const navigate = useNavigate();
  const params = useParams<CatalogModelDetailsParams>();
  const decodedParams = decodeParams(params);
  const { filterData, filterOptions } = React.useContext(ModelCatalogContext);

  // Get performance-specific filter params for the /performance_artifacts endpoint
  const targetRPS = filterData[ModelCatalogNumberFilterKey.MIN_RPS];
  const latencyProperty = getActiveLatencyFieldName(filterData);

  // Check if this is a validated model that needs performance insights
  const showValidatedInsights = shouldShowValidatedInsights(model, artifacts.items);

  // Only fetch from performance artifacts endpoint when on performance insights tab
  const shouldFetchPerformanceArtifacts =
    showValidatedInsights && tab === ModelDetailsTab.PERFORMANCE_INSIGHTS;

  // Fetch performance artifacts from server with filtering/sorting/pagination
  const [performanceArtifactsList, performanceArtifactsLoaded, performanceArtifactsError] =
    useCatalogPerformanceArtifacts(
      decodedParams.sourceId || '',
      encodeURIComponent(`${decodedParams.modelName}`),
      {
        targetRPS,
        latencyProperty,
        recommendations: true,
      },
      filterData,
      filterOptions,
      shouldFetchPerformanceArtifacts,
    );

  // Extract performance artifacts from the list
  const performanceArtifacts = React.useMemo(
    () =>
      filterArtifactsByType<CatalogPerformanceMetricsArtifact>(
        performanceArtifactsList.items,
        CatalogArtifactType.metricsArtifact,
        MetricsType.performanceMetrics,
      ),
    [performanceArtifactsList.items],
  );

  // Fallback: use locally filtered artifacts if server fetch is not active
  const localPerformanceArtifacts = React.useMemo(
    () =>
      filterArtifactsByType<CatalogPerformanceMetricsArtifact>(
        artifacts.items,
        CatalogArtifactType.metricsArtifact,
        MetricsType.performanceMetrics,
      ),
    [artifacts.items],
  );

  // Use server-fetched artifacts when available, otherwise fallback to local
  const displayPerformanceArtifacts = shouldFetchPerformanceArtifacts
    ? performanceArtifacts
    : localPerformanceArtifacts;
  const isPerformanceLoading = shouldFetchPerformanceArtifacts
    ? !performanceArtifactsLoaded
    : !artifactLoaded;

  const handleTabClick = (
    _event: React.MouseEvent<HTMLElement, MouseEvent>,
    tabIndex: string | number,
  ) => {
    const validTab = Object.values(ModelDetailsTab).find((t) => t === tabIndex);
    if (validTab) {
      navigate(`../${validTab}`, { relative: 'path' });
    }
  };

  if (!showValidatedInsights) {
    return (
      <PageSection
        hasBodyWrapper={false}
        isFilled
        data-testid="model-overview-tab-content"
        padding={{ default: 'noPadding' }}
      >
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
        <PageSection
          hasBodyWrapper={false}
          isFilled
          data-testid="model-overview-tab-content"
          padding={{ default: 'noPadding' }}
        >
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
        <PageSection
          hasBodyWrapper={false}
          isFilled
          data-testid="performance-insights-tab-content"
          padding={{ default: 'noPadding' }}
        >
          <PerformanceInsightsView
            performanceArtifacts={displayPerformanceArtifacts}
            isLoading={isPerformanceLoading}
            loadError={performanceArtifactsError}
          />
        </PageSection>
      </Tab>
    </Tabs>
  );
};

export default ModelDetailsTabs;

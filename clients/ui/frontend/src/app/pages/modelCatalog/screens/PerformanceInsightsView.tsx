import * as React from 'react';
import { useLocation } from 'react-router-dom';
import { PageSection, Card, CardBody, Title, Flex, FlexItem } from '@patternfly/react-core';
import HardwareConfigurationTable from '~/app/pages/modelCatalog/components/HardwareConfigurationTable';
import { CatalogPerformanceMetricsArtifact } from '~/app/modelCatalogTypes';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import {
  hasPerformanceFiltersApplied,
  deepEqual,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { ModelDetailsTab } from './ModelDetailsTabs';

type PerformanceInsightsViewProps = {
  performanceArtifacts: CatalogPerformanceMetricsArtifact[];
  isLoading?: boolean;
};

const PerformanceInsightsView = ({
  performanceArtifacts,
  isLoading = false,
}: PerformanceInsightsViewProps): React.JSX.Element => {
  const location = useLocation();
  const { filterData, setPerformanceFiltersChangedOnDetailsPage } =
    React.useContext(ModelCatalogContext);
  const initialFilterDataRef = React.useRef<typeof filterData | null>(null);
  const prevFilterDataRef = React.useRef(filterData);
  const prevLocationRef = React.useRef(location.pathname);

  React.useEffect(() => {
    const isOnPerformanceInsightsTab = location.pathname.includes(
      ModelDetailsTab.PERFORMANCE_INSIGHTS,
    );
    const wasOnPerformanceInsightsTab = prevLocationRef.current.includes(
      ModelDetailsTab.PERFORMANCE_INSIGHTS,
    );

    if (isOnPerformanceInsightsTab && !wasOnPerformanceInsightsTab) {
      initialFilterDataRef.current = JSON.parse(JSON.stringify(filterData));
      prevFilterDataRef.current = filterData;
      setPerformanceFiltersChangedOnDetailsPage(false);
    }

    if (!isOnPerformanceInsightsTab && wasOnPerformanceInsightsTab) {
      initialFilterDataRef.current = null;
    }

    prevLocationRef.current = location.pathname;
  }, [location.pathname, filterData, setPerformanceFiltersChangedOnDetailsPage]);

  React.useEffect(() => {
    const isOnPerformanceInsightsTab = location.pathname.includes(
      ModelDetailsTab.PERFORMANCE_INSIGHTS,
    );

    if (isOnPerformanceInsightsTab && initialFilterDataRef.current === null) {
      initialFilterDataRef.current = JSON.parse(JSON.stringify(filterData));
      prevFilterDataRef.current = filterData;
    }
  }, [location.pathname, filterData]);

  React.useEffect(() => {
    const isOnPerformanceInsightsTab = location.pathname.includes(
      ModelDetailsTab.PERFORMANCE_INSIGHTS,
    );

    if (!isOnPerformanceInsightsTab || !initialFilterDataRef.current) {
      return;
    }

    const prevFilters = prevFilterDataRef.current;
    const initialFilters = initialFilterDataRef.current;
    const filtersChanged = !deepEqual(prevFilters, filterData);
    const changedFromInitial = !deepEqual(initialFilters, filterData);

    if (filtersChanged && changedFromInitial && hasPerformanceFiltersApplied(filterData)) {
      setPerformanceFiltersChangedOnDetailsPage(true);
    }

    prevFilterDataRef.current = filterData;
  }, [filterData, location.pathname, setPerformanceFiltersChangedOnDetailsPage]);

  return (
    <PageSection padding={{ default: 'noPadding' }}>
      <Card>
        <CardBody>
          <Flex direction={{ default: 'column' }} gap={{ default: 'gapLg' }}>
            <FlexItem>
              <Flex direction={{ default: 'column' }} gap={{ default: 'gapSm' }}>
                <FlexItem>
                  <Title headingLevel="h2" size="lg">
                    Hardware Configuration
                  </Title>
                </FlexItem>
                <FlexItem>
                  <p>
                    Compare the performance metrics of hardware configuration to determine the most
                    suitable option for deployment.
                  </p>
                </FlexItem>
              </Flex>
            </FlexItem>
            <FlexItem>
              <HardwareConfigurationTable
                performanceArtifacts={performanceArtifacts}
                isLoading={isLoading}
              />
            </FlexItem>
          </Flex>
        </CardBody>
      </Card>
    </PageSection>
  );
};

export default PerformanceInsightsView;

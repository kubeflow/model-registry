import * as React from 'react';
import { PageSection, Card, CardBody, Title, Flex, FlexItem } from '@patternfly/react-core';
import HardwareConfigurationTable from '~/app/pages/modelCatalog/components/HardwareConfigurationTable';
import { CatalogPerformanceMetricsArtifact } from '~/app/modelCatalogTypes';

type PerformanceInsightsViewProps = {
  performanceArtifacts: CatalogPerformanceMetricsArtifact[];
  isLoading?: boolean;
};

const PerformanceInsightsView = ({
  performanceArtifacts,
  isLoading = false,
}: PerformanceInsightsViewProps): React.JSX.Element => (
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

export default PerformanceInsightsView;

import * as React from 'react';
import { PageSection, Card, CardBody, Title, Flex, FlexItem } from '@patternfly/react-core';
import HardwareConfigurationTable from '~/app/pages/modelCatalog/components/HardwareConfigurationTable';
import { mockPerformanceMetricsArtifacts } from '~/app/pages/modelCatalog/mocks/hardwareConfigurationMock';

const PerformanceInsightsView = (): React.JSX.Element => {
  const [configurations] = React.useState(mockPerformanceMetricsArtifacts);
  const [isLoading] = React.useState(false);

  return (
    <PageSection>
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
                performanceArtifacts={configurations}
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

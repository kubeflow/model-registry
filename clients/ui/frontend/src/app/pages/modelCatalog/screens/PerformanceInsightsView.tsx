import * as React from 'react';
import { PageSection, Title, Flex, FlexItem } from '@patternfly/react-core';
import HardwareConfigurationTable from '~/app/pages/modelCatalog/components/HardwareConfigurationTable';
import { mockHardwareConfigurations } from '~/app/pages/modelCatalog/mocks/hardwareConfigurationMock';

const PerformanceInsightsView = (): React.JSX.Element => {
  const [configurations] = React.useState(mockHardwareConfigurations);
  const [isLoading] = React.useState(false);

  return (
    <PageSection style={{ overflow: 'hidden', maxWidth: '100vw', width: '100%' }}>
      <style>
        {`
          .performance-insights-container {
            width: 100%;
            max-width: 100%;
            overflow-x: hidden;
          }
        `}
      </style>
      <div className="performance-insights-container">
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
          <FlexItem style={{ width: '100%', overflow: 'hidden' }}>
            <HardwareConfigurationTable configurations={configurations} isLoading={isLoading} />
          </FlexItem>
        </Flex>
      </div>
    </PageSection>
  );
};

export default PerformanceInsightsView;

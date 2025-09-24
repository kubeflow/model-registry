import * as React from 'react';
import { PageSection, Title, Flex, FlexItem } from '@patternfly/react-core';
import HardwareConfigurationTable from '~/app/pages/modelCatalog/components/HardwareConfigurationTable';
import { mockHardwareConfigurations } from '~/app/pages/modelCatalog/mocks/hardwareConfigurationMock';

const PerformanceInsightsView = (): React.JSX.Element => {
  const [configurations] = React.useState(mockHardwareConfigurations);
  const [isLoading] = React.useState(false);

  return (
    <PageSection>
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
          <HardwareConfigurationTable configurations={configurations} isLoading={isLoading} />
        </FlexItem>
      </Flex>
    </PageSection>
  );
};

export default PerformanceInsightsView;

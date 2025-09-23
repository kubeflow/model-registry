import * as React from 'react';
import { Flex, FlexItem, Content, ContentVariants } from '@patternfly/react-core';

const PerformanceInsightsView: React.FC = () => (
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
);

export default PerformanceInsightsView;

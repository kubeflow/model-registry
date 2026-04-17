import * as React from 'react';
import {
  Stack,
  StackItem,
  Switch,
  Content,
  ContentVariants,
  Card,
  CardBody,
  Flex,
  FlexItem,
} from '@patternfly/react-core';
import { ChartBarIcon } from '@patternfly/react-icons';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import './ModelPerformanceViewToggleCard.scss';

const ModelPerformanceViewToggleCard: React.FC = () => {
  const { performanceViewEnabled, setPerformanceViewEnabled, filterOptionsLoaded } =
    React.useContext(ModelCatalogContext);

  return (
    <Card style={{ minWidth: '280px' }}>
      <CardBody>
        <Stack hasGutter>
          <StackItem>
            <Flex
              alignItems={{ default: 'alignItemsCenter' }}
              spaceItems={{ default: 'spaceItemsXs' }}
            >
              <FlexItem>
                <ChartBarIcon color="var(--pf-t--global--icon--color--status--info--default)" />
              </FlexItem>
              <FlexItem>
                <Switch
                  id="model-performance-view-toggle"
                  label="Model performance view"
                  isChecked={performanceViewEnabled}
                  isReversed
                  isDisabled={!filterOptionsLoaded}
                  onChange={(_event, checked) => setPerformanceViewEnabled(checked)}
                  data-testid="model-performance-view-toggle"
                  className="model-performance-switch"
                />
              </FlexItem>
            </Flex>
          </StackItem>
          <StackItem>
            <Content component={ContentVariants.small}>
              Enable performance filters, display model benchmark data, and exclude unvalidated
              models.
            </Content>
          </StackItem>
        </Stack>
      </CardBody>
    </Card>
  );
};

export default ModelPerformanceViewToggleCard;

import * as React from 'react';
import {
  Stack,
  StackItem,
  Switch,
  Content,
  ContentVariants,
  Card,
  CardBody,
} from '@patternfly/react-core';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';

const ModelPerformanceViewToggleCard: React.FC = () => {
  const { performanceViewEnabled, setPerformanceViewEnabled, filterOptionsLoaded } =
    React.useContext(ModelCatalogContext);

  return (
    <Card>
      <CardBody>
        <Stack hasGutter>
          <StackItem>
            <Switch
              id="model-performance-view-toggle"
              label="Model performance view"
              isChecked={performanceViewEnabled}
              isDisabled={!filterOptionsLoaded}
              onChange={(_event, checked) => setPerformanceViewEnabled(checked)}
              data-testid="model-performance-view-toggle"
            />
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

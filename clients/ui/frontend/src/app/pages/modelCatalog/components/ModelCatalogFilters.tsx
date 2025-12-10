import * as React from 'react';
import {
  Stack,
  StackItem,
  Spinner,
  Alert,
  Switch,
  Content,
  ContentVariants,
  Card,
  CardBody,
} from '@patternfly/react-core';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { ModelCatalogStringFilterKey } from '~/concepts/modelCatalog/const';
import TaskFilter from './globalFilters/TaskFilter';
import ProviderFilter from './globalFilters/ProviderFilter';
import LicenseFilter from './globalFilters/LicenseFilter';
import LanguageFilter from './globalFilters/LanguageFilter';

const ModelCatalogFilters: React.FC = () => {
  const {
    filterOptions,
    filterOptionsLoaded,
    filterOptionsLoadError,
    performanceViewEnabled,
    setPerformanceViewEnabled,
  } = React.useContext(ModelCatalogContext);
  const filters = filterOptions?.filters;
  if (!filterOptionsLoaded) {
    return <Spinner />;
  }
  if (filterOptionsLoadError) {
    return (
      <Alert variant="danger" title="Failed to load filter options" isInline>
        {filterOptionsLoadError.message}
      </Alert>
    );
  }
  return (
    <Stack hasGutter>
      <StackItem>
        <Card>
          <CardBody>
            <Stack hasGutter>
              <StackItem>
                <Switch
                  id="model-performance-view-toggle"
                  label="Model performance view"
                  isChecked={performanceViewEnabled}
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
      </StackItem>
      <TaskFilter
        filters={filters && ModelCatalogStringFilterKey.TASK in filters ? filters : undefined}
      />
      <ProviderFilter
        filters={filters && ModelCatalogStringFilterKey.PROVIDER in filters ? filters : undefined}
      />
      <LicenseFilter
        filters={filters && ModelCatalogStringFilterKey.LICENSE in filters ? filters : undefined}
      />
      <LanguageFilter
        filters={filters && ModelCatalogStringFilterKey.LANGUAGE in filters ? filters : undefined}
      />
    </Stack>
  );
};

export default ModelCatalogFilters;

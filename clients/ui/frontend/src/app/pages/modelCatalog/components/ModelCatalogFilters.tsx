import * as React from 'react';
import { Stack, StackItem, Divider, Spinner, Alert } from '@patternfly/react-core';
import { useCatalogFilterOptionList } from '~/app/hooks/modelCatalog/useCatalogFilterOptionList';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { ModelCatalogStringFilterKey } from '~/concepts/modelCatalog/const';
import TaskFilter from './globalFilters/TaskFilter';
import ProviderFilter from './globalFilters/ProviderFilter';
import LicenseFilter from './globalFilters/LicenseFilter';
import LanguageFilter from './globalFilters/LanguageFilter';

const ModelCatalogFilters: React.FC = () => {
  const { apiState } = React.useContext(ModelCatalogContext);
  const [filterOptions, loaded, error] = useCatalogFilterOptionList(apiState);
  const filters = filterOptions?.filters;
  if (!loaded) {
    return <Spinner />;
  }
  if (error) {
    return (
      <Alert variant="danger" title="Failed to load filter options" isInline>
        {error.message}
      </Alert>
    );
  }
  return (
    <Stack hasGutter>
      <StackItem>
        <TaskFilter
          filters={filters && ModelCatalogStringFilterKey.TASK in filters ? filters : undefined}
        />
      </StackItem>
      <Divider />
      <StackItem>
        <ProviderFilter
          filters={filters && ModelCatalogStringFilterKey.PROVIDER in filters ? filters : undefined}
        />
      </StackItem>
      <Divider />
      <StackItem>
        <LicenseFilter
          filters={filters && ModelCatalogStringFilterKey.LICENSE in filters ? filters : undefined}
        />
      </StackItem>
      <Divider />
      <StackItem>
        <LanguageFilter
          filters={filters && ModelCatalogStringFilterKey.LANGUAGE in filters ? filters : undefined}
        />
      </StackItem>
    </Stack>
  );
};

export default ModelCatalogFilters;

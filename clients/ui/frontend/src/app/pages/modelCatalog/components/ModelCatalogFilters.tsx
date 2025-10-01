import * as React from 'react';
import { Stack, StackItem, Divider } from '@patternfly/react-core';
import { useCatalogFilterOptionList } from '~/app/hooks/modelCatalog/useCatalogFilterOptionList';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import TaskFilter from './globalFilters/TaskFilter';
import ProviderFilter from './globalFilters/ProviderFilter';
import LicenseFilter from './globalFilters/LicenseFilter';
import LanguageFilter from './globalFilters/LanguageFilter';

const ModelCatalogFilters: React.FC = () => {
  const { apiState } = React.useContext(ModelCatalogContext);
  const [filterOptions] = useCatalogFilterOptionList(apiState);
  const filters = filterOptions?.filters;
  return (
    <Stack hasGutter>
      <StackItem>
        <TaskFilter filters={filters} />
      </StackItem>
      <Divider />
      <StackItem>
        <ProviderFilter filters={filters} />
      </StackItem>
      <Divider />
      <StackItem>
        <LicenseFilter filters={filters} />
      </StackItem>
      <Divider />
      <StackItem>
        <LanguageFilter filters={filters} />
      </StackItem>
    </Stack>
  );
};

export default ModelCatalogFilters;

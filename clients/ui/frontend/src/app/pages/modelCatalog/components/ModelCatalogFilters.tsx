import * as React from 'react';
import { Stack, StackItem, Divider } from '@patternfly/react-core';
import { getModelCatalogFilters } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import TaskFilter from './globalFilters/TaskFilter';
import ProviderFilter from './globalFilters/ProviderFilter';
import LicenseFilter from './globalFilters/LicenseFilter';
import LanguageFilter from './globalFilters/LanguageFilter';

const ModelCatalogFilters: React.FC = () => {
  const { filters } = getModelCatalogFilters();
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

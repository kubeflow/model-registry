import * as React from 'react';
import { Stack, StackItem, Divider } from '@patternfly/react-core';
import { useCatalogFilterOptionList } from '~/app/hooks/modelCatalog/useCatalogFilterOptionList';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { ModelCatalogFilterKeys } from '~/concepts/modelCatalog/const';
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
        <TaskFilter
          filters={filters && ModelCatalogFilterKeys.TASK in filters ? filters : undefined}
        />
      </StackItem>
      <Divider />
      <StackItem>
        <ProviderFilter
          filters={filters && ModelCatalogFilterKeys.PROVIDER in filters ? filters : undefined}
        />
      </StackItem>
      <Divider />
      <StackItem>
        <LicenseFilter
          filters={filters && ModelCatalogFilterKeys.LICENSE in filters ? filters : undefined}
        />
      </StackItem>
      <Divider />
      <StackItem>
        <LanguageFilter
          filters={filters && ModelCatalogFilterKeys.LANGUAGE in filters ? filters : undefined}
        />
      </StackItem>
    </Stack>
  );
};

export default ModelCatalogFilters;

import * as React from 'react';
import { Divider, StackItem } from '@patternfly/react-core';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';
import {
  ModelCatalogStringFilterKey,
  MODEL_CATALOG_PROVIDER_NAME_MAPPING,
} from '~/concepts/modelCatalog/const';
import { CatalogFilterOptions, ModelCatalogStringFilterOptions } from '~/app/modelCatalogTypes';

const filterKey = ModelCatalogStringFilterKey.PROVIDER;

type ProviderFilterProps = {
  filters?: Extract<CatalogFilterOptions, Partial<ModelCatalogStringFilterOptions>>;
};

const ProviderFilter: React.FC<ProviderFilterProps> = ({ filters }) => {
  const provider = filters?.[filterKey];

  if (!provider) {
    return null;
  }

  return (
    <>
      <StackItem>
        <ModelCatalogStringFilter<ModelCatalogStringFilterKey.PROVIDER>
          title="Provider"
          filterKey={filterKey}
          filterToNameMapping={MODEL_CATALOG_PROVIDER_NAME_MAPPING}
          filters={provider}
        />
      </StackItem>
      <Divider />
    </>
  );
};

export default ProviderFilter;

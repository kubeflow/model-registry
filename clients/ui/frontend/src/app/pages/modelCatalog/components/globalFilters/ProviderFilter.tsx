import * as React from 'react';
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
    <ModelCatalogStringFilter<ModelCatalogStringFilterKey.PROVIDER>
      title="Provider"
      filterKey={filterKey}
      filterToNameMapping={MODEL_CATALOG_PROVIDER_NAME_MAPPING}
      filters={provider}
    />
  );
};

export default ProviderFilter;

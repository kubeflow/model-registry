import * as React from 'react';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';
import {
  ModelCatalogFilterKeys,
  MODEL_CATALOG_PROVIDER_NAME_MAPPING,
} from '~/concepts/modelCatalog/const';
import { CatalogFilterOptionsList, ModelCatalogFilterTypesByKey } from '~/app/modelCatalogTypes';

const filterKey = ModelCatalogFilterKeys.PROVIDER;

type ProviderFilterProps = {
  filters?: Extract<CatalogFilterOptionsList['filters'], Partial<ModelCatalogFilterTypesByKey>>;
};

const ProviderFilter: React.FC<ProviderFilterProps> = ({ filters }) => {
  const provider = filters?.[filterKey];

  if (!provider) {
    return null;
  }

  return (
    <ModelCatalogStringFilter<ModelCatalogFilterKeys.PROVIDER>
      title="Provider"
      filterKey={filterKey}
      filterToNameMapping={MODEL_CATALOG_PROVIDER_NAME_MAPPING}
      filters={provider}
    />
  );
};

export default ProviderFilter;

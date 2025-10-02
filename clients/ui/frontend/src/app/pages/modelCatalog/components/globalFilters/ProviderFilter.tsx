import * as React from 'react';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';
import {
  ModelCatalogFilterKey,
  MODEL_CATALOG_PROVIDER_NAME_MAPPING,
} from '~/concepts/modelCatalog/const';
import { CatalogFilterOptionsList, GlobalFilterTypes } from '~/app/modelCatalogTypes';

const filterKey = ModelCatalogFilterKey.PROVIDER;

type ProviderFilterProps = {
  filters?: Extract<CatalogFilterOptionsList['filters'], Partial<GlobalFilterTypes>>;
};

const ProviderFilter: React.FC<ProviderFilterProps> = ({ filters }) => {
  const provider = filters?.[filterKey];

  if (!provider) {
    return null;
  }

  return (
    <ModelCatalogStringFilter<ModelCatalogFilterKey.PROVIDER>
      title="Provider"
      filterKey={filterKey}
      filterToNameMapping={MODEL_CATALOG_PROVIDER_NAME_MAPPING}
      filters={provider}
    />
  );
};

export default ProviderFilter;

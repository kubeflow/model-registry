import * as React from 'react';
import {
  ModelCatalogFilterResponseType,
  ModelCatalogFilterStatesByKey,
} from '~/app/pages/modelCatalog/types';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import {
  ModelCatalogFilterKeys,
  MODEL_CATALOG_PROVIDER_NAME_MAPPING,
} from '~/concepts/modelCatalog/const';

const filterKey = ModelCatalogFilterKeys.PROVIDER;

type ProviderFilterProps = {
  filters: ModelCatalogFilterResponseType['filters'];
};

const ProviderFilter: React.FC<ProviderFilterProps> = ({ filters }) => {
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
  const provider = filters[filterKey];

  React.useEffect(() => {
    if (provider && !(filterKey in filterData)) {
      const state: ModelCatalogFilterStatesByKey[typeof filterKey] = {};
      provider.values.forEach((key) => {
        state[key] = false;
      });
      setFilterData(filterKey, state);
    }
  }, [provider, filterData, setFilterData]);

  if (!provider) {
    return null;
  }

  return (
    <ModelCatalogStringFilter<ModelCatalogFilterKeys.PROVIDER>
      title="Provider"
      filterToNameMapping={MODEL_CATALOG_PROVIDER_NAME_MAPPING}
      filters={provider}
      data={filterData[filterKey]}
      setData={(state) => setFilterData(filterKey, state)}
    />
  );
};

export default ProviderFilter;

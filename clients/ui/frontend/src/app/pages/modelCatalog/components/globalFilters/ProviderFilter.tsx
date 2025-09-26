import * as React from 'react';
import {
  ModelCatalogStringFilterStateType,
  ModelCatalogFilterResponseType,
} from '~/app/pages/modelCatalog/types';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { MODEL_CATALOG_PROVIDER_NAME_MAPPING } from '~/concepts/modelCatalog/const';

type ProviderFilterProps = {
  filters: ModelCatalogFilterResponseType['filters'];
};

const ProviderFilter: React.FC<ProviderFilterProps> = ({ filters }) => {
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
  const { provider } = filters;

  React.useEffect(() => {
    if (provider && !('provider' in filterData)) {
      const state: Record<string, boolean> = {};
      provider.values.forEach((key) => {
        state[key] = false;
      });
      setFilterData('provider', state);
    }
  }, [provider, filterData, setFilterData]);

  if (!provider) {
    return null;
  }

  return (
    <ModelCatalogStringFilter
      title="Provider"
      filterKey="provider"
      filterToNameMapping={MODEL_CATALOG_PROVIDER_NAME_MAPPING}
      filters={provider}
      data={filterData}
      setData={(state: ModelCatalogStringFilterStateType) => setFilterData('provider', state)}
    />
  );
};

export default ProviderFilter;

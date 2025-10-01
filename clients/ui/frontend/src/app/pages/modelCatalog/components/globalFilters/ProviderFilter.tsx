import * as React from 'react';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import {
  ModelCatalogFilterKeys,
  MODEL_CATALOG_PROVIDER_NAME_MAPPING,
} from '~/concepts/modelCatalog/const';
import {
  CatalogFilterOptionsList,
  ModelCatalogFilterTypesByKey,
  ModelCatalogProvidersFilterStateType,
} from '~/app/modelCatalogTypes';

const filterKey = ModelCatalogFilterKeys.PROVIDER;

type ProviderFilterProps = {
  filters?: Extract<CatalogFilterOptionsList['filters'], Partial<ModelCatalogFilterTypesByKey>>;
};

const ProviderFilter: React.FC<ProviderFilterProps> = ({ filters }) => {
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
  const provider = filters?.[filterKey];
  const currentState = filterData[filterKey];

  React.useEffect(() => {
    if (!provider) {
      return;
    }

    const filterKeys = provider.values;
    const hasMatchingKeys =
      currentState !== undefined &&
      filterKeys.length === Object.keys(currentState).length &&
      filterKeys.every((key) => key in currentState);

    if (hasMatchingKeys) {
      return;
    }

    const nextState: ModelCatalogProvidersFilterStateType = {};
    filterKeys.forEach((key) => {
      nextState[key] = currentState?.[key] ?? false;
    });

    setFilterData(filterKey, nextState);
  }, [provider, currentState, setFilterData]);

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

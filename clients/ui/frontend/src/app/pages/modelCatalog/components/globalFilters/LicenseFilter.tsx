import * as React from 'react';
import {
  ModelCatalogFilterResponseType,
  ModelCatalogFilterStatesByKey,
} from '~/app/pages/modelCatalog/types';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import {
  ModelCatalogFilterKeys,
  MODEL_CATALOG_LICENSE_NAME_MAPPING,
} from '~/concepts/modelCatalog/const';

const filterKey = ModelCatalogFilterKeys.LICENSE;

type LicenseFilterProps = {
  filters: ModelCatalogFilterResponseType['filters'];
};

const LicenseFilter: React.FC<LicenseFilterProps> = ({ filters }) => {
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
  const license = filters[filterKey];

  React.useEffect(() => {
    if (license && !(filterKey in filterData)) {
      const state: ModelCatalogFilterStatesByKey[typeof filterKey] = {};
      license.values.forEach((key) => {
        state[key] = false;
      });
      setFilterData(filterKey, state);
    }
  }, [license, filterData, setFilterData]);

  if (!license) {
    return null;
  }

  return (
    <ModelCatalogStringFilter<ModelCatalogFilterKeys.LICENSE>
      title="License"
      filterToNameMapping={MODEL_CATALOG_LICENSE_NAME_MAPPING}
      filters={license}
      data={filterData[filterKey]}
      setData={(state) => setFilterData(filterKey, state)}
    />
  );
};

export default LicenseFilter;

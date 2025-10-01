import * as React from 'react';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import {
  ModelCatalogFilterKeys,
  MODEL_CATALOG_LICENSE_NAME_MAPPING,
} from '~/concepts/modelCatalog/const';
import {
  CatalogFilterOptionsList,
  ModelCatalogLicensesFilterStateType,
} from '~/app/modelCatalogTypes';

const filterKey = ModelCatalogFilterKeys.LICENSE;

type LicenseFilterProps = {
  filters?: CatalogFilterOptionsList['filters'];
};

const LicenseFilter: React.FC<LicenseFilterProps> = ({ filters }) => {
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
  const license = filters?.[filterKey];
  const currentState = filterData[filterKey];

  React.useEffect(() => {
    if (!license) {
      return;
    }

    const filterKeys = license.values;
    const hasMatchingKeys =
      currentState !== undefined &&
      filterKeys.length === Object.keys(currentState).length &&
      filterKeys.every((key) => key in currentState);

    if (hasMatchingKeys) {
      return;
    }

    const nextState: ModelCatalogLicensesFilterStateType = {};
    filterKeys.forEach((key) => {
      nextState[key] = currentState?.[key] ?? false;
    });

    setFilterData(filterKey, nextState);
  }, [license, currentState, setFilterData]);

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

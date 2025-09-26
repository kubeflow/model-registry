import * as React from 'react';
import {
  ModelCatalogStringFilterStateType,
  ModelCatalogFilterResponseType,
} from '~/app/pages/modelCatalog/types';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { MODEL_CATALOG_LICENSE_NAME_MAPPING } from '~/concepts/modelCatalog/const';

type LicenseFilterProps = {
  filters: ModelCatalogFilterResponseType['filters'];
};

const LicenseFilter: React.FC<LicenseFilterProps> = ({ filters }) => {
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
  const { license } = filters;

  React.useEffect(() => {
    if (license && !('license' in filterData)) {
      const state: Record<string, boolean> = {};
      license.values.forEach((key) => {
        state[key] = false;
      });
      setFilterData('license', state);
    }
  }, [license, filterData, setFilterData]);

  if (!license) {
    return null;
  }

  return (
    <ModelCatalogStringFilter
      title="License"
      filterKey="license"
      filterToNameMapping={MODEL_CATALOG_LICENSE_NAME_MAPPING}
      filters={license}
      data={filterData}
      setData={(state: ModelCatalogStringFilterStateType) => setFilterData('license', state)}
    />
  );
};

export default LicenseFilter;

import * as React from 'react';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';
import {
  ModelCatalogFilterKeys,
  MODEL_CATALOG_LICENSE_NAME_MAPPING,
} from '~/concepts/modelCatalog/const';
import { CatalogFilterOptionsList, ModelCatalogFilterTypesByKey } from '~/app/modelCatalogTypes';

const filterKey = ModelCatalogFilterKeys.LICENSE;

type LicenseFilterProps = {
  filters?: Extract<CatalogFilterOptionsList['filters'], Partial<ModelCatalogFilterTypesByKey>>;
};

const LicenseFilter: React.FC<LicenseFilterProps> = ({ filters }) => {
  const license = filters?.[filterKey];

  if (!license) {
    return null;
  }

  return (
    <ModelCatalogStringFilter<ModelCatalogFilterKeys.LICENSE>
      title="License"
      filterKey={filterKey}
      filterToNameMapping={MODEL_CATALOG_LICENSE_NAME_MAPPING}
      filters={license}
    />
  );
};

export default LicenseFilter;

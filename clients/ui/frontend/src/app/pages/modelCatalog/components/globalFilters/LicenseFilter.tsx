import * as React from 'react';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';
import {
  ModelCatalogFilterKey,
  MODEL_CATALOG_LICENSE_NAME_MAPPING,
} from '~/concepts/modelCatalog/const';
import { CatalogFilterOptionsList, GlobalFilterTypes } from '~/app/modelCatalogTypes';

const filterKey = ModelCatalogFilterKey.LICENSE;

type LicenseFilterProps = {
  filters?: Extract<CatalogFilterOptionsList['filters'], Partial<GlobalFilterTypes>>;
};

const LicenseFilter: React.FC<LicenseFilterProps> = ({ filters }) => {
  const license = filters?.[filterKey];

  if (!license) {
    return null;
  }

  return (
    <ModelCatalogStringFilter<ModelCatalogFilterKey.LICENSE>
      title="License"
      filterKey={filterKey}
      filterToNameMapping={MODEL_CATALOG_LICENSE_NAME_MAPPING}
      filters={license}
    />
  );
};

export default LicenseFilter;

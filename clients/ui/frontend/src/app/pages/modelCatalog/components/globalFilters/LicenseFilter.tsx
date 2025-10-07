import * as React from 'react';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';
import {
  ModelCatalogStringFilterKey,
  MODEL_CATALOG_LICENSE_NAME_MAPPING,
} from '~/concepts/modelCatalog/const';
import { CatalogFilterOptions, ModelCatalogStringFilterOptions } from '~/app/modelCatalogTypes';

const filterKey = ModelCatalogStringFilterKey.LICENSE;

type LicenseFilterProps = {
  filters?: Extract<CatalogFilterOptions, Partial<ModelCatalogStringFilterOptions>>;
};

const LicenseFilter: React.FC<LicenseFilterProps> = ({ filters }) => {
  const license = filters?.[filterKey];

  if (!license) {
    return null;
  }

  return (
    <ModelCatalogStringFilter<ModelCatalogStringFilterKey.LICENSE>
      title="License"
      filterKey={filterKey}
      filterToNameMapping={MODEL_CATALOG_LICENSE_NAME_MAPPING}
      filters={license}
    />
  );
};

export default LicenseFilter;

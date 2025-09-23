import * as React from 'react';
import { ModelCatalogFilterCategoryType } from '~/app/pages/modelCatalog/types';
import ModelCatalogStringFilter from './ModelCatalogStringFilter';

type ModelCatalogFilterCategoryProps = {
  filters: ModelCatalogFilterCategoryType;
  setData: (property: string, checked: boolean) => void;
};

const ModelCatalogFilterCategory: React.FC<ModelCatalogFilterCategoryProps> = ({
  filters,
  setData,
}) => {
  if (filters.type === 'string') {
    return <ModelCatalogStringFilter filters={filters.values} setData={setData} />;
  }
  return null;
};

export default ModelCatalogFilterCategory;

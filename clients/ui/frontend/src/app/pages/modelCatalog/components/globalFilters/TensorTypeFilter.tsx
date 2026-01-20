import React from 'react';
import { StackItem } from '@patternfly/react-core';
import { CatalogFilterOptions, ModelCatalogStringFilterOptions } from '~/app/modelCatalogTypes';
import {
  MODEL_CATALOG_TENSOR_TYPE_MAPPING,
  ModelCatalogStringFilterKey,
} from '~/concepts/modelCatalog/const';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';

const filterKey = ModelCatalogStringFilterKey.TENSOR_TYPE;

type TensorTypeFilterProps = {
  filters?: Extract<CatalogFilterOptions, Partial<ModelCatalogStringFilterOptions>>;
};

const TensorTypeFilter: React.FC<TensorTypeFilterProps> = ({ filters }) => {
  const tensorType = filters?.[filterKey];

  if (!tensorType) {
    return null;
  }

  return (
    <StackItem>
      <ModelCatalogStringFilter<ModelCatalogStringFilterKey.TENSOR_TYPE>
        title="Tensor type"
        filterKey={filterKey}
        filters={tensorType}
        filterToNameMapping={MODEL_CATALOG_TENSOR_TYPE_MAPPING}
      />
    </StackItem>
  );
};

export default TensorTypeFilter;

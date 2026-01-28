import * as React from 'react';
import { StackItem } from '@patternfly/react-core';
import { CatalogFilterOptions, ModelCatalogStringFilterOptions } from '~/app/modelCatalogTypes';
import { ModelCatalogStringFilterKey, ModelCatalogTensorType } from '~/concepts/modelCatalog/const';
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
        filterToNameMapping={ModelCatalogTensorType}
      />
    </StackItem>
  );
};

export default TensorTypeFilter;

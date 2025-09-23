import * as React from 'react';
import useGenericObjectState from 'mod-arch-core/dist/utilities/useGenericObjectState';
import type { ModelCatalogFiltersType } from '~/app/pages/modelCatalog/types';
import {
  getModelCatalogFilters,
  processModelCatalogFilters,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import ModelCatalogFilterCategory from './ModelCatalogFilterCategory';

const ModelCatalogFilters: React.FC = () => {
  const { filters } = getModelCatalogFilters();
  const [data, setData] = useGenericObjectState<ModelCatalogFiltersType>(
    processModelCatalogFilters(filters),
  );
  return (
    <div data-testid="model-catalog-filters">
      {Object.entries(data).map(([key, value]) => (
        <ModelCatalogFilterCategory
          key={key}
          filters={value}
          setData={(property, checked) => setData(key, { ...value, [property]: checked })}
        />
      ))}
    </div>
  );
};

export default ModelCatalogFilters;

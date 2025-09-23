import * as React from 'react';
import { ModelCatalogFilterCategoryType } from '~/app/pages/modelCatalog/types';

type ModelCatalogStringFilterProps = {
  filters: ModelCatalogFilterCategoryType;
  setData: (key: string, value: string) => void;
};

const ModelCatalogStringFilter: React.FC<ModelCatalogStringFilterProps> = ({
  filters,
  setData,
}) => <div>ModelCatalogStringFilter</div>;

export default ModelCatalogStringFilter;

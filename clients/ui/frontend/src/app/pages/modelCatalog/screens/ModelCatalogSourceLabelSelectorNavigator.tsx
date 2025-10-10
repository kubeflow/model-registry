import React from 'react';
import ModelCatalogSourceLabelSelector from './ModelCatalogSourceLabelSelector';

type ModelCatalogSourceLabelSelectorNavigatorProps = {
  searchTerm?: string;
  onSearch?: (term: string) => void;
  onClearSearch?: () => void;
};

const ModelCatalogSourceLabelSelectorNavigator: React.FC<
  ModelCatalogSourceLabelSelectorNavigatorProps
> = ({ searchTerm, onSearch, onClearSearch }) => (
  <ModelCatalogSourceLabelSelector
    searchTerm={searchTerm}
    onSearch={onSearch}
    onClearSearch={onClearSearch}
  />
);
export default ModelCatalogSourceLabelSelectorNavigator;

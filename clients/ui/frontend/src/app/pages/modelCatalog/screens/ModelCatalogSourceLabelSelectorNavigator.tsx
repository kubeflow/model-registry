import React from 'react';
import ModelCatalogSourceLabelSelector from './ModelCatalogSourceLabelSelector';

type ModelCatalogSourceLabelSelectorNavigatorProps = {
  searchTerm?: string;
  onSearch?: (term: string) => void;
  onClearSearch?: () => void;
  onResetAllFilters?: () => void;
};

const ModelCatalogSourceLabelSelectorNavigator: React.FC<
  ModelCatalogSourceLabelSelectorNavigatorProps
> = ({ searchTerm, onSearch, onClearSearch, onResetAllFilters }) => (
  <ModelCatalogSourceLabelSelector
    searchTerm={searchTerm}
    onSearch={onSearch}
    onClearSearch={onClearSearch}
    onResetAllFilters={onResetAllFilters}
  />
);
export default ModelCatalogSourceLabelSelectorNavigator;

import React from 'react';
import ModelCatalogSourceLabelSelector from './ModelCatalogSourceLabelSelector';

type ModelCatalogSourceLabelSelectorNavigatorProps = {
  searchTerm?: string;
  onSearch?: (term: string) => void;
  onClearSearch?: () => void;
  isPrimary?: boolean;
};

const ModelCatalogSourceLabelSelectorNavigator: React.FC<
  ModelCatalogSourceLabelSelectorNavigatorProps
> = ({ searchTerm, onSearch, onClearSearch, isPrimary }) => (
  <ModelCatalogSourceLabelSelector
    searchTerm={searchTerm}
    onSearch={onSearch}
    onClearSearch={onClearSearch}
    primary={isPrimary}
  />
);
export default ModelCatalogSourceLabelSelectorNavigator;

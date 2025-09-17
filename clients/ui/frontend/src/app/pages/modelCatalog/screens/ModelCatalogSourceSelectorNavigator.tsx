import React from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { CatalogModelDetailsParams } from '~/app/modelCatalogTypes';
import { decodeParams } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import ModelCatalogSourceSelector from './ModelCatalogSourceSelector';

type ModelCatalogSourceSelectorNavigatorProps = {
  getRedirectPath: (sourceId: string) => string;
  searchTerm?: string;
  onSearch?: (term: string) => void;
  onClearSearch?: () => void;
  isPrimary?: boolean;
};

const ModelCatalogSourceSelectorNavigator: React.FC<ModelCatalogSourceSelectorNavigatorProps> = ({
  getRedirectPath,
  searchTerm,
  onSearch,
  onClearSearch,
  isPrimary,
}) => {
  const navigate = useNavigate();
  const params = useParams<CatalogModelDetailsParams>();
  const decodedParams = decodeParams(params);

  return (
    <ModelCatalogSourceSelector
      sourceId={decodedParams.sourceId ?? ''}
      onSelection={(id) => {
        if (id !== decodedParams.sourceId) {
          navigate(getRedirectPath(id || ''));
        }
      }}
      searchTerm={searchTerm}
      onSearch={onSearch}
      onClearSearch={onClearSearch}
      primary={isPrimary}
    />
  );
};

export default ModelCatalogSourceSelectorNavigator;

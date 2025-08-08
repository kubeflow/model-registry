import * as React from 'react';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { ModelCatalogContextType } from '~/app/modelCatalogTypes';
// TODO: This is a placeholder for the model catalog sources hook once API are implemented and BFF endpoints are there
export const useModelCatalogSources = (): ModelCatalogContextType => {
  const context = React.useContext(ModelCatalogContext);

  return context;
};

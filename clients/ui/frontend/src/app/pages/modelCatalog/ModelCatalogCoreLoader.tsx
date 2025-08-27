import * as React from 'react';
import { Outlet } from 'react-router-dom';
import { ModelCatalogContextProvider } from '~/app/context/modelCatalog/ModelCatalogContext';

const ModelCatalogCoreLoader: React.FC = () => (
  <ModelCatalogContextProvider>
    <Outlet />
  </ModelCatalogContextProvider>
);

export default ModelCatalogCoreLoader;

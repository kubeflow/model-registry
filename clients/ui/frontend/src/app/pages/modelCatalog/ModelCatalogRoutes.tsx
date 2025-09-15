import * as React from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import { ModelCatalogContextProvider } from '~/app/context/modelCatalog/ModelCatalogContext';
import { modelCatalogUrl } from '~/app/routes/modelCatalog/catalogModel';
import ModelCatalogCoreLoader from './ModelCatalogCoreLoader';
import ModelDetailsPage from './screens/ModelDetailsPage';
import ModelCatalogPage from './screens/ModelCatalogPage';
import RegisterCatalogModelPage from './screens/RegisterCatalogModelPage';

const ModelCatalogRoutes: React.FC = () => (
  <ModelCatalogContextProvider>
    <Routes>
      <Route
        path="/:sourceId?/*"
        element={
          <ModelCatalogCoreLoader
            getInvalidRedirectPath={(sourceId) => modelCatalogUrl(sourceId)}
          />
        }
      >
        <Route index element={<ModelCatalogPage />} />
        <Route path=":repositoryName/:modelName">
          <Route index element={<ModelDetailsPage />} />
          <Route path="register" element={<RegisterCatalogModelPage />} />
          <Route path="*" element={<Navigate to="." />} />
        </Route>
        <Route path="*" element={<Navigate to="." />} />
      </Route>
    </Routes>
  </ModelCatalogContextProvider>
);

export default ModelCatalogRoutes;

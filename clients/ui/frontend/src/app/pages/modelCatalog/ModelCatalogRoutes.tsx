import * as React from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import ModelCatalogCoreLoader from './ModelCatalogCoreLoader';
import ModelDetailsPage from './screens/ModelDetailsPage';
import ModelCatalogPage from './screens/ModelCatalogPage';
import RegisterCatalogModelPage from './screens/RegisterCatalogModelPage';

const ModelCatalogRoutes: React.FC = () => (
  <Routes>
    <Route path="/" element={<ModelCatalogCoreLoader />}>
      <Route index element={<ModelCatalogPage />} />
      {/* TODO: keep simple modelId param for now */}
      <Route path=":modelId" element={<ModelDetailsPage />} />
      <Route path=":modelId/register" element={<RegisterCatalogModelPage />} />
      <Route path="*" element={<Navigate to="." />} />
    </Route>
  </Routes>
);

export default ModelCatalogRoutes;

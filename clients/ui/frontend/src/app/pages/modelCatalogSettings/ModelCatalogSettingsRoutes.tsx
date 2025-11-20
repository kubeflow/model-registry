import * as React from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import { ModelCatalogSettingsContextProvider } from '~/app/context/modelCatalogSettings/ModelCatalogSettingsContext';
import ModelCatalogSettings from '~/app/pages/modelCatalogSettings/screens/ModelCatalogSettings';
import ManageSourcePage from '~/app/pages/modelCatalogSettings/screens/ManageSourcePage';

const ModelCatalogSettingsRoutes: React.FC = () => (
  <ModelCatalogSettingsContextProvider>
    <Routes>
      <Route path="/" element={<ModelCatalogSettings />} />
      <Route path="add-source" element={<ManageSourcePage />} />
      <Route path="manage-source/:catalogSourceId" element={<ManageSourcePage />} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  </ModelCatalogSettingsContextProvider>
);

export default ModelCatalogSettingsRoutes;

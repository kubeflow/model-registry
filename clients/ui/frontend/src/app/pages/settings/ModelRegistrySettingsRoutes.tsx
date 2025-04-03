import * as React from 'react';
import { Navigate, Routes, Route } from 'react-router-dom';
import ModelRegistrySettings from './ModelRegistrySettings';
import ModelRegistriesPermissions from './ModelRegistriesPermissions';

const ModelRegistrySettingsRoutes: React.FC = () => (
  <Routes>
    <Route path="/" element={<ModelRegistrySettings />} />
    <Route path="*" element={<Navigate to="/" />} />
    <Route path="/permissions" element={<ModelRegistriesPermissions />} />
  </Routes>
);

export default ModelRegistrySettingsRoutes;

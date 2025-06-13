import * as React from 'react';
import { Navigate, Routes, Route } from 'react-router-dom';
import ModelRegistrySettings from './ModelRegistrySettings';
import ModelRegistriesManagePermissions from './ModelRegistriesPermissions';

const ModelRegistrySettingsRoutes: React.FC = () => (
  <Routes>
    <Route path="/" element={<ModelRegistrySettings />} />
    <Route path="permissions/:mrName" element={<ModelRegistriesManagePermissions />} />
    <Route path="*" element={<Navigate to="/" />} />
  </Routes>
);

export default ModelRegistrySettingsRoutes; 
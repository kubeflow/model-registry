import * as React from 'react';
import { Navigate, Routes, Route } from 'react-router-dom';
import ModelRegistrySettings from './ModelRegistrySettings';

const ModelRegistrySettingsRoutes: React.FC = () => (
  <Routes>
    <Route path="/" element={<ModelRegistrySettings />} />
    <Route path="*" element={<Navigate to="/" />} />
  </Routes>
);

export default ModelRegistrySettingsRoutes;

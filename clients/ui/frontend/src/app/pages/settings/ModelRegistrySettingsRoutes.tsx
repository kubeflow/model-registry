import * as React from 'react';
import { Navigate, Routes, Route } from 'react-router-dom';
import ModelRegistrySettings from '~/app/pages/settings/ModelRegistrySettings';
import ModelRegistriesManagePermissions from '~/app/pages/modelRegistrySettings/ModelRegistriesPermissions';

const ModelRegistrySettingsRoutes: React.FC = () => (
  <Routes>
    <Route path="/" element={<ModelRegistrySettings />} />
    <Route path="permissions/:mrName" element={<ModelRegistriesManagePermissions />} />
    <Route path="*" element={<Navigate to="/" />} />
  </Routes>
);

export default ModelRegistrySettingsRoutes;

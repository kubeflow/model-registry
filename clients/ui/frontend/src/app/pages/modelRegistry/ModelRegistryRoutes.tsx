import * as React from 'react';
import { Route, Routes } from 'react-router-dom';
import ModelRegistry from './ModelRegistry';

const ModelRegistryRoutes: React.FC = () => (
  <Routes>
    <Route index element={<ModelRegistry empty={false} />} />
  </Routes>
);

export default ModelRegistryRoutes;

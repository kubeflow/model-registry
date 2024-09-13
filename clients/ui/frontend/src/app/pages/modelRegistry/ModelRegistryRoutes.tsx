import * as React from 'react';
import { Route, Routes } from 'react-router-dom';
import ModelRegistry from './screens/ModelRegistry';
import ModelRegistryCoreLoader from './ModelRegistryCoreLoader';
import { modelRegistryUrl } from './screens/routeUtils';

const ModelRegistryRoutes: React.FC = () => (
  <Routes>
    <Route
      path={'/:modelRegistry?/*'}
      element={
        <ModelRegistryCoreLoader
          getInvalidRedirectPath={(modelRegistry) => modelRegistryUrl(modelRegistry)}
        />
      }
    >
      <Route index element={<ModelRegistry empty={false} />} />
    </Route>
  </Routes>
);

export default ModelRegistryRoutes;

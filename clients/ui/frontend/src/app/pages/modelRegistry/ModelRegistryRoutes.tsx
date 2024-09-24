import * as React from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import ModelRegistry from './screens/ModelRegistry';
import ModelRegistryCoreLoader from './ModelRegistryCoreLoader';
import { modelRegistryUrl } from './screens/routeUtils';
import RegisteredModelsArchive from './screens/RegisteredModelsArchive/RegisteredModelsArchive';
import { ModelVersionsTab } from './screens/ModelVersions/const';
import ModelVersions from './screens/ModelVersions/ModelVersions';

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
      <Route path="registeredModels/:registeredModelId">
        <Route index element={<Navigate to={ModelVersionsTab.VERSIONS} replace />} />
        <Route
          path={ModelVersionsTab.VERSIONS}
          element={<ModelVersions tab={ModelVersionsTab.VERSIONS} empty={false} />}
        />
        <Route
          path={ModelVersionsTab.DETAILS}
          element={<ModelVersions tab={ModelVersionsTab.DETAILS} empty={false} />}
        />
        <Route path="*" element={<Navigate to="." />} />
      </Route>
      <Route path="registeredModels/archive">
        <Route index element={<RegisteredModelsArchive empty={false} />} />
        <Route path="*" element={<Navigate to="." />} />
      </Route>
      <Route path="*" element={<Navigate to="." />} />
    </Route>
  </Routes>
);

export default ModelRegistryRoutes;

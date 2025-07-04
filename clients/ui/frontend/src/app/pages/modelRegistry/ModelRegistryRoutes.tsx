import * as React from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import { DeploymentMode, useModularArchContext } from 'mod-arch-shared';
import ModelRegistry from './screens/ModelRegistry';
import ModelRegistryCoreLoader from './ModelRegistryCoreLoader';
import { modelRegistryUrl } from './screens/routeUtils';
import RegisteredModelsArchive from './screens/RegisteredModelsArchive/RegisteredModelsArchive';
import { ModelVersionsTab } from './screens/ModelVersions/const';
import ModelVersions from './screens/ModelVersions/ModelVersions';
import { ModelVersionDetailsTab } from './screens/ModelVersionDetails/const';
import ModelVersionsDetails from './screens/ModelVersionDetails/ModelVersionDetails';
import ModelVersionsArchive from './screens/ModelVersionsArchive/ModelVersionsArchive';
import ModelVersionsArchiveDetails from './screens/ModelVersionsArchive/ModelVersionArchiveDetails';
import ArchiveModelVersionDetails from './screens/ModelVersionsArchive/ArchiveModelVersionDetails';
import RegisteredModelsArchiveDetails from './screens/RegisteredModelsArchive/RegisteredModelArchiveDetails';
import RegisterModel from './screens/RegisterModel/RegisterModel';
import RegisterVersion from './screens/RegisterModel/RegisterVersion';

const ModelRegistryRoutes: React.FC = () => {
  const { config } = useModularArchContext();
  const { deploymentMode } = config;
  const isStandalone = deploymentMode === DeploymentMode.Standalone;
  const isFederated = deploymentMode === DeploymentMode.Federated;
  return (
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
          <Route path="registerVersion" element={<RegisterVersion />} />
          <Route path="versions/:modelVersionId">
            <Route index element={<Navigate to={ModelVersionDetailsTab.DETAILS} replace />} />
            <Route
              path={ModelVersionDetailsTab.DETAILS}
              element={<ModelVersionsDetails tab={ModelVersionDetailsTab.DETAILS} empty={false} />}
            />
            {(isStandalone || isFederated) && (
              <Route
                path={ModelVersionDetailsTab.DEPLOYMENTS}
                element={
                  <ModelVersionsDetails tab={ModelVersionDetailsTab.DEPLOYMENTS} empty={false} />
                }
              />
            )}
            <Route path="*" element={<Navigate to="." />} />
          </Route>
          <Route path="versions/archive">
            <Route index element={<ModelVersionsArchive empty={false} />} />
            <Route path=":modelVersionId">
              <Route index element={<Navigate to={ModelVersionDetailsTab.DETAILS} replace />} />
              <Route
                path={ModelVersionDetailsTab.DETAILS}
                element={
                  <ModelVersionsArchiveDetails tab={ModelVersionDetailsTab.DETAILS} empty={false} />
                }
              />
              {(isStandalone || isFederated) && (
                <Route
                  path={ModelVersionDetailsTab.DEPLOYMENTS}
                  element={
                    <ModelVersionsArchiveDetails
                      tab={ModelVersionDetailsTab.DEPLOYMENTS}
                      empty={false}
                    />
                  }
                />
              )}

              <Route path="*" element={<Navigate to="." />} />
            </Route>
            <Route path="*" element={<Navigate to="." />} />
          </Route>
          <Route path="*" element={<Navigate to="." />} />
        </Route>
        <Route path="registeredModels/archive">
          <Route index element={<RegisteredModelsArchive empty={false} />} />
          <Route path=":registeredModelId">
            <Route index element={<Navigate to={ModelVersionsTab.VERSIONS} replace />} />
            <Route
              path={ModelVersionsTab.DETAILS}
              element={
                <RegisteredModelsArchiveDetails tab={ModelVersionsTab.DETAILS} empty={false} />
              }
            />
            <Route
              path={ModelVersionsTab.VERSIONS}
              element={
                <RegisteredModelsArchiveDetails tab={ModelVersionsTab.VERSIONS} empty={false} />
              }
            />
            <Route path="versions/:modelVersionId">
              <Route index element={<Navigate to={ModelVersionDetailsTab.DETAILS} replace />} />
              <Route
                path={ModelVersionDetailsTab.DETAILS}
                element={
                  <ArchiveModelVersionDetails tab={ModelVersionDetailsTab.DETAILS} empty={false} />
                }
              />
              {(isStandalone || isFederated) && (
                <Route
                  path={ModelVersionDetailsTab.DEPLOYMENTS}
                  element={
                    <ArchiveModelVersionDetails
                      tab={ModelVersionDetailsTab.DEPLOYMENTS}
                      empty={false}
                    />
                  }
                />
              )}

              <Route path="*" element={<Navigate to="." />} />
            </Route>
            <Route path="*" element={<Navigate to="." />} />
          </Route>
          <Route path="*" element={<Navigate to="." />} />
        </Route>
        <Route path="registerModel" element={<RegisterModel />} />
        <Route path="registerVersion" element={<RegisterVersion />} />
        <Route path="*" element={<Navigate to="." />} />
      </Route>
    </Routes>
  );
};

export default ModelRegistryRoutes;

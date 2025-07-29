import * as React from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import { NotFound } from 'mod-arch-shared';
import { NavDataItem } from '~/app/standalone/types';
import ModelRegistrySettingsRoutes from '~/app/pages/settings/ModelRegistrySettingsRoutes';
import ModelRegistryRoutes from '~/app/pages/modelRegistry/ModelRegistryRoutes';
import ModelCatalogPage from './pages/modelCatalog/screens/ModelCatalogPage';
import { ModelCatalogContextProvider } from './context/modelCatalog/ModelCatalogContext';
import useUser from './hooks/useUser';
import { useModularArchContext } from 'mod-arch-shared';
import { DeploymentMode } from 'mod-arch-shared';

export const useAdminSettings = (): NavDataItem[] => {
  const { clusterAdmin } = useUser();

  if (!clusterAdmin) {
    return [];
  }

  return [
    {
      label: 'Settings',
      children: [{ label: 'Model Registry', path: '/model-registry-settings' }],
    },
  ];
};

export const useNavData = (): NavDataItem[] => {
  const { config } = useModularArchContext();
  const { deploymentMode } = config;
  const isStandalone = deploymentMode === DeploymentMode.Standalone;
  const isFederated = deploymentMode === DeploymentMode.Federated;

  const baseNavItems = [
    {
      label: 'Model Registry',
      path: '/model-registry',
    },
  ];

  // Only show Model Catalog in Standalone or Federated mode
  if (isStandalone || isFederated) {
    baseNavItems.push({
      label: 'Model Catalog',
      path: '/model-catalog',
    });
  }

  return [...baseNavItems, ...useAdminSettings()];
};

const AppRoutes: React.FC = () => {
  const { clusterAdmin } = useUser();
  const { config } = useModularArchContext();
  const { deploymentMode } = config;
  const isStandalone = deploymentMode === DeploymentMode.Standalone;
  const isFederated = deploymentMode === DeploymentMode.Federated;

  return (
    <Routes>
      <Route path="/" element={<Navigate to="/model-registry" replace />} />
      <Route path="/model-registry/*" element={<ModelRegistryRoutes />} />
      {(isStandalone || isFederated) && (
        <Route
          path="/model-catalog"
          element={
            <ModelCatalogContextProvider>
              <ModelCatalogPage />
            </ModelCatalogContextProvider>
          }
        />
      )}
      <Route path="*" element={<NotFound />} />
      {clusterAdmin && (
        <Route path="/model-registry-settings/*" element={<ModelRegistrySettingsRoutes />} />
      )}
    </Routes>
  );
};

export default AppRoutes;

import * as React from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import { NotFound } from 'mod-arch-shared';
import { useModularArchContext, DeploymentMode } from 'mod-arch-core';
import { NavDataItem } from '~/app/standalone/types';
import ModelRegistrySettingsRoutes from './pages/settings/ModelRegistrySettingsRoutes';
import { modelRegistryUrl } from './pages/modelRegistry/screens/routeUtils';
import ModelRegistryRoutes from './pages/modelRegistry/ModelRegistryRoutes';
import ModelCatalogRoutes from './pages/modelCatalog/ModelCatalogRoutes';
import useUser from './hooks/useUser';

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
      path: '/ai-hub/registry',
    },
  ];

  // Only show Model Catalog in Standalone or Federated mode
  if (isStandalone || isFederated) {
    baseNavItems.push({
      label: 'Model Catalog',
      path: '/ai-hub/catalog',
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
      <Route path="/" element={<Navigate to={modelRegistryUrl()} replace />} />
      <Route path="/ai-hub/registry/*" element={<ModelRegistryRoutes />} />
      {(isStandalone || isFederated) && (
        <Route path="/ai-hub/catalog/*" element={<ModelCatalogRoutes />} />
      )}
      <Route path="*" element={<NotFound />} />
      {/* TODO: [Conditional render] Follow up add testing and conditional rendering when in standalone mode */}
      {clusterAdmin && (
        <Route path="/model-registry-settings/*" element={<ModelRegistrySettingsRoutes />} />
      )}
    </Routes>
  );
};

export default AppRoutes;

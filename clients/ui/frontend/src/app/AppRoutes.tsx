import * as React from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import { NotFound } from 'mod-arch-shared';
import { useModularArchContext, DeploymentMode } from 'mod-arch-core';
import { NavDataItem } from '~/app/standalone/types';
import ModelRegistrySettingsRoutes from './pages/settings/ModelRegistrySettingsRoutes';
import ModelRegistryRoutes from './pages/modelRegistry/ModelRegistryRoutes';
import ModelCatalogRoutes from './pages/modelCatalog/ModelCatalogRoutes';
import ModelCatalogSettingsRoutes from './pages/modelCatalogSettings/ModelCatalogSettingsRoutes';
import { modelCatalogUrl } from './routes/modelCatalog/catalogModel';
import {
  catalogSettingsUrl,
  CATALOG_SETTINGS_PAGE_TITLE,
} from './routes/modelCatalogSettings/modelCatalogSettings';
import { modelRegistryUrl } from './pages/modelRegistry/screens/routeUtils';
import useUser from './hooks/useUser';

export const useAdminSettings = (): NavDataItem[] => {
  const { clusterAdmin } = useUser();
  const { config } = useModularArchContext();
  const { deploymentMode } = config;
  const isStandalone = deploymentMode === DeploymentMode.Standalone;
  const isFederated = deploymentMode === DeploymentMode.Federated;

  if (!clusterAdmin) {
    return [];
  }

  const settingsChildren = [{ label: 'Model Registry', path: '/model-registry-settings' }];
  // Only show Model Catalog Settings in Standalone or Federated mode
  if (isStandalone || isFederated) {
    settingsChildren.push({ label: CATALOG_SETTINGS_PAGE_TITLE, path: catalogSettingsUrl() });
  }

  return [
    {
      label: 'Settings',
      children: settingsChildren,
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
      path: modelRegistryUrl(),
    },
  ];

  // Only show Model Catalog in Standalone or Federated mode
  if (isStandalone || isFederated) {
    baseNavItems.push({
      label: 'Model Catalog',
      path: modelCatalogUrl(),
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
      <Route path={`${modelRegistryUrl()}/*`} element={<ModelRegistryRoutes />} />
      {(isStandalone || isFederated) && (
        <>
          <Route path={`${modelCatalogUrl()}/*`} element={<ModelCatalogRoutes />} />
          <Route path={`${catalogSettingsUrl()}/*`} element={<ModelCatalogSettingsRoutes />} />
        </>
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

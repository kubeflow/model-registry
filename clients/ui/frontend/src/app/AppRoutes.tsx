import React from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';

import { NotFound, NavDataItem } from 'mod-arch-shared';
import ModelRegistryRoutes from '~/app/pages/modelRegistry/ModelRegistryRoutes';
import ModelRegistrySettingsRoutes from '~/app/pages/settings/ModelRegistrySettingsRoutes';
import useUser from '~/app/hooks/useUser';

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

export const useNavData = (): NavDataItem[] => [
  {
    label: 'Model Registry',
    path: '/model-registry',
  },
  ...useAdminSettings(),
];

export const AppRoutes: React.FC = () => {
  const { clusterAdmin } = useUser();

  return (
    <Routes>
      <Route path="/" element={<Navigate to="/model-registry" replace />} />
      <Route path="/model-registry/*" element={<ModelRegistryRoutes />} />
      <Route path="*" element={<NotFound />} />
      {/* TODO: [Conditional render] Follow up add testing and conditional rendering when in standalone mode*/}
      {clusterAdmin && (
        <Route path="/model-registry-settings/*" element={<ModelRegistrySettingsRoutes />} />
      )}
    </Routes>
  );
};

export default AppRoutes;

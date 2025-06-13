import React from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';

import ModelRegistryRoutes from '~/app/pages/modelRegistry/ModelRegistryRoutes';
import ModelRegistrySettingsRoutes from '~/app/pages/settings/ModelRegistrySettingsRoutes';
import useUser from '~/app/hooks/useUser';
import { NotFound } from '~/app/components/NotFound';
import { AppLayout } from '~/app/AppLayout';

type NavDataItem = {
    label: string;
    path?: string;
    children?: NavDataItem[];
};

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
            <Route path="/" element={<AppLayout />}>
                <Route index element={<Navigate to="/model-registry" replace />} />
                <Route path="/model-registry/*" element={<ModelRegistryRoutes />} />
                {clusterAdmin && (
                    <Route path="/model-registry-settings/*" element={<ModelRegistrySettingsRoutes />} />
                )}
                <Route path="*" element={<NotFound />} />
            </Route>
        </Routes>
    );
};

export default AppRoutes;

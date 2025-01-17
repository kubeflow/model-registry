import * as React from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import { NotFound } from './pages/notFound/NotFound';
import ModelRegistrySettingsRoutes from './pages/settings/ModelRegistrySettingsRoutes';
import ModelRegistryRoutes from './pages/modelRegistry/ModelRegistryRoutes';
import useUser from './hooks/useUser';

export const isNavDataGroup = (navItem: NavDataItem): navItem is NavDataGroup =>
  'children' in navItem;

type NavDataCommon = {
  label: string;
};

export type NavDataHref = NavDataCommon & {
  path: string;
};

export type NavDataGroup = NavDataCommon & {
  children: NavDataHref[];
};

type NavDataItem = NavDataHref | NavDataGroup;

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

const AppRoutes: React.FC = () => {
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

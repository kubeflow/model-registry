import * as React from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import { NotFound } from './pages/notFound/NotFound';
import ModelRegistrySettingsRoutes from './pages/settings/ModelRegistrySettingsRoutes';
import ModelRegistryRoutes from './pages/modelRegistry/ModelRegistryRoutes';
import { useAppContext } from './AppContext';

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
  const { isAdmin } = useAppContext().user;

  if (!isAdmin) {
    return [];
  }

  return [
    {
      label: 'Settings',
      children: [{ label: 'Model Registry', path: '/modelRegistrySettings' }],
    },
  ];
};

export const useNavData = (): NavDataItem[] => [
  {
    label: 'Model Registry',
    path: '/modelRegistry',
  },
  ...useAdminSettings(),
];

const AppRoutes: React.FC = () => {
  const { isAdmin } = useAppContext().user;

  return (
    <Routes>
      <Route path="/" element={<Navigate to="/modelRegistry" replace />} />
      <Route path="/modelRegistry/*" element={<ModelRegistryRoutes />} />
      <Route path="*" element={<NotFound />} />
      {isAdmin && (
        <Route path="/modelRegistrySettings/*" element={<ModelRegistrySettingsRoutes />} />
      )}
    </Routes>
  );
};

export default AppRoutes;

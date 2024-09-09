import * as React from 'react';
import { Route, Routes } from 'react-router-dom';
import { NotFound } from './pages/notFound/NotFound';
import ModelRegistrySettingsRoutes from './pages/settings/ModelRegistrySettingsRoutes';
import ModelRegistryRoutes from './pages/modelRegistry/ModelRegistryRoutes';

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
  // get auth access for example set admin as true
  const isAdmin = true; //this should be a call to getting auth / role access

  // TODO: Remove the linter skip when we implement authentication
  // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
  if (!isAdmin) {
    return [];
  }

  return [
    {
      label: 'Settings',
      children: [{ label: 'Model Registry', path: '/settings' }],
    },
  ];
};

export const useNavData = (): NavDataItem[] => [
  {
    label: 'Model Registry',
    path: '/',
  },
  ...useAdminSettings(),
];

const AppRoutes: React.FC = () => {
  const isAdmin = true;

  return (
    <Routes>
      <Route path="/" element={<ModelRegistryRoutes />} />
      <Route path="*" element={<NotFound />} />
      {
        // TODO: Remove the linter skip when we implement authentication
        // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
        isAdmin && <Route path="/settings" element={<ModelRegistrySettingsRoutes />} />
      }
    </Routes>
  );
};

export default AppRoutes;

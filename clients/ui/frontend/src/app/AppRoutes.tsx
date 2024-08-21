import * as React from 'react';
import { Route, Routes } from 'react-router-dom';
import { Dashboard } from './Dashboard/Dashboard';
import { Support } from './Support/Support';
import { NotFound } from './NotFound/NotFound';
import { Admin } from './Settings/Admin';

export const isNavDataGroup = (navItem: NavDataItem): navItem is NavDataGroup => 'children' in navItem;

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

  if (!isAdmin) {
    return [];
  }

  return [{
    label: 'Settings',
    children: [
      { label: 'Setting 1', path: '/admin' },
      { label: 'Setting 2', path: '/admin' },
      { label: 'Setting 3', path: '/admin'}
    ]
  }]
}

export const useNavData = (): NavDataItem[] => {
  return ([
    {
      label: 'Dashboard',
      path: '/'
    },
    {
      label: 'Support',
      path: '/support'
    },
    ...useAdminSettings()
  ])
}

const AppRoutes: React.FC = () => {
  const isAdmin = true;

  return (
    <Routes>
      <Route path="/" element={<Dashboard />} />
      <Route path="/support" element={<Support />} />
      <Route path="*" element={<NotFound />} />
      {isAdmin &&
        <Route path="/admin" element={<Admin />} />
      }
    </Routes>
  );
};

export default AppRoutes;

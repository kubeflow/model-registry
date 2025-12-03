import React from 'react';
import ReactDOM from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import {
  BrowserStorageContextProvider,
  NotificationContextProvider,
  ModularArchContextProvider,
  ModularArchConfig,
} from 'mod-arch-core';
import { ThemeProvider } from 'mod-arch-kubeflow';
import {
  BFF_API_VERSION,
  DEPLOYMENT_MODE,
  MANDATORY_NAMESPACE,
  STYLE_THEME,
  URL_PREFIX,
} from '~/app/utilities/const';
import App from '~/app/App';

const root = ReactDOM.createRoot(document.getElementById('root')!);

const modularArchConfig: ModularArchConfig = {
  deploymentMode: DEPLOYMENT_MODE,
  URL_PREFIX,
  BFF_API_VERSION,
  mandatoryNamespace: MANDATORY_NAMESPACE,
};

// Wrapper component that provides all context providers
const RootLayout: React.FC = () => (
  <ModularArchContextProvider config={modularArchConfig}>
    <ThemeProvider theme={STYLE_THEME}>
      <BrowserStorageContextProvider>
        <NotificationContextProvider>
          <App />
        </NotificationContextProvider>
      </BrowserStorageContextProvider>
    </ThemeProvider>
  </ModularArchContextProvider>
);

// Use createBrowserRouter for data router features (useBlocker, etc.)
const router = createBrowserRouter([
  {
    path: '*',
    element: <RootLayout />,
  },
]);

root.render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
);

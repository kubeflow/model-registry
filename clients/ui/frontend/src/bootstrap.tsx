import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter as Router } from 'react-router-dom';
import {
  BrowserStorageContextProvider,
  NotificationContextProvider,
  NamespaceSelectorContextProvider,
  DashboardScriptLoader,
  ThemeProvider,
  ModularArchContextProvider,
} from 'mod-arch-shared';
import 'mod-arch-shared/style/MUI-theme.scss';
import App from './app/App';
import {
  BFF_API_VERSION,
  isIntegrated,
  isMUITheme,
  isPlatformKubeflow,
  URL_PREFIX,
} from './app/utilities/const';

const root = ReactDOM.createRoot(document.getElementById('root')!);

const modularArchConfig = {
  isMUITheme: isMUITheme(),
  isIntegrated: isIntegrated(),
  isPlatformKubeflow: isPlatformKubeflow(),
  URL_PREFIX,
  BFF_API_VERSION,
};

root.render(
  <React.StrictMode>
    <Router>
      <ModularArchContextProvider config={modularArchConfig}>
        <BrowserStorageContextProvider>
          <ThemeProvider>
            <NotificationContextProvider>
              <DashboardScriptLoader>
                <NamespaceSelectorContextProvider>
                  <App />
                </NamespaceSelectorContextProvider>
              </DashboardScriptLoader>
            </NotificationContextProvider>
          </ThemeProvider>
        </BrowserStorageContextProvider>
      </ModularArchContextProvider>
    </Router>
  </React.StrictMode>,
);

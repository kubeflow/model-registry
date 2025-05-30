import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter as Router } from 'react-router-dom';
import {
  BrowserStorageContextProvider,
  NotificationContextProvider,
  ThemeProvider,
  ModularArchContextProvider,
  ModularArchConfig,
} from 'mod-arch-shared';
import App from './app/App';
import {
  BFF_API_VERSION,
  DEPLOYMENT_MODE,
  PLATFORM_MODE,
  STYLE_THEME,
  URL_PREFIX,
} from './app/utilities/const';

const root = ReactDOM.createRoot(document.getElementById('root')!);

const modularArchConfig: ModularArchConfig = {
  platformMode: PLATFORM_MODE,
  deploymentMode: DEPLOYMENT_MODE,
  URL_PREFIX,
  BFF_API_VERSION,
};

root.render(
  <React.StrictMode>
    <Router>
      <ModularArchContextProvider config={modularArchConfig}>
        <ThemeProvider theme={STYLE_THEME}>
          <BrowserStorageContextProvider>
            <NotificationContextProvider>
              <App />
            </NotificationContextProvider>
          </BrowserStorageContextProvider>
        </ThemeProvider>
      </ModularArchContextProvider>
    </Router>
  </React.StrictMode>,
);

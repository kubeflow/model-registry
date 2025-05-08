import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter as Router } from 'react-router-dom';
import {
  BrowserStorageContextProvider,
  NotificationContextProvider,
  NamespaceSelectorContextProvider,
  DashboardScriptLoader,
  ThemeProvider,
} from 'mod-arch-shared';
import App from './app/App';

const root = ReactDOM.createRoot(document.getElementById('root')!);

root.render(
  <React.StrictMode>
    <Router>
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
    </Router>
  </React.StrictMode>,
);

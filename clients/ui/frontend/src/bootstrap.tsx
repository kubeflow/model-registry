import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter as Router } from 'react-router-dom';
import App from './app/App';
import { BrowserStorageContextProvider } from './shared/components/browserStorage/BrowserStorageContext';
import { NotificationContextProvider } from './app/context/NotificationContext';
import { NamespaceSelectorContextProvider } from './shared/context/NamespaceSelectorContext';
import DashboardScriptLoader from './shared/context/DashboardScriptLoader';
import ThemeProvider from './app/ThemeContext';

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

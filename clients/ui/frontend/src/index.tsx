import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter as Router } from 'react-router-dom';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import App from './app/App';
import { BrowserStorageContextProvider } from './shared/components/browserStorage/BrowserStorageContext';
import { NotificationContextProvider } from './app/context/NotificationContext';
import { NamespaceSelectorContextProvider } from './shared/context/NamespaceSelectorContext';
import DashboardScriptLoader from './shared/context/DashboardScriptLoader';

const theme = createTheme({ cssVariables: true });
const root = ReactDOM.createRoot(document.getElementById('root')!);

root.render(
  <React.StrictMode>
    <Router>
      <BrowserStorageContextProvider>
        <ThemeProvider theme={theme}>
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

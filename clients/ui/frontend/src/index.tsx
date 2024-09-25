import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter as Router } from 'react-router-dom';
import App from './app/App';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import { BrowserStorageContextProvider } from './components/browserStorage/BrowserStorageContext';

const theme = createTheme({ cssVariables: true });
const root = ReactDOM.createRoot(document.getElementById('root')!);

root.render(
  <React.StrictMode>
    <Router>
      <BrowserStorageContextProvider>
        <ThemeProvider theme={theme}>
          <App />
        </ThemeProvider>
      </BrowserStorageContextProvider>
    </Router>
  </React.StrictMode>,
);

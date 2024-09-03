import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter as Router } from 'react-router-dom';
import App from './app/App';
import { BrowserStorageContextProvider } from './components/browserStorage/BrowserStorageContext';

const root = ReactDOM.createRoot(document.getElementById('root')!);
root.render(
  <React.StrictMode>
    <Router>
      <BrowserStorageContextProvider>
        <App />
      </BrowserStorageContextProvider>
    </Router>
  </React.StrictMode>,
);

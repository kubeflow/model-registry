import * as React from 'react';
import { Route, Routes } from 'react-router-dom';
import { McpCatalogContextProvider } from '~/app/context/mcpCatalog/McpCatalogContext';
import McpCatalog from './screens/McpCatalog';

const McpCatalogRoutes: React.FC = () => (
  <McpCatalogContextProvider>
    <Routes>
      <Route path="/" element={<McpCatalog />} />
      <Route path="*" element={<McpCatalog />} />
    </Routes>
  </McpCatalogContextProvider>
);

export default McpCatalogRoutes;

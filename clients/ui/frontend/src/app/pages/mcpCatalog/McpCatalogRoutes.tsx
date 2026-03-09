import * as React from 'react';
import { Route, Routes } from 'react-router-dom';
import { McpCatalogContextProvider } from '~/app/context/mcpCatalog/McpCatalogContext';
import McpCatalogCoreLoader from './McpCatalogCoreLoader';
import McpCatalog from './screens/McpCatalog';

const McpCatalogRoutes: React.FC = () => (
  <McpCatalogContextProvider>
    <Routes>
      <Route path="/*" element={<McpCatalogCoreLoader />}>
        <Route index element={<McpCatalog />} />
        <Route path="*" element={<McpCatalog />} />
      </Route>
    </Routes>
  </McpCatalogContextProvider>
);

export default McpCatalogRoutes;

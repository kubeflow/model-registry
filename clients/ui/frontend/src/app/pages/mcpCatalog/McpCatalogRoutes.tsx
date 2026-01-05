import * as React from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import McpCatalogCoreLoader from './McpCatalogCoreLoader';
import McpCatalog from './screens/McpCatalog';
import McpServerDetailsPage from './screens/McpServerDetailsPage';

const McpCatalogRoutes: React.FC = () => (
  <Routes>
    <Route path="/*" element={<McpCatalogCoreLoader />}>
      <Route index element={<McpCatalog />} />
      <Route path=":serverId" element={<McpServerDetailsPage />} />
      <Route path="*" element={<Navigate to="." />} />
    </Route>
  </Routes>
);

export default McpCatalogRoutes;

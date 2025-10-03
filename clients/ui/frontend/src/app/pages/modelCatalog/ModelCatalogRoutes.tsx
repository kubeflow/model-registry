import * as React from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import { ModelCatalogContextProvider } from '~/app/context/modelCatalog/ModelCatalogContext';
import ModelCatalogCoreLoader from './ModelCatalogCoreLoader';
import ModelDetailsPage from './screens/ModelDetailsPage';
import RegisterCatalogModelPage from './screens/RegisterCatalogModelPage';
import ModelCatalog from './screens/ModelCatalog';
import { ModelDetailsTab } from './screens/ModelDetailsTabs';

const ModelCatalogRoutes: React.FC = () => (
  <ModelCatalogContextProvider>
    <Routes>
      <Route path="/:sourceId?/*" element={<ModelCatalogCoreLoader />}>
        <Route index element={<ModelCatalog />} />
        <Route path=":modelName">
          <Route index element={<Navigate to={ModelDetailsTab.OVERVIEW} replace />} />
          <Route
            path={ModelDetailsTab.OVERVIEW}
            element={<ModelDetailsPage tab={ModelDetailsTab.OVERVIEW} />}
          />
          <Route
            path={ModelDetailsTab.PERFORMANCE_INSIGHTS}
            element={<ModelDetailsPage tab={ModelDetailsTab.PERFORMANCE_INSIGHTS} />}
          />
          <Route path="register" element={<RegisterCatalogModelPage />} />
          <Route path="*" element={<Navigate to="." />} />
        </Route>
        <Route path="*" element={<Navigate to="." />} />
      </Route>
    </Routes>
  </ModelCatalogContextProvider>
);

export default ModelCatalogRoutes;

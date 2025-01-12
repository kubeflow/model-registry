import { Breadcrumb, BreadcrumbItem } from '@patternfly/react-core';
import React from 'react';
import { Link } from 'react-router-dom';
import { RegisteredModel } from '~/app/types';
import {
  modelVersionArchiveUrl,
  registeredModelUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';

type ModelVersionArchiveDetailsBreadcrumbProps = {
  preferredModelRegistry?: string;
  registeredModel: RegisteredModel | null;
  modelVersionName?: string;
};

const ModelVersionArchiveDetailsBreadcrumb: React.FC<ModelVersionArchiveDetailsBreadcrumbProps> = ({
  preferredModelRegistry,
  registeredModel,
  modelVersionName,
}) => (
  <Breadcrumb>
    <BreadcrumbItem
      render={() => <Link to="/model-registry">Model registry - {preferredModelRegistry}</Link>}
    />
    <BreadcrumbItem
      render={() => (
        <Link to={registeredModelUrl(registeredModel?.id, preferredModelRegistry)}>
          {registeredModel?.name || 'Loading...'}
        </Link>
      )}
    />
    <BreadcrumbItem
      render={() => (
        <Link to={modelVersionArchiveUrl(registeredModel?.id, preferredModelRegistry)}>
          Archived versions
        </Link>
      )}
    />
    <BreadcrumbItem isActive>{modelVersionName || 'Loading...'}</BreadcrumbItem>
  </Breadcrumb>
);

export default ModelVersionArchiveDetailsBreadcrumb;

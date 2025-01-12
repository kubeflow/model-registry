import { Breadcrumb, BreadcrumbItem } from '@patternfly/react-core';
import React from 'react';
import { Link } from 'react-router-dom';
import { RegisteredModel } from '~/app/types';
import {
  registeredModelArchiveDetailsUrl,
  registeredModelArchiveUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';

type ArchiveModelVersionDetailsBreadcrumbProps = {
  preferredModelRegistry?: string;
  registeredModel: RegisteredModel | null;
  modelVersionName?: string;
};

const ArchiveModelVersionDetailsBreadcrumb: React.FC<ArchiveModelVersionDetailsBreadcrumbProps> = ({
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
        <Link to={registeredModelArchiveUrl(preferredModelRegistry)}>Archived models</Link>
      )}
    />
    <BreadcrumbItem
      render={() => (
        <Link to={registeredModelArchiveDetailsUrl(registeredModel?.id, preferredModelRegistry)}>
          {registeredModel?.name || 'Loading...'}
        </Link>
      )}
    />
    <BreadcrumbItem isActive>{modelVersionName || 'Loading...'}</BreadcrumbItem>
  </Breadcrumb>
);

export default ArchiveModelVersionDetailsBreadcrumb;

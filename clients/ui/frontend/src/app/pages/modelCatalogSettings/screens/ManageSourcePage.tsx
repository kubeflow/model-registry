import * as React from 'react';
import { useParams, Link } from 'react-router-dom';
import { Breadcrumb, BreadcrumbItem } from '@patternfly/react-core';
import { ApplicationsPage, TitleWithIcon, ProjectObjectType } from 'mod-arch-shared';
import {
  CATALOG_SETTINGS_PAGE_TITLE,
  ADD_SOURCE_TITLE,
  ADD_SOURCE_DESCRIPTION,
  MANAGE_SOURCE_TITLE,
  MANAGE_SOURCE_DESCRIPTION,
  catalogSettingsUrl,
} from '~/app/routes/modelCatalogSettings/modelCatalogSettings';

const ManageSourcePage: React.FC = () => {
  const { catalogSourceId } = useParams<{ catalogSourceId?: string }>();
  const isAddMode = !catalogSourceId;
  const pageTitle = isAddMode ? ADD_SOURCE_TITLE : MANAGE_SOURCE_TITLE;
  const breadcrumbLabel = isAddMode ? ADD_SOURCE_TITLE : MANAGE_SOURCE_TITLE;
  const description = isAddMode ? ADD_SOURCE_DESCRIPTION : MANAGE_SOURCE_DESCRIPTION;

  return (
    <ApplicationsPage
      breadcrumb={
        <Breadcrumb>
          <BreadcrumbItem>
            <Link to={catalogSettingsUrl()}>{CATALOG_SETTINGS_PAGE_TITLE}</Link>
          </BreadcrumbItem>
          <BreadcrumbItem data-testid="breadcrumb-source-action" isActive>
            {breadcrumbLabel}
          </BreadcrumbItem>
        </Breadcrumb>
      }
      title={<TitleWithIcon title={pageTitle} objectType={ProjectObjectType.modelCatalog} />}
      description={description}
      empty={false}
      loaded
      provideChildrenPadding
    />
  );
};

export default ManageSourcePage;

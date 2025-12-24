import * as React from 'react';
import { useParams, Link } from 'react-router-dom';
import { Breadcrumb, BreadcrumbItem } from '@patternfly/react-core';
import { ApplicationsPage } from 'mod-arch-shared';
import {
  CATALOG_SETTINGS_PAGE_TITLE,
  ADD_SOURCE_TITLE,
  ADD_SOURCE_DESCRIPTION,
  MANAGE_SOURCE_TITLE,
  MANAGE_SOURCE_DESCRIPTION,
  catalogSettingsUrl,
} from '~/app/routes/modelCatalogSettings/modelCatalogSettings';
import ManageSourceForm from '~/app/pages/modelCatalogSettings/components/ManageSourceForm';
import { useCatalogSourceConfigBySourceId } from '~/app/hooks/modelCatalogSettings/useCatalogSourceConfigBySourceId';

const ManageSourcePage: React.FC = () => {
  const { catalogSourceId } = useParams<{ catalogSourceId?: string }>();
  const isAddMode = !catalogSourceId;
  const pageTitle = isAddMode ? ADD_SOURCE_TITLE : MANAGE_SOURCE_TITLE;
  const breadcrumbLabel = isAddMode ? ADD_SOURCE_TITLE : MANAGE_SOURCE_TITLE;
  const description = isAddMode ? ADD_SOURCE_DESCRIPTION : MANAGE_SOURCE_DESCRIPTION;

  const state = useCatalogSourceConfigBySourceId(catalogSourceId || '');
  const [existingSourceConfig, existingSourceConfigLoaded, existingSourceConfigLoadError] = state;

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
      title={pageTitle}
      description={description}
      errorMessage={catalogSourceId ? existingSourceConfigLoadError?.message : undefined}
      empty={catalogSourceId ? !existingSourceConfig : false}
      loaded={catalogSourceId ? existingSourceConfigLoaded : true}
      provideChildrenPadding
    >
      <ManageSourceForm
        existingSourceConfig={existingSourceConfig || undefined}
        isEditMode={!isAddMode}
      />
    </ApplicationsPage>
  );
};

export default ManageSourcePage;

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
import { ModelCatalogSettingsContext } from '~/app/context/modelCatalogSettings/ModelCatalogSettingsContext';
import { catalogSourceConfigToFormData } from '~/app/pages/modelCatalogSettings/utils/modelCatalogSettingsUtils';

const ManageSourcePage: React.FC = () => {
  const { catalogSourceId } = useParams<{ catalogSourceId?: string }>();
  const isAddMode = !catalogSourceId;
  const pageTitle = isAddMode ? ADD_SOURCE_TITLE : MANAGE_SOURCE_TITLE;
  const breadcrumbLabel = isAddMode ? ADD_SOURCE_TITLE : MANAGE_SOURCE_TITLE;
  const description = isAddMode ? ADD_SOURCE_DESCRIPTION : MANAGE_SOURCE_DESCRIPTION;

  const { catalogSourceConfigs, catalogSourceConfigsLoaded, catalogSourceConfigsLoadError } =
    React.useContext(ModelCatalogSettingsContext);

  const existingSourceConfig = catalogSourceConfigs?.catalogs.find(
    (sourceConfig) => sourceConfig.id === catalogSourceId,
  );
  const existingData = existingSourceConfig
    ? catalogSourceConfigToFormData(existingSourceConfig)
    : existingSourceConfig;

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
      errorMessage={catalogSourceConfigsLoadError?.message}
      empty={catalogSourceConfigs?.catalogs.length === 0}
      loaded={catalogSourceConfigsLoaded}
      provideChildrenPadding
    >
      <ManageSourceForm existingData={existingData} isEditMode={!isAddMode} />
    </ApplicationsPage>
  );
};

export default ManageSourcePage;

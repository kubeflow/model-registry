import * as React from 'react';
import { useParams, Link } from 'react-router-dom';
import { Breadcrumb, BreadcrumbItem, PageSection } from '@patternfly/react-core';
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
import ExpectedYamlFormatDrawer from '~/app/pages/modelCatalogSettings/components/expectedYamlFormatContent';
import { useCatalogSourceConfigBySourceId } from '~/app/hooks/modelCatalogSettings/useCatalogSourceConfigBySourceId';

const ManageSourcePage: React.FC = () => {
  const { catalogSourceId } = useParams<{ catalogSourceId?: string }>();
  const isAddMode = !catalogSourceId;
  const pageTitle = isAddMode ? ADD_SOURCE_TITLE : MANAGE_SOURCE_TITLE;
  const breadcrumbLabel = isAddMode ? ADD_SOURCE_TITLE : MANAGE_SOURCE_TITLE;
  const description = isAddMode ? ADD_SOURCE_DESCRIPTION : MANAGE_SOURCE_DESCRIPTION;

  const state = useCatalogSourceConfigBySourceId(catalogSourceId || '');
  const [existingSourceConfig, existingSourceConfigLoaded, existingSourceConfigLoadError] = state;
  const [isExpectedFormatDrawerOpen, setIsExpectedFormatDrawerOpen] = React.useState(false);

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
      <PageSection
        component="main"
        isFilled
        hasOverflowScroll={false}
        padding={{ default: 'noPadding' }}
        aria-label="Manage source content"
        style={{ minHeight: 0, display: 'flex', flexDirection: 'column', position: 'relative' }}
      >
        <div
          style={{
            position: 'relative',
            flex: 1,
            minHeight: 0,
            display: 'flex',
            flexDirection: 'column',
            width: '100%',
          }}
        >
          <ManageSourceForm
            existingSourceConfig={existingSourceConfig || undefined}
            isEditMode={!isAddMode}
            onOpenExpectedFormatDrawer={() => setIsExpectedFormatDrawerOpen((prev) => !prev)}
          />
        </div>
      </PageSection>
      <ExpectedYamlFormatDrawer
        isOpen={isExpectedFormatDrawerOpen}
        onClose={() => setIsExpectedFormatDrawerOpen(false)}
      />
    </ApplicationsPage>
  );
};

export default ManageSourcePage;

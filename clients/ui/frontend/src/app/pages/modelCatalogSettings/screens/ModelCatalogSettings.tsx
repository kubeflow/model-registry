import * as React from 'react';
import { Button, EmptyState, EmptyStateBody, EmptyStateVariant } from '@patternfly/react-core';
import { PlusCircleIcon } from '@patternfly/react-icons';
import { useNavigate } from 'react-router-dom';
import { ProjectObjectType, TitleWithIcon, ApplicationsPage } from 'mod-arch-shared';
import {
  CATALOG_SETTINGS_PAGE_TITLE,
  CATALOG_SETTINGS_DESCRIPTION,
  addSourceUrl,
} from '~/app/routes/modelCatalogSettings/modelCatalogSettings';
import { ModelCatalogSettingsContext } from '~/app/context/modelCatalogSettings/ModelCatalogSettingsContext';
import CatalogSourceConfigsTable from './CatalogSourceConfigsTable';

const ModelCatalogSettings: React.FC = () => {
  const navigate = useNavigate();
  const {
    catalogSourceConfigs,
    catalogSourceConfigsLoaded,
    catalogSourceConfigsLoadError,
    apiState,
    refreshCatalogSourceConfigs,
  } = React.useContext(ModelCatalogSettingsContext);

  const configs = catalogSourceConfigs?.catalogs || [];
  const isEmpty = catalogSourceConfigsLoaded && configs.length === 0;

  const handleDeleteSource = React.useCallback(
    async (sourceId: string): Promise<void> => {
      if (!apiState.apiAvailable) {
        throw new Error('API not available');
      }
      await apiState.api.deleteCatalogSourceConfig({}, sourceId);
      refreshCatalogSourceConfigs();
    },
    [apiState.api, apiState.apiAvailable, refreshCatalogSourceConfigs],
  );

  return (
    <ApplicationsPage
      title={
        <TitleWithIcon
          title={CATALOG_SETTINGS_PAGE_TITLE}
          objectType={ProjectObjectType.modelCatalog}
        />
      }
      description={CATALOG_SETTINGS_DESCRIPTION}
      empty={isEmpty}
      emptyStatePage={
        <EmptyState
          headingLevel="h5"
          icon={PlusCircleIcon}
          titleText="No catalog sources"
          variant={EmptyStateVariant.lg}
          data-testid="catalog-settings-empty-state"
        >
          <EmptyStateBody>
            No catalog sources have been configured. Add a source to get started.
          </EmptyStateBody>
          <Button
            variant="primary"
            icon={<PlusCircleIcon />}
            onClick={() => navigate(addSourceUrl())}
            data-testid="add-source-button-empty"
          >
            Add a source
          </Button>
        </EmptyState>
      }
      loaded={catalogSourceConfigsLoaded}
      loadError={catalogSourceConfigsLoadError}
      errorMessage="Unable to load catalog source configurations."
      provideChildrenPadding
    >
      <CatalogSourceConfigsTable
        catalogSourceConfigs={configs}
        onAddSource={() => navigate(addSourceUrl())}
        onDeleteSource={handleDeleteSource}
      />
    </ApplicationsPage>
  );
};

export default ModelCatalogSettings;

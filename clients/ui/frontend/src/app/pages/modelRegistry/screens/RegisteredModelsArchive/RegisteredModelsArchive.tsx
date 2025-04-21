import React from 'react';
import { Breadcrumb, BreadcrumbItem } from '@patternfly/react-core';
import { Link } from 'react-router-dom';
import { ApplicationsPage } from 'mod-arch-shared';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import { filterArchiveModels } from '~/app/utils';
import useRegisteredModels from '~/app/hooks/useRegisteredModels';
import useModelVersions from '~/app/hooks/useModelVersions';
import RegisteredModelsArchiveListView from './RegisteredModelsArchiveListView';

type RegisteredModelsArchiveProps = Omit<
  React.ComponentProps<typeof ApplicationsPage>,
  'breadcrumb' | 'title' | 'loadError' | 'loaded' | 'provideChildrenPadding'
>;

const RegisteredModelsArchive: React.FC<RegisteredModelsArchiveProps> = ({ ...pageProps }) => {
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);
  const [registeredModels, modelsLoaded, modelsLoadError, refreshModels] = useRegisteredModels();
  const [modelVersions, versionsLoaded, versionsLoadError, refreshVersions] = useModelVersions();

  const loaded = modelsLoaded && versionsLoaded;
  const loadError = modelsLoadError || versionsLoadError;

  const refresh = React.useCallback(() => {
    refreshModels();
    refreshVersions();
  }, [refreshModels, refreshVersions]);

  return (
    <ApplicationsPage
      {...pageProps}
      breadcrumb={
        <Breadcrumb>
          <BreadcrumbItem
            render={() => (
              <Link to="/model-registry">Model registry - {preferredModelRegistry?.name}</Link>
            )}
          />
          <BreadcrumbItem data-testid="archive-model-page-breadcrumb" isActive>
            Archived models
          </BreadcrumbItem>
        </Breadcrumb>
      }
      title={`Archived models of ${preferredModelRegistry?.name ?? ''}`}
      loadError={loadError}
      loaded={loaded}
      provideChildrenPadding
    >
      <RegisteredModelsArchiveListView
        registeredModels={filterArchiveModels(registeredModels.items)}
        modelVersions={modelVersions.items}
        refresh={refresh}
      />
    </ApplicationsPage>
  );
};

export default RegisteredModelsArchive;

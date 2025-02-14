import React from 'react';
import { Divider } from '@patternfly/react-core';
import ApplicationsPage from '~/shared/components/ApplicationsPage';
import useRegisteredModels from '~/app/hooks/useRegisteredModels';
import useModelVersions from '~/app/hooks/useModelVersions';
import { isMUITheme } from '~/shared/utilities/const';
import TitleWithIcon from '~/shared/components/design/TitleWithIcon';
import { ProjectObjectType } from '~/shared/components/design/utils';
import ModelRegistrySelectorNavigator from './ModelRegistrySelectorNavigator';
import RegisteredModelListView from './RegisteredModels/RegisteredModelListView';
import { modelRegistryUrl } from './routeUtils';

type ModelRegistryProps = Omit<
  React.ComponentProps<typeof ApplicationsPage>,
  | 'title'
  | 'description'
  | 'loadError'
  | 'loaded'
  | 'provideChildrenPadding'
  | 'removeChildrenTopPadding'
  | 'headerContent'
>;

const ModelRegistry: React.FC<ModelRegistryProps> = ({ ...pageProps }) => {
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
      title={
        !isMUITheme() ? (
          <TitleWithIcon title="Model Registry" objectType={ProjectObjectType.registeredModels} />
        ) : (
          'Model Registry'
        )
      }
      description={
        !isMUITheme() ? (
          'Select a model registry to view and manage your registered models. Model registries provide a structured and organized way to store, share, version, deploy, and track models.'
        ) : (
          <Divider />
        )
      }
      headerContent={
        <ModelRegistrySelectorNavigator
          getRedirectPath={(modelRegistryName) => modelRegistryUrl(modelRegistryName)}
        />
      }
      loadError={loadError}
      loaded={loaded}
      provideChildrenPadding
      removeChildrenTopPadding
    >
      <RegisteredModelListView
        registeredModels={registeredModels.items}
        modelVersions={modelVersions.items}
        refresh={refresh}
      />
    </ApplicationsPage>
  );
};

export default ModelRegistry;

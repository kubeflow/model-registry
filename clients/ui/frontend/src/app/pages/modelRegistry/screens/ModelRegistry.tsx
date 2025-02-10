import React from 'react';
import { Content, Divider, Flex, FlexItem } from '@patternfly/react-core';
import ApplicationsPage from '~/shared/components/ApplicationsPage';
import useRegisteredModels from '~/app/hooks/useRegisteredModels';
import useModelVersions from '~/app/hooks/useModelVersions';
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
      title="Model Registry"
      description={
        <Flex>
          <FlexItem>
            <Content>
              Select a model registry to view and manage your registered models. Model registries
              provide a structured and organized way to store, share, version, deploy, and track
              models.
            </Content>
          </FlexItem>
          <Divider />
        </Flex>
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

import React from 'react';
import ApplicationsPage from '~/app/components/ApplicationsPage';
import TitleWithIcon from '~/app/components/design/TitleWithIcon';
import { ProjectObjectType } from '~/app/components/design/utils';
import useRegisteredModels from '~/app/hooks/useRegisteredModels';
import ModelRegistrySelectorNavigator from './ModelRegistrySelectorNavigator';
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
  const [, loaded, loadError] = useRegisteredModels();

  return (
    <ApplicationsPage
      {...pageProps}
      title={
        <TitleWithIcon title="Model Registry" objectType={ProjectObjectType.registeredModels} />
      }
      description="Select a model registry to view and manage your registered models. Model registries provide a structured and organized way to store, share, version, deploy, and track models."
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
      TODO: Add table of registered models;
    </ApplicationsPage>
  );
};

export default ModelRegistry;

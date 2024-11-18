import React from 'react';
import ApplicationsPage from '~/shared/components/ApplicationsPage';
import TitleWithIcon from '~/shared/components/design/TitleWithIcon';
import { ProjectObjectType } from '~/shared/components/design/utils';
import useRegisteredModels from '~/app/hooks/useRegisteredModels';
import { filterLiveModels } from '~/app/pages/modelRegistry/screens/utils';
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
  const [registeredModels, loaded, loadError, refresh] = useRegisteredModels();

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
      <RegisteredModelListView
        registeredModels={filterLiveModels(registeredModels.items)}
        refresh={refresh}
      />
    </ApplicationsPage>
  );
};

export default ModelRegistry;

import React from 'react';
import ApplicationsPage from '~/app/components/ApplicationsPage';
import TitleWithIcon from '~/app/components/design/TitleWithIcon';
import { ProjectObjectType } from '~/app/components/design/utils';

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
  const [loaded, loadError] = [true, undefined]; // TODO: change with real usage

  return (
    <ApplicationsPage
      {...pageProps}
      title={
        <TitleWithIcon title="Model registry" objectType={ProjectObjectType.registeredModels} />
      }
      description="Select a model registry to view and manage your registered models. Model registries provide a structured and organized way to store, share, version, deploy, and track models."
      loadError={loadError}
      loaded={loaded}
      provideChildrenPadding
      removeChildrenTopPadding
    />
  );
};

export default ModelRegistry;

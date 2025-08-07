import * as React from 'react';
import { Navigate, Outlet, useParams } from 'react-router-dom';
import { Bullseye, Alert, Divider, Stack, StackItem } from '@patternfly/react-core';
import {
  ProjectObjectType,
  typedEmptyImage,
  TitleWithIcon,
  WhosMyAdministrator,
  KubeflowDocs,
  ApplicationsPage,
} from 'mod-arch-shared';
import { useThemeContext } from 'mod-arch-kubeflow';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import { ModelRegistryContextProvider } from '~/app/context/ModelRegistryContext';
import EmptyModelRegistryState from './screens/components/EmptyModelRegistryState';
import InvalidModelRegistry from './screens/InvalidModelRegistry';
import ModelRegistrySelectorNavigator from './screens/ModelRegistrySelectorNavigator';
import { modelRegistryUrl } from './screens/routeUtils';

type ApplicationPageProps = React.ComponentProps<typeof ApplicationsPage>;

type ModelRegistryCoreLoaderProps = {
  getInvalidRedirectPath: (modelRegistry: string) => string;
};

type ApplicationPageRenderState = Pick<
  ApplicationPageProps,
  'emptyStatePage' | 'empty' | 'headerContent'
>;

const ModelRegistryCoreLoader: React.FC<ModelRegistryCoreLoaderProps> = ({
  getInvalidRedirectPath,
}) => {
  const { modelRegistry } = useParams<{ modelRegistry: string }>();
  const {
    modelRegistriesLoaded,
    modelRegistriesLoadError,
    modelRegistries,
    preferredModelRegistry,
    updatePreferredModelRegistry,
  } = React.useContext(ModelRegistrySelectorContext);
  const { isMUITheme } = useThemeContext();

  const modelRegistryFromRoute = modelRegistries.find((mr) => mr.name === modelRegistry);

  React.useEffect(() => {
    if (modelRegistryFromRoute && preferredModelRegistry?.name !== modelRegistryFromRoute.name) {
      updatePreferredModelRegistry(modelRegistryFromRoute);
    }
  }, [modelRegistryFromRoute, updatePreferredModelRegistry, preferredModelRegistry?.name]);

  if (modelRegistriesLoadError) {
    return (
      <Bullseye>
        <Alert title="Model registry load error" variant="danger" isInline>
          {modelRegistriesLoadError.message}
        </Alert>
      </Bullseye>
    );
  }
  if (!modelRegistriesLoaded) {
    return <Bullseye>Loading model registries...</Bullseye>;
  }

  let renderStateProps: ApplicationPageRenderState & { children?: React.ReactNode };
  if (modelRegistries.length === 0) {
    renderStateProps = {
      empty: true,
      emptyStatePage: (
        <EmptyModelRegistryState
          testid="empty-model-registries-state"
          title={isMUITheme ? 'Deploy a model registry' : 'Request access to model registries'}
          description={
            isMUITheme
              ? 'To deploy a new model registry, follow the instructions in the docs below.'
              : 'To request a new model registry, or to request permission to access an existing model registry, contact your administrator.'
          }
          headerIcon={() => (
            <img src={typedEmptyImage(ProjectObjectType.registeredModels)} alt="" />
          )}
          customAction={isMUITheme ? <KubeflowDocs /> : <WhosMyAdministrator />}
        />
      ),
      headerContent: null,
    };
  } else if (modelRegistry) {
    const foundModelRegistry = modelRegistries.find((mr) => mr.name === modelRegistry);
    if (foundModelRegistry) {
      // Render the content
      return (
        <ModelRegistryContextProvider modelRegistryName={modelRegistry}>
          <Outlet />
        </ModelRegistryContextProvider>
      );
    }
    // They ended up on a non-valid project path
    renderStateProps = {
      empty: true,
      emptyStatePage: <InvalidModelRegistry modelRegistry={modelRegistry} />,
    };
  } else {
    // Redirect the namespace suffix into the URL
    const redirectModelRegistry = preferredModelRegistry ?? modelRegistries[0];
    return <Navigate to={getInvalidRedirectPath(redirectModelRegistry.name)} replace />;
  }

  return (
    <ApplicationsPage
      title={
        <TitleWithIcon title="Model Registry" objectType={ProjectObjectType.registeredModels} />
      }
      description={
        <Stack hasGutter>
          <StackItem>
            Select a model registry to view and manage your registered models. Model registries
            provide a structured and organized way to store, share, version, deploy, and track
            models.
          </StackItem>
          <StackItem>
            <Divider />
          </StackItem>
        </Stack>
      }
      headerContent={
        <ModelRegistrySelectorNavigator
          getRedirectPath={(modelRegistryName) => modelRegistryUrl(modelRegistryName)}
        />
      }
      {...renderStateProps}
      loaded
      provideChildrenPadding
    />
  );
};

export default ModelRegistryCoreLoader;

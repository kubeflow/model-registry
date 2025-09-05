import { Alert, Bullseye } from '@patternfly/react-core';
import { useThemeContext } from 'mod-arch-kubeflow';
import {
  ApplicationsPage,
  KubeflowDocs,
  ProjectObjectType,
  TitleWithIcon,
  typedEmptyImage,
  WhosMyAdministrator,
} from 'mod-arch-shared';
import * as React from 'react';
import { Navigate, Outlet, useParams } from 'react-router-dom';
import EmptyModelRegistryState from '../modelRegistry/screens/components/EmptyModelRegistryState';
import InvalidModelRegistry from '../modelRegistry/screens/InvalidModelRegistry';
import ModelRegistrySelectorNavigator from '../modelRegistry/screens/ModelRegistrySelectorNavigator';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { modelCatalogUrl } from './routeUtils';

type ApplicationPageProps = React.ComponentProps<typeof ApplicationsPage>;

type ApplicationPageRenderState = Pick<
  ApplicationPageProps,
  'emptyStatePage' | 'empty' | 'headerContent'
>;

type ModelCatalogCoreLoaderrProps = {
  getInvalidRedirectPath: (modelRegistry: string) => string;
};

const ModelCatalogCoreLoader: React.FC<ModelCatalogCoreLoaderrProps> = ({
  getInvalidRedirectPath,
}) => {
  const { sourceId } = useParams<{ sourceId: string }>();

  const {
    catalogSources,
    catalogSourcesLoaded,
    catalogSourcesLoadError,
    selectedSource,
    updateSelectedSource,
  } = React.useContext(ModelCatalogContext);

  const { isMUITheme } = useThemeContext();

  const modelCatalogFromRoute = catalogSources?.items.find((source) => source.id === sourceId);

  React.useEffect(() => {
    if (modelCatalogFromRoute && selectedSource?.name !== modelCatalogFromRoute.name) {
      updateSelectedSource(modelCatalogFromRoute);
    }
  }, [modelCatalogFromRoute, updateSelectedSource, selectedSource?.name]);

  if (catalogSourcesLoadError) {
    return (
      <Bullseye>
        <Alert title="Model catalog source load error" variant="danger" isInline>
          {catalogSourcesLoadError.message}
        </Alert>
      </Bullseye>
    );
  }

  if (!catalogSourcesLoaded) {
    return <Bullseye>Loading catalog sources...</Bullseye>;
  }
  let renderStateProps: ApplicationPageRenderState & { children?: React.ReactNode };
  if (catalogSources?.items.length === 0) {
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
  } else if (sourceId) {
    const foundCatalogSource = catalogSources?.items.find((source) => source.id === sourceId);
    if (foundCatalogSource) {
      // Render the content
      return <Outlet />;
    }
    // They ended up on a non-valid project path
    renderStateProps = {
      empty: true,
      emptyStatePage: <InvalidModelRegistry modelRegistry={sourceId} />,
    };
  } else {
    // Redirect the namespace suffix into the URL
    const redirectCatalogSource = selectedSource ?? catalogSources?.items[0];
    return <Navigate to={getInvalidRedirectPath(redirectCatalogSource?.id || '')} replace />;
  }

  return (
    <ApplicationsPage
      title={
        <TitleWithIcon title="Model catalog" objectType={ProjectObjectType.modelRegistrySettings} />
      }
      description={''}
      {...renderStateProps}
      loaded
      provideChildrenPadding
    />
  );
};

export default ModelCatalogCoreLoader;

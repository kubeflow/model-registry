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
import { Outlet } from 'react-router-dom';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import EmptyModelCatalogState from './EmptyModelCatalogState';

const ModelCatalogCoreLoader: React.FC = () => {
  const { catalogSources, catalogSourcesLoaded, catalogSourcesLoadError } =
    React.useContext(ModelCatalogContext);

  const { isMUITheme } = useThemeContext();

  if (catalogSourcesLoadError) {
    return (
      <ApplicationsPage
        title={<TitleWithIcon title="Model Catalog" objectType={ProjectObjectType.modelCatalog} />}
        description="Discover models that are available for your organization to register, deploy, and customize."
        headerContent={null}
        empty
        emptyStatePage={
          <Bullseye>
            <Alert title="Model catalog source load error" variant="danger" isInline>
              {catalogSourcesLoadError.message}
            </Alert>
          </Bullseye>
        }
        loaded
      />
    );
  }

  if (!catalogSourcesLoaded) {
    return (
      <ApplicationsPage
        title={<TitleWithIcon title="Model Catalog" objectType={ProjectObjectType.modelCatalog} />}
        description="Discover models that are available for your organization to register, deploy, and customize."
        headerContent={null}
        empty
        emptyStatePage={<Bullseye>Loading catalog sources...</Bullseye>}
        loaded={false}
      />
    );
  }

  if (catalogSources?.items?.length === 0) {
    return (
      <ApplicationsPage
        title={<TitleWithIcon title="Model Catalog" objectType={ProjectObjectType.modelCatalog} />}
        description="Discover models that are available for your organization to register, deploy, and customize."
        empty
        emptyStatePage={
          <EmptyModelCatalogState
            testid="empty-model-catalog-state"
            title={isMUITheme ? 'Deploy a model catalog' : 'Request access to model catalog'}
            description={
              isMUITheme
                ? 'To deploy model catalog, follow the instructions in the docs below.'
                : 'To request model catalog, or to request permission to access model catalog, contact your administrator.'
            }
            headerIcon={() => (
              <img src={typedEmptyImage(ProjectObjectType.modelRegistrySettings)} alt="" />
            )}
            customAction={isMUITheme ? <KubeflowDocs /> : <WhosMyAdministrator />}
          />
        }
        headerContent={null}
        loaded
        provideChildrenPadding
      />
    );
  }

  return <Outlet />;
};

export default ModelCatalogCoreLoader;

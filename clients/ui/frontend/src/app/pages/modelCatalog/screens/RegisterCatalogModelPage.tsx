import * as React from 'react';
import { Breadcrumb, BreadcrumbItem } from '@patternfly/react-core';
import { Link, useParams } from 'react-router-dom';
import { ApplicationsPage } from 'mod-arch-shared';
import { useModelCatalogSources } from '~/app/hooks/modelCatalog/useModelCatalogSources';
import { ModelCatalogItem } from '~/app/modelCatalogTypes';
import { ModelRegistryContextProvider } from '~/app/context/ModelRegistryContext';
import { getCatalogModelDetailsRoute } from '~/app/routes/modelCatalog/catalogModelDetails';
import {
  ModelRegistrySelectorContextProvider,
  ModelRegistrySelectorContext,
} from '~/app/context/ModelRegistrySelectorContext';
import RegisterCatalogModelForm from './RegisterCatalogModelForm';

type RouteParams = {
  modelId: string;
};

const RegisterCatalogModelPageInner: React.FC = () => {
  const { modelId } = useParams<RouteParams>();
  const { sources, loading, error } = useModelCatalogSources();
  const { modelRegistries, modelRegistriesLoaded } = React.useContext(ModelRegistrySelectorContext);

  const model: ModelCatalogItem | undefined = React.useMemo(() => {
    for (const source of sources) {
      const found = source.models?.find((m) => m.id === modelId);
      if (found) {
        return found;
      }
    }
    return undefined;
  }, [sources, modelId]);

  // Get the first available model registry from the context
  const preferredModelRegistry = modelRegistries.length > 0 ? modelRegistries[0] : null;

  // Check to see if data is loaded
  const isDataReady =
    !loading &&
    !error &&
    model !== undefined &&
    modelRegistriesLoaded &&
    modelRegistries.length > 0;

  return (
    <ApplicationsPage
      title={`Register ${model?.name || ''} model`}
      description="Create and register the first version of a new model."
      breadcrumb={
        <Breadcrumb>
          <BreadcrumbItem render={() => <Link to="/model-catalog">Model catalog</Link>} />
          <BreadcrumbItem
            data-testid="breadcrumb-model-name"
            render={() =>
              !model?.name ? (
                'Loading...'
              ) : (
                <Link to={getCatalogModelDetailsRoute({ modelName: model.name, tag: '' })}>
                  {model.name}
                </Link>
              )
            }
          />
          <BreadcrumbItem data-testid="breadcrumb-version-name" isActive>
            Register model
          </BreadcrumbItem>
        </Breadcrumb>
      }
      loaded={!loading}
      loadError={error}
      empty={false}
    >
      {isDataReady && preferredModelRegistry ? (
        <ModelRegistryContextProvider modelRegistryName={preferredModelRegistry.name}>
          <RegisterCatalogModelForm model={model} preferredModelRegistry={preferredModelRegistry} />
        </ModelRegistryContextProvider>
      ) : (
        <div>Loading...</div>
      )}
    </ApplicationsPage>
  );
};

const RegisterCatalogModelPage: React.FC = () => (
  <ModelRegistrySelectorContextProvider>
    <RegisterCatalogModelPageInner />
  </ModelRegistrySelectorContextProvider>
);

export default RegisterCatalogModelPage;

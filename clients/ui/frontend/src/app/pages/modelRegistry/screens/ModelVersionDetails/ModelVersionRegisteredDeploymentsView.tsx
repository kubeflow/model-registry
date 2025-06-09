import React from 'react';
import { ProjectObjectType, typedEmptyImage } from 'mod-arch-shared';
import ModelVersionDetailsTabs from '~/app/pages/modelRegistry/screens/ModelVersionDetails/ModelVersionDetailsTabs';
import EmptyModelRegistryState from '~/app/pages/modelRegistry/screens/components/EmptyModelRegistryState';

type ModelVersionRegisteredDeploymentsViewProps = Pick<
  React.ComponentProps<typeof ModelVersionDetailsTabs>,
  'inferenceServices' | 'servingRuntimes' | 'refresh'
>;

const ModelVersionRegisteredDeploymentsView: React.FC<
  ModelVersionRegisteredDeploymentsViewProps
  // TODO: [Model Serving] Remove this when model serving is available
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
> = ({ inferenceServices, servingRuntimes, refresh }) => {
  const isLoading = !inferenceServices.loaded || !servingRuntimes.loaded;

  if (!isLoading && !inferenceServices.data.length) {
    return (
      <EmptyModelRegistryState
        title="No deployments from model registry"
        headerIcon={() => (
          <img
            src={typedEmptyImage(ProjectObjectType.registeredModels, 'MissingDeployment')}
            alt="missing deployment"
          />
        )}
        description="No deployments initiated from model registry for this model version."
        testid="model-version-deployments-empty-state"
      />
    );
  }

  return (
    <EmptyModelRegistryState
      title="No deployments from model registry"
      headerIcon={() => (
        <img
          src={typedEmptyImage(ProjectObjectType.registeredModels, 'MissingDeployment')}
          alt="missing deployment"
        />
      )}
      description="No deployments initiated from model registry for this model version."
      testid="model-version-deployments-empty-state"
    />
  );

  // TODO: [Model Serving] Uncomment when model serving is available
  // return (
  //   <Stack hasGutter>
  //     <Alert variant="info" isInline title="Filtered list: Deployments from model registry only">
  //       This list includes only deployments that were initiated from the model registry. To view and
  //       manage all of your deployments, go to the <Link to="/modelServing">Model Serving</Link>{' '}
  //       page.
  //     </Alert>

  //     <InferenceServiceTable
  //       isGlobal
  //       getColumns={getVersionDetailsInferenceServiceColumns}
  //       inferenceServices={inferenceServices.data}
  //       servingRuntimes={servingRuntimes.data}
  //       isLoading={isLoading}
  //       refresh={refresh}
  //     />
  //   </Stack>
  // );
};

export default ModelVersionRegisteredDeploymentsView;

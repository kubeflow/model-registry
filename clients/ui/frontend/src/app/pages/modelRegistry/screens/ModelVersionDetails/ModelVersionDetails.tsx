import React, { useEffect } from 'react';
import { useNavigate, useParams } from 'react-router';
import { Breadcrumb, BreadcrumbItem, Flex, FlexItem, Truncate } from '@patternfly/react-core';
import { Link } from 'react-router-dom';
import {
  InferenceServiceKind,
  ServingRuntimeKind,
  FetchStateObject,
  ApplicationsPage,
} from 'mod-arch-shared';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import useRegisteredModelById from '~/app/hooks/useRegisteredModelById';
import useModelVersionById from '~/app/hooks/useModelVersionById';
import { ModelState } from '~/app/types';
import {
  archiveModelVersionDetailsUrl,
  modelVersionArchiveDetailsUrl,
  modelVersionUrl,
  registeredModelUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';
import { ModelVersionDetailsTab } from './const';
import ModelVersionSelector from './ModelVersionSelector';
import ModelVersionDetailsTabs from './ModelVersionDetailsTabs';
import ModelVersionsDetailsHeaderActions from './ModelVersionDetailsHeaderActions';

type ModelVersionsDetailProps = {
  tab: ModelVersionDetailsTab;
} & Omit<
  React.ComponentProps<typeof ApplicationsPage>,
  'breadcrumb' | 'title' | 'description' | 'loadError' | 'loaded' | 'provideChildrenPadding'
>;

const ModelVersionsDetails: React.FC<ModelVersionsDetailProps> = ({ tab, ...pageProps }) => {
  const navigate = useNavigate();

  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);

  const { modelVersionId: mvId, registeredModelId: rmId } = useParams();
  const [rm] = useRegisteredModelById(rmId);
  const [mv, mvLoaded, mvLoadError, refreshModelVersion] = useModelVersionById(mvId);

  const inferenceServices: FetchStateObject<InferenceServiceKind[]> = {
    data: [],
    loaded: false,
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    refresh: () => {},
  };
  const servingRuntimes: FetchStateObject<ServingRuntimeKind[]> = {
    data: [],
    loaded: false,
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    refresh: () => {},
  };

  const refresh = React.useCallback(() => {
    refreshModelVersion();
  }, [refreshModelVersion]);

  useEffect(() => {
    if (rm?.state === ModelState.ARCHIVED && mv?.id) {
      navigate(
        archiveModelVersionDetailsUrl(mv.id, mv.registeredModelId, preferredModelRegistry?.name),
      );
    } else if (mv?.state === ModelState.ARCHIVED) {
      navigate(
        modelVersionArchiveDetailsUrl(mv.id, mv.registeredModelId, preferredModelRegistry?.name),
      );
    }
  }, [rm?.state, mv?.id, mv?.state, mv?.registeredModelId, preferredModelRegistry?.name, navigate]);

  return (
    <ApplicationsPage
      {...pageProps}
      breadcrumb={
        <Breadcrumb>
          <BreadcrumbItem
            render={() => (
              <Link to="/model-registry">Model registry - {preferredModelRegistry?.name}</Link>
            )}
          />
          <BreadcrumbItem
            render={() => (
              <Link to={registeredModelUrl(rmId, preferredModelRegistry?.name)}>
                {rm?.name || 'Loading...'}
              </Link>
            )}
          />
          <BreadcrumbItem data-testid="breadcrumb-version-name" isActive>
            {mv?.name || 'Loading...'}
          </BreadcrumbItem>
        </Breadcrumb>
      }
      title={mv?.name}
      headerAction={
        mvLoaded &&
        mv && (
          <Flex
            spaceItems={{ default: 'spaceItemsMd' }}
            alignItems={{ default: 'alignItemsFlexStart' }}
          >
            <FlexItem style={{ width: '300px' }}>
              <ModelVersionSelector
                rmId={rmId}
                selection={mv}
                onSelect={(modelVersionId) =>
                  navigate(modelVersionUrl(modelVersionId, rmId, preferredModelRegistry?.name))
                }
              />
            </FlexItem>
            <FlexItem>
              <ModelVersionsDetailsHeaderActions
                mv={mv}
                hasDeployment={inferenceServices.data.length > 0}
                refresh={refresh}
              />
            </FlexItem>
          </Flex>
        )
      }
      description={<Truncate content={mv?.description || ''} />}
      loadError={mvLoadError}
      loaded={mvLoaded}
      provideChildrenPadding
    >
      {mv !== null && (
        <ModelVersionDetailsTabs
          tab={tab}
          modelVersion={mv}
          inferenceServices={inferenceServices}
          servingRuntimes={servingRuntimes}
          refresh={refresh}
        />
      )}
    </ApplicationsPage>
  );
};

export default ModelVersionsDetails;

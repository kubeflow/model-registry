import React, { useEffect } from 'react';
import { useNavigate, useParams } from 'react-router';
import { Breadcrumb, BreadcrumbItem, Truncate } from '@patternfly/react-core';
import { Link } from 'react-router-dom';
import { ModelVersionsTab } from '~/app/pages/modelRegistry/screens/ModelVersions/const';
import ApplicationsPage from '~/shared/components/ApplicationsPage';
import useModelVersionsByRegisteredModel from '~/app/hooks/useModelVersionsByRegisteredModel';
import useRegisteredModelById from '~/app/hooks/useRegisteredModelById';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import { filterLiveVersions } from '~/app/pages/modelRegistry/screens/utils';
import ModelVersionsHeaderActions from '~/app/pages/modelRegistry/screens/ModelVersions/ModelVersionsHeaderActions';
import { ModelState } from '~/app/types';
import { registeredModelArchiveDetailsUrl } from '~/app/pages/modelRegistry/screens/routeUtils';
import ModelVersionsTabs from './ModelVersionsTabs';

type ModelVersionsProps = {
  tab: ModelVersionsTab;
} & Omit<
  React.ComponentProps<typeof ApplicationsPage>,
  'breadcrumb' | 'title' | 'description' | 'loadError' | 'loaded' | 'provideChildrenPadding'
>;

const ModelVersions: React.FC<ModelVersionsProps> = ({ tab, ...pageProps }) => {
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);
  const { registeredModelId: rmId } = useParams();
  const [modelVersions, mvLoaded, mvLoadError, mvRefresh] = useModelVersionsByRegisteredModel(rmId);
  const [rm, rmLoaded, rmLoadError, rmRefresh] = useRegisteredModelById(rmId);
  const loadError = mvLoadError || rmLoadError;
  const loaded = mvLoaded && rmLoaded;
  const navigate = useNavigate();

  useEffect(() => {
    if (rm?.state === ModelState.ARCHIVED) {
      navigate(registeredModelArchiveDetailsUrl(rm.id, preferredModelRegistry?.name));
    }
  }, [rm?.state, rm?.id, preferredModelRegistry?.name, navigate]);

  return (
    <ApplicationsPage
      {...pageProps}
      breadcrumb={
        <Breadcrumb>
          <BreadcrumbItem
            render={() => (
              <Link to="/modelRegistry">Model registry - {preferredModelRegistry?.name}</Link>
            )}
          />
          <BreadcrumbItem data-testid="breadcrumb-model" isActive>
            {rm?.name || 'Loading...'}
          </BreadcrumbItem>
        </Breadcrumb>
      }
      title={rm?.name}
      headerAction={rm && <ModelVersionsHeaderActions rm={rm} />}
      description={<Truncate content={rm?.description || ''} />}
      loadError={loadError}
      loaded={loaded}
      provideChildrenPadding
    >
      {rm !== null && (
        <ModelVersionsTabs
          tab={tab}
          registeredModel={rm}
          refresh={rmRefresh}
          mvRefresh={mvRefresh}
          modelVersions={filterLiveVersions(modelVersions.items)}
        />
      )}
    </ApplicationsPage>
  );
};

export default ModelVersions;

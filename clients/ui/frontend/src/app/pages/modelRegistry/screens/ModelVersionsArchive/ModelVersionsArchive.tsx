import React from 'react';
import { useParams } from 'react-router';
import { Breadcrumb, BreadcrumbItem } from '@patternfly/react-core';
import { Link } from 'react-router-dom';
import { ApplicationsPage } from 'mod-arch-shared';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import useRegisteredModelById from '~/app/hooks/useRegisteredModelById';
import useModelVersionsByRegisteredModel from '~/app/hooks/useModelVersionsByRegisteredModel';
import { registeredModelUrl } from '~/app/pages/modelRegistry/screens/routeUtils';
import { filterArchiveVersions } from '~/app/utils';
import ModelVersionsArchiveListView from './ModelVersionsArchiveListView';

type ModelVersionsArchiveProps = Omit<
  React.ComponentProps<typeof ApplicationsPage>,
  'breadcrumb' | 'title' | 'description' | 'loadError' | 'loaded' | 'provideChildrenPadding'
>;

const ModelVersionsArchive: React.FC<ModelVersionsArchiveProps> = ({ ...pageProps }) => {
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);

  const { registeredModelId: rmId } = useParams();
  const [rm] = useRegisteredModelById(rmId);
  const [modelVersions, mvLoaded, mvLoadError, refresh] = useModelVersionsByRegisteredModel(rmId);

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
          <BreadcrumbItem data-testid="archive-version-page-breadcrumb" isActive>
            Archived versions
          </BreadcrumbItem>
        </Breadcrumb>
      }
      title={rm ? `Archived versions of ${rm.name}` : 'Archived versions'}
      loadError={mvLoadError}
      loaded={mvLoaded}
      provideChildrenPadding
    >
      {rm && (
        <ModelVersionsArchiveListView
          modelVersions={filterArchiveVersions(modelVersions.items)}
          refresh={refresh}
        />
      )}
    </ApplicationsPage>
  );
};

export default ModelVersionsArchive;

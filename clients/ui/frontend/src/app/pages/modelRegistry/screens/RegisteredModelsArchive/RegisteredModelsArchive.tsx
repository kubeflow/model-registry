import React from 'react';
import { Breadcrumb, BreadcrumbItem } from '@patternfly/react-core';
import { Link } from 'react-router-dom';
import ApplicationsPage from '~/shared/components/ApplicationsPage';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import { filterArchiveModels } from '~/app/pages/modelRegistry/screens/utils';
import useRegisteredModels from '~/app/hooks/useRegisteredModels';
import RegisteredModelsArchiveListView from './RegisteredModelsArchiveListView';

type RegisteredModelsArchiveProps = Omit<
  React.ComponentProps<typeof ApplicationsPage>,
  'breadcrumb' | 'title' | 'loadError' | 'loaded' | 'provideChildrenPadding'
>;

const RegisteredModelsArchive: React.FC<RegisteredModelsArchiveProps> = ({ ...pageProps }) => {
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);
  const [registeredModels, loaded, loadError, refresh] = useRegisteredModels();

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
          <BreadcrumbItem data-testid="archive-model-page-breadcrumb" isActive>
            Archived models
          </BreadcrumbItem>
        </Breadcrumb>
      }
      title={`Archived models of ${preferredModelRegistry?.name}`}
      loadError={loadError}
      loaded={loaded}
      provideChildrenPadding
    >
      <RegisteredModelsArchiveListView
        registeredModels={filterArchiveModels(registeredModels.items)}
        refresh={refresh}
      />
    </ApplicationsPage>
  );
};

export default RegisteredModelsArchive;

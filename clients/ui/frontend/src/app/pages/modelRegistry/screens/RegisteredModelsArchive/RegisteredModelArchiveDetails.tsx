import React, { useEffect } from 'react';
import { useNavigate, useParams } from 'react-router';
import { Button, Flex, FlexItem, Label, Truncate } from '@patternfly/react-core';
import { ApplicationsPage } from 'mod-arch-shared';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import { ModelRegistryContext } from '~/app/context/ModelRegistryContext';
import useRegisteredModelById from '~/app/hooks/useRegisteredModelById';
import useModelVersionsByRegisteredModel from '~/app/hooks/useModelVersionsByRegisteredModel';
import { ModelState } from '~/app/types';
import { ModelVersionsTab } from '~/app/pages/modelRegistry/screens/ModelVersions/const';
import ModelVersionsTabs from '~/app/pages/modelRegistry/screens/ModelVersions/ModelVersionsTabs';
import { RestoreRegisteredModelModal } from '~/app/pages/modelRegistry/screens/components/RestoreRegisteredModel';
import { registeredModelUrl } from '~/app/pages/modelRegistry/screens/routeUtils';
import RegisteredModelArchiveDetailsBreadcrumb from './RegisteredModelArchiveDetailsBreadcrumb';

type RegisteredModelsArchiveDetailsProps = {
  tab: ModelVersionsTab;
} & Omit<
  React.ComponentProps<typeof ApplicationsPage>,
  'breadcrumb' | 'title' | 'description' | 'loadError' | 'loaded' | 'provideChildrenPadding'
>;

const RegisteredModelsArchiveDetails: React.FC<RegisteredModelsArchiveDetailsProps> = ({
  tab,
  ...pageProps
}) => {
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);
  const { apiState } = React.useContext(ModelRegistryContext);

  const navigate = useNavigate();

  const { registeredModelId: rmId } = useParams();
  const [rm, rmLoaded, rmLoadError, rmRefresh] = useRegisteredModelById(rmId);
  const [modelVersions, mvLoaded, mvLoadError, refresh] = useModelVersionsByRegisteredModel(rmId);
  const [isRestoreModalOpen, setIsRestoreModalOpen] = React.useState(false);

  useEffect(() => {
    if (rm?.state === ModelState.LIVE) {
      navigate(registeredModelUrl(rm.id, preferredModelRegistry?.name));
    }
  }, [rm?.state, preferredModelRegistry?.name, rm?.id, navigate]);

  return (
    <>
      <ApplicationsPage
        {...pageProps}
        breadcrumb={
          <RegisteredModelArchiveDetailsBreadcrumb
            preferredModelRegistry={preferredModelRegistry?.name}
            registeredModel={rm}
          />
        }
        title={
          rm && (
            <Flex alignItems={{ default: 'alignItemsCenter' }}>
              <FlexItem>{rm.name}</FlexItem>
              <Label>Archived</Label>
            </Flex>
          )
        }
        headerAction={
          <Button data-testid="restore-button" onClick={() => setIsRestoreModalOpen(true)}>
            Restore model
          </Button>
        }
        description={<Truncate content={rm?.description || ''} />}
        loadError={rmLoadError}
        loaded={rmLoaded}
        provideChildrenPadding
      >
        {rm !== null && mvLoaded && !mvLoadError && (
          <ModelVersionsTabs
            tab={tab}
            isArchiveModel
            registeredModel={rm}
            modelVersions={modelVersions.items}
            refresh={rmRefresh}
            mvRefresh={refresh}
          />
        )}
      </ApplicationsPage>

      {rm !== null && isRestoreModalOpen ? (
        <RestoreRegisteredModelModal
          onCancel={() => setIsRestoreModalOpen(false)}
          onSubmit={() =>
            apiState.api
              .patchRegisteredModel(
                {},
                {
                  state: ModelState.LIVE,
                },
                rm.id,
              )
              .then(() => navigate(registeredModelUrl(rm.id, preferredModelRegistry?.name)))
          }
          registeredModelName={rm.name}
        />
      ) : null}
    </>
  );
};

export default RegisteredModelsArchiveDetails;

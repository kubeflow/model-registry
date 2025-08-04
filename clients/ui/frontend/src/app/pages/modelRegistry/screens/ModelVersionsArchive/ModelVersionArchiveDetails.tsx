import React, { useEffect } from 'react';
import { useNavigate, useParams } from 'react-router';
import { Button, Flex, FlexItem, Label, Truncate } from '@patternfly/react-core';
import { ApplicationsPage } from 'mod-arch-shared';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import { ModelRegistryContext } from '~/app/context/ModelRegistryContext';
import useRegisteredModelById from '~/app/hooks/useRegisteredModelById';
import useModelVersionById from '~/app/hooks/useModelVersionById';
import useModelArtifactsByVersionId from '~/app/hooks/useModelArtifactsByVersionId';
import { ModelState } from '~/app/types';
import {
  archiveModelVersionDetailsUrl,
  modelVersionUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';
import ModelVersionDetailsTabs from '~/app/pages/modelRegistry/screens/ModelVersionDetails/ModelVersionDetailsTabs';
import { RestoreModelVersionModal } from '~/app/pages/modelRegistry/screens/components/RestoreModelVersionModal';
import { ModelVersionDetailsTab } from '~/app/pages/modelRegistry/screens/ModelVersionDetails/const';
import ModelVersionArchiveDetailsBreadcrumb from './ModelVersionArchiveDetailsBreadcrumb';

type ModelVersionsArchiveDetailsProps = {
  tab: ModelVersionDetailsTab;
} & Omit<
  React.ComponentProps<typeof ApplicationsPage>,
  'breadcrumb' | 'title' | 'description' | 'loadError' | 'loaded' | 'provideChildrenPadding'
>;

const ModelVersionsArchiveDetails: React.FC<ModelVersionsArchiveDetailsProps> = ({
  tab,
  ...pageProps
}) => {
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);
  const { apiState } = React.useContext(ModelRegistryContext);

  const navigate = useNavigate();

  const { modelVersionId: mvId, registeredModelId: rmId } = useParams();
  const [rm] = useRegisteredModelById(rmId);
  const [mv, mvLoaded, mvLoadError, refreshModelVersion] = useModelVersionById(mvId);
  const [modelArtifacts, modelArtifactsLoaded, modelArtifactsLoadError, refreshModelArtifacts] =
    useModelArtifactsByVersionId(mvId);
  const [isRestoreModalOpen, setIsRestoreModalOpen] = React.useState(false);

  const refresh = React.useCallback(() => {
    refreshModelVersion();
    refreshModelArtifacts();
  }, [refreshModelVersion, refreshModelArtifacts]);

  const loaded = mvLoaded && modelArtifactsLoaded;
  const loadError = mvLoadError || modelArtifactsLoadError;

  useEffect(() => {
    if (rm?.state === ModelState.ARCHIVED && mv?.id) {
      navigate(
        archiveModelVersionDetailsUrl(mv.id, mv.registeredModelId, preferredModelRegistry?.name),
      );
    } else if (mv?.state === ModelState.LIVE) {
      navigate(modelVersionUrl(mv.id, mv.registeredModelId, preferredModelRegistry?.name));
    }
  }, [rm?.state, mv?.state, mv?.id, mv?.registeredModelId, preferredModelRegistry?.name, navigate]);

  return (
    <>
      <ApplicationsPage
        {...pageProps}
        breadcrumb={
          <ModelVersionArchiveDetailsBreadcrumb
            preferredModelRegistry={preferredModelRegistry?.name}
            registeredModel={rm}
            modelVersionName={mv?.name}
          />
        }
        title={
          mv && (
            <Flex alignItems={{ default: 'alignItemsCenter' }}>
              <FlexItem>{mv.name}</FlexItem>
              <Label>Archived</Label>
            </Flex>
          )
        }
        headerAction={
          <Button data-testid="restore-button" onClick={() => setIsRestoreModalOpen(true)}>
            Restore model version
          </Button>
        }
        description={<Truncate content={mv?.description || ''} />}
        loadError={loadError}
        loaded={loaded}
        provideChildrenPadding
      >
        {mv !== null && (
          <ModelVersionDetailsTabs
            isArchiveVersion
            tab={tab}
            modelVersion={mv}
            refresh={refresh}
            modelArtifacts={modelArtifacts}
          />
        )}
      </ApplicationsPage>
      {mv !== null && isRestoreModalOpen ? (
        <RestoreModelVersionModal
          onCancel={() => setIsRestoreModalOpen(false)}
          onSubmit={() =>
            apiState.api
              .patchModelVersion(
                {},
                {
                  state: ModelState.LIVE,
                },
                mv.id,
              )
              .then(() => navigate(modelVersionUrl(mv.id, rm?.id, preferredModelRegistry?.name)))
          }
          modelVersionName={mv.name}
        />
      ) : null}
    </>
  );
};

export default ModelVersionsArchiveDetails;

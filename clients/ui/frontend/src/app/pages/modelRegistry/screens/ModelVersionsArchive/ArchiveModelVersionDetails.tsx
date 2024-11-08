import React, { useEffect } from 'react';
import { useNavigate, useParams } from 'react-router';
import { Button, Flex, FlexItem, Label, Content, Tooltip, Truncate } from '@patternfly/react-core';

import ApplicationsPage from '~/shared/components/ApplicationsPage';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import useRegisteredModelById from '~/app/hooks/useRegisteredModelById';
import useModelVersionById from '~/app/hooks/useModelVersionById';
import { ModelState } from '~/app/types';
import { modelVersionUrl } from '~/app/pages/modelRegistry/screens/routeUtils';
import ModelVersionDetailsTabs from '~/app/pages/modelRegistry/screens/ModelVersionDetails/ModelVersionDetailsTabs';
import { ModelVersionDetailsTab } from '~/app/pages/modelRegistry/screens/ModelVersionDetails/const';
import ArchiveModelVersionDetailsBreadcrumb from './ArchiveModelVersionDetailsBreadcrumb';

type ArchiveModelVersionDetailsProps = {
  tab: ModelVersionDetailsTab;
} & Omit<
  React.ComponentProps<typeof ApplicationsPage>,
  'breadcrumb' | 'title' | 'description' | 'loadError' | 'loaded' | 'provideChildrenPadding'
>;

const ArchiveModelVersionDetails: React.FC<ArchiveModelVersionDetailsProps> = ({
  tab,
  ...pageProps
}) => {
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);
  const { modelVersionId: mvId, registeredModelId: rmId } = useParams();
  const [rm] = useRegisteredModelById(rmId);
  const [mv, mvLoaded, mvLoadError, refreshModelVersion] = useModelVersionById(mvId);
  const navigate = useNavigate();

  useEffect(() => {
    if (rm?.state === ModelState.LIVE && mv?.id) {
      navigate(modelVersionUrl(mv.id, mv.registeredModelId, preferredModelRegistry?.name));
    }
  }, [rm?.state, mv?.id, mv?.registeredModelId, preferredModelRegistry?.name, navigate]);

  return (
    <ApplicationsPage
      {...pageProps}
      breadcrumb={
        <ArchiveModelVersionDetailsBreadcrumb
          preferredModelRegistry={preferredModelRegistry?.name}
          registeredModel={rm}
          modelVersionName={mv?.name}
        />
      }
      title={
        mv && (
          <Flex>
            <FlexItem>
              <Content>{mv.name}</Content>
            </FlexItem>
            <FlexItem>
              <Label>Archived</Label>
            </FlexItem>
          </Flex>
        )
      }
      headerAction={
        <Tooltip content="The version of an archived model cannot be restored unless the model is restored.">
          <Button data-testid="restore-button" aria-label="restore version" isAriaDisabled>
            Restore version
          </Button>
        </Tooltip>
      }
      description={<Truncate content={mv?.description || ''} />}
      loadError={mvLoadError}
      loaded={mvLoaded}
      provideChildrenPadding
    >
      {mv !== null && (
        <ModelVersionDetailsTabs
          isArchiveVersion
          tab={tab}
          modelVersion={mv}
          refresh={refreshModelVersion}
        />
      )}
    </ApplicationsPage>
  );
};

export default ArchiveModelVersionDetails;

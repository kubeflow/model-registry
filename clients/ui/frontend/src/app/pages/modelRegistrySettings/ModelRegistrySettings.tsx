import * as React from 'react';
import {
  Button,
  Box,
  Typography,
} from '@mui/material';
import { AddCircleOutline } from '@mui/icons-material';
import { AreaContext } from '../../concepts/areas/AreaContext';
import ApplicationsPage from '../ApplicationsPage';
import RedirectErrorState from '../external/RedirectErrorState';
import TitleWithIcon from '../../concepts/design/TitleWithIcon';
import { ProjectObjectType } from '../../concepts/design/utils';
import useModelRegistriesBackend from '../../concepts/modelRegistrySettings/useModelRegistriesBackend';
import { ModelRegistriesContext } from '../../concepts/modelRegistry/context/ModelRegistriesContext';
import ModelRegistriesTable from './ModelRegistriesTable';
import CreateModal from './CreateModal';
import useModelRegistryRoleBindings from './useModelRegistryRoleBindings';

const ModelRegistrySettings: React.FC = () => {
  const { dscStatus } = React.useContext(AreaContext);
  const modelRegistryNamespace = dscStatus?.components?.modelregistry?.registriesNamespace;
  const [createModalOpen, setCreateModalOpen] = React.useState(false);

  const modelRegistriesResult = useModelRegistriesBackend();
  const roleBindings = useModelRegistryRoleBindings();
  const { refreshRulesReview } = React.useContext(ModelRegistriesContext);
  const loaded = modelRegistriesResult.loaded && roleBindings.loaded;

  const refreshAll = React.useCallback(
    () => Promise.all([modelRegistriesResult.refresh(), roleBindings.refresh(), refreshRulesReview()]),
    [modelRegistriesResult, roleBindings, refreshRulesReview],
  );

  const error = !modelRegistryNamespace
    ? new Error('No registries namespace could be found')
    : null;

  if (!modelRegistryNamespace) {
    return (
      <ApplicationsPage loaded empty={false}>
        <RedirectErrorState title="Could not load component state" errorMessage={error?.message} />
      </ApplicationsPage>
    );
  }
  return (
    <>
      <ApplicationsPage
        title={
          <TitleWithIcon
            title="Model registry settings"
            objectType={ProjectObjectType.modelRegistrySettings}
          />
        }
        description="Manage model registry settings for all users in your organization."
        loaded={loaded}
        loadError={modelRegistriesResult.error}
        errorMessage="Unable to load model registries."
        empty={modelRegistriesResult.data.length === 0}
        emptyStatePage={
          <Box sx={{ textAlign: 'center' }}>
            <AddCircleOutline sx={{ fontSize: 48 }} color="disabled" />
            <Typography variant="h5" component="h2">No model registries</Typography>
            <Typography variant="body1">
              To get started, create a model registry. You can manage permissions after creation.
            </Typography>
            <Button variant="contained" onClick={() => setCreateModalOpen(true)}>
              Create model registry
            </Button>
          </Box>
        }
        provideChildrenPadding
      >
        <ModelRegistriesTable
          modelRegistries={modelRegistriesResult.data}
          roleBindings={roleBindings}
          refresh={refreshAll}
          onCreateModelRegistryClick={() => {
            setCreateModalOpen(true);
          }}
        />
      </ApplicationsPage>
      {createModalOpen ? (
        <CreateModal onClose={() => setCreateModalOpen(false)} refresh={refreshAll} />
      ) : null}
    </>
  );
};

export default ModelRegistrySettings; 
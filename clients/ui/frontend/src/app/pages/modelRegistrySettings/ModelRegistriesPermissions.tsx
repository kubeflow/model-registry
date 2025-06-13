import React from 'react';
import {
  Breadcrumbs,
  Link as MUILink,
  Tabs,
  Tab,
  Box,
  Typography,
} from '@mui/material';
import { Link, Navigate, useParams } from 'react-router-dom';
import { ModelRegistryKind, RoleBindingKind } from '~/app/k8sTypes';
import { useGroups } from '~/app/api/k8s/groups';
import RoleBindingPermissions from '~/app/concepts/roleBinding/RoleBindingPermissions';
import ApplicationsPage from '~/app/pages/ApplicationsPage';
import { useModelRegistryNamespaceCR } from '~/app/concepts/modelRegistry/context/useModelRegistryNamespaceCR';
import { AreaContext } from '~/app/concepts/areas/AreaContext';
import {
  createModelRegistryRoleBinding,
  deleteModelRegistryRoleBinding,
} from '~/app/services/modelRegistrySettingsService';
import RedirectErrorState from '~/app/pages/external/RedirectErrorState';
import useModelRegistryRoleBindings from '~/app/pages/modelRegistrySettings/useModelRegistryRoleBindings';
import ProjectsSettingsTab from '~/app/pages/modelRegistrySettings/ProjectsTab/ProjectsSettingsTab';
import { RoleBindingPermissionsRoleType } from '~/app/pages/settings/roleBinding/types';

const ModelRegistriesManagePermissions: React.FC = () => {
  const { dscStatus } = React.useContext(AreaContext);
  const modelRegistryNamespace = dscStatus?.components?.modelregistry?.registriesNamespace;
  const [activeTabKey, setActiveTabKey] = React.useState(0);
  const [ownerReference, setOwnerReference] = React.useState<ModelRegistryKind>();
  const [groups] = useGroups();
  const roleBindings = useModelRegistryRoleBindings();
  const { mrName } = useParams<{ mrName: string }>();
  const state = useModelRegistryNamespaceCR(modelRegistryNamespace, mrName || '');
  const [modelRegistryCR, crLoaded] = state;
  const filteredRoleBindings = roleBindings.data.filter(
    (rb: RoleBindingKind) => rb.metadata.labels?.['app.kubernetes.io/name'] === mrName,
  );

  const error = !modelRegistryNamespace
    ? new Error('No registries namespace could be found')
    : null;

  React.useEffect(() => {
    if (modelRegistryCR) {
      setOwnerReference(modelRegistryCR);
    } else {
      setOwnerReference(undefined);
    }
  }, [modelRegistryCR]);

  if (!modelRegistryNamespace) {
    return (
      <ApplicationsPage loaded empty={false}>
        <RedirectErrorState title="Could not load component state" errorMessage={error?.message} />
      </ApplicationsPage>
    );
  }

  if (
    (roleBindings.loaded && filteredRoleBindings.length === 0) ||
    (crLoaded && !modelRegistryCR)
  ) {
    return <Navigate to="/modelRegistrySettings" replace />;
  }

  return (
    <ApplicationsPage
      title={`Manage ${mrName ?? ''} permissions`}
      description="Manage access to this model registry for individual users and user groups, and for service accounts in a project."
      breadcrumb={
        <Breadcrumbs>
          <MUILink component={Link} to="/modelRegistrySettings">
            Model registry settings
          </MUILink>
          <Typography color="text.primary">Manage Permissions</Typography>
        </Breadcrumbs>
      }
      loaded
      empty={false}
      provideChildrenPadding
    >
      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs value={activeTabKey} onChange={(e, newValue) => setActiveTabKey(newValue)}>
          <Tab label="Users" />
          <Tab label="Projects" />
        </Tabs>
      </Box>
      <Box sx={{ pt: 2 }}>
        {activeTabKey === 0 && (
          <RoleBindingPermissions
            ownerReference={ownerReference}
            roleBindingPermissionsRB={{ ...roleBindings, data: filteredRoleBindings }}
            groups={groups}
            createRoleBinding={createModelRegistryRoleBinding}
            deleteRoleBinding={deleteModelRegistryRoleBinding}
            projectName={modelRegistryNamespace}
            permissionOptions={[{
                type: RoleBindingPermissionsRoleType.DEFAULT,
                description: 'Default role for all users',
            }]}
            description="To enable access for all cluster users, add system:authenticated to the group list."
            roleRefKind="Role"
            roleRefName={`registry-user-${mrName ?? ''}`}
          />
        )}
        {activeTabKey === 1 && (
            <ProjectsSettingsTab
                ownerReference={ownerReference}
                projectName={modelRegistryNamespace}
                roleBindingPermissionsRB={{ ...roleBindings, data: filteredRoleBindings }}
                permissionOptions={[{
                    type: RoleBindingPermissionsRoleType.DEFAULT,
                    description: 'Default role for all projects',
                }]}
                description="To enable access for all service accounts in a project, add the project name to the projects list."
                roleRefName={`registry-user-${mrName ?? ''}`}
            />
        )}
      </Box>
    </ApplicationsPage>
  );
};

export default ModelRegistriesManagePermissions; 
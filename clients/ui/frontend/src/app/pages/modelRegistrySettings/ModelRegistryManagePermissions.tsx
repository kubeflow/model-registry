import React from 'react';
import { Breadcrumbs, Link as MUILink, Tabs, Tab, Box, Typography } from '@mui/material';
import { Link, Navigate, useParams } from 'react-router-dom';
import {
  ApplicationsPage,
  ModelRegistryKind,
  RoleBindingKind,
  useQueryParamNamespaces,
} from 'mod-arch-shared';
import { useGroups } from '~/app/hooks/useGroups';
import RoleBindingPermissions from '~/app/pages/settings/roleBinding/RoleBindingPermissions';
import { useModelRegistryCR } from '~/app/hooks/useModelRegistryCR';
import useModelRegistryRoleBindings from '~/app/pages/modelRegistrySettings/useModelRegistryRoleBindings';
import { RoleBindingPermissionsRoleType } from '~/app/pages/settings/roleBinding/types';
import {
  createModelRegistryRoleBindingWrapper,
  deleteModelRegistryRoleBindingWrapper,
} from '~/app/pages/settings/roleBindingUtils';
import RedirectErrorState from '~/app/shared/components/RedirectErrorState';

const ModelRegistriesManagePermissions: React.FC = () => {
  const modelRegistryNamespace = 'model-registry'; // TODO: This is a placeholder
  const [activeTabKey, setActiveTabKey] = React.useState(0);
  const [ownerReference, setOwnerReference] = React.useState<ModelRegistryKind>();
  const queryParams = useQueryParamNamespaces();
  const [groups] = useGroups(queryParams);
  const roleBindings = useModelRegistryRoleBindings(queryParams);
  const { mrName } = useParams<{ mrName: string }>();
  const [modelRegistryCR, crLoaded] = useModelRegistryCR(modelRegistryNamespace, queryParams);

  const filteredRoleBindings = roleBindings.data.filter(
    (rb: RoleBindingKind) => rb.metadata.labels?.['app.kubernetes.io/name'] === mrName,
  );

  React.useEffect(() => {
    if (modelRegistryCR) {
      setOwnerReference(modelRegistryCR);
    } else {
      setOwnerReference(undefined);
    }
  }, [modelRegistryCR]);

  if (!queryParams.namespace) {
    return (
      <ApplicationsPage loaded empty={false}>
        <RedirectErrorState title="Could not load component state" />
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
            createRoleBinding={createModelRegistryRoleBindingWrapper}
            deleteRoleBinding={deleteModelRegistryRoleBindingWrapper}
            projectName={modelRegistryNamespace}
            permissionOptions={[
              {
                type: RoleBindingPermissionsRoleType.DEFAULT,
                description: 'Default role for all users',
              },
            ]}
            description="To enable access for all cluster users, add system:authenticated to the group list."
            roleRefKind="Role"
            roleRefName={`registry-user-${mrName ?? ''}`}
          />
        )}
        {/* TODO: Projects tab */}
      </Box>
    </ApplicationsPage>
  );
};

export default ModelRegistriesManagePermissions;

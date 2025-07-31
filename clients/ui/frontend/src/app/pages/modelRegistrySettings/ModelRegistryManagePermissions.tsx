import React from 'react';
import { Breadcrumbs, Link as MUILink, Tabs, Tab, Box, Typography } from '@mui/material';
import { Link, Navigate } from 'react-router-dom';
import { ApplicationsPage } from 'mod-arch-shared';
import RoleBindingPermissions from '~/app/pages/settings/roleBinding/RoleBindingPermissions';
import RedirectErrorState from '~/app/shared/components/RedirectErrorState';
import { useModelRegistryPermissionsLogic } from './useModelRegistryPermissionsLogic';

const ModelRegistryManagePermissions: React.FC = () => {
  const {
    activeTabKey,
    setActiveTabKey,
    ownerReference,
    groups,
    filteredRoleBindings,
    filteredNamespaceRoleBindings,
    mrName,
    modelRegistryNamespace,
    roleBindings,
    userPermissionOptions,
    namespacePermissionOptions,
    createUserRoleBinding,
    deleteUserRoleBinding,
    createNamespaceRoleBinding,
    deleteNamespaceRoleBinding,
    userRoleRefName,
    namespaceRoleRefName,
    shouldShowError,
    shouldRedirect,
  } = useModelRegistryPermissionsLogic();

  // Handle error states
  if (shouldShowError) {
    return (
      <ApplicationsPage loaded empty={false}>
        <RedirectErrorState title="Could not load component state" />
      </ApplicationsPage>
    );
  }

  if (shouldRedirect) {
    return <Navigate to="/modelRegistrySettings" replace />;
  }

  return (
    <ApplicationsPage
      title={`Manage ${mrName ?? ''} permissions`}
      description="Manage access to this model registry for individual users and user groups, and for service accounts in a namespace."
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
          <Tab label="Namespaces" />
        </Tabs>
      </Box>
      <Box sx={{ pt: 2 }}>
        {activeTabKey === 0 && (
          <RoleBindingPermissions
            ownerReference={ownerReference}
            roleBindingPermissionsRB={{ ...roleBindings, data: filteredRoleBindings }}
            groups={groups}
            createRoleBinding={createUserRoleBinding}
            deleteRoleBinding={deleteUserRoleBinding}
            projectName={modelRegistryNamespace}
            permissionOptions={userPermissionOptions}
            description="To enable access for all cluster users, add system:authenticated to the group list."
            roleRefKind="ClusterRole"
            roleRefName={userRoleRefName}
          />
        )}
        {activeTabKey === 1 && (
          <RoleBindingPermissions
            ownerReference={ownerReference}
            roleBindingPermissionsRB={{ ...roleBindings, data: filteredNamespaceRoleBindings }}
            groups={[]} // Namespaces don't use groups
            createRoleBinding={createNamespaceRoleBinding}
            deleteRoleBinding={deleteNamespaceRoleBinding}
            projectName={modelRegistryNamespace}
            permissionOptions={namespacePermissionOptions}
            description="Grant access to model registry for service accounts within specific namespaces."
            roleRefKind="ClusterRole"
            roleRefName={namespaceRoleRefName}
            isProjectSubject
          />
        )}
      </Box>
    </ApplicationsPage>
  );
};

export default ModelRegistryManagePermissions;

import React from 'react';
import { Navigate } from 'react-router-dom';
import {
  Tabs,
  Tab,
  TabTitleText,
  EmptyState,
  EmptyStateBody,
  EmptyStateVariant,
  Title,
} from '@patternfly/react-core';
import { ExclamationCircleIcon } from '@patternfly/react-icons';
import RoleBindingPermissions from '~/app/pages/settings/roleBinding/RoleBindingPermissions';
import { useModelRegistryPermissionsLogic } from './useModelRegistryPermissionsLogic';

const ModelRegistriesPermissions: React.FC = () => {
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

  if (shouldShowError) {
    return (
      <EmptyState
        headingLevel="h2"
        icon={ExclamationCircleIcon}
        titleText="There was an issue loading permissions."
        variant={EmptyStateVariant.lg}
        data-id="error-empty-state"
        id="permissions"
      >
        <EmptyStateBody>
          Unable to load model registry permissions. Refresh the page to try again.
        </EmptyStateBody>
      </EmptyState>
    );
  }

  if (shouldRedirect) {
    return <Navigate to="/modelRegistrySettings" replace />;
  }

  return (
    <>
      <Title headingLevel="h2" size="xl">
        Manage {mrName} permissions
      </Title>
      <Tabs
        activeKey={activeTabKey}
        onSelect={(event, tabIndex) => setActiveTabKey(Number(tabIndex))}
        usePageInsets
        id="manage-permissions-tabs"
      >
        <Tab eventKey={0} title={<TabTitleText>Users</TabTitleText>}>
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
        </Tab>
        <Tab eventKey={1} title={<TabTitleText>Namespaces</TabTitleText>}>
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
        </Tab>
      </Tabs>
    </>
  );
};

export default ModelRegistriesPermissions;

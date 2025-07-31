import React from 'react';
import { Breadcrumb, BreadcrumbItem, PageSection, Tab, Tabs } from '@patternfly/react-core';
import { Link, Navigate } from 'react-router-dom';
import { ApplicationsPage } from 'mod-arch-shared';
import RoleBindingPermissions from '~/app/pages/settings/roleBinding/RoleBindingPermissions';
import RedirectErrorState from '~/app/shared/components/RedirectErrorState';
import { useModelRegistryPermissionsLogic } from './useModelRegistryPermissionsLogic';

const ModelRegistriesManagePermissions: React.FC = () => {
  const {
    activeTabKey,
    setActiveTabKey,
    ownerReference,
    groups,
    filteredRoleBindings,
    filteredProjectRoleBindings,
    mrName,
    modelRegistryNamespace,
    roleBindings,
    userPermissionOptions,
    projectPermissionOptions,
    createUserRoleBinding,
    deleteUserRoleBinding,
    createProjectRoleBinding,
    deleteProjectRoleBinding,
    userRoleRefName,
    projectRoleRefName,
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
      description="Manage access to this model registry for individual users and user groups, and for service accounts in a project."
      breadcrumb={
        <Breadcrumb>
          <BreadcrumbItem>
            <Link to="/modelRegistrySettings">Model registry settings</Link>
          </BreadcrumbItem>
          <BreadcrumbItem isActive>Manage Permissions</BreadcrumbItem>
        </Breadcrumb>
      }
      loaded
      empty={false}
    >
      <Tabs activeKey={activeTabKey} onSelect={(_e, tabIndex) => setActiveTabKey(Number(tabIndex))}>
        <Tab eventKey={0} title="Users" />
        <Tab eventKey={1} title="Projects" />
      </Tabs>
      <PageSection isFilled>
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
            roleRefKind="Role"
            roleRefName={userRoleRefName}
          />
        )}
        {activeTabKey === 1 && (
          <RoleBindingPermissions
            ownerReference={ownerReference}
            roleBindingPermissionsRB={{ ...roleBindings, data: filteredProjectRoleBindings }}
            groups={[]} // Projects don't use groups
            createRoleBinding={createProjectRoleBinding}
            deleteRoleBinding={deleteProjectRoleBinding}
            projectName={modelRegistryNamespace}
            permissionOptions={projectPermissionOptions}
            description="Grant access to model registry for service accounts within specific projects."
            roleRefKind="Role"
            roleRefName={projectRoleRefName}
            isProjectSubject
          />
        )}
      </PageSection>
    </ApplicationsPage>
  );
};

export default ModelRegistriesManagePermissions;

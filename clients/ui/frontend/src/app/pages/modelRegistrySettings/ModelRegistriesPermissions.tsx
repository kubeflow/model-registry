import React from 'react';
import { Breadcrumb, BreadcrumbItem, PageSection, Tab, Tabs } from '@patternfly/react-core';
import { Link, Navigate, useParams } from 'react-router-dom';
import {
  ApplicationsPage,
  ModelRegistryKind,
  RoleBindingKind,
  useQueryParamNamespaces,
} from 'mod-arch-shared';
import RoleBindingPermissions from '~/app/pages/settings/roleBinding/RoleBindingPermissions';
import { useModelRegistryCR } from '~/app/hooks/useModelRegistryCR';
import useModelRegistryRoleBindings from '~/app/pages/modelRegistrySettings/useModelRegistryRoleBindings';
import { RoleBindingPermissionsRoleType } from '~/app/pages/settings/roleBinding/types';
import RedirectErrorState from '~/app/shared/components/RedirectErrorState';
import {
  createModelRegistryRoleBindingWrapper,
  deleteModelRegistryRoleBindingWrapper,
} from '~/app/pages/settings/roleBindingUtils';
import { useGroups } from '~/app/hooks/useGroups';

const ModelRegistriesManagePermissions: React.FC = () => {
  const [activeTabKey, setActiveTabKey] = React.useState(0);
  const modelRegistryNamespace = 'model-registry'; //TODO this is a placeholder
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
      </PageSection>
    </ApplicationsPage>
  );
};

export default ModelRegistriesManagePermissions;

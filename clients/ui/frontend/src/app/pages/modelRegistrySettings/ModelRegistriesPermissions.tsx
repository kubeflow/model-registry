import React from 'react';
import {
  Breadcrumb,
  BreadcrumbItem,
  PageSection,
  Tab,
  Tabs,
} from '@patternfly/react-core';
import { Link, Navigate, useParams } from 'react-router-dom';
import { ApplicationsPage } from 'mod-arch-shared';
import { ModelRegistryKind, RoleBindingKind } from '~/app/k8sTypes';
import { useGroups } from '~/app/api/k8s/groups';
import RoleBindingPermissions from '~/app/pages/settings/roleBinding/RoleBindingPermissions';
import { useModelRegistryNamespaceCR } from '~/app/concepts/modelRegistry/context/useModelRegistryNamespaceCR';
// import { AreaContext } from '~/app/concepts/areas/AreaContext';
import {
  createModelRegistryRoleBinding,
  deleteModelRegistryRoleBinding,
} from '~/app/services/modelRegistrySettingsService';
import useModelRegistryRoleBindings from '~/app/pages/modelRegistrySettings/useModelRegistryRoleBindings';
import { RoleBindingPermissionsRoleType } from '~/app/pages/settings/roleBinding/types';
import RedirectErrorState from '~/app/pages/external/RedirectErrorState';

const ModelRegistriesManagePermissions: React.FC = () => {
  // const { dscStatus } = React.useContext(AreaContext);
  // const modelRegistryNamespace = dscStatus?.components?.modelregistry?.registriesNamespace;
  const modelRegistryNamespace = 'model-registry'; // TODO: This is a placeholder
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
      <Tabs activeKey={activeTabKey} onSelect={(_e, tabIndex) => setActiveTabKey(tabIndex as number)}>
        <Tab eventKey={0} title="Users" />
        <Tab eventKey={1} title="Projects" />
      </Tabs>
      <PageSection isFilled>
        {activeTabKey === 0 && (
          <RoleBindingPermissions
            ownerReference={ownerReference}
            roleBindingPermissionsRB={{ ...roleBindings, data: filteredRoleBindings }}
            groups={groups}
            createRoleBinding={createModelRegistryRoleBinding}
            deleteRoleBinding={deleteModelRegistryRoleBinding}
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

import * as React from 'react';
import {
  Alert,
  EmptyState,
  EmptyStateBody,
  EmptyStateVariant,
  PageSection,
  Spinner,
  Stack,
  StackItem,
} from '@patternfly/react-core';
import { ExclamationCircleIcon } from '@patternfly/react-icons';
import {
  FetchStateObject,
  GroupKind,
  K8sResourceCommon,
  K8sStatus,
  RoleBindingKind,
  RoleBindingRoleRef,
} from 'mod-arch-shared';
import RoleBindingPermissionsTableSection from './RoleBindingPermissionsTableSection';
import { RoleBindingPermissionsRBType, RoleBindingPermissionsRoleType } from './types';
import { filterRoleBindingSubjects } from './utils';

type RoleBindingPermissionsProps = {
  ownerReference?: K8sResourceCommon;
  roleBindingPermissionsRB: FetchStateObject<RoleBindingKind[]>;
  defaultRoleBindingName?: string;
  permissionOptions: {
    type: RoleBindingPermissionsRoleType;
    description: string;
  }[];
  createRoleBinding: (roleBinding: RoleBindingKind) => Promise<RoleBindingKind>;
  deleteRoleBinding: (name: string, namespace: string) => Promise<K8sStatus>;
  projectName: string;
  roleRefKind: RoleBindingRoleRef['kind'];
  roleRefName?: RoleBindingRoleRef['name'];
  labels?: { [key: string]: string };
  description: React.ReactElement | string;
  groups: GroupKind[];
  isGroupFirst?: boolean;
  isProjectSubject?: boolean;
};

const RoleBindingPermissions: React.FC<RoleBindingPermissionsProps> = ({
  ownerReference,
  roleBindingPermissionsRB,
  defaultRoleBindingName,
  permissionOptions,
  projectName,
  createRoleBinding,
  deleteRoleBinding,
  roleRefKind,
  roleRefName,
  labels,
  description,
  groups,
  isGroupFirst = false,
  isProjectSubject = false,
}) => {
  const {
    data: roleBindings,
    loaded,
    error: loadError,
    refresh: refreshRB,
  } = roleBindingPermissionsRB;
  if (loadError) {
    return (
      <EmptyState
        headingLevel="h2"
        icon={ExclamationCircleIcon}
        titleText="There was an issue loading permissions."
        variant={EmptyStateVariant.lg}
        data-id="error-empty-state"
        id="permissions"
      >
        <EmptyStateBody>{loadError.message}</EmptyStateBody>
      </EmptyState>
    );
  }

  if (!loaded) {
    return (
      <EmptyState
        headingLevel="h2"
        titleText="Loading"
        variant={EmptyStateVariant.lg}
        data-id="loading-empty-state"
        id="permissions"
      >
        <Spinner size="xl" />
      </EmptyState>
    );
  }

  const userTable = (
    <RoleBindingPermissionsTableSection
      ownerReference={ownerReference}
      defaultRoleBindingName={defaultRoleBindingName}
      projectName={projectName}
      roleRefKind={roleRefKind}
      roleRefName={roleRefName}
      labels={labels}
      permissionOptions={permissionOptions}
      roleBindings={filterRoleBindingSubjects(roleBindings, RoleBindingPermissionsRBType.USER)}
      subjectKind={RoleBindingPermissionsRBType.USER}
      refresh={refreshRB}
      typeModifier="user"
      createRoleBinding={createRoleBinding}
      deleteRoleBinding={deleteRoleBinding}
      isProjectSubject={isProjectSubject}
    />
  );

  const groupTable = (
    <RoleBindingPermissionsTableSection
      ownerReference={ownerReference}
      defaultRoleBindingName={defaultRoleBindingName}
      projectName={projectName}
      roleRefKind={roleRefKind}
      roleRefName={roleRefName}
      permissionOptions={permissionOptions}
      labels={labels}
      roleBindings={filterRoleBindingSubjects(roleBindings, RoleBindingPermissionsRBType.GROUP)}
      subjectKind={RoleBindingPermissionsRBType.GROUP}
      refresh={refreshRB}
      typeAhead={
        groups.length > 0 ? groups.map((group: GroupKind) => group.metadata.name) : undefined
      }
      typeModifier="group"
      createRoleBinding={createRoleBinding}
      deleteRoleBinding={deleteRoleBinding}
      isProjectSubject={isProjectSubject}
    />
  );

  return (
    <PageSection
      hasBodyWrapper={false}
      isFilled
      aria-label="project-sharing-page-section"
      id="permissions"
    >
      <Stack hasGutter>
        <Alert variant="warning" title="Warning" isInline>
          Changing user or group permissions may remove their access to this resource.
        </Alert>
        <StackItem>{description}</StackItem>
        <StackItem>{isGroupFirst ? groupTable : userTable}</StackItem>
        <StackItem>{isGroupFirst ? userTable : groupTable}</StackItem>
      </Stack>
    </PageSection>
  );
};

export default RoleBindingPermissions;

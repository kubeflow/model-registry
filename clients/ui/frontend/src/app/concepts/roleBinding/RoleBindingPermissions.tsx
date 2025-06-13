import * as React from 'react';
import {
  Alert,
  Box,
  CircularProgress,
  Stack,
  Typography,
} from '@mui/material';
import { Error } from '@mui/icons-material';
import { K8sResourceCommon, K8sStatus } from '@openshift/dynamic-plugin-sdk-utils';
import { GroupKind, RoleBindingKind, RoleBindingRoleRef } from '~/app/k8sTypes';
import { FetchState } from '~/app/utils/useFetch';
import RoleBindingPermissionsTableSection from '~/app/pages/settings/roleBinding/RoleBindingPermissionsTableSection';
import { RoleBindingPermissionsRBType, RoleBindingPermissionsRoleType } from '~/app/pages/settings/roleBinding/types';
import { filterRoleBindingSubjects, tryPatchRoleBinding } from '~/app/concepts/roleBinding/utils';

type RoleBindingPermissionsProps = {
  ownerReference?: K8sResourceCommon;
  roleBindingPermissionsRB: FetchState<RoleBindingKind[]>;
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
}) => {
  const {
    data: roleBindings,
    loaded,
    error: loadError,
    refresh: refreshRB,
  } = roleBindingPermissionsRB;
  if (loadError) {
    return (
      <Box sx={{ textAlign: 'center' }}>
        <Error color="error" sx={{ fontSize: 48 }} />
        <Typography variant="h5" component="h2">
            There was an issue loading permissions.
        </Typography>
        <Typography variant="body1">{loadError.message}</Typography>
      </Box>
    );
  }

  if (!loaded) {
    return (
        <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
            <CircularProgress />
        </Box>
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
      tryPatchRoleBinding={tryPatchRoleBinding}
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
      tryPatchRoleBinding={tryPatchRoleBinding}
    />  
  );

  return (
    <Stack spacing={2}>
        <Alert severity="warning">
            Changing user or group permissions may remove their access to this resource.
        </Alert>
        <Typography>{description}</Typography>
        {isGroupFirst ? groupTable : userTable}
        {isGroupFirst ? userTable : groupTable}
    </Stack>
  );
};

export default RoleBindingPermissions; 
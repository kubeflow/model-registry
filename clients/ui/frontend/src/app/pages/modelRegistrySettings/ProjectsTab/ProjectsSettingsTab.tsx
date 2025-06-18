import { K8sResourceCommon } from 'mod-arch-shared';
import { Box, CircularProgress, Stack, Typography } from '@mui/material';
import { Error } from '@mui/icons-material';
import React from 'react';
import {
  RoleBindingPermissionsRBType,
  RoleBindingPermissionsRoleType,
} from '~/app/pages/settings/roleBinding/types';
import {
  filterRoleBindingSubjects,
  removePrefix,
  tryPatchRoleBinding,
} from '~/app/pages/settings/roleBinding/utils';
import { RoleBindingKind, RoleBindingRoleRef, ProjectKind } from '~/app/k8sTypes';
import { ProjectsContext } from '~/app/concepts/projects/ProjectsContext';
import RoleBindingPermissionsTableSection from '~/app/pages/settings/roleBinding/RoleBindingPermissionsTableSection';
import {
  createModelRegistryRoleBinding,
  deleteModelRegistryRoleBinding,
} from '~/app/services/modelRegistrySettingsService';
import { FetchState } from '~/app/utils/useFetch';

type RoleBindingProjectPermissionsProps = {
  ownerReference?: K8sResourceCommon;
  roleBindingPermissionsRB: FetchState<RoleBindingKind[]>;
  permissionOptions: {
    type: RoleBindingPermissionsRoleType;
    description: string;
  }[];
  projectName: string;
  roleRefName?: RoleBindingRoleRef['name'];
  labels?: { [key: string]: string };
  isProjectSubject?: boolean;
  description: string;
};

const ProjectsSettingsTab: React.FC<RoleBindingProjectPermissionsProps> = ({
  ownerReference,
  roleBindingPermissionsRB,
  permissionOptions,
  projectName,
  roleRefName,
  labels,
  description,
}) => {
  const {
    data: roleBindings,
    loaded,
    error: loadError,
    refresh: refreshRB,
  } = roleBindingPermissionsRB;

  const { projects } = React.useContext(ProjectsContext) as { projects: ProjectKind[] };
  const filteredProjects = projects.filter(
    (project: ProjectKind) => !removePrefix(roleBindings).includes(project.metadata.name),
  );

  if (loadError) {
    return (
      <Box sx={{ textAlign: 'center' }}>
        <Error color="error" sx={{ fontSize: 48 }} />
        <Typography variant="h5" component="h2">
          There was an issue loading projects
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

  return (
    <Stack spacing={2}>
      <Typography>{description}</Typography>
      <RoleBindingPermissionsTableSection
        ownerReference={ownerReference}
        roleBindings={filterRoleBindingSubjects(roleBindings, RoleBindingPermissionsRBType.GROUP)}
        projectName={projectName}
        roleRefKind="Role"
        roleRefName={roleRefName}
        labels={labels}
        subjectKind={RoleBindingPermissionsRBType.GROUP}
        permissionOptions={permissionOptions}
        typeAhead={
          filteredProjects.length > 0
            ? filteredProjects.map((project: ProjectKind) => project.metadata.name)
            : undefined
        }
        refresh={refreshRB}
        typeModifier="project"
        createRoleBinding={createModelRegistryRoleBinding}
        deleteRoleBinding={deleteModelRegistryRoleBinding}
        tryPatchRoleBinding={tryPatchRoleBinding}
      />
    </Stack>
  );
};

export default ProjectsSettingsTab;

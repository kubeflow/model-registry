import { ProjectKind } from '~/app/k8sTypes';
import { getDisplayNameFromK8sResource, isOOTB } from '~/app/concepts/k8s/utils';

export const namespaceToProjectDisplayName = (namespace: string, projects: ProjectKind[]): string =>
  projects.find((p) => p.metadata.name === namespace)?.metadata.annotations?.[
    'openshift.io/display-name'
  ] || namespace;

export const isAvailableProject = (projectName: string, dashboardNamespace: string): boolean =>
  projectName !== dashboardNamespace;

export const isProjectSharing = (project: ProjectKind): boolean =>
  project.metadata.annotations?.['opendatahub.io/project-sharing'] === 'true';

export const getProjectOwner = (project: ProjectKind): string =>
  project.metadata.annotations?.['openshift.io/requester'] || '';

export const getProjectCreationTime = (project: ProjectKind): number => {
  const time = project.metadata.creationTimestamp;
  if (!time) {
    return 0;
  }
  return new Date(time).getTime();
};

export const getProjectDisplayName = (project: ProjectKind): string =>
  getDisplayNameFromK8sResource(project);

export const sortProjectsByDisplayName = (p1: ProjectKind, p2: ProjectKind): number =>
  getProjectDisplayName(p1).localeCompare(getProjectDisplayName(p2));

/**
 * An OOTB project is one that has the label `platform.opendatahub.io/part-of`.
 * The data-science-cluster project and the dashboard's namespace are not considered OOTB projects.
 */
export const isOotbProject = (project: ProjectKind, dashboardNamespace: string): boolean =>
  isOOTB(project) && isAvailableProject(project.metadata.name, dashboardNamespace);

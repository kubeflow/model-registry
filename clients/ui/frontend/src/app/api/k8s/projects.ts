import { k8sListResource } from 'mod-arch-shared';
import { ProjectKind } from '~/app/k8sTypes';
import { ProjectModel } from '~/app/api/models';

export const listProjects = (): Promise<ProjectKind[]> =>
  k8sListResource<ProjectKind>({
    model: ProjectModel,
    queryOptions: {},
  }).then((projects) => projects.items);

import { GroupKind } from '~/app/k8sTypes';

export const useGroups = (): [GroupKind[], boolean, Error | undefined] => [[], true, undefined];

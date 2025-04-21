import { Namespace } from 'mod-arch-shared';

type MockNamespace = {
  name?: string;
};

export const mockNamespace = ({ name = 'kubeflow' }: MockNamespace): Namespace => ({
  name,
});

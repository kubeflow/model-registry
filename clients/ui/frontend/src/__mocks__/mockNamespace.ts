import { Namespace } from 'mod-arch-core';

type MockNamespace = {
  name?: string;
};

export const mockNamespace = ({ name = 'kubeflow' }: MockNamespace): Namespace => ({
  name,
});

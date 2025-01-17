import { Namespace } from '~/shared/types';

type MockNamespace = {
  name?: string;
};

export const mockNamespace = ({ name = 'kubeflow' }: MockNamespace): Namespace => ({
  name,
});

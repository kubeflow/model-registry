import { ModelRegistry } from '~/app/types';

type MockModelRegistry = {
  name?: string;
  description?: string;
  displayName?: string;
};

export const mockModelRegistry = ({
  name = 'modelregistry-sample',
  description = 'New model registry',
  displayName = 'Model Registry Sample',
}: MockModelRegistry): ModelRegistry => ({
  name,
  description,
  displayName,
});

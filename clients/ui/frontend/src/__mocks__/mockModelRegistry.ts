import { ModelRegistry } from '~/app/types';

type MockModelRegistry = {
  name?: string;
  description?: string;
  displayName?: string;
  isAvailable?: boolean;
};

export const mockModelRegistry = ({
  name = 'modelregistry-sample',
  description = 'Model registry description',
  displayName = 'Model Registry Sample',
  isAvailable = true,
}: MockModelRegistry): ModelRegistry => ({
  name,
  description,
  displayName,
  isAvailable,
});

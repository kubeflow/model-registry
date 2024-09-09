/* eslint-disable camelcase */
import { ModelRegistryResponse } from '~/app/types';

export const mockModelRegistryResponse = ({
  model_registry = [],
}: Partial<ModelRegistryResponse>): ModelRegistryResponse => ({
  model_registry,
});

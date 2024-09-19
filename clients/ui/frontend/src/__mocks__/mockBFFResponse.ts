import { ModelRegistryResponse } from '~/app/types';

export const mockBFFResponse = <T>(data: T): ModelRegistryResponse<T> => ({
  data,
});

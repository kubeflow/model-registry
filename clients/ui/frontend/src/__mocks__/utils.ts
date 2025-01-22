import { ModelRegistryBody } from '~/app/types';

export const mockBFFResponse = <T>(data: T): ModelRegistryBody<T> => ({
  data,
});

import { APIError } from '~/shared/api/types';
import { isCommonStateError } from '~/shared/utilities/useFetchState';

const isError = (e: unknown): e is APIError => typeof e === 'object' && e !== null && 'error' in e;

export const handleRestFailures = <T>(promise: Promise<T>): Promise<T> =>
  promise
    .then((result) => {
      if (isError(result)) {
        throw result;
      }
      return result;
    })
    .catch((e) => {
      if (isError(e)) {
        throw new Error(e.error.message);
      }
      if (isCommonStateError(e)) {
        // Common state errors are handled by useFetchState at storage level, let them deal with it
        throw e;
      }
      // eslint-disable-next-line no-console
      console.error('Unknown API error', e);
      throw new Error('Error communicating with server');
    });

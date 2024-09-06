import { NotReadyError } from '~/utilities/useFetchState';
import { APIError } from '~/types';
import { handleRestFailures } from '~/app/api/errorUtils';
import { mockRegisteredModel } from '~/__mocks__/mockRegisteredModel';

describe('handleRestFailures', () => {
  it('should successfully return registered models', async () => {
    const modelRegistryMock = mockRegisteredModel({});
    const result = await handleRestFailures(Promise.resolve(modelRegistryMock));
    expect(result).toStrictEqual(modelRegistryMock);
  });

  it('should handle and throw model registry errors', async () => {
    const statusMock: APIError = {
      code: '',
      message: 'error',
    };

    await expect(handleRestFailures(Promise.resolve(statusMock))).rejects.toThrow('error');
  });

  it('should handle common state errors ', async () => {
    await expect(handleRestFailures(Promise.reject(new NotReadyError('error')))).rejects.toThrow(
      'error',
    );
  });

  it('should handle other errors', async () => {
    await expect(handleRestFailures(Promise.reject(new Error('error')))).rejects.toThrow(
      'Error communicating with server',
    );
  });
});

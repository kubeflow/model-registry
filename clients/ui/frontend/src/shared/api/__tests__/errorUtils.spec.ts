import { NotReadyError } from '~/shared/utilities/useFetchState';
import { APIError } from '~/shared/api/types';
import { handleRestFailures } from '~/shared/api/errorUtils';
import { mockRegisteredModel } from '~/__mocks__/mockRegisteredModel';
import { mockBFFResponse } from '~/__mocks__/utils';

describe('handleRestFailures', () => {
  it('should successfully return registered models', async () => {
    const modelRegistryMock = mockRegisteredModel({});
    const result = await handleRestFailures(Promise.resolve(mockBFFResponse(modelRegistryMock)));
    expect(result.data).toStrictEqual(modelRegistryMock);
  });

  it('should handle and throw model registry errors', async () => {
    const statusMock: APIError = {
      error: {
        code: '',
        message: 'error',
      },
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

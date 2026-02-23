import { waitFor } from '@testing-library/react';
import { useCheckNamespaceRegistryAccess } from '~/app/hooks/useCheckNamespaceRegistryAccess';
import { checkNamespaceRegistryAccess } from '~/app/api/k8s';
import { testHook } from '~/__tests__/unit/testUtils/hooks';

jest.mock('~/app/api/k8s', () => ({
  checkNamespaceRegistryAccess: jest.fn(),
}));

const mockCheckNamespaceRegistryAccess = jest.mocked(checkNamespaceRegistryAccess);

describe('useCheckNamespaceRegistryAccess', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockCheckNamespaceRegistryAccess.mockReturnValue(
      jest.fn().mockResolvedValue({ hasAccess: true }),
    );
  });

  it('should return undefined hasAccess and no loading when jobNamespace is missing', async () => {
    const { result } = testHook(useCheckNamespaceRegistryAccess)('my-registry', 'registry-ns', '');

    expect(result.current.hasAccess).toBeUndefined();
    expect(result.current.isLoading).toBe(false);
    expect(result.current.error).toBeUndefined();
    expect(mockCheckNamespaceRegistryAccess).not.toHaveBeenCalled();
  });

  it('should return undefined hasAccess when registryName is missing', async () => {
    const { result } = testHook(useCheckNamespaceRegistryAccess)(
      undefined,
      'registry-ns',
      'job-ns',
    );

    expect(result.current.hasAccess).toBeUndefined();
    expect(mockCheckNamespaceRegistryAccess).not.toHaveBeenCalled();
  });

  it('should return undefined hasAccess when registryNamespace is missing', async () => {
    const { result } = testHook(useCheckNamespaceRegistryAccess)(
      'my-registry',
      undefined,
      'job-ns',
    );

    expect(result.current.hasAccess).toBeUndefined();
    expect(mockCheckNamespaceRegistryAccess).not.toHaveBeenCalled();
  });

  it('should call API and set hasAccess to true when all params are provided', async () => {
    const apiMock = jest.fn().mockResolvedValue({ hasAccess: true });
    mockCheckNamespaceRegistryAccess.mockReturnValue(apiMock);

    const { result } = testHook(useCheckNamespaceRegistryAccess)(
      'my-registry',
      'registry-ns',
      'job-ns',
    );

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    expect(apiMock).toHaveBeenCalledWith(
      {},
      {
        namespace: 'job-ns',
        registryName: 'my-registry',
        registryNamespace: 'registry-ns',
      },
    );
    expect(result.current.hasAccess).toBe(true);
    expect(result.current.error).toBeUndefined();
  });

  it('should set hasAccess to false when API returns hasAccess false', async () => {
    const apiMock = jest.fn().mockResolvedValue({ hasAccess: false });
    mockCheckNamespaceRegistryAccess.mockReturnValue(apiMock);

    const { result } = testHook(useCheckNamespaceRegistryAccess)(
      'my-registry',
      'registry-ns',
      'job-ns',
    );

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    expect(result.current.hasAccess).toBe(false);
  });

  it('should set error when API throws', async () => {
    const apiError = new Error('Network error');
    const apiMock = jest.fn().mockRejectedValue(apiError);
    mockCheckNamespaceRegistryAccess.mockReturnValue(apiMock);

    const { result } = testHook(useCheckNamespaceRegistryAccess)(
      'my-registry',
      'registry-ns',
      'job-ns',
    );

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeDefined();
    });

    expect(result.current.error).toEqual(apiError);
    expect(result.current.hasAccess).toBeUndefined();
  });
});

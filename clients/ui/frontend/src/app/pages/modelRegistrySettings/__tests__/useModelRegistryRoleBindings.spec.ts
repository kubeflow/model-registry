import { waitFor } from '@testing-library/react';
import { RoleBindingKind } from 'mod-arch-shared';
import { useFetchState, POLL_INTERVAL, useDeepCompareMemoize } from 'mod-arch-core';
import { testHook } from '~/__tests__/unit/testUtils/hooks';
import * as k8sApi from '~/app/api/k8s';
import useModelRegistryRoleBindings from '~/app/pages/modelRegistrySettings/useModelRegistryRoleBindings';

// Mock mod-arch-core
jest.mock('mod-arch-core', () => ({
  useFetchState: jest.fn(),
  useDeepCompareMemoize: jest.fn(),
  POLL_INTERVAL: 5000,
  asEnumMember: jest.fn((value, enumObj) => {
    if (!enumObj) {
      return value;
    }
    return value || Object.values(enumObj)[0];
  }),
  DeploymentMode: {
    Federated: 'federated',
    Standalone: 'standalone',
  },
}));

// Mock mod-arch-shared
jest.mock('mod-arch-shared', () => ({
  asEnumMember: jest.fn((value, enumObj) => value || Object.values(enumObj)[0]),
  Theme: {
    Patternfly: 'patternfly',
    Material: 'material',
  },
  DeploymentMode: {
    Federated: 'federated',
    Standalone: 'standalone',
  },
}));

// Mock the k8s API
jest.mock('~/app/api/k8s');

const mockUseFetchState = jest.mocked(useFetchState);
const mockUseDeepCompareMemoize = jest.mocked(useDeepCompareMemoize);
const mockK8sApi = jest.mocked(k8sApi);

const mockRoleBinding: RoleBindingKind = {
  apiVersion: 'rbac.authorization.k8s.io/v1',
  kind: 'RoleBinding',
  metadata: {
    name: 'test-role-binding',
    namespace: 'test-namespace',
    labels: {
      'app.kubernetes.io/name': 'test-app',
    },
  },
  subjects: [
    {
      kind: 'User',
      name: 'test-user',
      apiGroup: 'rbac.authorization.k8s.io',
    },
  ],
  roleRef: {
    kind: 'ClusterRole',
    name: 'test-role',
    apiGroup: 'rbac.authorization.k8s.io',
  },
};

const mockQueryParams = {
  namespace: 'test-namespace',
  // eslint-disable-next-line camelcase
  some_param: 'test-value',
};

describe('useModelRegistryRoleBindings', () => {
  const mockRefresh = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
    mockUseDeepCompareMemoize.mockReturnValue(mockQueryParams);
  });

  it('should return initial state when loading', () => {
    mockUseFetchState.mockReturnValue([[], false, undefined, mockRefresh]);

    const { result } = testHook(useModelRegistryRoleBindings)(mockQueryParams);

    expect(result.current).toEqual({
      data: [],
      loaded: false,
      error: undefined,
      refresh: mockRefresh,
    });
  });

  it('should return role bindings data when loaded successfully', async () => {
    const mockData = [mockRoleBinding];
    mockUseFetchState.mockReturnValue([mockData, true, undefined, mockRefresh]);

    const { result } = testHook(useModelRegistryRoleBindings)(mockQueryParams);

    await waitFor(() => {
      expect(result.current).toEqual({
        data: mockData,
        loaded: true,
        error: undefined,
        refresh: mockRefresh,
      });
    });
  });

  it('should return error state when fetch fails', async () => {
    const mockError = new Error('Failed to fetch role bindings');
    mockUseFetchState.mockReturnValue([[], false, mockError, mockRefresh]);

    const { result } = testHook(useModelRegistryRoleBindings)(mockQueryParams);

    await waitFor(() => {
      expect(result.current).toEqual({
        data: [],
        loaded: false,
        error: mockError,
        refresh: mockRefresh,
      });
    });
  });

  it('should use deep compare memoization for query params', () => {
    mockUseFetchState.mockReturnValue([[], false, undefined, mockRefresh]);

    testHook(useModelRegistryRoleBindings)(mockQueryParams);

    expect(mockUseDeepCompareMemoize).toHaveBeenCalledWith(mockQueryParams);
  });

  it('should configure useFetchState with correct options', () => {
    mockUseFetchState.mockReturnValue([[], false, undefined, mockRefresh]);

    testHook(useModelRegistryRoleBindings)(mockQueryParams);

    expect(mockUseFetchState).toHaveBeenCalledWith(
      expect.any(Function), // fetchRoleBindings callback
      [], // initial data
      { refreshRate: POLL_INTERVAL },
    );
  });

  it('should update when query params change', () => {
    const initialParams = { namespace: 'initial' };
    const updatedParams = { namespace: 'updated' };

    mockUseDeepCompareMemoize
      .mockReturnValueOnce(initialParams)
      .mockReturnValueOnce({ $params: updatedParams });

    mockUseFetchState.mockReturnValue([[], false, undefined, mockRefresh]);

    const { rerender } = testHook(useModelRegistryRoleBindings)(initialParams);

    expect(mockUseDeepCompareMemoize).toHaveBeenCalledWith(initialParams);

    rerender({ $params: updatedParams });

    // The hook should be called again with the updated params wrapped in $params
    expect(mockUseDeepCompareMemoize).toHaveBeenCalledTimes(2);
    expect(mockUseDeepCompareMemoize).toHaveBeenNthCalledWith(2, { $params: updatedParams });
  });

  it('should call getRoleBindings with memoized params', () => {
    const mockGetRoleBindingsInner = jest.fn().mockResolvedValue([]);
    mockK8sApi.getRoleBindings.mockReturnValue(() => mockGetRoleBindingsInner());

    mockUseFetchState.mockReturnValue([[], false, undefined, mockRefresh]);

    testHook(useModelRegistryRoleBindings)(mockQueryParams);

    // Verify that the callback was created with the correct dependencies
    expect(mockUseFetchState).toHaveBeenCalledWith(expect.any(Function), [], {
      refreshRate: POLL_INTERVAL,
    });

    // The callback should use the memoized params
    expect(mockUseDeepCompareMemoize).toHaveBeenCalledWith(mockQueryParams);
  });
});

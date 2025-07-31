import { RoleBindingKind, K8sStatus } from 'mod-arch-shared';
import * as k8sApi from '~/app/api/k8s';
import {
  createModelRegistryRoleBindingWrapper,
  deleteModelRegistryRoleBindingWrapper,
  createModelRegistryProjectRoleBinding,
  deleteModelRegistryProjectRoleBinding,
} from '~/app/pages/settings/roleBindingUtils';

// Mock the k8s API functions
jest.mock('~/app/api/k8s');
const mockedK8sApi = jest.mocked(k8sApi);

// Mock window.location.origin
Object.defineProperty(window, 'location', {
  value: {
    origin: 'https://example.com',
  },
  writable: true,
});

describe('roleBindingUtils', () => {
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

  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('createModelRegistryRoleBindingWrapper', () => {
    it('should call createRoleBinding and return the result', async () => {
      // Mock the curried function: createRoleBinding(host, opts)(apiOpts, payload)
      const mockFinalCall = jest.fn().mockResolvedValue(mockRoleBinding);
      mockedK8sApi.createRoleBinding.mockReturnValue((opts, data) => mockFinalCall(opts, data));

      const result = await createModelRegistryRoleBindingWrapper(mockRoleBinding);

      expect(mockedK8sApi.createRoleBinding).toHaveBeenCalledWith('https://example.com', {});
      expect(mockFinalCall).toHaveBeenCalledWith({}, mockRoleBinding);
      expect(result).toEqual(mockRoleBinding);
    });

    it('should propagate errors from createRoleBinding', async () => {
      const error = new Error('Create failed');
      const mockFinalCall = jest.fn().mockRejectedValue(error);
      mockedK8sApi.createRoleBinding.mockReturnValue((opts, data) => mockFinalCall(opts, data));

      await expect(createModelRegistryRoleBindingWrapper(mockRoleBinding)).rejects.toThrow(
        'Create failed',
      );
    });
  });

  describe('deleteModelRegistryRoleBindingWrapper', () => {
    it('should call deleteRoleBinding and return success status', async () => {
      // Mock the curried function: deleteRoleBinding(host, opts)(apiOpts, name)
      const mockFinalCall = jest.fn().mockResolvedValue(undefined);
      mockedK8sApi.deleteRoleBinding.mockReturnValue((opts, name) => mockFinalCall(opts, name));

      const result = await deleteModelRegistryRoleBindingWrapper(
        'test-role-binding',
        'test-namespace',
      );

      expect(mockedK8sApi.deleteRoleBinding).toHaveBeenCalledWith('https://example.com', {});
      expect(mockFinalCall).toHaveBeenCalledWith({}, 'test-role-binding');

      const expectedStatus: K8sStatus = {
        apiVersion: 'v1',
        kind: 'Status',
        status: 'Success',
        code: 200,
        message: 'Role binding deleted successfully',
        reason: 'Deleted',
      };
      expect(result).toEqual(expectedStatus);
    });

    it('should propagate errors from deleteRoleBinding', async () => {
      const error = new Error('Delete failed');
      const mockFinalCall = jest.fn().mockRejectedValue(error);
      mockedK8sApi.deleteRoleBinding.mockReturnValue((opts, name) => mockFinalCall(opts, name));

      await expect(
        deleteModelRegistryRoleBindingWrapper('test-role-binding', 'test-namespace'),
      ).rejects.toThrow('Delete failed');
    });
  });

  describe('createModelRegistryProjectRoleBinding', () => {
    it('should add project-specific label before calling createRoleBinding', async () => {
      const mockFinalCall = jest.fn().mockResolvedValue(mockRoleBinding);
      mockedK8sApi.createRoleBinding.mockReturnValue((opts, data) => mockFinalCall(opts, data));

      const result = await createModelRegistryProjectRoleBinding(mockRoleBinding);

      expect(mockedK8sApi.createRoleBinding).toHaveBeenCalledWith('https://example.com', {});

      // Verify that the project-specific label was added
      const callArgs = mockFinalCall.mock.calls[0];
      const modifiedRoleBinding = callArgs[1];
      expect(modifiedRoleBinding.metadata.labels['app.kubernetes.io/component']).toBe(
        'model-registry-project-rbac',
      );
      expect(result).toEqual(mockRoleBinding);
    });

    it('should preserve existing labels when adding project label', async () => {
      const roleBindingWithLabels: RoleBindingKind = {
        ...mockRoleBinding,
        metadata: {
          ...mockRoleBinding.metadata,
          labels: {
            'existing.label': 'value',
            'another.label': 'another-value',
          },
        },
      };

      const mockFinalCall = jest.fn().mockResolvedValue(roleBindingWithLabels);
      mockedK8sApi.createRoleBinding.mockReturnValue((opts, data) => mockFinalCall(opts, data));

      await createModelRegistryProjectRoleBinding(roleBindingWithLabels);

      const callArgs = mockFinalCall.mock.calls[0];
      const modifiedRoleBinding = callArgs[1];

      expect(modifiedRoleBinding.metadata.labels['existing.label']).toBe('value');
      expect(modifiedRoleBinding.metadata.labels['another.label']).toBe('another-value');
      expect(modifiedRoleBinding.metadata.labels['app.kubernetes.io/component']).toBe(
        'model-registry-project-rbac',
      );
    });

    it('should handle role binding without existing labels', async () => {
      const roleBindingWithoutLabels: RoleBindingKind = {
        ...mockRoleBinding,
        metadata: {
          name: 'test-role-binding',
          namespace: 'test-namespace',
        },
      };

      const mockFinalCall = jest.fn().mockResolvedValue(roleBindingWithoutLabels);
      mockedK8sApi.createRoleBinding.mockReturnValue((opts, data) => mockFinalCall(opts, data));

      await createModelRegistryProjectRoleBinding(roleBindingWithoutLabels);

      const callArgs = mockFinalCall.mock.calls[0];
      const modifiedRoleBinding = callArgs[1];

      expect(modifiedRoleBinding.metadata.labels['app.kubernetes.io/component']).toBe(
        'model-registry-project-rbac',
      );
    });
  });

  describe('deleteModelRegistryProjectRoleBinding', () => {
    it('should call deleteRoleBinding and return project-specific success status', async () => {
      const mockFinalCall = jest.fn().mockResolvedValue(undefined);
      mockedK8sApi.deleteRoleBinding.mockReturnValue((opts, name) => mockFinalCall(opts, name));

      const result = await deleteModelRegistryProjectRoleBinding(
        'test-project-role-binding',
        'test-namespace',
      );

      expect(mockedK8sApi.deleteRoleBinding).toHaveBeenCalledWith('https://example.com', {});
      expect(mockFinalCall).toHaveBeenCalledWith({}, 'test-project-role-binding');

      const expectedStatus: K8sStatus = {
        apiVersion: 'v1',
        kind: 'Status',
        status: 'Success',
        code: 200,
        message: 'Project role binding deleted successfully',
        reason: 'Deleted',
      };
      expect(result).toEqual(expectedStatus);
    });

    it('should propagate errors from deleteRoleBinding', async () => {
      const error = new Error('Delete project failed');
      const mockFinalCall = jest.fn().mockRejectedValue(error);
      mockedK8sApi.deleteRoleBinding.mockReturnValue((opts, name) => mockFinalCall(opts, name));

      await expect(
        deleteModelRegistryProjectRoleBinding('test-project-role-binding', 'test-namespace'),
      ).rejects.toThrow('Delete project failed');
    });
  });
});

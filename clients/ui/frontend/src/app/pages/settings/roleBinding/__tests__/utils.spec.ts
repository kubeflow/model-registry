import { mockRoleBindingK8sResource } from '~/__mocks__/mockRoleBindingK8sResource';
import { RoleBindingPermissionsRoleType } from '~/app/pages/settings/roleBinding/types';
import { NamespaceKind } from '~/app/shared/components/types';
import {
  castRoleBindingPermissionsRoleType,
  firstSubject,
  isCurrentUserChanging,
  roleLabel,
  namespaceToDisplayName,
  displayNameToNamespace,
} from '~/app/pages/settings/roleBinding/utils';

describe('firstSubject', () => {
  it('should return name', () => {
    const roleBinding = mockRoleBindingK8sResource({
      name: 'test-user',
      subjects: [{ kind: 'User', apiGroup: 'rbac.authorization.k8s.io', name: 'test-user' }],
    });
    const result = firstSubject(roleBinding);
    expect(result).toBe('test-user');
  });
});
describe('isCurrentUserChanging', () => {
  it('should return true when role binding subject matches current username', () => {
    const roleBinding = mockRoleBindingK8sResource({
      name: 'test-user',
      subjects: [{ kind: 'User', apiGroup: 'rbac.authorization.k8s.io', name: 'test-user' }],
    });
    expect(isCurrentUserChanging(roleBinding, 'test-user')).toBe(true);
  });

  it('should return false when role binding subject does not match current username', () => {
    const roleBinding = mockRoleBindingK8sResource({
      name: 'other-user',
      subjects: [{ kind: 'User', apiGroup: 'rbac.authorization.k8s.io', name: 'other-user' }],
    });
    expect(isCurrentUserChanging(roleBinding, 'test-user')).toBe(false);
  });

  it('should return false when role binding is undefined', () => {
    expect(isCurrentUserChanging(undefined, 'test-user')).toBe(false);
  });
});

describe('castRoleBindingPermissionsRoleType', () => {
  it('should return default when role includes registry-user', () => {
    expect(castRoleBindingPermissionsRoleType('registry-user')).toBe(
      RoleBindingPermissionsRoleType.DEFAULT,
    );
  });

  it('should return admin when role is admin', () => {
    expect(castRoleBindingPermissionsRoleType('admin')).toBe(RoleBindingPermissionsRoleType.ADMIN);
  });

  it('should return edit when role is edit', () => {
    expect(castRoleBindingPermissionsRoleType('edit')).toBe(RoleBindingPermissionsRoleType.EDIT);
  });

  it('should return custom when role is not admin, edit, or registry-user', () => {
    expect(castRoleBindingPermissionsRoleType('custom')).toBe(
      RoleBindingPermissionsRoleType.CUSTOM,
    );
  });
});

describe('roleLabel', () => {
  it('should return contributor, when the RoleBindingPermissionsRoleType is Edit', () => {
    const result = roleLabel(RoleBindingPermissionsRoleType.EDIT);
    expect(result).toBe('Contributor');
  });

  it('should return Default, when the RoleBindingPermissionsRoleType is other than default', () => {
    const result = roleLabel(RoleBindingPermissionsRoleType.DEFAULT);
    expect(result).toBe('Default');
  });

  it('should return Custom, when the RoleBindingPermissionsRoleType is other than custom', () => {
    const result = roleLabel(RoleBindingPermissionsRoleType.CUSTOM);
    expect(result).toBe('Custom');
  });

  it('should return Admin, when the RoleBindingPermissionsRoleType is Admin', () => {
    const result = roleLabel(RoleBindingPermissionsRoleType.ADMIN);
    expect(result).toBe('Admin');
  });
});

describe('namespaceToDisplayName', () => {
  const mockNamespaces: NamespaceKind[] = [
    { name: 'default', 'display-name': 'Default Namespace' },
    { name: 'kube-system', 'display-name': 'Kubernetes System' },
    { name: 'my-project', 'display-name': 'My Project' },
  ];

  it('should return display name when namespace exists', () => {
    const result = namespaceToDisplayName('default', mockNamespaces);
    expect(result).toBe('Default Namespace');
  });

  it('should return original namespace name when namespace not found', () => {
    const result = namespaceToDisplayName('non-existent', mockNamespaces);
    expect(result).toBe('non-existent');
  });

  it('should return namespace name when namespaces array is empty', () => {
    const result = namespaceToDisplayName('default', []);
    expect(result).toBe('default');
  });

  it('should handle namespace with same name and display-name', () => {
    const namespaces: NamespaceKind[] = [{ name: 'test', 'display-name': 'test' }];
    const result = namespaceToDisplayName('test', namespaces);
    expect(result).toBe('test');
  });

  it('should handle empty display name gracefully', () => {
    const namespaces: NamespaceKind[] = [{ name: 'test', 'display-name': '' }];
    const result = namespaceToDisplayName('test', namespaces);
    expect(result).toBe('test');
  });
});

describe('displayNameToNamespace', () => {
  const mockNamespaces: NamespaceKind[] = [
    { name: 'default', 'display-name': 'Default Namespace' },
    { name: 'kube-system', 'display-name': 'Kubernetes System' },
    { name: 'my-project', 'display-name': 'My Project' },
  ];

  it('should return namespace name when display name exists', () => {
    const result = displayNameToNamespace('Default Namespace', mockNamespaces);
    expect(result).toBe('default');
  });

  it('should return original display name when namespace not found', () => {
    const result = displayNameToNamespace('Non Existent', mockNamespaces);
    expect(result).toBe('Non Existent');
  });

  it('should return display name when namespaces array is empty', () => {
    const result = displayNameToNamespace('Default Namespace', []);
    expect(result).toBe('Default Namespace');
  });

  it('should handle namespace with same name and display-name', () => {
    const namespaces: NamespaceKind[] = [{ name: 'test', 'display-name': 'test' }];
    const result = displayNameToNamespace('test', namespaces);
    expect(result).toBe('test');
  });

  it('should be case sensitive', () => {
    const result = displayNameToNamespace('default namespace', mockNamespaces);
    expect(result).toBe('default namespace');
  });

  it('should handle exact matches only', () => {
    const result = displayNameToNamespace('Default', mockNamespaces);
    expect(result).toBe('Default');
  });
});

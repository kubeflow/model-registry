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
  it('should return display name when namespace exists', () => {
    const namespaces: NamespaceKind[] = [
      { name: 'default', displayName: 'Default Namespace' },
      { name: 'kube-system', displayName: 'Kubernetes System' },
      { name: 'my-project', displayName: 'My Project' },
    ];

    expect(namespaceToDisplayName('default', namespaces)).toBe('Default Namespace');
    expect(namespaceToDisplayName('kube-system', namespaces)).toBe('Kubernetes System');
    expect(namespaceToDisplayName('my-project', namespaces)).toBe('My Project');
  });

  it('should return namespace name when namespace not found', () => {
    const namespaces: NamespaceKind[] = [{ name: 'default', displayName: 'Default Namespace' }];

    expect(namespaceToDisplayName('non-existent', namespaces)).toBe('non-existent');
  });

  it('should handle namespace with same name and displayName', () => {
    const namespaces: NamespaceKind[] = [{ name: 'test', displayName: 'test' }];

    expect(namespaceToDisplayName('test', namespaces)).toBe('test');
  });

  it('should handle namespace with empty displayName', () => {
    const namespaces: NamespaceKind[] = [{ name: 'test', displayName: '' }];

    expect(namespaceToDisplayName('test', namespaces)).toBe('test');
  });
});

describe('displayNameToNamespace', () => {
  it('should return namespace name when display name exists', () => {
    const namespaces: NamespaceKind[] = [
      { name: 'default', displayName: 'Default Namespace' },
      { name: 'kube-system', displayName: 'Kubernetes System' },
      { name: 'my-project', displayName: 'My Project' },
    ];

    expect(displayNameToNamespace('Default Namespace', namespaces)).toBe('default');
    expect(displayNameToNamespace('Kubernetes System', namespaces)).toBe('kube-system');
    expect(displayNameToNamespace('My Project', namespaces)).toBe('my-project');
  });

  it('should return display name when namespace not found', () => {
    const namespaces: NamespaceKind[] = [{ name: 'default', displayName: 'Default Namespace' }];

    expect(displayNameToNamespace('Non-existent Display Name', namespaces)).toBe(
      'Non-existent Display Name',
    );
  });

  it('should handle namespace with same name and displayName', () => {
    const namespaces: NamespaceKind[] = [{ name: 'test', displayName: 'test' }];

    expect(displayNameToNamespace('test', namespaces)).toBe('test');
  });
});

import {
  ModelRegistryKind,
  GroupKind,
  RoleBindingKind,
  K8sResourceCommon,
  RoleBindingSubject,
  RoleBindingRoleRef,
  genRandomChars,
} from 'mod-arch-shared';
import {
  APIOptions,
  handleRestFailures,
  UserSettings,
  assembleModArchBody,
  isModArchResponse,
  restCREATE,
  restDELETE,
  restGET,
  restPATCH,
} from 'mod-arch-core';
import { ModelRegistry, ModelRegistryPayload } from '~/app/types';
import { BFF_API_VERSION, URL_PREFIX } from '~/app/utilities/const';
import { RoleBindingPermissionsRoleType } from '~/app/pages/settings/roleBinding/types';
import { ListConfigSecretsResponse, NamespaceKind } from '~/app/shared/components/types';

export type ModelRegistryAndCredentials = {
  modelRegistry: ModelRegistryKind;
  databasePassword?: string;
  newDatabaseCACertificate?: string;
};

export const getListModelRegistries =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions): Promise<ModelRegistry[]> =>
    handleRestFailures(
      restGET(hostPath, `${URL_PREFIX}/api/${BFF_API_VERSION}/model_registry`, queryParams, opts),
    ).then((response) => {
      if (isModArchResponse<ModelRegistry[]>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getUser =
  (hostPath: string) =>
  (opts: APIOptions): Promise<UserSettings> =>
    handleRestFailures(
      restGET(hostPath, `${URL_PREFIX}/api/${BFF_API_VERSION}/user`, {}, opts),
    ).then((response) => {
      if (isModArchResponse<UserSettings>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getNamespaces =
  (hostPath: string) =>
  (opts: APIOptions): Promise<NamespaceKind[]> =>
    handleRestFailures(
      restGET(hostPath, `${URL_PREFIX}/api/${BFF_API_VERSION}/namespaces`, {}, opts),
    ).then((response) => {
      if (isModArchResponse<NamespaceKind[]>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getNamespacesForSettings =
  (hostPath: string) =>
  (opts: APIOptions): Promise<NamespaceKind[]> =>
    handleRestFailures(
      restGET(hostPath, `${URL_PREFIX}/api/${BFF_API_VERSION}/settings/namespaces`, {}, opts),
    ).then((response) => {
      if (isModArchResponse<NamespaceKind[]>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getGroups =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions): Promise<GroupKind[]> =>
    handleRestFailures(
      restGET(hostPath, `${URL_PREFIX}/api/${BFF_API_VERSION}/groups`, queryParams, opts),
    ).then((response) => {
      if (isModArchResponse<GroupKind[]>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getRoleBindings =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions): Promise<RoleBindingKind[]> =>
    handleRestFailures(
      restGET(
        hostPath,
        `${URL_PREFIX}/api/${BFF_API_VERSION}/settings/role_bindings`,
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModArchResponse<{ items: RoleBindingKind[] }>(response)) {
        return response.data.items;
      }
      throw new Error('Invalid response format');
    });

export const getModelRegistrySettings =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, modelRegistryId: string): Promise<ModelRegistryAndCredentials> =>
    handleRestFailures(
      restGET(
        hostPath,
        `${URL_PREFIX}/api/${BFF_API_VERSION}/settings/model_registry/${modelRegistryId}`,
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModArchResponse<ModelRegistryAndCredentials>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const listModelRegistrySettings =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions): Promise<ModelRegistryKind[]> =>
    handleRestFailures(
      restGET(
        hostPath,
        `${URL_PREFIX}/api/${BFF_API_VERSION}/settings/model_registry`,
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModArchResponse<ModelRegistryKind[]>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const createModelRegistrySettings =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, data: ModelRegistryPayload): Promise<ModelRegistryKind> =>
    handleRestFailures(
      restCREATE(
        hostPath,
        `${URL_PREFIX}/api/${BFF_API_VERSION}/settings/model_registry`,
        assembleModArchBody(data),
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModArchResponse<ModelRegistryKind>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const deleteModelRegistrySettings =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (
    opts: APIOptions,
    data: ModelRegistryKind,
    modelRegistryId: string,
  ): Promise<ModelRegistryKind[]> =>
    handleRestFailures(
      restDELETE(
        hostPath,
        `${URL_PREFIX}/api/${BFF_API_VERSION}/settings/model_registry/${modelRegistryId}`,
        assembleModArchBody(data),
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModArchResponse<ModelRegistryKind[]>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const patchModelRegistrySettings =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (
    opts: APIOptions,
    data: ModelRegistryKind,
    modelRegistryId: string,
  ): Promise<ModelRegistryKind[]> =>
    handleRestFailures(
      restPATCH(
        hostPath,
        `${URL_PREFIX}/api/${BFF_API_VERSION}/settings/model_registry/${modelRegistryId}`,
        assembleModArchBody(data),
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModArchResponse<ModelRegistryKind[]>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const createRoleBinding =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, data: RoleBindingKind): Promise<RoleBindingKind> =>
    handleRestFailures(
      restCREATE(
        hostPath,
        `${URL_PREFIX}/api/${BFF_API_VERSION}/settings/role_bindings`,
        assembleModArchBody(data),
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModArchResponse<RoleBindingKind>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const patchRoleBinding =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, data: RoleBindingKind, roleBindingName: string): Promise<RoleBindingKind> =>
    handleRestFailures(
      restPATCH(
        hostPath,
        `${URL_PREFIX}/api/${BFF_API_VERSION}/settings/role_bindings/${roleBindingName}`,
        assembleModArchBody(data),
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModArchResponse<RoleBindingKind>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const deleteRoleBinding =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, roleBindingName: string): Promise<void> =>
    handleRestFailures(
      restDELETE(
        hostPath,
        `${URL_PREFIX}/api/${BFF_API_VERSION}/settings/role_bindings/${roleBindingName}`,
        {},
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModArchResponse<void>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

//TODO : migrate this to shared library
export const addOwnerReference = <R extends K8sResourceCommon>(
  resource: R,
  owner?: K8sResourceCommon,
  blockOwnerDeletion = false,
): R => {
  if (!owner) {
    return resource;
  }
  const ownerReferences = resource.metadata?.ownerReferences || [];
  if (
    owner.metadata?.uid &&
    owner.metadata.name &&
    !ownerReferences.find((r) => r.uid === owner.metadata?.uid)
  ) {
    ownerReferences.push({
      uid: owner.metadata.uid,
      name: owner.metadata.name,
      apiVersion: owner.apiVersion,
      kind: owner.kind,
      blockOwnerDeletion,
    });
  }
  return {
    ...resource,
    metadata: {
      ...resource.metadata,
      ownerReferences,
    },
  };
};

export const generateRoleBindingPermissions = (
  namespace: string,
  rbSubjectKind: RoleBindingSubject['kind'],
  rbSubjectName: RoleBindingSubject['name'],
  rbRoleRefName: RoleBindingPermissionsRoleType | string, //string because with MR this can include MR name
  rbRoleRefKind: RoleBindingRoleRef['kind'],
  rbLabels?: { [key: string]: string },
  ownerReference?: K8sResourceCommon,
): RoleBindingKind => {
  const roleBindingObject: RoleBindingKind = {
    apiVersion: 'rbac.authorization.k8s.io/v1',
    kind: 'RoleBinding',
    metadata: {
      name: `dashboard-permissions-${genRandomChars()}`,
      namespace,
      labels: rbLabels,
    },
    roleRef: {
      apiGroup: 'rbac.authorization.k8s.io',
      kind: rbRoleRefKind,
      name: rbRoleRefName,
    },
    subjects: [
      {
        apiGroup: 'rbac.authorization.k8s.io',
        kind: rbSubjectKind,
        name: rbSubjectName,
      },
    ],
  };
  return addOwnerReference(roleBindingObject, ownerReference);
};

export const listModelRegistryCertificateNames =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions): Promise<ListConfigSecretsResponse> =>
    handleRestFailures(
      restGET(
        hostPath,
        `${URL_PREFIX}/api/${BFF_API_VERSION}/settings/certificates`,
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModArchResponse<{ items: ListConfigSecretsResponse }>(response)) {
        return response.data.items;
      }
      throw new Error('Invalid response format');
    });

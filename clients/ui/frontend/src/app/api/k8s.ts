import {
  APIOptions,
  handleRestFailures,
  Namespace,
  UserSettings,
  ModelRegistryKind,
  assembleModArchBody,
  isModArchResponse,
  restCREATE,
  restDELETE,
  restGET,
  restPATCH,
  GroupKind,
  RoleBindingKind,
} from 'mod-arch-shared';
import { ModelRegistry } from '~/app/types';
import { BFF_API_VERSION, URL_PREFIX } from '~/app/utilities/const';

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
  (opts: APIOptions): Promise<Namespace[]> =>
    handleRestFailures(
      restGET(hostPath, `${URL_PREFIX}/api/${BFF_API_VERSION}/namespaces`, {}, opts),
    ).then((response) => {
      if (isModArchResponse<Namespace[]>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getNamespacesForSettings =
  (hostPath: string) =>
  (opts: APIOptions): Promise<Namespace[]> =>
    handleRestFailures(
      restGET(hostPath, `${URL_PREFIX}/api/${BFF_API_VERSION}/settings/namespace`, {}, opts),
    ).then((response) => {
      if (isModArchResponse<Namespace[]>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getGroups =
  (hostPath: string) =>
  (opts: APIOptions): Promise<GroupKind[]> =>
    handleRestFailures(
      restGET(hostPath, `${URL_PREFIX}/api/${BFF_API_VERSION}/groups`, {}, opts),
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
      if (isModArchResponse<RoleBindingKind[]>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getModelRegistrySettings =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, modelRegistryId: string): Promise<ModelRegistryKind> =>
    handleRestFailures(
      restGET(
        hostPath,
        `${URL_PREFIX}/api/${BFF_API_VERSION}/settings/model_registry/${modelRegistryId}`,
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModArchResponse<ModelRegistryKind>(response)) {
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
  (opts: APIOptions, data: ModelRegistryKind): Promise<ModelRegistryKind[]> =>
    handleRestFailures(
      restCREATE(
        hostPath,
        `${URL_PREFIX}/api/${BFF_API_VERSION}/settings/model_registry`,
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

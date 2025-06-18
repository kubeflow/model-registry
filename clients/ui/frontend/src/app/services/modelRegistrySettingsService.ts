import * as _ from 'lodash-es';
import { K8sStatus } from 'mod-arch-shared';
import axios from '~/app/utils/axios';
import { ListConfigSecretsResponse, ModelRegistryKind, RoleBindingKind } from '~/app/k8sTypes';
import { RecursivePartial } from '~/typeHelpers';

const registriesUrl = '/api/modelRegistries';
const mrRoleBindingsUrl = '/api/modelRegistryRoleBindings';
const configSecretsUrl = '/api/modelRegistryCertificates';

export type ModelRegistryAndCredentials = {
  modelRegistry: ModelRegistryKind;
  databasePassword?: string;
  newDatabaseCACertificate?: string;
};

export const listModelRegistriesBackend = (labelSelector?: string): Promise<ModelRegistryKind[]> =>
  axios
    .get(registriesUrl, { params: { labelSelector } })
    .then((response) => response.data.items)
    .catch((e) => {
      throw new Error(e.response.data.message);
    });

export const createModelRegistryBackend = (
  data: ModelRegistryAndCredentials,
): Promise<ModelRegistryKind> =>
  axios
    .post(registriesUrl, data)
    .then((response) => response.data)
    .catch((e) => {
      throw new Error(e.response.data.message);
    });

export const getModelRegistryBackend = (
  modelRegistryName: string,
): Promise<ModelRegistryAndCredentials> =>
  axios
    .get(`${registriesUrl}/${modelRegistryName}`)
    .then((response) => response.data)
    .catch((e) => {
      throw new Error(e.response.data.message);
    });

export const updateModelRegistryBackend = (
  modelRegistryName: string,
  patch: RecursivePartial<ModelRegistryAndCredentials>,
): Promise<ModelRegistryAndCredentials> =>
  axios
    .patch(`${registriesUrl}/${modelRegistryName}`, patch)
    .then((response) => response.data)
    .catch((e) => {
      throw new Error(e.response.data.message);
    });

export const deleteModelRegistryBackend = (modelRegistryName: string): Promise<ModelRegistryKind> =>
  axios
    .delete(`${registriesUrl}/${modelRegistryName}`)
    .then((response) => response.data)
    .catch((e) => {
      throw new Error(e.response.data.message);
    });

export const listModelRegistryRoleBindings = (): Promise<RoleBindingKind[]> =>
  axios
    .get(mrRoleBindingsUrl)
    .then((response) => response.data.items)
    .catch((e) => {
      throw new Error(e.response.data.message);
    });

export const createModelRegistryRoleBinding = (data: RoleBindingKind): Promise<RoleBindingKind> => {
  const roleBindingWithoutNamespace: RoleBindingKind & { metadata: { namespace?: string } } =
    _.omit(data, 'metadata.namespace');
  return axios
    .post(mrRoleBindingsUrl, roleBindingWithoutNamespace)
    .then((response) => response.data)
    .catch((e) => {
      throw new Error(e.response.data.message);
    });
};

export const deleteModelRegistryRoleBinding = (roleBindingName: string): Promise<K8sStatus> =>
  axios
    .delete(`${mrRoleBindingsUrl}/${roleBindingName}`)
    .then((response) => response.data)
    .catch((e) => {
      throw new Error(e.response.data.message);
    });

export const listModelRegistryCertificateNames = (): Promise<ListConfigSecretsResponse> =>
  axios
    .get<ListConfigSecretsResponse>(configSecretsUrl)
    .then((response) => response.data)
    .catch((e) => {
      throw new Error(e.response.data.message);
    });

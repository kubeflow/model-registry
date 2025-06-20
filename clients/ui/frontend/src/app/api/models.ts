import { K8sModel } from '~/app/k8sTypes';

export const ModelRegistryModel: K8sModel = {
  apiGroup: 'modelregistry.opendatahub.io',
  apiVersion: 'v1alpha1',
  kind: 'ModelRegistry',
  plural: 'modelregistries',
  abbr: 'mr',
  label: 'ModelRegistry',
  labelPlural: 'ModelRegistries',
  crd: true,
};

import { K8sModel } from '~/app/k8sTypes';

export const ProjectModel: K8sModel = {
  apiGroup: 'project.openshift.io',
  apiVersion: 'v1',
  kind: 'Project',
  plural: 'projects',
  abbr: 'p',
  label: 'Project',
  labelPlural: 'Projects',
  crd: false,
};

export const RoleBindingModel: K8sModel = {
  apiGroup: 'rbac.authorization.k8s.io',
  apiVersion: 'v1',
  kind: 'RoleBinding',
  plural: 'rolebindings',
  abbr: 'rb',
  label: 'RoleBinding',
  labelPlural: 'RoleBindings',
  crd: false,
};

export const SelfSubjectAccessReviewModel: K8sModel = {
  apiGroup: 'authorization.k8s.io',
  apiVersion: 'v1',
  kind: 'SelfSubjectAccessReview',
  plural: 'selfsubjectaccessreviews',
  abbr: 'ssar',
  label: 'Self Subject Access Review',
  labelPlural: 'Self Subject Access Reviews',
  crd: false,
};

export const GroupModel: K8sModel = {
    apiGroup: 'user.openshift.io',
    apiVersion: 'v1',
    kind: 'Group',
    plural: 'groups',
    abbr: 'g',
    label: 'Group',
    labelPlural: 'Groups',
    crd: false,
};

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

export const DataScienceClusterModel: K8sModel = {
    apiGroup: 'datasciencecluster.opendatahub.io',
    apiVersion: 'v1',
    kind: 'DataScienceCluster',
    plural: 'datascienceclusters',
    abbr: 'dsc',
    label: 'DataScienceCluster',
    labelPlural: 'DataScienceClusters',
    crd: true,
}; 
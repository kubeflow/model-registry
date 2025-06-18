import { ModelRegistryKind } from '~/app/k8sTypes';

const useModelRegistryNamespaceCR = (
  namespace: string | undefined,
  name: string,
): [ModelRegistryKind | null, boolean] => {
  if (!name) {
    return [null, true];
  }
  const mockModelRegistry: ModelRegistryKind = {
    apiVersion: 'modelregistry.opendatahub.io/v1alpha1',
    kind: 'ModelRegistry',
    metadata: {
      name,
      namespace: namespace || 'opendatahub',
      uid: '1234',
      resourceVersion: '1',
    },
    spec: {
      grpc: {},
      rest: {},
      oauthProxy: {},
      mysql: {
        database: 'db',
        host: 'host',
      },
    },
  };
  return [mockModelRegistry, true];
};

export { useModelRegistryNamespaceCR };

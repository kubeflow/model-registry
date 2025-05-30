import { DatabaseType, ModelRegistryKind } from 'mod-arch-shared';

type MockModelRegistryKind = {
  name?: string;
  description?: string;
  displayName?: string;
};

export const mockModelRegistryKind = ({
  name = 'modelregistry-sample',
  description = 'Model registry description',
  displayName = 'Model Registry Sample',
}: MockModelRegistryKind): ModelRegistryKind => ({
  apiVersion: 'modelregistry.opendatahub.io/v1alpha1',
  kind: 'ModelRegistry',
  metadata: {
    name,
    namespace: 'default',
    displayName,
    description,
    creationTimestamp: '2024-01-01T00:00:00Z',
    generation: 1,
    resourceVersion: '12345',
    uid: 'abc-123-def-456',
    labels: {
      'app.kubernetes.io/name': 'model-registry',
      'app.kubernetes.io/component': 'model-registry',
    },
  },
  spec: {
    grpc: {},
    rest: {},
    istio: {
      gateway: {
        grpc: {
          tls: {},
        },
        rest: {
          tls: {},
        },
      },
    },
  },
  databaseConfig: {
    databaseType: DatabaseType.MySQL,
    host: 'localhost',
    port: 5432,
    database: 'model_registry',
    username: 'mlmd',
    skipDBCreation: false,
  },
  status: {
    conditions: [
      {
        type: 'Ready',
        status: 'True',
        lastTransitionTime: '2024-01-01T00:00:00Z',
        reason: 'Available',
        message: 'Model registry is ready',
      },
    ],
  },
  data: {},
});

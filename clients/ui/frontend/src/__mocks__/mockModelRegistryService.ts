export type ServiceKind = {
  kind: 'Service';
  apiVersion: string;
  metadata: {
    name: string;
    namespace: string;
    annotations: {
      'openshift.io/description': string;
      'openshift.io/display-name': string;
      'routing.opendatahub.io/external-address-rest': string;
    };
  };
  spec: {
    selector: {
      app: string;
      component: string;
    };
    ports: {
      name: string;
      protocol: string;
      port: number;
      targetPort: number;
    }[];
  };
};

type MockServiceType = {
  name?: string;
  namespace?: string;
  description?: string;
  serverUrl?: string;
};

export const mockModelRegistryService = ({
  name = 'modelregistry-sample',
  namespace = 'odh-model-registries',
  description = 'Model registry description',
  serverUrl = 'modelregistry-sample-rest.com:443',
}: MockServiceType): ServiceKind => ({
  kind: 'Service',
  apiVersion: 'v1',
  metadata: {
    name,
    namespace,
    annotations: {
      'openshift.io/description': description,
      'openshift.io/display-name': name,
      'routing.opendatahub.io/external-address-rest': serverUrl,
    },
  },
  spec: {
    selector: {
      app: name,
      component: 'model-registry',
    },
    ports: [
      {
        name: 'grpc-api',
        protocol: 'TCP',
        port: 9090,
        targetPort: 9090,
      },
      {
        name: 'http-api',
        protocol: 'TCP',
        port: 8080,
        targetPort: 8080,
      },
    ],
  },
});

import { EitherNotBoth } from './typeHelpers';
import {
  ContainerResources,
  K8sResourceCommon,
  NodeSelector,
  PodAffinity,
  TolerationEffect,
  TolerationOperator,
  Volume,
  VolumeMount,
} from './types';

export type K8sCondition = {
  type: string;
  status: string;
  reason?: string;
  message?: string;
  lastProbeTime?: string | null;
  lastTransitionTime?: string;
  lastHeartbeatTime?: string;
};

export type ModelRegistryKind = K8sResourceCommon & {
  metadata: {
    name: string;
    namespace: string;
  };
  spec: {
    grpc: Record<string, never>; // Empty object at create time, properties here aren't used by the UI
    rest: Record<string, never>; // Empty object at create time, properties here aren't used by the UI
    istio: {
      gateway: {
        grpc: {
          tls: Record<string, never>; // Empty object at create time, properties here aren't used by the UI
        };
        rest: {
          tls: Record<string, never>; // Empty object at create time, properties here aren't used by the UI
        };
      };
    };
  };
  databaseConfig: DatabaseConfig;
  status?: {
    conditions?: K8sCondition[];
  };
};

export enum DatabaseType {
  MySQL = 'MySQL',
  Postgres = 'Postgres',
}

export type PasswordSecret = {
  key: string;
  name: string;
};

export type DatabaseConfig = {
  databaseType: DatabaseType;
  database: string;
  host: string;
  passwordSecret?: PasswordSecret;
  port: number;
  skipDBCreation: boolean;
  username: string;
  sslRootCertificateConfigMap?: string;
  sslRootCertificateSecret?: string;
} & EitherNotBoth<
  {
    sslRootCertificateConfigMap?: {
      name: string;
      key: string;
    } | null;
  },
  {
    sslRootCertificateSecret?: {
      name: string;
      key: string;
    } | null;
  }
>;

export enum DeploymentMode {
  ModelMesh = 'ModelMesh',
  RawDeployment = 'RawDeployment',
  Serverless = 'Serverless',
}

export type InferenceServiceAnnotations = Partial<{
  'security.kubeflow.io/enable-auth': string;
}>;

export type InferenceServiceLabels = Partial<{
  'networking.knative.dev/visibility': string;
  'networking.kserve.io/visibility': 'exposed';
}>;

export type ImagePullSecret = {
  name: string;
};

export type InferenceServiceKind = K8sResourceCommon & {
  metadata: {
    name: string;
    namespace: string;
    annotations?: InferenceServiceAnnotations &
      Partial<{
        'serving.kserve.io/deploymentMode': DeploymentMode;
        'sidecar.istio.io/inject': 'true';
        'sidecar.istio.io/rewriteAppHTTPProbers': 'true';
      }>;
    labels?: InferenceServiceLabels;
  };
  spec: {
    predictor: {
      tolerations?: Toleration[];
      nodeSelector?: NodeSelector;
      model?: {
        modelFormat?: {
          name: string;
          version?: string;
        };
        resources?: ContainerResources;
        runtime?: string;
        storageUri?: string;
        storage?: {
          key?: string;
          parameters?: Record<string, string>;
          path?: string;
          schemaPath?: string;
        };
        args?: ServingContainer['args'];
        env?: ServingContainer['env'];
      };
      maxReplicas?: number;
      minReplicas?: number;
      imagePullSecrets?: ImagePullSecret[];
    };
  };
  status?: {
    components?: {
      predictor?: {
        grpcUrl?: string;
        restUrl?: string;
        url?: string;
      };
    };
    conditions?: {
      lastTransitionTime?: string;
      status: string;
      type: string;
    }[];
    modelStatus?: {
      copies?: {
        failedCopies?: number;
        totalCopies?: number;
      };
      lastFailureInfo?: {
        location?: string;
        message?: string;
        modelRevisionName?: string;
        reason?: string;
        time?: string;
      };
      states?: {
        activeModelState: string;
        targetModelState?: string;
      };
      transitionStatus: string;
    };
    url: string;
    address?: {
      CACerts?: string;
      audience?: string;
      name?: string;
      url?: string;
    };
  };
};

export type Toleration = {
  key: string;
  operator?: TolerationOperator;
  value?: string;
  effect?: TolerationEffect;
  tolerationSeconds?: number;
};

export type ServingContainer = {
  name: string;
  args?: string[];
  image?: string;
  affinity?: PodAffinity;
  resources?: ContainerResources;
  volumeMounts?: VolumeMount[];
  env?: {
    name: string;
    value?: string;
    valueFrom?: {
      secretKeyRef?: {
        name: string;
        key: string;
      };
    };
  }[];
};

export type ServingRuntimeAnnotations = Partial<{
  'enable-route': string;
  'enable-auth': string;
  'modelmesh-enabled': 'true' | 'false';
}>;

export type SupportedModelFormats = {
  name: string;
  version?: string;
  autoSelect?: boolean;
};

export type ServingRuntimeKind = K8sResourceCommon & {
  metadata: {
    annotations?: ServingRuntimeAnnotations;
    name: string;
    namespace: string;
  };
  spec: {
    builtInAdapter?: {
      serverType?: string;
      runtimeManagementPort?: number;
      memBufferBytes?: number;
      modelLoadingTimeoutMillis?: number;
    };
    containers: ServingContainer[];
    supportedModelFormats?: SupportedModelFormats[];
    replicas?: number;
    tolerations?: Toleration[];
    nodeSelector?: NodeSelector;
    volumes?: Volume[];
    imagePullSecrets?: ImagePullSecret[];
  };
};

export type ServiceKind = K8sResourceCommon & {
  metadata: {
    name: string;
    namespace: string;
    labels?: Partial<{
      'kubeflow.io/user': string;
      component: string;
    }>;
  };
  spec: {
    selector: {
      app: string;
      component: string;
    };
    ports: {
      name?: string;
      protocol?: string;
      appProtocol?: string;
      port?: number;
      targetPort?: number | string;
    }[];
  };
};

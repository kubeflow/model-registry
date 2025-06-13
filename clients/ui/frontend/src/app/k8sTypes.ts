import { K8sResourceCommon, K8sModelCommon, MatchExpression } from '@openshift/dynamic-plugin-sdk-utils';
import { EitherNotBoth } from '../typeHelpers';

export type K8sModel = {
  apiGroup: string;
  apiVersion: string;
  kind: string;
  plural: string;
  abbr: string;
  label: string;
  labelPlural: string;
  crd: boolean;
};

export enum KnownLabels {
  DASHBOARD_RESOURCE = 'opendatahub.io/dashboard',
  PROJECT_SHARING = 'opendatahub.io/project-sharing',
  MODEL_SERVING_PROJECT = 'modelmesh-enabled',
  DATA_CONNECTION_AWS = 'opendatahub.io/managed',
  LABEL_SELECTOR_MODEL_REGISTRY = 'component=model-registry',
  PROJECT_SUBJECT = 'opendatahub.io/rb-project-subject',
  REGISTERED_MODEL_ID = 'modelregistry.opendatahub.io/registered-model-id',
  MODEL_VERSION_ID = 'modelregistry.opendatahub.io/model-version-id',
  MODEL_REGISTRY_NAME = 'modelregistry.opendatahub.io/name',
}

export type K8sVerb =
  | 'create'
  | 'get'
  | 'list'
  | 'update'
  | 'patch'
  | 'delete'
  | 'deletecollection'
  | 'watch';

export type DisplayNameAnnotations = Partial<{
  'openshift.io/description': string;
  'openshift.io/display-name': string;
}>;

export type K8sDSGResource = K8sResourceCommon & {
  metadata: {
    annotations?: DisplayNameAnnotations &
      Partial<{
        'opendatahub.io/recommended-accelerators': string;
      }>;
    name: string;
  };
};

export type K8sCondition = {
  type: string;
  status: string;
  reason?: string;
  message?: string;
  lastProbeTime?: string | null;
  lastTransitionTime?: string;
  lastHeartbeatTime?: string;
};

export type DashboardLabels = {
    [KnownLabels.DASHBOARD_RESOURCE]: 'true';
};

export type ModelServingProjectLabels = {
    [KnownLabels.MODEL_SERVING_PROJECT]: 'true' | 'false';
};

export type ProjectKind = K8sResourceCommon & {
  metadata: {
    annotations?: DisplayNameAnnotations &
      Partial<{
        'openshift.io/requester': string;
      }>;
    labels?: Partial<DashboardLabels> & Partial<ModelServingProjectLabels>;
    name: string;
  };
  status?: {
    phase: 'Active' | 'Terminating';
  };
};

export type RoleBindingSubject = {
  kind: string;
  apiGroup?: string;
  name: string;
};

export type RoleBindingRoleRef = {
  kind: 'Role' | 'ClusterRole';
  apiGroup?: string;
  name: string;
};

export type RoleBindingKind = K8sResourceCommon & {
  metadata: {
    name: string;
    namespace: string;
  };
  subjects: RoleBindingSubject[];
  roleRef: RoleBindingRoleRef;
};

export type GroupKind = K8sResourceCommon & {
  metadata: {
    name: string;
  };
  users: string[];
};

export type ModelRegistryKind = K8sResourceCommon & {
  metadata: {
    name: string;
    namespace: string;
    annotations?: DisplayNameAnnotations;
  };
  spec: {
    grpc: Record<string, never>; 
    rest: Record<string, never>;
    oauthProxy: Record<string, never>;
  } & EitherNotBoth<
    {
      mysql?: {
        database: string;
        host: string;
        passwordSecret?: {
          key: string;
          name: string;
        };
        port?: number;
        skipDBCreation?: boolean;
        username?: string;
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
    },
    {
      postgres?: {
        database: string;
        host?: string;
        passwordSecret?: {
          key: string;
          name: string;
        };
        port: number;
        skipDBCreation?: boolean;
        sslMode?: string;
        username?: string;
      };
    }
  >;
  status?: {
    conditions?: K8sCondition[];
  };
};

export enum DeploymentMode {
  ModelMesh = 'ModelMesh',
  RawDeployment = 'RawDeployment',
  Serverless = 'Serverless',
}

export type DataScienceClusterComponentStatus = {
  managementState?: 'Managed' | 'Removed';
};

export type DataScienceClusterKindStatus = {
  components?: {
    modelregistry?: {
      registriesNamespace?: string;
    };
  };
  conditions: K8sCondition[];
  installedComponents?: { [key: string]: boolean };
  phase?: string;
  release?: {
    name: string;
    version: string;
  };
};

export type DataScienceClusterInitializationKindStatus = {
  conditions: K8sCondition[];
  release?: {
    name?: string;
    version?: string;
  };
  components?: Record<string, never>;
  phase?: string;
};

export type AccessReviewResourceAttributes = {
  group?: '*' | string;
  resource?: string;
  subresource?: '' | 'spec' | 'status';
  verb: '*' | K8sVerb;
  name?: string;
  namespace?: string;
};

export type SelfSubjectAccessReviewKind = K8sResourceCommon & {
  spec: {
    resourceAttributes?: AccessReviewResourceAttributes;
  };
  status?: {
    allowed: boolean;
    denied?: boolean;
    reason?: string;
    evaluationError?: string;
  };
};

export type DashboardCommonConfig = {
  enablement: boolean;
  disableInfo: boolean;
  disableSupport: boolean;
  disableClusterManager: boolean;
  disableTracking: boolean;
  disableBYONImageStream: boolean;
  disableISVBadges: boolean;
  disableAppLauncher: boolean;
  disableUserManagement: boolean;
  disableHome: boolean;
  disableProjects: boolean;
  disableModelServing: boolean;
  disableProjectScoped: boolean;
  disableProjectSharing: boolean;
  disableCustomServingRuntimes: boolean;
  disablePipelines: boolean;
  disableTrustyBiasMetrics: boolean;
  disablePerformanceMetrics: boolean;
  disableKServe: boolean;
  disableKServeAuth: boolean;
  disableKServeMetrics: boolean;
  disableKServeRaw: boolean;
  disableModelMesh: boolean;
  disableAcceleratorProfiles: boolean;
  disableHardwareProfiles: boolean;
  disableDistributedWorkloads: boolean;
  disableModelCatalog: boolean;
  disableModelRegistry: boolean;
  disableModelRegistrySecureDB: boolean;
  disableServingRuntimeParams: boolean;
  disableStorageClasses: boolean;
  disableNIMModelServing: boolean;
  disableAdminConnectionTypes: boolean;
  disableFineTuning: boolean;
  disableLlamaStackChatBot: boolean;
  disableLMEval: boolean;
};

export type DashboardConfigKind = K8sResourceCommon & {
  spec: {
    dashboardConfig: DashboardCommonConfig;
  };
};

export type DataScienceClusterKind = K8sResourceCommon & {
  metadata: {
    name: string;
  };
  spec: {
    components?: {
      [key in string]?: DataScienceClusterComponentStatus;
    } & {
      modelregistry?: {
        registriesNamespace: string;
      };
    };
  };
  status?: DataScienceClusterKindStatus;
};

export type ServiceKind = K8sResourceCommon & {
  metadata: {
    name: string;
    namespace: string;
    labels?: { [key: string]: string };
  };
  spec: {
    selector: { [key: string]: string };
    ports: {
      name?: string;
      port: number;
      protocol: string;
      targetPort: number;
    }[];
  };
};

export type ListConfigSecretsResponse = {
  secrets: { name: string; keys: string[] }[];
  configMaps: { name: string; keys: string[] }[];
};

export type ConfigSecretItem = {
    name: string;
    keys: string[];
};

export type OdhDocumentType = {
  [key: string]: any;
}

export type CustomWatchK8sResult<T> = [data: T, loaded: boolean, error: Error | undefined];

export type K8sAPIOptions = {
  dryRun?: boolean;
  signal?: AbortSignal;
  parseJSON?: boolean;
}; 
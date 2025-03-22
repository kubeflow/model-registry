import { EitherNotBoth } from './typeHelpers';

export type UserSettings = {
  userId: string;
  clusterAdmin?: boolean;
};

export type ConfigSettings = {
  common: CommonConfig;
};

export type CommonConfig = {
  featureFlags: FeatureFlag;
};

export type FeatureFlag = {
  modelRegistry: boolean;
};

export type KeyValuePair = {
  key: string;
  value: string;
};

export type Namespace = {
  name: string;
};

export declare type K8sResourceIdentifier = {
  apiGroup?: string;
  apiVersion: string;
  kind: string;
};

export declare type K8sResourceCommon = K8sResourceIdentifier &
  Partial<{
    metadata: Partial<{
      annotations: Record<string, string>;
      clusterName: string;
      creationTimestamp: string;
      deletionGracePeriodSeconds: number;
      deletionTimestamp: string;
      finalizers: string[];
      generateName: string;
      generation: number;
      labels: Record<string, string>;
      managedFields: unknown[];
      name: string;
      namespace: string;
      ownerReferences: OwnerReference[];
      resourceVersion: string;
      uid: string;
    }>;
    spec: {
      selector?: Selector | MatchLabels;
      [key: string]: unknown;
    };
    status: {
      [key: string]: unknown;
    };
    data: {
      [key: string]: unknown;
    };
  }>;

export declare type OwnerReference = {
  apiVersion: string;
  kind: string;
  name: string;
  uid: string;
  controller?: boolean;
  blockOwnerDeletion?: boolean;
};

export declare type Selector = Partial<{
  matchLabels: MatchLabels;
  matchExpressions: MatchExpression[];
  [key: string]: unknown;
}>;

export declare type MatchExpression = {
  key: string;
  operator: Operator | string;
  values?: string[];
  value?: string;
};

export declare type MatchLabels = {
  [key: string]: string;
};

export declare enum Operator {
  Exists = 'Exists',
  DoesNotExist = 'DoesNotExist',
  In = 'In',
  NotIn = 'NotIn',
  Equals = 'Equals',
  NotEqual = 'NotEqual',
  GreaterThan = 'GreaterThan',
  LessThan = 'LessThan',
  NotEquals = 'NotEquals',
}

export type UpdateObjectAtPropAndValue<T> = <K extends keyof T>(
  propKey: K,
  propValue: T[K],
) => void;

export type FetchStateObject<T, E = Error> = {
  data: T;
  loaded: boolean;
  error?: E;
  refresh: () => void;
};

export type NodeSelector = Record<string, string>;

export enum ContainerResourceAttributes {
  CPU = 'cpu',
  MEMORY = 'memory',
}

export type ContainerResources = {
  requests?: {
    [key: string]: number | string | undefined;
    cpu?: string | number;
    memory?: string;
  };
  limits?: {
    [key: string]: number | string | undefined;
    cpu?: string | number;
    memory?: string;
  };
};

export type EnvironmentVariable = EitherNotBoth<
  { value: string | number },
  { valueFrom: Record<string, unknown> }
> & {
  name: string;
};

export type PodAffinity = {
  nodeAffinity?: { [key: string]: unknown };
};

export type VolumeMount = { mountPath: string; name: string };

export enum TolerationOperator {
  EXISTS = 'Exists',
  EQUAL = 'Equal',
}

export enum TolerationEffect {
  NO_SCHEDULE = 'NoSchedule',
  PREFER_NO_SCHEDULE = 'PreferNoSchedule',
  NO_EXECUTE = 'NoExecute',
}

export type Volume = {
  name: string;
  emptyDir?: Record<string, unknown>;
  persistentVolumeClaim?: {
    claimName: string;
  };
  secret?: {
    secretName: string;
    optional?: boolean;
  };
};

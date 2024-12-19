import { ValueOf } from '~/shared/typeHelpers';

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

export type UpdateObjectAtPropAndValue<T> = (propKey: keyof T, propValue: ValueOf<T>) => void;

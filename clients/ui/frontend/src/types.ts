//  TODO: Get the status config params
export type UserSettings = {
  username: string;
  isAdmin: boolean;
  isAllowed: boolean;
};

// TODO: Add more config parameters
export type ConfigSettings = {
  common: CommonConfig;
};

// TODO: Add more config parameters
export type CommonConfig = {
  featureFlags: FeatureFlag;
};

// TODO: Add more config parameters
export type FeatureFlag = {
  modelRegistry: boolean;
};

export type KeyValuePair = {
  key: string;
  value: string;
};

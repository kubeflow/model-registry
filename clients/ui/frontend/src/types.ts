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

export type APIOptions = {
  dryRun?: boolean;
  signal?: AbortSignal;
  parseJSON?: boolean;
};

export type APIError = {
  code: string;
  message: string;
};

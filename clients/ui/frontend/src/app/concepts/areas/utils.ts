import { IsAreaAvailableStatus } from './types';

export const isAreaAvailable = (): IsAreaAvailableStatus => ({
  status: true,
  devFlags: null,
  featureFlags: null,
  reliantAreas: null,
  requiredComponents: null,
  requiredCapabilities: null,
  customCondition: () => true,
});

export type FlagState = { [key: string]: boolean };

export const getFlags = (): FlagState => ({});

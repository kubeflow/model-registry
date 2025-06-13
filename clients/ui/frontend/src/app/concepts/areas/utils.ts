import {
  DashboardConfigKind,
  DataScienceClusterInitializationKindStatus,
  DataScienceClusterKindStatus,
} from '~/app/k8sTypes';
import {
  CustomConditionFunction,
  IsAreaAvailableStatus,
  SupportedAreasState,
  SupportedAreaType,
} from './types';

export const isAreaAvailable = (
  area: SupportedAreaType,
  dashboardConfigSpec: DashboardConfigKind['spec'],
  dscStatus: DataScienceClusterKindStatus | null,
  dsciStatus: DataScienceClusterInitializationKindStatus | null,
  {
    internalStateMap,
    flagState,
  }: {
    internalStateMap: SupportedAreasState;
    flagState: FlagState;
  },
): IsAreaAvailableStatus => {
    return {
        status: true,
        devFlags: null,
        featureFlags: null,
        reliantAreas: null,
        requiredComponents: null,
        requiredCapabilities: null,
        customCondition: () => true,
    }
};

export type FlagState = { [key: string]: boolean };

export const getFlags = (spec: DashboardConfigKind['spec']): FlagState => {
    return {};
}; 
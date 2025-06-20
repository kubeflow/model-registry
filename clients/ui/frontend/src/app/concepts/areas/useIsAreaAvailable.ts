// import { AreaContext } from '~/app/concepts/areas/AreaContext';
import { IsAreaAvailableStatus } from '~/app/concepts/areas/types';

const useIsAreaAvailable = (): IsAreaAvailableStatus =>
  // (React.useContext(AreaContext).dscStatus.components[area] as IsAreaAvailableStatus) ?? {
  ({
    status: true,
    devFlags: null,
    featureFlags: null,
    reliantAreas: null,
    requiredComponents: null,
    requiredCapabilities: null,
    customCondition: () => true,
  });

export default useIsAreaAvailable;

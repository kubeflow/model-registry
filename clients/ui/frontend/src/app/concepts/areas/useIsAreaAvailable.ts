import * as React from 'react';
import { AreaContext } from '~/app/concepts/areas/AreaContext';
import { IsAreaAvailableStatus, SupportedArea } from '~/app/concepts/areas/types';

const useIsAreaAvailable = (area: SupportedArea): IsAreaAvailableStatus =>
  (React.useContext(AreaContext).dscStatus.components[area] as IsAreaAvailableStatus) ?? {
    status: false,
    devFlags: null,
    featureFlags: null,
    reliantAreas: null,
    requiredComponents: null,
    requiredCapabilities: null,
    customCondition: () => false,
  };

export default useIsAreaAvailable;

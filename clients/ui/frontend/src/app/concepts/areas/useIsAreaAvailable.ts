import * as React from 'react';
import { AreaContext } from '~/app/concepts/areas/AreaContext';
import { IsAreaAvailableStatus, SupportedArea } from './types';

const useIsAreaAvailable = (area: SupportedArea): IsAreaAvailableStatus =>
  React.useContext(AreaContext).dscStatus.components[area] ?? {
    status: false,
    devFlags: null,
    featureFlags: null,
    reliantAreas: null,
    requiredComponents: null,
    requiredCapabilities: null,
    customCondition: () => false,
  };

export default useIsAreaAvailable;

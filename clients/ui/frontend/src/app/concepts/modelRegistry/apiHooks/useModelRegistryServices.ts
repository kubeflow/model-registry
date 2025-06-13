import { ServiceKind } from '~/app/k8sTypes';
import useFetch, { FetchState } from '~/app/utils/useFetch';
import { POLL_INTERVAL } from '~/app/utils/const';
import * as React from 'react';

export const useModelRegistryServices = (namespace?: string): FetchState<ServiceKind[]> => {
    const getServices = React.useCallback(
        () => Promise.resolve([]), // This is a mock, replace with actual implementation
    [],
  );

  return useFetch<ServiceKind[]>(getServices, [], {
    refreshRate: POLL_INTERVAL,
  });
}; 
import * as React from 'react';
import { DataScienceClusterInitializationKindStatus } from '~/app/k8sTypes';
import useFetch, { FetchState } from '~/app/utils/useFetch';
import { POLL_INTERVAL } from '~/app/utils/const';

const useFetchDsciStatus = (): FetchState<DataScienceClusterInitializationKindStatus | null> => {
  const getDsci = React.useCallback(
    () => Promise.resolve(null), // This is a mock, replace with actual implementation
    [],
  );

  return useFetch<DataScienceClusterInitializationKindStatus | null>(getDsci, null, {
    refreshRate: POLL_INTERVAL,
  });
};

export default useFetchDsciStatus;

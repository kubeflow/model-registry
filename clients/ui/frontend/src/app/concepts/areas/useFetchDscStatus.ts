import * as React from 'react';
import { DataScienceClusterKindStatus } from '~/app/k8sTypes';
import { listDataScienceClusters } from '~/app/api/k8s/dsc';
import useFetch, { FetchState } from '~/app/utils/useFetch';
import { POLL_INTERVAL } from '~/app/utils/const';

const useFetchDscStatus = (): FetchState<DataScienceClusterKindStatus | null> => {
  const getDsc = React.useCallback(
    () =>
      listDataScienceClusters().then((dataScienceClusters) => {
        if (dataScienceClusters.length === 0) {
          return null;
        }
        return dataScienceClusters[0].status || null;
      }),
    [],
  );

  return useFetch<DataScienceClusterKindStatus | null>(getDsc, null, {
    refreshRate: POLL_INTERVAL,
  });
};

export default useFetchDscStatus;

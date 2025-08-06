import * as React from 'react';
import {
  useDeepCompareMemoize,
  FetchState,
  FetchStateCallbackPromise,
  useFetchState,
} from 'mod-arch-core';
import { ModelRegistry } from '~/app/types';
import { getListModelRegistries } from '~/app/api/k8s';

const useModelRegistries = (queryParams: Record<string, unknown>): FetchState<ModelRegistry[]> => {
  const paramsMemo = useDeepCompareMemoize(queryParams);

  const listModelRegistries = React.useMemo(
    () => getListModelRegistries('', paramsMemo),
    [paramsMemo],
  );
  const callback = React.useCallback<FetchStateCallbackPromise<ModelRegistry[]>>(
    (opts) => listModelRegistries(opts),
    [listModelRegistries],
  );
  return useFetchState(callback, [], { initialPromisePurity: true });
};

export default useModelRegistries;

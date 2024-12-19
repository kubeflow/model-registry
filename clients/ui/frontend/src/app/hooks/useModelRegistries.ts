import * as React from 'react';
import useFetchState, {
  FetchState,
  FetchStateCallbackPromise,
} from '~/shared/utilities/useFetchState';
import { ModelRegistry } from '~/app/types';
import { getListModelRegistries } from '~/shared/api/k8s';
import { useDeepCompareMemoize } from '~/shared/utilities/useDeepCompareMemoize';

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

import * as React from 'react';
import {
  useDeepCompareMemoize,
  FetchState,
  FetchStateCallbackPromise,
  useFetchState,
} from 'mod-arch-core';
import { ModelRegistryKind } from 'mod-arch-shared';
import { listModelRegistrySettings } from '~/app/api/k8s';

const useModelRegistriesSettings = (
  queryParams: Record<string, unknown>,
): FetchState<ModelRegistryKind[]> => {
  const paramsMemo = useDeepCompareMemoize(queryParams);

  const listModelRegistriesSettings = React.useMemo(
    () => listModelRegistrySettings('', paramsMemo),
    [paramsMemo],
  );
  const callback = React.useCallback<FetchStateCallbackPromise<ModelRegistryKind[]>>(
    (opts) => listModelRegistriesSettings(opts),
    [listModelRegistriesSettings],
  );
  return useFetchState(callback, [], { initialPromisePurity: true });
};

export default useModelRegistriesSettings;

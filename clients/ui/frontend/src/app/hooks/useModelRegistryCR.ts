import React from 'react';
import {
  APIOptions,
  useFetchState,
  useDeepCompareMemoize,
  FetchState,
  FetchStateCallbackPromise,
} from 'mod-arch-core';
import { ModelRegistryKind } from 'mod-arch-shared';
import { getModelRegistrySettings } from '~/app/api/k8s';

const useModelRegistryCR = (
  name: string,
  queryParams: Record<string, unknown>,
): FetchState<ModelRegistryKind | null> => {
  const paramsMemo = useDeepCompareMemoize(queryParams);
  const getModelRegistry = React.useMemo(
    () => getModelRegistrySettings('', paramsMemo),
    [paramsMemo],
  );

  const callback = React.useCallback<FetchStateCallbackPromise<ModelRegistryKind | null>>(
    (opts: APIOptions) =>
      name ? getModelRegistry(opts, name).then((res) => res.modelRegistry) : Promise.resolve(null),
    [getModelRegistry, name],
  );

  return useFetchState(callback, null, { initialPromisePurity: true });
};

export { useModelRegistryCR };

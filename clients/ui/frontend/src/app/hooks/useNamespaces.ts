import { useFetchState, APIOptions, FetchStateCallbackPromise } from 'mod-arch-core';
import React from 'react';
import { getNamespaces } from '~/app/api/k8s';
import { NamespaceKind } from '~/app/shared/components/types';

export const useNamespaces = (): [NamespaceKind[], boolean, Error | undefined] => {
  const callback = React.useCallback<FetchStateCallbackPromise<NamespaceKind[]>>(
    (opts: APIOptions) => getNamespaces('')(opts),
    [],
  );
  const [namespaces, loaded, error] = useFetchState<NamespaceKind[]>(callback, []);

  return [namespaces, loaded, error];
};

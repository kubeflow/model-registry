import { GroupKind } from 'mod-arch-shared';
import { useFetchState, APIOptions, FetchStateCallbackPromise } from 'mod-arch-core';
import React from 'react';
import { getGroups } from '~/app/api/k8s';

export const useGroups = (
  queryParams: Record<string, unknown> = {},
): [GroupKind[], boolean, Error | undefined] => {
  const callback = React.useCallback<FetchStateCallbackPromise<GroupKind[]>>(
    (opts: APIOptions) => getGroups('', queryParams)(opts),
    [queryParams],
  );
  const [groups, loaded, error] = useFetchState<GroupKind[]>(callback, []);

  return [groups, loaded, error];
};

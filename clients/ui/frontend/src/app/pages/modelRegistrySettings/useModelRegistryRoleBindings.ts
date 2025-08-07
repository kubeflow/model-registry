import * as React from 'react';
import {
  POLL_INTERVAL,
  useFetchState,
  FetchStateObject,
  useDeepCompareMemoize,
} from 'mod-arch-core';
import { RoleBindingKind } from 'mod-arch-shared';
import { getRoleBindings } from '~/app/api/k8s';

const useModelRegistryRoleBindings = (
  queryParams: Record<string, unknown>,
): FetchStateObject<RoleBindingKind[]> => {
  const paramsMemo = useDeepCompareMemoize(queryParams);
  const fetchRoleBindings = React.useCallback(
    () =>
      getRoleBindings(
        '',
        paramsMemo,
      )({}).catch((e) => {
        if (e.response?.status === 404) {
          throw new Error('No rolebindings found.');
        }
        throw e;
      }),
    [paramsMemo],
  );

  const [data, loaded, error, refresh] = useFetchState<RoleBindingKind[]>(fetchRoleBindings, [], {
    refreshRate: POLL_INTERVAL,
  });
  return { data, loaded, error, refresh };
};

export default useModelRegistryRoleBindings;

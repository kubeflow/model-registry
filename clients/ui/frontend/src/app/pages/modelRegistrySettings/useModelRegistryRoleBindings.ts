import * as React from 'react';
import { RoleBindingKind, POLL_INTERVAL, useFetchState, FetchStateObject } from 'mod-arch-shared';
import { getRoleBindings } from '~/app/api/k8s';

const useModelRegistryRoleBindings = (): FetchStateObject<RoleBindingKind[]> => {
  const fetchRoleBindings = React.useCallback(
    () =>
      getRoleBindings('')({}).catch((e) => {
        if (e.response?.status === 404) {
          throw new Error('No rolebindings found.');
        }
        throw e;
      }),
    [],
  );

  const [data, loaded, error, refresh] = useFetchState<RoleBindingKind[]>(fetchRoleBindings, [], {
    refreshRate: POLL_INTERVAL,
  });
  return { data, loaded, error, refresh };
};

export default useModelRegistryRoleBindings;

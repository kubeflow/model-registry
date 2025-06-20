import * as React from 'react';
import { RoleBindingKind } from '~/app/k8sTypes';
import { listModelRegistryRoleBindings } from '~/app/services/modelRegistrySettingsService';
import { POLL_INTERVAL } from 'mod-arch-shared';
import { useFetchState, FetchStateObject } from 'mod-arch-shared';

const useModelRegistryRoleBindings = (): FetchStateObject<RoleBindingKind[]> => {
  const getRoleBindings = React.useCallback(
    () =>
      listModelRegistryRoleBindings().catch((e) => {
        if (e.response?.status === 404) {
          throw new Error('No rolebindings found.');
        }
        throw e;
      }),
    [],
  );

  const [data, loaded, error, refresh] = useFetchState<RoleBindingKind[]>(getRoleBindings, [], {
    refreshRate: POLL_INTERVAL,
  });
  return { data, loaded, error, refresh };
};

export default useModelRegistryRoleBindings;

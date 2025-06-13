import * as React from 'react';
import { RoleBindingKind } from '~/app/k8sTypes';
import { listModelRegistryRoleBindings } from '~/app/services/modelRegistrySettingsService';
import { POLL_INTERVAL } from '~/app/utils/const';
import useFetch, { FetchState } from '~/app/utils/useFetch';

const useModelRegistryRoleBindings = (): FetchState<RoleBindingKind[]> => {
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

  return useFetch<RoleBindingKind[]>(getRoleBindings, [], { refreshRate: POLL_INTERVAL });
};

export default useModelRegistryRoleBindings;

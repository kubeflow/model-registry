import * as React from 'react';
import { useFetchState } from 'mod-arch-shared';
import { GroupKind } from '~/app/k8sTypes';

const getGroups = (): Promise<GroupKind[]> =>
  fetch('/api/v1/settings/groups')
    .then((res) => {
      if (res.ok) {
        return res.json();
      }
      throw new Error(res.statusText);
    })
    .then((data) => data.items);

export const useGroups = (): [GroupKind[], boolean, Error | undefined] => {
  const [groups, loaded, error] = useFetchState<GroupKind[]>(getGroups, []);

  return [groups, loaded, error];
};

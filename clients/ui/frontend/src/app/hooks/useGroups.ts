import { GroupKind, useFetchState } from 'mod-arch-shared';

const getGroupsForHook = (): Promise<GroupKind[]> =>
  fetch('${URL_PREFIX}/api/${BFF_API_VERSION}/')
    .then((res) => {
      if (res.ok) {
        return res.json();
      }
      throw new Error(res.statusText);
    })
    .then((data) => data.items);

export const useGroups = (): [GroupKind[], boolean, Error | undefined] => {
  const [groups, loaded, error] = useFetchState<GroupKind[]>(getGroupsForHook, []);

  return [groups, loaded, error];
};

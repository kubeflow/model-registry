import { useFetchState, APIOptions, FetchStateCallbackPromise, Namespace } from 'mod-arch-shared';
import React from 'react';
import { getNamespaces } from '~/app/api/k8s';
import { ProjectKind } from '~/app/shared/components/types';

export const useProjects = (): [ProjectKind[], boolean, Error | undefined] => {
  const callback = React.useCallback<FetchStateCallbackPromise<Namespace[]>>(
    (opts: APIOptions) => getNamespaces('')(opts),
    [],
  );
  const [namespaces, loaded, error] = useFetchState<Namespace[]>(callback, []);

  // Convert namespaces to projects format
  const projects = React.useMemo(
    () =>
      namespaces.map(
        (namespace): ProjectKind => ({
          metadata: {
            name: namespace.name || '',
            namespace: namespace.name || '',
          },
          kind: 'Project',
          apiVersion: 'project.openshift.io/v1',
          spec: {},
          status: { phase: 'Active' },
        }),
      ),
    [namespaces],
  );

  return [projects, loaded, error];
};

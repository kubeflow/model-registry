import * as React from 'react';
import useFetchState, {
  FetchState,
  FetchStateCallbackPromise,
} from '~/shared/utilities/useFetchState';
import { Namespace } from '~/shared/types';
import { AUTH_HEADER, isStandalone, MOCK_AUTH, KUBEFLOW_USERNAME } from '~/shared/utilities/const';
import { getNamespaces } from '~/app/api/k8s';

const useNamespaces = (): FetchState<Namespace[]> => {
  const listNamespaces = React.useMemo(() => getNamespaces(''), []);
  const callback = React.useCallback<FetchStateCallbackPromise<Namespace[]>>(
    (opts) => {
      if (!isStandalone()) {
        return Promise.resolve([]);
      }
      const headers = MOCK_AUTH ? { [AUTH_HEADER]: KUBEFLOW_USERNAME } : undefined;
      return listNamespaces({
        ...opts,
        headers,
      });
    },
    [listNamespaces],
  );
  return useFetchState(callback, [], { initialPromisePurity: true });
};

export default useNamespaces;

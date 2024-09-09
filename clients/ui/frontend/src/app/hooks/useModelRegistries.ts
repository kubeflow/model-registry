import * as React from 'react';
import { BFF_API_VERSION } from '~/app/const';
import useFetchState, { FetchState, FetchStateCallbackPromise } from '~/utilities/useFetchState';
import { ModelRegistryList } from '~/app/types';
import { getListModelRegistries } from '~/app/api/k8s';

const useModelRegistries = (): FetchState<ModelRegistryList> => {
  const listModelRegistries = getListModelRegistries(`/api/${BFF_API_VERSION}/model_registry`);
  const callback = React.useCallback<FetchStateCallbackPromise<ModelRegistryList>>(
    (opts) => listModelRegistries(opts),
    [listModelRegistries],
  );
  return useFetchState(callback, [], { initialPromisePurity: true });
};

export default useModelRegistries;

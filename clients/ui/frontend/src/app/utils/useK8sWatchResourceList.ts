import {
  K8sModelCommon,
  K8sResourceCommon,
  WatchK8sResource,
  WebSocketOptions,
  useK8sWatchResource,
} from '@openshift/dynamic-plugin-sdk-utils';
import React from 'react';
import { CustomWatchK8sResult } from '~/app/k8sTypes';

const useK8sWatchResourceList = <T extends K8sResourceCommon[]>(
  initResource: WatchK8sResource | null,
  initModel?: K8sModelCommon,
  options?: Partial<WebSocketOptions & RequestInit & { wsPrefix?: string; pathPrefix?: string }>,
): CustomWatchK8sResult<T> => {
  const initListResource = React.useMemo(
    () => (initResource != null ? { ...initResource, isList: true } : null),
    [initResource],
  );

  const [data, loaded, error] = useK8sWatchResource<T>(initListResource, initModel, options);

  const loadError = React.useMemo(() => {
    if (error instanceof Error) {
      return error;
    }

    if (!error) {
      return undefined;
    }

    return new Error('Unknown error occured');
  }, [error]);

  // disable as data can be `undefined` by the type in the SDK is incorrect
  // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
  return [React.useMemo(() => data ?? [], [data]), loaded, loadError];
};

export default useK8sWatchResourceList; 
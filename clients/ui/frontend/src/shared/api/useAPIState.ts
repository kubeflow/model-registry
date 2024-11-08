import * as React from 'react';
import { APIState } from '~/shared/api/types';

const useAPIState = <T>(
  hostPath: string | null,
  createAPI: (path: string) => T,
): [apiState: APIState<T>, refreshAPIState: () => void] => {
  const [internalAPIToggleState, setInternalAPIToggleState] = React.useState(false);

  const refreshAPIState = React.useCallback(() => {
    setInternalAPIToggleState((v) => !v);
  }, []);

  const apiState = React.useMemo<APIState<T>>(() => {
    let path = hostPath;
    if (!path) {
      // TODO: we need to figure out maybe a stopgap or something
      path = '';
    }
    const api = createAPI(path);

    return {
      apiAvailable: !!path,
      api,
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [createAPI, hostPath, internalAPIToggleState]);

  return [apiState, refreshAPIState];
};

export default useAPIState;

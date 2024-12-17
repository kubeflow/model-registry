import * as React from 'react';
import { useBrowserStorage } from '~/shared/components/browserStorage';
import { logout } from '~/shared/utilities/appUtils';

export type SetTime = (refreshDateMarker: Date) => void;

const useTimeBasedRefresh = (): SetTime => {
  const KEY_NAME = 'kf.dashboard.last.auto.refresh';
  const [lastRefreshTimestamp, setLastRefreshTimestamp] = useBrowserStorage(
    KEY_NAME,
    '0',
    false,
    true,
  );
  const ref = React.useRef<{
    lastRefreshTimestamp: string;
    setLastRefreshTimestamp: (newValue: string) => void;
  }>({ lastRefreshTimestamp, setLastRefreshTimestamp });
  ref.current = { lastRefreshTimestamp, setLastRefreshTimestamp };

  return React.useCallback<SetTime>((refreshDateMarker) => {
    // Intentionally avoid referential changes. We want the value at call time.
    // Recomputing the ref is not needed and will impact usage in hooks if it does.
    const lastDate = new Date(ref.current.lastRefreshTimestamp);
    const setNewDateString = ref.current.setLastRefreshTimestamp;

    /* eslint-disable no-console */
    // Print into the console in case we are not refreshing or the browser has preserve log enabled
    console.warn('Attempting to re-trigger an auto refresh');
    console.log('Last refresh was on:', lastDate);
    console.log('Refreshing requested after:', refreshDateMarker);

    lastDate.setHours(lastDate.getHours() + 1);
    if (lastDate < refreshDateMarker) {
      setNewDateString(refreshDateMarker.toString());
      console.log('Logging out and refreshing');
      logout().then(() => window.location.reload());
    } else {
      console.error(
        `We should have refreshed but it appears the last time we auto-refreshed was less than an hour ago. '${KEY_NAME}' session storage setting can be cleared for this to refresh again within the hour from the last refresh.`,
      );
    }
    /* eslint-enable no-console */
  }, []);
};

export default useTimeBasedRefresh;

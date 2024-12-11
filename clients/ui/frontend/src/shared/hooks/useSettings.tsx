import * as React from 'react';
import { mockedUsername, POLL_INTERVAL, USER_ID } from '~/shared/utilities/const';
import { useDeepCompareMemoize } from '~/shared/utilities/useDeepCompareMemoize';
import { ConfigSettings, UserSettings } from '~/shared/types';
import useTimeBasedRefresh from '~/shared/hooks/useTimeBasedRefresh';
import { getUser } from '~/shared/api/k8s';

export const useSettings = (): {
  configSettings: ConfigSettings | null;
  userSettings: UserSettings | null;
  loaded: boolean;
  loadError: Error | undefined;
} => {
  const [loaded, setLoaded] = React.useState(false);
  const [loadError, setLoadError] = React.useState<Error>();
  const [config, setConfig] = React.useState<ConfigSettings | null>(null);
  const [user, setUser] = React.useState<UserSettings | null>(null);
  const userSettings = React.useMemo(() => getUser(''), []);
  const setRefreshMarker = useTimeBasedRefresh();

  React.useEffect(() => {
    let watchHandle: ReturnType<typeof setTimeout>;
    let cancelled = false;
    const watchConfig = () => {
      // TODO: [Env Handling] Add mocked mode for frontend in dev
      // const headers = process.env.mocked === 'true' ? { [USER_ID]: mockedUsername } : undefined;
      const headers = { [USER_ID]: mockedUsername };
      Promise.all([fetchConfig(), userSettings({ headers })])
        .then(([fetchedConfig, fetchedUser]) => {
          if (cancelled) {
            return;
          }
          setConfig(fetchedConfig);
          setUser(fetchedUser);
          setLoaded(true);
          setLoadError(undefined);
        })
        .catch((e) => {
          if (e?.message?.includes('Error getting Oauth Info for user')) {
            // NOTE: this endpoint only requests oauth because of the security layer, this is not an ironclad use-case
            // Something went wrong on the server with the Oauth, let us just log them out and refresh for them
            /* eslint-disable-next-line no-console */
            console.error(
              'Something went wrong with the oauth token, please log out...',
              e.message,
              e,
            );
            setRefreshMarker(new Date());
            return;
          }
          setLoadError(e);
        });
      watchHandle = setTimeout(watchConfig, POLL_INTERVAL);
    };
    watchConfig();

    return () => {
      cancelled = true;
      clearTimeout(watchHandle);
    };
  }, [setRefreshMarker, userSettings]);

  const retConfig = useDeepCompareMemoize<ConfigSettings | null>(config);
  const retUser = useDeepCompareMemoize<UserSettings | null>(user);

  return { configSettings: retConfig, userSettings: retUser, loaded, loadError };
};

// Mock a settings config call
// TODO: [Data Flow] replace with the actual call once we have the endpoint
export const fetchConfig = async (): Promise<ConfigSettings> => ({
  common: {
    featureFlags: {
      modelRegistry: true,
    },
  },
});

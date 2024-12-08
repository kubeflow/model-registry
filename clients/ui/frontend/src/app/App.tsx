import * as React from 'react';
import '@patternfly/react-core/dist/styles/base.css';
import './app.css';
import {
  Alert,
  Bullseye,
  Button,
  Masthead,
  MastheadContent,
  MastheadMain,
  Page,
  PageSection,
  Spinner,
  Stack,
  StackItem,
} from '@patternfly/react-core';
import ToastNotifications from '~/shared/components/ToastNotifications';
import { useSettings } from '~/shared/hooks/useSettings';
import { isMUITheme, Theme } from '~/shared/utilities/const';
import NavSidebar from './NavSidebar';
import AppRoutes from './AppRoutes';
import { AppContext } from './AppContext';
import { ModelRegistrySelectorContextProvider } from './context/ModelRegistrySelectorContext';

const App: React.FC = () => {
  const {
    configSettings,
    userSettings,
    loaded: configLoaded,
    loadError: configError,
  } = useSettings();

  React.useEffect(() => {
    // Apply the theme based on the value of STYLE_THEME
    if (isMUITheme()) {
      document.documentElement.classList.add(Theme.MUI);
    } else {
      document.documentElement.classList.remove(Theme.MUI);
    }
  }, []);

  const contextValue = React.useMemo(
    () =>
      configSettings && userSettings
        ? {
            config: configSettings!,
            user: userSettings!,
          }
        : null,
    [configSettings, userSettings],
  );

  // We lack the critical data to startup the app
  if (configError) {
    // There was an error fetching critical data
    return (
      <Page>
        <PageSection>
          <Stack hasGutter>
            <StackItem>
              <Alert variant="danger" isInline title="General loading error">
                <p>{configError.message || 'Unknown error occurred during startup.'}</p>
                <p>Logging out and logging back in may solve the issue.</p>
              </Alert>
            </StackItem>
            <StackItem>
              <Button
                variant="secondary"
                onClick={() => {
                  // TODO: [Auth-enablement] Logout when auth is enabled
                }}
              >
                Logout
              </Button>
            </StackItem>
          </Stack>
        </PageSection>
      </Page>
    );
  }

  // Waiting on the API to finish
  const loading = !configLoaded || !userSettings || !configSettings || !contextValue;

  const masthead = (
    <Masthead>
      <MastheadMain />
      <MastheadContent>
        {/* TODO: [Auth-enablement] Add logout and user status once we enable itNavigates to register page from table toolbar */}
      </MastheadContent>
    </Masthead>
  );

  return loading ? (
    <Bullseye>
      <Spinner />
    </Bullseye>
  ) : (
    <AppContext.Provider value={contextValue}>
      <Page
        mainContainerId="primary-app-container"
        masthead={masthead}
        isManagedSidebar
        sidebar={<NavSidebar />}
      >
        <ModelRegistrySelectorContextProvider>
          <AppRoutes />
        </ModelRegistrySelectorContextProvider>
        <ToastNotifications />
      </Page>
    </AppContext.Provider>
  );
};

export default App;

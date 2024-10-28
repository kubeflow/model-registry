import * as React from 'react';
import '@patternfly/react-core/dist/styles/base.css';
import './app.css';
import {
  Alert,
  Brand,
  Bullseye,
  Button,
  Masthead,
  MastheadBrand,
  MastheadContent,
  MastheadMain,
  MastheadToggle,
  Page,
  PageSection,
  PageToggleButton,
  Spinner,
  Stack,
  StackItem,
} from '@patternfly/react-core';
import { BarsIcon } from '@patternfly/react-icons';
import ToastNotifications from '~/components/ToastNotifications';
import NavSidebar from './NavSidebar';
import AppRoutes from './AppRoutes';
import { AppContext } from './AppContext';
import { useSettings } from './useSettings';
import { ModelRegistrySelectorContextProvider } from './context/ModelRegistrySelectorContext';

const App: React.FC = () => {
  const {
    configSettings,
    userSettings,
    loaded: configLoaded,
    loadError: configError,
  } = useSettings();

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
                  // TODO: logout
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
      <MastheadMain>
        <MastheadToggle>
          <PageToggleButton id="page-nav-toggle" variant="plain" aria-label="Dashboard navigation">
            <BarsIcon />
          </PageToggleButton>
        </MastheadToggle>
        <MastheadBrand>
          <Brand
            className="kubeflow_brand"
            src={`${window.location.origin}/images/logo.svg`}
            alt="Kubeflow Logo"
          />
        </MastheadBrand>
      </MastheadMain>

      <MastheadContent>{/* TODO: Change this into a component for Header Tools */}</MastheadContent>
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

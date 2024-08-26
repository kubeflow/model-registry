import * as React from 'react';
import '@patternfly/react-core/dist/styles/base.css';
import AppRoutes from './AppRoutes';
import './app.css';
import {
  Alert,
  Bullseye,
  Button,
  Flex,
  Masthead,
  MastheadContent,
  MastheadToggle,
  Page,
  PageSection,
  PageToggleButton,
  Spinner,
  Stack,
  StackItem,
  Title,
} from "@patternfly/react-core";
import NavSidebar from "./NavSidebar";
import { BarsIcon } from "@patternfly/react-icons";
import { AppContext } from "./AppContext";
import { useSettings } from "./useSettings";

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
    [configSettings, userSettings]
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
                <p>
                  {configError.message ||
                    "Unknown error occurred during startup."}
                </p>
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
      <MastheadToggle>
        <PageToggleButton
          id="page-nav-toggle"
          variant="plain"
          aria-label="Dashboard navigation"
        >
          <BarsIcon />
        </PageToggleButton>
      </MastheadToggle>

      <MastheadContent>
        <Flex>
          <Title headingLevel="h2" size="3xl">
            Kubeflow Model Registry UI
          </Title>
        </Flex>
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
        <AppRoutes />
      </Page>
    </AppContext.Provider>
  );
};

export default App;

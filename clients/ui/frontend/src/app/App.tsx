import * as React from 'react';
import '@patternfly/react-core/dist/styles/base.css';
import './app.css';
import {
  Alert,
  Bullseye,
  Button,
  Dropdown,
  DropdownItem,
  DropdownList,
  Masthead,
  MastheadContent,
  MastheadMain,
  MenuToggle,
  MenuToggleElement,
  Page,
  PageSection,
  Spinner,
  Stack,
  StackItem,
  Toolbar,
  ToolbarContent,
  ToolbarGroup,
  ToolbarItem,
} from '@patternfly/react-core';
import ToastNotifications from '~/shared/components/ToastNotifications';
import { useSettings } from '~/shared/hooks/useSettings';
import { isMUITheme, Theme, USER_ID } from '~/shared/utilities/const';
import NavSidebar from './NavSidebar';
import AppRoutes from './AppRoutes';
import { AppContext } from './AppContext';
import { ModelRegistrySelectorContextProvider } from './context/ModelRegistrySelectorContext';
import { Select } from '@mui/material';
import { SimpleSelect, SimpleSelectOption } from '@patternfly/react-templates';

const App: React.FC = () => {
  const {
    configSettings,
    userSettings,
    loaded: configLoaded,
    loadError: configError,
  } = useSettings();

  const username = userSettings?.username;

  React.useEffect(() => {
    // Apply the theme based on the value of STYLE_THEME
    if (isMUITheme()) {
      document.documentElement.classList.add(Theme.MUI);
    } else {
      document.documentElement.classList.remove(Theme.MUI);
    }
  }, []);

  React.useEffect(() => {
    // Add the user to localStorage if in PoC
    // TODO: [Env Handling] Remove this when auth is enabled
    if (username) {
      localStorage.setItem(USER_ID, username);
    } else {
      localStorage.removeItem(USER_ID);
    }
  }, [username]);

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

  const handleLogout = () => {
    setUserMenuOpen(false);
    // TODO: [Auth-enablement] Logout when auth is enabled
  };


  const [userMenuOpen, setUserMenuOpen] = React.useState(false);
  const userMenuItems = [
    <DropdownItem key="logout" onClick={handleLogout}>
      Log out
    </DropdownItem>,
  ];

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

  const Options: SimpleSelectOption[] = [
    { content: 'All Namespaces', value: 'All' },
  ];

  const [selected, setSelected] = React.useState<string | undefined>('All');

  const initialOptions = React.useMemo<SimpleSelectOption[]>(
    () => Options.map((o) => ({ ...o, selected: o.value === selected })),
    [selected]
  );

  const masthead = (
    <Masthead>
      <MastheadMain />
      <MastheadContent>
        <Toolbar>
          <ToolbarContent>
            <ToolbarGroup variant="action-group-plain" align={{ default: 'alignStart' }}>
              <ToolbarItem>
                <SimpleSelect
                  isDisabled
                  initialOptions={initialOptions}
                  onSelect={(_ev, selection) => setSelected(String(selection))}
                >
                </SimpleSelect>
              </ToolbarItem>
            </ToolbarGroup>
            <ToolbarGroup variant="action-group-plain" align={{ default: 'alignEnd' }}>
              <ToolbarItem>
                {/* TODO: [Auth-enablement] Add logout button */}
                <Dropdown
                  popperProps={{ position: 'right' }}
                  onOpenChange={(isOpen) => setUserMenuOpen(isOpen)}
                  toggle={(toggleRef: React.Ref<MenuToggleElement>) => (
                    <MenuToggle
                      aria-label="User menu"
                      id="user-menu-toggle"
                      ref={toggleRef}
                      onClick={() => setUserMenuOpen(!userMenuOpen)}
                      isExpanded={userMenuOpen}
                    >
                      {username}
                    </MenuToggle>
                  )}
                  isOpen={userMenuOpen}
                >
                  <DropdownList>{userMenuItems}</DropdownList>
                </Dropdown>
              </ToolbarItem>
            </ToolbarGroup>
          </ToolbarContent>
        </Toolbar>
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

import React from 'react';
import {
  Brand,
  Dropdown,
  DropdownItem,
  DropdownList,
  Masthead,
  MastheadBrand,
  MastheadContent,
  MastheadLogo,
  MastheadMain,
  MastheadToggle,
  MenuToggle,
  MenuToggleElement,
  PageToggleButton,
  Toolbar,
  ToolbarContent,
  ToolbarGroup,
  ToolbarItem,
} from '@patternfly/react-core';
import { SimpleSelect } from '@patternfly/react-templates';
import { BarsIcon } from '@patternfly/react-icons';
import { useNamespaceSelector, useModularArchContext } from 'mod-arch-core';
import { useThemeContext } from 'mod-arch-kubeflow';
import { images as sharedImages } from 'mod-arch-shared';

interface NavBarProps {
  username?: string;
  onLogout: () => void;
}

const NavBar: React.FC<NavBarProps> = ({ username, onLogout }) => {
  const { namespaces, preferredNamespace, updatePreferredNamespace } = useNamespaceSelector();
  const { config } = useModularArchContext();
  const { isMUITheme } = useThemeContext();

  const [userMenuOpen, setUserMenuOpen] = React.useState(false);

  // Check if mandatory namespace is configured
  const isMandatoryNamespace = Boolean(config.mandatoryNamespace);

  const options = namespaces.map((namespace) => ({
    content: namespace.name,
    value: namespace.name,
    selected: namespace.name === preferredNamespace?.name,
  }));

  const handleLogout = () => {
    setUserMenuOpen(false);
    onLogout();
  };

  const userMenuItems = [
    <DropdownItem key="logout" onClick={handleLogout}>
      Log out
    </DropdownItem>,
  ];

  return (
    <Masthead>
      <MastheadMain>
        <MastheadToggle>
          <PageToggleButton id="page-nav-toggle" variant="plain" aria-label="Dashboard navigation">
            <BarsIcon />
          </PageToggleButton>
        </MastheadToggle>
        {!isMUITheme ? (
          <MastheadBrand>
            <MastheadLogo component="a">
              <Brand
                src={sharedImages.logoLightThemePath}
                alt="Kubeflow"
                heights={{ default: '36px' }}
              />
            </MastheadLogo>
          </MastheadBrand>
        ) : null}
      </MastheadMain>
      <MastheadContent>
        <Toolbar>
          <ToolbarContent>
            <ToolbarGroup variant="action-group-plain" align={{ default: 'alignStart' }}>
              <ToolbarItem className="kubeflow-u-namespace-select">
                <SimpleSelect
                  initialOptions={options}
                  isDisabled={isMandatoryNamespace} // Disable selection when mandatory namespace is set
                  onSelect={(_ev, selection) => {
                    // Only allow selection if not mandatory namespace
                    if (!isMandatoryNamespace) {
                      updatePreferredNamespace({ name: String(selection) });
                    }
                  }}
                />
              </ToolbarItem>
            </ToolbarGroup>
            {username && (
              <ToolbarGroup variant="action-group-plain" align={{ default: 'alignEnd' }}>
                <ToolbarItem>
                  <Dropdown
                    popperProps={{ position: 'right' }}
                    onOpenChange={(isOpen) => setUserMenuOpen(isOpen)}
                    toggle={(toggleRef: React.Ref<MenuToggleElement>) => (
                      <MenuToggle
                        aria-label="User menu"
                        id="user-menu-toggle"
                        data-testid="user-menu-toggle-button"
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
            )}
          </ToolbarContent>
        </Toolbar>
      </MastheadContent>
    </Masthead>
  );
};

export default NavBar;

import React from 'react';
import {
  Dropdown,
  DropdownItem,
  DropdownList,
  Masthead,
  MastheadContent,
  MastheadMain,
  MenuToggle,
  MenuToggleElement,
  Toolbar,
  ToolbarContent,
  ToolbarGroup,
  ToolbarItem,
} from '@patternfly/react-core';
import { SimpleSelect } from '@patternfly/react-templates';
import { NamespaceSelectorContext } from '~/shared/context/NamespaceSelectorContext';

interface NavBarProps {
  username?: string;
  onLogout: () => void;
}

const NavBar: React.FC<NavBarProps> = ({ username, onLogout }) => {
  const { namespaces, preferredNamespace, updatePreferredNamespace } =
    React.useContext(NamespaceSelectorContext);

  const [userMenuOpen, setUserMenuOpen] = React.useState(false);

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
      <MastheadMain />
      <MastheadContent>
        <Toolbar>
          <ToolbarContent>
            <ToolbarGroup variant="action-group-plain" align={{ default: 'alignStart' }}>
              <ToolbarItem>
                <SimpleSelect
                  initialOptions={options}
                  onSelect={(_ev, selection) => {
                    updatePreferredNamespace({ name: String(selection) });
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

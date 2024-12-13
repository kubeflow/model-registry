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
import { SimpleSelect, SimpleSelectOption } from '@patternfly/react-templates';

interface NavBarProps {
  username?: string;
  onLogout: () => void;
}

const Options: SimpleSelectOption[] = [{ content: 'All Namespaces', value: 'All' }];

const NavBar: React.FC<NavBarProps> = ({ username, onLogout }) => {
  const [selected, setSelected] = React.useState<string | undefined>('All');
  const [userMenuOpen, setUserMenuOpen] = React.useState(false);

  const initialOptions = React.useMemo<SimpleSelectOption[]>(
    () => Options.map((o) => ({ ...o, selected: o.value === selected })),
    [selected],
  );

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
                  isDisabled
                  initialOptions={initialOptions}
                  onSelect={(_ev, selection) => setSelected(String(selection))}
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

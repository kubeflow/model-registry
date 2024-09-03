import * as React from 'react';
import { NavLink } from 'react-router-dom';
import {
  Nav,
  NavExpandable,
  NavItem,
  NavList,
  PageSidebar,
  PageSidebarBody,
} from '@patternfly/react-core';
import { useNavData, isNavDataGroup, NavDataHref, NavDataGroup } from './AppRoutes';

const NavHref: React.FC<{ item: NavDataHref }> = ({ item }) => (
  <NavItem key={item.label} data-id={item.label} itemId={item.label}>
    <NavLink to={item.path}>{item.label}</NavLink>
  </NavItem>
);

const NavGroup: React.FC<{ item: NavDataGroup }> = ({ item }) => {
  const { children } = item;
  const [expanded, setExpanded] = React.useState(false);

  return (
    <NavExpandable
      data-id={item.label}
      key={item.label}
      id={item.label}
      title={item.label}
      groupId={item.label}
      isExpanded={expanded}
      onExpand={(e, val) => setExpanded(val)}
      aria-label={item.label}
    >
      {children.map((childItem) => (
        <NavHref key={childItem.label} data-id={childItem.label} item={childItem} />
      ))}
    </NavExpandable>
  );
};

const NavSidebar: React.FC = () => {
  const navData = useNavData();

  return (
    <PageSidebar>
      <PageSidebarBody>
        <Nav id="nav-primary-simple">
          <NavList id="nav-list-simple">
            {navData.map((item) =>
              isNavDataGroup(item) ? (
                <NavGroup key={item.label} item={item} />
              ) : (
                <NavHref key={item.label} item={item} />
              ),
            )}
          </NavList>
        </Nav>
      </PageSidebarBody>
    </PageSidebar>
  );
};

export default NavSidebar;

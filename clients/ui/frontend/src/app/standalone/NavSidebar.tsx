import * as React from 'react';
import { NavLink } from 'react-router-dom';
import {
  Brand,
  Nav,
  NavExpandable,
  NavItem,
  NavList,
  PageSidebar,
  PageSidebarBody,
} from '@patternfly/react-core';
import { useThemeContext, images as kubeflowImages } from 'mod-arch-kubeflow';
import { isNavDataGroup, NavDataHref, NavDataGroup, NavDataItem } from '~/app/standalone/types';

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

export type NavSidebarProps = {
  navData: NavDataItem[];
};

const NavSidebar: React.FC<NavSidebarProps> = ({ navData }) => {
  const { isMUITheme } = useThemeContext();
  return (
    <PageSidebar>
      <PageSidebarBody>
        <Nav id="nav-primary-simple">
          <NavList id="nav-list-simple">
            {isMUITheme ? (
              <NavItem>
                <Brand
                  className="kubeflow_brand"
                  src={kubeflowImages.logoKubeflowLight}
                  alt="Kubeflow Logo"
                />
              </NavItem>
            ) : null}
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

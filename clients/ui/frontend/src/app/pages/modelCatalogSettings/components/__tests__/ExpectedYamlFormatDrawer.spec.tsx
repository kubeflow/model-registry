import * as React from 'react';
import { screen, render } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';
import '@testing-library/jest-dom';
import {
  Drawer,
  DrawerContent,
  DrawerContentBody,
  DrawerPanelContent,
} from '@patternfly/react-core';
import { ExpectedYamlFormatDrawerPanel } from '~/app/pages/modelCatalogSettings/components/ExpectedYamlFormatDrawer';
import { EXPECTED_YAML_FORMAT_LABEL } from '~/app/pages/modelCatalogSettings/constants';

const DRAWER_TITLE = EXPECTED_YAML_FORMAT_LABEL;

const renderWithDrawer = (onClose: jest.Mock, isExpanded = true) =>
  render(
    <Drawer isExpanded={isExpanded} isInline>
      <DrawerContent
        panelContent={
          <DrawerPanelContent>
            <ExpectedYamlFormatDrawerPanel onClose={onClose} />
          </DrawerPanelContent>
        }
      >
        <DrawerContentBody>Page content</DrawerContentBody>
      </DrawerContent>
    </Drawer>,
  );

describe('ExpectedYamlFormatDrawerPanel', () => {
  const onClose = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders drawer with title and close button when expanded', () => {
    renderWithDrawer(onClose);
    expect(screen.getByTestId('expected-format-drawer-title')).toHaveTextContent(DRAWER_TITLE);
    expect(screen.getByRole('button', { name: 'Close drawer' })).toBeInTheDocument();
  });

  it('calls onClose when close button is clicked', async () => {
    const user = userEvent.setup();
    renderWithDrawer(onClose);
    await user.click(screen.getByRole('button', { name: 'Close drawer' }));
    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('displays expected YAML format content when expanded', () => {
    renderWithDrawer(onClose);
    expect(document.body.textContent).toContain('source:');
    expect(document.body.textContent).toContain('models:');
  });

  it('renders page content alongside the drawer panel', () => {
    renderWithDrawer(onClose);
    expect(screen.getByText('Page content')).toBeInTheDocument();
  });
});

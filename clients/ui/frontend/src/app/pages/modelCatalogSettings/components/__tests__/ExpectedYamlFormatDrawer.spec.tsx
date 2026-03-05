import React from 'react';
import { screen, render, waitFor } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';
import '@testing-library/jest-dom';
import ExpectedYamlFormatDrawer from '~/app/pages/modelCatalogSettings/components/ExpectedYamlFormatDrawer';
import { EXPECTED_YAML_FORMAT_LABEL } from '~/app/pages/modelCatalogSettings/constants';

const DRAWER_TITLE = EXPECTED_YAML_FORMAT_LABEL;
const PRIMARY_APP_CONTAINER_ID = 'primary-app-container';

describe('ExpectedYamlFormatDrawer', () => {
  let container: HTMLElement;
  const onClose = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
    global.ResizeObserver = jest.fn().mockImplementation(() => ({
      observe: jest.fn(),
      unobserve: jest.fn(),
      disconnect: jest.fn(),
    }));
    container = document.createElement('div');
    container.id = PRIMARY_APP_CONTAINER_ID;
    document.body.appendChild(container);
  });

  afterEach(() => {
    if (container.parentNode) {
      document.body.removeChild(container);
    }
  });

  it('does not render drawer panel when isOpen is false', () => {
    render(
      <ExpectedYamlFormatDrawer isOpen={false} onClose={onClose}>
        <div>Page content</div>
      </ExpectedYamlFormatDrawer>,
    );
    expect(screen.queryByRole('region', { name: DRAWER_TITLE })).not.toBeInTheDocument();
    expect(screen.getByText('Page content')).toBeInTheDocument();
  });

  it('renders drawer with title and close button when isOpen is true', async () => {
    render(
      <ExpectedYamlFormatDrawer isOpen onClose={onClose}>
        <div>Page content</div>
      </ExpectedYamlFormatDrawer>,
    );

    await waitFor(() => {
      expect(screen.getByRole('region', { name: DRAWER_TITLE })).toBeInTheDocument();
    });

    expect(screen.getByTestId('expected-format-drawer-title')).toHaveTextContent(DRAWER_TITLE);
    expect(screen.getByRole('button', { name: 'Close drawer' })).toBeInTheDocument();
  });

  it('calls onClose when close button is clicked', async () => {
    const user = userEvent.setup();
    render(
      <ExpectedYamlFormatDrawer isOpen onClose={onClose}>
        <div>Page content</div>
      </ExpectedYamlFormatDrawer>,
    );

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Close drawer' })).toBeInTheDocument();
    });

    await user.click(screen.getByRole('button', { name: 'Close drawer' }));

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('displays expected YAML format content when open', async () => {
    render(
      <ExpectedYamlFormatDrawer isOpen onClose={onClose}>
        <div>Page content</div>
      </ExpectedYamlFormatDrawer>,
    );

    await waitFor(() => {
      expect(screen.getByRole('region', { name: DRAWER_TITLE })).toBeInTheDocument();
    });

    expect(document.body.textContent).toContain('source:');
    expect(document.body.textContent).toContain('models:');
  });
});

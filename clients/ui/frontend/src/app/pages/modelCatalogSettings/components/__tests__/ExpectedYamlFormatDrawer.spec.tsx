import React from 'react';
import { screen, render, waitFor } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';
import '@testing-library/jest-dom';
import ExpectedYamlFormatDrawer from '~/app/pages/modelCatalogSettings/components/ExpectedYamlFormatDrawer';
import {
  EXPECTED_YAML_FORMAT_LABEL,
  PRIMARY_APP_CONTAINER_ID,
} from '~/app/pages/modelCatalogSettings/constants';

const DRAWER_TITLE = EXPECTED_YAML_FORMAT_LABEL;

const mockUseThemeContext = jest.fn(() => ({ isMUITheme: false }));
jest.mock('mod-arch-kubeflow', () => ({
  useThemeContext: () => mockUseThemeContext(),
}));

describe('ExpectedYamlFormatDrawer', () => {
  let container: HTMLElement;
  const onClose = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
    mockUseThemeContext.mockReturnValue({ isMUITheme: false });
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

  it('when #primary-app-container is missing, renders only children, does not render drawer, and warns', () => {
    const warnSpy = jest.spyOn(console, 'warn').mockImplementation();
    container.remove();
    render(
      <ExpectedYamlFormatDrawer isOpen onClose={onClose}>
        <div>Page content</div>
      </ExpectedYamlFormatDrawer>,
    );
    expect(screen.queryByRole('region', { name: DRAWER_TITLE })).not.toBeInTheDocument();
    expect(screen.getByText('Page content')).toBeInTheDocument();
    expect(warnSpy).toHaveBeenCalledWith(
      expect.stringContaining(PRIMARY_APP_CONTAINER_ID),
    );
    warnSpy.mockRestore();
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

  it('when MUI theme is enabled, drawer renders and MutationObserver is used for panel max-width override', async () => {
    mockUseThemeContext.mockReturnValue({ isMUITheme: true });
    const observeSpy = jest.spyOn(MutationObserver.prototype, 'observe');

    render(
      <ExpectedYamlFormatDrawer isOpen onClose={onClose}>
        <div>Page content</div>
      </ExpectedYamlFormatDrawer>,
    );

    await waitFor(() => {
      expect(screen.getByRole('region', { name: DRAWER_TITLE })).toBeInTheDocument();
    });

    expect(observeSpy).toHaveBeenCalled();

    observeSpy.mockRestore();
  });
});

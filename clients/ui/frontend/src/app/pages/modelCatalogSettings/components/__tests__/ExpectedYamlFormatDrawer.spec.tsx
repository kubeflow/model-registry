import React from 'react';
import { screen, render } from '@testing-library/react';
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
    container = document.createElement('div');
    container.id = PRIMARY_APP_CONTAINER_ID;
    document.body.appendChild(container);
  });

  afterEach(() => {
    container.remove();
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

  it('when container is missing, renders only children and does not render drawer', () => {
    container.remove();
    render(
      <ExpectedYamlFormatDrawer isOpen onClose={onClose}>
        <div>Page content</div>
      </ExpectedYamlFormatDrawer>,
    );
    expect(screen.queryByRole('region', { name: DRAWER_TITLE })).not.toBeInTheDocument();
    expect(screen.getByText('Page content')).toBeInTheDocument();
  });

  it('renders drawer with title and close button when isOpen is true', () => {
    render(
      <ExpectedYamlFormatDrawer isOpen onClose={onClose}>
        <div>Page content</div>
      </ExpectedYamlFormatDrawer>,
    );
    expect(screen.getByRole('region', { name: DRAWER_TITLE })).toBeInTheDocument();
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
    await user.click(screen.getByRole('button', { name: 'Close drawer' }));
    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('displays expected YAML format content when open', () => {
    render(
      <ExpectedYamlFormatDrawer isOpen onClose={onClose}>
        <div>Page content</div>
      </ExpectedYamlFormatDrawer>,
    );
    expect(screen.getByRole('region', { name: DRAWER_TITLE })).toBeInTheDocument();
    expect(document.body.textContent).toContain('source:');
    expect(document.body.textContent).toContain('models:');
  });

  it('renders drawer when MUI theme is enabled', () => {
    mockUseThemeContext.mockReturnValue({ isMUITheme: true });
    render(
      <ExpectedYamlFormatDrawer isOpen onClose={onClose}>
        <div>Page content</div>
      </ExpectedYamlFormatDrawer>,
    );
    expect(screen.getByRole('region', { name: DRAWER_TITLE })).toBeInTheDocument();
    expect(screen.getByTestId('expected-format-drawer-title')).toHaveTextContent(DRAWER_TITLE);
  });
});

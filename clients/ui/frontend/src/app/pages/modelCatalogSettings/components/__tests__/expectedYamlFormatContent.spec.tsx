import React from 'react';
import { screen, render, waitFor } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';
import '@testing-library/jest-dom';
import ExpectedYamlFormatDrawer from '~/app/pages/modelCatalogSettings/components/expectedYamlFormatContent';

const DRAWER_TITLE = 'View expected file format';
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
    container.getBoundingClientRect = jest.fn().mockReturnValue({
      top: 0,
      left: 0,
      width: 800,
      height: 600,
      bottom: 600,
      right: 800,
      x: 0,
      y: 0,
      toJSON: () => ({}),
    });
  });

  afterEach(() => {
    if (container.parentNode) {
      document.body.removeChild(container);
    }
  });

  it('returns null when isOpen is false', () => {
    const { container: renderContainer } = render(
      <ExpectedYamlFormatDrawer isOpen={false} onClose={onClose} />,
    );
    expect(renderContainer.firstChild).toBeNull();
    expect(screen.queryByRole('region', { name: DRAWER_TITLE })).not.toBeInTheDocument();
  });

  it('renders drawer with title and close button when isOpen is true', async () => {
    render(<ExpectedYamlFormatDrawer isOpen onClose={onClose} />);

    await waitFor(() => {
      expect(screen.getByRole('region', { name: DRAWER_TITLE })).toBeInTheDocument();
    });

    expect(screen.getByTestId('expected-format-drawer-title')).toHaveTextContent(DRAWER_TITLE);
    expect(screen.getByRole('button', { name: 'Close drawer' })).toBeInTheDocument();
  });

  it('calls onClose when close button is clicked', async () => {
    const user = userEvent.setup();
    render(<ExpectedYamlFormatDrawer isOpen onClose={onClose} />);

    await waitFor(() => {
      expect(screen.getByTestId('expected-format-drawer-close')).toBeInTheDocument();
    });

    await user.click(screen.getByTestId('expected-format-drawer-close'));

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('displays expected YAML format content when open', async () => {
    render(<ExpectedYamlFormatDrawer isOpen onClose={onClose} />);

    await waitFor(() => {
      expect(screen.getByRole('region', { name: DRAWER_TITLE })).toBeInTheDocument();
    });

    expect(document.body.textContent).toContain('source:');
    expect(document.body.textContent).toContain('models:');
  });
});

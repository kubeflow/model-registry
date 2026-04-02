import React from 'react';

import { screen, render, waitFor } from '@testing-library/react';

import { userEvent } from '@testing-library/user-event';
import '@testing-library/jest-dom';

import { ModelRegistryTableRowStatus } from '~/app/pages/settings/ModelRegistryTableRowStatus';

describe('ModelRegistryTableRowStatus', () => {
  it('renders "Unavailable" status with correct popover for Istio and Gateway conditions', async () => {
    const user = userEvent.setup();

    render(
      <ModelRegistryTableRowStatus
        conditions={[
          {
            status: 'False',
            type: 'Available',
            message: 'Service is unavailable',
          },
          {
            status: 'False',
            type: 'IstioAvailable',
            message: 'Istio is unavailable',
          },
          {
            status: 'False',
            type: 'GatewayAvailable',
            message: 'Gateway is unavailable',
          },
        ]}
      />,
    );

    const label = screen.getByText('Unavailable');
    expect(label).toBeVisible();

    await user.click(label);

    await waitFor(() => {
      // Check for the popover title
      expect(
        screen.getByText('Istio resources and Istio Gateway resources are both unavailable'),
      ).toBeInTheDocument();

      // Check for the condition messages
      expect(screen.getByText('Service is unavailable')).toBeInTheDocument();
      expect(screen.getByText('Istio is unavailable')).toBeInTheDocument();
      expect(screen.getByText('Gateway is unavailable')).toBeInTheDocument();
    });
  });
  it('renders "Istio resources and Istio Gateway resources are both unavailable" as popover title', async () => {
    const user = userEvent.setup();

    render(
      <ModelRegistryTableRowStatus
        conditions={[
          { status: 'False', type: 'Available', message: 'Some unavailable message' },
          {
            status: 'False',
            type: 'IstioAvailable',
            message: 'Some istio unavailable message',
          },
          {
            status: 'False',
            type: 'GatewayAvailable',
            message: 'Some gateway unavailable message',
          },
        ]}
      />,
    );

    await user.click(screen.getByText('Unavailable'));
    expect(
      screen.getByRole('heading', {
        name: 'danger alert: Istio resources and Istio Gateway resources are both unavailable',
      }),
    ).toBeVisible();
  });

  it('renders "Istio resources are unavailable" as popover title', async () => {
    const user = userEvent.setup();

    render(
      <ModelRegistryTableRowStatus
        conditions={[
          { status: 'False', type: 'Available', message: 'Some unavailable message' },
          {
            status: 'False',
            type: 'IstioAvailable',
            message: 'Some istio unavailable message',
          },
          {
            status: 'True',
            type: 'GatewayAvailable',
          },
        ]}
      />,
    );

    await user.click(screen.getByText('Unavailable'));
    expect(
      screen.getByRole('heading', { name: 'danger alert: Istio resources are unavailable' }),
    ).toBeVisible();
  });

  it('renders "Istio Gateway resources are unavailable" as popover title', async () => {
    const user = userEvent.setup();

    render(
      <ModelRegistryTableRowStatus
        conditions={[
          { status: 'False', type: 'Available', message: 'Some unavailable message' },
          {
            status: 'True',
            type: 'IstioAvailable',
          },
          {
            status: 'False',
            type: 'GatewayAvailable',
            message: 'Some gateway unavailable message',
          },
        ]}
      />,
    );

    await user.click(screen.getByText('Unavailable'));
    expect(
      screen.getByRole('heading', {
        name: 'danger alert: Istio Gateway resources are unavailable',
      }),
    ).toBeVisible();
  });

  it('renders "Deployment is unavailable" as popover title', async () => {
    const user = userEvent.setup();

    render(
      <ModelRegistryTableRowStatus
        conditions={[
          { status: 'False', type: 'Available', message: 'Some unavailable message' },
          {
            status: 'True',
            type: 'IstioAvailable',
          },
          {
            status: 'True',
            type: 'GatewayAvailable',
          },
        ]}
      />,
    );

    await user.click(screen.getByText('Unavailable'));
    expect(
      screen.getByRole('heading', { name: 'danger alert: Deployment is unavailable' }),
    ).toBeVisible();
  });

  it('renders "Ready" status', () => {
    render(
      <ModelRegistryTableRowStatus
        conditions={[
          { status: 'False', type: 'Degraded' },
          { status: 'True', type: 'Available' },
        ]}
      />,
    );
    expect(screen.getByText('Ready')).toBeVisible();
  });
  it('renders "Unavailable" status', async () => {
    const user = userEvent.setup();

    render(
      <ModelRegistryTableRowStatus
        conditions={[
          {
            status: 'True',
            type: 'Degraded',
          },
          {
            status: 'False',
            type: 'Available',
            message: 'Some unavailable message',
          },
        ]}
      />,
    );

    const label = screen.getByText('Unavailable');
    expect(label).toBeVisible();

    await user.click(label);

    expect(
      screen.getByRole('heading', { name: 'danger alert: Service is unavailable' }),
    ).toBeVisible();
    expect(screen.getByText('Some unavailable message')).toBeVisible();
  });
  it('renders "Degrading" status', async () => {
    const user = userEvent.setup();

    render(
      <ModelRegistryTableRowStatus
        conditions={[
          {
            status: 'True',
            type: 'Degraded',
          },
        ]}
      />,
    );

    const label = screen.getByText('Degrading');
    expect(label).toBeVisible();

    await user.click(label);

    const degradingText = screen.getByText(/degrading/i, { exact: false });
    expect(degradingText).toBeInTheDocument();
  });

  it('renders "Starting" status when popover message contains "ContainerCreating"', async () => {
    render(
      <ModelRegistryTableRowStatus
        conditions={[
          {
            status: 'False',
            type: 'Unavailable',
            message:
              'Deployment is unavailable: pod test has unready containers [grpc-container: {waiting: {reason: ContainerCreating, message: }}',
          },
        ]}
      />,
    );

    expect(screen.getByText('Starting')).toBeVisible();
  });

  it('renders "Starting" status when conditions are empty', () => {
    render(<ModelRegistryTableRowStatus conditions={[]} />);
    expect(screen.getByText('Starting')).toBeVisible();
  });

  it('renders "Starting" status when conditions are undefined', () => {
    render(<ModelRegistryTableRowStatus conditions={undefined} />);
    expect(screen.getByText('Starting')).toBeVisible();
  });

  it('renders "Unavailable" with multiple messages in popover', async () => {
    const user = userEvent.setup();

    render(
      <ModelRegistryTableRowStatus
        conditions={[
          {
            status: 'False',
            type: 'Progressing',
            message: 'Some unavailable message 1',
          },
          {
            status: 'False',
            type: 'Degraded',
            message: 'Some unavailable message 2',
          },
          {
            status: 'False',
            type: 'Available',
            message: 'Some unavailable message 3',
          },
        ]}
      />,
    );

    const label = screen.getByText('Unavailable');
    expect(label).toBeVisible();

    await user.click(label);

    expect(
      screen.getByRole('heading', { name: 'danger alert: Service is unavailable' }),
    ).toBeVisible();
    expect(screen.getByText('Some unavailable message 1')).toBeVisible();
    expect(screen.getByText('Some unavailable message 2')).toBeVisible();
    expect(screen.getByText('Some unavailable message 3')).toBeVisible();
  });

  it('renders "Stopping" status when deletionTimestamp is set', () => {
    render(
      <ModelRegistryTableRowStatus
        conditions={[{ status: 'True', type: 'Available' }]}
        deletionTimestamp="2024-03-22T10:00:00Z"
      />,
    );
    expect(screen.getByText('Stopping')).toBeVisible();
  });

  it('renders "Stopping" status even when conditions indicate Ready', () => {
    render(
      <ModelRegistryTableRowStatus
        conditions={[
          { status: 'True', type: 'Available' },
          { status: 'False', type: 'Degraded' },
        ]}
        deletionTimestamp="2024-03-22T10:00:00Z"
      />,
    );
    expect(screen.getByText('Stopping')).toBeVisible();
    expect(screen.queryByText('Ready')).not.toBeInTheDocument();
  });

  describe('label variant styling', () => {
    it('renders "Ready" label with outline variant (non-interactive)', () => {
      render(
        <ModelRegistryTableRowStatus
          conditions={[
            { status: 'False', type: 'Degraded' },
            { status: 'True', type: 'Available' },
          ]}
        />,
      );
      const labelEl = screen.getByTestId('model-registry-label');
      expect(labelEl.className).toMatch(/outline/);
      expect(labelEl.className).not.toMatch(/filled/);
    });

    it('renders "Starting" label with outline variant (non-interactive)', () => {
      render(<ModelRegistryTableRowStatus conditions={[]} />);
      const labelEl = screen.getByTestId('model-registry-label');
      expect(labelEl.className).toMatch(/outline/);
    });

    it('renders "Stopping" label with outline variant (non-interactive)', () => {
      render(
        <ModelRegistryTableRowStatus
          conditions={[{ status: 'True', type: 'Available' }]}
          deletionTimestamp="2024-03-22T10:00:00Z"
        />,
      );
      const labelEl = screen.getByTestId('model-registry-label');
      expect(labelEl.className).toMatch(/outline/);
    });

    it('renders "Unavailable" label with filled variant (interactive with popover)', () => {
      render(
        <ModelRegistryTableRowStatus
          conditions={[
            { status: 'False', type: 'Available', message: 'Service is unavailable' },
            { status: 'False', type: 'IstioAvailable', message: 'Istio is unavailable' },
          ]}
        />,
      );
      const labelEl = screen.getByTestId('model-registry-label');
      expect(labelEl.className).toMatch(/filled/);
      expect(labelEl.className).not.toMatch(/outline/);
    });
  });

  it('renders "Unavailable" with an unknown status', async () => {
    const user = userEvent.setup();

    render(
      <ModelRegistryTableRowStatus
        conditions={[
          {
            status: 'False',
            type: 'Available',
            message: 'Some unknown status message',
          },
          {
            status: 'True',
            type: 'Unknown',
          },
        ]}
      />,
    );

    const label = screen.getByText('Unavailable');
    expect(label).toBeVisible();

    await user.click(label);

    expect(
      screen.getByRole('heading', { name: 'danger alert: Service is unavailable' }),
    ).toBeVisible();
    expect(screen.getByText('Some unknown status message')).toBeVisible();
  });
});

import React from 'react';
import { screen, render, waitFor } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import '@testing-library/jest-dom';
import ManageSourcePage from '~/app/pages/modelCatalogSettings/screens/ManageSourcePage';
import {
  EXPECTED_YAML_FORMAT_LABEL,
  PRIMARY_APP_CONTAINER_ID,
} from '~/app/pages/modelCatalogSettings/constants';

jest.mock('mod-arch-shared', () => ({
  ApplicationsPage: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}));

jest.mock('mod-arch-kubeflow', () => ({
  useThemeContext: () => ({ isMUITheme: false }),
}));

jest.mock('~/app/hooks/modelCatalogSettings/useCatalogSourceConfigBySourceId', () => ({
  useCatalogSourceConfigBySourceId: () => [undefined, true, undefined],
}));

jest.mock('~/app/pages/modelCatalogSettings/components/ManageSourceForm', () => ({
  __esModule: true,
  default: function MockManageSourceForm({
    onToggleExpectedFormatDrawer,
  }: {
    onToggleExpectedFormatDrawer?: () => void;
  }) {
    return (
      <div>
        {onToggleExpectedFormatDrawer && (
          <button
            type="button"
            onClick={onToggleExpectedFormatDrawer}
            data-testid="open-expected-format-drawer"
          >
            View expected format
          </button>
        )}
      </div>
    );
  },
}));

describe('ManageSourcePage', () => {
  let portalContainer: HTMLElement;

  beforeEach(() => {
    jest.clearAllMocks();
    portalContainer = document.createElement('div');
    portalContainer.id = PRIMARY_APP_CONTAINER_ID;
    document.body.appendChild(portalContainer);
  });

  afterEach(() => {
    portalContainer.remove();
  });

  it('form link button opens drawer and close button closes it (open/close state wiring)', async () => {
    const user = userEvent.setup();
    render(
      <MemoryRouter>
        <ManageSourcePage />
      </MemoryRouter>,
    );

    expect(
      screen.queryByRole('region', { name: EXPECTED_YAML_FORMAT_LABEL }),
    ).not.toBeInTheDocument();

    await user.click(screen.getByTestId('open-expected-format-drawer'));

    await waitFor(() => {
      expect(screen.getByRole('region', { name: EXPECTED_YAML_FORMAT_LABEL })).toBeInTheDocument();
    });

    await user.click(screen.getByRole('button', { name: 'Close drawer' }));

    await waitFor(() => {
      expect(
        screen.queryByRole('region', { name: EXPECTED_YAML_FORMAT_LABEL }),
      ).not.toBeInTheDocument();
    });
  });
});

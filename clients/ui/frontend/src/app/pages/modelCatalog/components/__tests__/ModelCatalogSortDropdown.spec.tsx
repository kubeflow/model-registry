import '@testing-library/jest-dom';
import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import {
  ModelCatalogContextProvider,
  ModelCatalogContext,
} from '~/app/context/modelCatalog/ModelCatalogContext';
import { ModelCatalogSortOption } from '~/concepts/modelCatalog/const';
import ModelCatalogSortDropdown from '~/app/pages/modelCatalog/components/ModelCatalogSortDropdown';

// Mock mod-arch-core - must be before any imports that use it
jest.mock('mod-arch-core', () => {
  const actual = jest.requireActual('mod-arch-core');
  return {
    ...actual,
    useModularArchContext: jest.fn(() => ({
      config: {
        deploymentMode: 'standalone',
        URL_PREFIX: '/',
        BFF_API_VERSION: 'v1',
      },
    })),
    useNamespaceSelector: jest.fn(() => ({
      selectedNamespace: undefined,
      setSelectedNamespace: jest.fn(),
    })),
    useQueryParamNamespaces: jest.fn(() => ({})),
    asEnumMember: jest.fn((value, enumObj) =>
      value && enumObj && Object.values(enumObj).includes(value) ? value : undefined,
    ),
  };
});

jest.mock('react-router-dom', () => ({
  useLocation: jest.fn(() => ({ pathname: '/model-catalog' })),
}));

jest.mock('~/app/hooks/modelCatalog/useCatalogSources', () => ({
  useCatalogSources: jest.fn(() => [
    { items: [], size: 0, pageSize: 0, nextPageToken: '' },
    true,
    undefined,
  ]),
}));

jest.mock('~/app/hooks/modelCatalog/useCatalogFilterOptionList', () => ({
  useCatalogFilterOptionList: jest.fn(() => [null, true, undefined]),
}));

jest.mock('~/app/hooks/modelCatalog/useModelCatalogAPIState', () => ({
  __esModule: true,
  default: jest.fn(() => [
    {
      apiAvailable: true,
      api: {
        getCatalogModelsBySource: jest.fn(),
        getCatalogSourceModel: jest.fn(),
        getCatalogSourceModelArtifacts: jest.fn(),
        getCatalogModelPerformanceArtifacts: jest.fn(),
        getAllCatalogSources: jest.fn(),
        getCatalogFilterOptions: jest.fn(),
        createCatalogSourcePreview: jest.fn(),
      },
    },
    jest.fn(),
  ]),
}));

// Mock the getActiveLatencyFieldName utility
jest.mock('~/app/pages/modelCatalog/utils/modelCatalogUtils', () => ({
  ...jest.requireActual('~/app/pages/modelCatalog/utils/modelCatalogUtils'),
  getActiveLatencyFieldName: jest.fn(),
}));

import { getActiveLatencyFieldName } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

const mockGetActiveLatencyFieldName = getActiveLatencyFieldName as jest.MockedFunction<
  typeof getActiveLatencyFieldName
>;

describe('ModelCatalogSortDropdown', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockGetActiveLatencyFieldName.mockReturnValue(undefined);
  });

  describe('Visibility', () => {
    it('should not render when performance view is disabled', () => {
      render(
        <ModelCatalogContextProvider>
          <ModelCatalogSortDropdown performanceViewEnabled={false} />
        </ModelCatalogContextProvider>,
      );

      expect(screen.queryByText('Sort:')).not.toBeInTheDocument();
      expect(screen.queryByTestId('model-catalog-sort-dropdown')).not.toBeInTheDocument();
    });

    it('should render when performance view is enabled', () => {
      render(
        <ModelCatalogContextProvider>
          <ModelCatalogSortDropdown performanceViewEnabled />
        </ModelCatalogContextProvider>,
      );

      expect(screen.getByText('Sort:')).toBeInTheDocument();
      expect(screen.getByTestId('model-catalog-sort-dropdown')).toBeInTheDocument();
    });
  });

  describe('Display Value', () => {
    it('should display "Recent publish" when sortBy is null', () => {
      render(
        <ModelCatalogContextProvider>
          <ModelCatalogSortDropdown performanceViewEnabled />
        </ModelCatalogContextProvider>,
      );

      const toggle = screen.getByTestId('model-catalog-sort-dropdown');
      expect(toggle).toHaveTextContent('Recent publish');
    });

    it('should display "Recent publish" when sortBy is recent_publish', () => {
      // We'll need to set the context value - for now, test the default behavior
      render(
        <ModelCatalogContextProvider>
          <ModelCatalogSortDropdown performanceViewEnabled />
        </ModelCatalogContextProvider>,
      );

      const toggle = screen.getByTestId('model-catalog-sort-dropdown');
      expect(toggle).toHaveTextContent('Recent publish');
    });

    it('should display "Lowest latency" when sortBy is lowest_latency', async () => {
      const TestComponent: React.FC = () => {
        const { setSortBy } = React.useContext(ModelCatalogContext);

        React.useEffect(() => {
          setSortBy(ModelCatalogSortOption.LOWEST_LATENCY);
        }, [setSortBy]);

        return <ModelCatalogSortDropdown performanceViewEnabled />;
      };

      render(
        <ModelCatalogContextProvider>
          <TestComponent />
        </ModelCatalogContextProvider>,
      );

      await waitFor(() => {
        const toggle = screen.getByTestId('model-catalog-sort-dropdown');
        expect(toggle).toHaveTextContent('Lowest latency');
      });
    });
  });

  describe('Dropdown Options', () => {
    it('should show both sort options when dropdown is opened', async () => {
      render(
        <ModelCatalogContextProvider>
          <ModelCatalogSortDropdown performanceViewEnabled />
        </ModelCatalogContextProvider>,
      );

      const toggle = screen.getByTestId('model-catalog-sort-dropdown');
      fireEvent.click(toggle);

      await waitFor(() => {
        expect(screen.getByTestId('sort-option-recent-publish')).toBeInTheDocument();
        expect(screen.getByTestId('sort-option-lowest-latency')).toBeInTheDocument();
      });
    });

    it('should disable "Lowest latency" option when there is no active latency field', async () => {
      mockGetActiveLatencyFieldName.mockReturnValue(undefined);

      render(
        <ModelCatalogContextProvider>
          <ModelCatalogSortDropdown performanceViewEnabled />
        </ModelCatalogContextProvider>,
      );

      const toggle = screen.getByTestId('model-catalog-sort-dropdown');
      fireEvent.click(toggle);

      await waitFor(() => {
        const lowestLatencyOption = screen.getByTestId('sort-option-lowest-latency');
        const button = lowestLatencyOption.querySelector('button');
        expect(button).toBeDisabled();
      });
    });

    it('should enable "Lowest latency" when active latency field exists', async () => {
      mockGetActiveLatencyFieldName.mockReturnValue('artifacts.ttft_p90.double_value');

      render(
        <ModelCatalogContextProvider>
          <ModelCatalogSortDropdown performanceViewEnabled />
        </ModelCatalogContextProvider>,
      );

      const toggle = screen.getByTestId('model-catalog-sort-dropdown');
      fireEvent.click(toggle);

      await waitFor(() => {
        const lowestLatencyOption = screen.getByTestId('sort-option-lowest-latency');
        expect(lowestLatencyOption).not.toHaveAttribute('disabled');
      });
    });
  });

  describe('Selection', () => {
    it('should call setSortBy when "Recent publish" is selected', async () => {
      const TestComponent: React.FC = () => {
        const { sortBy, setSortBy } = React.useContext(ModelCatalogContext);
        const [selectedValue, setSelectedValue] = React.useState<string>('');

        React.useEffect(() => {
          setSortBy(ModelCatalogSortOption.LOWEST_LATENCY);
        }, [setSortBy]);

        React.useEffect(() => {
          if (sortBy) {
            setSelectedValue(sortBy);
          }
        }, [sortBy]);

        return (
          <div>
            <ModelCatalogSortDropdown performanceViewEnabled />
            <div data-testid="selected-value">{selectedValue}</div>
          </div>
        );
      };

      render(
        <ModelCatalogContextProvider>
          <TestComponent />
        </ModelCatalogContextProvider>,
      );

      // Wait for initial state
      await waitFor(() => {
        expect(screen.getByTestId('model-catalog-sort-dropdown')).toHaveTextContent(
          'Lowest latency',
        );
      });

      // Open dropdown
      const toggle = screen.getByTestId('model-catalog-sort-dropdown');
      fireEvent.click(toggle);

      // Wait for dropdown to open and select "Recent publish"
      await waitFor(() => {
        const recentPublishOption = screen.getByTestId('sort-option-recent-publish');
        expect(recentPublishOption).toBeInTheDocument();
      });

      const recentPublishOption = screen.getByTestId('sort-option-recent-publish');
      // Click the button inside the option
      const button = recentPublishOption.querySelector('button');
      if (button) {
        fireEvent.click(button);
      } else {
        fireEvent.click(recentPublishOption);
      }

      // Verify the sort changed
      await waitFor(
        () => {
          const selectedValueElement = screen.getByTestId('selected-value');
          expect(selectedValueElement.textContent).toBe(ModelCatalogSortOption.RECENT_PUBLISH);
        },
        { timeout: 3000 },
      );
    });

    it('should call setSortBy when "Lowest latency" is selected', async () => {
      // Mock an active latency field so the option is enabled
      mockGetActiveLatencyFieldName.mockReturnValue('artifacts.ttft_p90.double_value');

      const TestComponent: React.FC = () => {
        const { sortBy } = React.useContext(ModelCatalogContext);
        const [selectedValue, setSelectedValue] = React.useState<string>('');

        React.useEffect(() => {
          if (sortBy) {
            setSelectedValue(sortBy);
          }
        }, [sortBy]);

        return (
          <div>
            <ModelCatalogSortDropdown performanceViewEnabled />
            <div data-testid="selected-value">{selectedValue}</div>
          </div>
        );
      };

      render(
        <ModelCatalogContextProvider>
          <TestComponent />
        </ModelCatalogContextProvider>,
      );

      // Open dropdown
      const toggle = screen.getByTestId('model-catalog-sort-dropdown');
      fireEvent.click(toggle);

      // Wait for dropdown to open and select "Lowest latency"
      await waitFor(() => {
        const lowestLatencyOption = screen.getByTestId('sort-option-lowest-latency');
        expect(lowestLatencyOption).toBeInTheDocument();
      });

      const lowestLatencyOption = screen.getByTestId('sort-option-lowest-latency');
      // Click the button inside the option
      const button = lowestLatencyOption.querySelector('button');
      if (button) {
        fireEvent.click(button);
      } else {
        fireEvent.click(lowestLatencyOption);
      }

      // Verify the sort changed
      await waitFor(
        () => {
          const selectedValueElement = screen.getByTestId('selected-value');
          expect(selectedValueElement.textContent).toBe(ModelCatalogSortOption.LOWEST_LATENCY);
        },
        { timeout: 3000 },
      );
    });
  });

  describe('Layout', () => {
    it('should render "Sort:" label next to the dropdown', () => {
      render(
        <ModelCatalogContextProvider>
          <ModelCatalogSortDropdown performanceViewEnabled />
        </ModelCatalogContextProvider>,
      );

      expect(screen.getByText('Sort:')).toBeInTheDocument();
      expect(screen.getByTestId('model-catalog-sort-dropdown')).toBeInTheDocument();
    });
  });
});

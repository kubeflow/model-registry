import '@testing-library/jest-dom';
import React from 'react';
import { renderHook, act } from '@testing-library/react';
import {
  ModelCatalogContextProvider,
  ModelCatalogContext,
} from '~/app/context/modelCatalog/ModelCatalogContext';
import { ModelCatalogSortOption } from '~/concepts/modelCatalog/const';

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

describe('ModelCatalogContext - Sorting Behavior', () => {
  const wrapper = ({ children }: { children: React.ReactNode }) => (
    <ModelCatalogContextProvider>{children}</ModelCatalogContextProvider>
  );

  describe('setPerformanceViewEnabled - Default Sorting', () => {
    it('should set sortBy to recent_publish when toggle is turned off', () => {
      const { result } = renderHook(() => React.useContext(ModelCatalogContext), { wrapper });

      // First, turn on the toggle (which sets default to lowest_latency)
      act(() => {
        result.current.setPerformanceViewEnabled(true);
      });

      // Verify it's set to lowest_latency
      expect(result.current.sortBy).toBe(ModelCatalogSortOption.LOWEST_LATENCY);

      // Now turn off the toggle
      act(() => {
        result.current.setPerformanceViewEnabled(false);
      });

      // Should default to recent_publish
      expect(result.current.sortBy).toBe(ModelCatalogSortOption.RECENT_PUBLISH);
    });

    it('should set sortBy to lowest_latency when toggle is turned on', () => {
      const { result } = renderHook(() => React.useContext(ModelCatalogContext), { wrapper });

      // Initially sortBy should be null
      expect(result.current.sortBy).toBeNull();

      // Turn on the toggle
      act(() => {
        result.current.setPerformanceViewEnabled(true);
      });

      // Should default to lowest_latency
      expect(result.current.sortBy).toBe(ModelCatalogSortOption.LOWEST_LATENCY);
    });

    it('should preserve existing sortBy when toggle is turned off if sortBy is already recent_publish', () => {
      const { result } = renderHook(() => React.useContext(ModelCatalogContext), { wrapper });

      // Set sortBy to recent_publish manually
      act(() => {
        result.current.setSortBy(ModelCatalogSortOption.RECENT_PUBLISH);
      });

      expect(result.current.sortBy).toBe(ModelCatalogSortOption.RECENT_PUBLISH);

      // Turn on toggle (should change to lowest_latency)
      act(() => {
        result.current.setPerformanceViewEnabled(true);
      });

      expect(result.current.sortBy).toBe(ModelCatalogSortOption.LOWEST_LATENCY);

      // Turn off toggle (should change back to recent_publish)
      act(() => {
        result.current.setPerformanceViewEnabled(false);
      });

      expect(result.current.sortBy).toBe(ModelCatalogSortOption.RECENT_PUBLISH);
    });

    it('should preserve existing sortBy when toggle is turned on if sortBy is already lowest_latency', () => {
      const { result } = renderHook(() => React.useContext(ModelCatalogContext), { wrapper });

      // Set sortBy to lowest_latency manually
      act(() => {
        result.current.setSortBy(ModelCatalogSortOption.LOWEST_LATENCY);
      });

      expect(result.current.sortBy).toBe(ModelCatalogSortOption.LOWEST_LATENCY);

      // Turn on toggle (should preserve lowest_latency)
      act(() => {
        result.current.setPerformanceViewEnabled(true);
      });

      expect(result.current.sortBy).toBe(ModelCatalogSortOption.LOWEST_LATENCY);
    });

    it('should change from lowest_latency to recent_publish when toggle is turned off', () => {
      const { result } = renderHook(() => React.useContext(ModelCatalogContext), { wrapper });

      // Turn on toggle first
      act(() => {
        result.current.setPerformanceViewEnabled(true);
      });

      expect(result.current.sortBy).toBe(ModelCatalogSortOption.LOWEST_LATENCY);

      // Turn off toggle
      act(() => {
        result.current.setPerformanceViewEnabled(false);
      });

      expect(result.current.sortBy).toBe(ModelCatalogSortOption.RECENT_PUBLISH);
    });

    it('should change from recent_publish to lowest_latency when toggle is turned on', () => {
      const { result } = renderHook(() => React.useContext(ModelCatalogContext), { wrapper });

      // Set to recent_publish first
      act(() => {
        result.current.setSortBy(ModelCatalogSortOption.RECENT_PUBLISH);
      });

      expect(result.current.sortBy).toBe(ModelCatalogSortOption.RECENT_PUBLISH);

      // Turn on toggle
      act(() => {
        result.current.setPerformanceViewEnabled(true);
      });

      expect(result.current.sortBy).toBe(ModelCatalogSortOption.LOWEST_LATENCY);
    });
  });

  describe('setSortBy', () => {
    it('should update sortBy when setSortBy is called', () => {
      const { result } = renderHook(() => React.useContext(ModelCatalogContext), { wrapper });

      expect(result.current.sortBy).toBeNull();

      act(() => {
        result.current.setSortBy(ModelCatalogSortOption.RECENT_PUBLISH);
      });

      expect(result.current.sortBy).toBe(ModelCatalogSortOption.RECENT_PUBLISH);

      act(() => {
        result.current.setSortBy(ModelCatalogSortOption.LOWEST_LATENCY);
      });

      expect(result.current.sortBy).toBe(ModelCatalogSortOption.LOWEST_LATENCY);
    });

    it('should allow setting sortBy to null', () => {
      const { result } = renderHook(() => React.useContext(ModelCatalogContext), { wrapper });

      act(() => {
        result.current.setSortBy(ModelCatalogSortOption.RECENT_PUBLISH);
      });

      expect(result.current.sortBy).toBe(ModelCatalogSortOption.RECENT_PUBLISH);

      act(() => {
        result.current.setSortBy(null);
      });

      expect(result.current.sortBy).toBeNull();
    });
  });
});

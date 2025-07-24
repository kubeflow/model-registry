import * as React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { ModelCatalogItem, ModelCatalogSource } from '~/app/modelCatalogTypes';
import ModelCatalogPage from '~/app/pages/modelCatalog/screens/ModelCatalogPage';

const mockRefreshSources = jest.fn();

const mockSources: ModelCatalogSource[] = [
  {
    name: 'test-source',
    displayName: 'Test Source',
    models: [
      {
        id: 'test-model',
        name: 'test-model',
        displayName: 'Test Model',
        description: 'Test model description',
        provider: 'Test Provider',
        url: 'https://test.com/model',
        tags: ['test', 'mock'],
        framework: 'Test Framework',
        task: 'test-task',
        license: 'MIT',
        metrics: {
          accuracy: 0.95,
        },
      },
    ],
  },
];

const renderWithContext = (
  ui: React.ReactElement,
  { loading = false, error = undefined, sources = mockSources }: RenderContextOptions = {},
) =>
  render(
    <ModelCatalogContext.Provider
      value={{ sources, loading, error, refreshSources: mockRefreshSources }}
    >
      {ui}
    </ModelCatalogContext.Provider>,
  );

type RenderContextOptions = {
  loading?: boolean;
  error?: Error | undefined;
  sources?: ModelCatalogSource[];
};

describe('ModelCatalogPage', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders loading state', () => {
    renderWithContext(<ModelCatalogPage />, { loading: true });
    expect(screen.getByText('Loading model catalog...')).toBeInTheDocument();
  });

  it('renders error state', () => {
    const error = new Error('Test error');
    renderWithContext(<ModelCatalogPage />, { error });
    expect(screen.getByText('Failed to load model catalog')).toBeInTheDocument();
    expect(screen.getByText(error.message)).toBeInTheDocument();
  });

  it('renders empty state', () => {
    renderWithContext(<ModelCatalogPage />, { sources: [] });
    expect(screen.getByText('No models available')).toBeInTheDocument();
  });

  it('renders model cards', async () => {
    renderWithContext(<ModelCatalogPage />);

    await waitFor(() => {
      mockSources.forEach((source: ModelCatalogSource) => {
        source.models?.forEach((model: ModelCatalogItem) => {
          expect(screen.getByText(model.displayName)).toBeInTheDocument();
        });
      });
    });
  });
});

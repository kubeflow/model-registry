import * as React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import { MemoryRouter } from 'react-router-dom';
import { useModelCatalogSources } from '~/app/hooks/modelCatalog/useModelCatalogSources';
import ModelCatalogPage from '~/app/pages/modelCatalog/screens/ModelCatalogPage';

jest.mock('~/app/hooks/modelCatalog/useModelCatalogSources');

const mockUseModelCatalogSources = useModelCatalogSources as jest.Mock;

const mockSources = [
  {
    name: 'test-source',
    displayName: 'Test Source',
    models: [
      {
        id: 'test-model',
        name: 'test-model',
        displayName: 'Test Model',
        description: 'Test Description',
        tags: ['1.0.0'],
      },
    ],
  },
];

const renderWithRouter = (ui: React.ReactElement) => render(ui, { wrapper: MemoryRouter });

describe('ModelCatalogPage', () => {
  beforeEach(() => {
    mockUseModelCatalogSources.mockReturnValue({
      sources: mockSources,
      loading: false,
      error: null,
    });
  });

  it('renders model cards', () => {
    renderWithRouter(<ModelCatalogPage />);

    expect(screen.getByText('Model Catalog')).toBeInTheDocument();
    expect(screen.getByText('test-model')).toBeInTheDocument();
    expect(screen.getByText('Test Description')).toBeInTheDocument();
    expect(screen.getByText('Test Source')).toBeInTheDocument();
  });

  it('shows loading state', () => {
    mockUseModelCatalogSources.mockReturnValue({
      sources: [],
      loading: true,
      error: null,
    });

    renderWithRouter(<ModelCatalogPage />);
    expect(screen.getByText('Loading model catalog...')).toBeInTheDocument();
  });

  it('shows error state', () => {
    mockUseModelCatalogSources.mockReturnValue({
      sources: [],
      loading: false,
      error: new Error('Test error'),
    });

    renderWithRouter(<ModelCatalogPage />);
    expect(screen.getByText('Failed to load model catalog')).toBeInTheDocument();
  });

  it('shows empty state', () => {
    mockUseModelCatalogSources.mockReturnValue({
      sources: [],
      loading: false,
      error: null,
    });

    renderWithRouter(<ModelCatalogPage />);
    expect(screen.getByText('No models available')).toBeInTheDocument();
  });
});

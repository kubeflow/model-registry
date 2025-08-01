import * as React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import { MemoryRouter } from 'react-router-dom';
import { ModelCatalogItem } from '~/app/modelCatalogTypes';
import ModelCatalogCard from '~/app/pages/modelCatalog/components/ModelCatalogCard';

const mockModel: ModelCatalogItem = {
  id: 'test-model',
  name: 'test-model',
  displayName: 'Test Model',
  description: 'Test model description',
  provider: 'Test Provider',
  url: 'https://test.com/model',
  tags: ['1.0.0', 'feature1', 'feature2'],
  framework: 'Test Framework',
  task: 'test-task',
  license: 'MIT',
};

const renderWithRouter = (ui: React.ReactElement) => render(ui, { wrapper: MemoryRouter });

describe('ModelCatalogCard', () => {
  it('renders model information correctly', () => {
    renderWithRouter(<ModelCatalogCard model={mockModel} source="Test Source" />);

    expect(screen.getByText('test-model')).toBeInTheDocument();
    expect(screen.getByText('Test model description')).toBeInTheDocument();
    expect(screen.getByText('Test Source')).toBeInTheDocument();
    expect(screen.getByText('Test Framework')).toBeInTheDocument();
    expect(screen.getByText('test-task')).toBeInTheDocument();
    expect(screen.getByText('MIT')).toBeInTheDocument();
    expect(screen.getByText('1.0.0')).toBeInTheDocument();
    expect(screen.getByText('feature1')).toBeInTheDocument();
    expect(screen.getByText('feature2')).toBeInTheDocument();
  });

  it('calls onSelect when select button is clicked', () => {
    const onSelect = jest.fn();
    renderWithRouter(
      <ModelCatalogCard model={mockModel} source="Test Source" onSelect={onSelect} />,
    );

    fireEvent.click(screen.getByTestId('model-catalog-detail-link'));
    expect(onSelect).toHaveBeenCalledWith(mockModel);
  });

  it('navigates when no onSelect is provided', () => {
    renderWithRouter(<ModelCatalogCard model={mockModel} source="Test Source" />);

    const link = screen.getByTestId('model-catalog-detail-link');
    expect(link).toBeInTheDocument();
  });

  it('shows version from tags', () => {
    renderWithRouter(<ModelCatalogCard model={mockModel} source="Test Source" />);
    expect(screen.getByText('1.0.0')).toBeInTheDocument();
  });

  it('shows "No version" when no version tag is present', () => {
    const modelWithoutVersion = { ...mockModel, tags: ['feature1', 'feature2'] };
    renderWithRouter(<ModelCatalogCard model={modelWithoutVersion} source="Test Source" />);
    expect(screen.getByText('No version')).toBeInTheDocument();
  });
});

import * as React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import { ModelCatalogItem } from '~/app/modelCatalogTypes';
import ModelCatalogCard from '~/app/pages/modelCatalog/components/ModelCatalogCard';

const mockModel: ModelCatalogItem = {
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
};

describe('ModelCatalogCard', () => {
  it('renders model information correctly', () => {
    render(<ModelCatalogCard model={mockModel} />);

    expect(screen.getByText(mockModel.displayName)).toBeInTheDocument();
    expect(screen.getByText(mockModel.description!)).toBeInTheDocument();
    expect(screen.getByText(mockModel.provider!)).toBeInTheDocument();
    expect(screen.getByText(mockModel.framework!)).toBeInTheDocument();
    expect(screen.getByText(mockModel.task!)).toBeInTheDocument();
    expect(screen.getByText(mockModel.license!)).toBeInTheDocument();
  });

  it('calls onSelect when select button is clicked', () => {
    const onSelect = jest.fn();
    render(<ModelCatalogCard model={mockModel} onSelect={onSelect} />);

    fireEvent.click(screen.getByTestId('select-model-button'));
    expect(onSelect).toHaveBeenCalledWith(mockModel);
  });

  it('renders view model link when url is provided', () => {
    render(<ModelCatalogCard model={mockModel} />);

    const link = screen.getByTestId('view-model-link');
    expect(link).toHaveAttribute('href', mockModel.url);
  });

  it('renders metrics when provided', () => {
    render(<ModelCatalogCard model={mockModel} />);

    Object.entries(mockModel.metrics!).forEach(([key, value]) => {
      expect(screen.getByText(new RegExp(`${key}.*${value}`))).toBeInTheDocument();
    });
  });
});

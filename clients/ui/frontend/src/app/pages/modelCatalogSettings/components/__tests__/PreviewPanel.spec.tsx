import React from 'react';
import { screen, render } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';
import '@testing-library/jest-dom';

import { CatalogSourcePreviewModel, CatalogSourcePreviewSummary } from '~/app/modelCatalogTypes';
import PreviewPanel from '~/app/pages/modelCatalogSettings/components/PreviewPanel';

const mockSummary: CatalogSourcePreviewSummary = {
  totalModels: 20,
  includedModels: 15,
  excludedModels: 5,
};

const mockIncludedItems: CatalogSourcePreviewModel[] = [
  { name: 'model-1', included: true },
  { name: 'model-2', included: true },
  { name: 'model-3', included: true },
];

const mockExcludedItems: CatalogSourcePreviewModel[] = [
  { name: 'excluded-model-1', included: false },
  { name: 'excluded-model-2', included: false },
];

const defaultProps = {
  isPreviewEnabled: true,
  isLoading: false,
  onPreview: jest.fn(),
  previewError: undefined,
  hasFormChanged: false,
  activeTab: 'included' as const,
  items: mockIncludedItems,
  summary: mockSummary,
  hasMore: false,
  isLoadingMore: false,
  onTabChange: jest.fn(),
  onLoadMore: jest.fn(),
};

describe('PreviewPanel', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders empty state when no items and no summary', () => {
    render(<PreviewPanel {...defaultProps} items={[]} summary={undefined} />);

    expect(screen.getByText('Preview models')).toBeInTheDocument();
    expect(
      screen.getByText(/To view the models from this source that will appear in the model catalog/),
    ).toBeInTheDocument();
  });

  it('renders loading spinner when isLoading is true', () => {
    render(<PreviewPanel {...defaultProps} isLoading items={[]} />);

    expect(screen.getByLabelText('Loading preview')).toBeInTheDocument();
  });

  it('renders error state with retry button when previewError exists', () => {
    const error = new Error('Failed to fetch preview');
    render(<PreviewPanel {...defaultProps} previewError={error} items={[]} />);

    expect(screen.getByText('Failed to preview the results')).toBeInTheDocument();
    expect(screen.getByText('Failed to fetch preview')).toBeInTheDocument();
    expect(screen.getByTestId('preview-button-panel-retry')).toBeInTheDocument();
  });

  it('renders tabs for included/excluded models', () => {
    render(<PreviewPanel {...defaultProps} />);

    expect(screen.getByText('Models included')).toBeInTheDocument();
    expect(screen.getByText('Models excluded')).toBeInTheDocument();
  });

  it('calls onTabChange when switching tabs', async () => {
    const user = userEvent.setup();
    const onTabChange = jest.fn();

    render(<PreviewPanel {...defaultProps} onTabChange={onTabChange} />);

    await user.click(screen.getByText('Models excluded'));

    expect(onTabChange).toHaveBeenCalledWith('excluded');
  });

  it('displays correct count text for included tab', () => {
    render(<PreviewPanel {...defaultProps} activeTab="included" />);

    expect(screen.getByText('15 of 20 models included:')).toBeInTheDocument();
  });

  it('displays correct count text for excluded tab', () => {
    render(<PreviewPanel {...defaultProps} activeTab="excluded" items={mockExcludedItems} />);

    expect(screen.getByText('5 of 20 models excluded:')).toBeInTheDocument();
  });

  it('renders Load more button when hasMore is true', () => {
    render(<PreviewPanel {...defaultProps} hasMore />);

    expect(screen.getByText('Load more')).toBeInTheDocument();
  });

  it('does not render Load more button when hasMore is false', () => {
    render(<PreviewPanel {...defaultProps} hasMore={false} />);

    expect(screen.queryByText('Load more')).not.toBeInTheDocument();
  });

  it('calls onLoadMore when Load more button clicked', async () => {
    const user = userEvent.setup();
    const onLoadMore = jest.fn();

    render(<PreviewPanel {...defaultProps} hasMore onLoadMore={onLoadMore} />);

    await user.click(screen.getByText('Load more'));

    expect(onLoadMore).toHaveBeenCalled();
  });

  it('shows loading state on Load more button when isLoadingMore is true', () => {
    render(<PreviewPanel {...defaultProps} hasMore isLoadingMore />);

    expect(screen.getByText('Loading...')).toBeInTheDocument();
  });

  it('shows refresh alert when hasFormChanged is true', () => {
    render(<PreviewPanel {...defaultProps} hasFormChanged />);

    expect(
      screen.getByText('The preview needs to be refreshed after any changes are made'),
    ).toBeInTheDocument();
    expect(screen.getByTestId('refresh-preview-link')).toBeInTheDocument();
  });

  it('calls onPreview when refresh link clicked', async () => {
    const user = userEvent.setup();
    const onPreview = jest.fn();

    render(<PreviewPanel {...defaultProps} hasFormChanged onPreview={onPreview} />);

    await user.click(screen.getByTestId('refresh-preview-link'));

    expect(onPreview).toHaveBeenCalled();
  });

  it('renders empty state for included tab with no items but with summary', () => {
    render(
      <PreviewPanel
        {...defaultProps}
        items={[]}
        summary={{ ...mockSummary, includedModels: 0 }}
        activeTab="included"
      />,
    );

    expect(screen.getByText('No models included')).toBeInTheDocument();
    expect(screen.getByText('No models from this source match this filter')).toBeInTheDocument();
  });

  it('renders empty state for excluded tab with no items', () => {
    render(
      <PreviewPanel
        {...defaultProps}
        items={[]}
        summary={{ ...mockSummary, excludedModels: 0 }}
        activeTab="excluded"
      />,
    );

    expect(screen.getByText('No models excluded')).toBeInTheDocument();
    expect(
      screen.getByText('No models from this source are excluded by this filter'),
    ).toBeInTheDocument();
  });

  it('renders model items in list', () => {
    render(<PreviewPanel {...defaultProps} />);

    expect(screen.getByText('model-1')).toBeInTheDocument();
    expect(screen.getByText('model-2')).toBeInTheDocument();
    expect(screen.getByText('model-3')).toBeInTheDocument();
  });

  it('disables preview button when isPreviewEnabled is false', () => {
    render(
      <PreviewPanel {...defaultProps} isPreviewEnabled={false} items={[]} summary={undefined} />,
    );

    const previewButton = screen.getByTestId('preview-button-panel');
    expect(previewButton).toBeDisabled();
  });
});

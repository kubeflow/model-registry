import * as React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import ModelDetailsView from '~/app/pages/modelCatalog/screens/ModelDetailsView';
import { CatalogArtifactList, CatalogModel } from '~/app/modelCatalogTypes';

jest.mock('mod-arch-shared', () => ({
  InlineTruncatedClipboardCopy: ({
    textToCopy,
    testId,
  }: {
    textToCopy: string;
    testId: string;
  }) => <span data-testid={testId}>{textToCopy}</span>,
  relativeTime: () => '2 days ago',
}));

jest.mock('~/app/shared/markdown/MarkdownComponent', () => ({
  __esModule: true,
  default: ({ data, dataTestId }: { data: string; dataTestId: string }) => (
    <div data-testid={dataTestId}>{data}</div>
  ),
}));

jest.mock('~/app/pages/modelCatalog/components/ModelCatalogLabels', () => ({
  __esModule: true,
  default: () => <span data-testid="model-catalog-labels">Labels</span>,
}));

jest.mock('~/app/pages/modelRegistry/screens/components/ModelTimestamp', () => ({
  __esModule: true,
  default: ({ timeSinceEpoch }: { timeSinceEpoch?: string }) => (
    <span data-testid="model-timestamp">{timeSinceEpoch || '--'}</span>
  ),
}));

const mockModel: CatalogModel = {
  name: 'Test Model',
  description: 'A test model description',
  provider: 'Test Provider',
  readme: '# Test README',
  licenseLink: 'https://example.com/license',
  license: 'Apache-2.0',
  tasks: ['text-generation'],
  createTimeSinceEpoch: '1700000000000',
  lastUpdateTimeSinceEpoch: '1700100000000',
  customProperties: {},
};

const mockArtifacts: CatalogArtifactList = {
  items: [],
  size: 0,
  pageSize: 0,
  nextPageToken: '',
};

describe('ModelDetailsView', () => {
  it('renders model description', () => {
    render(
      <ModelDetailsView
        model={mockModel}
        artifacts={mockArtifacts}
        artifactLoaded
        artifactsLoadError={undefined}
      />,
    );
    expect(screen.getByTestId('model-long-description')).toHaveTextContent(
      'A test model description',
    );
  });

  it('renders license link inside a DescriptionListDescription (a11y compliance)', () => {
    const { container } = render(
      <ModelDetailsView
        model={mockModel}
        artifacts={mockArtifacts}
        artifactLoaded
        artifactsLoadError={undefined}
      />,
    );

    const licenseLink = screen.getByTestId('model-license-link');
    expect(licenseLink).toBeInTheDocument();

    // Verify the ExternalLink is wrapped in a DescriptionListDescription (rendered as <dd>)
    const ddElement = licenseLink.closest('dd');
    expect(ddElement).not.toBeNull();

    // Verify the dd is inside a DescriptionListGroup (rendered as a div within dl)
    const dlElement = container.querySelector('.pf-v6-c-description-list');
    expect(dlElement).toBeInTheDocument();
  });

  it('renders provider information', () => {
    render(
      <ModelDetailsView
        model={mockModel}
        artifacts={mockArtifacts}
        artifactLoaded
        artifactsLoadError={undefined}
      />,
    );
    expect(screen.getByText('Test Provider')).toBeInTheDocument();
  });

  it('renders "No description" when description is empty', () => {
    render(
      <ModelDetailsView
        model={{ ...mockModel, description: undefined }}
        artifacts={mockArtifacts}
        artifactLoaded
        artifactsLoadError={undefined}
      />,
    );
    expect(screen.getByTestId('model-long-description')).toHaveTextContent('No description');
  });

  it('renders model card markdown when readme exists', () => {
    render(
      <ModelDetailsView
        model={mockModel}
        artifacts={mockArtifacts}
        artifactLoaded
        artifactsLoadError={undefined}
      />,
    );
    expect(screen.getByTestId('model-card-markdown')).toBeInTheDocument();
  });

  it('renders "No model card" when readme is absent', () => {
    render(
      <ModelDetailsView
        model={{ ...mockModel, readme: undefined }}
        artifacts={mockArtifacts}
        artifactLoaded
        artifactsLoadError={undefined}
      />,
    );
    expect(screen.getByText('No model card')).toBeInTheDocument();
  });

  it('renders N/A for provider when provider is not set', () => {
    render(
      <ModelDetailsView
        model={{ ...mockModel, provider: undefined }}
        artifacts={mockArtifacts}
        artifactLoaded
        artifactsLoadError={undefined}
      />,
    );
    const providerValues = screen.getAllByText('N/A');
    expect(providerValues.length).toBeGreaterThan(0);
  });

  it('renders artifact load error when present', () => {
    const error = new Error('Failed to load artifacts');
    error.name = 'LoadError';
    render(
      <ModelDetailsView
        model={mockModel}
        artifacts={mockArtifacts}
        artifactLoaded={false}
        artifactsLoadError={error}
      />,
    );
    expect(screen.getByText('Failed to load artifacts')).toBeInTheDocument();
  });

  it('renders spinner when artifacts are loading', () => {
    render(
      <ModelDetailsView
        model={mockModel}
        artifacts={mockArtifacts}
        artifactLoaded={false}
        artifactsLoadError={undefined}
      />,
    );
    expect(screen.getByRole('progressbar')).toBeInTheDocument();
  });
});

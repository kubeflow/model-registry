import '@testing-library/jest-dom';
import * as React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import McpCatalogStringFilter from '~/app/pages/mcpCatalog/components/McpCatalogStringFilter';
import { McpCatalogContextProvider } from '~/app/context/mcpCatalog/McpCatalogContext';

jest.mock('mod-arch-core', () => ({ useQueryParamNamespaces: () => ({}) }));
jest.mock('~/app/utilities/const', () => ({
  BFF_API_VERSION: 'v1',
  URL_PREFIX: '/model-registry',
}));
jest.mock('~/app/hooks/modelCatalog/useModelCatalogAPIState', () => ({
  __esModule: true,
  default: () => [{ apiAvailable: false, api: {} }, jest.fn()],
}));
jest.mock('~/app/hooks/modelCatalog/useCatalogSources', () => ({
  useCatalogSources: () => [{ items: [] }, true, undefined],
}));
jest.mock('~/app/hooks/modelCatalog/useCatalogLabels', () => ({
  useCatalogLabels: () => [null, true, undefined],
}));
jest.mock('~/app/hooks/mcpServerCatalog/useMcpServerFilterOptionList', () => ({
  useMcpServerFilterOptionListWithAPI: () => [null, true, undefined],
}));

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <MemoryRouter>
    <McpCatalogContextProvider>{children}</McpCatalogContextProvider>
  </MemoryRouter>
);

describe('McpCatalogStringFilter', () => {
  it('renders title and filter options', () => {
    render(
      <McpCatalogStringFilter
        title="Deployment mode"
        filterKey="deploymentMode"
        filters={{ type: 'string', values: ['Local', 'Remote'] }}
      />,
      { wrapper },
    );
    expect(screen.getByText('Deployment mode')).toBeInTheDocument();
    expect(screen.getByLabelText('Local')).toBeInTheDocument();
    expect(screen.getByLabelText('Remote')).toBeInTheDocument();
  });

  it('shows empty state when no values match search', () => {
    render(
      <McpCatalogStringFilter
        title="License"
        filterKey="license"
        filters={{
          type: 'string',
          values: ['MIT', 'Apache-2.0', 'GPL', 'BSD', 'LGPL', 'AGPL'],
        }}
      />,
      { wrapper },
    );
    const searchWrapper = screen.getByTestId('mcp-filter-license-search');
    const input = searchWrapper.querySelector('input');
    expect(input).toBeTruthy();
    if (input) {
      fireEvent.change(input, { target: { value: 'nonexistent' } });
    }
    expect(screen.getByTestId('mcp-filter-license-empty')).toHaveTextContent('No results found');
  });

  it('toggles checkbox and updates context', () => {
    render(
      <McpCatalogStringFilter
        title="Deployment mode"
        filterKey="deploymentMode"
        filters={{ type: 'string', values: ['Local', 'Remote'] }}
      />,
      { wrapper },
    );
    const localCheckbox = screen.getByTestId('mcp-filter-deploymentMode-Local');
    expect(localCheckbox).not.toBeChecked();
    fireEvent.click(localCheckbox);
    expect(localCheckbox).toBeChecked();
  });

  it('shows Show more when values exceed max visible', () => {
    const values = ['a', 'b', 'c', 'd', 'e', 'f'];
    render(
      <McpCatalogStringFilter
        title="Labels"
        filterKey="labels"
        filters={{ type: 'string', values }}
      />,
      { wrapper },
    );
    expect(screen.getByTestId('mcp-filter-labels-show-more')).toBeInTheDocument();
    fireEvent.click(screen.getByTestId('mcp-filter-labels-show-more'));
    expect(screen.getByTestId('mcp-filter-labels-show-less')).toBeInTheDocument();
  });

  it('renders with data-testid for filter key', () => {
    render(
      <McpCatalogStringFilter
        title="Deployment mode"
        filterKey="deploymentMode"
        filters={{ type: 'string', values: ['Local'] }}
      />,
      { wrapper },
    );
    expect(screen.getByTestId('mcp-filter-deploymentMode')).toBeInTheDocument();
  });
});

import '@testing-library/jest-dom';
import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import McpCatalogStringFilter from '~/app/pages/mcpCatalog/components/McpCatalogStringFilter';
import { McpCatalogContextProvider } from '~/app/context/mcpCatalog/McpCatalogContext';

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <McpCatalogContextProvider>{children}</McpCatalogContextProvider>
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
        filters={{ type: 'string', values: ['MIT', 'Apache-2.0'] }}
        showSearch
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

  it('shows Show more when values exceed MCP_FILTER_MAX_VISIBLE', () => {
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

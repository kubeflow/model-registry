import '@testing-library/jest-dom';
import React from 'react';
import { render, screen } from '@testing-library/react';
import McpCatalogFilters from '~/app/pages/mcpCatalog/components/McpCatalogFilters';
import { McpCatalogContextProvider } from '~/app/context/mcpCatalog/McpCatalogContext';

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <McpCatalogContextProvider>{children}</McpCatalogContextProvider>
);

describe('McpCatalogFilters', () => {
  it('renders all filter sections from mock options', () => {
    render(<McpCatalogFilters />, { wrapper });
    expect(screen.getByTestId('mcp-filter-deploymentMode')).toBeInTheDocument();
    expect(screen.getByTestId('mcp-filter-supportedTransports')).toBeInTheDocument();
    expect(screen.getByTestId('mcp-filter-license')).toBeInTheDocument();
    expect(screen.getByTestId('mcp-filter-labels')).toBeInTheDocument();
    expect(screen.getByTestId('mcp-filter-securityVerification')).toBeInTheDocument();
  });

  it('renders Deployment mode filter with Local and Remote options', () => {
    render(<McpCatalogFilters />, { wrapper });
    expect(screen.getByLabelText('Local')).toBeInTheDocument();
    expect(screen.getByLabelText('Remote')).toBeInTheDocument();
  });
});

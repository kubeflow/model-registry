import '@testing-library/jest-dom';
import * as React from 'react';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { testFilterOptions as mockMcpCatalogFilterOptions } from '~/__mocks__';
import McpCatalogFilters from '~/app/pages/mcpCatalog/components/McpCatalogFilters';
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
jest.mock('~/app/hooks/mcpServerCatalog/useMcpServersBySourceLabel', () => ({
  useMcpServersBySourceLabelWithAPI: () => ({
    mcpServers: { items: [] },
    mcpServersLoaded: true,
    mcpServersLoadError: undefined,
  }),
}));
jest.mock('~/app/hooks/mcpServerCatalog/useMcpServerFilterOptionList', () => ({
  useMcpServerFilterOptionListWithAPI: () => [mockMcpCatalogFilterOptions, true, undefined],
}));

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <MemoryRouter>
    <McpCatalogContextProvider>{children}</McpCatalogContextProvider>
  </MemoryRouter>
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

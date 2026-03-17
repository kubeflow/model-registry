import '@testing-library/jest-dom';
import * as React from 'react';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import McpCatalogCard from '~/app/pages/mcpCatalog/components/McpCatalogCard';
import type { McpServer } from '~/app/mcpServerCatalogTypes';

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <MemoryRouter>{children}</MemoryRouter>
);

const mockServer: McpServer = {
  id: '1',
  name: 'Test MCP Server',
  description: 'Test description for the server.',
  deploymentMode: 'local',
  securityIndicators: { verifiedSource: true, sast: true },
  toolCount: 0,
};

describe('McpCatalogCard', () => {
  it('renders server name and description', () => {
    render(<McpCatalogCard server={mockServer} />, { wrapper });
    expect(screen.getByTestId('mcp-catalog-card-name-1')).toHaveTextContent('Test MCP Server');
    expect(screen.getByTestId('mcp-catalog-card-description-1')).toHaveTextContent(
      'Test description for the server.',
    );
  });

  it('does not render deployment chip for local mode', () => {
    render(<McpCatalogCard server={mockServer} />, { wrapper });
    expect(screen.queryByTestId('mcp-catalog-card-deployment-1')).not.toBeInTheDocument();
  });

  it('renders deployment chip only for remote mode', () => {
    render(<McpCatalogCard server={{ ...mockServer, id: '2', deploymentMode: 'remote' }} />, {
      wrapper,
    });
    expect(screen.getByTestId('mcp-catalog-card-deployment-2')).toHaveTextContent('Remote');
  });

  it('renders security verification tags', () => {
    render(<McpCatalogCard server={mockServer} />, { wrapper });
    expect(screen.getByText('Verified source')).toBeInTheDocument();
    expect(screen.getByText('SAST')).toBeInTheDocument();
  });

  it('does not render security section when securityIndicators is empty', () => {
    render(<McpCatalogCard server={{ ...mockServer, id: '3', securityIndicators: undefined }} />, {
      wrapper,
    });
    expect(screen.queryByText('Verified source')).not.toBeInTheDocument();
  });

  it('renders card with data-testid for the server id', () => {
    render(<McpCatalogCard server={mockServer} />, { wrapper });
    expect(screen.getByTestId('mcp-catalog-card-1')).toBeInTheDocument();
  });

  it('renders clickable card name as link to details page', () => {
    render(<McpCatalogCard server={mockServer} />, { wrapper });
    const link = screen.getByTestId('mcp-catalog-card-detail-link-1');
    expect(link).toBeInTheDocument();
    expect(link.tagName).toBe('A');
    expect(link).toHaveAttribute('href', '/mcp-catalog/1');
  });

  it('renders description with TruncatedText wrapper', () => {
    render(<McpCatalogCard server={mockServer} />, { wrapper });
    const desc = screen.getByTestId('mcp-catalog-card-description-1');
    expect(desc).toBeInTheDocument();
    expect(desc.style.display).toBe('-webkit-box');
  });

  it('renders empty string when description is undefined', () => {
    render(<McpCatalogCard server={{ ...mockServer, id: '4', description: undefined }} />, {
      wrapper,
    });
    const desc = screen.getByTestId('mcp-catalog-card-description-4');
    expect(desc).toBeInTheDocument();
  });
});

import '@testing-library/jest-dom';
import React from 'react';
import { render, screen } from '@testing-library/react';
import McpCatalogCard from '~/app/pages/mcpCatalog/components/McpCatalogCard';
import type { McpServer } from '~/app/mcpServerCatalogTypes';

jest.mock('react-router-dom', () => ({
  Link: ({ to, children }: { to: string; children: React.ReactNode }) => (
    <a href={to}>{children}</a>
  ),
}));

const mockServer: McpServer = {
  id: 1,
  name: 'Test MCP Server',
  description: 'Test description for the server.',
  deploymentMode: 'local',
  securityIndicators: { verifiedSource: true, sast: true },
  toolCount: 0,
};

describe('McpCatalogCard', () => {
  it('renders server name and description', () => {
    render(<McpCatalogCard server={mockServer} />);
    expect(screen.getByTestId('mcp-catalog-card-name-1')).toHaveTextContent('Test MCP Server');
    expect(screen.getByTestId('mcp-catalog-card-description-1')).toHaveTextContent(
      'Test description for the server.',
    );
  });

  it('renders deployment mode label for local', () => {
    render(<McpCatalogCard server={mockServer} />);
    expect(screen.getByTestId('mcp-catalog-card-deployment-1')).toHaveTextContent(
      'Local to cluster',
    );
  });

  it('renders deployment mode label for remote', () => {
    render(<McpCatalogCard server={{ ...mockServer, id: 2, deploymentMode: 'remote' }} />);
    expect(screen.getByTestId('mcp-catalog-card-deployment-2')).toHaveTextContent('Remote');
  });

  it('renders security verification tags', () => {
    render(<McpCatalogCard server={mockServer} />);
    expect(screen.getByText('Verified source')).toBeInTheDocument();
    expect(screen.getByText('SAST')).toBeInTheDocument();
  });

  it('does not render security section when securityIndicators is empty', () => {
    render(<McpCatalogCard server={{ ...mockServer, id: 3, securityIndicators: undefined }} />);
    expect(screen.queryByText('Verified source')).not.toBeInTheDocument();
  });

  it('renders card with data-testid for the server id', () => {
    render(<McpCatalogCard server={mockServer} />);
    expect(screen.getByTestId('mcp-catalog-card-1')).toBeInTheDocument();
  });
});

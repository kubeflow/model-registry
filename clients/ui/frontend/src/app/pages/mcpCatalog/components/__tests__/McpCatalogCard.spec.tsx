import '@testing-library/jest-dom';
import React from 'react';
import { render, screen } from '@testing-library/react';
import McpCatalogCard from '~/app/pages/mcpCatalog/components/McpCatalogCard';
import type { McpServerMock } from '~/app/pages/mcpCatalog/types/mcpServer';

jest.mock('react-router-dom', () => ({
  Link: ({ to, children }: { to: string; children: React.ReactNode }) => (
    <a href={to}>{children}</a>
  ),
}));

const mockServer: McpServerMock = {
  id: 'test-id',
  name: 'Test MCP Server',
  description: 'Test description for the server.',
  deploymentMode: 'Local',
  securityVerification: ['Verified source', 'SAST'],
  category: 'sample',
};

describe('McpCatalogCard', () => {
  it('renders server name and description', () => {
    render(<McpCatalogCard server={mockServer} />);
    expect(screen.getByTestId('mcp-catalog-card-name-test-id')).toHaveTextContent(
      'Test MCP Server',
    );
    expect(screen.getByTestId('mcp-catalog-card-description-test-id')).toHaveTextContent(
      'Test description for the server.',
    );
  });

  it('renders deployment mode label for Local', () => {
    render(<McpCatalogCard server={mockServer} />);
    expect(screen.getByTestId('mcp-catalog-card-deployment-test-id')).toHaveTextContent(
      'Local to cluster',
    );
  });

  it('renders deployment mode label for Remote', () => {
    render(<McpCatalogCard server={{ ...mockServer, id: 'r1', deploymentMode: 'Remote' }} />);
    expect(screen.getByTestId('mcp-catalog-card-deployment-r1')).toHaveTextContent('Remote');
  });

  it('renders security verification tags', () => {
    render(<McpCatalogCard server={mockServer} />);
    expect(screen.getByText('Verified source')).toBeInTheDocument();
    expect(screen.getByText('SAST')).toBeInTheDocument();
  });

  it('does not render security section when securityVerification is empty', () => {
    render(<McpCatalogCard server={{ ...mockServer, id: 'no-sec', securityVerification: [] }} />);
    expect(screen.queryByText('Verified source')).not.toBeInTheDocument();
  });

  it('renders card with data-testid for the server id', () => {
    render(<McpCatalogCard server={mockServer} />);
    expect(screen.getByTestId('mcp-catalog-card-test-id')).toBeInTheDocument();
  });
});

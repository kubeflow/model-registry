/* eslint-disable camelcase */
import type { McpServer } from '~/app/mcpServerCatalogTypes';

export const mockMcpServers: McpServer[] = [
  {
    id: 1,
    name: 'Kubernetes',
    description:
      'Control and inspect Kubernetes clusters using natural language queries for health, resources, and deployments.',
    deploymentMode: 'local',
    securityIndicators: { verifiedSource: true, sast: true },
    source_id: 'sample',
    toolCount: 0,
  },
  {
    id: 2,
    name: 'GitHub',
    description:
      'Integrate with GitHub repositories, issues, and pull requests using natural language.',
    deploymentMode: 'remote',
    securityIndicators: { verifiedSource: true, secureEndpoint: true },
    source_id: 'sample',
    toolCount: 0,
  },
  {
    id: 3,
    name: 'Slack',
    description: 'Search and interact with Slack workspaces, channels, and messages.',
    deploymentMode: 'remote',
    securityIndicators: { verifiedSource: true },
    source_id: 'sample',
    toolCount: 0,
  },
  {
    id: 4,
    name: 'PostgreSQL',
    description: 'Query and manage PostgreSQL databases using natural language.',
    deploymentMode: 'local',
    securityIndicators: { readOnlyTools: true },
    source_id: 'other',
    toolCount: 0,
  },
  {
    id: 5,
    name: 'Custom MCP Server',
    description: 'A custom MCP server for extended integrations and workflows.',
    deploymentMode: 'remote',
    source_id: 'other',
    toolCount: 0,
  },
];

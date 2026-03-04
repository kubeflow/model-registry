import type { McpServerMock } from '~/app/pages/mcpCatalog/types/mcpServer';

export const mockMcpServers: McpServerMock[] = [
  {
    id: 'kubernetes',
    name: 'Kubernetes',
    description:
      'Control and inspect Kubernetes clusters using natural language queries for health, resources, and deployments.',
    deploymentMode: 'Local',
    securityVerification: ['Verified source', 'SAST'],
    category: 'sample',
  },
  {
    id: 'github',
    name: 'GitHub',
    description:
      'Integrate with GitHub repositories, issues, and pull requests using natural language.',
    deploymentMode: 'Remote',
    securityVerification: ['Verified source', 'Secure endpoint'],
    category: 'sample',
  },
  {
    id: 'slack',
    name: 'Slack',
    description: 'Search and interact with Slack workspaces, channels, and messages.',
    deploymentMode: 'Remote',
    securityVerification: ['Verified source'],
    category: 'sample',
  },
  {
    id: 'postgresql',
    name: 'PostgreSQL',
    description: 'Query and manage PostgreSQL databases using natural language.',
    deploymentMode: 'Local',
    securityVerification: ['Read only tools'],
    category: 'other',
  },
  {
    id: 'custom-mcp',
    name: 'Custom MCP Server',
    description: 'A custom MCP server for extended integrations and workflows.',
    deploymentMode: 'Remote',
    securityVerification: [],
    category: 'other',
  },
];

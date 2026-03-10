export const testMcpServers = [
  {
    id: 1,
    name: 'Kubernetes',
    description: 'Control and inspect Kubernetes clusters.',
    deploymentMode: 'local' as const,
    securityIndicators: { verifiedSource: true, sast: true },
    source_id: 'sample', // eslint-disable-line camelcase
    toolCount: 0,
    license: 'Apache-2.0',
    licenseLink: 'https://opensource.org/licenses/Apache-2.0',
    version: '1.0.0',
    provider: 'Kubernetes',
    sourceCode: 'https://github.com/kubernetes/mcp-server',
    repositoryUrl: 'https://github.com/kubernetes/mcp-server',
    artifacts: [{ uri: 'ghcr.io/kubernetes/mcp-server:latest' }],
    transports: ['http-streaming' as const],
    tags: ['kubernetes', 'infrastructure'],
    readme: '# Kubernetes MCP Server\n\n### Overview\n\nManage clusters with `kubectl`.',
  },
  {
    id: 2,
    name: 'GitHub',
    description: 'Integrate with GitHub repositories.',
    deploymentMode: 'remote' as const,
    securityIndicators: { verifiedSource: true, secureEndpoint: true },
    source_id: 'sample', // eslint-disable-line camelcase
    toolCount: 0,
    license: 'MIT',
    version: '2.1.0',
    provider: 'GitHub',
    sourceCode: 'https://github.com/github/mcp-server',
    artifacts: [{ uri: 'ghcr.io/github/mcp-server:latest' }],
    transports: ['sse' as const],
    tags: ['github', 'vcs'],
  },
  {
    id: 3,
    name: 'Custom MCP Server',
    description: 'A custom MCP server without README.',
    deploymentMode: 'local' as const,
    securityIndicators: {},
    source_id: 'sample', // eslint-disable-line camelcase
    toolCount: 0,
  },
];

export const testFilterOptions = {
  filters: {
    deploymentMode: { type: 'string' as const, values: ['Remote', 'Local'] },
    supportedTransports: { type: 'string' as const, values: ['SSE', 'http-streaming'] },
    license: { type: 'string' as const, values: ['MIT', 'Apache-2.0'] },
    labels: {
      type: 'string' as const,
      values: ['kubernetes', 'github', 'database', 'monitoring', 'security', 'automation'],
    },
    securityVerification: {
      type: 'string' as const,
      values: ['Verified source', 'Secure endpoint', 'SAST', 'Read only tools'],
    },
  },
};

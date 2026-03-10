/* eslint-disable camelcase */
import type { McpServer, McpTool } from '~/app/mcpServerCatalogTypes';

const dynatraceTools: McpTool[] = [
  {
    name: 'create_maintenance_window',
    description: 'Create a maintenance window to suppress alerts',
    accessType: 'read_write',
    parameters: [
      { name: 'name', type: 'string', description: 'Maintenance window name', required: true },
      { name: 'start_time', type: 'string', description: 'Start time (ISO format)', required: true },
      { name: 'end_time', type: 'string', description: 'End time (ISO format)', required: true },
      { name: 'entity_ids', type: 'array', description: 'Entity IDs to include in maintenance', required: false },
      { name: 'description', type: 'string', description: 'Maintenance window description', required: false },
    ],
  },
  {
    name: 'execute_dql',
    description: 'Execute Dynatrace Query Language (DQL) queries',
    accessType: 'read_only',
    parameters: [
      { name: 'query', type: 'string', description: 'DQL query string', required: true },
      { name: 'max_results', type: 'integer', description: 'Maximum number of results to return', required: false },
    ],
  },
  {
    name: 'get_problems',
    description: 'Retrieve current problems and incidents',
    accessType: 'read_only',
    parameters: [],
  },
  {
    name: 'get_service_health',
    description: 'Get health status of services',
    accessType: 'read_only',
    parameters: [
      { name: 'service_id', type: 'string', description: 'Service identifier', required: false },
    ],
  },
  {
    name: 'get_vulnerabilities',
    description: 'Retrieve security vulnerability data',
    accessType: 'read_only',
    parameters: [
      { name: 'severity', type: 'string', description: 'Filter by severity level (critical, high, medium, low)', required: false },
      { name: 'status', type: 'string', description: 'Filter by vulnerability status (open, resolved, suppressed)', required: false },
    ],
  },
  {
    name: 'query',
    description: 'Execute PromQL queries against the Prometheus time-series database',
    accessType: 'read_only',
    parameters: [
      { name: 'query', type: 'string', description: 'PromQL query expression', required: true },
      { name: 'time', type: 'string', description: 'Evaluation timestamp (ISO format)', required: false },
      { name: 'timeout', type: 'string', description: 'Evaluation timeout duration', required: false },
    ],
  },
  {
    name: 'query_range',
    description: 'Execute a PromQL range query over a time window',
    accessType: 'read_only',
    parameters: [
      { name: 'query', type: 'string', description: 'PromQL query expression', required: true },
      { name: 'start', type: 'string', description: 'Start timestamp (ISO format)', required: true },
      { name: 'end', type: 'string', description: 'End timestamp (ISO format)', required: true },
      { name: 'step', type: 'string', description: 'Query resolution step width (e.g. 15s, 1m, 5m)', required: false },
    ],
  },
  {
    name: 'get_metric_metadata',
    description: 'Retrieve metadata about a specific Prometheus metric',
    accessType: 'read_only',
    parameters: [
      { name: 'metric', type: 'string', description: 'Metric name', required: true },
    ],
  },
  {
    name: 'list_targets',
    description: 'List all active and dropped scrape targets',
    accessType: 'read_only',
    parameters: [
      { name: 'state', type: 'string', description: 'Filter by target state (active, dropped, any)', required: false },
    ],
  },
  {
    name: 'delete_alert_silence',
    description: 'Delete an alert silence by ID',
    accessType: 'execute',
    parameters: [
      { name: 'silence_id', type: 'string', description: 'ID of the silence to delete', required: true },
    ],
  },
  {
    name: 'deploy_model',
    description: 'Deploy a machine learning model to a Kubernetes cluster',
    accessType: 'execute',
    parameters: [
      { name: 'model_name', type: 'string', description: 'Name of the model to deploy', required: true },
      { name: 'namespace', type: 'string', description: 'Target Kubernetes namespace', required: true },
      { name: 'replicas', type: 'integer', description: 'Number of replicas', required: false },
      { name: 'image', type: 'string', description: 'Container image URI for the model server', required: true },
    ],
  },
  {
    name: 'get_alerts',
    description: 'Retrieve active alerts from the alerting system',
    accessType: 'read_only',
    parameters: [
      { name: 'severity', type: 'string', description: 'Filter alerts by severity', required: false },
      { name: 'active_only', type: 'boolean', description: 'Show only active alerts', required: false },
    ],
  },
];

export const mockMcpServers: McpServer[] = [
  {
    id: 1,
    name: 'Kubernetes',
    description:
      'Control and inspect Kubernetes clusters using natural language queries for health, resources, and deployments.',
    deploymentMode: 'local',
    securityIndicators: { verifiedSource: true, sast: true },
    source_id: 'sample',
    toolCount: 12,
    license: 'Apache-2.0',
    licenseLink: 'https://opensource.org/licenses/Apache-2.0',
    version: '1.2.0',
    provider: 'kubernetes-sigs',
    tags: ['kubernetes', 'infrastructure', 'containers', 'orchestration'],
    transports: ['http'],
    artifacts: [{ uri: 'quay.io/kubernetes-sigs/mcp-kubernetes:1.2.0' }],
    sourceCode: 'kubernetes-sigs/mcp-kubernetes',
    repositoryUrl: 'https://github.com/kubernetes-sigs/mcp-kubernetes',
    readme:
      '# Kubernetes MCP Server\n\nThe Kubernetes MCP Server allows AI Assistants to interact with Kubernetes clusters.\n\n## Quickstart\n\nInstall via `npx`:\n\n```bash\nnpx @kubernetes-sigs/mcp-kubernetes\n```\n\n## Use Cases\n\n- **Cluster inspection** - Query pod, deployment, and service status\n- **Resource management** - Create, update, and delete resources\n- **Health monitoring** - Check cluster health and resource utilization\n',
    lastUpdated: '1709913600000',
    tools: dynatraceTools,
  },
  {
    id: 2,
    name: 'GitHub',
    description:
      'Integrate with GitHub repositories, issues, and pull requests using natural language.',
    deploymentMode: 'remote',
    securityIndicators: { verifiedSource: true, secureEndpoint: true },
    source_id: 'sample',
    toolCount: 12,
    license: 'MIT',
    licenseLink: 'https://opensource.org/licenses/MIT',
    version: '3.0.1',
    provider: 'github',
    tags: ['github', 'vcs', 'devops'],
    transports: ['http', 'sse'],
    artifacts: [{ uri: 'quay.io/github/mcp-github:3.0.1' }],
    sourceCode: 'github/mcp-server',
    repositoryUrl: 'https://github.com/github/mcp-server',
    tools: dynatraceTools,
  },
  {
    id: 3,
    name: 'Slack',
    description: 'Search and interact with Slack workspaces, channels, and messages.',
    deploymentMode: 'remote',
    securityIndicators: { verifiedSource: true },
    source_id: 'sample',
    toolCount: 12,
    version: '2.1.0',
    provider: 'slack',
    tags: ['slack', 'messaging'],
    transports: ['sse'],
    tools: dynatraceTools,
  },
  {
    id: 4,
    name: 'PostgreSQL',
    description: 'Query and manage PostgreSQL databases using natural language.',
    deploymentMode: 'local',
    securityIndicators: { readOnlyTools: true },
    source_id: 'other',
    toolCount: 12,
    version: '1.0.0',
    provider: 'postgres-community',
    transports: ['stdio'],
    tools: dynatraceTools,
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

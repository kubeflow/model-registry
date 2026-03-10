import type { McpTool } from '~/app/mcpServerCatalogTypes';

export const mockMcpTools: McpTool[] = [
  {
    name: 'create_maintenance_window',
    description: 'Create a maintenance window to suppress alerts',
    accessType: 'read_write',
    parameters: [
      { name: 'name', type: 'string', description: 'Maintenance window name', required: true },
      { name: 'start_time', type: 'string', description: 'Start time (ISO format)', required: true },
      { name: 'end_time', type: 'string', description: 'End time (ISO format)', required: true },
      {
        name: 'entity_ids',
        type: 'array',
        description: 'Entity IDs to include in maintenance',
        required: false,
      },
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

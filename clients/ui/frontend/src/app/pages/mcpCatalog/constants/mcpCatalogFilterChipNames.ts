import type { McpFilterCategoryKey } from '~/app/pages/mcpCatalog/types/mcpCatalogFilterOptions';

export const MCP_FILTER_CATEGORY_NAMES: Record<McpFilterCategoryKey, string> = {
  deploymentMode: 'Deployment mode',
  supportedTransports: 'Supported transports',
  license: 'License',
  labels: 'Labels',
  securityVerification: 'Security & Verification',
};

export const MCP_FILTER_KEYS: McpFilterCategoryKey[] = [
  'deploymentMode',
  'supportedTransports',
  'license',
  'labels',
  'securityVerification',
];

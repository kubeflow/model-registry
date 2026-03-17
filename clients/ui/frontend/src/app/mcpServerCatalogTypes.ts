import { PaginationParams } from './modelCatalogTypes';

export type McpDeploymentMode = 'local' | 'remote';

export type McpTransportType = 'stdio' | 'sse' | 'http';

export type McpToolAccessType = 'read_only' | 'read_write' | 'execute';

export type McpEndpoints = {
  http?: string;
  sse?: string;
};

export type McpArtifact = {
  uri: string;
  createTimeSinceEpoch?: string;
  lastUpdateTimeSinceEpoch?: string;
};

export type McpSecurityIndicator = {
  verifiedSource?: boolean;
  secureEndpoint?: boolean;
  sast?: boolean;
  readOnlyTools?: boolean;
};

export type McpToolParameter = {
  name: string;
  type: string;
  description?: string;
  required: boolean;
};

export type McpTool = {
  name: string;
  description?: string;
  accessType: McpToolAccessType;
  parameters?: McpToolParameter[];
  revoked?: boolean;
  revokedReason?: string;
};

export type McpServer = {
  id: string;
  name: string;
  source_id?: string;
  description?: string;
  logo?: string;
  license?: string;
  licenseLink?: string;
  provider?: string;
  version?: string;
  tags?: string[];
  toolCount: number;
  tools?: McpTool[];
  securityIndicators?: McpSecurityIndicator;
  documentationUrl?: string;
  repositoryUrl?: string;
  sourceCode?: string;
  lastUpdated?: string;
  publishedDate?: string;
  artifacts?: McpArtifact[];
  transports?: McpTransportType[];
  readme?: string;
  deploymentMode?: McpDeploymentMode;
  endpoints?: McpEndpoints;
};

export type McpServerList = PaginationParams & { items?: McpServer[] };

export type McpToolWithServer = {
  serverId: string;
  serverName: string;
  tool: McpTool;
};

export type McpToolList = PaginationParams & { items?: McpToolWithServer[] };

export type McpServerListParams = {
  sourceLabel?: string;
  pageSize?: number | string;
  nextPageToken?: string;
  filterQuery?: string;
  namedQuery?: string;
  includeTools?: boolean;
  toolLimit?: number;
  orderBy?: string;
  sortOrder?: string;
  name?: string;
  q?: string;
};

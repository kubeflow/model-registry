/**
 * MCP (Model Context Protocol) Catalog Types
 *
 * These types define the structure of MCP servers displayed in the catalog.
 */

import { APIOptions } from 'mod-arch-core';

export enum McpToolAccessType {
  READ_ONLY = 'read_only',
  READ_WRITE = 'read_write',
  EXECUTE = 'execute',
}

export type McpSecurityIndicator = {
  verifiedSource: boolean;
  secureEndpoint: boolean;
  sast: boolean;
  readOnlyTools: boolean;
};

export type McpToolParameter = {
  name: string;
  type: string;
  description: string;
  required: boolean;
};

/**
 * MetadataValue types following Model Registry patterns.
 * Tags are stored as MetadataStringValue entries with empty string_value (label pattern).
 * Security indicators are stored as MetadataBoolValue entries.
 */
export type MetadataStringValue = {
  string_value: string;
  metadataType: 'MetadataStringValue';
};

export type MetadataBoolValue = {
  bool_value: boolean;
  metadataType: 'MetadataBoolValue';
};

export type MetadataIntValue = {
  int_value: string;
  metadataType: 'MetadataIntValue';
};

export type MetadataDoubleValue = {
  double_value: number;
  metadataType: 'MetadataDoubleValue';
};

export type MetadataStructValue = {
  struct_value: string;
  metadataType: 'MetadataStructValue';
};

export type McpMetadataValue =
  | MetadataStringValue
  | MetadataBoolValue
  | MetadataIntValue
  | MetadataDoubleValue
  | MetadataStructValue;

/**
 * Custom properties map following Model Registry patterns.
 * Keys are property names, values are MetadataValue discriminated union types.
 */
export type McpCustomProperties = Record<string, McpMetadataValue>;

export type McpTool = {
  name: string;
  description: string;
  accessType: McpToolAccessType;
  parameters?: McpToolParameter[];
  /** Whether this tool has been revoked. Revoked tools should not be invoked by AI agents. */
  revoked?: boolean;
  /** Human-readable reason why the tool was revoked. */
  revokedReason?: string;
  customProperties?: McpCustomProperties;
};

export enum McpTransportType {
  STDIO = 'stdio',
  SSE = 'sse',
  HTTP = 'http',
}

export enum McpDeploymentMode {
  LOCAL = 'local',
  REMOTE = 'remote',
}

export type McpEndpoints = {
  http?: string;
  sse?: string;
};

/**
 * Artifact for an MCP server (e.g., OCI image for local deployment).
 */
export type McpArtifact = {
  uri: string;
  createTimeSinceEpoch?: string;
  lastUpdateTimeSinceEpoch?: string;
};

export type McpServer = {
  id: string;
  name: string;
  description: string;
  source_id?: string;
  logo?: string;
  license?: string;
  license_link?: string;
  provider?: string;
  version?: string;
  tags?: string[];
  tools?: McpTool[];
  securityIndicators?: McpSecurityIndicator;
  documentationUrl?: string;
  repositoryUrl?: string;
  sourceCode?: string;
  lastUpdated?: string;
  publishedDate?: string;
  artifacts?: McpArtifact[];
  transports?: McpTransportType[];
  apiKey?: string;
  readme?: string;
  deploymentMode?: McpDeploymentMode;
  endpoints?: McpEndpoints;
  /**
   * Custom properties following Model Registry patterns.
   * Tags are stored as MetadataStringValue entries with empty string_value (label pattern).
   * Security indicators are stored as MetadataBoolValue entries.
   */
  customProperties?: McpCustomProperties;
};

/**
 * Asset type for catalog sources
 */
export enum CatalogAssetType {
  MODELS = 'models',
  MCP_SERVERS = 'mcp_servers',
}

/**
 * MCP Catalog Source - represents a source of MCP servers
 */
export type McpCatalogSource = {
  id: string;
  name: string;
  labels: string[];
  enabled?: boolean;
  assetType?: CatalogAssetType;
  status?: 'available' | 'error' | 'disabled';
  error?: string;
};

export type McpCatalogSourceList = {
  items?: McpCatalogSource[];
  size: number;
  pageSize: number;
  nextPageToken?: string;
};

/**
 * Category names for MCP server organization
 */
export enum McpCategoryName {
  allServers = 'All MCP servers',
  communityAndCustomServers = 'Community and custom',
}

/**
 * Special label for sources without labels
 */
export enum McpSourceLabel {
  other = 'null',
}

export type McpServerList = {
  items: McpServer[];
  size: number;
  pageSize: number;
  nextPageToken?: string;
};

/**
 * Filter key constants for MCP Catalog filters.
 * Following Model Catalog pattern with string filter keys.
 */
export enum McpFilterKey {
  PROVIDER = 'provider',
  LICENSE = 'license',
  TAGS = 'tags',
  TRANSPORTS = 'transports',
  DEPLOYMENT_MODE = 'deploymentMode',
}

/**
 * Filter state for MCP Catalog.
 * Each filter key maps to an array of selected values.
 */
export type McpFilterState = {
  [McpFilterKey.PROVIDER]: string[];
  [McpFilterKey.LICENSE]: string[];
  [McpFilterKey.TAGS]: string[];
  [McpFilterKey.TRANSPORTS]: string[];
  [McpFilterKey.DEPLOYMENT_MODE]: string[];
};

/**
 * Filter option structure matching backend response.
 * Each filter option has a type and available values.
 */
export type McpFilterOption = {
  type: string;
  values?: unknown[];
};

/**
 * Filter options list response from the backend.
 * Keys are filter field names, values describe available options.
 */
export type McpFilterOptionsList = {
  filters?: Record<string, McpFilterOption>;
};

/**
 * API interface for MCP Catalog operations
 */
export type McpCatalogAPIs = {
  getMcpServers: (
    opts: APIOptions,
    sourceLabel?: string,
    pageSize?: number,
    filterQuery?: string,
    searchTerm?: string,
  ) => Promise<McpServerList>;
  getMcpServer: (opts: APIOptions, serverId: string) => Promise<McpServer>;
  getMcpSources: (opts: APIOptions) => Promise<McpCatalogSourceList>;
  getMcpFilterOptions: (opts: APIOptions) => Promise<McpFilterOptionsList>;
};

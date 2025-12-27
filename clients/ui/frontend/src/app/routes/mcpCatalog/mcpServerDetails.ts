import { mcpCatalogUrl } from './mcpCatalog';

export type McpServerDetailsParams = {
  serverId: string;
};

export const getMcpServerDetailsRoute = (serverId: string): string =>
  `${mcpCatalogUrl()}/${encodeURIComponent(serverId)}`;

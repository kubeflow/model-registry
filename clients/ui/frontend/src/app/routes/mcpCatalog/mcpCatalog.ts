export const mcpCatalogUrl = (): string => '/mcp-catalog';

export const mcpServerDetailsUrl = (serverId: string | number): string =>
  `${mcpCatalogUrl()}/${encodeURIComponent(String(serverId))}`;

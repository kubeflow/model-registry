import type { McpCatalogFiltersState } from '~/app/context/mcpCatalog/McpCatalogContext';
import { MCP_FILTER_KEYS } from '~/app/pages/mcpCatalog/constants/mcpCatalogFilterChipNames';
import type { McpSecurityIndicator } from '~/app/mcpServerCatalogTypes';

const SECURITY_INDICATOR_LABELS: Record<keyof McpSecurityIndicator, string> = {
  verifiedSource: 'Verified source',
  secureEndpoint: 'Secure endpoint',
  sast: 'SAST',
  readOnlyTools: 'Read only tools',
};

const SECURITY_INDICATOR_KEYS: (keyof McpSecurityIndicator)[] = [
  'verifiedSource',
  'secureEndpoint',
  'sast',
  'readOnlyTools',
];

export const getSecurityIndicatorLabels = (
  securityIndicators?: McpSecurityIndicator | null,
): string[] => {
  if (!securityIndicators) {
    return [];
  }
  return SECURITY_INDICATOR_KEYS.filter((key) => Boolean(securityIndicators[key])).map(
    (key) => SECURITY_INDICATOR_LABELS[key],
  );
};

export const hasMcpFiltersApplied = (
  filters: McpCatalogFiltersState,
  searchQuery: string,
): boolean => {
  if (searchQuery && searchQuery.trim().length > 0) {
    return true;
  }
  for (const key of MCP_FILTER_KEYS) {
    const value = filters[key];
    if (Array.isArray(value) && value.length > 0) {
      return true;
    }
  }
  return false;
};

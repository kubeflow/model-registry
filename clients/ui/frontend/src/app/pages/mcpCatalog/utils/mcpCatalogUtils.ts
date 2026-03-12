import type { McpCatalogFiltersState } from '~/app/pages/mcpCatalog/types/mcpCatalogFilterOptions';
import { MCP_FILTER_KEYS } from '~/app/pages/mcpCatalog/const';
import type {
  McpSecurityIndicator,
  McpServer,
  McpTransportType,
} from '~/app/mcpServerCatalogTypes';

const SECURITY_INDICATOR_LABELS: Record<keyof McpSecurityIndicator, string> = {
  verifiedSource: 'Verified source',
  secureEndpoint: 'Secure endpoint',
  sast: 'SAST',
  readOnlyTools: 'Read only tools',
};

const SECURITY_LABEL_TO_KEY: Record<string, keyof McpSecurityIndicator | undefined> = {
  [SECURITY_INDICATOR_LABELS.verifiedSource]: 'verifiedSource',
  [SECURITY_INDICATOR_LABELS.secureEndpoint]: 'secureEndpoint',
  [SECURITY_INDICATOR_LABELS.sast]: 'sast',
  [SECURITY_INDICATOR_LABELS.readOnlyTools]: 'readOnlyTools',
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

export function filterMcpServersBySearchQuery(
  items: McpServer[],
  searchQuery: string,
): McpServer[] {
  const q = searchQuery.trim().toLowerCase();
  if (!q) {
    return items;
  }
  return items.filter((server) => {
    const name = server.name.toLowerCase();
    const description = server.description?.toLowerCase() ?? '';
    return name.includes(q) || description.includes(q);
  });
}

function isMcpTransport(s: string): s is McpTransportType {
  return s === 'stdio' || s === 'sse' || s === 'http';
}

function matchesDeploymentMode(server: McpServer, selected: string[]): boolean {
  if (selected.length === 0) {
    return true;
  }
  const mode = server.deploymentMode?.toLowerCase();
  if (!mode) {
    return false;
  }
  return selected.some((v) => v.toLowerCase() === mode);
}

function matchesLicense(server: McpServer, selected: string[]): boolean {
  if (selected.length === 0) {
    return true;
  }
  const license = server.license?.trim();
  if (!license) {
    return false;
  }
  return selected.some((v) => v.trim().toLowerCase() === license.toLowerCase());
}

function matchesLabels(server: McpServer, selected: string[]): boolean {
  if (selected.length === 0) {
    return true;
  }
  const tags = server.tags ?? [];
  return selected.some((s) => tags.includes(s));
}

function matchesTransports(server: McpServer, selected: string[]): boolean {
  if (selected.length === 0) {
    return true;
  }
  const transports: McpTransportType[] = server.transports ?? [];
  return selected.some((s) => isMcpTransport(s) && transports.includes(s));
}

function matchesSecurityVerification(server: McpServer, selected: string[]): boolean {
  if (selected.length === 0) {
    return true;
  }
  const ind = server.securityIndicators;
  if (!ind) {
    return false;
  }
  const selectedKeys = selected
    .map((label) => SECURITY_LABEL_TO_KEY[label])
    .filter((k): k is keyof McpSecurityIndicator => k !== undefined);
  if (selectedKeys.length === 0) {
    return false;
  }
  return selectedKeys.some((key) => Boolean(ind[key]));
}

export function filterMcpServersByFilters(
  items: McpServer[],
  filters: McpCatalogFiltersState,
): McpServer[] {
  const {
    deploymentMode: deploymentModeFilter,
    license: licenseFilter,
    labels: labelsFilter,
    supportedTransports: transportsFilter,
    securityIndicators: securityFilter,
  } = filters;
  return items.filter((server) => {
    if (deploymentModeFilter?.length && !matchesDeploymentMode(server, deploymentModeFilter)) {
      return false;
    }
    if (licenseFilter?.length && !matchesLicense(server, licenseFilter)) {
      return false;
    }
    if (labelsFilter?.length && !matchesLabels(server, labelsFilter)) {
      return false;
    }
    if (transportsFilter?.length && !matchesTransports(server, transportsFilter)) {
      return false;
    }
    if (securityFilter?.length && !matchesSecurityVerification(server, securityFilter)) {
      return false;
    }
    return true;
  });
}

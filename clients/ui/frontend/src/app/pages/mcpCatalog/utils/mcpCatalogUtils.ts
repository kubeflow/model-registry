import {
  McpCatalogSource,
  McpCatalogSourceList,
  McpCustomProperties,
  McpDeploymentMode,
  McpMetadataValue,
  McpSecurityIndicator,
  McpServer,
  MetadataBoolValue,
} from '~/app/pages/mcpCatalog/types';

/**
 * Filter state for MCP server sidebar filters
 */
export type McpServerFilterState = {
  selectedProviders: string[];
  selectedLicenses: string[];
  selectedTags: string[];
  selectedTransports: string[];
  selectedDeploymentModes: string[];
};

// Helper to wrap value in quotes for filterQuery
const wrapInQuotes = (v: string): string => `'${v}'`;

// Helper to create equality filter (matching Model Catalog pattern)
const eqFilter = (k: string, v: string): string => `${k}=${wrapInQuotes(v)}`;

// Helper to create IN filter (matching Model Catalog pattern)
const inFilter = (k: string, values: string[]): string =>
  `${k} IN (${values.map((v) => wrapInQuotes(v)).join(',')})`;

/**
 * Convert MCP server filter state to a filterQuery string for server-side filtering.
 * Uses SQL-like syntax matching Model Catalog: field='value', field IN ('a','b')
 *
 * @param filters - The current filter state
 * @returns A filterQuery string or empty string if no filters are active
 */
export const mcpFiltersToFilterQuery = (filters: McpServerFilterState): string => {
  const queryParts: string[] = [];

  // Provider filter
  if (filters.selectedProviders.length === 1) {
    queryParts.push(eqFilter('provider', filters.selectedProviders[0]));
  } else if (filters.selectedProviders.length > 1) {
    queryParts.push(inFilter('provider', filters.selectedProviders));
  }

  // License filter
  if (filters.selectedLicenses.length === 1) {
    queryParts.push(eqFilter('license', filters.selectedLicenses[0]));
  } else if (filters.selectedLicenses.length > 1) {
    queryParts.push(inFilter('license', filters.selectedLicenses));
  }

  // Tags filter (array field - backend handles array matching)
  if (filters.selectedTags.length === 1) {
    queryParts.push(eqFilter('tags', filters.selectedTags[0]));
  } else if (filters.selectedTags.length > 1) {
    queryParts.push(inFilter('tags', filters.selectedTags));
  }

  // Transport filter (array field - backend handles array matching)
  if (filters.selectedTransports.length === 1) {
    queryParts.push(eqFilter('transports', filters.selectedTransports[0]));
  } else if (filters.selectedTransports.length > 1) {
    queryParts.push(inFilter('transports', filters.selectedTransports));
  }

  // Deployment mode filter
  if (filters.selectedDeploymentModes.length === 1) {
    queryParts.push(eqFilter('deploymentMode', filters.selectedDeploymentModes[0]));
  } else if (filters.selectedDeploymentModes.length > 1) {
    queryParts.push(inFilter('deploymentMode', filters.selectedDeploymentModes));
  }

  return queryParts.length === 0 ? '' : queryParts.join(' AND ');
};

/**
 * Check if any filters are currently active
 */
export const hasMcpFiltersActive = (filters: McpServerFilterState): boolean =>
  filters.selectedProviders.length > 0 ||
  filters.selectedLicenses.length > 0 ||
  filters.selectedTags.length > 0 ||
  filters.selectedTransports.length > 0 ||
  filters.selectedDeploymentModes.length > 0;

/**
 * Format a transport type for display
 */
export const formatTransportType = (transport: string | undefined): string => {
  if (!transport) {
    return 'N/A';
  }
  return transport.toUpperCase();
};

/**
 * Format transports array for display
 */
export const formatTransports = (transports: string[] | undefined): string => {
  if (!transports || transports.length === 0) {
    return 'N/A';
  }
  return transports.map((t) => t.toUpperCase()).join(', ');
};

/**
 * Format deployment mode for display
 */
export const formatDeploymentMode = (deploymentMode: McpDeploymentMode | undefined): string => {
  switch (deploymentMode) {
    case McpDeploymentMode.REMOTE:
      return 'Remote';
    case McpDeploymentMode.LOCAL:
      return 'Local';
    default:
      return 'Local';
  }
};

/**
 * Check if the server is a remote deployment
 */
export const isRemoteMcpServer = (deploymentMode: McpDeploymentMode | undefined): boolean =>
  deploymentMode === McpDeploymentMode.REMOTE;

/**
 * Filter catalog sources to only include enabled sources
 */
export const filterEnabledMcpSources = (
  sources: McpCatalogSourceList | null,
): McpCatalogSourceList | null => {
  if (!sources) {
    return null;
  }

  const filteredItems = sources.items?.filter((source) => source.enabled !== false);

  return {
    ...sources,
    items: filteredItems || [],
    size: filteredItems?.length || 0,
  };
};

/**
 * Get unique source labels from catalog sources
 */
export const getUniqueMcpSourceLabels = (sources: McpCatalogSourceList | null): string[] => {
  if (!sources || !sources.items) {
    return [];
  }

  const allLabels = new Set<string>();

  sources.items.forEach((source) => {
    if (source.enabled !== false && source.labels.length > 0) {
      source.labels.forEach((label) => {
        if (label.trim()) {
          allLabels.add(label.trim());
        }
      });
    }
  });

  return Array.from(allLabels);
};

/**
 * Check if any sources exist without labels
 */
export const hasMcpSourcesWithoutLabels = (sources: McpCatalogSourceList | null): boolean => {
  if (!sources || !sources.items) {
    return false;
  }

  return sources.items.some((source) => {
    if (source.enabled !== false) {
      return source.labels.length === 0 || source.labels.every((label) => !label.trim());
    }
    return false;
  });
};

/**
 * Get source from source ID
 */
export const getMcpSourceFromSourceId = (
  sourceId: string,
  sources: McpCatalogSourceList | null,
): McpCatalogSource | undefined => {
  if (!sources || !sourceId || !sources.items) {
    return undefined;
  }

  return sources.items.find((source) => source.id === sourceId);
};

// ============================================================================
// CustomProperties Utility Functions
// Following Model Registry patterns for label and property handling
// ============================================================================

/**
 * Type guard to check if a MetadataValue is a string value
 */
const isMetadataStringValue = (value: McpMetadataValue): boolean =>
  value.metadataType === 'MetadataStringValue';

/**
 * Type guard to check if a MetadataValue is a bool value
 */
const isMetadataBoolValue = (value: McpMetadataValue): value is MetadataBoolValue =>
  value.metadataType === 'MetadataBoolValue';

/**
 * Extract tags from customProperties.
 * Tags are stored as MetadataStringValue entries with EMPTY string_value (label pattern).
 * This follows the Model Registry pattern where labels are empty string properties.
 *
 * @param customProperties - The customProperties map from McpServer
 * @returns Array of tag names (the keys of entries with empty string_value)
 */
export const getMcpServerTags = (customProperties: McpCustomProperties | undefined): string[] => {
  if (!customProperties) {
    return [];
  }

  return Object.entries(customProperties)
    .filter(([, value]) => {
      if (!isMetadataStringValue(value)) {
        return false;
      }
      // Tags have empty string_value (label pattern)
      return 'string_value' in value && value.string_value === '';
    })
    .map(([key]) => key);
};

/**
 * Extract security indicators from customProperties.
 * Security indicators are stored as MetadataBoolValue entries with specific keys.
 *
 * @param customProperties - The customProperties map from McpServer
 * @returns McpSecurityIndicator object or undefined if no security properties exist
 */
export const getMcpSecurityIndicatorsFromCustomProperties = (
  customProperties: McpCustomProperties | undefined,
): McpSecurityIndicator | undefined => {
  if (!customProperties) {
    return undefined;
  }

  const securityKeys = ['verifiedSource', 'secureEndpoint', 'sast', 'readOnlyTools'];
  const hasAnySecurityIndicator = securityKeys.some(
    (key) => key in customProperties && isMetadataBoolValue(customProperties[key]),
  );

  if (!hasAnySecurityIndicator) {
    return undefined;
  }

  const getBoolValue = (key: string): boolean => {
    const value = customProperties[key];
    return isMetadataBoolValue(value) ? value.bool_value : false;
  };

  return {
    verifiedSource: getBoolValue('verifiedSource'),
    secureEndpoint: getBoolValue('secureEndpoint'),
    sast: getBoolValue('sast'),
    readOnlyTools: getBoolValue('readOnlyTools'),
  };
};

/**
 * Extract string properties from customProperties (non-empty string_value entries).
 * This follows the Model Registry pattern where properties have non-empty string values.
 *
 * @param customProperties - The customProperties map from McpServer
 * @returns Record of property name to string value
 */
export const getMcpServerProperties = (
  customProperties: McpCustomProperties | undefined,
): Record<string, string> => {
  if (!customProperties) {
    return {};
  }

  const result: Record<string, string> = {};

  Object.entries(customProperties).forEach(([key, value]) => {
    if (isMetadataStringValue(value) && 'string_value' in value && value.string_value !== '') {
      result[key] = value.string_value;
    }
  });

  return result;
};

/**
 * Get a specific string property value from customProperties.
 *
 * @param customProperties - The customProperties map from McpServer
 * @param key - The property key to retrieve
 * @returns The string value or undefined if not found
 */
export const getMcpCustomPropertyString = (
  customProperties: McpCustomProperties | undefined,
  key: string,
): string | undefined => {
  if (!customProperties || !(key in customProperties)) {
    return undefined;
  }

  const value = customProperties[key];
  if (isMetadataStringValue(value) && 'string_value' in value) {
    return value.string_value;
  }

  return undefined;
};

/**
 * Get tags from an MCP server, preferring customProperties but falling back to tags array.
 * This supports backward compatibility where both formats may exist.
 *
 * @param server - The MCP server
 * @returns Array of tag names
 */
export const getServerTags = (server: McpServer): string[] => {
  // First try to get from customProperties (Model Registry aligned)
  const customPropsTags = getMcpServerTags(server.customProperties);
  if (customPropsTags.length > 0) {
    return customPropsTags;
  }

  // Fall back to first-class tags field (backward compatibility)
  return server.tags ?? [];
};

/**
 * Get security indicators from an MCP server, preferring customProperties but falling back.
 * This supports backward compatibility where both formats may exist.
 *
 * @param server - The MCP server
 * @returns McpSecurityIndicator object or undefined
 */
export const getServerSecurityIndicators = (
  server: McpServer,
): McpSecurityIndicator | undefined => {
  // First try to get from customProperties (Model Registry aligned)
  const customPropsIndicators = getMcpSecurityIndicatorsFromCustomProperties(
    server.customProperties,
  );
  if (customPropsIndicators) {
    return customPropsIndicators;
  }

  // Fall back to first-class securityIndicators field (backward compatibility)
  return server.securityIndicators;
};

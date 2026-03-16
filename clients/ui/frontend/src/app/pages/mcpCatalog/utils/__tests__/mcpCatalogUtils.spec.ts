import {
  filterMcpServersByFilters,
  filterMcpServersBySearchQuery,
  getSecurityIndicatorLabels,
  hasMcpFiltersApplied,
} from '~/app/pages/mcpCatalog/utils/mcpCatalogUtils';
import type { McpCatalogFiltersState } from '~/app/pages/mcpCatalog/types/mcpCatalogFilterOptions';
import type { McpServer } from '~/app/mcpServerCatalogTypes';

describe('getSecurityIndicatorLabels', () => {
  it('returns empty array when securityIndicators is undefined or null', () => {
    expect(getSecurityIndicatorLabels(undefined)).toEqual([]);
    expect(getSecurityIndicatorLabels(null)).toEqual([]);
  });

  it('returns labels for true boolean flags', () => {
    expect(getSecurityIndicatorLabels({ verifiedSource: true, sast: true })).toEqual([
      'Verified source',
      'SAST',
    ]);
    expect(getSecurityIndicatorLabels({ secureEndpoint: true })).toEqual(['Secure endpoint']);
    expect(getSecurityIndicatorLabels({ readOnlyTools: true })).toEqual(['Read only tools']);
  });

  it('ignores false or undefined flags', () => {
    expect(
      getSecurityIndicatorLabels({
        verifiedSource: false,
        secureEndpoint: true,
        sast: undefined,
      }),
    ).toEqual(['Secure endpoint']);
  });
});

describe('hasMcpFiltersApplied', () => {
  it('returns false when filters are empty and searchQuery is empty', () => {
    expect(hasMcpFiltersApplied({}, '')).toBe(false);
    expect(hasMcpFiltersApplied({}, '   ')).toBe(false);
  });

  it('returns true when searchQuery has non-empty trimmed content', () => {
    expect(hasMcpFiltersApplied({}, 'q')).toBe(true);
    expect(hasMcpFiltersApplied({}, '  query  ')).toBe(true);
  });

  it('returns false when all filter keys have empty arrays or are missing', () => {
    const filters: McpCatalogFiltersState = {
      deploymentMode: [],
      supportedTransports: [],
      license: [],
      labels: [],
      securityIndicators: [],
    };
    expect(hasMcpFiltersApplied(filters, '')).toBe(false);
  });

  it('returns true when deploymentMode has values', () => {
    expect(hasMcpFiltersApplied({ deploymentMode: ['Local'] }, '')).toBe(true);
  });

  it('returns true when supportedTransports has values', () => {
    expect(hasMcpFiltersApplied({ supportedTransports: ['stdio'] }, '')).toBe(true);
  });

  it('returns true when license has values', () => {
    expect(hasMcpFiltersApplied({ license: ['MIT'] }, '')).toBe(true);
  });

  it('returns true when labels has values', () => {
    expect(hasMcpFiltersApplied({ labels: ['Red Hat'] }, '')).toBe(true);
  });

  it('returns true when securityIndicators has values', () => {
    expect(hasMcpFiltersApplied({ securityIndicators: ['Verified'] }, '')).toBe(true);
  });

  it('returns true when multiple filter keys have values', () => {
    const filters: McpCatalogFiltersState = {
      deploymentMode: ['Local'],
      license: ['Apache-2.0'],
    };
    expect(hasMcpFiltersApplied(filters, '')).toBe(true);
  });

  it('ignores non-array filter values', () => {
    expect(hasMcpFiltersApplied({ deploymentMode: 'Local' as unknown as string[] }, '')).toBe(
      false,
    );
  });
});

describe('filterMcpServersByFilters', () => {
  const servers: McpServer[] = [
    { id: 1, name: 'A', deploymentMode: 'local', toolCount: 0 },
    { id: 2, name: 'B', deploymentMode: 'remote', toolCount: 0 },
    { id: 3, name: 'C', deploymentMode: 'remote', toolCount: 0 },
  ];

  it('returns all items when filters are empty', () => {
    expect(filterMcpServersByFilters(servers, {})).toEqual(servers);
    expect(filterMcpServersByFilters(servers, { deploymentMode: [] })).toEqual(servers);
  });

  it('filters by deploymentMode (Remote)', () => {
    const result = filterMcpServersByFilters(servers, { deploymentMode: ['Remote'] });
    expect(result).toHaveLength(2);
    expect(result.map((s) => s.name)).toEqual(['B', 'C']);
  });

  it('filters by deploymentMode (Local) case-insensitive', () => {
    const result = filterMcpServersByFilters(servers, { deploymentMode: ['local'] });
    expect(result).toHaveLength(1);
    expect(result[0].name).toBe('A');
  });

  it('filters by license when server has license', () => {
    const withLicense: McpServer[] = [
      { id: 1, name: 'X', license: 'MIT', toolCount: 0 },
      { id: 2, name: 'Y', license: 'Apache-2.0', toolCount: 0 },
    ];
    expect(filterMcpServersByFilters(withLicense, { license: ['MIT'] })).toHaveLength(1);
    expect(filterMcpServersByFilters(withLicense, { license: ['MIT'] })[0].name).toBe('X');
  });

  it('filters by securityIndicators label', () => {
    const withSecurity: McpServer[] = [
      { id: 1, name: 'S1', securityIndicators: { verifiedSource: true }, toolCount: 0 },
      { id: 2, name: 'S2', securityIndicators: { sast: true }, toolCount: 0 },
    ];
    const result = filterMcpServersByFilters(withSecurity, {
      securityIndicators: ['Verified source'],
    });
    expect(result).toHaveLength(1);
    expect(result[0].name).toBe('S1');
  });
});

describe('filterMcpServersBySearchQuery', () => {
  const servers: McpServer[] = [
    { id: 1, name: 'GitHub', description: 'Integrate with GitHub', toolCount: 0 },
    { id: 2, name: 'Slack', description: 'Search and interact with Slack', toolCount: 0 },
    { id: 3, name: 'PostgreSQL', description: 'Query databases', toolCount: 0 },
  ];

  it('returns all items when search is empty or whitespace', () => {
    expect(filterMcpServersBySearchQuery(servers, '')).toEqual(servers);
    expect(filterMcpServersBySearchQuery(servers, '   ')).toEqual(servers);
  });

  it('filters by name (case-insensitive)', () => {
    const result = filterMcpServersBySearchQuery(servers, 'git');
    expect(result).toHaveLength(1);
    expect(result[0].name).toBe('GitHub');
  });

  it('filters by description', () => {
    const result = filterMcpServersBySearchQuery(servers, 'databases');
    expect(result).toHaveLength(1);
    expect(result[0].name).toBe('PostgreSQL');
  });
});

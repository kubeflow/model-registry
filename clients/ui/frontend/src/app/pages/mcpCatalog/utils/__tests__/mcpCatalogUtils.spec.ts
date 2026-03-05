import {
  getSecurityIndicatorLabels,
  hasMcpFiltersApplied,
} from '~/app/pages/mcpCatalog/utils/mcpCatalogUtils';
import type { McpCatalogFiltersState } from '~/app/pages/mcpCatalog/types/mcpCatalogFilterOptions';

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
      securityVerification: [],
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

  it('returns true when securityVerification has values', () => {
    expect(hasMcpFiltersApplied({ securityVerification: ['Verified'] }, '')).toBe(true);
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

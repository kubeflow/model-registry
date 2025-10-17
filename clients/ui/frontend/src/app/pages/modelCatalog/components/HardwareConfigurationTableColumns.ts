import { SortableData } from 'mod-arch-shared';
import {
  CatalogPerformanceMetricsArtifact,
  PerformanceMetricsCustomProperties,
} from '~/app/modelCatalogTypes';
import {
  getHardwareConfiguration,
  getTotalRps,
} from '~/app/pages/modelCatalog/utils/performanceMetricsUtils';
import { getDoubleValue, getStringValue } from '~/app/utils';

export type HardwareConfigColumnField = keyof PerformanceMetricsCustomProperties | 'total_rps';

export type HardwareConfigColumn = Omit<
  SortableData<CatalogPerformanceMetricsArtifact>,
  'field'
> & { field: HardwareConfigColumnField };

// Note: The labels of most columns here include a non-breaking space (U+00a0) to selectively control word wrap.
// Your editor may highlight these and warn you that "The character U+00a0 is invisible."
// This character is ideally represented by the HTML entity &nbsp; but these strings can't contain HTML entities.
export const hardwareConfigColumns: HardwareConfigColumn[] = [
  {
    field: 'hardware_type',
    label: 'Hardware Configuration',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number => getHardwareConfiguration(a).localeCompare(getHardwareConfiguration(b)),
    isStickyColumn: true,
    stickyMinWidth: '162px',
    stickyLeftOffset: '0',
    hasRightBorder: true,
    modifier: 'wrap',
  },
  {
    field: 'requests_per_second',
    label: 'RPS per Replica',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'requests_per_second') -
      getDoubleValue(b.customProperties, 'requests_per_second'),
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'total_rps',
    label: 'Total RPS',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number => getTotalRps(a.customProperties) - getTotalRps(b.customProperties),
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'ttft_mean',
    label: 'TTFT Latency Mean',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'ttft_mean') -
      getDoubleValue(b.customProperties, 'ttft_mean'),
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'ttft_p90',
    label: 'TTFT Latency P90',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'ttft_p90') -
      getDoubleValue(b.customProperties, 'ttft_p90'),
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'ttft_p95',
    label: 'TTFT Latency P95',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'ttft_p95') -
      getDoubleValue(b.customProperties, 'ttft_p95'),
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'ttft_p99',
    label: 'TTFT Latency P99',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'ttft_p99') -
      getDoubleValue(b.customProperties, 'ttft_p99'),
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'e2e_mean',
    label: 'E2E Latency Mean',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'e2e_mean') -
      getDoubleValue(b.customProperties, 'e2e_mean'),
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'e2e_p90',
    label: 'E2E Latency P90',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'e2e_p90') - getDoubleValue(b.customProperties, 'e2e_p90'),
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'e2e_p95',
    label: 'E2E Latency P95',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'e2e_p95') - getDoubleValue(b.customProperties, 'e2e_p95'),
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'e2e_p99',
    label: 'E2E Latency P99',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'e2e_p99') - getDoubleValue(b.customProperties, 'e2e_p99'),
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'tps_mean',
    label: 'TPS Latency Mean',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'tps_mean') -
      getDoubleValue(b.customProperties, 'tps_mean'),
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'tps_p90',
    label: 'TPS Latency P90',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'tps_p90') - getDoubleValue(b.customProperties, 'tps_p90'),
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'tps_p95',
    label: 'TPS Latency P95',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'tps_p95') - getDoubleValue(b.customProperties, 'tps_p95'),
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'tps_p99',
    label: 'TPS Latency P99',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'tps_p99') - getDoubleValue(b.customProperties, 'tps_p99'),
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'itl_mean',
    label: 'ITL Latency Mean',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'itl_mean') -
      getDoubleValue(b.customProperties, 'itl_mean'),
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'itl_p90',
    label: 'ITL Latency P90',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'itl_p90') - getDoubleValue(b.customProperties, 'itl_p90'),
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'itl_p95',
    label: 'ITL Latency P95',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'itl_p95') - getDoubleValue(b.customProperties, 'itl_p95'),
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'itl_p99',
    label: 'ITL Latency P99',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'itl_p99') - getDoubleValue(b.customProperties, 'itl_p99'),
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'mean_input_tokens',
    label: 'Mean Input Tokens',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'mean_input_tokens') -
      getDoubleValue(b.customProperties, 'mean_input_tokens'),
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'mean_output_tokens',
    label: 'Mean Output Tokens',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'mean_output_tokens') -
      getDoubleValue(b.customProperties, 'mean_output_tokens'),
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'use_case',
    label: 'Use Case',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number => {
      const useCaseA = getStringValue(a.customProperties, 'use_case');
      const useCaseB = getStringValue(b.customProperties, 'use_case');
      return useCaseA.localeCompare(useCaseB);
    },
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'framework_version',
    label: 'vLLM Version',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number => {
      const versionA = getStringValue(a.customProperties, 'framework_version');
      const versionB = getStringValue(b.customProperties, 'framework_version');
      return versionA.localeCompare(versionB);
    },
    width: 20,
    modifier: 'wrap',
  },
];

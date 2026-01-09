import { SortableData } from 'mod-arch-shared';
import {
  CatalogPerformanceMetricsArtifact,
  PerformanceMetricsCustomProperties,
} from '~/app/modelCatalogTypes';
import {
  getHardwareConfiguration,
  getWorkloadType,
} from '~/app/pages/modelCatalog/utils/performanceMetricsUtils';
import { getDoubleValue, getIntValue, getStringValue } from '~/app/utils';
import { PerformancePropertyKey } from '~/concepts/modelCatalog/const';

export type HardwareConfigColumnField = keyof PerformanceMetricsCustomProperties;

export type HardwareConfigColumn = Omit<
  SortableData<CatalogPerformanceMetricsArtifact>,
  'field'
> & { field: HardwareConfigColumnField };

/*Non-breaking space constant (U+00A0) used to selectively control word wrap in column labels.
This prevents word wrapping into 3 lines (e.g., keeps "TTFT Latency" together instead of "TTFT\nLatency\nMean").
*/
const NBSP = '\u00A0';

export const hardwareConfigColumns: HardwareConfigColumn[] = [
  {
    field: PerformancePropertyKey.HARDWARE_TYPE,
    label: 'Hardware Configuration',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number => getHardwareConfiguration(a).localeCompare(getHardwareConfiguration(b)),
    isStickyColumn: true,
    stickyMinWidth: '162px',
    stickyLeftOffset: '0',
    modifier: 'wrap',
  },
  {
    field: PerformancePropertyKey.USE_CASE,
    label: 'Workload type',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number => getWorkloadType(a).localeCompare(getWorkloadType(b)),
    isStickyColumn: true,
    stickyMinWidth: '132px',
    stickyLeftOffset: '162px',
    modifier: 'wrap',
    hasRightBorder: true,
  },
  {
    field: PerformancePropertyKey.REQUESTS_PER_SECOND,
    label: `RPS${NBSP}per Replica`,
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, PerformancePropertyKey.REQUESTS_PER_SECOND) -
      getDoubleValue(b.customProperties, PerformancePropertyKey.REQUESTS_PER_SECOND),
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'replicas',
    label: 'Replicas',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getIntValue(a.customProperties, 'replicas') - getIntValue(b.customProperties, 'replicas'),
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'total_requests_per_second',
    label: `Total${NBSP}RPS`,
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'total_requests_per_second') -
      getDoubleValue(b.customProperties, 'total_requests_per_second'),
    width: 20,
    modifier: 'wrap',
  },
  {
    field: 'ttft_mean',
    label: `TTFT${NBSP}Latency Mean`,
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
    label: `TTFT${NBSP}Latency P90`,
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
    label: `TTFT${NBSP}Latency P95`,
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
    label: `TTFT${NBSP}Latency P99`,
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
    label: `E2E${NBSP}Latency Mean`,
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
    label: `E2E${NBSP}Latency P90`,
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
    label: `E2E${NBSP}Latency P95`,
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
    label: `E2E${NBSP}Latency P99`,
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
    label: `TPS${NBSP}Latency Mean`,
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
    label: `TPS${NBSP}Latency P90`,
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
    label: `TPS${NBSP}Latency P95`,
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
    label: `TPS${NBSP}Latency P99`,
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
    label: `ITL${NBSP}Latency Mean`,
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
    label: `ITL${NBSP}Latency P90`,
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
    label: `ITL${NBSP}Latency P95`,
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
    label: `ITL${NBSP}Latency P99`,
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
    label: `Mean${NBSP}Input Tokens`,
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
    label: `Mean${NBSP}Output Tokens`,
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

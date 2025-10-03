import { SortableData } from 'mod-arch-shared';
import { CatalogPerformanceMetricsArtifact } from '~/app/modelCatalogTypes';
import {
  getDoubleValue,
  getHardwareConfiguration,
  getIntValue,
  getStringValue,
  getTotalRps,
} from '~/app/pages/modelCatalog/utils/performanceMetricsUtils';

export const hardwareConfigColumns: SortableData<CatalogPerformanceMetricsArtifact>[] = [
  {
    field: 'hardware',
    label: 'Hardware\nConfiguration',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number => getHardwareConfiguration(a).localeCompare(getHardwareConfiguration(b)),
    width: 25,
  },
  {
    field: 'hardware_count',
    label: 'Total\nHardware',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getIntValue(a.customProperties, 'hardware_count') -
      getIntValue(b.customProperties, 'hardware_count'),
    width: 20,
  },
  {
    field: 'requests_per_second',
    label: 'RPS per\nReplica',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'requests_per_second') -
      getDoubleValue(b.customProperties, 'requests_per_second'),
    width: 20,
  },
  {
    field: 'total_rps',
    label: 'Total\nRPS',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number => getTotalRps(a.customProperties) - getTotalRps(b.customProperties),
    width: 20,
  },
  {
    field: 'ttft_mean',
    label: 'TTFT Latency\nMean',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'ttft_mean') -
      getDoubleValue(b.customProperties, 'ttft_mean'),
    width: 20,
  },
  {
    field: 'ttft_p90',
    label: 'TTFT Latency\nP90',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'ttft_p90') -
      getDoubleValue(b.customProperties, 'ttft_p90'),
    width: 20,
  },
  {
    field: 'ttft_p95',
    label: 'TTFT Latency\nP95',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'ttft_p95') -
      getDoubleValue(b.customProperties, 'ttft_p95'),
    width: 20,
  },
  {
    field: 'ttft_p99',
    label: 'TTFT Latency\nP99',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'ttft_p99') -
      getDoubleValue(b.customProperties, 'ttft_p99'),
    width: 20,
  },
  {
    field: 'e2e_mean',
    label: 'E2E Latency\nMean',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'e2e_mean') -
      getDoubleValue(b.customProperties, 'e2e_mean'),
    width: 20,
  },
  {
    field: 'e2e_p90',
    label: 'E2E Latency\nP90',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'e2e_p90') - getDoubleValue(b.customProperties, 'e2e_p90'),
    width: 20,
  },
  {
    field: 'e2e_p95',
    label: 'E2E Latency\nP95',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'e2e_p95') - getDoubleValue(b.customProperties, 'e2e_p95'),
    width: 20,
  },
  {
    field: 'e2e_p99',
    label: 'E2E Latency\nP99',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'e2e_p99') - getDoubleValue(b.customProperties, 'e2e_p99'),
    width: 20,
  },
  {
    field: 'tps_mean',
    label: 'TPS Latency\nMean',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'tps_mean') -
      getDoubleValue(b.customProperties, 'tps_mean'),
    width: 20,
  },
  {
    field: 'tps_p90',
    label: 'TPS Latency\nP90',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'tps_p90') - getDoubleValue(b.customProperties, 'tps_p90'),
    width: 20,
  },
  {
    field: 'tps_p95',
    label: 'TPS Latency\nP95',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'tps_p95') - getDoubleValue(b.customProperties, 'tps_p95'),
    width: 20,
  },
  {
    field: 'tps_p99',
    label: 'TPS Latency\nP99',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'tps_p99') - getDoubleValue(b.customProperties, 'tps_p99'),
    width: 20,
  },
  {
    field: 'itl_mean',
    label: 'ITL Latency\nMean',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'itl_mean') -
      getDoubleValue(b.customProperties, 'itl_mean'),
    width: 20,
  },
  {
    field: 'itl_p90',
    label: 'ITL Latency\nP90',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'itl_p90') - getDoubleValue(b.customProperties, 'itl_p90'),
    width: 20,
  },
  {
    field: 'itl_p95',
    label: 'ITL Latency\nP95',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'itl_p95') - getDoubleValue(b.customProperties, 'itl_p95'),
    width: 20,
  },
  {
    field: 'itl_p99',
    label: 'ITL Latency\nP99',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'itl_p99') - getDoubleValue(b.customProperties, 'itl_p99'),
    width: 20,
  },
  {
    field: 'max_input_tokens',
    label: 'Max Input\nTokens',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'max_input_tokens') -
      getDoubleValue(b.customProperties, 'max_input_tokens'),
    width: 20,
  },
  {
    field: 'max_output_tokens',
    label: 'Max Output\nTokens',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'max_output_tokens') -
      getDoubleValue(b.customProperties, 'max_output_tokens'),
    width: 20,
  },
  {
    field: 'mean_input_tokens',
    label: 'Mean Input\nTokens',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'mean_input_tokens') -
      getDoubleValue(b.customProperties, 'mean_input_tokens'),
    width: 20,
  },
  {
    field: 'mean_output_tokens',
    label: 'Mean Output\nTokens',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number =>
      getDoubleValue(a.customProperties, 'mean_output_tokens') -
      getDoubleValue(b.customProperties, 'mean_output_tokens'),
    width: 20,
  },
  {
    field: 'framework_version',
    label: 'vLLM\nVersion',
    sortable: (
      a: CatalogPerformanceMetricsArtifact,
      b: CatalogPerformanceMetricsArtifact,
    ): number => {
      const versionA = getStringValue(a.customProperties, 'framework_version');
      const versionB = getStringValue(b.customProperties, 'framework_version');
      return versionA.localeCompare(versionB);
    },
    width: 20,
  },
];

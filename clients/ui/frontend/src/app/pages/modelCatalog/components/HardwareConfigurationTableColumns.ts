import { SortableData } from 'mod-arch-shared';
import { HardwareConfiguration } from '~/app/pages/modelCatalog/types/hardwareConfiguration';

export const hardwareConfigColumns: SortableData<HardwareConfiguration>[] = [
  {
    field: 'hardwareConfiguration',
    label: 'Hardware\nConfiguration',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.hardwareConfiguration.localeCompare(b.hardwareConfiguration),
    width: 25,
  },
  {
    field: 'totalHardware',
    label: 'Total\nHardware',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.totalHardware - b.totalHardware,
    width: 20,
  },
  {
    field: 'rpsPerReplica',
    label: 'RPS per\nReplica',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.rpsPerReplica - b.rpsPerReplica,
    width: 20,
  },
  {
    field: 'totalRps',
    label: 'Total\nRPS',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.totalRps - b.totalRps,
    width: 20,
  },
  {
    field: 'ttftLatencyMean',
    label: 'TTFT Latency\nMean',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.ttftLatencyMean - b.ttftLatencyMean,
    width: 20,
  },
  {
    field: 'ttftLatencyP90',
    label: 'TTFT Latency\nP90',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.ttftLatencyP90 - b.ttftLatencyP90,
    width: 20,
  },
  {
    field: 'ttftLatencyP95',
    label: 'TTFT Latency\nP95',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.ttftLatencyP95 - b.ttftLatencyP95,
    width: 20,
  },
  {
    field: 'ttftLatencyP99',
    label: 'TTFT Latency\nP99',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.ttftLatencyP99 - b.ttftLatencyP99,
    width: 20,
  },
  {
    field: 'e2eLatencyMean',
    label: 'E2E Latency\nMean',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.e2eLatencyMean - b.e2eLatencyMean,
    width: 20,
  },
  {
    field: 'e2eLatencyP90',
    label: 'E2E Latency\nP90',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.e2eLatencyP90 - b.e2eLatencyP90,
    width: 20,
  },
  {
    field: 'e2eLatencyP95',
    label: 'E2E Latency\nP95',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.e2eLatencyP95 - b.e2eLatencyP95,
    width: 20,
  },
  {
    field: 'e2eLatencyP99',
    label: 'E2E Latency\nP99',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.e2eLatencyP99 - b.e2eLatencyP99,
    width: 20,
  },
  {
    field: 'tpsLatencyMean',
    label: 'TPS Latency\nMean',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.tpsLatencyMean - b.tpsLatencyMean,
    width: 20,
  },
  {
    field: 'tpsLatencyP90',
    label: 'TPS Latency\nP90',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.tpsLatencyP90 - b.tpsLatencyP90,
    width: 20,
  },
  {
    field: 'tpsLatencyP95',
    label: 'TPS Latency\nP95',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.tpsLatencyP95 - b.tpsLatencyP95,
    width: 20,
  },
  {
    field: 'tpsLatencyP99',
    label: 'TPS Latency\nP99',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.tpsLatencyP99 - b.tpsLatencyP99,
    width: 20,
  },
  {
    field: 'itlLatencyMean',
    label: 'ITL Latency\nMean',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.itlLatencyMean - b.itlLatencyMean,
    width: 20,
  },
  {
    field: 'itlLatencyP90',
    label: 'ITL Latency\nP90',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.itlLatencyP90 - b.itlLatencyP90,
    width: 20,
  },
  {
    field: 'itlLatencyP95',
    label: 'ITL Latency\nP95',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.itlLatencyP95 - b.itlLatencyP95,
    width: 20,
  },
  {
    field: 'itlLatencyP99',
    label: 'ITL Latency\nP99',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.itlLatencyP99 - b.itlLatencyP99,
    width: 20,
  },
  {
    field: 'maxInputTokens',
    label: 'Max Input\nTokens',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.maxInputTokens - b.maxInputTokens,
    width: 20,
  },
  {
    field: 'maxOutputTokens',
    label: 'Max Output\nTokens',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.maxOutputTokens - b.maxOutputTokens,
    width: 20,
  },
  {
    field: 'meanInputTokens',
    label: 'Mean Input\nTokens',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.meanInputTokens - b.meanInputTokens,
    width: 20,
  },
  {
    field: 'meanOutputTokens',
    label: 'Mean Output\nTokens',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.meanOutputTokens - b.meanOutputTokens,
    width: 20,
  },
  {
    field: 'vllmVersion',
    label: 'vLLM\nVersion',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.vllmVersion.localeCompare(b.vllmVersion),
    width: 20,
  },
  {
    field: 'guideLLMVersion',
    label: 'GuideLLM\nVersion',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.guideLLMVersion.localeCompare(b.guideLLMVersion),
    width: 20,
  },
  {
    field: 'rhaiisVersion',
    label: 'RHAIIS\nVersion',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.rhaiisVersion.localeCompare(b.rhaiisVersion),
    width: 20,
  },
];

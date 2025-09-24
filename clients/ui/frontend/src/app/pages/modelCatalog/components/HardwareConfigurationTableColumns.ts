import { SortableData } from 'mod-arch-shared';
import { HardwareConfiguration } from '~/app/pages/modelCatalog/types/hardwareConfiguration';

export const hardwareConfigColumns: SortableData<HardwareConfiguration>[] = [
  {
    field: 'hardwareConfiguration',
    label: 'Hardware configuration',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number => {
      const aConfig = `${a.hardwareCount} x ${a.hardwareType}`;
      const bConfig = `${b.hardwareCount} x ${b.hardwareType}`;
      return aConfig.localeCompare(bConfig);
    },
    width: 15,
  },
  {
    field: 'latency',
    label: 'Latency',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number => a.latency - b.latency,
    width: 10,
  },
  {
    field: 'throughput',
    label: 'Throughput',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      a.throughput - b.throughput,
    width: 10,
  },
  {
    field: 'tps',
    label: 'TPS',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number => a.tps - b.tps,
    width: 10,
  },
  {
    field: 'guideLLMVersion',
    label: 'GuideLLM version',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number => {
      const aVersion = a.guideLLMVersion || '';
      const bVersion = b.guideLLMVersion || '';
      return aVersion.localeCompare(bVersion);
    },
    width: 15,
  },
  {
    field: 'rhaiisVersion',
    label: 'RHAIIS version',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number => {
      const aVersion = a.rhaiisVersion || '';
      const bVersion = b.rhaiisVersion || '';
      return aVersion.localeCompare(bVersion);
    },
    width: 15,
  },
  {
    field: 'framework',
    label: 'Framework',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number => {
      const aFramework = a.framework || '';
      const bFramework = b.framework || '';
      return aFramework.localeCompare(bFramework);
    },
    width: 10,
  },
  {
    field: 'precision',
    label: 'Precision',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number => {
      const aPrecision = a.precision || '';
      const bPrecision = b.precision || '';
      return aPrecision.localeCompare(bPrecision);
    },
    width: 10,
  },
  {
    field: 'batchSize',
    label: 'Batch size',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      (a.batchSize || 0) - (b.batchSize || 0),
    width: 10,
  },
  {
    field: 'memoryUsage',
    label: 'Memory usage',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      (a.memoryUsage || 0) - (b.memoryUsage || 0),
    width: 10,
  },
  {
    field: 'utilization',
    label: 'Utilization',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      (a.utilization || 0) - (b.utilization || 0),
    width: 10,
  },
  {
    field: 'accuracy',
    label: 'Accuracy',
    sortable: (a: HardwareConfiguration, b: HardwareConfiguration): number =>
      (a.accuracy || 0) - (b.accuracy || 0),
    width: 10,
  },
];

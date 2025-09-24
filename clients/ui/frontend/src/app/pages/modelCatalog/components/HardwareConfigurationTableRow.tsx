import * as React from 'react';
import { Td, Tr } from '@patternfly/react-table';
import { HardwareConfiguration } from '~/app/pages/modelCatalog/types/hardwareConfiguration';
import {
  formatValue,
  formatHardwareConfiguration,
} from '~/app/pages/modelCatalog/utils/hardwareConfigurationUtils';

type HardwareConfigurationTableRowProps = {
  configuration: HardwareConfiguration;
};

const HardwareConfigurationTableRow = ({
  configuration,
}: HardwareConfigurationTableRowProps): React.JSX.Element => (
  <Tr>
    <Td dataLabel="Hardware configuration">
      {formatHardwareConfiguration(configuration.hardwareCount, configuration.hardwareType)}
    </Td>
    <Td dataLabel="Latency">{formatValue(configuration.latency, 'ms')}</Td>
    <Td dataLabel="Throughput">{formatValue(configuration.throughput, 'req/s')}</Td>
    <Td dataLabel="TPS">{formatValue(configuration.tps, 'tps')}</Td>
    <Td dataLabel="GuideLLM version">{configuration.guideLLMVersion || '-'}</Td>
    <Td dataLabel="RHAIIS version">{configuration.rhaiisVersion || '-'}</Td>
    <Td dataLabel="Framework">{configuration.framework || '-'}</Td>
    <Td dataLabel="Precision">{configuration.precision || '-'}</Td>
    <Td dataLabel="Batch size">{configuration.batchSize || '-'}</Td>
    <Td dataLabel="Memory usage">{formatValue(configuration.memoryUsage, 'GB')}</Td>
    <Td dataLabel="Utilization">{formatValue(configuration.utilization, '%')}</Td>
    <Td dataLabel="Accuracy">{formatValue(configuration.accuracy, '%')}</Td>
  </Tr>
);

export default HardwareConfigurationTableRow;

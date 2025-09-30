import * as React from 'react';
import { Td, Tr } from '@patternfly/react-table';
import { CatalogPerformanceMetricsArtifact } from '~/app/modelCatalogTypes';
import {
  formatLatency,
  formatTokenValue,
  getDoubleValue,
  getHardwareConfiguration,
  getIntValue,
  getStringValue,
  getTotalRps,
} from '~/app/pages/modelCatalog/utils/performanceMetricsUtils';
import { hardwareConfigColumns } from './HardwareConfigurationTableColumns';

type HardwareConfigurationTableRowProps = {
  configuration: CatalogPerformanceMetricsArtifact;
};

const HardwareConfigurationTableRow = ({
  configuration,
}: HardwareConfigurationTableRowProps): React.JSX.Element => {
  const getCellValue = (field: string): string | number => {
    const { customProperties } = configuration;

    switch (field) {
      case 'hardware':
        return getHardwareConfiguration(configuration);
      case 'hardware_count':
        return getIntValue(customProperties, 'hardware_count');
      case 'requests_per_second':
        return getDoubleValue(customProperties, 'requests_per_second');
      case 'total_rps':
        return getTotalRps(customProperties);
      case 'ttft_mean':
      case 'ttft_p90':
      case 'ttft_p95':
      case 'ttft_p99':
      case 'e2e_mean':
      case 'e2e_p90':
      case 'e2e_p95':
      case 'e2e_p99':
      case 'tps_mean':
      case 'tps_p90':
      case 'tps_p95':
      case 'tps_p99':
      case 'itl_mean':
      case 'itl_p90':
      case 'itl_p95':
      case 'itl_p99':
        return formatLatency(getDoubleValue(customProperties, field));
      case 'max_input_tokens':
      case 'max_output_tokens':
      case 'mean_input_tokens':
      case 'mean_output_tokens':
        return formatTokenValue(getDoubleValue(customProperties, field));
      case 'framework_version':
        return getStringValue(customProperties, field);
      default:
        return '-';
    }
  };

  return (
    <Tr>
      {hardwareConfigColumns.map((column, index) => (
        <Td
          key={column.field}
          dataLabel={column.label.replace('\n', ' ')}
          isStickyColumn={index < 2}
          stickyMinWidth={index < 2 ? `${column.width}ch` : undefined}
          modifier="fitContent"
        >
          {getCellValue(column.field)}
        </Td>
      ))}
    </Tr>
  );
};

export default HardwareConfigurationTableRow;

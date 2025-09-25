import * as React from 'react';
import { Td, Tr } from '@patternfly/react-table';
import { HardwareConfiguration } from '~/app/pages/modelCatalog/types/hardwareConfiguration';
import { hardwareConfigColumns } from './HardwareConfigurationTableColumns';

type HardwareConfigurationTableRowProps = {
  configuration: HardwareConfiguration;
};

const HardwareConfigurationTableRow = ({
  configuration,
}: HardwareConfigurationTableRowProps): React.JSX.Element => (
  <Tr>
    {hardwareConfigColumns.map((column, index) => {
      const getCellValue = () => {
        // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
        const field = column.field as keyof HardwareConfiguration;
        const value = configuration[field];
        if (column.field.includes('Latency') && typeof value === 'number') {
          return `${value} ms`;
        }
        return value;
      };

      return (
        <Td
          key={column.field}
          dataLabel={column.label.replace('\n', ' ')}
          className={
            index < 2 ? 'pf-v6-c-table__sticky-column pf-v6-c-table__sticky-column--left' : ''
          }
          style={{
            width: `${column.width}ch`,
            minWidth: `${column.width}ch`,
            textAlign: 'left',
          }}
        >
          {getCellValue()}
        </Td>
      );
    })}
  </Tr>
);

export default HardwareConfigurationTableRow;

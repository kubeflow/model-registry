/* eslint-disable @typescript-eslint/consistent-type-assertions */
import * as React from 'react';
import { Table, Thead, Tbody, Tr, Th } from '@patternfly/react-table';
import { Spinner } from '@patternfly/react-core';
import {
  HardwareConfiguration,
  HardwareConfigurationTableProps,
} from '~/app/pages/modelCatalog/types/hardwareConfiguration';
import { hardwareConfigColumns } from './HardwareConfigurationTableColumns';
import HardwareConfigurationTableRow from './HardwareConfigurationTableRow';

const HardwareConfigurationTable = ({
  configurations,
  isLoading = false,
}: HardwareConfigurationTableProps): React.JSX.Element => {
  const [sortBy, setSortBy] = React.useState<{ index: number; direction: 'asc' | 'desc' }>({
    index: 0,
    direction: 'asc',
  });
  const [sortedConfigurations, setSortedConfigurations] = React.useState(configurations);

  React.useEffect(() => {
    const column = hardwareConfigColumns[sortBy.index];
    // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
    if (column && column.sortable && typeof column.sortable === 'function') {
      const sorted = [...configurations].toSorted((a, b) => {
        const result = (
          column.sortable as (
            a: HardwareConfiguration,
            b: HardwareConfiguration,
            keyField: string,
          ) => number
        )(a, b, column.field);
        return sortBy.direction === 'asc' ? result : -result;
      });
      setSortedConfigurations(sorted);
    } else {
      setSortedConfigurations(configurations);
    }
  }, [configurations, sortBy]);

  const onSort = React.useCallback(
    (_event: React.MouseEvent, index: number, direction: 'asc' | 'desc') => {
      setSortBy({ index, direction });
    },
    [],
  );

  if (isLoading) {
    return <Spinner size="lg" />;
  }

  return (
    <div className="hardware-table-container">
      <style>
        {`
          .hardware-table-container {
            width: 100%;
            max-width: 100%;
            overflow-x: auto;
          }
          .hardware-table-container table {
            width: max-content;
            min-width: 100%;
            table-layout: auto;
          }
        `}
      </style>
      <Table data-testid="hardware-configuration-table" variant="compact" isStickyHeader>
        <Thead>
          <Tr>
            {hardwareConfigColumns.map((column, index) => {
              const renderLabel = () => {
                if (column.label.includes('\n')) {
                  const parts = column.label.split('\n');
                  return (
                    <>
                      {parts[0]}
                      <br />
                      {parts[1]}
                    </>
                  );
                }
                return column.label;
              };

              return (
                <Th
                  key={column.field}
                  className={
                    index < 2
                      ? 'pf-v6-c-table__sticky-column pf-v6-c-table__sticky-column--left'
                      : ''
                  }
                  style={{
                    width: `${column.width}ch`,
                    minWidth: `${column.width}ch`,
                    verticalAlign: 'middle',
                    height: 'auto',
                    lineHeight: '1.2',
                    textAlign: 'left',
                  }}
                  sort={column.sortable ? { onSort, columnIndex: index, sortBy } : undefined}
                >
                  {renderLabel()}
                </Th>
              );
            })}
          </Tr>
        </Thead>
        <Tbody>
          {sortedConfigurations.map((config) => (
            <HardwareConfigurationTableRow key={config.id} configuration={config} />
          ))}
        </Tbody>
      </Table>
    </div>
  );
};

export default HardwareConfigurationTable;

import * as React from 'react';
import { TbodyProps } from '@patternfly/react-table';
import { EitherNotBoth } from '~/shared/typeHelpers';
import TableBase, { MIN_PAGE_SIZE } from './TableBase';
import useTableColumnSort from './useTableColumnSort';

type TableProps<DataType> = Omit<
  React.ComponentProps<typeof TableBase<DataType>>,
  'itemCount' | 'onPerPageSelect' | 'onSetPage' | 'page' | 'perPage'
> &
  EitherNotBoth<
    { disableRowRenderSupport?: boolean },
    { tbodyProps?: TbodyProps & { ref?: React.Ref<HTMLTableSectionElement> } }
  >;

const Table = <T,>({
  data: allData,
  columns,
  subColumns,
  enablePagination,
  defaultSortColumn = 0,
  truncateRenderingAt = 0,
  ...props
}: TableProps<T>): React.ReactElement => {
  const [page, setPage] = React.useState(1);
  const [pageSize, setPageSize] = React.useState(MIN_PAGE_SIZE);
  const sort = useTableColumnSort<T>(columns, subColumns || [], defaultSortColumn);
  const sortedData = sort.transformData(allData);

  let data: T[];
  if (truncateRenderingAt) {
    data = sortedData.slice(0, truncateRenderingAt);
  } else if (enablePagination) {
    data = sortedData.slice(pageSize * (page - 1), pageSize * page);
  } else {
    data = sortedData;
  }

  // update page to 1 if data changes (common when filter is applied)
  React.useEffect(() => {
    if (data.length === 0) {
      setPage(1);
    }
  }, [data.length]);

  return (
    <TableBase
      data={data}
      columns={columns}
      subColumns={subColumns}
      enablePagination={enablePagination}
      itemCount={allData.length}
      perPage={pageSize}
      page={page}
      onSetPage={(e, newPage) => setPage(newPage)}
      onPerPageSelect={(e, newSize, newPage) => {
        setPageSize(newSize);
        setPage(newPage);
      }}
      getColumnSort={sort.getColumnSort}
      {...props}
    />
  );
};

export default Table;

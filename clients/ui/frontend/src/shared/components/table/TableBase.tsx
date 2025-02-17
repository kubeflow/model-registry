import * as React from 'react';
import {
  Pagination,
  PaginationProps,
  Skeleton,
  Toolbar,
  ToolbarContent,
  ToolbarGroup,
  ToolbarItem,
  Tooltip,
} from '@patternfly/react-core';
import {
  Table,
  Thead,
  Tr,
  Th,
  TableProps,
  Caption,
  Tbody,
  Td,
  TbodyProps,
  InnerScrollContainer,
  TrProps,
} from '@patternfly/react-table';
import { EitherNotBoth } from '~/shared/typeHelpers';
import { GetColumnSort, SortableData } from './types';
import { CHECKBOX_FIELD_ID, EXPAND_FIELD_ID, KEBAB_FIELD_ID } from './const';

type Props<DataType> = {
  loading?: boolean;
  skeletonRowCount?: number;
  skeletonRowProps?: TrProps;
  data: DataType[];
  columns: SortableData<DataType>[];
  subColumns?: SortableData<DataType>[];
  hasNestedHeader?: boolean;
  defaultSortColumn?: number;
  rowRenderer: (data: DataType, rowIndex: number) => React.ReactNode;
  enablePagination?: boolean | 'compact';
  truncateRenderingAt?: number;
  toolbarContent?: React.ReactElement<typeof ToolbarItem | typeof ToolbarGroup>;
  onClearFilters?: () => void;
  bottomToolbarContent?: React.ReactElement<typeof ToolbarItem | typeof ToolbarGroup>;
  emptyTableView?: React.ReactNode;
  caption?: string;
  footerRow?: (pageNumber: number) => React.ReactElement<typeof Tr> | null;
  selectAll?: {
    onSelect: (value: boolean) => void;
    selected: boolean;
    disabled?: boolean;
    tooltip?: string;
  };
  getColumnSort?: GetColumnSort;
  disableItemCount?: boolean;
  hasStickyColumns?: boolean;
} & EitherNotBoth<
  { disableRowRenderSupport?: boolean },
  { tbodyProps?: TbodyProps & { ref?: React.Ref<HTMLTableSectionElement> } }
> &
  Omit<TableProps, 'ref' | 'data'> &
  Pick<
    PaginationProps,
    | 'itemCount'
    | 'onPerPageSelect'
    | 'onSetPage'
    | 'page'
    | 'perPage'
    | 'perPageOptions'
    | 'toggleTemplate'
    | 'onNextClick'
    | 'onPreviousClick'
  >;

export const MIN_PAGE_SIZE = 10;

const defaultPerPageOptions = [
  {
    title: '10',
    value: 10,
  },
  {
    title: '20',
    value: 20,
  },
  {
    title: '30',
    value: 30,
  },
];

const TableBase = <T,>({
  data,
  columns,
  subColumns,
  hasNestedHeader,
  rowRenderer,
  enablePagination,
  toolbarContent,
  onClearFilters,
  bottomToolbarContent,
  emptyTableView,
  caption,
  disableRowRenderSupport,
  selectAll,
  footerRow,
  tbodyProps,
  perPage = 10,
  page = 1,
  perPageOptions = defaultPerPageOptions,
  onSetPage,
  onNextClick,
  onPreviousClick,
  onPerPageSelect,
  getColumnSort,
  itemCount = 0,
  loading,
  skeletonRowCount,
  skeletonRowProps = {},
  toggleTemplate,
  disableItemCount = false,
  hasStickyColumns,
  ...props
}: Props<T>): React.ReactElement => {
  const selectAllRef = React.useRef(null);
  const showPagination = enablePagination;

  const pagination = (variant: 'top' | 'bottom') => (
    <Pagination
      isCompact
      {...(!disableItemCount && { itemCount })}
      perPage={perPage}
      page={page}
      onSetPage={onSetPage}
      onNextClick={onNextClick}
      onPreviousClick={onPreviousClick}
      onPerPageSelect={onPerPageSelect}
      toggleTemplate={toggleTemplate}
      variant={variant}
      widgetId="table-pagination"
      perPageOptions={perPageOptions}
      menuAppendTo="inline"
      titles={{
        paginationAriaLabel: `${variant} pagination`,
      }}
    />
  );

  // Use a reference to store the heights of table rows once loaded
  const tableRef = React.useRef<HTMLTableElement>(null);
  const rowHeightsRef = React.useRef<number[] | undefined>();
  React.useLayoutEffect(() => {
    if (!loading || rowHeightsRef.current == null) {
      const heights: number[] = [];
      const rows = tableRef.current?.querySelectorAll<HTMLTableRowElement>(':scope > tbody > tr');
      rows?.forEach((r) => heights.push(r.offsetHeight));
      rowHeightsRef.current = heights;
    }
  }, [loading]);

  const renderColumnHeader = (col: SortableData<T>, i: number, isSubheader?: boolean) => {
    if (col.field === CHECKBOX_FIELD_ID && selectAll) {
      return (
        <React.Fragment key={`checkbox-${i}`}>
          <Tooltip
            key="select-all-checkbox"
            content={selectAll.tooltip ?? 'Select all page items'}
            triggerRef={selectAllRef}
          />
          <Th
            ref={selectAllRef}
            colSpan={col.colSpan}
            rowSpan={col.rowSpan}
            isStickyColumn={col.isStickyColumn}
            stickyMinWidth={col.stickyMinWidth}
            select={{
              isSelected: selectAll.selected,
              onSelect: (e, value) => selectAll.onSelect(value),
              isDisabled: selectAll.disabled,
            }}
            // TODO: Log PF bug -- when there are no rows this gets truncated
            style={{ minWidth: '45px' }}
            isSubheader={isSubheader}
            aria-label="Select all"
          />
        </React.Fragment>
      );
    }

    return col.label ? (
      <Th
        key={col.field + i}
        colSpan={col.colSpan}
        rowSpan={col.rowSpan}
        sort={getColumnSort && col.sortable ? getColumnSort(i) : undefined}
        width={col.width}
        info={col.info}
        isSubheader={isSubheader}
        isStickyColumn={col.isStickyColumn}
        stickyMinWidth={col.stickyMinWidth}
        stickyLeftOffset={col.stickyLeftOffset}
        hasRightBorder={col.hasRightBorder}
        modifier={col.modifier}
        visibility={col.visibility}
        className={col.className}
      >
        {col.label}
      </Th>
    ) : (
      // Table headers cannot be empty for a11y, table cells can -- https://dequeuniversity.com/rules/axe/4.0/empty-table-header
      <Td key={col.field + i} width={col.width} />
    );
  };

  const renderRows = () =>
    loading
      ? // compute the number of items in the upcoming page
        new Array(
          itemCount === 0
            ? skeletonRowCount || rowHeightsRef.current?.length || MIN_PAGE_SIZE
            : Math.max(0, Math.min(perPage, itemCount - perPage * (page - 1))),
        )
          .fill(undefined)
          .map((_, i) => {
            // Set the height to the last known row height or otherwise the same height as the first row.
            // When going to a previous page, the number of rows may be greater than the current.

            const getRow = () => (
              <Tr
                key={`skeleton-${i}`}
                {...skeletonRowProps}
                style={{
                  ...(skeletonRowProps.style || {}),
                  height: rowHeightsRef.current?.[i] || rowHeightsRef.current?.[0],
                }}
              >
                {columns.map((col) => (
                  <Td
                    key={col.field}
                    // assign classes to reserve space
                    className={
                      col.field === CHECKBOX_FIELD_ID || col.field === EXPAND_FIELD_ID
                        ? 'pf-c-table__toggle'
                        : col.field === KEBAB_FIELD_ID
                          ? 'pf-c-table__action'
                          : undefined
                    }
                  >
                    {
                      // render placeholders to reserve space
                      col.field === EXPAND_FIELD_ID || col.field === KEBAB_FIELD_ID ? (
                        <div style={{ width: 46 }} />
                      ) : col.field === CHECKBOX_FIELD_ID ? (
                        <div style={{ width: 13 }} />
                      ) : (
                        <Skeleton width="50%" />
                      )
                    }
                  </Td>
                ))}
              </Tr>
            );
            return disableRowRenderSupport ? (
              <Tbody key={`skeleton-tbody-${i}`}>{getRow()}</Tbody>
            ) : (
              getRow()
            );
          })
      : data.map((row, rowIndex) => rowRenderer(row, rowIndex));

  const table = (
    <Table {...props} {...(hasStickyColumns && { gridBreakPoint: '' })} ref={tableRef}>
      {caption && <Caption>{caption}</Caption>}
      <Thead noWrap hasNestedHeader={hasNestedHeader}>
        {/* Note from PF: following custom style can be removed when we can resolve misalignment issue natively */}
        <Tr>{columns.map((col, i) => renderColumnHeader(col, i))}</Tr>
        {subColumns?.length ? (
          <Tr>{subColumns.map((col, i) => renderColumnHeader(col, columns.length + i, true))}</Tr>
        ) : null}
      </Thead>
      {disableRowRenderSupport ? renderRows() : <Tbody {...tbodyProps}>{renderRows()}</Tbody>}
      {footerRow && footerRow(page)}
    </Table>
  );

  return (
    <>
      {(toolbarContent || showPagination) && (
        <Toolbar
          inset={{ default: 'insetNone' }}
          className="pf-v6-u-w-100"
          customLabelGroupContent={onClearFilters ? undefined : <></>}
          clearAllFilters={onClearFilters}
        >
          <ToolbarContent>
            {toolbarContent}
            {showPagination && (
              <ToolbarItem
                variant="pagination"
                align={{ default: 'alignEnd' }}
                className="pf-v6-u-pr-lg"
              >
                {pagination('top')}
              </ToolbarItem>
            )}
          </ToolbarContent>
        </Toolbar>
      )}

      {hasStickyColumns ? <InnerScrollContainer>{table}</InnerScrollContainer> : table}

      {!loading && emptyTableView && data.length === 0 && (
        <div style={{ padding: 'var(--pf-global--spacer--2xl) 0', textAlign: 'center' }}>
          {emptyTableView}
        </div>
      )}

      {bottomToolbarContent && (
        <Toolbar inset={{ default: 'insetNone' }} className="pf-v6-u-w-100">
          <ToolbarContent alignItems="center">{bottomToolbarContent}</ToolbarContent>
        </Toolbar>
      )}
      {showPagination && <>{pagination('bottom')}</>}
    </>
  );
};

export default TableBase;

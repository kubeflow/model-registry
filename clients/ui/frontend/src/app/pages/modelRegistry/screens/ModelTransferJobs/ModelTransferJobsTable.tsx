import * as React from 'react';
import { Table, DashboardEmptyTableView } from 'mod-arch-shared';
import { ModelTransferJob } from '~/app/types';
import { modelTransferJobsColumns } from './ModelTransferJobsTableColumns';
import ModelTransferJobTableRow from './ModelTransferJobTableRow';

type ModelTransferJobsTableProps = {
  jobs: ModelTransferJob[];
  clearFilters: () => void;
  toolbarContent?: React.ComponentProps<typeof Table>['toolbarContent'];
  onRequestDelete?: (job: ModelTransferJob) => void;
};

const ModelTransferJobsTable: React.FC<ModelTransferJobsTableProps> = ({
  jobs,
  clearFilters,
  toolbarContent,
  onRequestDelete,
}) => {
  const defaultSortColumnIndex = React.useMemo(
    () => modelTransferJobsColumns.findIndex((col) => col.field === 'created'),
    [],
  );

  const sortedJobs = React.useMemo(() => {
    const createdColumn = modelTransferJobsColumns.find((col) => col.field === 'created');
    if (createdColumn?.sortable && typeof createdColumn.sortable === 'function') {
      const sortFn = createdColumn.sortable;
      return [...jobs].toSorted((a, b) => sortFn(a, b, 'created'));
    }
    return jobs;
  }, [jobs]);

  return (
    <Table
      data-testid="model-transfer-jobs-table"
      data={sortedJobs}
      columns={modelTransferJobsColumns}
      toolbarContent={toolbarContent}
      defaultSortColumn={defaultSortColumnIndex}
      onClearFilters={clearFilters}
      enablePagination
      emptyTableView={<DashboardEmptyTableView onClearFilters={clearFilters} />}
      rowRenderer={(job) => (
        <ModelTransferJobTableRow key={job.id} job={job} onRequestDelete={onRequestDelete} />
      )}
    />
  );
};

export default ModelTransferJobsTable;

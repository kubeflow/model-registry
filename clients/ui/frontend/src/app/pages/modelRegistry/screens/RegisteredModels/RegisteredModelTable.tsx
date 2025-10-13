import * as React from 'react';
import { Table, DashboardEmptyTableView } from 'mod-arch-shared';
import { ModelVersion, RegisteredModel } from '~/app/types';
import { getLatestVersionForRegisteredModel } from '~/app/pages/modelRegistry/screens/utils';
import { rmColumns } from './RegisteredModelsTableColumns';
import RegisteredModelTableRow from './RegisteredModelTableRow';

type RegisteredModelTableProps = {
  clearFilters: () => void;
  registeredModels: RegisteredModel[];
  modelVersions: ModelVersion[];
  refresh: () => void;
} & Partial<Pick<React.ComponentProps<typeof Table>, 'toolbarContent'>>;

const RegisteredModelTable: React.FC<RegisteredModelTableProps> = ({
  clearFilters,
  registeredModels,
  modelVersions,
  toolbarContent,
  refresh,
}) => {
  const { defaultSortColumnIndex } = React.useMemo(() => {
    const columns = [...rmColumns];

    // Find the index of the "last_modified" column after any insertions
    const lastModifiedIndex = columns.findIndex((col) => col.field === 'last_modified');

    return {
      extendedColumns: columns,
      defaultSortColumnIndex: lastModifiedIndex,
    };
  }, []);
  // Pre-sort the data by last modified to ensure initial sort
  const sortedRegisteredModels = React.useMemo(() => {
    const lastModifiedColumn = rmColumns.find((col) => col.field === 'last_modified');
    if (lastModifiedColumn?.sortable && typeof lastModifiedColumn.sortable === 'function') {
      const sortFn = lastModifiedColumn.sortable;
      return [...registeredModels].toSorted((a, b) => sortFn(a, b, 'last_modified'));
    }
    return registeredModels;
  }, [registeredModels]);

  return (
    <Table
      data-testid="registered-model-table"
      data={sortedRegisteredModels}
      columns={rmColumns}
      toolbarContent={toolbarContent}
      defaultSortColumn={defaultSortColumnIndex}
      onClearFilters={clearFilters}
      enablePagination
      emptyTableView={<DashboardEmptyTableView onClearFilters={clearFilters} />}
      rowRenderer={(rm: RegisteredModel) => (
        <RegisteredModelTableRow
          key={rm.name}
          hasDeploys={false}
          registeredModel={rm}
          latestModelVersion={getLatestVersionForRegisteredModel(modelVersions, rm.id)}
          refresh={refresh}
        />
      )}
    />
  );
};

export default RegisteredModelTable;

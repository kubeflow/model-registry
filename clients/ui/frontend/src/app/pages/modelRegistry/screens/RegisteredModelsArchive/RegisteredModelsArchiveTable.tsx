import * as React from 'react';
import { Table } from '~/shared/components/table';
import { RegisteredModel } from '~/app/types';
import { rmColumns } from '~/app/pages/modelRegistry/screens/RegisteredModels/RegisteredModelsTableColumns';
import DashboardEmptyTableView from '~/shared/components/DashboardEmptyTableView';
import RegisteredModelTableRow from '~/app/pages/modelRegistry/screens/RegisteredModels/RegisteredModelTableRow';

type RegisteredModelsArchiveTableProps = {
  clearFilters: () => void;
  registeredModels: RegisteredModel[];
  refresh: () => void;
} & Partial<Pick<React.ComponentProps<typeof Table>, 'toolbarContent'>>;

const RegisteredModelsArchiveTable: React.FC<RegisteredModelsArchiveTableProps> = ({
  clearFilters,
  registeredModels,
  toolbarContent,
  refresh,
}) => (
  <Table
    data-testid="registered-models-archive-table"
    data={registeredModels}
    columns={rmColumns}
    toolbarContent={toolbarContent}
    defaultSortColumn={2}
    onClearFilters={clearFilters}
    enablePagination
    emptyTableView={<DashboardEmptyTableView onClearFilters={clearFilters} />}
    rowRenderer={(rm) => (
      <RegisteredModelTableRow key={rm.name} registeredModel={rm} isArchiveRow refresh={refresh} />
    )}
  />
);

export default RegisteredModelsArchiveTable;

import * as React from 'react';
import { Table } from '~/shared/components/table';
import { RegisteredModel } from '~/app/types';
import DashboardEmptyTableView from '~/shared/components/DashboardEmptyTableView';
import { rmColumns } from './RegisteredModelsTableColumns';
import RegisteredModelTableRow from './RegisteredModelTableRow';

type RegisteredModelTableProps = {
  clearFilters: () => void;
  registeredModels: RegisteredModel[];
  refresh: () => void;
} & Partial<Pick<React.ComponentProps<typeof Table>, 'toolbarContent'>>;

const RegisteredModelTable: React.FC<RegisteredModelTableProps> = ({
  clearFilters,
  registeredModels,
  toolbarContent,
  refresh,
}) => (
  <Table
    data-testid="registered-model-table"
    data={registeredModels}
    columns={rmColumns}
    toolbarContent={toolbarContent}
    defaultSortColumn={2}
    enablePagination="compact"
    emptyTableView={<DashboardEmptyTableView onClearFilters={clearFilters} />}
    rowRenderer={(rm) => (
      <RegisteredModelTableRow key={rm.name} registeredModel={rm} refresh={refresh} />
    )}
  />
);

export default RegisteredModelTable;

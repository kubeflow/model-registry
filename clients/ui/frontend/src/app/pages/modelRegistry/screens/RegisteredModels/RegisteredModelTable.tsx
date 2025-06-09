import * as React from 'react';
import { Table, DashboardEmptyTableView } from 'mod-arch-shared';
import { RegisteredModel } from '~/app/types';
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
    onClearFilters={clearFilters}
    enablePagination
    emptyTableView={<DashboardEmptyTableView onClearFilters={clearFilters} />}
    rowRenderer={(rm: RegisteredModel) => (
      <RegisteredModelTableRow
        key={rm.name}
        hasDeploys={false}
        registeredModel={rm}
        refresh={refresh}
      />
    )}
  />
);

export default RegisteredModelTable;

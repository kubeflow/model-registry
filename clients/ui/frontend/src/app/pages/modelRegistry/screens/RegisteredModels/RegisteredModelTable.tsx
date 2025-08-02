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
        latestModelVersion={getLatestVersionForRegisteredModel(modelVersions, rm.id)}
        refresh={refresh}
      />
    )}
  />
);

export default RegisteredModelTable;

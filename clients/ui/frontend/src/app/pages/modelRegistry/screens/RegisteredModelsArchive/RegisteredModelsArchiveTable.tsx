import * as React from 'react';
import { Table, DashboardEmptyTableView } from 'mod-arch-shared';
import { ModelVersion, RegisteredModel } from '~/app/types';
import { rmColumns } from '~/app/pages/modelRegistry/screens/RegisteredModels/RegisteredModelsTableColumns';
import RegisteredModelTableRow from '~/app/pages/modelRegistry/screens/RegisteredModels/RegisteredModelTableRow';
import { getLatestVersionForRegisteredModel } from '~/app/pages/modelRegistry/screens/utils';

type RegisteredModelsArchiveTableProps = {
  clearFilters: () => void;
  registeredModels: RegisteredModel[];
  modelVersions: ModelVersion[];
  refresh: () => void;
} & Partial<Pick<React.ComponentProps<typeof Table>, 'toolbarContent'>>;

const RegisteredModelsArchiveTable: React.FC<RegisteredModelsArchiveTableProps> = ({
  clearFilters,
  registeredModels,
  modelVersions,
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
    rowRenderer={(rm: RegisteredModel) => (
      <RegisteredModelTableRow
        key={rm.name}
        registeredModel={rm}
        latestModelVersion={getLatestVersionForRegisteredModel(modelVersions, rm.id)}
        isArchiveRow
        refresh={refresh}
      />
    )}
  />
);

export default RegisteredModelsArchiveTable;

import * as React from 'react';
import { DashboardEmptyTableView, Table } from 'mod-arch-shared';
import { ModelVersion } from '~/app/types';
import { mvColumns } from '~/app/pages/modelRegistry/screens/ModelVersions/ModelVersionsTableColumns';
import ModelVersionsTableRow from '~/app/pages/modelRegistry/screens/ModelVersions/ModelVersionsTableRow';

type ModelVersionsTableProps = {
  clearFilters: () => void;
  modelVersions: ModelVersion[];
  isArchiveModel?: boolean;
  refresh: () => void;
} & Partial<Pick<React.ComponentProps<typeof Table>, 'toolbarContent'>>;

const ModelVersionsTable: React.FC<ModelVersionsTableProps> = ({
  clearFilters,
  modelVersions,
  toolbarContent,
  isArchiveModel,
  refresh,
}) => (
  <Table
    data-testid="model-versions-table"
    data={modelVersions}
    columns={mvColumns}
    toolbarContent={toolbarContent}
    defaultSortColumn={3}
    enablePagination
    onClearFilters={clearFilters}
    emptyTableView={<DashboardEmptyTableView onClearFilters={clearFilters} />}
    rowRenderer={(mv: ModelVersion) => (
      <ModelVersionsTableRow
        hasDeployment={false}
        key={mv.name}
        modelVersion={mv}
        isArchiveModel={isArchiveModel}
        refresh={refresh}
      />
    )}
  />
);

export default ModelVersionsTable;

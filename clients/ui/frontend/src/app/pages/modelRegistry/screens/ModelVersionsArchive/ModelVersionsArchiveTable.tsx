import * as React from 'react';
import { Table, DashboardEmptyTableView } from 'mod-arch-shared';
import { ModelVersion } from '~/app/types';
import ModelVersionsTableRow from '~/app/pages/modelRegistry/screens/ModelVersions/ModelVersionsTableRow';
import { mvColumns } from '~/app/pages/modelRegistry/screens/ModelVersions/ModelVersionsTableColumns';

type ModelVersionsArchiveTableProps = {
  clearFilters: () => void;
  modelVersions: ModelVersion[];
  refresh: () => void;
} & Partial<Pick<React.ComponentProps<typeof Table>, 'toolbarContent'>>;

const ModelVersionsArchiveTable: React.FC<ModelVersionsArchiveTableProps> = ({
  clearFilters,
  modelVersions,
  toolbarContent,
  refresh,
}) => (
  <Table
    data-testid="model-versions-archive-table"
    data={modelVersions}
    columns={mvColumns}
    toolbarContent={toolbarContent}
    enablePagination
    onClearFilters={clearFilters}
    emptyTableView={<DashboardEmptyTableView onClearFilters={clearFilters} />}
    defaultSortColumn={1}
    rowRenderer={(mv: ModelVersion) => (
      <ModelVersionsTableRow key={mv.name} modelVersion={mv} isArchiveRow refresh={refresh} />
    )}
  />
);

export default ModelVersionsArchiveTable;

import * as React from 'react';
import { Table } from '~/app/components/table';
import { ModelVersion } from '~/app/types';
import { mvColumns } from '~/app/pages/modelRegistry/screens/ModelVersions/ModelVersionsTableColumns';
import DashboardEmptyTableView from '~/app/components/DashboardEmptyTableView';
import ModelVersionsTableRow from '~/app/pages/modelRegistry/screens/ModelVersions/ModelVersionsTableRow';

type ModelVersionsTableProps = {
  clearFilters: () => void;
  modelVersions: ModelVersion[];
  refresh: () => void;
} & Partial<Pick<React.ComponentProps<typeof Table>, 'toolbarContent'>>;

const ModelVersionsTable: React.FC<ModelVersionsTableProps> = ({
  clearFilters,
  modelVersions,
  toolbarContent,
  refresh,
}) => (
  <Table
    data-testid="model-versions-table"
    data={modelVersions}
    columns={mvColumns}
    toolbarContent={toolbarContent}
    defaultSortColumn={3}
    enablePagination
    emptyTableView={<DashboardEmptyTableView onClearFilters={clearFilters} />}
    rowRenderer={(mv) => (
      <ModelVersionsTableRow key={mv.name} modelVersion={mv} refresh={refresh} />
    )}
  />
);

export default ModelVersionsTable;

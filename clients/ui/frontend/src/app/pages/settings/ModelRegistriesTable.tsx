import React from 'react';
import { ModelRegistry } from '~/app/types';
import { Table } from '~/shared/components/table';
import { modelRegistryColumns } from './columns';
import ModelRegistriesTableRow from './ModelRegistriesTableRow';

type ModelRegistriesTableProps = {
  modelRegistries: ModelRegistry[];
};

const ModelRegistriesTable: React.FC<ModelRegistriesTableProps> = ({ modelRegistries }) => (
  // TODO: [Midstream] Complete once we have permissions
  <Table
    data-testid="model-registries-table"
    data={modelRegistries}
    columns={modelRegistryColumns}
    rowRenderer={(mr) => (
      <ModelRegistriesTableRow
        key={mr.name}
        modelRegistry={mr}
        // eslint-disable-next-line @typescript-eslint/no-empty-function
        onDeleteRegistry={() => {}}
        // eslint-disable-next-line @typescript-eslint/no-empty-function
        onEditRegistry={() => {}}
      />
    )}
    variant="compact"
  />
);

export default ModelRegistriesTable;

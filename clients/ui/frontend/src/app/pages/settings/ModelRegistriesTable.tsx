import React from 'react';
import { Button, Toolbar, ToolbarContent, ToolbarItem } from '@patternfly/react-core';
import { ModelRegistry } from '~/app/types';
import { Table } from '~/shared/components/table';
import { modelRegistryColumns } from './columns';
import ModelRegistriesTableRow from './ModelRegistriesTableRow';

type ModelRegistriesTableProps = {
  modelRegistries: ModelRegistry[];
  onCreateModelRegistryClick: () => void;
};

const ModelRegistriesTable: React.FC<ModelRegistriesTableProps> = ({
  modelRegistries,
  onCreateModelRegistryClick,
}) => (
  // TODO: [Midstream] Complete once we have permissions
  <Table
    data-testid="model-registries-table"
    data={modelRegistries}
    columns={modelRegistryColumns}
    toolbarContent={
      <Toolbar>
        <ToolbarContent>
          <ToolbarItem>
            <Button
              data-testid="create-model-registry-button"
              variant="primary"
              onClick={onCreateModelRegistryClick}
            >
              Create model registry
            </Button>
          </ToolbarItem>
        </ToolbarContent>
      </Toolbar>
    }
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

import React from 'react';
import { Button, Toolbar, ToolbarContent, ToolbarItem } from '@patternfly/react-core';
import { Table, Thead, Tbody, Tr, Th } from '@patternfly/react-table';
import { ModelRegistryKind } from 'mod-arch-shared';
import CreateModal from '~/app/pages/modelRegistrySettings/CreateModal';
import ModelRegistriesTableRow from './ModelRegistriesTableRow';
import DeleteModelRegistryModal from './DeleteModelRegistryModal';

type ModelRegistriesTableProps = {
  modelRegistries: ModelRegistryKind[];
  onCreateModelRegistryClick: () => void;
  refresh: () => void;
};

const ModelRegistriesTable: React.FC<ModelRegistriesTableProps> = ({
  modelRegistries,
  onCreateModelRegistryClick,
}) => {
  // TODO: [Midstream] Complete once we have permissions

  const [deleteRegistry, setDeleteRegistry] = React.useState<ModelRegistryKind>();
  const [editRegistry, setEditRegistry] = React.useState<ModelRegistryKind>();

  const columns = ['Name', 'Owner', 'Created', ''];

  return (
    <>
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
      <Table aria-label="Model Registries Table" variant="compact">
        <Thead>
          <Tr>
            {columns.map((column, index) => (
              <Th key={index}>{column}</Th>
            ))}
          </Tr>
        </Thead>
        <Tbody>
          {modelRegistries.map((mr) => (
            <ModelRegistriesTableRow
              key={mr.metadata.name}
              modelRegistry={mr}
              onDeleteRegistry={setDeleteRegistry}
              onEditRegistry={setEditRegistry}
              roleBindings={{
                data: [],
                loaded: true,
                error: undefined,
                refresh: () => Promise.resolve([]),
              }}
            />
          ))}
        </Tbody>
      </Table>
      {editRegistry ? (
        <CreateModal
          modelRegistry={editRegistry}
          onClose={() => setEditRegistry(undefined)}
          refresh={() => Promise.resolve(undefined)}
        />
      ) : null}
      {deleteRegistry ? (
        <DeleteModelRegistryModal
          modelRegistry={deleteRegistry}
          onClose={() => setDeleteRegistry(undefined)}
          refresh={() => Promise.resolve(undefined)}
        />
      ) : null}
    </>
  );
};

export default ModelRegistriesTable;

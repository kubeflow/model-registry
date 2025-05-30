import React from 'react';
import { Button, Toolbar, ToolbarContent, ToolbarItem } from '@patternfly/react-core';
import { ModelRegistryKind, Table } from 'mod-arch-shared';
import { modelRegistryColumns } from './columns';
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
  refresh,
}) => {
  // TODO: [Midstream] Complete once we have permissions

  const [deleteRegistry, setDeleteRegistry] = React.useState<ModelRegistryKind>();

  return (
    <>
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
        rowRenderer={(mr: ModelRegistryKind) => (
          <ModelRegistriesTableRow
            key={mr.metadata.name}
            modelRegistry={mr}
            onDeleteRegistry={(i) => setDeleteRegistry(i)}
            // eslint-disable-next-line @typescript-eslint/no-empty-function
            onEditRegistry={() => {}}
          />
        )}
        variant="compact"
      />
      {/* TODO: implement when CRD endpoint is ready */}
      {/* {editRegistry ? (
        <CreateModal
          modelRegistry={editRegistry}
          onClose={() => setEditRegistry(undefined)}
          refresh={refresh}
        />
      ) : null} */}
      {deleteRegistry ? (
        <DeleteModelRegistryModal
          modelRegistry={deleteRegistry}
          onClose={() => setDeleteRegistry(undefined)}
          refresh={refresh}
        />
      ) : null}
    </>
  );
};

export default ModelRegistriesTable;

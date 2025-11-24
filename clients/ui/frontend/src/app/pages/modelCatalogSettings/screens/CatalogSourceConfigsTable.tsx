import * as React from 'react';
import { Button, Toolbar, ToolbarContent, ToolbarItem } from '@patternfly/react-core';
import { PlusCircleIcon } from '@patternfly/react-icons';
import { Table } from 'mod-arch-shared';
import { CatalogSourceConfig } from '~/app/modelCatalogTypes';
import { catalogSourceConfigsColumns } from './CatalogSourceConfigsTableColumns';
import CatalogSourceConfigsTableRow from './CatalogSourceConfigsTableRow';

type CatalogSourceConfigsTableProps = {
  catalogSourceConfigs: CatalogSourceConfig[];
  onAddSource: () => void;
  onDeleteSource?: (config: CatalogSourceConfig) => void;
};

const CatalogSourceConfigsTable: React.FC<CatalogSourceConfigsTableProps> = ({
  catalogSourceConfigs,
  onAddSource,
  onDeleteSource,
}) => (
  <Table
    data-testid="catalog-source-configs-table"
    data={catalogSourceConfigs}
    columns={catalogSourceConfigsColumns}
    toolbarContent={
      <Toolbar>
        <ToolbarContent>
          <ToolbarItem>
            <Button
              variant="primary"
              icon={<PlusCircleIcon />}
              onClick={onAddSource}
              data-testid="add-source-button"
            >
              Add a source
            </Button>
          </ToolbarItem>
        </ToolbarContent>
      </Toolbar>
    }
    rowRenderer={(config) => (
      <CatalogSourceConfigsTableRow
        key={config.id}
        catalogSourceConfig={config}
        onDelete={onDeleteSource}
      />
    )}
    variant="compact"
  />
);

export default CatalogSourceConfigsTable;

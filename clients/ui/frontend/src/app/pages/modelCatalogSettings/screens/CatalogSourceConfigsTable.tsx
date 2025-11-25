import * as React from 'react';
import {
  Alert,
  AlertActionCloseButton,
  Button,
  Flex,
  FlexItem,
  Stack,
  StackItem,
  Toolbar,
  ToolbarContent,
  ToolbarItem,
} from '@patternfly/react-core';
import { PlusCircleIcon } from '@patternfly/react-icons';
import { Table } from 'mod-arch-shared';
import { CatalogSourceConfig } from '~/app/modelCatalogTypes';
import { ModelCatalogSettingsContext } from '~/app/context/modelCatalogSettings/ModelCatalogSettingsContext';
import { catalogSourceConfigsColumns } from './CatalogSourceConfigsTableColumns';
import CatalogSourceConfigsTableRow from './CatalogSourceConfigsTableRow';

type CatalogSourceConfigsTableProps = {
  catalogSourceConfigs: CatalogSourceConfig[];
  onAddSource: () => void;
  onDeleteSource: (sourceId: string) => Promise<void>;
};

const CatalogSourceConfigsTable: React.FC<CatalogSourceConfigsTableProps> = ({
  catalogSourceConfigs,
  onAddSource,
  onDeleteSource,
}) => {
  const [toggleError, setToggleError] = React.useState<Error | undefined>(undefined);
  const [isUpdatingToggle, setIsUpdatingToggle] = React.useState(false);
  const { apiState, refreshCatalogSourceConfigs, catalogSourcesLoadError } = React.useContext(
    ModelCatalogSettingsContext,
  );
  const [showAlert, setShowAlert] = React.useState<boolean>(false);

  const handleEnableToggle = async (checked: boolean, catalogSourceConfig: CatalogSourceConfig) => {
    if (!apiState.apiAvailable) {
      setToggleError(new Error('API is not available'));
      setShowAlert(true);
      return;
    }
    setIsUpdatingToggle(true);
    setShowAlert(false);

    try {
      await apiState.api.updateCatalogSourceConfig({}, catalogSourceConfig.id, {
        enabled: checked,
      });
      setToggleError(undefined);
      refreshCatalogSourceConfigs();
    } catch (e) {
      if (e instanceof Error) {
        setToggleError(new Error(`Error enabling/disabling source ${catalogSourceConfig.name}`));
        setShowAlert(true);
      }
    } finally {
      setIsUpdatingToggle(false);
    }
  };

  return (
    <Stack hasGutter>
      {catalogSourcesLoadError && (
        <StackItem>
          <Alert
            variant="danger"
            isInline
            title="Error fetching source statuses"
            data-testid="source-status-error-alert"
          >
            {catalogSourcesLoadError.message}
          </Alert>
        </StackItem>
      )}
      <StackItem>
        <Table
          data-testid="catalog-source-configs-table"
          data={catalogSourceConfigs}
          columns={catalogSourceConfigsColumns}
          toolbarContent={
            <Flex direction={{ default: 'column' }}>
              <FlexItem>
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
              </FlexItem>
              {toggleError && showAlert && (
                <FlexItem>
                  <Alert
                    variant="danger"
                    data-testid="toggle-alert"
                    title={toggleError.message}
                    actionClose={<AlertActionCloseButton onClose={() => setShowAlert(false)} />}
                  />
                </FlexItem>
              )}
            </Flex>
          }
          rowRenderer={(config) => (
            <CatalogSourceConfigsTableRow
              key={config.id}
              catalogSourceConfig={config}
              isUpdatingToggle={isUpdatingToggle}
              onToggleUpdate={handleEnableToggle}
              onDeleteSource={onDeleteSource}
            />
          )}
          variant="compact"
        />
      </StackItem>
    </Stack>
  );
};

export default CatalogSourceConfigsTable;

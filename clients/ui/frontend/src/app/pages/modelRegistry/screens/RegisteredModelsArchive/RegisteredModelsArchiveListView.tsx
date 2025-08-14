import * as React from 'react';
import { Toolbar, ToolbarContent, ToolbarGroup, ToolbarToggleGroup } from '@patternfly/react-core';
import { FilterIcon, SearchIcon } from '@patternfly/react-icons';
import { ModelVersion, RegisteredModel } from '~/app/types';
import { filterRegisteredModels } from '~/app/pages/modelRegistry/screens/utils';
import EmptyModelRegistryState from '~/app/pages/modelRegistry/screens/components/EmptyModelRegistryState';
import {
  ModelRegistryFilterDataType,
  ModelRegistryFilterOptions,
  initialModelRegistryFilterData,
  modelRegistryFilterOptions,
} from '~/app/pages/modelRegistry/screens/const';
import FilterToolbar from '~/app/shared/components/FilterToolbar';
import ThemeAwareSearchInput from '~/app/pages/modelRegistry/screens/components/ThemeAwareSearchInput';
import RegisteredModelsArchiveTable from './RegisteredModelsArchiveTable';

type RegisteredModelsArchiveListViewProps = {
  registeredModels: RegisteredModel[];
  modelVersions: ModelVersion[];
  refresh: () => void;
};

const RegisteredModelsArchiveListView: React.FC<RegisteredModelsArchiveListViewProps> = ({
  registeredModels: unfilteredRegisteredModels,
  modelVersions,
  refresh,
}) => {
  const [filterData, setFilterData] = React.useState<ModelRegistryFilterDataType>(
    initialModelRegistryFilterData,
  );

  const onFilterUpdate = React.useCallback(
    (key: string, value: string | { label: string; value: string } | undefined) =>
      setFilterData((prevValues) => ({ ...prevValues, [key]: value })),
    [setFilterData],
  );

  const onClearFilters = React.useCallback(
    () => setFilterData(initialModelRegistryFilterData),
    [setFilterData],
  );

  const filteredRegisteredModels = filterRegisteredModels(
    unfilteredRegisteredModels,
    modelVersions,
    filterData,
  );

  if (unfilteredRegisteredModels.length === 0) {
    return (
      <EmptyModelRegistryState
        headerIcon={SearchIcon}
        testid="empty-archive-model-state"
        title="No archived models"
        description="You can archive the active models that you no longer use. You can restore an archived
      model to make it active."
      />
    );
  }

  return (
    <RegisteredModelsArchiveTable
      refresh={refresh}
      clearFilters={onClearFilters}
      registeredModels={filteredRegisteredModels}
      modelVersions={modelVersions}
      toolbarContent={
        <Toolbar
          data-testid="registered-models-archive-table-toolbar"
          clearAllFilters={onClearFilters}
        >
          <ToolbarContent>
            <ToolbarToggleGroup toggleIcon={<FilterIcon />} breakpoint="xl">
              <ToolbarGroup variant="filter-group">
                <FilterToolbar
                  filterOptions={modelRegistryFilterOptions}
                  filterOptionRenders={{
                    [ModelRegistryFilterOptions.keyword]: ({ onChange, ...props }) => (
                      <ThemeAwareSearchInput
                        {...props}
                        fieldLabel="Filter by keyword"
                        placeholder="Filter by keyword"
                        className="toolbar-fieldset-wrapper"
                        style={{ minWidth: '270px' }}
                        onChange={(value) => onChange(value)}
                      />
                    ),
                    [ModelRegistryFilterOptions.owner]: ({ onChange, ...props }) => (
                      <ThemeAwareSearchInput
                        {...props}
                        fieldLabel="Filter by owner"
                        placeholder="Filter by owner"
                        className="toolbar-fieldset-wrapper"
                        style={{ minWidth: '270px' }}
                        onChange={(value) => onChange(value)}
                      />
                    ),
                  }}
                  filterData={filterData}
                  onFilterUpdate={onFilterUpdate}
                />
              </ToolbarGroup>
            </ToolbarToggleGroup>
          </ToolbarContent>
        </Toolbar>
      }
    />
  );
};

export default RegisteredModelsArchiveListView;

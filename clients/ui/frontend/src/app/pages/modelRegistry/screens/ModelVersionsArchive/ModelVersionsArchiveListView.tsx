import * as React from 'react';
import { Toolbar, ToolbarContent, ToolbarGroup, ToolbarToggleGroup } from '@patternfly/react-core';
import { FilterIcon, SearchIcon } from '@patternfly/react-icons';
import { ModelVersion } from '~/app/types';
import { filterModelVersions } from '~/app/pages/modelRegistry/screens/utils';
import EmptyModelRegistryState from '~/app/pages/modelRegistry/screens/components/EmptyModelRegistryState';
import FilterToolbar from '~/app/shared/components/FilterToolbar';
import {
  initialModelRegistryVersionsFilterData,
  ModelRegistryVersionsFilterDataType,
  modelRegistryVersionsFilterOptions,
  ModelRegistryVersionsFilterOptions,
} from '~/app/pages/modelRegistry/screens/const';
import ThemeAwareSearchInput from '~/app/pages/modelRegistry/screens/components/ThemeAwareSearchInput';
import ModelVersionsArchiveTable from './ModelVersionsArchiveTable';

type ModelVersionsArchiveListViewProps = {
  modelVersions: ModelVersion[];
  refresh: () => void;
};

const ModelVersionsArchiveListView: React.FC<ModelVersionsArchiveListViewProps> = ({
  modelVersions: unfilteredmodelVersions,
  refresh,
}) => {
  const [filterData, setFilterData] = React.useState<ModelRegistryVersionsFilterDataType>(
    initialModelRegistryVersionsFilterData,
  );

  const onFilterUpdate = React.useCallback(
    (key: string, value: string | { label: string; value: string } | undefined) =>
      setFilterData((prevValues) => ({ ...prevValues, [key]: value })),
    [setFilterData],
  );

  const onClearFilters = React.useCallback(
    () => setFilterData(initialModelRegistryVersionsFilterData),
    [setFilterData],
  );

  const filteredModelVersions = filterModelVersions(unfilteredmodelVersions, filterData);

  if (unfilteredmodelVersions.length === 0) {
    return (
      <EmptyModelRegistryState
        headerIcon={SearchIcon}
        testid="empty-archive-state"
        title="No archived versions"
        description="You can archive the active versions that you no longer use. You can restore archived versions to make it active."
      />
    );
  }

  return (
    <ModelVersionsArchiveTable
      refresh={refresh}
      clearFilters={onClearFilters}
      modelVersions={filteredModelVersions}
      toolbarContent={
        <Toolbar
          data-testid="model-versions-archive-table-toolbar"
          clearAllFilters={onClearFilters}
        >
          <ToolbarContent>
            <ToolbarToggleGroup toggleIcon={<FilterIcon />} breakpoint="xl">
              <ToolbarGroup variant="filter-group">
                <FilterToolbar
                  filterOptions={modelRegistryVersionsFilterOptions}
                  filterOptionRenders={{
                    [ModelRegistryVersionsFilterOptions.keyword]: ({ onChange, ...props }) => (
                      <ThemeAwareSearchInput
                        {...props}
                        fieldLabel="Filter by keyword"
                        placeholder="Filter by keyword"
                        className="toolbar-fieldset-wrapper"
                        style={{ minWidth: '270px' }}
                        onChange={(value) => onChange(value)}
                      />
                    ),
                    [ModelRegistryVersionsFilterOptions.author]: ({ onChange, ...props }) => (
                      <ThemeAwareSearchInput
                        {...props}
                        fieldLabel="Filter by author"
                        placeholder="Filter by author"
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

export default ModelVersionsArchiveListView;

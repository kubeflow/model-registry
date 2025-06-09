import * as React from 'react';
import {
  Toolbar,
  ToolbarContent,
  ToolbarFilter,
  ToolbarGroup,
  ToolbarItem,
  ToolbarToggleGroup,
} from '@patternfly/react-core';
import { FilterIcon, SearchIcon } from '@patternfly/react-icons';
import { asEnumMember, SimpleSelect } from 'mod-arch-shared';
import { SearchType } from 'mod-arch-shared/dist/components/DashboardSearchField';
import { ModelVersion, RegisteredModel } from '~/app/types';
import { filterRegisteredModels } from '~/app/pages/modelRegistry/screens/utils';
import EmptyModelRegistryState from '~/app/pages/modelRegistry/screens/components/EmptyModelRegistryState';
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
  const [searchType, setSearchType] = React.useState<SearchType>(SearchType.KEYWORD);
  const [search, setSearch] = React.useState('');

  const searchTypes = [SearchType.KEYWORD, SearchType.OWNER];
  const filteredRegisteredModels = filterRegisteredModels(
    unfilteredRegisteredModels,
    modelVersions,
    search,
    searchType,
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

  const resetFilters = () => setSearch('');

  return (
    <RegisteredModelsArchiveTable
      refresh={refresh}
      clearFilters={resetFilters}
      registeredModels={filteredRegisteredModels}
      toolbarContent={
        <Toolbar>
          <ToolbarContent>
            <ToolbarToggleGroup toggleIcon={<FilterIcon />} breakpoint="xl">
              <ToolbarGroup variant="filter-group">
                <ToolbarFilter
                  labels={search === '' ? [] : [search]}
                  deleteLabel={resetFilters}
                  deleteLabelGroup={resetFilters}
                  categoryName="Keyword"
                >
                  <SimpleSelect
                    options={searchTypes.map((key) => ({
                      key,
                      label: key,
                    }))}
                    value={searchType}
                    onChange={(newSearchType) => {
                      const newSearchTypeInput = asEnumMember(newSearchType, SearchType);
                      if (newSearchTypeInput !== null) {
                        setSearchType(newSearchTypeInput);
                      }
                    }}
                    icon={<FilterIcon />}
                  />
                </ToolbarFilter>
                <ToolbarItem>
                  <ThemeAwareSearchInput
                    value={search}
                    onChange={setSearch}
                    onClear={resetFilters}
                    placeholder={`Find by ${searchType.toLowerCase()}`}
                    fieldLabel={`Find by ${searchType.toLowerCase()}`}
                    className="toolbar-fieldset-wrapper"
                    style={{ minWidth: '200px' }}
                    data-testid="registered-models-archive-table-search"
                  />
                </ToolbarItem>
              </ToolbarGroup>
            </ToolbarToggleGroup>
          </ToolbarContent>
        </Toolbar>
      }
    />
  );
};

export default RegisteredModelsArchiveListView;

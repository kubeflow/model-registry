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
import { ModelVersion } from '~/app/types';
import { filterModelVersions } from '~/app/pages/modelRegistry/screens/utils';
import EmptyModelRegistryState from '~/app/pages/modelRegistry/screens/components/EmptyModelRegistryState';
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
  const [searchType, setSearchType] = React.useState<SearchType>(SearchType.KEYWORD);
  const [search, setSearch] = React.useState('');

  const searchTypes = [SearchType.KEYWORD, SearchType.AUTHOR];

  const filteredModelVersions = filterModelVersions(unfilteredmodelVersions, search, searchType);

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

  const resetFilters = () => setSearch('');

  return (
    <ModelVersionsArchiveTable
      refresh={refresh}
      clearFilters={resetFilters}
      modelVersions={filteredModelVersions}
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
                      const enumMember = asEnumMember(newSearchType, SearchType);
                      if (enumMember) {
                        setSearchType(enumMember);
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
                    data-testid="model-versions-archive-table-search"
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

export default ModelVersionsArchiveListView;

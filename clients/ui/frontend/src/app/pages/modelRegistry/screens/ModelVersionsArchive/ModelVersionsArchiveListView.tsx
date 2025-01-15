import * as React from 'react';
import {
  SearchInput,
  TextInput,
  ToolbarContent,
  ToolbarFilter,
  ToolbarGroup,
  ToolbarItem,
  ToolbarToggleGroup,
} from '@patternfly/react-core';
import { FilterIcon, SearchIcon } from '@patternfly/react-icons';
import { ModelVersion } from '~/app/types';
import { SearchType } from '~/shared/components/DashboardSearchField';
import SimpleSelect from '~/shared/components/SimpleSelect';
import { asEnumMember } from '~/shared/utilities/utils';
import { filterModelVersions } from '~/app/pages/modelRegistry/screens/utils';
import EmptyModelRegistryState from '~/app/pages/modelRegistry/screens/components/EmptyModelRegistryState';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import { isMUITheme } from '~/shared/utilities/const';
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

  return (
    <ModelVersionsArchiveTable
      refresh={refresh}
      clearFilters={() => setSearch('')}
      modelVersions={filteredModelVersions}
      toolbarContent={
        <ToolbarContent>
          <ToolbarToggleGroup toggleIcon={<FilterIcon />} breakpoint="xl">
            <ToolbarGroup variant="filter-group">
              <ToolbarFilter
                labels={search === '' ? [] : [search]}
                deleteLabel={() => setSearch('')}
                deleteLabelGroup={() => setSearch('')}
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
                {isMUITheme() ? (
                  <FormFieldset
                    className="toolbar-fieldset-wrapper"
                    component={
                      <TextInput
                        value={search}
                        type="text"
                        onChange={(_, searchValue) => {
                          setSearch(searchValue);
                        }}
                        style={{ minWidth: '200px' }}
                        data-testid="model-versions-archive-table-search"
                        aria-label="Search"
                      />
                    }
                    field={`Find by ${searchType.toLowerCase()}`}
                  />
                ) : (
                  <SearchInput
                    placeholder={`Find by ${searchType.toLowerCase()}`}
                    value={search}
                    onChange={(_, searchValue) => {
                      setSearch(searchValue);
                    }}
                    onClear={() => setSearch('')}
                    style={{ minWidth: '200px' }}
                    data-testid="model-versions-archive-table-search"
                  />
                )}
              </ToolbarItem>
            </ToolbarGroup>
          </ToolbarToggleGroup>
        </ToolbarContent>
      }
    />
  );
};

export default ModelVersionsArchiveListView;

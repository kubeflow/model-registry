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
import { FilterIcon } from '@patternfly/react-icons';
import { RegisteredModel } from '~/app/types';
import { SearchType } from '~/shared/components/DashboardSearchField';
import { filterRegisteredModels } from '~/app/pages/modelRegistry/screens/utils';
import EmptyModelRegistryState from '~/app/pages/modelRegistry/screens/components/EmptyModelRegistryState';
import SimpleSelect from '~/shared/components/SimpleSelect';
import { asEnumMember } from '~/app/utils';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import { isMUITheme } from '~/shared/utilities/const';
import RegisteredModelsArchiveTable from './RegisteredModelsArchiveTable';

type RegisteredModelsArchiveListViewProps = {
  registeredModels: RegisteredModel[];
  refresh: () => void;
};

const RegisteredModelsArchiveListView: React.FC<RegisteredModelsArchiveListViewProps> = ({
  registeredModels: unfilteredRegisteredModels,
  refresh,
}) => {
  const [searchType, setSearchType] = React.useState<SearchType>(SearchType.KEYWORD);
  const [search, setSearch] = React.useState('');

  const searchTypes = [SearchType.KEYWORD, SearchType.AUTHOR];

  const filteredRegisteredModels = filterRegisteredModels(
    unfilteredRegisteredModels,
    search,
    searchType,
  );

  if (unfilteredRegisteredModels.length === 0) {
    return (
      <EmptyModelRegistryState
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
      clearFilters={() => setSearch('')}
      registeredModels={filteredRegisteredModels}
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
                    const newSearchTypeInput = asEnumMember(newSearchType, SearchType);
                    if (newSearchTypeInput !== null) {
                      setSearchType(newSearchTypeInput);
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
                        data-testid="registered-models-archive-table-search"
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
                    data-testid="registered-models-archive-table-search"
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

export default RegisteredModelsArchiveListView;

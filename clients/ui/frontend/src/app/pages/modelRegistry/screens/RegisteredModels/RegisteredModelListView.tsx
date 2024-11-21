import * as React from 'react';
import {
  SearchInput,
  TextInput,
  ToolbarFilter,
  ToolbarGroup,
  ToolbarItem,
} from '@patternfly/react-core';
import { FilterIcon } from '@patternfly/react-icons';
import { useNavigate } from 'react-router-dom';
import { RegisteredModel } from '~/app/types';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import { SearchType } from '~/shared/components/DashboardSearchField';
import { ProjectObjectType, typedEmptyImage } from '~/shared/components/design/utils';
import { asEnumMember, filterRegisteredModels } from '~/app/utils';
import SimpleSelect from '~/shared/components/SimpleSelect';
import {
  registeredModelArchiveUrl,
  registerModelUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';
import EmptyModelRegistryState from '~/app/pages/modelRegistry/screens/components/EmptyModelRegistryState';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import { isMUITheme } from '~/shared/utilities/const';
import RegisteredModelTable from './RegisteredModelTable';
import RegisteredModelsTableToolbar from './RegisteredModelsTableToolbar';

type RegisteredModelListViewProps = {
  registeredModels: RegisteredModel[];
  refresh: () => void;
};

const RegisteredModelListView: React.FC<RegisteredModelListViewProps> = ({
  registeredModels: unfilteredRegisteredModels,
  refresh,
}) => {
  const navigate = useNavigate();
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);
  const [searchType, setSearchType] = React.useState<SearchType>(SearchType.KEYWORD);
  const [search, setSearch] = React.useState('');

  const searchTypes = React.useMemo(() => [SearchType.KEYWORD, SearchType.OWNER], []);

  if (unfilteredRegisteredModels.length === 0) {
    return (
      <EmptyModelRegistryState
        testid="empty-registered-models"
        title="No models in selected registry"
        headerIcon={() => (
          <img
            src={typedEmptyImage(ProjectObjectType.registeredModels, 'MissingModel')}
            alt="missing model"
          />
        )}
        description={`${preferredModelRegistry?.name} has no active registered models. Register a model in this registry, or select a different registry.`}
        primaryActionText="Register model"
        secondaryActionText="View archived models"
        primaryActionOnClick={() => {
          navigate(registerModelUrl(preferredModelRegistry?.name));
        }}
        secondaryActionOnClick={() => {
          navigate(registeredModelArchiveUrl(preferredModelRegistry?.name));
        }}
      />
    );
  }

  const filteredRegisteredModels = filterRegisteredModels(
    unfilteredRegisteredModels,
    search,
    searchType,
  );

  const resetFilters = () => {
    setSearch('');
  };

  const toggleGroupItems = (
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
      <ToolbarItem variant="label-group">
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
                data-testid="registered-model-table-search"
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
            data-testid="registered-model-table-search"
          />
        )}
      </ToolbarItem>
    </ToolbarGroup>
  );

  return (
    <RegisteredModelTable
      refresh={refresh}
      clearFilters={resetFilters}
      registeredModels={filteredRegisteredModels}
      toolbarContent={<RegisteredModelsTableToolbar toggleGroupItems={toggleGroupItems} />}
    />
  );
};

export default RegisteredModelListView;

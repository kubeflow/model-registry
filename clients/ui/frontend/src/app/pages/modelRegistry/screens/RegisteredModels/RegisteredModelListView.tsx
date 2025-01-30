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
import { ModelVersion, RegisteredModel } from '~/app/types';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import { SearchType } from '~/shared/components/DashboardSearchField';
import { ProjectObjectType, typedEmptyImage } from '~/shared/components/design/utils';
import SimpleSelect from '~/shared/components/SimpleSelect';
import {
  registeredModelArchiveUrl,
  registerModelUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';
import EmptyModelRegistryState from '~/app/pages/modelRegistry/screens/components/EmptyModelRegistryState';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import { isMUITheme } from '~/shared/utilities/const';
import { filterRegisteredModels } from '~/app/pages/modelRegistry/screens/utils';
import { asEnumMember } from '~/shared/utilities/utils';
import { filterArchiveModels, filterLiveModels } from '~/app/utils';
import RegisteredModelTable from './RegisteredModelTable';
import RegisteredModelsTableToolbar from './RegisteredModelsTableToolbar';

type RegisteredModelListViewProps = {
  registeredModels: RegisteredModel[];
  modelVersions: ModelVersion[];
  refresh: () => void;
};

const RegisteredModelListView: React.FC<RegisteredModelListViewProps> = ({
  registeredModels,
  modelVersions,
  refresh,
}) => {
  const navigate = useNavigate();
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);
  const [searchType, setSearchType] = React.useState<SearchType>(SearchType.KEYWORD);
  const [search, setSearch] = React.useState('');
  const unfilteredRegisteredModels = filterLiveModels(registeredModels);
  const archiveRegisteredModels = filterArchiveModels(registeredModels);
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
        description={`${
          preferredModelRegistry?.name ?? ''
        } has no active registered models. Register a model in this registry, or select a different registry.`}
        primaryActionText="Register model"
        secondaryActionText={
          archiveRegisteredModels.length !== 0 ? 'View archived models' : undefined
        }
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
    modelVersions,
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
            onClear={resetFilters}
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
      toolbarContent={
        <RegisteredModelsTableToolbar
          toggleGroupItems={toggleGroupItems}
          onClearAllFilters={resetFilters}
        />
      }
    />
  );
};

export default RegisteredModelListView;

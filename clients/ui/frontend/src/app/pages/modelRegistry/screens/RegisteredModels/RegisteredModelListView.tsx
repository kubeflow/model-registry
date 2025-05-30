import * as React from 'react';
import { ToolbarFilter, ToolbarGroup, ToolbarItem } from '@patternfly/react-core';
import { FilterIcon } from '@patternfly/react-icons';
import { useNavigate } from 'react-router-dom';
import { ProjectObjectType, typedEmptyImage, asEnumMember, SimpleSelect } from 'mod-arch-shared';
import { SearchType } from 'mod-arch-shared/dist/components/DashboardSearchField';
import { ModelVersion, RegisteredModel } from '~/app/types';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import {
  registeredModelArchiveUrl,
  registerModelUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';
import EmptyModelRegistryState from '~/app/pages/modelRegistry/screens/components/EmptyModelRegistryState';
import { filterRegisteredModels } from '~/app/pages/modelRegistry/screens/utils';
import { filterArchiveModels, filterLiveModels } from '~/app/utils';
import ThemeAwareSearchInput from '~/app/pages/modelRegistry/screens/components/ThemeAwareSearchInput';
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
          toggleProps={{ style: { minWidth: '150px' } }}
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
          data-testid="registered-model-table-search"
        />
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

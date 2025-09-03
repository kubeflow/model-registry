import * as React from 'react';
import { ToolbarGroup } from '@patternfly/react-core';
import { useNavigate } from 'react-router-dom';
import { ProjectObjectType, typedEmptyImage } from 'mod-arch-shared';
import { ModelVersion, RegisteredModel } from '~/app/types';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import {
  registeredModelArchiveUrl,
  registerModelUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';
import EmptyModelRegistryState from '~/app/pages/modelRegistry/screens/components/EmptyModelRegistryState';
import { filterRegisteredModels } from '~/app/pages/modelRegistry/screens/utils';
import { filterArchiveModels, filterLiveModels } from '~/app/utils';
import {
  initialModelRegistryFilterData,
  ModelRegistryFilterDataType,
  modelRegistryFilterOptions,
  ModelRegistryFilterOptions,
} from '~/app/pages/modelRegistry/screens/const';
import FilterToolbar from '~/app/shared/components/FilterToolbar';
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
  const [filterData, setFilterData] = React.useState<ModelRegistryFilterDataType>(
    initialModelRegistryFilterData,
  );
  const unfilteredRegisteredModels = filterLiveModels(registeredModels);
  const archiveRegisteredModels = filterArchiveModels(registeredModels);

  const onFilterUpdate = React.useCallback(
    (key: string, value: string | { label: string; value: string } | undefined) =>
      setFilterData((prevValues) => ({ ...prevValues, [key]: value })),
    [setFilterData],
  );

  const onClearFilters = React.useCallback(
    () => setFilterData(initialModelRegistryFilterData),
    [setFilterData],
  );

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
    filterData,
  );

  const toggleGroupItems = (
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
  );

  return (
    <RegisteredModelTable
      refresh={refresh}
      clearFilters={onClearFilters}
      registeredModels={filteredRegisteredModels}
      modelVersions={modelVersions}
      toolbarContent={
        <RegisteredModelsTableToolbar
          toggleGroupItems={toggleGroupItems}
          onClearAllFilters={onClearFilters}
        />
      }
    />
  );
};

export default RegisteredModelListView;

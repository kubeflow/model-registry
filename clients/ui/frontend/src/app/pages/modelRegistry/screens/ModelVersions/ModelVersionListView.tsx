import * as React from 'react';
import {
  Alert,
  Button,
  Dropdown,
  DropdownItem,
  DropdownList,
  Flex,
  MenuToggle,
  MenuToggleElement,
  Toolbar,
  ToolbarContent,
  ToolbarGroup,
  ToolbarItem,
  ToolbarToggleGroup,
} from '@patternfly/react-core';
import { EllipsisVIcon, FilterIcon } from '@patternfly/react-icons';
import { useNavigate } from 'react-router-dom';
import { ProjectObjectType, typedEmptyImage } from 'mod-arch-shared';
import { ModelVersion, RegisteredModel } from '~/app/types';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import EmptyModelRegistryState from '~/app/pages/modelRegistry/screens/components/EmptyModelRegistryState';
import {
  modelVersionArchiveUrl,
  registerVersionForModelUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';
import {
  filterModelVersions,
  sortModelVersionsByCreateTime,
} from '~/app/pages/modelRegistry/screens/utils';
import ModelVersionsTable from '~/app/pages/modelRegistry/screens/ModelVersions/ModelVersionsTable';
import { filterArchiveVersions, filterLiveVersions } from '~/app/utils';
import {
  initialModelRegistryVersionsFilterData,
  ModelRegistryVersionsFilterDataType,
  modelRegistryVersionsFilterOptions,
  ModelRegistryVersionsFilterOptions,
} from '~/app/pages/modelRegistry/screens/const';
import FilterToolbar from '~/app/shared/components/FilterToolbar';
import ThemeAwareSearchInput from '~/app/pages/modelRegistry/screens/components/ThemeAwareSearchInput';

type ModelVersionListViewProps = {
  modelVersions: ModelVersion[];
  registeredModel: RegisteredModel;
  isArchiveModel?: boolean;
  refresh: () => void;
};

const ModelVersionListView: React.FC<ModelVersionListViewProps> = ({
  modelVersions,
  registeredModel: rm,
  isArchiveModel,
  refresh,
}) => {
  const unfilteredModelVersions = isArchiveModel
    ? modelVersions
    : filterLiveVersions(modelVersions);

  const archiveModelVersions = filterArchiveVersions(modelVersions);
  const navigate = useNavigate();
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);
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

  const [isArchivedModelVersionKebabOpen, setIsArchivedModelVersionKebabOpen] =
    React.useState(false);

  const filteredModelVersions = filterModelVersions(unfilteredModelVersions, filterData);
  const date = rm.lastUpdateTimeSinceEpoch && new Date(parseInt(rm.lastUpdateTimeSinceEpoch));

  if (unfilteredModelVersions.length === 0) {
    if (isArchiveModel) {
      return (
        <EmptyModelRegistryState
          testid="empty-archive-model-versions"
          title="No versions"
          headerIcon={() => (
            <img
              src={typedEmptyImage(ProjectObjectType.registeredModels, 'MissingVersion')}
              alt="missing version"
            />
          )}
          description={`${rm.name} has no registered versions.`}
        />
      );
    }
    return (
      <EmptyModelRegistryState
        testid="empty-model-versions"
        title="No versions"
        headerIcon={() => (
          <img
            src={typedEmptyImage(ProjectObjectType.registeredModels, 'MissingVersion')}
            alt="missing version"
          />
        )}
        description={`${rm.name} has no registered versions. Register a version to this model.`}
        primaryActionText="Register new version"
        primaryActionOnClick={() => {
          navigate(registerVersionForModelUrl(rm.id, preferredModelRegistry?.name));
        }}
        secondaryActionText={
          archiveModelVersions.length !== 0 ? 'View archived versions' : undefined
        }
        secondaryActionOnClick={() => {
          navigate(modelVersionArchiveUrl(rm.id, preferredModelRegistry?.name));
        }}
      />
    );
  }

  return (
    <>
      {isArchiveModel && (
        <Alert
          variant="warning"
          isInline
          title={`All the versions have been archived along with the model on ${
            date
              ? `${date.toLocaleString('en-US', {
                  month: 'long',
                  timeZone: 'UTC',
                })} ${date.getUTCDate()}, ${date.getUTCFullYear()}`
              : '--'
          }. They are now read-only and can only be restored together with the model.`}
        />
      )}
      <ModelVersionsTable
        refresh={refresh}
        isArchiveModel={isArchiveModel}
        clearFilters={onClearFilters}
        modelVersions={sortModelVersionsByCreateTime(filteredModelVersions)}
        toolbarContent={
          <Toolbar data-testid="model-versions-table-toolbar" clearAllFilters={onClearFilters}>
            <ToolbarContent>
              {/* TODO: Remove this Flex after the ToolbarContent can center the children elements */}
              <Flex>
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

                {!isArchiveModel && (
                  <ToolbarGroup>
                    <ToolbarItem>
                      <Button
                        variant="primary"
                        onClick={() => {
                          navigate(registerVersionForModelUrl(rm.id, preferredModelRegistry?.name));
                        }}
                      >
                        Register new version
                      </Button>
                    </ToolbarItem>
                    <ToolbarItem>
                      <Dropdown
                        isOpen={isArchivedModelVersionKebabOpen}
                        onSelect={() => setIsArchivedModelVersionKebabOpen(false)}
                        onOpenChange={(isOpen: boolean) =>
                          setIsArchivedModelVersionKebabOpen(isOpen)
                        }
                        toggle={(tr: React.Ref<MenuToggleElement>) => (
                          <MenuToggle
                            data-testid="model-versions-table-kebab-action"
                            ref={tr}
                            variant="plain"
                            onClick={() =>
                              setIsArchivedModelVersionKebabOpen(!isArchivedModelVersionKebabOpen)
                            }
                            isExpanded={isArchivedModelVersionKebabOpen}
                            aria-label="View archived versions"
                          >
                            <EllipsisVIcon />
                          </MenuToggle>
                        )}
                        shouldFocusToggleOnSelect
                        popperProps={{ appendTo: 'inline' }}
                      >
                        <DropdownList>
                          <DropdownItem
                            onClick={() =>
                              navigate(modelVersionArchiveUrl(rm.id, preferredModelRegistry?.name))
                            }
                          >
                            View archived versions
                          </DropdownItem>
                        </DropdownList>
                      </Dropdown>
                    </ToolbarItem>
                  </ToolbarGroup>
                )}
              </Flex>
            </ToolbarContent>
          </Toolbar>
        }
      />
    </>
  );
};

export default ModelVersionListView;

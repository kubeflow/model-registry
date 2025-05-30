import * as React from 'react';
import {
  Alert,
  Button,
  Dropdown,
  DropdownItem,
  DropdownList,
  MenuToggle,
  MenuToggleElement,
  Toolbar,
  ToolbarContent,
  ToolbarFilter,
  ToolbarGroup,
  ToolbarItem,
  ToolbarToggleGroup,
} from '@patternfly/react-core';
import { EllipsisVIcon, FilterIcon } from '@patternfly/react-icons';
import { useNavigate } from 'react-router-dom';
import { ProjectObjectType, typedEmptyImage, SimpleSelect, asEnumMember } from 'mod-arch-shared';
import { SearchType } from 'mod-arch-shared/dist/components/DashboardSearchField';
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

  const [searchType, setSearchType] = React.useState<SearchType>(SearchType.KEYWORD);
  const [search, setSearch] = React.useState('');

  const searchTypes = [SearchType.KEYWORD, SearchType.AUTHOR];

  const [isArchivedModelVersionKebabOpen, setIsArchivedModelVersionKebabOpen] =
    React.useState(false);

  const filteredModelVersions = filterModelVersions(unfilteredModelVersions, search, searchType);
  const date = rm.lastUpdateTimeSinceEpoch && new Date(parseInt(rm.lastUpdateTimeSinceEpoch));

  const resetFilters = () => setSearch('');

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
        clearFilters={resetFilters}
        modelVersions={sortModelVersionsByCreateTime(filteredModelVersions)}
        toolbarContent={
          <Toolbar data-testid="model-versions-table-toolbar">
            <ToolbarContent>
              <ToolbarToggleGroup toggleIcon={<FilterIcon />} breakpoint="xl">
                <ToolbarGroup variant="filter-group">
                  <ToolbarFilter
                    labels={search === '' ? [] : [search]}
                    deleteLabel={resetFilters}
                    deleteLabelGroup={resetFilters}
                    categoryName={searchType}
                  >
                    <SimpleSelect
                      dataTestId="model-versions-table-filter"
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
                      data-testid="model-versions-table-search"
                    />
                  </ToolbarItem>
                </ToolbarGroup>

                {!isArchiveModel && (
                  <>
                    <ToolbarGroup>
                      <ToolbarItem>
                        <Button
                          variant="primary"
                          onClick={() => {
                            navigate(
                              registerVersionForModelUrl(rm.id, preferredModelRegistry?.name),
                            );
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
                                navigate(
                                  modelVersionArchiveUrl(rm.id, preferredModelRegistry?.name),
                                )
                              }
                            >
                              View archived versions
                            </DropdownItem>
                          </DropdownList>
                        </Dropdown>
                      </ToolbarItem>
                    </ToolbarGroup>
                  </>
                )}
              </ToolbarToggleGroup>
            </ToolbarContent>
          </Toolbar>
        }
      />
    </>
  );
};

export default ModelVersionListView;

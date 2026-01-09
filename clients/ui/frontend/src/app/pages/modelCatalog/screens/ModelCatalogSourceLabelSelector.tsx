import {
  Alert,
  AlertActionCloseButton,
  Button,
  Content,
  Flex,
  Stack,
  StackItem,
  Toolbar,
  ToolbarContent,
  ToolbarGroup,
  ToolbarItem,
  ToolbarToggleGroup,
} from '@patternfly/react-core';
import { ArrowRightIcon, FilterIcon } from '@patternfly/react-icons';
import React from 'react';
import { useThemeContext } from 'mod-arch-kubeflow';
import { getBasicFiltersToShow } from '~/concepts/modelCatalog/const';
import ModelCatalogActiveFilters from '~/app/pages/modelCatalog/components/ModelCatalogActiveFilters';
import HardwareConfigurationFilterToolbar from '~/app/pages/modelCatalog/components/HardwareConfigurationFilterToolbar';
import ThemeAwareSearchInput from '~/app/pages/modelRegistry/screens/components/ThemeAwareSearchInput';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { hasFiltersApplied } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import ModelCatalogSourceLabelBlocks from './ModelCatalogSourceLabelBlocks';

type ModelCatalogSourceLabelSelectorProps = {
  searchTerm?: string;
  onSearch?: (term: string) => void;
  onClearSearch?: () => void;
  onResetAllFilters?: () => void;
};

const ModelCatalogSourceLabelSelector: React.FC<ModelCatalogSourceLabelSelectorProps> = ({
  searchTerm,
  onSearch,
  onClearSearch,
  onResetAllFilters,
}) => {
  const [inputValue, setInputValue] = React.useState(searchTerm || '');
  const { isMUITheme } = useThemeContext();
  const {
    filterData,
    performanceViewEnabled,
    performanceFiltersChangedOnDetailsPage,
    setPerformanceFiltersChangedOnDetailsPage,
  } = React.useContext(ModelCatalogContext);
  const filtersApplied = React.useMemo(() => hasFiltersApplied(filterData), [filterData]);
  const hasActiveFilters = React.useMemo(
    () => filtersApplied || (searchTerm && searchTerm.trim().length > 0),
    [filtersApplied, searchTerm],
  );

  // Only show basic filters in the main chip bar - performance filters have their own section
  const filtersToShow = React.useMemo(() => getBasicFiltersToShow(), []);

  const shouldShowAlert = performanceViewEnabled && performanceFiltersChangedOnDetailsPage;

  const handleClearAllFilters = React.useCallback(() => {
    if (hasActiveFilters && onResetAllFilters) {
      onResetAllFilters();
    }
  }, [hasActiveFilters, onResetAllFilters]);

  React.useEffect(() => {
    setInputValue(searchTerm || '');
  }, [searchTerm]);

  const handleModelSearch = () => {
    if (onSearch && inputValue.trim() !== searchTerm) {
      onSearch(inputValue.trim());
    }
  };

  const handleClear = () => {
    if (onClearSearch) {
      onClearSearch();
    }
  };

  const handleSearchInputChange = (value: string) => {
    setInputValue(value);
  };

  const handleSearchInputSearch = (_: React.SyntheticEvent<HTMLButtonElement>, value: string) => {
    if (onSearch) {
      onSearch(value.trim());
    }
  };

  return (
    <Stack hasGutter>
      <StackItem>
        <Toolbar
          key={`toolbar-${hasActiveFilters}`}
          {...(onResetAllFilters
            ? {
                clearAllFilters: handleClearAllFilters,
                clearFiltersButtonText: hasActiveFilters ? 'Reset all filters' : '',
              }
            : {})}
        >
          <ToolbarContent>
            <Flex>
              <ToolbarToggleGroup breakpoint="md" toggleIcon={<FilterIcon />}>
                <ToolbarGroup variant="filter-group" gap={{ default: 'gapMd' }} alignItems="center">
                  <ToolbarItem>
                    <ThemeAwareSearchInput
                      dara-testid="search-input"
                      fieldLabel="Filter by name, description and provider"
                      aria-label="Search with submit button"
                      className="toolbar-fieldset-wrapper"
                      placeholder="Filter by name, description and provider"
                      value={inputValue}
                      style={{
                        minWidth: '600px',
                      }}
                      onChange={handleSearchInputChange}
                      onSearch={handleSearchInputSearch}
                      onClear={handleClear}
                    />
                  </ToolbarItem>
                  <ToolbarItem>
                    {isMUITheme && (
                      <Button
                        isInline
                        aria-label="arrow-right-button"
                        data-testid="search-button"
                        variant="link"
                        icon={<ArrowRightIcon />}
                        iconPosition="right"
                        onClick={handleModelSearch}
                      />
                    )}
                  </ToolbarItem>
                </ToolbarGroup>
              </ToolbarToggleGroup>
              {/* When toggle is OFF, show basic filter chips here. When toggle is ON, chips are shown in HardwareConfigurationFilterToolbar */}
              {!performanceViewEnabled && onResetAllFilters && hasActiveFilters && (
                <ModelCatalogActiveFilters filtersToShow={filtersToShow} />
              )}
            </Flex>
          </ToolbarContent>
        </Toolbar>
      </StackItem>
      {performanceViewEnabled && (
        <>
          <StackItem>
            <Content component="h2" className="pf-v6-u-font-weight-bold">
              Workload and performance constraints
            </Content>
          </StackItem>
          <StackItem>
            <HardwareConfigurationFilterToolbar onResetAllFilters={onResetAllFilters} />
          </StackItem>
        </>
      )}
      <StackItem>
        <ModelCatalogSourceLabelBlocks />
      </StackItem>
      {shouldShowAlert && (
        <StackItem>
          <Alert
            variant="info"
            isInline
            className="pf-v6-u-mb-lg"
            title="The results list has been updated to match the latest performance criteria set on the details page."
            actionClose={
              <AlertActionCloseButton
                onClose={() => {
                  setPerformanceFiltersChangedOnDetailsPage(false);
                }}
              />
            }
            data-testid="performance-filters-updated-alert"
          />
        </StackItem>
      )}
    </Stack>
  );
};

export default ModelCatalogSourceLabelSelector;

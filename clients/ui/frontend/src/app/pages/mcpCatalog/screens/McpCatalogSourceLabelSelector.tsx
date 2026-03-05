import * as React from 'react';
import {
  Button,
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
import { useThemeContext } from 'mod-arch-kubeflow';
import ThemeAwareSearchInput from '~/app/pages/modelRegistry/screens/components/ThemeAwareSearchInput';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import { hasMcpFiltersApplied } from '~/app/pages/mcpCatalog/utils/mcpCatalogUtils';
import McpCatalogActiveFilters from '~/app/pages/mcpCatalog/components/McpCatalogActiveFilters';
import McpCatalogSourceLabelBlocks from './McpCatalogSourceLabelBlocks';

type McpCatalogSourceLabelSelectorProps = {
  searchTerm: string;
  onSearch: (term: string) => void;
  onClearSearch: () => void;
  onResetAllFilters: () => void;
};

const McpCatalogSourceLabelSelector: React.FC<McpCatalogSourceLabelSelectorProps> = ({
  searchTerm,
  onSearch,
  onClearSearch,
  onResetAllFilters,
}) => {
  const [inputValue, setInputValue] = React.useState(searchTerm || '');
  const { isMUITheme } = useThemeContext();
  const { filters } = React.useContext(McpCatalogContext);

  const hasFiltersAppliedValue = hasMcpFiltersApplied(filters, searchTerm);

  React.useEffect(() => {
    setInputValue(searchTerm || '');
  }, [searchTerm]);

  const handleClearAllFilters = React.useCallback(() => {
    if (hasFiltersAppliedValue) {
      onResetAllFilters();
    }
  }, [hasFiltersAppliedValue, onResetAllFilters]);

  const handleSearch = React.useCallback(() => {
    if (inputValue.trim() !== searchTerm) {
      onSearch(inputValue.trim());
    }
  }, [inputValue, searchTerm, onSearch]);

  const handleClear = React.useCallback(() => {
    onClearSearch();
  }, [onClearSearch]);

  const handleSearchInputChange = React.useCallback((value: string) => {
    setInputValue(value);
  }, []);

  const handleSearchInputSearch = React.useCallback(
    (_: React.SyntheticEvent<HTMLButtonElement>, value: string) => {
      onSearch(value.trim());
    },
    [onSearch],
  );

  const toolbarClearAllProps = hasFiltersAppliedValue
    ? {
        clearAllFilters: handleClearAllFilters,
        clearFiltersButtonText: 'Reset all filters' as const,
      }
    : undefined;

  return (
    <Stack hasGutter>
      <StackItem>
        <Toolbar
          key={hasFiltersAppliedValue ? 'has-filters' : 'no-filters'}
          {...(toolbarClearAllProps ?? {})}
        >
          <ToolbarContent rowWrap={{ default: 'wrap' }}>
            <Flex>
              <ToolbarToggleGroup breakpoint="md" toggleIcon={<FilterIcon />}>
                <ToolbarGroup variant="filter-group" gap={{ default: 'gapMd' }} alignItems="center">
                  <ToolbarItem>
                    <ThemeAwareSearchInput
                      data-testid="mcp-catalog-search-input"
                      aria-label="Search with submit button"
                      className="toolbar-fieldset-wrapper"
                      placeholder="Filter by name, description and provider"
                      value={inputValue}
                      style={{ minWidth: '600px' }}
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
                        data-testid="mcp-search-button"
                        variant="link"
                        icon={<ArrowRightIcon />}
                        iconPosition="right"
                        onClick={handleSearch}
                      />
                    )}
                  </ToolbarItem>
                </ToolbarGroup>
              </ToolbarToggleGroup>
              {hasFiltersAppliedValue && <McpCatalogActiveFilters />}
            </Flex>
          </ToolbarContent>
        </Toolbar>
      </StackItem>
      <StackItem>
        <Flex
          justifyContent={{ default: 'justifyContentSpaceBetween' }}
          alignItems={{ default: 'alignItemsCenter' }}
        >
          <McpCatalogSourceLabelBlocks />
        </Flex>
      </StackItem>
    </Stack>
  );
};

export default McpCatalogSourceLabelSelector;

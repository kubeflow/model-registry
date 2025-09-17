import {
  Button,
  Dropdown,
  DropdownItem,
  DropdownList,
  Flex,
  MenuToggle,
  Toolbar,
  ToolbarContent,
  ToolbarGroup,
  ToolbarItem,
  ToolbarToggleGroup,
} from '@patternfly/react-core';
import { ArrowRightIcon, FilterIcon } from '@patternfly/react-icons';
import React from 'react';
import { useThemeContext } from 'mod-arch-kubeflow';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import ThemeAwareSearchInput from '~/app/pages/modelRegistry/screens/components/ThemeAwareSearchInput';
import { filterEnabledCatalogSources } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

type ModelCatalogSourceSelectorProps = {
  sourceId: string;
  onSelection: (sourceId: string) => void;
  searchTerm?: string;
  onSearch?: (term: string) => void;
  onClearSearch?: () => void;
  primary?: boolean;
};

const ModelCatalogSourceSelector: React.FC<ModelCatalogSourceSelectorProps> = ({
  sourceId,
  onSelection,
  searchTerm,
  onSearch,
  onClearSearch,
  primary,
}) => {
  const [open, setOpen] = React.useState(false);
  const [inputValue, setInputValue] = React.useState(searchTerm || '');
  const { catalogSources, updateSelectedSource, selectedSource } =
    React.useContext(ModelCatalogContext);
  const selection = catalogSources?.items.find((source) => source.id === sourceId);
  const { isMUITheme } = useThemeContext();

  const enabledCatalogSources = filterEnabledCatalogSources(catalogSources);

  const handleModelSearch = () => {
    if (onSearch && inputValue.trim() !== searchTerm) {
      onSearch(inputValue.trim());
    }
  };

  const handleClear = () => {
    setInputValue('');
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

  const selector = (
    <Dropdown
      shouldFocusToggleOnSelect
      toggle={(toggleRef) => (
        <MenuToggle
          data-testid="source-selector"
          id="source-selector-toggle-button"
          ref={toggleRef}
          aria-label="Filter source"
          onClick={() => setOpen(!open)}
          isExpanded={open}
          icon={<FilterIcon />}
        >
          {selection?.name}
        </MenuToggle>
      )}
      isOpen={open}
      popperProps={{ appendTo: 'inline' }}
    >
      <DropdownList>
        {enabledCatalogSources?.items.map((source) => (
          <DropdownItem
            isSelected={source.id === selectedSource?.id}
            key={source.id}
            id={source.id}
            onClick={() => {
              setOpen(false);
              const catalogSource = enabledCatalogSources.items.find(
                (enabledSource) => enabledSource.id === source.id,
              );
              updateSelectedSource(catalogSource);
              onSelection(source.id);
            }}
          >
            {source.name}
          </DropdownItem>
        ))}
      </DropdownList>
    </Dropdown>
  );

  if (primary) {
    return selector;
  }

  return (
    <Toolbar>
      <ToolbarContent>
        <Flex>
          <ToolbarToggleGroup breakpoint="md" toggleIcon={<FilterIcon />}>
            <ToolbarGroup variant="filter-group" gap={{ default: 'gapMd' }} alignItems="center">
              <ToolbarItem>{selector}</ToolbarItem>
              <ToolbarItem>
                <ThemeAwareSearchInput
                  fieldLabel="Filter by name, description and provider"
                  aria-label="Search with submit button"
                  className="toolbar-fieldset-wrapper"
                  placeholder="Filter by name, description and provider"
                  value={inputValue}
                  style={{
                    minWidth: '400px',
                  }}
                  onChange={handleSearchInputChange}
                  onSearch={handleSearchInputSearch}
                  onClear={handleClear}
                  onClick={() => setOpen(false)}
                />
              </ToolbarItem>
              <ToolbarItem>
                {isMUITheme && (
                  <Button
                    isInline
                    aria-label="arrow-right-button"
                    data-testid="versions-route-link"
                    variant="link"
                    icon={<ArrowRightIcon />}
                    iconPosition="right"
                    onClick={handleModelSearch}
                  />
                )}
              </ToolbarItem>
            </ToolbarGroup>
          </ToolbarToggleGroup>
        </Flex>
      </ToolbarContent>
    </Toolbar>
  );
};

export default ModelCatalogSourceSelector;

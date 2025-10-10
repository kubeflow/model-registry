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
import React from 'react';
import { useThemeContext } from 'mod-arch-kubeflow';
import ThemeAwareSearchInput from '~/app/pages/modelRegistry/screens/components/ThemeAwareSearchInput';
import ModelCatalogSourceLabelBlocks from './ModelCatalogSourceLabelBlocks';

type ModelCatalogSourceLabelSelectorProps = {
  searchTerm?: string;
  onSearch?: (term: string) => void;
  onClearSearch?: () => void;
};

const ModelCatalogSourceLabelSelector: React.FC<ModelCatalogSourceLabelSelectorProps> = ({
  searchTerm,
  onSearch,
  onClearSearch,
}) => {
  const [inputValue, setInputValue] = React.useState(searchTerm || '');
  const { isMUITheme } = useThemeContext();

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
    <Stack>
      <StackItem>
        <Toolbar className="pf-v6-u-pb-0">
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
            </Flex>
          </ToolbarContent>
        </Toolbar>
      </StackItem>
      <StackItem>
        <ModelCatalogSourceLabelBlocks />
      </StackItem>
    </Stack>
  );
};

export default ModelCatalogSourceLabelSelector;

import * as React from 'react';
import {
  Button,
  Dropdown,
  Flex,
  FlexItem,
  MenuToggle,
  MenuToggleElement,
  Popover,
  Slider,
} from '@patternfly/react-core';
import { HelpIcon } from '@patternfly/react-icons';
import { ModelCatalogNumberFilterKey } from '~/concepts/modelCatalog/const';
import { CatalogFilterOptionsList } from '~/app/modelCatalogTypes';
import { useCatalogNumberFilterState } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

const filterKey = ModelCatalogNumberFilterKey.MIN_RPS;

type MinRpsFilterProps = {
  filterOptions?: CatalogFilterOptionsList | null;
};

const MinRpsFilter: React.FC<MinRpsFilterProps> = () => {
  const { value: savedFilterValue, setValue: setSavedFilterValue } =
    useCatalogNumberFilterState(filterKey);
  const [isOpen, setIsOpen] = React.useState(false);

  // Local state for the filter configuration
  const [localValue, setLocalValue] = React.useState<number>(savedFilterValue || 1);

  // Update local value when dropdown opens to show current applied value
  React.useEffect(() => {
    if (isOpen) {
      setLocalValue(savedFilterValue || 1);
    }
  }, [isOpen, savedFilterValue]);

  // Parse saved value if it exists
  const hasActiveFilter = savedFilterValue !== undefined;

  const getDisplayText = (): string => {
    if (hasActiveFilter) {
      return `Min RPS: ${savedFilterValue}`;
    }
    return 'Min RPS';
  };

  const handleApplyFilter = () => {
    // Apply the local value to the actual filter state
    setSavedFilterValue(localValue);
    setIsOpen(false);
  };

  const handleReset = () => {
    setSavedFilterValue(undefined);
    setLocalValue(1);
    setIsOpen(false);
  };

  // Get min/max values from filter options or use defaults
  // TODO: Use real min/max values when available from API
  const minValue = 1; // Default minimum value
  const maxValue = 300; // Default maximum value

  const toggle = (toggleRef: React.Ref<MenuToggleElement>) => (
    <MenuToggle
      ref={toggleRef}
      onClick={() => setIsOpen(!isOpen)}
      isExpanded={isOpen}
      style={{ minWidth: '200px', width: 'fit-content' }}
    >
      {getDisplayText()}
    </MenuToggle>
  );

  const filterContent = (
    <Flex
      direction={{ default: 'column' }}
      spaceItems={{ default: 'spaceItemsSm' }}
      flexWrap={{ default: 'wrap' }}
      style={{ width: '500px', padding: '16px' }}
    >
      <FlexItem>
        <Flex alignItems={{ default: 'alignItemsCenter' }} spaceItems={{ default: 'spaceItemsXs' }}>
          <FlexItem>Min requests per second (RPS)</FlexItem>
          <FlexItem>
            <Popover
              bodyContent="Only show models that can handle at least this many requests per second (RPS). Hardware configurations performing below this value will be hidden."
              appendTo={() => document.body}
            >
              <Button
                variant="plain"
                aria-label="More info for min RPS"
                className="pf-v6-u-p-xs"
                onClick={(e) => e.stopPropagation()}
              >
                <HelpIcon />
              </Button>
            </Popover>
          </FlexItem>
        </Flex>
      </FlexItem>
      <FlexItem>
        <Flex alignItems={{ default: 'alignItemsCenter' }} spaceItems={{ default: 'spaceItemsMd' }}>
          <FlexItem flex={{ default: 'flex_1' }}>
            <Slider
              min={minValue}
              max={maxValue}
              value={localValue}
              onChange={(_, value) => {
                const clampedValue = Math.max(minValue, Math.min(maxValue, value));
                setLocalValue(clampedValue);
              }}
              isInputVisible
              inputValue={localValue}
            />
          </FlexItem>
        </Flex>
      </FlexItem>
      <FlexItem>
        <Flex spaceItems={{ default: 'spaceItemsSm' }}>
          <FlexItem>
            <Button variant="primary" onClick={handleApplyFilter}>
              Apply filter
            </Button>
          </FlexItem>
          <FlexItem>
            <Button variant="secondary" onClick={handleReset}>
              Reset
            </Button>
          </FlexItem>
        </Flex>
      </FlexItem>
    </Flex>
  );

  return (
    <Dropdown
      isOpen={isOpen}
      onOpenChange={setIsOpen}
      toggle={toggle}
      shouldFocusToggleOnSelect={false}
    >
      <div
        onMouseDown={(e) => e.stopPropagation()}
        onClick={(e) => e.stopPropagation()}
        onKeyDown={(e) => e.stopPropagation()}
        role="presentation"
      >
        {filterContent}
      </div>
    </Dropdown>
  );
};

export default MinRpsFilter;

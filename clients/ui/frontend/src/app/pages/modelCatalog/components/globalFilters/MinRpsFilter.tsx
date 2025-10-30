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
import { useCatalogNumberFilterState } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { CatalogPerformanceMetricsArtifact } from '~/app/modelCatalogTypes';
import { getDoubleValue } from '~/app/utils';

const filterKey = ModelCatalogNumberFilterKey.MIN_RPS;

type MinRpsFilterProps = {
  performanceArtifacts: CatalogPerformanceMetricsArtifact[];
};

const MinRpsFilter: React.FC<MinRpsFilterProps> = ({ performanceArtifacts }) => {
  const { value: rpsFilterValue, setValue: setRpsFilterValue } =
    useCatalogNumberFilterState(filterKey);
  const [isOpen, setIsOpen] = React.useState(false);

  // Local state for editing (initialized with current filter value or reasonable default)
  const [localValue, setLocalValue] = React.useState<number>(rpsFilterValue || 2);

  // Parse saved value if it exists
  const hasActiveFilter = rpsFilterValue !== undefined;

  const getDisplayText = (): React.ReactNode => {
    if (hasActiveFilter) {
      return (
        <>
          <strong>Min RPS:</strong> {rpsFilterValue}
        </>
      );
    }
    return 'Min RPS';
  };

  const handleApplyFilter = () => {
    // Apply the local value to the actual filter state
    setRpsFilterValue(localValue);
    setIsOpen(false);
  };

  const handleReset = () => {
    setRpsFilterValue(undefined);
    setLocalValue(minValue); // Reset to calculated minimum
    setIsOpen(false);
  };

  // Calculate min/max values from performance artifacts
  // TODO: Use real min/max values when available from API
  const { minValue, maxValue } = React.useMemo(() => {
    if (performanceArtifacts.length === 0) {
      return { minValue: 1, maxValue: 300 }; // Default values when no artifacts
    }

    const rpsValues = performanceArtifacts
      .map((artifact) => getDoubleValue(artifact.customProperties, 'requests_per_second'))
      .filter((rps) => rps > 0); // Filter out invalid values

    if (rpsValues.length === 0) {
      return { minValue: 1, maxValue: 300 }; // Default values when no valid RPS values
    }

    return {
      minValue: Math.min(...rpsValues),
      maxValue: Math.max(...rpsValues),
    };
  }, [performanceArtifacts]);

  // Ensure local value is within the calculated range
  const clampedLocalValue = Math.min(Math.max(localValue, minValue), maxValue);

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
                icon={<HelpIcon />}
              />
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
              value={clampedLocalValue}
              onChange={(_, value) => {
                const clampedValue = Math.max(minValue, Math.min(maxValue, value));
                setLocalValue(clampedValue);
              }}
              isInputVisible
              inputValue={clampedLocalValue}
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
            <Button variant="link" onClick={handleReset}>
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
      {filterContent}
    </Dropdown>
  );
};

export default MinRpsFilter;

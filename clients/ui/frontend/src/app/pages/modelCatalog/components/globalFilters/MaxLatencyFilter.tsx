import * as React from 'react';
import {
  Button,
  Dropdown,
  Flex,
  FlexItem,
  FormGroup,
  MenuToggle,
  MenuToggleElement,
  Popover,
  Select,
  SelectList,
  SelectOption,
  Slider,
} from '@patternfly/react-core';
import { HelpIcon } from '@patternfly/react-icons';
import { ModelCatalogNumberFilterKey } from '~/concepts/modelCatalog/const';
import { useCatalogNumberFilterState } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { CatalogFilterOptionsList } from '~/app/modelCatalogTypes';

const filterKey = ModelCatalogNumberFilterKey.MAX_LATENCY;

type LatencyMetric = 'E2E' | 'TTFT' | 'TPS' | 'ITL';
type LatencyPercentile = 'Mean' | 'P90' | 'P95' | 'P99';

type LatencyFilterState = {
  metric: LatencyMetric;
  percentile: LatencyPercentile;
  value: number;
};

type MaxLatencyFilterProps = {
  filterOptions?: CatalogFilterOptionsList | null;
};

const METRIC_OPTIONS: { value: LatencyMetric; label: string }[] = [
  { value: 'E2E', label: 'E2E' },
  { value: 'TTFT', label: 'TTFT' },
  { value: 'TPS', label: 'TPS' },
  { value: 'ITL', label: 'ITL' },
];

const PERCENTILE_OPTIONS: { value: LatencyPercentile; label: string }[] = [
  { value: 'Mean', label: 'Mean' },
  { value: 'P90', label: 'P90' },
  { value: 'P95', label: 'P95' },
  { value: 'P99', label: 'P99' },
];

const MaxLatencyFilter: React.FC<MaxLatencyFilterProps> = ({
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  filterOptions,
}) => {
  const { value: savedFilterValue, setValue: setSavedFilterValue } =
    useCatalogNumberFilterState(filterKey);
  const [isOpen, setIsOpen] = React.useState(false);
  const [isMetricOpen, setIsMetricOpen] = React.useState(false);
  const [isPercentileOpen, setIsPercentileOpen] = React.useState(false);

  // Local state for the filter configuration (persistent state)
  const [appliedFilter, setAppliedFilter] = React.useState<LatencyFilterState>({
    metric: 'E2E',
    percentile: 'P90',
    value: 893,
  });

  // Working state while editing the filter
  const [localFilter, setLocalFilter] = React.useState<LatencyFilterState>(appliedFilter);

  // Initialize local filter from applied filter only once when component mounts
  React.useEffect(() => {
    setLocalFilter(appliedFilter);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []); // Only run once on mount

  // Parse saved value if it exists (we'll encode metric|percentile|value as a number somehow)
  // For now, let's use the existing value as just the latency value
  const hasActiveFilter = savedFilterValue !== undefined;

  const getDisplayText = (): string => {
    if (hasActiveFilter) {
      // When there's an active filter, show the full specification with actual selected values
      return `Max latency: ${appliedFilter.metric} | ${appliedFilter.percentile} | ${savedFilterValue}ms`;
    }
    return 'Max latency';
  };

  const handleApplyFilter = () => {
    // Store the current local filter values as the applied filter
    setAppliedFilter(localFilter);
    setSavedFilterValue(localFilter.value);
    setIsOpen(false);
  };

  const handleReset = () => {
    const defaultFilter: LatencyFilterState = {
      metric: 'E2E',
      percentile: 'P90',
      value: 893,
    };
    setSavedFilterValue(undefined);
    setAppliedFilter(defaultFilter);
    setLocalFilter(defaultFilter);
    setIsOpen(false);
  };

  // Get min/max values from filter options or use defaults
  const minValue = 20; // Default minimum value
  const maxValue = 893; // Default maximum value

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
      {/* Title with help popover */}
      <FlexItem>
        <Flex alignItems={{ default: 'alignItemsCenter' }} spaceItems={{ default: 'spaceItemsXs' }}>
          <FlexItem>Max latency</FlexItem>
          <FlexItem>
            <Popover
              bodyContent="Set your maximum acceptable latency. Hardware configurations that respond slower than this value will be hidden."
              appendTo={() => document.body}
            >
              <Button
                variant="plain"
                aria-label="More info for max latency"
                className="pf-v6-u-p-xs"
                onClick={(e) => e.stopPropagation()}
              >
                <HelpIcon />
              </Button>
            </Popover>
          </FlexItem>
        </Flex>
      </FlexItem>

      {/* Metric and Percentile on the same line */}
      <FlexItem>
        <Flex spaceItems={{ default: 'spaceItemsMd' }}>
          <FlexItem flex={{ default: 'flex_1' }}>
            <FormGroup>
              <Select
                isOpen={isMetricOpen}
                selected={localFilter.metric}
                onSelect={(_, value) => {
                  if (
                    typeof value === 'string' &&
                    METRIC_OPTIONS.some((opt) => opt.value === value)
                  ) {
                    const selectedMetric = METRIC_OPTIONS.find((opt) => opt.value === value);
                    if (selectedMetric) {
                      setLocalFilter({ ...localFilter, metric: selectedMetric.value });
                    }
                  }
                  setIsMetricOpen(false);
                }}
                onOpenChange={(isSelectOpen) => {
                  setIsMetricOpen(isSelectOpen);
                  // Prevent parent dropdown from closing when this select opens/closes
                  if (isSelectOpen) {
                    setIsOpen(true);
                  }
                }}
                toggle={(toggleRef) => (
                  <MenuToggle
                    ref={toggleRef}
                    onClick={(e) => {
                      e.stopPropagation();
                      setIsMetricOpen(!isMetricOpen);
                    }}
                    isExpanded={isMetricOpen}
                    className="pf-v6-u-w-100"
                  >
                    <span>Metric: {localFilter.metric}</span>
                  </MenuToggle>
                )}
              >
                <SelectList>
                  {METRIC_OPTIONS.map((option) => (
                    <SelectOption key={option.value} value={option.value}>
                      {option.label}
                    </SelectOption>
                  ))}
                </SelectList>
              </Select>
            </FormGroup>
          </FlexItem>

          <FlexItem flex={{ default: 'flex_1' }}>
            <FormGroup>
              <Select
                isOpen={isPercentileOpen}
                selected={localFilter.percentile}
                onSelect={(_, value) => {
                  if (
                    typeof value === 'string' &&
                    PERCENTILE_OPTIONS.some((opt) => opt.value === value)
                  ) {
                    const selectedPercentile = PERCENTILE_OPTIONS.find(
                      (opt) => opt.value === value,
                    );
                    if (selectedPercentile) {
                      setLocalFilter({ ...localFilter, percentile: selectedPercentile.value });
                    }
                  }
                  setIsPercentileOpen(false);
                }}
                onOpenChange={(isSelectOpen) => {
                  setIsPercentileOpen(isSelectOpen);
                  // Prevent parent dropdown from closing when this select opens/closes
                  if (isSelectOpen) {
                    setIsOpen(true);
                  }
                }}
                toggle={(toggleRef) => (
                  <MenuToggle
                    ref={toggleRef}
                    onClick={(e) => {
                      e.stopPropagation();
                      setIsPercentileOpen(!isPercentileOpen);
                    }}
                    isExpanded={isPercentileOpen}
                    className="pf-v6-u-w-100"
                  >
                    <span>Percentile: {localFilter.percentile}</span>
                  </MenuToggle>
                )}
              >
                <SelectList>
                  {PERCENTILE_OPTIONS.map((option) => (
                    <SelectOption key={option.value} value={option.value}>
                      {option.label}
                    </SelectOption>
                  ))}
                </SelectList>
              </Select>
            </FormGroup>
          </FlexItem>
        </Flex>
      </FlexItem>

      {/* Slider with value display */}
      <FlexItem>
        <Slider
          min={minValue}
          max={maxValue}
          value={localFilter.value}
          onChange={(_, value) => {
            const clampedValue = Math.max(minValue, Math.min(maxValue, value));
            setLocalFilter({ ...localFilter, value: clampedValue });
          }}
          isInputVisible
          inputValue={localFilter.value}
          inputLabel="ms"
        />
      </FlexItem>

      {/* Buttons: Apply filter first, then Reset */}
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

export default MaxLatencyFilter;

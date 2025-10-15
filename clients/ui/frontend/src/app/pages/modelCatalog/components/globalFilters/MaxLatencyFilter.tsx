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
import {
  CatalogFilterOptionsList,
  CatalogPerformanceMetricsArtifact,
} from '~/app/modelCatalogTypes';
import { getDoubleValue } from '~/app/utils';

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
  performanceArtifacts: CatalogPerformanceMetricsArtifact[];
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

// Helper function to generate field name from metric and percentile
const getLatencyFieldName = (metric: string, percentile: string): string => {
  const metricPrefix = metric.toLowerCase();
  const percentileSuffix = percentile === 'Mean' ? '_mean' : `_${percentile.toLowerCase()}`;
  return `${metricPrefix}${percentileSuffix}`;
};

const MaxLatencyFilter: React.FC<MaxLatencyFilterProps> = ({
  filterOptions,
  performanceArtifacts,
}) => {
  const { value: savedFilterValue, setValue: setSavedFilterValue } =
    useCatalogNumberFilterState(filterKey);
  const [isOpen, setIsOpen] = React.useState(false);
  const [isMetricOpen, setIsMetricOpen] = React.useState(false);
  const [isPercentileOpen, setIsPercentileOpen] = React.useState(false);

  // Filter available metrics based on what's available in filterOptions
  const availableMetrics = React.useMemo(() => {
    if (!filterOptions?.filters) {
      return METRIC_OPTIONS; // Return all options if no filter options available
    }

    const filteredMetrics = METRIC_OPTIONS.filter((metric) =>
      // Check if at least one percentile option exists for this metric
      PERCENTILE_OPTIONS.some((percentile) => {
        const fieldName = getLatencyFieldName(metric.value, percentile.value);
        return fieldName in filterOptions.filters;
      }),
    );

    // Fallback: if no metrics are available, return all metrics to prevent empty dropdown
    return filteredMetrics.length > 0 ? filteredMetrics : METRIC_OPTIONS;
  }, [filterOptions]);

  // Filter available percentiles based on selected metric and filterOptions
  const getAvailablePercentiles = React.useCallback(
    (selectedMetric: LatencyMetric) => {
      if (!filterOptions?.filters) {
        return PERCENTILE_OPTIONS; // Return all options if no filter options available
      }

      const filteredOptions = PERCENTILE_OPTIONS.filter((percentile) => {
        const fieldName = getLatencyFieldName(selectedMetric, percentile.value);
        return fieldName in filterOptions.filters;
      });

      // Fallback: if no options are available, return all options to prevent empty dropdown
      return filteredOptions.length > 0 ? filteredOptions : PERCENTILE_OPTIONS;
    },
    [filterOptions],
  );

  // Local state for the filter configuration (persistent state)
  const [appliedFilter, setAppliedFilter] = React.useState<LatencyFilterState>(() => {
    // Initialize with first available options to ensure consistency
    const firstAvailableMetric = availableMetrics.length > 0 ? availableMetrics[0].value : 'E2E';
    const firstAvailablePercentile = getAvailablePercentiles(firstAvailableMetric);
    const defaultPercentile =
      firstAvailablePercentile.length > 0 ? firstAvailablePercentile[0].value : 'P90';

    return {
      metric: firstAvailableMetric,
      percentile: defaultPercentile,
      value: 30, // Reasonable default within typical TTFT range
    };
  });

  // Working state while editing the filter
  const [localFilter, setLocalFilter] = React.useState<LatencyFilterState>(appliedFilter);

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
    // Use first available options instead of hardcoded defaults
    const firstAvailableMetric = availableMetrics.length > 0 ? availableMetrics[0].value : 'E2E';
    const firstAvailablePercentile = getAvailablePercentiles(firstAvailableMetric);
    const defaultPercentile =
      firstAvailablePercentile.length > 0 ? firstAvailablePercentile[0].value : 'P90';

    const defaultFilter: LatencyFilterState = {
      metric: firstAvailableMetric,
      percentile: defaultPercentile,
      value: Math.min(maxValue, 100), // Use calculated maxValue but cap at reasonable level
    };

    setSavedFilterValue(undefined);
    setAppliedFilter(defaultFilter);
    setLocalFilter(defaultFilter);
    setIsOpen(false);
  };

  // Calculate min/max latency values from performance artifacts
  const { minValue, maxValue } = React.useMemo((): { minValue: number; maxValue: number } => {
    if (performanceArtifacts.length === 0) {
      return { minValue: 20, maxValue: 893 }; // Default values when no artifacts
    }

    // Get all latency values for the currently selected metric/percentile
    const fieldName = getLatencyFieldName(appliedFilter.metric, appliedFilter.percentile);
    const latencyValues = performanceArtifacts
      .map((artifact) =>
        getDoubleValue(
          artifact.customProperties,
          // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
          fieldName as keyof typeof artifact.customProperties,
        ),
      )
      .filter((latency) => latency > 0); // Filter out invalid values

    if (latencyValues.length === 0) {
      return { minValue: 20, maxValue: 893 }; // Default values when no valid latency values
    }

    return {
      minValue: Math.min(...latencyValues),
      maxValue: Math.max(...latencyValues),
    };
  }, [performanceArtifacts, appliedFilter.metric, appliedFilter.percentile]);

  // Helper to ensure value is within bounds
  const clampedLocalValue = Math.min(Math.max(localFilter.value, minValue), maxValue);

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
                icon={<HelpIcon />}
              />
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
                onClick={(e) => e.stopPropagation()}
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
                    onClick={() => setIsMetricOpen(!isMetricOpen)}
                    isExpanded={isMetricOpen}
                    className="pf-v6-u-w-100"
                  >
                    <span>Metric: {localFilter.metric}</span>
                  </MenuToggle>
                )}
              >
                <SelectList>
                  {availableMetrics.map((option) => (
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
                onClick={(e) => e.stopPropagation()}
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
                    onClick={() => setIsPercentileOpen(!isPercentileOpen)}
                    isExpanded={isPercentileOpen}
                    className="pf-v6-u-w-100"
                  >
                    <span>Percentile: {localFilter.percentile}</span>
                  </MenuToggle>
                )}
              >
                <SelectList>
                  {getAvailablePercentiles(localFilter.metric).map((option) => (
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
          value={clampedLocalValue}
          onChange={(_, value) => {
            const clampedValue = Math.max(minValue, Math.min(maxValue, value));
            setLocalFilter({ ...localFilter, value: clampedValue });
          }}
          isInputVisible
          inputValue={clampedLocalValue}
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

export default MaxLatencyFilter;

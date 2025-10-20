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
import {
  LatencyMetric,
  LatencyPercentile,
  ModelCatalogNumberFilterKey,
} from '~/concepts/modelCatalog/const';
import { useCatalogNumberFilterState } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { CatalogPerformanceMetricsArtifact } from '~/app/modelCatalogTypes';
import { getDoubleValue } from '~/app/utils';
import { getLatencyFieldName } from '~/app/pages/modelCatalog/utils/hardwareConfigurationFilterUtils';
import { useLatencyFilterConfig } from '~/app/pages/modelCatalog/utils/latencyFilterState';

const filterKey = ModelCatalogNumberFilterKey.MAX_LATENCY;

type LatencyFilterState = {
  metric: LatencyMetric;
  percentile: LatencyPercentile;
  value: number;
};

type MaxLatencyFilterProps = {
  performanceArtifacts: CatalogPerformanceMetricsArtifact[];
};

const METRIC_OPTIONS: { value: LatencyMetric; label: LatencyMetric }[] = Object.values(
  LatencyMetric,
).map((metric) => ({ value: metric, label: metric }));

const PERCENTILE_OPTIONS: { value: LatencyPercentile; label: LatencyPercentile }[] = Object.values(
  LatencyPercentile,
).map((percentile) => ({ value: percentile, label: percentile }));

const MaxLatencyFilter: React.FC<MaxLatencyFilterProps> = ({ performanceArtifacts }) => {
  const { value: savedFilterValue, setValue: setSavedFilterValue } =
    useCatalogNumberFilterState(filterKey);
  const { updateConfig: updateSharedLatencyConfig } = useLatencyFilterConfig();
  const [isOpen, setIsOpen] = React.useState(false);
  const [isMetricOpen, setIsMetricOpen] = React.useState(false);
  const [isPercentileOpen, setIsPercentileOpen] = React.useState(false);

  // Show all available metrics - in production this could be filtered based on backend data
  const availableMetrics = React.useMemo(() => METRIC_OPTIONS, []);

  // Show all available percentiles - in production this could be filtered based on backend data
  const getAvailablePercentiles = React.useCallback(() => PERCENTILE_OPTIONS, []);

  const defaultFilterState = React.useMemo(() => {
    // Initialize with first available options to ensure consistency
    const firstAvailableMetric =
      availableMetrics.length > 0 ? availableMetrics[0].value : LatencyMetric.E2E;
    const firstAvailablePercentile = getAvailablePercentiles();
    const defaultPercentile =
      firstAvailablePercentile.length > 0
        ? firstAvailablePercentile[0].value
        : LatencyPercentile.P90;

    return {
      metric: firstAvailableMetric,
      percentile: defaultPercentile,
      value: 30, // Reasonable default within typical TTFT range
    };
  }, [availableMetrics, getAvailablePercentiles]);

  // Local state for the filter configuration (persistent state)
  // TODO when we eventually move from ModelCatalogNumberFilterKey.MAX_LATENCY to using the LatencyMetricFieldNames
  // as separate filter keys, we won't need appliedFilter to be state, we can use the filterKey as the dropdown state and
  // use parseLatencyFieldName to identify the metric and the percentile from the filterKey.
  const [appliedFilter, setAppliedFilter] = React.useState<LatencyFilterState>(defaultFilterState);

  // Working state while editing the filter
  const [localFilter, setLocalFilter] = React.useState<LatencyFilterState>(appliedFilter);

  const hasActiveFilter = savedFilterValue !== undefined;

  const getDisplayText = (): React.ReactNode => {
    if (hasActiveFilter) {
      // When there's an active filter, show the full specification with actual selected values
      return (
        <>
          <strong>Max latency:</strong> {appliedFilter.metric} | {appliedFilter.percentile} |{' '}
          {savedFilterValue}ms
        </>
      );
    }
    return 'Max latency';
  };

  const handleApplyFilter = () => {
    // Store the current local filter values as the applied filter
    setAppliedFilter(localFilter);
    setSavedFilterValue(localFilter.value);
    // Update the shared config so filtering logic can use it
    updateSharedLatencyConfig({ metric: localFilter.metric, percentile: localFilter.percentile });
    setIsOpen(false);
  };

  const handleReset = () => {
    const defaultFilterWithCalculatedMax: LatencyFilterState = {
      ...defaultFilterState,
      value: Math.min(maxValue, 100), // Use calculated maxValue but cap at reasonable level
    };
    setSavedFilterValue(undefined);
    setAppliedFilter(defaultFilterWithCalculatedMax);
    setLocalFilter(defaultFilterWithCalculatedMax);
    setIsOpen(false);
  };

  // Calculate min/max latency values from performance artifacts
  const { minValue, maxValue } = React.useMemo((): { minValue: number; maxValue: number } => {
    if (performanceArtifacts.length === 0) {
      return { minValue: 20, maxValue: 893 }; // Default values when no artifacts
    }

    // Get all latency values for the currently selected metric/percentile (use localFilter for immediate updates)
    const fieldName = getLatencyFieldName(localFilter.metric, localFilter.percentile);
    const latencyValues = performanceArtifacts
      .map((artifact) => getDoubleValue(artifact.customProperties, fieldName))
      .filter((latency) => latency > 0); // Filter out invalid values

    if (latencyValues.length === 0) {
      return { minValue: 20, maxValue: 893 }; // Default values when no valid latency values
    }

    return {
      minValue: Math.round(Math.min(...latencyValues)),
      maxValue: Math.round(Math.max(...latencyValues)),
    };
  }, [performanceArtifacts, localFilter.metric, localFilter.percentile]);

  // Helper to ensure value is within bounds and rounded to integer
  const clampedLocalValue = Math.round(Math.min(Math.max(localFilter.value, minValue), maxValue));

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
      style={{ width: '550px', padding: '16px' }}
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
                  {getAvailablePercentiles().map((option) => (
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
        <div style={{ width: '100%', minWidth: '400px' }}>
          <Slider
            min={minValue}
            max={maxValue}
            value={clampedLocalValue}
            onChange={(_, value) => {
              const clampedValue = Math.round(Math.max(minValue, Math.min(maxValue, value)));
              setLocalFilter({ ...localFilter, value: clampedValue });
            }}
            isInputVisible
            inputValue={clampedLocalValue}
            inputLabel="ms"
          />
        </div>
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
      popperProps={{
        position: 'left',
        enableFlip: true,
      }}
    >
      {filterContent}
    </Dropdown>
  );
};

export default MaxLatencyFilter;

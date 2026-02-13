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
} from '@patternfly/react-core';
import { HelpIcon } from '@patternfly/react-icons';
import {
  LatencyMetric,
  LatencyPercentile,
  getLatencyFilterKey,
  ALL_LATENCY_FILTER_KEYS,
  parseLatencyFilterKey,
  LatencyMetricLabels,
  latencyMetricDescriptions,
} from '~/concepts/modelCatalog/const';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import {
  FALLBACK_LATENCY_RANGE,
  SliderRange,
  formatLatency,
} from '~/app/pages/modelCatalog/utils/performanceMetricsUtils';
import { getDefaultPerformanceFilters } from '~/app/pages/modelCatalog/utils/performanceFilterUtils';
import SliderWithInput from './SliderWithInput';

type LatencyFilterState = {
  metric: LatencyMetric;
  percentile: LatencyPercentile;
  value: number;
};
type MetricOption = {
  value: LatencyMetric;
  label: string;
  description: string;
};

const METRIC_OPTIONS: MetricOption[] = [
  LatencyMetric.TTFT,
  LatencyMetric.E2E,
  LatencyMetric.ITL,
].map((metric) => ({
  value: metric,
  label: LatencyMetricLabels[metric] ?? metric,
  description: latencyMetricDescriptions[metric] ?? '',
}));

const PERCENTILE_OPTIONS: { value: LatencyPercentile; label: LatencyPercentile }[] = Object.values(
  LatencyPercentile,
).map((percentile) => ({ value: percentile, label: percentile }));

const LatencyFilter: React.FC = () => {
  const { filterData, setFilterData, filterOptions } = React.useContext(ModelCatalogContext);
  const [isOpen, setIsOpen] = React.useState(false);
  const [isMetricOpen, setIsMetricOpen] = React.useState(false);
  const [isPercentileOpen, setIsPercentileOpen] = React.useState(false);

  // Show all available metrics - in production this could be filtered based on backend data
  const availableMetrics = React.useMemo(() => METRIC_OPTIONS, []);

  // Show all available percentiles - in production this could be filtered based on backend data
  const getAvailablePercentiles = React.useCallback(() => PERCENTILE_OPTIONS, []);

  // Find the currently active latency filter (if any)
  const currentActiveFilter = React.useMemo(() => {
    for (const metric of Object.values(LatencyMetric)) {
      for (const percentile of Object.values(LatencyPercentile)) {
        const filterKey = getLatencyFilterKey(metric, percentile);
        const value = filterData[filterKey];
        if (value !== undefined && typeof value === 'number') {
          return { fieldName: filterKey, metric, percentile, value };
        }
      }
    }
    return null;
  }, [filterData]);

  const defaultFilterState = React.useMemo(() => {
    // Find the default latency filter from namedQueries
    const defaults = getDefaultPerformanceFilters(filterOptions);
    for (const latencyKey of ALL_LATENCY_FILTER_KEYS) {
      const defaultValue = defaults[latencyKey];
      if (typeof defaultValue === 'number') {
        const { metric, percentile } = parseLatencyFilterKey(latencyKey);
        return { metric, percentile, value: defaultValue };
      }
    }
    // Fallback if no default found in namedQueries
    return {
      metric: LatencyMetric.TTFT,
      percentile: LatencyPercentile.P90,
      value: 30,
    };
  }, [filterOptions]);

  // Working state while editing the filter
  const [localFilter, setLocalFilter] = React.useState<LatencyFilterState>(() => {
    if (currentActiveFilter) {
      return {
        metric: currentActiveFilter.metric,
        percentile: currentActiveFilter.percentile,
        value: currentActiveFilter.value,
      };
    }
    return defaultFilterState;
  });

  React.useEffect(() => {
    if (isOpen) {
      // Use currentActiveFilter or defaultFilterState
      const initialState = currentActiveFilter
        ? {
            metric: currentActiveFilter.metric,
            percentile: currentActiveFilter.percentile,
            value: currentActiveFilter.value,
          }
        : defaultFilterState;
      setLocalFilter(initialState);
    }
  }, [isOpen, currentActiveFilter, defaultFilterState]);

  const { minValue, maxValue, isSliderDisabled } = React.useMemo((): SliderRange => {
    const filterKey = getLatencyFilterKey(localFilter.metric, localFilter.percentile);

    const latencyFilter = filterOptions?.filters?.[filterKey];
    if (latencyFilter && 'range' in latencyFilter && latencyFilter.range) {
      return {
        minValue: Math.round(latencyFilter.range.min ?? FALLBACK_LATENCY_RANGE.minValue),
        maxValue: Math.round(latencyFilter.range.max ?? FALLBACK_LATENCY_RANGE.maxValue),
        isSliderDisabled: false,
      };
    }
    return FALLBACK_LATENCY_RANGE;
  }, [localFilter.metric, localFilter.percentile, filterOptions]);

  // Reset value to max when metric or percentile changes (range changes)
  // This ensures the value is always valid for the current range
  const prevMetricRef = React.useRef(localFilter.metric);
  const prevPercentileRef = React.useRef(localFilter.percentile);

  React.useEffect(() => {
    const metricChanged = prevMetricRef.current !== localFilter.metric;
    const percentileChanged = prevPercentileRef.current !== localFilter.percentile;

    if (metricChanged || percentileChanged) {
      setLocalFilter((prev) => ({ ...prev, value: maxValue }));
      prevMetricRef.current = localFilter.metric;
      prevPercentileRef.current = localFilter.percentile;
    }
  }, [localFilter.metric, localFilter.percentile, maxValue]);

  const clampedValue = React.useMemo(
    () => Math.min(Math.max(localFilter.value, minValue), maxValue),
    [localFilter.value, minValue, maxValue],
  );
  const getDisplayText = (): React.ReactNode => {
    if (currentActiveFilter) {
      // When there's an active filter, show the full specification with actual selected values
      return (
        <>
          <strong>Latency:</strong> {currentActiveFilter.metric} at {currentActiveFilter.percentile}{' '}
          ≤ {formatLatency(currentActiveFilter.value)}
        </>
      );
    }
    return 'Latency';
  };

  const handleApplyFilter = () => {
    // Clear any existing latency filter
    if (currentActiveFilter) {
      setFilterData(currentActiveFilter.fieldName, undefined);
    }

    // Set the new latency filter using the dynamic filter key
    const newFilterKey = getLatencyFilterKey(localFilter.metric, localFilter.percentile);
    setFilterData(newFilterKey, localFilter.value);
    setIsOpen(false);
  };

  const handleReset = () => {
    // Reset to the filter state that was active when menu opened
    const resetState = currentActiveFilter
      ? {
          metric: currentActiveFilter.metric,
          percentile: currentActiveFilter.percentile,
          value: currentActiveFilter.value,
        }
      : defaultFilterState;
    setLocalFilter(resetState);
    // Update the refs so the useEffect doesn't think metric/percentile changed
    // This prevents the value from being reset to maxValue
    prevMetricRef.current = resetState.metric;
    prevPercentileRef.current = resetState.percentile;
  };

  const toggle = (toggleRef: React.Ref<MenuToggleElement>) => (
    <MenuToggle
      ref={toggleRef}
      data-testid="latency-filter"
      onClick={() => setIsOpen(!isOpen)}
      isExpanded={isOpen}
      isFullHeight
      style={{ minWidth: '200px', width: 'fit-content', height: '56px' }}
    >
      {getDisplayText()}
    </MenuToggle>
  );

  const filterContent = (
    <Flex
      data-testid="latency-filter-content"
      direction={{ default: 'column' }}
      spaceItems={{ default: 'spaceItemsSm' }}
      flexWrap={{ default: 'wrap' }}
      style={{ width: '550px', padding: '16px' }}
    >
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
                      setLocalFilter((prev) => ({ ...prev, metric: selectedMetric.value }));
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
                    data-testid="latency-metric-select"
                    onClick={() => setIsMetricOpen(!isMetricOpen)}
                    isExpanded={isMetricOpen}
                    className="pf-v6-u-w-100"
                  >
                    <span>Metric: {localFilter.metric}</span>
                  </MenuToggle>
                )}
              >
                <SelectList data-testid="latency-metric-options">
                  {availableMetrics.map((option) => (
                    <SelectOption
                      key={option.value}
                      value={option.value}
                      data-testid={`latency-metric-option-${option.value}`}
                      description={option.description}
                    >
                      {option.label}
                    </SelectOption>
                  ))}
                </SelectList>
              </Select>
            </FormGroup>
          </FlexItem>

          <FlexItem>
            <Flex
              alignItems={{ default: 'alignItemsCenter' }}
              spaceItems={{ default: 'spaceItemsXs' }}
            >
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
                          setLocalFilter((prev) => ({
                            ...prev,
                            percentile: selectedPercentile.value,
                          }));
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
                        data-testid="latency-percentile-select"
                        onClick={() => setIsPercentileOpen(!isPercentileOpen)}
                        isExpanded={isPercentileOpen}
                        className="pf-v6-u-w-100"
                      >
                        <span>Percentile: {localFilter.percentile}</span>
                      </MenuToggle>
                    )}
                  >
                    <SelectList data-testid="latency-percentile-options">
                      {getAvailablePercentiles().map((option) => (
                        <SelectOption
                          key={option.value}
                          value={option.value}
                          data-testid={`latency-percentile-option-${option.value}`}
                        >
                          {option.label}
                        </SelectOption>
                      ))}
                    </SelectList>
                  </Select>
                </FormGroup>
              </FlexItem>
              <FlexItem>
                <Popover
                  bodyContent={
                    <>
                      Select the latency measure used for benchmarking - percentile or mean.
                      <ul style={{ marginTop: '8px' }}>
                        <li>
                          <strong>P90, P95, P99:</strong> The selected percentage of requests must
                          meet the latency threshold.
                        </li>
                        <li>
                          <strong>Mean:</strong> The average latency across all requests.
                        </li>
                      </ul>
                    </>
                  }
                  appendTo={() => document.body}
                >
                  <Button
                    variant="plain"
                    aria-label="More info for Percentile"
                    className="pf-v6-u-p-xs"
                    icon={<HelpIcon />}
                  />
                </Popover>
              </FlexItem>
            </Flex>
          </FlexItem>
        </Flex>
      </FlexItem>

      {/* Slider with value display */}
      <FlexItem>
        <SliderWithInput
          value={clampedValue}
          min={minValue}
          max={maxValue}
          isDisabled={isSliderDisabled}
          onChange={(value) => setLocalFilter({ ...localFilter, value })}
          suffix="ms"
          ariaLabel="Latency value input"
          shouldRound
          showBoundaries={!isSliderDisabled}
          hasTooltipOverThumb={isSliderDisabled}
        />
      </FlexItem>

      {/* Buttons: Apply filter first, then Reset */}
      <FlexItem>
        <Flex spaceItems={{ default: 'spaceItemsSm' }}>
          <FlexItem>
            <Button
              data-testid="latency-apply-filter"
              variant="primary"
              onClick={handleApplyFilter}
              isDisabled={isSliderDisabled}
            >
              Apply
            </Button>
          </FlexItem>
          <FlexItem>
            <Button data-testid="latency-reset-filter" variant="link" onClick={handleReset}>
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

export default LatencyFilter;

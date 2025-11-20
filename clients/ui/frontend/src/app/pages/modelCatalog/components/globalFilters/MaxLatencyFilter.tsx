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
import { LatencyMetric, LatencyPercentile } from '~/concepts/modelCatalog/const';
import { CatalogPerformanceMetricsArtifact } from '~/app/modelCatalogTypes';
import { getDoubleValue } from '~/app/utils';
import { getLatencyFieldName } from '~/app/pages/modelCatalog/utils/hardwareConfigurationFilterUtils';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import {
  getSliderRange,
  FALLBACK_LATENCY_RANGE,
  SliderRange,
} from '~/app/pages/modelCatalog/utils/performanceMetricsUtils';
import SliderWithInput from './SliderWithInput';

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
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
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
        const fieldName = getLatencyFieldName(metric, percentile);
        const value = filterData[fieldName];
        if (value !== undefined && typeof value === 'number') {
          return { fieldName, metric, percentile, value };
        }
      }
    }
    return null;
  }, [filterData]);

  const defaultFilterState = React.useMemo(() => {
    // Initialize with first available options to ensure consistency
    const firstAvailableMetric =
      availableMetrics.length > 0 ? availableMetrics[0].value : LatencyMetric.TTFT;
    const firstAvailablePercentile = getAvailablePercentiles();
    const defaultPercentile =
      firstAvailablePercentile.length > 0
        ? firstAvailablePercentile[0].value
        : LatencyPercentile.Mean;

    return {
      metric: firstAvailableMetric,
      percentile: defaultPercentile,
      value: 30, // Reasonable default within typical TTFT range
    };
  }, [availableMetrics, getAvailablePercentiles]);

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

  // Update local filter when active filter changes
  React.useEffect(() => {
    if (currentActiveFilter) {
      setLocalFilter({
        metric: currentActiveFilter.metric,
        percentile: currentActiveFilter.percentile,
        value: currentActiveFilter.value,
      });
    }
  }, [currentActiveFilter]);

  const { minValue, maxValue, isSliderDisabled } = React.useMemo((): SliderRange => {
    const fieldName = getLatencyFieldName(localFilter.metric, localFilter.percentile);

    return getSliderRange({
      performanceArtifacts,
      getArtifactFilterValue: (artifact) => getDoubleValue(artifact.customProperties, fieldName),
      fallbackRange: FALLBACK_LATENCY_RANGE,
      shouldRound: true,
    });
  }, [performanceArtifacts, localFilter.metric, localFilter.percentile]);

  const clampedValue = React.useMemo(
    () => Math.min(Math.max(localFilter.value, minValue), maxValue),
    [localFilter.value, minValue, maxValue],
  );
  const getDisplayText = (): React.ReactNode => {
    if (currentActiveFilter) {
      // When there's an active filter, show the full specification with actual selected values
      return (
        <>
          <strong>Max latency:</strong> {currentActiveFilter.metric} |{' '}
          {currentActiveFilter.percentile} | {currentActiveFilter.value}ms
        </>
      );
    }
    return 'Max latency';
  };

  const handleApplyFilter = () => {
    // Clear any existing latency filter
    if (currentActiveFilter) {
      setFilterData(currentActiveFilter.fieldName, undefined);
    }

    // Set the new latency filter using the dynamic field name
    const newFieldName = getLatencyFieldName(localFilter.metric, localFilter.percentile);
    setFilterData(newFieldName, localFilter.value);
    setIsOpen(false);
  };

  const handleReset = () => {
    // Clear any existing latency filter
    if (currentActiveFilter) {
      setFilterData(currentActiveFilter.fieldName, undefined);
    }

    // Reset local filter to default
    setLocalFilter(defaultFilterState);
    setIsOpen(false);
  };

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
            <Button variant="primary" onClick={handleApplyFilter} isDisabled={isSliderDisabled}>
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

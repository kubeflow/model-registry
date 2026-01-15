import * as React from 'react';
import {
  Button,
  Dropdown,
  Flex,
  FlexItem,
  MenuToggle,
  MenuToggleElement,
  Popover,
} from '@patternfly/react-core';
import { HelpIcon } from '@patternfly/react-icons';
import { ModelCatalogNumberFilterKey } from '~/concepts/modelCatalog/const';
import { useCatalogNumberFilterState } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import {
  FALLBACK_RPS_RANGE,
  SliderRange,
} from '~/app/pages/modelCatalog/utils/performanceMetricsUtils';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import SliderWithInput from './SliderWithInput';

const filterKey = ModelCatalogNumberFilterKey.MAX_RPS;

const MaxRpsFilter: React.FC = () => {
  const { value: rpsFilterValue, setValue: setRpsFilterValue } =
    useCatalogNumberFilterState(filterKey);
  const { filterOptions, getPerformanceFilterDefaultValue } = React.useContext(ModelCatalogContext);
  const [isOpen, setIsOpen] = React.useState(false);

  const { minValue, maxValue, isSliderDisabled } = React.useMemo((): SliderRange => {
    // Always get range from filterOptions (which provides the full range across all artifacts)
    // Don't use performanceArtifacts since we may not have all of them in memory when paginating
    const filterValue = filterOptions?.filters?.[ModelCatalogNumberFilterKey.MAX_RPS];
    if (filterValue && 'range' in filterValue && filterValue.range) {
      return {
        minValue: filterValue.range.min ?? FALLBACK_RPS_RANGE.minValue,
        maxValue: filterValue.range.max ?? FALLBACK_RPS_RANGE.maxValue,
        isSliderDisabled: false,
      };
    }
    return FALLBACK_RPS_RANGE;
  }, [filterOptions]);

  const [localValue, setLocalValue] = React.useState<number>(() => rpsFilterValue ?? maxValue);

  const clampedValue = React.useMemo(
    () => Math.min(Math.max(localValue, minValue), maxValue),
    [localValue, minValue, maxValue],
  );

  React.useEffect(() => {
    if (isOpen) {
      setLocalValue(rpsFilterValue ?? maxValue);
    }
  }, [isOpen, rpsFilterValue, maxValue]);

  const hasActiveFilter = rpsFilterValue !== undefined;

  const getDisplayText = (): React.ReactNode => {
    if (hasActiveFilter) {
      return (
        <>
          <strong>Max RPS:</strong> {rpsFilterValue}
        </>
      );
    }
    return 'Max RPS';
  };

  const handleApplyFilter = () => {
    setRpsFilterValue(localValue);
    setIsOpen(false);
  };

  const handleReset = () => {
    // Get default value from namedQueries, fallback to maxValue
    const defaultValue = getPerformanceFilterDefaultValue(filterKey);
    const value = typeof defaultValue === 'number' ? defaultValue : maxValue;
    setRpsFilterValue(value);
    setLocalValue(value);
  };

  const toggle = (toggleRef: React.Ref<MenuToggleElement>) => (
    <MenuToggle
      ref={toggleRef}
      data-testid="max-rps-filter"
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
      direction={{ default: 'column' }}
      spaceItems={{ default: 'spaceItemsSm' }}
      flexWrap={{ default: 'wrap' }}
      style={{ width: '500px', padding: '16px' }}
    >
      <FlexItem>
        <Flex alignItems={{ default: 'alignItemsCenter' }} spaceItems={{ default: 'spaceItemsXs' }}>
          <FlexItem>Max requests per second (RPS)</FlexItem>
          <FlexItem>
            <Popover
              bodyContent="Set the maximum requests per second (RPS) target. This will be used to filter hardware configurations that can meet your throughput requirements."
              appendTo={() => document.body}
            >
              <Button
                variant="plain"
                aria-label="More info for max RPS"
                className="pf-v6-u-p-xs"
                onClick={(e) => e.stopPropagation()}
                icon={<HelpIcon />}
              />
            </Popover>
          </FlexItem>
        </Flex>
      </FlexItem>
      <FlexItem>
        <SliderWithInput
          value={clampedValue}
          min={minValue}
          max={maxValue}
          isDisabled={isSliderDisabled}
          onChange={setLocalValue}
          ariaLabel="RPS value input"
        />
      </FlexItem>
      <FlexItem>
        <Flex spaceItems={{ default: 'spaceItemsSm' }}>
          <FlexItem>
            <Button
              variant="primary"
              onClick={handleApplyFilter}
              isDisabled={isSliderDisabled}
              data-testid="max-rps-apply-filter"
            >
              Apply filter
            </Button>
          </FlexItem>
          <FlexItem>
            <Button variant="link" onClick={handleReset} data-testid="max-rps-reset-filter">
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

export default MaxRpsFilter;

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
import { CatalogPerformanceMetricsArtifact } from '~/app/modelCatalogTypes';
import { getDoubleValue } from '~/app/utils';
import {
  getSliderRange,
  FALLBACK_RPS_RANGE,
  SliderRange,
} from '~/app/pages/modelCatalog/utils/performanceMetricsUtils';
import SliderWithInput from './SliderWithInput';

const filterKey = ModelCatalogNumberFilterKey.MIN_RPS;

type MinRpsFilterProps = {
  performanceArtifacts: CatalogPerformanceMetricsArtifact[];
};

const MinRpsFilter: React.FC<MinRpsFilterProps> = ({ performanceArtifacts }) => {
  const { value: rpsFilterValue, setValue: setRpsFilterValue } =
    useCatalogNumberFilterState(filterKey);
  const [isOpen, setIsOpen] = React.useState(false);

  const { minValue, maxValue, isSliderDisabled } = React.useMemo(
    (): SliderRange =>
      getSliderRange({
        performanceArtifacts,
        getArtifactFilterValue: (artifact) =>
          getDoubleValue(artifact.customProperties, 'requests_per_second'),
        fallbackRange: FALLBACK_RPS_RANGE,
      }),
    [performanceArtifacts],
  );

  const [localValue, setLocalValue] = React.useState<number>(
    () => rpsFilterValue ?? FALLBACK_RPS_RANGE.minValue,
  );

  const clampedValue = React.useMemo(
    () => Math.min(Math.max(localValue, minValue), maxValue),
    [localValue, minValue, maxValue],
  );

  React.useEffect(() => {
    if (isOpen) {
      setLocalValue(rpsFilterValue ?? minValue);
    }
  }, [isOpen, rpsFilterValue, minValue]);

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
    setRpsFilterValue(localValue);
    setIsOpen(false);
  };

  const handleReset = () => {
    setRpsFilterValue(undefined);
    setLocalValue(minValue);
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
        <SliderWithInput
          value={clampedValue}
          min={minValue}
          max={maxValue}
          isDisabled={isSliderDisabled}
          onChange={setLocalValue}
          suffix="RPS"
          ariaLabel="RPS value input"
        />
      </FlexItem>
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
    >
      {filterContent}
    </Dropdown>
  );
};

export default MinRpsFilter;

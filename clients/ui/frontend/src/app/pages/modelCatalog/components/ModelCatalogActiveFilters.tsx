import React from 'react';
import {
  Button,
  Flex,
  FlexItem,
  Label,
  LabelGroup,
  ToolbarFilter,
  ToolbarItem,
  ToolbarLabelGroup,
  ToolbarLabel,
} from '@patternfly/react-core';
import { UndoIcon } from '@patternfly/react-icons';
import { isEnumMember } from 'mod-arch-core';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import {
  ModelCatalogStringFilterKey,
  MODEL_CATALOG_PROVIDER_NAME_MAPPING,
  MODEL_CATALOG_LICENSE_NAME_MAPPING,
  MODEL_CATALOG_TASK_NAME_MAPPING,
  AllLanguageCodesMap,
  MODEL_CATALOG_FILTER_CATEGORY_NAMES,
  MODEL_CATALOG_FILTER_CHIP_PREFIXES,
  ModelCatalogProvider,
  ModelCatalogLicense,
  ModelCatalogTask,
  AllLanguageCode,
  ModelCatalogNumberFilterKey,
  isCatalogFilterKey,
  isPerformanceFilterKey,
  parseLatencyFilterKey,
  isLatencyFilterKey,
  LatencyFilterKey,
} from '~/concepts/modelCatalog/const';
import { ModelCatalogFilterKey } from '~/app/modelCatalogTypes';
import {
  isUseCaseOptionValue,
  getUseCaseDisplayLabel,
} from '~/app/pages/modelCatalog/utils/workloadTypeUtils';
import { isValueDifferentFromDefault } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { formatLatency } from '~/app/pages/modelCatalog/utils/performanceMetricsUtils';

type ModelCatalogActiveFiltersProps = {
  filtersToShow: ModelCatalogFilterKey[];
};

/**
 * Custom close button using PatternFly's UndoIcon.
 * Used for performance filters that reset to default instead of clearing.
 */
const undoCloseButton = (onClick: (event: React.MouseEvent) => void) => (
  <Button variant="plain" aria-label="Reset to default" onClick={onClick} icon={<UndoIcon />} />
);

const ModelCatalogActiveFilters: React.FC<ModelCatalogActiveFiltersProps> = ({ filtersToShow }) => {
  const {
    filterData,
    setFilterData,
    resetSinglePerformanceFilterToDefault,
    getPerformanceFilterDefaultValue,
  } = React.useContext(ModelCatalogContext);

  const handleRemoveFilter = (categoryKey: string, labelKey: string) => {
    if (!isCatalogFilterKey(categoryKey)) {
      return;
    }

    // Performance filters always reset to default (they should always have a value)
    if (isPerformanceFilterKey(categoryKey)) {
      resetSinglePerformanceFilterToDefault(categoryKey);
      return;
    }

    // Basic filters: remove the specific value
    if (isEnumMember(categoryKey, ModelCatalogStringFilterKey)) {
      const currentValues = filterData[categoryKey];
      if (Array.isArray(currentValues)) {
        const newValues = currentValues.filter((v) => String(v) !== String(labelKey));
        setFilterData(categoryKey, newValues);
      }
    } else {
      // For number filters, clear the value
      setFilterData(categoryKey, undefined);
    }
  };

  const handleClearCategory = (categoryKey: string) => {
    if (!isCatalogFilterKey(categoryKey)) {
      return;
    }

    // Performance filters always reset to default (they should always have a value)
    if (isPerformanceFilterKey(categoryKey)) {
      resetSinglePerformanceFilterToDefault(categoryKey);
      return;
    }

    // Basic filters: clear completely
    if (isEnumMember(categoryKey, ModelCatalogStringFilterKey)) {
      setFilterData(categoryKey, []);
    } else {
      // For number filters, clear the value
      setFilterData(categoryKey, undefined);
    }
  };

  /**
   * Gets the display label for a filter value based on the filter key type
   */
  const getFilterLabel = (filterKey: ModelCatalogFilterKey, value: string | number): string => {
    // Handle string filter keys
    if (isEnumMember(filterKey, ModelCatalogStringFilterKey)) {
      const valueStr = String(value);
      switch (filterKey) {
        case ModelCatalogStringFilterKey.PROVIDER: {
          return isEnumMember(valueStr, ModelCatalogProvider)
            ? MODEL_CATALOG_PROVIDER_NAME_MAPPING[valueStr]
            : valueStr;
        }
        case ModelCatalogStringFilterKey.LICENSE: {
          return isEnumMember(valueStr, ModelCatalogLicense)
            ? MODEL_CATALOG_LICENSE_NAME_MAPPING[valueStr]
            : valueStr;
        }
        case ModelCatalogStringFilterKey.TASK: {
          return isEnumMember(valueStr, ModelCatalogTask)
            ? MODEL_CATALOG_TASK_NAME_MAPPING[valueStr]
            : valueStr;
        }
        case ModelCatalogStringFilterKey.LANGUAGE: {
          return isEnumMember(valueStr, AllLanguageCode) ? AllLanguageCodesMap[valueStr] : valueStr;
        }
        case ModelCatalogStringFilterKey.USE_CASE: {
          // Show same format as menu toggle but without bold
          // Reuse getUseCaseDisplayLabel for consistency with dropdown
          if (isUseCaseOptionValue(valueStr)) {
            return `${MODEL_CATALOG_FILTER_CHIP_PREFIXES.SCENARIO} ${getUseCaseDisplayLabel(valueStr)}`;
          }
          return valueStr;
        }
        default:
          return valueStr;
      }
    }

    // Handle number filter keys
    if (isEnumMember(filterKey, ModelCatalogNumberFilterKey)) {
      switch (filterKey) {
        // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
        case ModelCatalogNumberFilterKey.MAX_RPS:
          return `${MODEL_CATALOG_FILTER_CHIP_PREFIXES.MAX_RPS} ${value}`;
        default:
          return String(value);
      }
    }

    // Handle latency field names - kept for backwards compatibility
    const parsed = parseLatencyFilterKey(filterKey);
    const formattedValue = typeof value === 'number' ? formatLatency(value) : `${value}ms`;
    return `${parsed.metric} | ${parsed.percentile} | ${formattedValue}`;
  };

  /**
   * Checks if a filter should be skipped (no value or matches default).
   */
  const shouldSkipFilter = (filterKey: ModelCatalogFilterKey): boolean => {
    const filterValue = filterData[filterKey];

    if (!filterValue) {
      return true;
    }

    if (Array.isArray(filterValue) && filterValue.length === 0) {
      return true;
    }

    const defaultValue = getPerformanceFilterDefaultValue(filterKey);
    if (defaultValue !== undefined && !isValueDifferentFromDefault(filterValue, defaultValue)) {
      return true;
    }

    return false;
  };

  /**
   * Renders a single-chip performance filter (Workload Type, Max RPS)
   * using PatternFly Label with UndoIcon close button.
   */
  const renderSingleChipPerformanceFilter = (filterKey: ModelCatalogFilterKey) => {
    const rawValue = filterData[filterKey];
    const firstValue = Array.isArray(rawValue) ? rawValue[0] : rawValue;

    if (firstValue === undefined) {
      return null;
    }

    const labelText = getFilterLabel(filterKey, firstValue);

    return (
      <ToolbarItem key={filterKey}>
        <LabelGroup data-testid={`${filterKey}-filter-container`}>
          <Label
            data-testid={`${filterKey}-filter-chip-${firstValue}`}
            onClose={() => resetSinglePerformanceFilterToDefault(filterKey)}
            closeBtn={undoCloseButton(() => resetSinglePerformanceFilterToDefault(filterKey))}
            closeBtnAriaLabel="Reset to default"
          >
            {labelText}
          </Label>
        </LabelGroup>
      </ToolbarItem>
    );
  };

  /**
   * Renders the latency filter as a group of 3 chips (Metric, Percentile, Threshold)
   * with a single group-level UndoIcon close button.
   */
  const renderLatencyChipGroup = (filterKey: LatencyFilterKey) => {
    const filterValue = filterData[filterKey];
    const parsed = parseLatencyFilterKey(filterKey);
    const formattedValue =
      typeof filterValue === 'number' ? formatLatency(filterValue) : `${filterValue}ms`;

    const chips = [
      {
        key: 'metric',
        label: `${MODEL_CATALOG_FILTER_CHIP_PREFIXES.LATENCY_METRIC} ${parsed.metric}`,
      },
      {
        key: 'percentile',
        label: `${MODEL_CATALOG_FILTER_CHIP_PREFIXES.LATENCY_PERCENTILE} ${parsed.percentile}`,
      },
      {
        key: 'threshold',
        label: `${MODEL_CATALOG_FILTER_CHIP_PREFIXES.LATENCY_THRESHOLD} ${formattedValue}`,
      },
    ];

    return (
      <ToolbarItem key={filterKey}>
        <Flex alignItems={{ default: 'alignItemsCenter' }} spaceItems={{ default: 'spaceItemsXs' }}>
          <FlexItem>
            <LabelGroup
              categoryName={MODEL_CATALOG_FILTER_CATEGORY_NAMES[filterKey]}
              data-testid={`${filterKey}-filter-container`}
            >
              {chips.map((chip) => (
                <Label key={chip.key} data-testid={`${filterKey}-filter-chip-${chip.key}`}>
                  {chip.label}
                </Label>
              ))}
            </LabelGroup>
          </FlexItem>
          <FlexItem>
            {undoCloseButton(() => resetSinglePerformanceFilterToDefault(filterKey))}
          </FlexItem>
        </Flex>
      </ToolbarItem>
    );
  };

  return (
    <>
      {filtersToShow.map((filterKey) => {
        if (shouldSkipFilter(filterKey)) {
          return null;
        }

        // Latency: 3 chips in a group with group-level undo icon
        if (isLatencyFilterKey(filterKey)) {
          return renderLatencyChipGroup(filterKey);
        }

        // Workload Type and Max RPS: single chip with undo icon
        if (
          filterKey === ModelCatalogStringFilterKey.USE_CASE ||
          filterKey === ModelCatalogNumberFilterKey.MAX_RPS
        ) {
          return renderSingleChipPerformanceFilter(filterKey);
        }

        // Basic filters and Hardware: standard chips with X icons via ToolbarFilter
        const filterValue = filterData[filterKey];
        const filterValues = Array.isArray(filterValue) ? filterValue : [filterValue];

        const labels: ToolbarLabel[] = filterValues.map((value) => {
          const valueStr = String(value);
          const labelText = getFilterLabel(filterKey, value);
          return {
            key: valueStr,
            node: <span data-testid={`${filterKey}-filter-chip-${valueStr}`}>{labelText}</span>,
          };
        });

        const categoryLabelGroup: ToolbarLabelGroup = {
          key: filterKey,
          name: MODEL_CATALOG_FILTER_CATEGORY_NAMES[filterKey],
        };

        return (
          <ToolbarFilter
            key={filterKey}
            categoryName={categoryLabelGroup}
            labels={labels}
            deleteLabel={(category, label) => {
              const categoryKeyValue = typeof category === 'string' ? category : category.key;
              const labelKey = typeof label === 'string' ? label : label.key;
              handleRemoveFilter(categoryKeyValue, labelKey);
            }}
            deleteLabelGroup={(category) => {
              const categoryKeyValue = typeof category === 'string' ? category : category.key;
              handleClearCategory(categoryKeyValue);
            }}
            data-testid={`${filterKey}-filter-container`}
          >
            {null}
          </ToolbarFilter>
        );
      })}
    </>
  );
};

export default ModelCatalogActiveFilters;

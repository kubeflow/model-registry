import React from 'react';
import { ToolbarFilter, ToolbarLabelGroup, ToolbarLabel } from '@patternfly/react-core';
import { isEnumMember } from 'mod-arch-core';
import { Theme } from 'mod-arch-kubeflow';
import { STYLE_THEME } from '~/app/utilities/const';
import './ModelCatalogActiveFilters.css';
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

const ModelCatalogActiveFilters: React.FC<ModelCatalogActiveFiltersProps> = ({ filtersToShow }) => {
  const {
    filterData,
    setFilterData,
    resetSinglePerformanceFilterToDefault,
    getPerformanceFilterDefaultValue,
  } = React.useContext(ModelCatalogContext);

  const isPatternfly = STYLE_THEME === Theme.Patternfly;

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

  const createLatencyChipLabels = (
    filterKey: LatencyFilterKey,
    parsed: ReturnType<typeof parseLatencyFilterKey>,
    formattedValue: string,
  ): ToolbarLabel[] => [
    {
      key: `${filterKey}-metric`,
      node: (
        <span data-testid={`${filterKey}-filter-chip-metric`} data-has-default="true">
          {MODEL_CATALOG_FILTER_CHIP_PREFIXES.LATENCY_METRIC} {parsed.metric}
        </span>
      ),
    },
    {
      key: `${filterKey}-percentile`,
      node: (
        <span data-testid={`${filterKey}-filter-chip-percentile`} data-has-default="true">
          {MODEL_CATALOG_FILTER_CHIP_PREFIXES.LATENCY_PERCENTILE} {parsed.percentile}
        </span>
      ),
    },
    {
      key: `${filterKey}-threshold`,
      node: (
        <span data-testid={`${filterKey}-filter-chip-threshold`} data-has-default="true">
          {MODEL_CATALOG_FILTER_CHIP_PREFIXES.LATENCY_THRESHOLD} {formattedValue}
        </span>
      ),
    },
  ];

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
    // TODO: Remove this condition if we add more number filter keys
    if (isEnumMember(filterKey, ModelCatalogNumberFilterKey)) {
      switch (filterKey) {
        // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
        case ModelCatalogNumberFilterKey.MAX_RPS:
          return `${MODEL_CATALOG_FILTER_CHIP_PREFIXES.MAX_RPS} ${value}`;
        default:
          return String(value);
      }
    }

    // Handle latency field names - type is already narrowed to LatencyMetricFieldName
    // Latency chips are handled separately - this shouldn't be called for latency filters
    // but keeping for backwards compatibility
    const parsed = parseLatencyFilterKey(filterKey);
    const formattedValue = typeof value === 'number' ? formatLatency(value) : `${value}ms`;
    return `${parsed.metric} | ${parsed.percentile} | ${formattedValue}`;
  };

  return (
    <>
      {filtersToShow.map((filterKey) => {
        const filterValue = filterData[filterKey];

        // Skip if no value is set
        if (!filterValue) {
          return null;
        }

        // For array values (string filters), skip if empty
        if (Array.isArray(filterValue) && filterValue.length === 0) {
          return null;
        }

        // For any filter with a default value, skip if value matches the default
        // This ensures consistent behavior across all pages (landing, details, etc.)
        const defaultValue = getPerformanceFilterDefaultValue(filterKey);
        if (defaultValue !== undefined) {
          if (!isValueDifferentFromDefault(filterValue, defaultValue)) {
            return null;
          }
        }

        // Special handling for latency filters - show as 3 separate chips
        if (isLatencyFilterKey(filterKey)) {
          // Type narrowing: filterKey is now LatencyFilterKey
          const latencyFilterKey: LatencyFilterKey = filterKey;
          const parsed = parseLatencyFilterKey(latencyFilterKey);
          const formattedValue =
            typeof filterValue === 'number' ? formatLatency(filterValue) : `${filterValue}ms`;

          // Create 3 separate chips for latency group using extracted helper
          const latencyLabels = createLatencyChipLabels(latencyFilterKey, parsed, formattedValue);

          const categoryLabelGroup: ToolbarLabelGroup = {
            key: filterKey,
            name: MODEL_CATALOG_FILTER_CATEGORY_NAMES[filterKey],
          };

          // For latency chips: clicking any individual chip OR the group reset should reset the entire filter
          return (
            <ToolbarFilter
              key={filterKey}
              categoryName={categoryLabelGroup}
              labels={latencyLabels}
              deleteLabel={(category) => {
                // Individual chip click - reset entire latency filter to default
                // Note: label param not used since all latency chips reset the entire group
                const categoryKeyValue = typeof category === 'string' ? category : category.key;
                handleClearCategory(categoryKeyValue);
              }}
              deleteLabelGroup={(category) => {
                // Group reset button click - reset entire latency filter to default
                const categoryKeyValue = typeof category === 'string' ? category : category.key;
                handleClearCategory(categoryKeyValue);
              }}
              data-testid={`${filterKey}-filter-container`}
            >
              {null}
            </ToolbarFilter>
          );
        }

        // Normalize to array for consistent handling (non-latency filters)
        const filterValues = Array.isArray(filterValue) ? filterValue : [filterValue];

        const categoryName = MODEL_CATALOG_FILTER_CATEGORY_NAMES[filterKey];

        // Check if this filter has a default value AND is a performance filter
        // If so, the filter group gets special styling (fa-undo on group, no X on labels)
        // This indicates to the user that clicking will reset to default, not clear
        // Note: HARDWARE_CONFIGURATION is not a performance filter, so it should show normal X icons
        const filterHasDefault =
          isPatternfly &&
          isPerformanceFilterKey(filterKey) &&
          getPerformanceFilterDefaultValue(filterKey) !== undefined;

        // Build labels for ToolbarFilter
        const labels: ToolbarLabel[] = filterValues.map((value) => {
          const valueStr = String(value);
          const labelText = getFilterLabel(filterKey, value);
          return {
            key: valueStr,
            node: (
              <span
                data-testid={`${filterKey}-filter-chip-${valueStr}`}
                {...(filterHasDefault && { 'data-has-default': 'true' })}
              >
                {labelText}
              </span>
            ),
          };
        });

        const categoryLabelGroup: ToolbarLabelGroup = {
          key: filterKey,
          name: categoryName,
        };

        // Use ToolbarFilter for all filters (both basic and performance)
        // This ensures proper integration with Toolbar's clearAllFilters button
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
            {/* ToolbarFilter requires children but we only render labels, not filter controls */}
            {null}
          </ToolbarFilter>
        );
      })}
    </>
  );
};

export default ModelCatalogActiveFilters;

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
  ModelCatalogProvider,
  ModelCatalogLicense,
  ModelCatalogTask,
  AllLanguageCode,
  ModelCatalogNumberFilterKey,
  isCatalogFilterKey,
  isPerformanceFilterKey,
  parseLatencyFilterKey,
} from '~/concepts/modelCatalog/const';
import { ModelCatalogFilterKey } from '~/app/modelCatalogTypes';
import {
  USE_CASE_OPTIONS,
  isUseCaseOptionValue,
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
          if (isUseCaseOptionValue(valueStr)) {
            const option = USE_CASE_OPTIONS.find((opt) => opt.value === valueStr);
            if (option) {
              return `${option.label} (${option.inputTokens} input | ${option.outputTokens} output tokens)`;
            }
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
          return `≤${value} requests/sec`;
        default:
          return String(value);
      }
    }

    // Handle latency field names - type is already narrowed to LatencyMetricFieldName
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

        // Normalize to array for consistent handling
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

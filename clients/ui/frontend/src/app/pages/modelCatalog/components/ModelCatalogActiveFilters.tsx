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

    if (isPerformanceFilterKey(categoryKey)) {
      resetSinglePerformanceFilterToDefault(categoryKey);
      return;
    }

    if (isEnumMember(categoryKey, ModelCatalogStringFilterKey)) {
      const currentValues = filterData[categoryKey];
      if (Array.isArray(currentValues)) {
        const newValues = currentValues.filter((v) => String(v) !== String(labelKey));
        setFilterData(categoryKey, newValues);
      }
    } else {
      setFilterData(categoryKey, undefined);
    }
  };

  const handleClearCategory = (categoryKey: string) => {
    if (!isCatalogFilterKey(categoryKey)) {
      return;
    }

    if (isPerformanceFilterKey(categoryKey)) {
      resetSinglePerformanceFilterToDefault(categoryKey);
      return;
    }

    if (isEnumMember(categoryKey, ModelCatalogStringFilterKey)) {
      setFilterData(categoryKey, []);
    } else {
      setFilterData(categoryKey, undefined);
    }
  };

  /**
   * Gets the display label for a filter value based on the filter key type
   */
  const getFilterLabel = (filterKey: ModelCatalogFilterKey, value: string | number): string => {
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
          if (isUseCaseOptionValue(valueStr)) {
            return `${MODEL_CATALOG_FILTER_CHIP_PREFIXES.WORKLOAD_TYPE} ${getUseCaseDisplayLabel(valueStr)}`;
          }
          return valueStr;
        }
        default:
          return valueStr;
      }
    }

    if (isEnumMember(filterKey, ModelCatalogNumberFilterKey)) {
      switch (filterKey) {
        // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
        case ModelCatalogNumberFilterKey.MAX_RPS:
          return `${MODEL_CATALOG_FILTER_CHIP_PREFIXES.MAX_RPS} ${value}`;
        default:
          return String(value);
      }
    }

    const parsed = parseLatencyFilterKey(filterKey);
    const formattedValue = typeof value === 'number' ? formatLatency(value) : `${value}ms`;
    return `${parsed.metric} | ${parsed.percentile} | ${formattedValue}`;
  };

  return (
    <>
      {filtersToShow.map((filterKey) => {
        const filterValue = filterData[filterKey];

        if (!filterValue) {
          return null;
        }

        if (Array.isArray(filterValue) && filterValue.length === 0) {
          return null;
        }

        const defaultValue = getPerformanceFilterDefaultValue(filterKey);
        if (defaultValue !== undefined && !isValueDifferentFromDefault(filterValue, defaultValue)) {
          return null;
        }

        // Performance filter chips use data-has-default to trigger undo icon styling via CSS
        const filterHasDefault =
          isPatternfly &&
          isPerformanceFilterKey(filterKey) &&
          getPerformanceFilterDefaultValue(filterKey) !== undefined;

        // Latency: 3 separate chips in a group
        if (isLatencyFilterKey(filterKey)) {
          const latencyFilterKey: LatencyFilterKey = filterKey;
          const parsed = parseLatencyFilterKey(latencyFilterKey);
          const formattedValue =
            typeof filterValue === 'number' ? formatLatency(filterValue) : `${filterValue}ms`;

          const latencyLabels: ToolbarLabel[] = [
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

          return (
            <ToolbarFilter
              key={filterKey}
              categoryName={{
                key: filterKey,
                name: MODEL_CATALOG_FILTER_CATEGORY_NAMES[filterKey],
              }}
              labels={latencyLabels}
              deleteLabel={(category) => {
                const categoryKeyValue = typeof category === 'string' ? category : category.key;
                handleClearCategory(categoryKeyValue);
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
        }

        // All other filters
        const filterValues = Array.isArray(filterValue) ? filterValue : [filterValue];

        // Single-value performance filters (Workload Type, Max RPS) use "single" to show
        // the undo icon on the chip's own close button (inside the chip).
        // Multi-chip groups (Latency) use "group" to show undo on the group close button.
        const isSingleValuePerformanceFilter =
          filterKey === ModelCatalogStringFilterKey.USE_CASE ||
          filterKey === ModelCatalogNumberFilterKey.MAX_RPS;

        const hasDefaultAttr = filterHasDefault
          ? isSingleValuePerformanceFilter
            ? 'single'
            : 'group'
          : undefined;

        const labels: ToolbarLabel[] = filterValues.map((value) => {
          const valueStr = String(value);
          const labelText = getFilterLabel(filterKey, value);
          return {
            key: valueStr,
            node: (
              <span
                data-testid={`${filterKey}-filter-chip-${valueStr}`}
                {...(hasDefaultAttr && { 'data-has-default': hasDefaultAttr })}
              >
                {labelText}
              </span>
            ),
          };
        });

        // Empty name removes the category box/border. PF's LabelGroup only applies
        // the category modifier class when categoryName is truthy.
        const categoryLabelGroup: ToolbarLabelGroup = {
          key: filterKey,
          name: isSingleValuePerformanceFilter
            ? ''
            : MODEL_CATALOG_FILTER_CATEGORY_NAMES[filterKey],
        };

        // Single-value performance filters: use deleteLabel only (undo icon on chip itself).
        // Multi-chip/basic filters: use both deleteLabel and deleteLabelGroup.
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
            {...(!isSingleValuePerformanceFilter && {
              deleteLabelGroup: (category: string | ToolbarLabelGroup) => {
                const categoryKeyValue = typeof category === 'string' ? category : category.key;
                handleClearCategory(categoryKeyValue);
              },
            })}
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

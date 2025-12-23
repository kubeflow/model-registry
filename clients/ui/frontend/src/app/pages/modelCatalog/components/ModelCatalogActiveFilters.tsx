import React from 'react';
import { ToolbarFilter, ToolbarLabelGroup, ToolbarLabel } from '@patternfly/react-core';
import { isEnumMember } from 'mod-arch-core';
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
} from '~/concepts/modelCatalog/const';
import { ModelCatalogFilterKey } from '~/app/modelCatalogTypes';
import { parseLatencyFieldName } from '~/app/pages/modelCatalog/utils/hardwareConfigurationFilterUtils';

type ModelCatalogActiveFiltersProps = {
  filtersToShow: ModelCatalogFilterKey[];
};

const ModelCatalogActiveFilters: React.FC<ModelCatalogActiveFiltersProps> = ({ filtersToShow }) => {
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);

  const handleRemoveFilter = (categoryKey: string, labelKey: string) => {
    if (!isCatalogFilterKey(categoryKey)) {
      return;
    }
    if (isEnumMember(categoryKey, ModelCatalogStringFilterKey)) {
      const currentValues = filterData[categoryKey];
      if (Array.isArray(currentValues)) {
        const newValues = currentValues.filter((v) => String(v) !== String(labelKey));
        setFilterData(categoryKey, newValues);
      }
    } else {
      // For number filters and latency fields, clear the value
      setFilterData(categoryKey, undefined);
    }
  };

  const handleClearCategory = (categoryKey: string) => {
    if (!isCatalogFilterKey(categoryKey)) {
      return;
    }
    if (isEnumMember(categoryKey, ModelCatalogStringFilterKey)) {
      setFilterData(categoryKey, []);
    } else {
      // For number filters and latency fields, clear the value
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
        default:
          return valueStr;
      }
    }

    // Handle number filter keys
    // TODO: Remove this condition if we add more number filter keys
    if (isEnumMember(filterKey, ModelCatalogNumberFilterKey)) {
      switch (filterKey) {
        // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
        case ModelCatalogNumberFilterKey.MIN_RPS:
          return `≥${value} requests/sec`;
        default:
          return String(value);
      }
    }

    // Handle latency field names - type is already narrowed to LatencyMetricFieldName
    const parsed = parseLatencyFieldName(filterKey);
    if (parsed) {
      return `${parsed.metric} ${parsed.percentile}: ≤${value}ms`;
    }
    return `${filterKey}: ≤${value}ms`;
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

        // Normalize to array for consistent handling
        const filterValues = Array.isArray(filterValue) ? filterValue : [filterValue];

        const categoryName = MODEL_CATALOG_FILTER_CATEGORY_NAMES[filterKey];

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
          name: categoryName,
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
            {/* ToolbarFilter requires children but we only render labels, not filter controls */}
            {null}
          </ToolbarFilter>
        );
      })}
    </>
  );
};

export default ModelCatalogActiveFilters;

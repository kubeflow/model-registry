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

  const filterHasDefaultByName = (filterLabel?: string) => filterLabel === 'Task'; // TODO somehow look up the filter key by name and check if it has defaults in the filter_options

  // Extremely hacky way to replace the default close icon with a redo icon on each render
  // Note that this also requires adding className="model-catalog-filter-toolbar" to the ancestor Toolbar component everywhere we render ModelCatalogActiveFilters
  React.useLayoutEffect(() => {
    setTimeout(() => {
      document
        .querySelectorAll('.model-catalog-filter-toolbar .pf-v6-c-toolbar__item.pf-m-label-group')
        .forEach((element) => {
          const filterLabel = element.querySelector('.pf-v6-c-label-group__label')?.textContent;
          const filterHasDefault = filterHasDefaultByName(filterLabel);
          if (filterHasDefault) {
            // Replace the default close icon with a redo icon if the filter is resettable to default
            const closeIcons = element.querySelectorAll('.pf-v6-c-button__icon');
            closeIcons.forEach((icon) => {
              icon.setHTMLUnsafe(
                '<svg class="pf-v6-svg" viewBox="0 0 512 512" fill="currentColor" aria-hidden="true" role="img" width="1em" height="1em"><path d="M500.33 0h-47.41a12 12 0 0 0-12 12.57l4 82.76A247.42 247.42 0 0 0 256 8C119.34 8 7.9 119.53 8 256.19 8.1 393.07 119.1 504 256 504a247.1 247.1 0 0 0 166.18-63.91 12 12 0 0 0 .48-17.43l-34-34a12 12 0 0 0-16.38-.55A176 176 0 1 1 402.1 157.8l-101.53-4.87a12 12 0 0 0-12.57 12v47.41a12 12 0 0 0 12 12h200.33a12 12 0 0 0 12-12V12a12 12 0 0 0-12-12z"></path></svg>',
              );
            });
          }
        });
    }, 0);
  });
  // Note also that we need to entirely prevent rendering ToolbarFilter for filters whose values are already set to default. the chips only represent things the user has changed.

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

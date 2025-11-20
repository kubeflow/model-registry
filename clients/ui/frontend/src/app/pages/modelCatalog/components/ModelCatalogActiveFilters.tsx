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
} from '~/concepts/modelCatalog/const';
import { ModelCatalogFilterKey } from '~/app/modelCatalogTypes';

type ModelCatalogActiveFiltersProps = {
  filtersToShow: ModelCatalogFilterKey[];
};

const ModelCatalogActiveFilters: React.FC<ModelCatalogActiveFiltersProps> = ({ filtersToShow }) => {
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);

  const handleRemoveFilter = (categoryKey: string, labelKey: string) => {
    if (isEnumMember(categoryKey, ModelCatalogStringFilterKey)) {
      const currentValues = filterData[categoryKey];
      if (Array.isArray(currentValues)) {
        const newValues = currentValues.filter((v) => String(v) !== String(labelKey));
        setFilterData(categoryKey, newValues);
      }
    }
  };

  const handleClearCategory = (categoryKey: string) => {
    if (isEnumMember(categoryKey, ModelCatalogStringFilterKey)) {
      setFilterData(categoryKey, []);
    }
  };

  const getFilterLabel = (filterKey: ModelCatalogStringFilterKey, value: string): string => {
    switch (filterKey) {
      case ModelCatalogStringFilterKey.PROVIDER: {
        return isEnumMember(value, ModelCatalogProvider)
          ? MODEL_CATALOG_PROVIDER_NAME_MAPPING[value]
          : value;
      }
      case ModelCatalogStringFilterKey.LICENSE: {
        return isEnumMember(value, ModelCatalogLicense)
          ? MODEL_CATALOG_LICENSE_NAME_MAPPING[value]
          : value;
      }
      case ModelCatalogStringFilterKey.TASK: {
        return isEnumMember(value, ModelCatalogTask)
          ? MODEL_CATALOG_TASK_NAME_MAPPING[value]
          : value;
      }
      case ModelCatalogStringFilterKey.LANGUAGE: {
        return isEnumMember(value, AllLanguageCode) ? AllLanguageCodesMap[value] : value;
      }
      default:
        return value;
    }
  };

  return (
    <>
      {filtersToShow.map((filterKey) => {
        // Only process string filter keys that are arrays
        if (!isEnumMember(filterKey, ModelCatalogStringFilterKey)) {
          return null;
        }

        const filterValues = filterData[filterKey];
        if (!Array.isArray(filterValues) || filterValues.length === 0) {
          return null;
        }

        const categoryName = MODEL_CATALOG_FILTER_CATEGORY_NAMES[filterKey];
        const labels: ToolbarLabel[] = filterValues.map((value) => {
          const valueStr = String(value);
          const labelText = getFilterLabel(filterKey, valueStr);
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
              const categoryKey = typeof category === 'string' ? category : category.key;
              const labelKey = typeof label === 'string' ? label : label.key;
              handleRemoveFilter(categoryKey, labelKey);
            }}
            deleteLabelGroup={(category) => {
              const categoryKey = typeof category === 'string' ? category : category.key;
              handleClearCategory(categoryKey);
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

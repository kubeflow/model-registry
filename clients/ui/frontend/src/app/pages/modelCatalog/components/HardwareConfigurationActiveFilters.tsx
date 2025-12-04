import React from 'react';
import { ToolbarFilter, ToolbarLabelGroup, ToolbarLabel } from '@patternfly/react-core';
import { isEnumMember } from 'mod-arch-core';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import {
  ModelCatalogStringFilterKey,
  ModelCatalogNumberFilterKey,
  MODEL_CATALOG_FILTER_CATEGORY_NAMES,
  LatencyMetric,
  LatencyPercentile,
  type LatencyMetricFieldName,
} from '~/concepts/modelCatalog/const';
import {
  getLatencyFieldName,
  parseLatencyFieldName,
  type LatencyFilterConfig,
} from '~/app/pages/modelCatalog/utils/hardwareConfigurationFilterUtils';
import type { ModelCatalogFilterStates } from '~/app/modelCatalogTypes';

const HardwareConfigurationActiveFilters: React.FC = () => {
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);

  const handleRemoveFilter = (categoryKey: string, labelKey: string) => {
    if (isEnumMember(categoryKey, ModelCatalogStringFilterKey)) {
      const currentValues = filterData[categoryKey];
      if (Array.isArray(currentValues)) {
        const newValues = currentValues.filter((v) => String(v) !== String(labelKey));
        setFilterData(categoryKey, newValues);
      }
    } else if (isEnumMember(categoryKey, ModelCatalogNumberFilterKey)) {
      setFilterData(categoryKey, undefined);
    } else if (categoryKey === 'maxLatency') {
      // Remove specific latency field
      setFilterData(labelKey as keyof ModelCatalogFilterStates, undefined);
    }
  };

  const handleClearCategory = (categoryKey: string) => {
    if (isEnumMember(categoryKey, ModelCatalogStringFilterKey)) {
      setFilterData(categoryKey, []);
    } else if (isEnumMember(categoryKey, ModelCatalogNumberFilterKey)) {
      setFilterData(categoryKey, undefined);
    } else if (categoryKey === 'maxLatency') {
      // Clear all latency filters
      for (const metric of Object.values(LatencyMetric)) {
        for (const percentile of Object.values(LatencyPercentile)) {
          const fieldName = getLatencyFieldName(metric, percentile);
          setFilterData(fieldName, undefined);
        }
      }
    }
  };

  // Collect active latency filters
  const activeLatencyFilters: Array<{ fieldName: LatencyMetricFieldName; value: number }> = [];
  for (const metric of Object.values(LatencyMetric)) {
    for (const percentile of Object.values(LatencyPercentile)) {
      const fieldName = getLatencyFieldName(metric, percentile);
      const filterValue = filterData[fieldName];
      if (filterValue !== undefined && typeof filterValue === 'number') {
        activeLatencyFilters.push({ fieldName, value: filterValue });
      }
    }
  }

  // Helper to format latency label
  const formatLatencyLabel = (fieldName: LatencyMetricFieldName, value: number): string => {
    const parsed = parseLatencyFieldName(fieldName);
    if (parsed) {
      return `${parsed.metric} ${parsed.percentile}: ≤${value}ms`;
    }
    return `${fieldName}: ≤${value}ms`;
  };

  return (
    <>
      {/* USE_CASE (Workload) Filter */}
      {(() => {
        const filterValues = filterData[ModelCatalogStringFilterKey.USE_CASE];
        if (!Array.isArray(filterValues) || filterValues.length === 0) {
          return null;
        }

        const categoryName = MODEL_CATALOG_FILTER_CATEGORY_NAMES[ModelCatalogStringFilterKey.USE_CASE];
        const labels: ToolbarLabel[] = filterValues.map((value) => {
          const valueStr = String(value);
          return {
            key: valueStr,
            node: <span data-testid={`use-case-filter-chip-${valueStr}`}>{valueStr}</span>,
          };
        });

        const categoryLabelGroup: ToolbarLabelGroup = {
          key: ModelCatalogStringFilterKey.USE_CASE,
          name: categoryName,
        };

        return (
          <ToolbarFilter
            key={ModelCatalogStringFilterKey.USE_CASE}
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
            data-testid="use-case-filter-container"
          >
            {null}
          </ToolbarFilter>
        );
      })()}

      {/* MAX LATENCY Filters */}
      {activeLatencyFilters.length > 0 && (
        <ToolbarFilter
          key="maxLatency"
          categoryName={{
            key: 'maxLatency',
            name: 'Max latency',
          }}
          labels={activeLatencyFilters.map(({ fieldName, value }) => ({
            key: fieldName,
            node: (
              <span data-testid={`max-latency-filter-chip-${fieldName}`}>
                {formatLatencyLabel(fieldName, value)}
              </span>
            ),
          }))}
          deleteLabel={(category, label) => {
            const categoryKey = typeof category === 'string' ? category : category.key;
            const labelKey = typeof label === 'string' ? label : label.key;
            handleRemoveFilter(categoryKey, labelKey);
          }}
          deleteLabelGroup={(category) => {
            const categoryKey = typeof category === 'string' ? category : category.key;
            handleClearCategory(categoryKey);
          }}
          data-testid="max-latency-filter-container"
        >
          {null}
        </ToolbarFilter>
      )}

      {/* MIN_RPS Filter */}
      {(() => {
        const filterValue = filterData[ModelCatalogNumberFilterKey.MIN_RPS];
        if (filterValue === undefined) {
          return null;
        }

        const labels: ToolbarLabel[] = [
          {
            key: String(filterValue),
            node: (
              <span data-testid="min-rps-filter-chip">≥{filterValue} requests/sec</span>
            ),
          },
        ];

        return (
          <ToolbarFilter
            key={ModelCatalogNumberFilterKey.MIN_RPS}
            categoryName={{
              key: ModelCatalogNumberFilterKey.MIN_RPS,
              name: 'Min RPS',
            }}
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
            data-testid="min-rps-filter-container"
          >
            {null}
          </ToolbarFilter>
        );
      })()}

      {/* HARDWARE_TYPE Filter */}
      {(() => {
        const filterValues = filterData[ModelCatalogStringFilterKey.HARDWARE_TYPE];
        if (!Array.isArray(filterValues) || filterValues.length === 0) {
          return null;
        }

        const categoryName = MODEL_CATALOG_FILTER_CATEGORY_NAMES[ModelCatalogStringFilterKey.HARDWARE_TYPE];
        const labels: ToolbarLabel[] = filterValues.map((value) => {
          const valueStr = String(value);
          return {
            key: valueStr,
            node: <span data-testid={`hardware-type-filter-chip-${valueStr}`}>{valueStr}</span>,
          };
        });

        const categoryLabelGroup: ToolbarLabelGroup = {
          key: ModelCatalogStringFilterKey.HARDWARE_TYPE,
          name: categoryName,
        };

        return (
          <ToolbarFilter
            key={ModelCatalogStringFilterKey.HARDWARE_TYPE}
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
            data-testid="hardware-type-filter-container"
          >
            {null}
          </ToolbarFilter>
        );
      })()}
    </>
  );
};

export default HardwareConfigurationActiveFilters;


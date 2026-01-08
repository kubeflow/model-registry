import * as React from 'react';
import {
  Button,
  Popover,
  Toolbar,
  ToolbarContent,
  ToolbarGroup,
  ToolbarItem,
} from '@patternfly/react-core';
import { HelpIcon } from '@patternfly/react-icons';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import {
  CatalogPerformanceMetricsArtifact,
  ModelCatalogFilterKey,
  ModelCatalogFilterStates,
} from '~/app/modelCatalogTypes';
import { hasFiltersApplied } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import {
  ModelCatalogStringFilterKey,
  ModelCatalogNumberFilterKey,
  ALL_LATENCY_FIELD_NAMES,
} from '~/concepts/modelCatalog/const';
import WorkloadTypeFilter from './globalFilters/WorkloadTypeFilter';
import HardwareTypeFilter from './globalFilters/HardwareTypeFilter';
import MaxRpsFilter from './globalFilters/MaxRpsFilter';
import LatencyFilter from './globalFilters/LatencyFilter';
import ModelCatalogActiveFilters from './ModelCatalogActiveFilters';

type HardwareConfigurationFilterToolbarProps = {
  performanceArtifacts?: CatalogPerformanceMetricsArtifact[];
};

/**
 * Basic filter keys that appear on the catalog landing page.
 */
const BASIC_FILTER_KEYS: ModelCatalogFilterKey[] = [
  ModelCatalogStringFilterKey.PROVIDER,
  ModelCatalogStringFilterKey.LICENSE,
  ModelCatalogStringFilterKey.TASK,
  ModelCatalogStringFilterKey.LANGUAGE,
];

/**
 * Performance filter keys that are shown on the performance/hardware configuration page.
 */
const PERFORMANCE_FILTER_KEYS: ModelCatalogFilterKey[] = [
  ModelCatalogStringFilterKey.USE_CASE,
  ModelCatalogStringFilterKey.HARDWARE_TYPE,
  ModelCatalogNumberFilterKey.MAX_RPS,
];

/**
 * Gets all filter keys including basic filters, performance filters, and any active latency filters.
 * When performance view is enabled, we show both basic and performance filter chips together.
 */
const getAllActiveFilterKeys = (filterData: ModelCatalogFilterStates): ModelCatalogFilterKey[] => {
  const activeLatencyKeys = ALL_LATENCY_FIELD_NAMES.filter((key) => filterData[key] !== undefined);
  return [...BASIC_FILTER_KEYS, ...PERFORMANCE_FILTER_KEYS, ...activeLatencyKeys];
};

const HardwareConfigurationFilterToolbar: React.FC<HardwareConfigurationFilterToolbarProps> = ({
  performanceArtifacts,
}) => {
  const {
    filterOptions,
    filterOptionsLoaded,
    filterOptionsLoadError,
    filterData,
    resetPerformanceFiltersToDefaults,
  } = React.useContext(ModelCatalogContext);

  // Get all filter keys including basic filters, performance filters, and active latency filters
  const filtersToShow = React.useMemo(() => getAllActiveFilterKeys(filterData), [filterData]);

  // Check if any performance filters are active (only checking performance-related filters)
  const hasActiveFilters = React.useMemo(
    () => hasFiltersApplied(filterData, filtersToShow),
    [filterData, filtersToShow],
  );

  if (!filterOptionsLoaded || filterOptionsLoadError || !filterOptions) {
    return null;
  }

  const handleClearAllFilters = () => {
    // Reset performance filters to default values from namedQueries
    resetPerformanceFiltersToDefaults();
  };

  return (
    <Toolbar
      key={`toolbar-${hasActiveFilters}`}
      clearAllFilters={handleClearAllFilters}
      clearFiltersButtonText={hasActiveFilters ? 'Reset all filters' : ''}
    >
      <ToolbarContent>
        <ToolbarGroup>
          <ToolbarItem>
            <WorkloadTypeFilter />
            <Popover
              bodyContent="Select a workload type to view performance under specific input and output token lengths."
              appendTo={() => document.body}
            >
              <Button
                variant="plain"
                aria-label="More info for workload type"
                className="pf-v6-u-p-xs"
                icon={<HelpIcon />}
              />
            </Popover>
          </ToolbarItem>
          <ToolbarItem variant="separator" />
          <ToolbarItem>
            <LatencyFilter performanceArtifacts={performanceArtifacts ?? []} />
            <Popover
              bodyContent={
                <>
                  Filter models performance benchmarks by measured latency.
                  <ul style={{ listStyleType: 'disc', paddingLeft: '20px', marginTop: '8px' }}>
                    <li>
                      <strong>Metric:</strong> Select the latency metric (TTFT, E2E, or ITL) to
                      evaluate.
                    </li>
                    <li>
                      <strong>Percentile:</strong> Choose how strictly the model must meet the
                      target. For example, P90 means 90% of requests must meet the selected
                      threshold.
                    </li>
                    <li>
                      <strong>Threshold:</strong> Set the maximum latency in milliseconds. Models
                      exceeding this value are excluded.
                    </li>
                  </ul>
                </>
              }
              appendTo={() => document.body}
            >
              <Button
                variant="plain"
                aria-label="More info for latency"
                className="pf-v6-u-p-xs"
                icon={<HelpIcon />}
              />
            </Popover>
          </ToolbarItem>
          <ToolbarItem>
            <MaxRpsFilter performanceArtifacts={performanceArtifacts ?? []} />
          </ToolbarItem>
          <ToolbarItem variant="separator" />
          <ToolbarItem>
            <HardwareTypeFilter performanceArtifacts={performanceArtifacts ?? []} />
          </ToolbarItem>
        </ToolbarGroup>
        <ModelCatalogActiveFilters filtersToShow={filtersToShow} />
      </ToolbarContent>
    </Toolbar>
  );
};

export default HardwareConfigurationFilterToolbar;

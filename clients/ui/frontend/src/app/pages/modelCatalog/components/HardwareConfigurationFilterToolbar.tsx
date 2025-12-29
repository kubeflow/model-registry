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
import { clearAllFilters } from '~/app/pages/modelCatalog/utils/hardwareConfigurationFilterUtils';
import { hasFiltersApplied } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import {
  ModelCatalogStringFilterKey,
  ModelCatalogNumberFilterKey,
  ALL_LATENCY_FIELD_NAMES,
} from '~/concepts/modelCatalog/const';
import WorkloadTypeFilter from './globalFilters/WorkloadTypeFilter';
import HardwareTypeFilter from './globalFilters/HardwareTypeFilter';
import MinRpsFilter from './globalFilters/MinRpsFilter';
import LatencyFilter from './globalFilters/LatencyFilter';
import ModelCatalogActiveFilters from './ModelCatalogActiveFilters';

type HardwareConfigurationFilterToolbarProps = {
  performanceArtifacts: CatalogPerformanceMetricsArtifact[];
};

/**
 * Filter keys that are shown on the performance/hardware configuration page.
 * This is used to determine which filters to show in the active filters chips
 * and which filters to clear when "Reset all filters" is clicked.
 */
const PERFORMANCE_FILTER_KEYS: ModelCatalogFilterKey[] = [
  ModelCatalogStringFilterKey.USE_CASE,
  ModelCatalogStringFilterKey.HARDWARE_TYPE,
  ModelCatalogNumberFilterKey.MIN_RPS,
];

/**
 * Gets the active filter keys including any active latency filters from filterData
 */
const getActivePerformanceFilterKeys = (
  filterData: ModelCatalogFilterStates,
): ModelCatalogFilterKey[] => {
  const activeLatencyKeys = ALL_LATENCY_FIELD_NAMES.filter((key) => filterData[key] !== undefined);
  return [...PERFORMANCE_FILTER_KEYS, ...activeLatencyKeys];
};

const HardwareConfigurationFilterToolbar: React.FC<HardwareConfigurationFilterToolbarProps> = ({
  performanceArtifacts,
}) => {
  const { filterOptions, filterOptionsLoaded, filterOptionsLoadError, filterData, setFilterData } =
    React.useContext(ModelCatalogContext);

  // Get all performance filter keys including active latency filters
  const filtersToShow = React.useMemo(
    () => getActivePerformanceFilterKeys(filterData),
    [filterData],
  );

  // Check if any performance filters are active (only checking performance-related filters)
  const hasActiveFilters = React.useMemo(
    () => hasFiltersApplied(filterData, filtersToShow),
    [filterData, filtersToShow],
  );

  if (!filterOptionsLoaded || filterOptionsLoadError || !filterOptions) {
    return null;
  }

  const handleClearAllFilters = () => {
    // Only clear performance-related filters, not the basic catalog filters
    clearAllFilters(setFilterData, filtersToShow);
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
            <LatencyFilter performanceArtifacts={performanceArtifacts} />
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
            <MinRpsFilter performanceArtifacts={performanceArtifacts} />
          </ToolbarItem>
          <ToolbarItem variant="separator" />
          <ToolbarItem>
            <HardwareTypeFilter performanceArtifacts={performanceArtifacts} />
          </ToolbarItem>
        </ToolbarGroup>
        <ModelCatalogActiveFilters filtersToShow={filtersToShow} />
      </ToolbarContent>
    </Toolbar>
  );
};

export default HardwareConfigurationFilterToolbar;

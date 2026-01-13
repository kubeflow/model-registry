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
  getPerformanceFiltersToShow,
  getAllFiltersToShow,
  BASIC_FILTER_KEYS,
  isPerformanceFilterKey,
} from '~/concepts/modelCatalog/const';
import { isValueDifferentFromDefault } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import WorkloadTypeFilter from './globalFilters/WorkloadTypeFilter';
import HardwareTypeFilter from './globalFilters/HardwareTypeFilter';
import MaxRpsFilter from './globalFilters/MaxRpsFilter';
import LatencyFilter from './globalFilters/LatencyFilter';
import ModelCatalogActiveFilters from './ModelCatalogActiveFilters';

type HardwareConfigurationFilterToolbarProps = {
  onResetAllFilters?: () => void;
  /** If true, shows basic filter chips. Defaults to false (only show on landing page when toggle is ON). */
  includeBasicFilters?: boolean;
  /** If true, shows performance filter chips. Landing page passes performanceViewEnabled, details page passes true. */
  includePerformanceFilters?: boolean;
};

const HardwareConfigurationFilterToolbar: React.FC<HardwareConfigurationFilterToolbarProps> = ({
  onResetAllFilters,
  includeBasicFilters = false,
  includePerformanceFilters = true,
}) => {
  const {
    filterOptions,
    filterOptionsLoaded,
    filterOptionsLoadError,
    filterData,
    getPerformanceFilterDefaultValue,
  } = React.useContext(ModelCatalogContext);

  // Get filter keys to show in the chip bar based on props
  // - includeBasicFilters: show basic filters (Task, Provider, License, Language)
  // - includePerformanceFilters: show performance filters (Workload type, Hardware type, Max RPS, Latency)
  const filtersToShow = React.useMemo(() => {
    if (includeBasicFilters && includePerformanceFilters) {
      return getAllFiltersToShow(filterData);
    }
    if (includePerformanceFilters) {
      return getPerformanceFiltersToShow(filterData);
    }
    if (includeBasicFilters) {
      return BASIC_FILTER_KEYS;
    }
    return [];
  }, [filterData, includeBasicFilters, includePerformanceFilters]);

  // Check if there are any visible filter chips (to control "Clear all filters" button visibility)
  // A chip is visible if:
  // - For basic filters: has a non-empty value
  // - For performance filters: has a value different from the default (regardless of toggle state)
  const hasVisibleChips = React.useMemo(
    () =>
      filtersToShow.some((filterKey) => {
        const filterValue = filterData[filterKey];

        // Skip if no value is set
        if (!filterValue) {
          return false;
        }

        // For array values (string filters), skip if empty
        if (Array.isArray(filterValue) && filterValue.length === 0) {
          return false;
        }

        // For performance filters, always check if value differs from default
        // (performance filters always have values - defaults when not explicitly set)
        if (isPerformanceFilterKey(filterKey)) {
          const defaultValue = getPerformanceFilterDefaultValue(filterKey);
          return isValueDifferentFromDefault(filterValue, defaultValue);
        }

        // For basic filters, any non-empty value means visible
        return true;
      }),
    [filtersToShow, filterData, getPerformanceFilterDefaultValue],
  );

  if (!filterOptionsLoaded || filterOptionsLoadError || !filterOptions) {
    return null;
  }

  return (
    <Toolbar
      // Only show "Clear all filters" button when there are visible chips to clear
      {...(onResetAllFilters && hasVisibleChips ? { clearAllFilters: onResetAllFilters } : {})}
    >
      <ToolbarContent rowWrap={{ default: 'wrap' }}>
        <ToolbarGroup rowWrap={{ default: 'wrap' }}>
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
            <LatencyFilter />
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
            <MaxRpsFilter />
          </ToolbarItem>
          <ToolbarItem variant="separator" />
          <ToolbarItem>
            <HardwareTypeFilter />
          </ToolbarItem>
        </ToolbarGroup>
        <ModelCatalogActiveFilters filtersToShow={filtersToShow} />
      </ToolbarContent>
    </Toolbar>
  );
};

export default HardwareConfigurationFilterToolbar;

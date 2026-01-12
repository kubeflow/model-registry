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
import { CatalogPerformanceMetricsArtifact } from '~/app/modelCatalogTypes';
import { getAllFiltersToShow, isPerformanceFilterKey } from '~/concepts/modelCatalog/const';
import { hasVisibleFilterChips } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import WorkloadTypeFilter from './globalFilters/WorkloadTypeFilter';
import HardwareTypeFilter from './globalFilters/HardwareTypeFilter';
import MaxRpsFilter from './globalFilters/MaxRpsFilter';
import LatencyFilter from './globalFilters/LatencyFilter';
import ModelCatalogActiveFilters from './ModelCatalogActiveFilters';

type HardwareConfigurationFilterToolbarProps = {
  performanceArtifacts?: CatalogPerformanceMetricsArtifact[];
  onResetAllFilters?: () => void;
};

const HardwareConfigurationFilterToolbar: React.FC<HardwareConfigurationFilterToolbarProps> = ({
  performanceArtifacts,
  onResetAllFilters,
}) => {
  const {
    filterOptions,
    filterOptionsLoaded,
    filterOptionsLoadError,
    filterData,
    performanceViewEnabled,
    getPerformanceFilterDefaultValue,
  } = React.useContext(ModelCatalogContext);

  // Get all filter keys (basic + performance) to show in the chip bar
  const allFiltersToShow = React.useMemo(() => getAllFiltersToShow(filterData), [filterData]);

  // Check if there are visible filter chips (accounting for defaults)
  const hasVisibleChips = React.useMemo(
    () =>
      hasVisibleFilterChips(
        filterData,
        allFiltersToShow,
        getPerformanceFilterDefaultValue,
        performanceViewEnabled,
        isPerformanceFilterKey,
      ),
    [filterData, allFiltersToShow, getPerformanceFilterDefaultValue, performanceViewEnabled],
  );

  if (!filterOptionsLoaded || filterOptionsLoadError || !filterOptions) {
    return null;
  }

  // Custom clear filters button with test ID for Cypress tests
  const customLabelGroupContent = onResetAllFilters && hasVisibleChips && (
    <ToolbarItem>
      <Button variant="link" onClick={onResetAllFilters} data-testid="clear-all-filters-button">
        Clear all filters
      </Button>
    </ToolbarItem>
  );

  return (
    <Toolbar
      {...(onResetAllFilters
        ? {
            clearAllFilters: onResetAllFilters,
            customLabelGroupContent,
          }
        : {})}
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
        {hasVisibleChips && <ModelCatalogActiveFilters filtersToShow={allFiltersToShow} />}
      </ToolbarContent>
    </Toolbar>
  );
};

export default HardwareConfigurationFilterToolbar;

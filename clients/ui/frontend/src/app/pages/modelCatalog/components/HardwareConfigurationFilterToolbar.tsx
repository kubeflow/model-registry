import * as React from 'react';
import { Toolbar, ToolbarContent, ToolbarGroup, ToolbarItem } from '@patternfly/react-core';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { CatalogPerformanceMetricsArtifact } from '~/app/modelCatalogTypes';
import { clearAllFilters } from '~/app/pages/modelCatalog/utils/hardwareConfigurationFilterUtils';
import { hasFiltersApplied } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import WorkloadTypeFilter from './globalFilters/WorkloadTypeFilter';
import HardwareTypeFilter from './globalFilters/HardwareTypeFilter';
import MinRpsFilter from './globalFilters/MinRpsFilter';
import MaxLatencyFilter from './globalFilters/MaxLatencyFilter';
import HardwareConfigurationActiveFilters from './HardwareConfigurationActiveFilters';

type HardwareConfigurationFilterToolbarProps = {
  performanceArtifacts: CatalogPerformanceMetricsArtifact[];
};

const HardwareConfigurationFilterToolbar: React.FC<HardwareConfigurationFilterToolbarProps> = ({
  performanceArtifacts,
}) => {
  const { filterOptions, filterOptionsLoaded, filterOptionsLoadError, filterData, setFilterData } =
    React.useContext(ModelCatalogContext);

  const hasActiveFilters = React.useMemo(() => hasFiltersApplied(filterData), [filterData]);

  if (!filterOptionsLoaded || filterOptionsLoadError || !filterOptions) {
    return null;
  }

  const handleClearAllFilters = () => {
    clearAllFilters(setFilterData);
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
          </ToolbarItem>
          <ToolbarItem variant="separator" />
          <ToolbarItem>
            <MaxLatencyFilter performanceArtifacts={performanceArtifacts} />
          </ToolbarItem>
          <ToolbarItem>
            <MinRpsFilter performanceArtifacts={performanceArtifacts} />
          </ToolbarItem>
          <ToolbarItem variant="separator" />
          <ToolbarItem>
            <HardwareTypeFilter performanceArtifacts={performanceArtifacts} />
          </ToolbarItem>
        </ToolbarGroup>
        <HardwareConfigurationActiveFilters />
      </ToolbarContent>
    </Toolbar>
  );
};

export default HardwareConfigurationFilterToolbar;

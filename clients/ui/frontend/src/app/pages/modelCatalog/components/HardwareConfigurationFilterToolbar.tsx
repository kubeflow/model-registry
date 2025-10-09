import * as React from 'react';
import { Toolbar, ToolbarContent, ToolbarGroup, ToolbarItem } from '@patternfly/react-core';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { CatalogPerformanceMetricsArtifact } from '~/app/modelCatalogTypes';
import WorkloadTypeFilter from './globalFilters/WorkloadTypeFilter';
import HardwareTypeFilter from './globalFilters/HardwareTypeFilter';
import MinRpsFilter from './globalFilters/MinRpsFilter';
import MaxLatencyFilter from './globalFilters/MaxLatencyFilter';
import { clearAllFilters } from '../utils/hardwareConfigurationFilterUtils';

type HardwareConfigurationFilterToolbarProps = {
  performanceArtifacts: CatalogPerformanceMetricsArtifact[];
  appliedHardwareTypes: string[];
  onApplyHardwareFilters: (types: string[]) => void;
  onResetHardwareFilters: () => void;
};

const HardwareConfigurationFilterToolbar: React.FC<HardwareConfigurationFilterToolbarProps> = ({
  performanceArtifacts,
  appliedHardwareTypes,
  onApplyHardwareFilters,
  onResetHardwareFilters,
}) => {
  const { filterOptions, filterOptionsLoaded, filterOptionsLoadError, setFilterData } =
    React.useContext(ModelCatalogContext);

  if (!filterOptionsLoaded || filterOptionsLoadError || !filterOptions) {
    return null;
  }

  const { filters } = filterOptions;

  const handleClearAllFilters = () => {
    clearAllFilters(setFilterData);
  };

  return (
    <Toolbar clearAllFilters={handleClearAllFilters} clearFiltersButtonText="Reset all filters">
      <ToolbarContent>
        <ToolbarGroup>
          <ToolbarItem>
            <WorkloadTypeFilter filterOptions={filters} />
          </ToolbarItem>
          <ToolbarItem
            style={{ borderLeft: '1px solid #d2d2d2', paddingLeft: '16px', marginLeft: '16px' }}
          >
            <MaxLatencyFilter filterOptions={filterOptions} />
          </ToolbarItem>
          <ToolbarItem>
            <MinRpsFilter filterOptions={filterOptions} />
          </ToolbarItem>
          <ToolbarItem
            style={{ borderLeft: '1px solid #d2d2d2', paddingLeft: '16px', marginLeft: '16px' }}
          >
            <HardwareTypeFilter
              performanceArtifacts={performanceArtifacts}
              appliedHardwareTypes={appliedHardwareTypes}
              onApplyHardwareFilters={onApplyHardwareFilters}
              onResetHardwareFilters={onResetHardwareFilters}
            />
          </ToolbarItem>
        </ToolbarGroup>
      </ToolbarContent>
    </Toolbar>
  );
};

export default HardwareConfigurationFilterToolbar;

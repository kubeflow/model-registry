import * as React from 'react';
import {
  Badge,
  Checkbox,
  Dropdown,
  Flex,
  FlexItem,
  MenuToggle,
  MenuToggleElement,
  Panel,
  PanelMain,
} from '@patternfly/react-core';
import { CatalogPerformanceMetricsArtifact } from '~/app/modelCatalogTypes';
import { getUniqueHardwareTypes } from '~/app/pages/modelCatalog/utils/hardwareConfigurationFilterUtils';
import { useHardwareTypeFilterState } from '~/app/pages/modelCatalog/utils/hardwareTypeFilterState';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { ModelCatalogStringFilterKey } from '~/concepts/modelCatalog/const';

type HardwareTypeFilterProps = {
  performanceArtifacts: CatalogPerformanceMetricsArtifact[];
};

type HardwareTypeOption = {
  value: string;
  label: string;
};

const HardwareTypeFilter: React.FC<HardwareTypeFilterProps> = ({ performanceArtifacts }) => {
  const { appliedHardwareTypes, setAppliedHardwareTypes } = useHardwareTypeFilterState();
  const { filterOptions } = React.useContext(ModelCatalogContext);
  const [isOpen, setIsOpen] = React.useState(false);

  // Get unique hardware types from performance artifacts, or fall back to filterOptions
  const hardwareOptions: HardwareTypeOption[] = React.useMemo(() => {
    // First try to get from performance artifacts
    if (performanceArtifacts.length > 0) {
      const uniqueTypes = getUniqueHardwareTypes(performanceArtifacts);
      return uniqueTypes.map((type) => ({
        value: type,
        label: type,
      }));
    }
    // Fall back to filterOptions from context
    const filterValue = filterOptions?.filters?.[ModelCatalogStringFilterKey.HARDWARE_TYPE];
    if (filterValue && 'values' in filterValue && Array.isArray(filterValue.values)) {
      return filterValue.values.map((type: string) => ({
        value: type,
        label: type,
      }));
    }
    return [];
  }, [performanceArtifacts, filterOptions]);

  const selectedCount = appliedHardwareTypes.length;

  const isHardwareSelected = (value: string): boolean => appliedHardwareTypes.includes(value);

  const toggleHardwareSelection = (value: string, selected: boolean) => {
    if (selected) {
      setAppliedHardwareTypes([...appliedHardwareTypes, value]);
    } else {
      setAppliedHardwareTypes(appliedHardwareTypes.filter((item) => item !== value));
    }
  };

  const toggle = (toggleRef: React.Ref<MenuToggleElement>) => (
    <MenuToggle
      ref={toggleRef}
      onClick={() => setIsOpen(!isOpen)}
      isExpanded={isOpen}
      isFullHeight
      style={{ minWidth: '200px', width: 'fit-content', height: '56px' }}
      badge={selectedCount > 0 ? <Badge>{selectedCount} selected</Badge> : undefined}
    >
      Hardware type
    </MenuToggle>
  );

  const filterContent = (
    <Panel>
      <PanelMain className="pf-v6-u-p-md">
        <Flex direction={{ default: 'column' }} spaceItems={{ default: 'spaceItemsSm' }}>
          {/* Hardware type checkboxes */}
          <FlexItem>
            <Flex direction={{ default: 'column' }} spaceItems={{ default: 'spaceItemsXs' }}>
              {hardwareOptions.map((option) => (
                <FlexItem key={option.value}>
                  <Flex alignItems={{ default: 'alignItemsCenter' }}>
                    <FlexItem flex={{ default: 'flex_1' }}>
                      <Checkbox
                        label={option.label}
                        id={option.value}
                        isChecked={isHardwareSelected(option.value)}
                        onChange={(_, checked) => toggleHardwareSelection(option.value, checked)}
                      />
                    </FlexItem>
                  </Flex>
                </FlexItem>
              ))}
            </Flex>
          </FlexItem>
        </Flex>
      </PanelMain>
    </Panel>
  );

  return (
    <Dropdown
      isOpen={isOpen}
      onOpenChange={setIsOpen}
      toggle={toggle}
      shouldFocusToggleOnSelect={false}
    >
      {filterContent}
    </Dropdown>
  );
};

export default HardwareTypeFilter;

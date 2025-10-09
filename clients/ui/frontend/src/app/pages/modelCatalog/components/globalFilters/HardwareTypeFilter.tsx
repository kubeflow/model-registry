import * as React from 'react';
import {
  Badge,
  Button,
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

type HardwareTypeFilterProps = {
  performanceArtifacts: CatalogPerformanceMetricsArtifact[];
  appliedHardwareTypes: string[];
  onApplyHardwareFilters: (types: string[]) => void;
  onResetHardwareFilters: () => void;
};

type HardwareTypeOption = {
  value: string;
  label: string;
};

const HardwareTypeFilter: React.FC<HardwareTypeFilterProps> = ({
  performanceArtifacts,
  appliedHardwareTypes,
  onApplyHardwareFilters,
  onResetHardwareFilters,
}) => {
  const [isOpen, setIsOpen] = React.useState(false);
  const [localSelectedHardwareTypes, setLocalSelectedHardwareTypes] = React.useState<string[]>([]);

  // Get unique hardware types from actual performance artifacts
  const hardwareOptions: HardwareTypeOption[] = React.useMemo(() => {
    const uniqueTypes = getUniqueHardwareTypes(performanceArtifacts);
    return uniqueTypes.map((type) => ({
      value: type,
      label: type,
    }));
  }, [performanceArtifacts]);

  const selectedCount = localSelectedHardwareTypes.length;

  // Initialize local state from applied state when dropdown opens
  React.useEffect(() => {
    if (isOpen) {
      setLocalSelectedHardwareTypes([...appliedHardwareTypes]);
    }
  }, [isOpen, appliedHardwareTypes]);

  const isHardwareSelected = (value: string): boolean => localSelectedHardwareTypes.includes(value);

  const toggleHardwareSelection = (value: string, selected: boolean) => {
    if (selected) {
      setLocalSelectedHardwareTypes([...localSelectedHardwareTypes, value]);
    } else {
      setLocalSelectedHardwareTypes(localSelectedHardwareTypes.filter((item) => item !== value));
    }
  };

  const handleApplyFilter = () => {
    onApplyHardwareFilters(localSelectedHardwareTypes);
    setIsOpen(false);
  };

  const handleReset = () => {
    setLocalSelectedHardwareTypes([]);
    onResetHardwareFilters();
    setIsOpen(false);
  };

  const toggle = (toggleRef: React.Ref<MenuToggleElement>) => (
    <MenuToggle
      ref={toggleRef}
      onClick={() => setIsOpen(!isOpen)}
      isExpanded={isOpen}
      style={{ minWidth: '200px', width: 'fit-content' }}
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

          {/* Buttons: Apply filter first, then Reset */}
          <FlexItem>
            <Flex spaceItems={{ default: 'spaceItemsSm' }}>
              <FlexItem>
                <Button variant="primary" onClick={handleApplyFilter}>
                  Apply filter
                </Button>
              </FlexItem>
              <FlexItem>
                <Button variant="secondary" onClick={handleReset}>
                  Reset
                </Button>
              </FlexItem>
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
      <div
        onMouseDown={(e) => e.stopPropagation()}
        onClick={(e) => e.stopPropagation()}
        onKeyDown={(e) => e.stopPropagation()}
        role="presentation"
      >
        {filterContent}
      </div>
    </Dropdown>
  );
};

export default HardwareTypeFilter;

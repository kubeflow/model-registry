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
import { ModelCatalogStringFilterKey, UseCaseOptionValue } from '~/concepts/modelCatalog/const';
import { USE_CASE_OPTIONS } from '~/app/pages/modelCatalog/utils/workloadTypeUtils';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';

const WorkloadTypeFilter: React.FC = () => {
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
  const selectedUseCases = filterData[ModelCatalogStringFilterKey.USE_CASE];
  const [isOpen, setIsOpen] = React.useState(false);

  const selectedCount = selectedUseCases.length;

  const isUseCaseSelected = (value: UseCaseOptionValue): boolean =>
    selectedUseCases.includes(value);

  const toggleUseCaseSelection = (value: UseCaseOptionValue, selected: boolean) => {
    if (selected) {
      setFilterData(ModelCatalogStringFilterKey.USE_CASE, [...selectedUseCases, value]);
    } else {
      setFilterData(
        ModelCatalogStringFilterKey.USE_CASE,
        selectedUseCases.filter((item) => item !== value),
      );
    }
  };

  const toggle = (toggleRef: React.Ref<MenuToggleElement>) => (
    <MenuToggle
      ref={toggleRef}
      data-testid="workload-type-filter"
      onClick={() => setIsOpen(!isOpen)}
      isExpanded={isOpen}
      style={{ minWidth: '200px', width: 'fit-content' }}
      badge={selectedCount > 0 ? <Badge>{selectedCount} selected</Badge> : undefined}
    >
      Workload type
    </MenuToggle>
  );

  const filterContent = (
    <Panel>
      <PanelMain className="pf-v6-u-p-md">
        <Flex direction={{ default: 'column' }} spaceItems={{ default: 'spaceItemsSm' }}>
          {/* Workload type checkboxes */}
          <FlexItem>
            <Flex direction={{ default: 'column' }} spaceItems={{ default: 'spaceItemsXs' }}>
              {USE_CASE_OPTIONS.map((option) => (
                <FlexItem key={option.value}>
                  <Flex alignItems={{ default: 'alignItemsCenter' }}>
                    <FlexItem flex={{ default: 'flex_1' }}>
                      <Checkbox
                        label={`${option.label} (${option.inputTokens} input | ${option.outputTokens} output tokens)`}
                        id={option.value}
                        data-testid={`workload-type-filter-${option.value}`}
                        isChecked={isUseCaseSelected(option.value)}
                        onChange={(_, checked) => toggleUseCaseSelection(option.value, checked)}
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

export default WorkloadTypeFilter;

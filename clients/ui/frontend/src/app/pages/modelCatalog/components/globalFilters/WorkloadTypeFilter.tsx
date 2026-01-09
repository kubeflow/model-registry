import * as React from 'react';
import {
  Dropdown,
  DropdownItem,
  DropdownList,
  MenuToggle,
  MenuToggleElement,
} from '@patternfly/react-core';
import { ModelCatalogStringFilterKey, UseCaseOptionValue } from '~/concepts/modelCatalog/const';
import {
  USE_CASE_OPTIONS,
  isUseCaseOptionValue,
  getUseCaseDisplayLabel,
} from '~/app/pages/modelCatalog/utils/workloadTypeUtils';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';

const WorkloadTypeFilter: React.FC = () => {
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
  const selectedUseCases = filterData[ModelCatalogStringFilterKey.USE_CASE];
  const [isOpen, setIsOpen] = React.useState(false);

  // Use static USE_CASE_OPTIONS - these are the only valid workload types we support
  const availableWorkloadTypes: UseCaseOptionValue[] = React.useMemo(
    () => USE_CASE_OPTIONS.map((opt) => opt.value),
    [],
  );

  // The currently selected workload type (single-select, so take first element)
  const selectedValue = selectedUseCases.length > 0 ? selectedUseCases[0] : undefined;

  const handleSelect = (value: string | undefined) => {
    if (value && isUseCaseOptionValue(value)) {
      setFilterData(ModelCatalogStringFilterKey.USE_CASE, [value]);
    } else {
      setFilterData(ModelCatalogStringFilterKey.USE_CASE, []);
    }
    setIsOpen(false);
  };

  const toggle = (toggleRef: React.Ref<MenuToggleElement>) => (
    <MenuToggle
      ref={toggleRef}
      data-testid="workload-type-filter"
      onClick={() => setIsOpen(!isOpen)}
      isExpanded={isOpen}
      isFullHeight
      style={{ minWidth: '200px', width: 'fit-content', height: '56px' }}
    >
      {selectedValue ? (
        <>
          <strong>Workload type:</strong> {getUseCaseDisplayLabel(selectedValue)}
        </>
      ) : (
        'Workload type'
      )}
    </MenuToggle>
  );

  return (
    <Dropdown isOpen={isOpen} onOpenChange={setIsOpen} toggle={toggle} shouldFocusToggleOnSelect>
      <DropdownList>
        {availableWorkloadTypes.map((value) => (
          <DropdownItem
            key={value}
            onClick={() => handleSelect(value)}
            isSelected={selectedValue === value}
            data-testid={`workload-type-filter-${value}`}
          >
            {getUseCaseDisplayLabel(value)}
          </DropdownItem>
        ))}
      </DropdownList>
    </Dropdown>
  );
};

export default WorkloadTypeFilter;

import * as React from 'react';
import {
  Dropdown,
  DropdownItem,
  DropdownList,
  MenuToggle,
  MenuToggleElement,
} from '@patternfly/react-core';
import { asEnumMember } from 'mod-arch-core';
import { ModelCatalogStringFilterKey, UseCaseOptionValue } from '~/concepts/modelCatalog/const';
import { USE_CASE_OPTIONS } from '~/app/pages/modelCatalog/utils/workloadTypeUtils';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';

const UseCaseFilter: React.FC = () => {
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
  const selectedUseCase = filterData[ModelCatalogStringFilterKey.USE_CASE];
  const [isOpen, setIsOpen] = React.useState(false);

  const handleUseCaseChange = (useCase: string) => {
    const useCaseValue = asEnumMember(useCase, UseCaseOptionValue);
    if (useCaseValue) {
      const newValue = useCaseValue === selectedUseCase ? undefined : useCaseValue;
      setFilterData(ModelCatalogStringFilterKey.USE_CASE, newValue);
    }
    setIsOpen(false);
  };

  // Get the display text for the toggle
  const getToggleText = () => {
    if (selectedUseCase) {
      const selectedOption = USE_CASE_OPTIONS.find((option) => option.value === selectedUseCase);
      return selectedOption ? (
        <>
          <strong>Workload type:</strong> {selectedOption.label}
        </>
      ) : (
        'Workload type'
      );
    }
    return 'Workload type';
  };

  const toggle = (toggleRef: React.Ref<MenuToggleElement>) => (
    <MenuToggle
      ref={toggleRef}
      onClick={() => setIsOpen(!isOpen)}
      isExpanded={isOpen}
      style={{ minWidth: '200px', width: 'fit-content' }}
    >
      {getToggleText()}
    </MenuToggle>
  );

  return (
    <Dropdown
      isOpen={isOpen}
      onOpenChange={setIsOpen}
      toggle={toggle}
      shouldFocusToggleOnSelect={false}
    >
      <DropdownList>
        {USE_CASE_OPTIONS.map((option) => (
          <DropdownItem
            key={option.value}
            onClick={() => handleUseCaseChange(option.value)}
            isSelected={selectedUseCase === option.value}
          >
            {option.label}
          </DropdownItem>
        ))}
      </DropdownList>
    </Dropdown>
  );
};

export default UseCaseFilter;

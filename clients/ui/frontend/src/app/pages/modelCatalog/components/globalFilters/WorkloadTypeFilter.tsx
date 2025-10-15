import * as React from 'react';
import {
  Button,
  MenuToggle,
  MenuToggleElement,
  Popover,
  Select,
  SelectList,
  SelectOption,
} from '@patternfly/react-core';
import { HelpIcon } from '@patternfly/react-icons';
import {
  ModelCatalogNumberFilterKey,
  WorkloadTypeOptionValue,
} from '~/concepts/modelCatalog/const';
import { useCatalogNumberFilterState } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import {
  WORKLOAD_TYPE_OPTIONS,
  maxInputOutputTokensToWorkloadType,
  workloadTypeToMaxInputOutputTokens,
} from '~/app/pages/modelCatalog/utils/workloadTypeUtils';

const WorkloadTypeFilter: React.FC = () => {
  const { value: maxInputTokens, setValue: setMaxInputTokens } = useCatalogNumberFilterState(
    ModelCatalogNumberFilterKey.MAX_INPUT_TOKENS,
  );
  const { value: maxOutputTokens, setValue: setMaxOutputTokens } = useCatalogNumberFilterState(
    ModelCatalogNumberFilterKey.MAX_OUTPUT_TOKENS,
  );
  const [isOpen, setIsOpen] = React.useState(false);

  // Derive the selected workload type from the token values
  const selectedWorkloadType = maxInputOutputTokensToWorkloadType(maxInputTokens, maxOutputTokens);
  const selectedOption = WORKLOAD_TYPE_OPTIONS.find(
    (option) => option.value === selectedWorkloadType,
  );
  const displayText = selectedOption ? selectedOption.label : 'Select workload type';

  const toggle = (toggleRef: React.Ref<MenuToggleElement>) => (
    <MenuToggle
      ref={toggleRef}
      onClick={() => setIsOpen(!isOpen)}
      isExpanded={isOpen}
      style={{ minWidth: '300px', width: 'fit-content' }}
    >
      <span className="pf-v6-u-mr-sm">Workload type:</span>
      {displayText}
    </MenuToggle>
  );

  return (
    <>
      <Select
        isOpen={isOpen}
        selected={selectedWorkloadType}
        onSelect={(_, selectedVal) => {
          // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
          const workloadType = selectedVal as WorkloadTypeOptionValue;
          const tokenValues = workloadTypeToMaxInputOutputTokens(workloadType);

          if (tokenValues) {
            setMaxInputTokens(tokenValues.maxInputTokens);
            setMaxOutputTokens(tokenValues.maxOutputTokens);
          }

          setIsOpen(false);
        }}
        onOpenChange={setIsOpen}
        toggle={toggle}
        shouldFocusToggleOnSelect
      >
        <SelectList>
          {WORKLOAD_TYPE_OPTIONS.map((option) => (
            <SelectOption
              key={option.value}
              value={option.value}
              isSelected={selectedWorkloadType === option.value}
            >
              {option.label}
            </SelectOption>
          ))}
        </SelectList>
      </Select>
      <Popover
        aria-label="Workload type information"
        bodyContent="Select a workload type to view performance under specific input and output token lengths."
        appendTo={() => document.body}
      >
        <Button
          variant="plain"
          aria-label="More info for workload type"
          onClick={(e) => e.stopPropagation()}
          icon={<HelpIcon />}
        />
      </Popover>
    </>
  );
};

export default WorkloadTypeFilter;

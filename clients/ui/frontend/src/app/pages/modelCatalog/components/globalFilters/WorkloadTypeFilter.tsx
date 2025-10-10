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
import { ModelCatalogNumberFilterKey } from '~/concepts/modelCatalog/const';
import { useCatalogNumberFilterState } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { CatalogFilterOptions } from '~/app/modelCatalogTypes';

const filterKey = ModelCatalogNumberFilterKey.WORKLOAD_TYPE;

type WorkloadTypeFilterProps = {
  filterOptions?: CatalogFilterOptions | null;
};

type WorkloadTypeOption = {
  value: string;
  label: string;
  description: string;
};

const WORKLOAD_TYPE_OPTIONS: WorkloadTypeOption[] = [
  {
    value: 'chat',
    label: 'Chat (512 input | 256 output tokens)',
    description: 'Conversational AI workload with moderate input/output token lengths',
  },
  {
    value: 'rag',
    label: 'RAG (4096 input | 512 output tokens)',
    description: 'Retrieval-Augmented Generation with larger context windows',
  },
  {
    value: 'summarization',
    label: 'Summarization (2048 input | 256 output tokens)',
    description: 'Text summarization tasks with long input documents',
  },
  {
    value: 'code_generation',
    label: 'Code Generation (1024 input | 512 output tokens)',
    description: 'Code generation and completion tasks',
  },
];

const WorkloadTypeFilter: React.FC<WorkloadTypeFilterProps> = () => {
  const { value: appliedValue, setValue } = useCatalogNumberFilterState(filterKey);
  const [isOpen, setIsOpen] = React.useState(false);
  const [localSelectedValue, setLocalSelectedValue] = React.useState<string | undefined>(undefined);

  const workloadOptions = WORKLOAD_TYPE_OPTIONS;

  // Initialize local state from applied state when dropdown opens
  React.useEffect(() => {
    if (isOpen) {
      const applied = workloadOptions.find((opt) => opt.value === appliedValue?.toString());
      setLocalSelectedValue(applied?.value);
    }
  }, [isOpen, appliedValue, workloadOptions]);

  const appliedOption = workloadOptions.find((option) => option.value === appliedValue?.toString());
  const displayText = appliedOption ? appliedOption.label : 'Select workload type';

  const toggle = (toggleRef: React.Ref<MenuToggleElement>) => (
    <MenuToggle
      ref={toggleRef}
      onClick={() => setIsOpen(!isOpen)}
      isExpanded={isOpen}
      style={{ minWidth: '300px', width: 'fit-content' }}
    >
      <span style={{ marginRight: '8px' }}>Workload type:</span>
      {displayText}
    </MenuToggle>
  );

  return (
    <>
      <Select
        isOpen={isOpen}
        selected={localSelectedValue}
        onSelect={(_, selectedVal) => {
          // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
          setLocalSelectedValue(selectedVal as string);
          // Apply immediately when selecting from dropdown
          const selectedOption = workloadOptions.find((opt) => opt.value === selectedVal);
          setValue(
            selectedOption
              ? // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
                (selectedVal as unknown as number)
              : undefined,
          );
          setIsOpen(false);
        }}
        onOpenChange={setIsOpen}
        toggle={toggle}
        shouldFocusToggleOnSelect
      >
        <SelectList>
          {workloadOptions.map((option) => (
            <SelectOption
              key={option.value}
              value={option.value}
              isSelected={localSelectedValue === option.value}
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
        >
          <HelpIcon />
        </Button>
      </Popover>
    </>
  );
};

export default WorkloadTypeFilter;

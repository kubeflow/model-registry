import * as React from 'react';
import {
  Truncate,
  MenuToggle,
  Select,
  SelectList,
  SelectOption,
  SelectGroup,
  Divider,
  MenuToggleProps,
} from '@patternfly/react-core';

import './SimpleSelect.scss';

export type SimpleSelectOption = {
  key: string;
  label: string;
  description?: React.ReactNode;
  dropdownLabel?: React.ReactNode;
  isPlaceholder?: boolean;
  isDisabled?: boolean;
};

export type SimpleGroupSelectOption = {
  key: string;
  label: string;
  options: SimpleSelectOption[];
};

type SimpleSelectProps = {
  options?: SimpleSelectOption[];
  groupedOptions?: SimpleGroupSelectOption[];
  value?: string;
  toggleLabel?: React.ReactNode;
  placeholder?: string;
  onChange: (key: string, isPlaceholder: boolean) => void;
  isFullWidth?: boolean;
  toggleProps?: MenuToggleProps;
  isDisabled?: boolean;
  icon?: React.ReactNode;
  dataTestId?: string;
} & Omit<
  React.ComponentProps<typeof Select>,
  'isOpen' | 'toggle' | 'dropdownItems' | 'onChange' | 'selected'
>;

const SimpleSelect: React.FC<SimpleSelectProps> = ({
  isDisabled,
  onChange,
  options,
  groupedOptions,
  placeholder = 'Select...',
  value,
  toggleLabel,
  isFullWidth,
  icon,
  dataTestId,
  toggleProps,
  ...props
}) => {
  const [open, setOpen] = React.useState(false);

  const findOptionForKey = (key: string) =>
    options?.find((option) => option.key === key) ||
    groupedOptions
      ?.reduce<SimpleSelectOption[]>((acc, group) => [...acc, ...group.options], [])
      .find((o) => o.key === key);

  const selectedOption = value ? findOptionForKey(value) : undefined;
  const selectedLabel = selectedOption?.label ?? placeholder;

  return (
    <Select
      {...props}
      isOpen={open}
      selected={value || toggleLabel}
      onSelect={(e, selectValue) => {
        onChange(
          String(selectValue),
          !!selectValue && (findOptionForKey(String(selectValue))?.isPlaceholder ?? false),
        );
        setOpen(false);
      }}
      onOpenChange={setOpen}
      toggle={(toggleRef) => (
        <MenuToggle
          ref={toggleRef}
          data-testid={dataTestId}
          aria-label="Options menu"
          onClick={() => setOpen(!open)}
          icon={icon}
          isExpanded={open}
          isDisabled={isDisabled}
          isFullWidth={isFullWidth}
          {...toggleProps}
        >
          {toggleLabel || <Truncate content={selectedLabel} className="truncate-no-min-width" />}
        </MenuToggle>
      )}
      shouldFocusToggleOnSelect
    >
      {groupedOptions?.map((group, index) => (
        <>
          {index > 0 ? <Divider /> : null}
          <SelectGroup key={group.key} label={group.label}>
            <SelectList>
              {group.options.map(
                ({ key, label, dropdownLabel, description, isDisabled: optionDisabled }) => (
                  <SelectOption
                    key={key}
                    value={key}
                    description={description}
                    isDisabled={optionDisabled}
                    data-testid={key}
                  >
                    {dropdownLabel || label}
                  </SelectOption>
                ),
              )}
            </SelectList>
          </SelectGroup>
        </>
      )) ?? null}
      {options?.length ? (
        <SelectList>
          {groupedOptions?.length ? <Divider /> : null}
          {options.map(({ key, label, dropdownLabel, description, isDisabled: optionDisabled }) => (
            <SelectOption
              key={key}
              value={key}
              description={description}
              isDisabled={optionDisabled}
              data-testid={key}
            >
              {dropdownLabel || label}
            </SelectOption>
          ))}
        </SelectList>
      ) : null}
    </Select>
  );
};

export default SimpleSelect;

import * as React from 'react';
import { RoleBindingSubject, TypeaheadSelect } from 'mod-arch-shared';
import { useNamespaces } from '~/app/hooks/useNamespaces';

type RoleBindingPermissionsNameInputProps = {
  subjectKind: RoleBindingSubject['kind'];
  value: string;
  onChange?: (selection: string) => void;
  onClear?: () => void;
  placeholderText?: string;
  typeAhead?: string[];
  isProjectSubject?: boolean;
};

export const RoleBindingPermissionsNameInput: React.FC<RoleBindingPermissionsNameInputProps> = ({
  subjectKind,
  value,
  onChange,
  onClear,
  placeholderText,
  typeAhead,
  isProjectSubject,
}) => {
  const [namespaces] = useNamespaces();

  const selectOptions = React.useMemo(() => {
    let options: Array<{ value: string; content: string }> = [];

    if (subjectKind === 'Group' && typeAhead) {
      options = typeAhead.map((name) => ({ value: name, content: name }));
    }

    if (isProjectSubject) {
      const namespaceOptions = namespaces.map((namespace) => ({
        value: namespace.name,
        content: namespace.displayName || namespace.name,
      }));
      options = [...options, ...namespaceOptions];
    }

    // If we've selected an option that doesn't exist via isCreatable, include it in the options so it remains selected
    if (value && !options.some((option) => option.value === value)) {
      options.push({ value, content: value });
    }

    return options;
  }, [subjectKind, typeAhead, isProjectSubject, namespaces, value]);

  return (
    <TypeaheadSelect
      dataTestId={`role-binding-name-select ${value}`}
      isScrollable
      selectOptions={selectOptions}
      selected={value}
      isCreatable
      onClearSelection={onClear}
      onSelect={(_ev, selectedValue) => {
        if (typeof selectedValue === 'string') {
          onChange?.(selectedValue);
        }
      }}
      placeholder={placeholderText || `Enter ${subjectKind.toLowerCase()} name`}
      createOptionMessage={(newValue) => `Select "${newValue}"`}
    />
  );
};

import React from 'react';
import { useNamespaceSelector, useModularArchContext } from 'mod-arch-core';
import { SimpleSelect } from 'mod-arch-shared';
import { SimpleSelectOption } from 'mod-arch-shared/dist/components/SimpleSelect';

interface NamespaceSelectorProps {
  onSelect?: (namespace: string) => void;
  className?: string;
  isDisabled?: boolean;
  selectedNamespace?: string;
  placeholderText?: string;
  isFullWidth?: boolean;
  isGlobalSelector?: boolean;
  ignoreMandatoryNamespace?: boolean;
}

const NamespaceSelector: React.FC<NamespaceSelectorProps> = ({
  placeholderText,
  onSelect,
  className,
  isDisabled: externalDisabled,
  selectedNamespace,
  isFullWidth,
  isGlobalSelector,
  ignoreMandatoryNamespace,
}) => {
  const { namespaces = [], preferredNamespace, updatePreferredNamespace } = useNamespaceSelector();
  const { config } = useModularArchContext();

  // Check if mandatory namespace is configured
  const isMandatoryNamespace = Boolean(config.mandatoryNamespace);

  const baseDisabled = externalDisabled || namespaces.length === 0;
  const isDisabled = ignoreMandatoryNamespace ? baseDisabled : baseDisabled || isMandatoryNamespace;
  const options: SimpleSelectOption[] = namespaces.map((namespace) => ({
    key: namespace.name,
    label: namespace.name,
  }));

  const selectedValue = !isGlobalSelector
    ? selectedNamespace || ''
    : preferredNamespace?.name || namespaces[0]?.name || '';

  const handleChange = (key: string, isPlaceholder: boolean) => {
    if (isPlaceholder || !key) {
      return;
    }

    if (!isMandatoryNamespace && isGlobalSelector) {
      updatePreferredNamespace({ name: key });
    }

    if (onSelect) {
      onSelect(key);
    }
  };

  return (
    <SimpleSelect
      options={options}
      value={selectedValue}
      className={className}
      onChange={handleChange}
      placeholder={placeholderText}
      isDisabled={isDisabled}
      isFullWidth={isFullWidth}
      popperProps={{ maxWidth: '400px' }}
      dataTestId={isGlobalSelector ? 'navbar-namespace-selector' : 'form-namespace-selector'}
    />
  );
};

export default NamespaceSelector;

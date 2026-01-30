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
}

const NamespaceSelector: React.FC<NamespaceSelectorProps> = ({
  placeholderText,
  onSelect,
  className,
  isDisabled: externalDisabled,
  selectedNamespace,
  isFullWidth,
}) => {
  const { namespaces, preferredNamespace, updatePreferredNamespace } = useNamespaceSelector();
  const { config } = useModularArchContext();

  // Check if mandatory namespace is configured
  const isMandatoryNamespace = Boolean(config.mandatoryNamespace);

  const isDisabled = externalDisabled || isMandatoryNamespace || namespaces.length === 0;
  const options: SimpleSelectOption[] = namespaces.map((namespace) => ({
    key: namespace.name,
    label: namespace.name,
  }));

  const selectedValue = placeholderText
    ? selectedNamespace || ''
    : preferredNamespace?.name || namespaces[0]?.name || '';

  const handleChange = (key: string, isPlaceholder: boolean) => {
    if (isPlaceholder || !key) {
      return;
    }

    if (!isMandatoryNamespace) {
      if (!placeholderText) {
        updatePreferredNamespace({ name: key });
      }

      if (onSelect) {
        onSelect(key);
      }
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
      dataTestId={placeholderText ? 'form-namespace-selector' : 'navbar-namespace-selector'}
    />
  );
};

export default NamespaceSelector;

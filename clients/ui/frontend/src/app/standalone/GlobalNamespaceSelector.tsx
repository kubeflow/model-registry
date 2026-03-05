import React from 'react';
import { useNamespaceSelector, useModularArchContext } from 'mod-arch-core';
import { SimpleSelect } from 'mod-arch-shared';
import { SimpleSelectOption } from 'mod-arch-shared/dist/components/SimpleSelect';

interface GlobalNamespaceSelectorProps {
  onSelect?: (namespace: string) => void;
  className?: string;
  isDisabled?: boolean;
}

const GlobalNamespaceSelector: React.FC<GlobalNamespaceSelectorProps> = ({
  onSelect,
  className,
  isDisabled: externalDisabled,
}) => {
  const { namespaces = [], preferredNamespace, updatePreferredNamespace } = useNamespaceSelector();
  const { config } = useModularArchContext();

  const isMandatoryNamespace = Boolean(config.mandatoryNamespace);
  const isDisabled = externalDisabled || isMandatoryNamespace || namespaces.length === 0;

  const options: SimpleSelectOption[] = namespaces.map((namespace) => ({
    key: namespace.name,
    label: namespace.name,
  }));

  const selectedValue = preferredNamespace?.name || namespaces[0]?.name || '';

  const handleChange = (key: string, isPlaceholder: boolean) => {
    if (isPlaceholder || !key) {
      return;
    }

    if (!isMandatoryNamespace) {
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
      isDisabled={isDisabled}
      popperProps={{ maxWidth: '400px' }}
      dataTestId="navbar-namespace-selector"
    />
  );
};

export default GlobalNamespaceSelector;

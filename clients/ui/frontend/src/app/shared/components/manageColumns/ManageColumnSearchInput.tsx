// TODO this component was copied from odh-dashboard temporarily and should be abstracted out into mod-arch-shared.

import React from 'react';
import { SearchInput } from '@patternfly/react-core';

interface ManageColumnSearchInputProps {
  value: string;
  placeholder?: string;
  onSearch: (value: string) => void;
  dataTestId?: string;
}

export const ManageColumnSearchInput: React.FC<ManageColumnSearchInputProps> = ({
  value,
  placeholder = 'Filter by column name',
  onSearch,
  dataTestId = 'manage-column-search',
}) => {
  const handleChange = React.useCallback(
    (_: React.FormEvent<HTMLInputElement> | null, newValue: string) => {
      onSearch(newValue);
    },
    [onSearch],
  );

  return (
    <SearchInput
      data-testid={dataTestId}
      placeholder={placeholder}
      value={value}
      onChange={handleChange}
      onClear={() => handleChange(null, '')}
    />
  );
};

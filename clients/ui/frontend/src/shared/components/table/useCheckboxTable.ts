import * as React from 'react';
import { useCheckboxTableBase } from '~/shared/components/table';
import type Table from './Table';

type UseCheckboxTable = {
  selections: string[];
  tableProps: Required<Pick<React.ComponentProps<typeof Table>, 'selectAll'>>;
  toggleSelection: (id: string) => void;
  isSelected: (id: string) => boolean;
  setSelections: (selections: string[]) => void;
};

const useCheckboxTable = (
  dataIds: string[],
  defaultSelectedIds?: string[],
  persistSelections?: boolean,
): UseCheckboxTable => {
  const [selectedIds, setSelectedIds] = React.useState<string[]>(defaultSelectedIds ?? []);

  return useCheckboxTableBase<string>(
    dataIds,
    selectedIds,
    setSelectedIds,
    React.useCallback((d) => d, []),
    { persistSelections },
  );
};

export default useCheckboxTable;

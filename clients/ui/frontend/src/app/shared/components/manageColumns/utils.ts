// TODO this util was copied from odh-dashboard temporarily and should be abstracted out into mod-arch-shared.

import { ManagedColumn } from './useManageColumns';

/**
 * Reorders columns while preserving the positions of non-matching columns.
 *
 * When a user reorders columns during a search filter, only the filtered (matching)
 * columns should change their relative order. Non-matching columns stay in their
 * original positions.
 *
 * Example:
 * - columns: [A, B, C, D, E]
 * - matchingIds: [B, D] (search matches B and D)
 * - reorderedIds: [D, B] (user dragged D before B)
 * - result: [A, D, C, B, E] (B and D swap positions, others unchanged)
 *
 * @param columns - The current column array
 * @param matchingIds - Set of column IDs that match the current search filter
 * @param reorderedIds - The new order of matching column IDs after drag-drop
 * @returns The reordered column array
 */
export const reorderColumns = (
  columns: ManagedColumn[],
  matchingIds: Set<string>,
  reorderedIds: string[],
): ManagedColumn[] => {
  const columnMap = new Map(columns.map((col) => [col.id, col]));

  // Find indices where matching columns are located
  const matchingIndices: number[] = [];
  columns.forEach((col, index) => {
    if (matchingIds.has(col.id)) {
      matchingIndices.push(index);
    }
  });

  // Start with a copy of columns (non-matching columns stay in place)
  const result = [...columns];

  // Place reordered columns at the original matching indices
  reorderedIds.forEach((id, i) => {
    const col = columnMap.get(id);
    if (col && i < matchingIndices.length) {
      result[matchingIndices[i]] = col;
    }
  });

  return result;
};

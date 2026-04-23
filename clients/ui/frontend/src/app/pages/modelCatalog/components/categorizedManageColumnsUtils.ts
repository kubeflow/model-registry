import type { ManagedColumn } from 'mod-arch-shared';

/**
 * Reorders columns while preserving the positions of non-matching columns.
 *
 * When a user reorders columns via drag-drop within a category section, only the
 * columns in that section should change their relative order. Columns outside the
 * section stay in their original positions.
 *
 * Example:
 * - columns: [A, B, C, D, E]
 * - matchingIds: {B, D} (columns in the dragged section)
 * - reorderedIds: [D, B] (user dragged D before B)
 * - result: [A, D, C, B, E] (B and D swap positions, others unchanged)
 */
export const reorderColumns = (
  columns: ManagedColumn[],
  matchingIds: Set<string>,
  reorderedIds: string[],
): ManagedColumn[] => {
  const columnMap = new Map(columns.map((col) => [col.id, col]));

  const matchingIndices: number[] = [];
  columns.forEach((col, index) => {
    if (matchingIds.has(col.id)) {
      matchingIndices.push(index);
    }
  });

  const result = [...columns];

  reorderedIds.forEach((id, i) => {
    const col = columnMap.get(id);
    if (col && i < matchingIndices.length) {
      result[matchingIndices[i]] = col;
    }
  });

  return result;
};

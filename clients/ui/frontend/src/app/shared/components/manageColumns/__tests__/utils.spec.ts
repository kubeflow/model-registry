import { ManagedColumn } from '~/app/shared/components/manageColumns/useManageColumns';
import { reorderColumns } from '~/app/shared/components/manageColumns/utils';

const createColumn = (id: string, isVisible = true): ManagedColumn => ({
  id,
  label: `Column ${id}`,
  isVisible,
});

describe('reorderColumns', () => {
  describe('when all columns match the search', () => {
    it('should reorder all columns according to the new order', () => {
      const columns = [createColumn('A'), createColumn('B'), createColumn('C')];
      const matchingIds = new Set(['A', 'B', 'C']);
      const reorderedIds = ['C', 'A', 'B'];

      const result = reorderColumns(columns, matchingIds, reorderedIds);

      expect(result.map((c) => c.id)).toEqual(['C', 'A', 'B']);
    });

    it('should handle reversing the order', () => {
      const columns = [createColumn('A'), createColumn('B'), createColumn('C')];
      const matchingIds = new Set(['A', 'B', 'C']);
      const reorderedIds = ['C', 'B', 'A'];

      const result = reorderColumns(columns, matchingIds, reorderedIds);

      expect(result.map((c) => c.id)).toEqual(['C', 'B', 'A']);
    });
  });

  describe('when only some columns match the search', () => {
    it('should swap matching columns while keeping non-matching in place', () => {
      // Example from the docstring: [A, B, C, D, E] with B,D matching, reorder to [D, B]
      const columns = [
        createColumn('A'),
        createColumn('B'),
        createColumn('C'),
        createColumn('D'),
        createColumn('E'),
      ];
      const matchingIds = new Set(['B', 'D']);
      const reorderedIds = ['D', 'B'];

      const result = reorderColumns(columns, matchingIds, reorderedIds);

      // B was at index 1, D was at index 3
      // After reorder: D goes to index 1, B goes to index 3
      expect(result.map((c) => c.id)).toEqual(['A', 'D', 'C', 'B', 'E']);
    });

    it('should handle matching columns at the start', () => {
      const columns = [createColumn('A'), createColumn('B'), createColumn('C'), createColumn('D')];
      const matchingIds = new Set(['A', 'B']);
      const reorderedIds = ['B', 'A'];

      const result = reorderColumns(columns, matchingIds, reorderedIds);

      expect(result.map((c) => c.id)).toEqual(['B', 'A', 'C', 'D']);
    });

    it('should handle matching columns at the end', () => {
      const columns = [createColumn('A'), createColumn('B'), createColumn('C'), createColumn('D')];
      const matchingIds = new Set(['C', 'D']);
      const reorderedIds = ['D', 'C'];

      const result = reorderColumns(columns, matchingIds, reorderedIds);

      expect(result.map((c) => c.id)).toEqual(['A', 'B', 'D', 'C']);
    });

    it('should handle a single matching column (no actual reorder)', () => {
      const columns = [createColumn('A'), createColumn('B'), createColumn('C')];
      const matchingIds = new Set(['B']);
      const reorderedIds = ['B'];

      const result = reorderColumns(columns, matchingIds, reorderedIds);

      expect(result.map((c) => c.id)).toEqual(['A', 'B', 'C']);
    });

    it('should handle three matching columns among five', () => {
      const columns = [
        createColumn('A'),
        createColumn('B'),
        createColumn('C'),
        createColumn('D'),
        createColumn('E'),
      ];
      const matchingIds = new Set(['A', 'C', 'E']);
      const reorderedIds = ['E', 'C', 'A'];

      const result = reorderColumns(columns, matchingIds, reorderedIds);

      // A at 0, C at 2, E at 4 -> E goes to 0, C stays at 2, A goes to 4
      expect(result.map((c) => c.id)).toEqual(['E', 'B', 'C', 'D', 'A']);
    });
  });

  describe('when no columns match the search', () => {
    it('should return columns unchanged', () => {
      const columns = [createColumn('A'), createColumn('B'), createColumn('C')];
      const matchingIds = new Set<string>();
      const reorderedIds: string[] = [];

      const result = reorderColumns(columns, matchingIds, reorderedIds);

      expect(result.map((c) => c.id)).toEqual(['A', 'B', 'C']);
    });
  });

  describe('edge cases', () => {
    it('should handle empty columns array', () => {
      const columns: ManagedColumn[] = [];
      const matchingIds = new Set<string>();
      const reorderedIds: string[] = [];

      const result = reorderColumns(columns, matchingIds, reorderedIds);

      expect(result).toEqual([]);
    });

    it('should handle reorderedIds with unknown column IDs gracefully', () => {
      const columns = [createColumn('A'), createColumn('B')];
      const matchingIds = new Set(['A', 'B']);
      const reorderedIds = ['B', 'A', 'UNKNOWN'];

      const result = reorderColumns(columns, matchingIds, reorderedIds);

      // Unknown ID is ignored, only valid columns are placed
      expect(result.map((c) => c.id)).toEqual(['B', 'A']);
    });

    it('should preserve column properties during reorder', () => {
      const columns = [
        { id: 'A', label: 'Alpha', isVisible: true },
        { id: 'B', label: 'Beta', isVisible: false },
        { id: 'C', label: 'Gamma', isVisible: true },
      ];
      const matchingIds = new Set(['A', 'C']);
      const reorderedIds = ['C', 'A'];

      const result = reorderColumns(columns, matchingIds, reorderedIds);

      expect(result).toEqual([
        { id: 'C', label: 'Gamma', isVisible: true },
        { id: 'B', label: 'Beta', isVisible: false },
        { id: 'A', label: 'Alpha', isVisible: true },
      ]);
    });
  });
});

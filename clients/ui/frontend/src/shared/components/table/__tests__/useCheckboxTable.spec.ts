import { act } from 'react';
import { testHook } from '~/__tests__/unit/testUtils/hooks';
import { useCheckboxTable } from '~/shared/components/table';

describe('useCheckboxTable', () => {
  it('should select/unselect all', () => {
    const renderResult = testHook(useCheckboxTable)(['a', 'b', 'c']);

    act(() => {
      renderResult.result.current.tableProps.selectAll.onSelect(true);
    });
    expect(renderResult.result.current.selections).toStrictEqual(['a', 'b', 'c']);
    expect(renderResult.result.current.tableProps.selectAll.selected).toBe(true);

    act(() => {
      renderResult.result.current.tableProps.selectAll.onSelect(false);
    });
    expect(renderResult.result.current.selections).toStrictEqual([]);
    expect(renderResult.result.current.tableProps.selectAll.selected).toBe(false);
  });

  it('should select/unselect ids', () => {
    const renderResult = testHook(useCheckboxTable)(['a', 'b', 'c']);

    act(() => {
      renderResult.result.current.toggleSelection('a');
      renderResult.result.current.toggleSelection('b');
    });
    expect(renderResult.result.current.selections).toStrictEqual(['a', 'b']);
    expect(renderResult.result.current.isSelected('a')).toBe(true);
    expect(renderResult.result.current.isSelected('b')).toBe(true);
    expect(renderResult.result.current.isSelected('c')).toBe(false);

    act(() => {
      renderResult.result.current.toggleSelection('a');
    });
    expect(renderResult.result.current.selections).toStrictEqual(['b']);
    expect(renderResult.result.current.isSelected('a')).toBe(false);
    expect(renderResult.result.current.isSelected('b')).toBe(true);
    expect(renderResult.result.current.isSelected('c')).toBe(false);
  });

  it('should remove selected ids that no longer exist', () => {
    const renderResult = testHook(useCheckboxTable)(['a', 'b', 'c']);

    act(() => {
      renderResult.result.current.tableProps.selectAll.onSelect(true);
    });

    renderResult.rerender(['c']);

    expect(renderResult.result.current.selections).toStrictEqual(['c']);
  });
});

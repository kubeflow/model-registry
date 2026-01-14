import { act } from 'react';
import { useBrowserStorage } from 'mod-arch-core';
import { SortableData, checkboxTableColumn, kebabTableColumn } from 'mod-arch-shared';
import { testHook } from '~/__tests__/unit/testUtils/hooks';
import { useManageColumns } from '~/app/shared/components/manageColumns/useManageColumns';

jest.mock('mod-arch-core', () => ({
  useBrowserStorage: jest.fn(),
}));

const useBrowserStorageMock = jest.mocked(useBrowserStorage);

type TestData = { id: string; name: string };

const createMockColumns = (): SortableData<TestData>[] => [
  checkboxTableColumn(),
  { label: 'Column A', field: 'col_a', sortable: false },
  { label: 'Column B', field: 'col_b', sortable: false },
  { label: 'Column C', field: 'col_c', sortable: false },
  kebabTableColumn(),
];

describe('useManageColumns', () => {
  let mockSetStoredIds: jest.Mock;

  beforeEach(() => {
    mockSetStoredIds = jest.fn(() => true);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('should return default visible columns when no localStorage value exists', () => {
    useBrowserStorageMock.mockReturnValue([['col_a', 'col_b'], mockSetStoredIds]);

    const renderResult = testHook(useManageColumns<TestData>)({
      allColumns: createMockColumns(),
      storageKey: 'test-columns',
      defaultVisibleColumnIds: ['col_a', 'col_b'],
    });

    expect(renderResult.result.current.visibleColumnIds).toEqual(['col_a', 'col_b']);
  });

  it('should use first 2 columns as default when defaultVisibleColumnIds not provided', () => {
    useBrowserStorageMock.mockReturnValue([['col_a', 'col_b'], mockSetStoredIds]);

    testHook(useManageColumns<TestData>)({
      allColumns: createMockColumns(),
      storageKey: 'test-columns',
    });

    // The hook should have called useBrowserStorage with defaults
    expect(useBrowserStorageMock).toHaveBeenCalledWith('test-columns', ['col_a', 'col_b'], true);
  });

  it('should filter out non-manageable columns (checkbox, kebab)', () => {
    useBrowserStorageMock.mockReturnValue([['col_a'], mockSetStoredIds]);

    const renderResult = testHook(useManageColumns<TestData>)({
      allColumns: createMockColumns(),
      storageKey: 'test-columns',
    });

    // managedColumns should not include checkbox or kebab
    const managedIds = renderResult.result.current.managedColumns.map((c) => c.id);
    expect(managedIds).not.toContain('checkbox');
    expect(managedIds).not.toContain('kebab');
    expect(managedIds).toContain('col_a');
    expect(managedIds).toContain('col_b');
    expect(managedIds).toContain('col_c');
  });

  it('should preserve column order from storage', () => {
    // Storage has columns in order: c, a, b
    useBrowserStorageMock.mockReturnValue([['col_c', 'col_a', 'col_b'], mockSetStoredIds]);

    const renderResult = testHook(useManageColumns<TestData>)({
      allColumns: createMockColumns(),
      storageKey: 'test-columns',
    });

    const managedIds = renderResult.result.current.managedColumns.map((c) => c.id);
    expect(managedIds).toEqual(['col_c', 'col_a', 'col_b']);
  });

  it('should add new columns to end of list', () => {
    // Storage only has col_a, but allColumns has col_a, col_b, col_c
    useBrowserStorageMock.mockReturnValue([['col_a'], mockSetStoredIds]);

    const renderResult = testHook(useManageColumns<TestData>)({
      allColumns: createMockColumns(),
      storageKey: 'test-columns',
    });

    const managedIds = renderResult.result.current.managedColumns.map((c) => c.id);
    // col_a should be first (from storage), then col_b and col_c appended
    expect(managedIds[0]).toBe('col_a');
    expect(managedIds).toContain('col_b');
    expect(managedIds).toContain('col_c');
  });

  it('should remove deleted columns from stored list', () => {
    // Storage has col_a, col_deleted, col_b but col_deleted is not in allColumns
    useBrowserStorageMock.mockReturnValue([['col_a', 'col_deleted', 'col_b'], mockSetStoredIds]);

    const renderResult = testHook(useManageColumns<TestData>)({
      allColumns: createMockColumns(),
      storageKey: 'test-columns',
    });

    const managedIds = renderResult.result.current.managedColumns.map((c) => c.id);
    expect(managedIds).not.toContain('col_deleted');
    expect(managedIds).toContain('col_a');
    expect(managedIds).toContain('col_b');
  });

  it('should mark columns as visible based on stored IDs', () => {
    useBrowserStorageMock.mockReturnValue([['col_a', 'col_c'], mockSetStoredIds]);

    const renderResult = testHook(useManageColumns<TestData>)({
      allColumns: createMockColumns(),
      storageKey: 'test-columns',
    });

    const colA = renderResult.result.current.managedColumns.find((c) => c.id === 'col_a');
    const colB = renderResult.result.current.managedColumns.find((c) => c.id === 'col_b');
    const colC = renderResult.result.current.managedColumns.find((c) => c.id === 'col_c');

    expect(colA?.isVisible).toBe(true);
    expect(colB?.isVisible).toBe(false);
    expect(colC?.isVisible).toBe(true);
  });

  it('should call setStoredVisibleIds when setVisibleColumnIds is called', () => {
    useBrowserStorageMock.mockReturnValue([['col_a'], mockSetStoredIds]);

    const renderResult = testHook(useManageColumns<TestData>)({
      allColumns: createMockColumns(),
      storageKey: 'test-columns',
    });

    act(() => {
      renderResult.result.current.setVisibleColumnIds(['col_b', 'col_c']);
    });

    expect(mockSetStoredIds).toHaveBeenCalledWith(['col_b', 'col_c']);
  });

  it('should build visibleColumns with checkbox first, visible columns in order, kebab last', () => {
    useBrowserStorageMock.mockReturnValue([['col_b', 'col_a'], mockSetStoredIds]);

    const renderResult = testHook(useManageColumns<TestData>)({
      allColumns: createMockColumns(),
      storageKey: 'test-columns',
    });

    const visibleFields = renderResult.result.current.visibleColumns.map((c) => c.field);

    // Checkbox should be first
    expect(visibleFields[0]).toBe('checkbox');
    // Then visible manageable columns in stored order
    expect(visibleFields[1]).toBe('col_b');
    expect(visibleFields[2]).toBe('col_a');
    // Kebab should be last
    expect(visibleFields[visibleFields.length - 1]).toBe('kebab');
  });

  it('should return empty visibleColumnIds when storage is empty', () => {
    useBrowserStorageMock.mockReturnValue([[], mockSetStoredIds]);

    const renderResult = testHook(useManageColumns<TestData>)({
      allColumns: createMockColumns(),
      storageKey: 'test-columns',
    });

    expect(renderResult.result.current.visibleColumnIds).toEqual([]);
  });
});

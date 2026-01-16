import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import { ManageColumnsModal } from '~/app/shared/components/manageColumns/ManageColumnsModal';
import { ManagedColumn } from '~/app/shared/components/manageColumns/useManageColumns';

// Mock DragDropSort to avoid portal issues in tests
jest.mock('@patternfly/react-drag-drop', () => ({
  DragDropSort: ({ items }: { items: { id: string; content: React.ReactNode }[] }) => (
    <div data-testid="mock-drag-drop-sort">
      {items.map((item) => (
        <div key={item.id} data-testid={`draggable-${item.id}`}>
          {item.content}
        </div>
      ))}
    </div>
  ),
}));

const createMockColumns = (count = 3): ManagedColumn[] =>
  Array.from({ length: count }, (_, i) => ({
    id: `col_${i + 1}`,
    label: `Column ${i + 1}`,
    isVisible: i < 2, // First 2 are visible by default
  }));

const createMockManageColumnsResult = (
  columns: ManagedColumn[],
  setVisibleColumnIds: jest.Mock,
  closeModal: jest.Mock,
  isModalOpen = true,
  defaultVisibleColumnIds?: string[],
) => ({
  managedColumns: columns,
  setVisibleColumnIds,
  defaultVisibleColumnIds,
  isModalOpen,
  closeModal,
});

describe('ManageColumnsModal', () => {
  const mockCloseModal = jest.fn();
  const mockSetVisibleColumnIds = jest.fn();

  beforeEach(() => {
    mockCloseModal.mockClear();
    mockSetVisibleColumnIds.mockClear();
  });

  it('should render with correct title', () => {
    render(
      <ManageColumnsModal
        manageColumnsResult={createMockManageColumnsResult(
          createMockColumns(),
          mockSetVisibleColumnIds,
          mockCloseModal,
        )}
        title="Custom Title"
      />,
    );

    expect(screen.getByText('Custom Title')).toBeInTheDocument();
  });

  it('should render with default title when not provided', () => {
    render(
      <ManageColumnsModal
        manageColumnsResult={createMockManageColumnsResult(
          createMockColumns(),
          mockSetVisibleColumnIds,
          mockCloseModal,
        )}
      />,
    );

    expect(screen.getByText('Customize columns')).toBeInTheDocument();
  });

  it('should render description when provided', () => {
    render(
      <ManageColumnsModal
        manageColumnsResult={createMockManageColumnsResult(
          createMockColumns(),
          mockSetVisibleColumnIds,
          mockCloseModal,
        )}
        description="Select columns to display"
      />,
    );

    expect(screen.getByText('Select columns to display')).toBeInTheDocument();
  });

  it('should display all columns with labels', () => {
    render(
      <ManageColumnsModal
        manageColumnsResult={createMockManageColumnsResult(
          createMockColumns(),
          mockSetVisibleColumnIds,
          mockCloseModal,
        )}
      />,
    );

    expect(screen.getByText('Column 1')).toBeInTheDocument();
    expect(screen.getByText('Column 2')).toBeInTheDocument();
    expect(screen.getByText('Column 3')).toBeInTheDocument();
  });

  it('should show correct selected count', () => {
    const columns = createMockColumns(5);
    // 2 visible out of 5
    render(
      <ManageColumnsModal
        manageColumnsResult={createMockManageColumnsResult(
          columns,
          mockSetVisibleColumnIds,
          mockCloseModal,
        )}
      />,
    );

    expect(screen.getByText('2 / total 5 selected')).toBeInTheDocument();
  });

  it('should call closeModal when Cancel button is clicked', () => {
    render(
      <ManageColumnsModal
        manageColumnsResult={createMockManageColumnsResult(
          createMockColumns(),
          mockSetVisibleColumnIds,
          mockCloseModal,
        )}
        dataTestId="test-modal"
      />,
    );

    fireEvent.click(screen.getByTestId('test-modal-cancel-button'));

    expect(mockCloseModal).toHaveBeenCalled();
    expect(mockSetVisibleColumnIds).not.toHaveBeenCalled();
  });

  it('should call setVisibleColumnIds with visible column IDs when Update button is clicked', () => {
    const columns: ManagedColumn[] = [
      { id: 'col_1', label: 'Column 1', isVisible: true },
      { id: 'col_2', label: 'Column 2', isVisible: false },
      { id: 'col_3', label: 'Column 3', isVisible: true },
    ];

    render(
      <ManageColumnsModal
        manageColumnsResult={createMockManageColumnsResult(
          columns,
          mockSetVisibleColumnIds,
          mockCloseModal,
        )}
        dataTestId="test-modal"
      />,
    );

    fireEvent.click(screen.getByTestId('test-modal-update-button'));

    expect(mockSetVisibleColumnIds).toHaveBeenCalledWith(['col_1', 'col_3']);
    expect(mockCloseModal).toHaveBeenCalled();
  });

  it('should render search input with custom placeholder', () => {
    render(
      <ManageColumnsModal
        manageColumnsResult={createMockManageColumnsResult(
          createMockColumns(),
          mockSetVisibleColumnIds,
          mockCloseModal,
        )}
        searchPlaceholder="Find columns..."
      />,
    );

    expect(screen.getByPlaceholderText('Find columns...')).toBeInTheDocument();
  });

  it('should filter columns when searching', () => {
    const columns: ManagedColumn[] = [
      { id: 'name', label: 'Name', isVisible: true },
      { id: 'status', label: 'Status', isVisible: true },
      { id: 'created', label: 'Created Date', isVisible: false },
    ];

    render(
      <ManageColumnsModal
        manageColumnsResult={createMockManageColumnsResult(
          columns,
          mockSetVisibleColumnIds,
          mockCloseModal,
        )}
        dataTestId="test-modal"
      />,
    );

    const searchInput = screen.getByPlaceholderText('Filter by column name');
    fireEvent.change(searchInput, { target: { value: 'status' } });

    // Only Status should be visible after filtering
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.queryByText('Name')).not.toBeInTheDocument();
    expect(screen.queryByText('Created Date')).not.toBeInTheDocument();
  });

  it('should not render when isModalOpen is false', () => {
    render(
      <ManageColumnsModal
        manageColumnsResult={createMockManageColumnsResult(
          createMockColumns(),
          mockSetVisibleColumnIds,
          mockCloseModal,
          false, // isModalOpen = false
        )}
        title="Should Not Appear"
      />,
    );

    expect(screen.queryByText('Should Not Appear')).not.toBeInTheDocument();
  });

  it('should toggle column visibility when checkbox is clicked', () => {
    const columns: ManagedColumn[] = [
      { id: 'col_1', label: 'Column 1', isVisible: true },
      { id: 'col_2', label: 'Column 2', isVisible: false },
    ];

    render(
      <ManageColumnsModal
        manageColumnsResult={createMockManageColumnsResult(
          columns,
          mockSetVisibleColumnIds,
          mockCloseModal,
        )}
        dataTestId="test-modal"
      />,
    );

    // Find checkboxes by their id attribute
    const checkboxes = screen.getAllByRole('checkbox');
    const col2CheckboxElement = checkboxes.find(
      (cb) => cb.getAttribute('id') === 'col_2',
    ) as HTMLInputElement;

    expect(col2CheckboxElement).not.toBeChecked();
    fireEvent.click(col2CheckboxElement);

    // Now click Update and verify col_2 is included
    fireEvent.click(screen.getByTestId('test-modal-update-button'));

    expect(mockSetVisibleColumnIds).toHaveBeenCalledWith(['col_1', 'col_2']);
  });

  it('should disable unchecked columns when maxSelections is reached', () => {
    const columns: ManagedColumn[] = [
      { id: 'col_1', label: 'Column 1', isVisible: true },
      { id: 'col_2', label: 'Column 2', isVisible: true },
      { id: 'col_3', label: 'Column 3', isVisible: false },
    ];

    render(
      <ManageColumnsModal
        manageColumnsResult={createMockManageColumnsResult(
          columns,
          mockSetVisibleColumnIds,
          mockCloseModal,
        )}
        maxSelections={2}
      />,
    );

    // col_3 should be disabled since we're at max
    const checkboxes = screen.getAllByRole('checkbox');
    const col3Checkbox = checkboxes.find(
      (cb) => cb.getAttribute('id') === 'col_3',
    ) as HTMLInputElement;

    expect(col3Checkbox).toBeDisabled();
  });

  it('should use custom dataTestId prefix', () => {
    render(
      <ManageColumnsModal
        manageColumnsResult={createMockManageColumnsResult(
          createMockColumns(),
          mockSetVisibleColumnIds,
          mockCloseModal,
        )}
        dataTestId="my-custom-modal"
      />,
    );

    expect(screen.getByTestId('my-custom-modal')).toBeInTheDocument();
    expect(screen.getByTestId('my-custom-modal-update-button')).toBeInTheDocument();
    expect(screen.getByTestId('my-custom-modal-cancel-button')).toBeInTheDocument();
  });

  it('should render restore defaults button when defaultVisibleColumnIds is provided', () => {
    render(
      <ManageColumnsModal
        manageColumnsResult={createMockManageColumnsResult(
          createMockColumns(),
          mockSetVisibleColumnIds,
          mockCloseModal,
          true,
          ['col_1', 'col_2'],
        )}
        dataTestId="test-modal"
      />,
    );

    expect(screen.getByTestId('test-modal-restore-defaults')).toBeInTheDocument();
    expect(screen.getByText('Restore default columns')).toBeInTheDocument();
  });

  it('should not render restore defaults button when defaultVisibleColumnIds is not provided', () => {
    render(
      <ManageColumnsModal
        manageColumnsResult={createMockManageColumnsResult(
          createMockColumns(),
          mockSetVisibleColumnIds,
          mockCloseModal,
        )}
        dataTestId="test-modal"
      />,
    );

    expect(screen.queryByTestId('test-modal-restore-defaults')).not.toBeInTheDocument();
    expect(screen.queryByText('Restore default columns')).not.toBeInTheDocument();
  });

  it('should restore default columns when restore defaults button is clicked', () => {
    const columns: ManagedColumn[] = [
      { id: 'col_1', label: 'Column 1', isVisible: false },
      { id: 'col_2', label: 'Column 2', isVisible: false },
      { id: 'col_3', label: 'Column 3', isVisible: true },
    ];

    render(
      <ManageColumnsModal
        manageColumnsResult={createMockManageColumnsResult(
          columns,
          mockSetVisibleColumnIds,
          mockCloseModal,
          true,
          ['col_1', 'col_2'],
        )}
        dataTestId="test-modal"
      />,
    );

    // Click restore defaults
    fireEvent.click(screen.getByTestId('test-modal-restore-defaults'));

    // Now click Update and verify only default columns are visible
    fireEvent.click(screen.getByTestId('test-modal-update-button'));

    expect(mockSetVisibleColumnIds).toHaveBeenCalledWith(['col_1', 'col_2']);
  });
});

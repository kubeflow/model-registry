import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import { ManageColumnsModal } from '#~/components/table/manageColumns/ManageColumnsModal';
import { ManagedColumn } from '#~/components/table/manageColumns/types';

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

describe('ManageColumnsModal', () => {
  const mockOnClose = jest.fn();
  const mockOnUpdate = jest.fn();

  beforeEach(() => {
    mockOnClose.mockClear();
    mockOnUpdate.mockClear();
  });

  it('should render with correct title', () => {
    render(
      <ManageColumnsModal
        isOpen
        columns={createMockColumns()}
        onClose={mockOnClose}
        onUpdate={mockOnUpdate}
        title="Custom Title"
      />,
    );

    expect(screen.getByText('Custom Title')).toBeInTheDocument();
  });

  it('should render with default title when not provided', () => {
    render(
      <ManageColumnsModal
        isOpen
        columns={createMockColumns()}
        onClose={mockOnClose}
        onUpdate={mockOnUpdate}
      />,
    );

    expect(screen.getByText('Manage columns')).toBeInTheDocument();
  });

  it('should render description when provided', () => {
    render(
      <ManageColumnsModal
        isOpen
        columns={createMockColumns()}
        onClose={mockOnClose}
        onUpdate={mockOnUpdate}
        description="Select columns to display"
      />,
    );

    expect(screen.getByText('Select columns to display')).toBeInTheDocument();
  });

  it('should display all columns with labels', () => {
    render(
      <ManageColumnsModal
        isOpen
        columns={createMockColumns()}
        onClose={mockOnClose}
        onUpdate={mockOnUpdate}
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
      <ManageColumnsModal isOpen columns={columns} onClose={mockOnClose} onUpdate={mockOnUpdate} />,
    );

    expect(screen.getByText('2 / total 5 selected')).toBeInTheDocument();
  });

  it('should call onClose when Cancel button is clicked', () => {
    render(
      <ManageColumnsModal
        isOpen
        columns={createMockColumns()}
        onClose={mockOnClose}
        onUpdate={mockOnUpdate}
        dataTestId="test-modal"
      />,
    );

    fireEvent.click(screen.getByTestId('test-modal-cancel-button'));

    expect(mockOnClose).toHaveBeenCalled();
    expect(mockOnUpdate).not.toHaveBeenCalled();
  });

  it('should call onUpdate with visible column IDs when Update button is clicked', () => {
    const columns: ManagedColumn[] = [
      { id: 'col_1', label: 'Column 1', isVisible: true },
      { id: 'col_2', label: 'Column 2', isVisible: false },
      { id: 'col_3', label: 'Column 3', isVisible: true },
    ];

    render(
      <ManageColumnsModal
        isOpen
        columns={columns}
        onClose={mockOnClose}
        onUpdate={mockOnUpdate}
        dataTestId="test-modal"
      />,
    );

    fireEvent.click(screen.getByTestId('test-modal-update-button'));

    expect(mockOnUpdate).toHaveBeenCalledWith(['col_1', 'col_3']);
    expect(mockOnClose).toHaveBeenCalled();
  });

  it('should render search input with custom placeholder', () => {
    render(
      <ManageColumnsModal
        isOpen
        columns={createMockColumns()}
        onClose={mockOnClose}
        onUpdate={mockOnUpdate}
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
        isOpen
        columns={columns}
        onClose={mockOnClose}
        onUpdate={mockOnUpdate}
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

  it('should not render when isOpen is false', () => {
    render(
      <ManageColumnsModal
        isOpen={false}
        columns={createMockColumns()}
        onClose={mockOnClose}
        onUpdate={mockOnUpdate}
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
        isOpen
        columns={columns}
        onClose={mockOnClose}
        onUpdate={mockOnUpdate}
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

    expect(mockOnUpdate).toHaveBeenCalledWith(['col_1', 'col_2']);
  });

  it('should disable unchecked columns when maxSelections is reached', () => {
    const columns: ManagedColumn[] = [
      { id: 'col_1', label: 'Column 1', isVisible: true },
      { id: 'col_2', label: 'Column 2', isVisible: true },
      { id: 'col_3', label: 'Column 3', isVisible: false },
    ];

    render(
      <ManageColumnsModal
        isOpen
        columns={columns}
        onClose={mockOnClose}
        onUpdate={mockOnUpdate}
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
        isOpen
        columns={createMockColumns()}
        onClose={mockOnClose}
        onUpdate={mockOnUpdate}
        dataTestId="my-custom-modal"
      />,
    );

    expect(screen.getByTestId('my-custom-modal')).toBeInTheDocument();
    expect(screen.getByTestId('my-custom-modal-update-button')).toBeInTheDocument();
    expect(screen.getByTestId('my-custom-modal-cancel-button')).toBeInTheDocument();
  });
});

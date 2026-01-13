import React from 'react';
import { Checkbox, Flex, FlexItem, Label, Stack, StackItem, Tooltip } from '@patternfly/react-core';
import {
  DragDropSort,
  DragDropSortDragEndEvent,
  DraggableObject,
} from '@patternfly/react-drag-drop';
import ContentModal, { ButtonAction } from '#~/components/modals/ContentModal';
import { ManageColumnSearchInput } from './ManageColumnSearchInput';
import { ManagedColumn } from './useManageColumns';
import { reorderColumns } from './utils';

/**
 * Configuration for the ManageColumnsModal
 */
export interface ManageColumnsModalProps {
  /** Whether the modal is open */
  isOpen: boolean;
  /** Callback when the modal closes (cancel or update) */
  onClose: () => void;
  /** Callback when columns are updated - receives the new ordered list of visible column IDs */
  onUpdate: (visibleColumnIds: string[]) => void;
  /** All available columns that can be managed */
  columns: ManagedColumn[];
  /** Modal title - defaults to "Manage columns" */
  title?: string;
  /** Description text shown above the column list */
  description?: string;
  /** Maximum number of columns that can be selected, undefined = unlimited */
  maxSelections?: number;
  /** Tooltip text when max is reached */
  maxSelectionsTooltip?: string;
  /** Placeholder for the search input */
  searchPlaceholder?: string;
  /** Test ID prefix for data-testid attributes */
  dataTestId?: string;
}

export const ManageColumnsModal: React.FC<ManageColumnsModalProps> = ({
  isOpen,
  onClose,
  onUpdate,
  columns: initialColumns,
  title = 'Manage columns',
  description,
  maxSelections,
  maxSelectionsTooltip = 'Maximum columns selected.',
  searchPlaceholder = 'Filter by column name',
  dataTestId = 'manage-columns-modal',
}) => {
  const [columns, setColumns] = React.useState<ManagedColumn[]>(initialColumns);
  const [searchValue, setSearchValue] = React.useState('');

  // Derive filtered columns from search
  const columnsMatchingSearch = React.useMemo(
    () =>
      searchValue
        ? columns.filter((col) => col.label.toLowerCase().includes(searchValue.toLowerCase()))
        : columns,
    [columns, searchValue],
  );

  // Reset state when modal opens with new columns
  React.useEffect(() => {
    if (isOpen) {
      setColumns(initialColumns);
      setSearchValue('');
    }
  }, [isOpen, initialColumns]);

  const selectedCount = columns.filter((col) => col.isVisible).length;
  const isMaxReached = maxSelections !== undefined && selectedCount >= maxSelections;

  const handleUpdate = React.useCallback(() => {
    const visibleColumnIds = columns.filter((col) => col.isVisible).map((col) => col.id);
    onUpdate(visibleColumnIds);
    onClose();
  }, [columns, onUpdate, onClose]);

  const handleSearch = React.useCallback((value: string) => {
    setSearchValue(value);
  }, []);

  const handleToggleColumn = React.useCallback((columnId: string, isChecked: boolean) => {
    setColumns((prev) =>
      prev.map((col) => (col.id === columnId ? { ...col, isVisible: isChecked } : col)),
    );
  }, []);

  const handleDrop = React.useCallback(
    (_: DragDropSortDragEndEvent, newItems: DraggableObject[]) => {
      const reorderedIds = newItems.map((item) => String(item.id));
      const matchingIds = new Set(columnsMatchingSearch.map((c) => c.id));
      setColumns((prev) => reorderColumns(prev, matchingIds, reorderedIds));
    },
    [columnsMatchingSearch],
  );

  if (!isOpen) {
    return null;
  }

  const buttonActions: ButtonAction[] = [
    {
      label: 'Update',
      onClick: handleUpdate,
      variant: 'primary',
      dataTestId: `${dataTestId}-update-button`,
    },
    {
      label: 'Cancel',
      onClick: onClose,
      variant: 'link',
      dataTestId: `${dataTestId}-cancel-button`,
    },
  ];

  const descriptionContent = (
    <Stack hasGutter className="pf-v6-u-pb-md">
      {description && <StackItem className="pf-v6-u-mt-sm">{description}</StackItem>}
      <StackItem>
        <ManageColumnSearchInput
          value={searchValue}
          placeholder={searchPlaceholder}
          onSearch={handleSearch}
          dataTestId={`${dataTestId}-search`}
        />
      </StackItem>
      <StackItem>
        <Label>
          {selectedCount} / total {columns.length} selected
        </Label>
      </StackItem>
    </Stack>
  );

  const draggableItems: DraggableObject[] = columnsMatchingSearch.map((col) => {
    const currentCol = columns.find((c) => c.id === col.id) ?? col;
    const isDisabled = isMaxReached && !currentCol.isVisible;

    const checkbox = (
      <Checkbox
        id={col.id}
        isChecked={currentCol.isVisible}
        isDisabled={isDisabled}
        onChange={(_, checked) => handleToggleColumn(col.id, checked)}
      />
    );

    return {
      id: col.id,
      content: (
        <div className="pf-v6-u-display-inline-block">
          <Flex
            alignItems={{ default: 'alignItemsCenter' }}
            flexWrap={{ default: 'nowrap' }}
            spaceItems={{ default: 'spaceItemsSm' }}
          >
            <FlexItem>
              {isDisabled ? <Tooltip content={maxSelectionsTooltip}>{checkbox}</Tooltip> : checkbox}
            </FlexItem>
            <FlexItem>{col.label}</FlexItem>
          </Flex>
        </div>
      ),
      props: { checked: currentCol.isVisible },
    };
  });

  return (
    <ContentModal
      onClose={onClose}
      title={title}
      description={descriptionContent}
      contents={<DragDropSort items={draggableItems} variant="default" onDrop={handleDrop} />}
      buttonActions={buttonActions}
      dataTestId={dataTestId}
      variant="small"
      bodyLabel="Column names"
    />
  );
};

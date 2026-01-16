// TODO this component was copied from odh-dashboard temporarily and should be abstracted out into mod-arch-shared.

import React from 'react';
import {
  Button,
  Checkbox,
  EmptyState,
  EmptyStateBody,
  EmptyStateVariant,
  Flex,
  FlexItem,
  Label,
  Stack,
  StackItem,
  Tooltip,
} from '@patternfly/react-core';
import { SearchIcon } from '@patternfly/react-icons';
import {
  DragDropSort,
  DragDropSortDragEndEvent,
  DraggableObject,
} from '@patternfly/react-drag-drop';
import ContentModal, { ButtonAction } from '~/app/shared/components/modals/ContentModal';
import { ManageColumnSearchInput } from './ManageColumnSearchInput';
import { ManagedColumn, UseManageColumnsResult } from './useManageColumns';
import { reorderColumns } from './utils';

/**
 * Configuration for the ManageColumnsModal
 */
export interface ManageColumnsModalProps {
  /** Result from useManageColumns hook - provides columns, callbacks, modal state, and defaults */
  manageColumnsResult: Pick<
    UseManageColumnsResult<unknown>,
    | 'managedColumns'
    | 'setVisibleColumnIds'
    | 'defaultVisibleColumnIds'
    | 'isModalOpen'
    | 'closeModal'
  >;
  /** Modal title - defaults to "Customize columns" */
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
  manageColumnsResult,
  title = 'Customize columns',
  description,
  maxSelections,
  maxSelectionsTooltip = 'Maximum columns selected.',
  searchPlaceholder = 'Filter by column name',
  dataTestId = 'manage-columns-modal',
}) => {
  const {
    managedColumns: initialColumns,
    setVisibleColumnIds,
    defaultVisibleColumnIds,
    isModalOpen,
    closeModal,
  } = manageColumnsResult;
  const [columns, setColumns] = React.useState<ManagedColumn[]>(initialColumns);
  const [searchValue, setSearchValue] = React.useState('');

  // Normalize whitespace (including NBSP) to regular spaces for search comparison
  const normalizeWhitespace = (str: string): string => str.replace(/\s+/g, ' ');

  // Derive filtered columns from search
  const columnsMatchingSearch = React.useMemo(
    () =>
      searchValue
        ? columns.filter((col) =>
            normalizeWhitespace(col.label.toLowerCase()).includes(
              normalizeWhitespace(searchValue.toLowerCase()),
            ),
          )
        : columns,
    [columns, searchValue],
  );

  // Reset state when modal opens with new columns
  React.useEffect(() => {
    if (isModalOpen) {
      setColumns(initialColumns);
      setSearchValue('');
    }
  }, [isModalOpen, initialColumns]);

  const selectedCount = columns.filter((col) => col.isVisible).length;
  const isMaxReached = maxSelections !== undefined && selectedCount >= maxSelections;

  const handleUpdate = React.useCallback(() => {
    const visibleColumnIds = columns.filter((col) => col.isVisible).map((col) => col.id);
    setVisibleColumnIds(visibleColumnIds);
    closeModal();
  }, [columns, setVisibleColumnIds, closeModal]);

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

  const handleRestoreDefaults = React.useCallback(() => {
    if (!defaultVisibleColumnIds) {
      return;
    }
    // Update visibility based on default column IDs
    // Also reorder to put default columns first in their original order
    setColumns((prev) => {
      const defaultSet = new Set(defaultVisibleColumnIds);
      const defaultColumns = defaultVisibleColumnIds
        .map((id) => prev.find((col) => col.id === id))
        .filter((col): col is ManagedColumn => col !== undefined)
        .map((col) => ({ ...col, isVisible: true }));
      const nonDefaultColumns = prev
        .filter((col) => !defaultSet.has(col.id))
        .map((col) => ({ ...col, isVisible: false }));
      return [...defaultColumns, ...nonDefaultColumns];
    });
  }, [defaultVisibleColumnIds]);

  if (!isModalOpen) {
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
      onClick: closeModal,
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
        <Flex
          justifyContent={{ default: 'justifyContentSpaceBetween' }}
          alignItems={{ default: 'alignItemsCenter' }}
        >
          <FlexItem>
            <Label>
              {selectedCount} / total {columns.length} selected
            </Label>
          </FlexItem>
          {defaultVisibleColumnIds && (
            <FlexItem>
              <Button
                variant="link"
                isInline
                onClick={handleRestoreDefaults}
                data-testid={`${dataTestId}-restore-defaults`}
              >
                Restore default columns
              </Button>
            </FlexItem>
          )}
        </Flex>
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
      onClose={closeModal}
      title={title}
      description={descriptionContent}
      contents={
        columnsMatchingSearch.length === 0 && searchValue ? (
          <EmptyState
            headingLevel="h4"
            icon={SearchIcon}
            titleText="No results found"
            variant={EmptyStateVariant.sm}
          >
            <EmptyStateBody>
              No columns match your search. Try adjusting your search terms.
            </EmptyStateBody>
          </EmptyState>
        ) : (
          <DragDropSort items={draggableItems} variant="default" onDrop={handleDrop} />
        )
      }
      buttonActions={buttonActions}
      dataTestId={dataTestId}
      variant="small"
      bodyLabel="Column names"
    />
  );
};
